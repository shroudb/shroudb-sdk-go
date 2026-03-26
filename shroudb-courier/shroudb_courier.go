// Package shroudb_courier provides a client for the ShroudbCourier Secure notification delivery pipeline.
//
// Auto-generated from shroudb-courier protocol spec. Do not edit.
//
// Connect using a ShroudbCourier URI:
//
//	client, err := shroudb_courier.Connect("shroudb-courier://localhost")
//	result, err := client.Issue(ctx, "my-keyspace", &shroudb_courier.IssueOptions{TTL: 3600})
//	fmt.Println(result.CredentialID, result.Token)
//	client.Close()
package shroudb_courier

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Client is a ShroudbCourier client backed by a connection pool.
type Client struct {
	pool *pool
	host string
	port int
	tls  bool
	auth string
}

// Connect creates a new Client from a ShroudbCourier URI.
//
// Supported URI formats:
//
//	shroudb-courier://localhost
//	shroudb-courier://localhost:6999
//	shroudb-courier+tls://prod.example.com
//	shroudb-courier://mytoken@localhost:6999
//	shroudb-courier://mytoken@localhost/sessions
//	shroudb-courier+tls://tok@host:6999/keys
func Connect(uri string, opts ...PoolConfig) (*Client, error) {
	cfg, err := parseURI(uri)
	if err != nil {
		return nil, err
	}
	var poolCfg PoolConfig
	if len(opts) > 0 {
		poolCfg = opts[0]
	}
	p := newPool(cfg.host, cfg.port, cfg.tls, cfg.authToken, poolCfg)
	// Verify connectivity by getting and returning one connection
	c, err := p.get()
	if err != nil {
		return nil, err
	}
	p.put(c)
	return &Client{pool: p, host: cfg.host, port: cfg.port, tls: cfg.tls, auth: cfg.authToken}, nil
}

// Close shuts down the client and all pooled connections.
func (c *Client) Close() {
	c.pool.close()
}

// Pipeline creates a new pipeline for batching commands.
func (c *Client) Pipeline() *Pipeline {
	return &Pipeline{pool: c.pool}
}

func (c *Client) exec(args ...string) (any, error) {
	conn, err := c.pool.get()
	if err != nil {
		return nil, err
	}
	result, err := conn.execute(args...)
	if err != nil {
		conn.close()
		return nil, err
	}
	c.pool.put(conn)
	return result, nil
}

func (c *Client) execMap(args ...string) (map[string]any, error) {
	result, err := c.exec(args...)
	if err != nil {
		return nil, err
	}
	m, ok := result.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("shroudb_courier: expected map response, got %T", result)
	}
	return m, nil
}

// Subscription represents an active streaming subscription.
type Subscription struct {
	conn   *connection
	events chan SubscriptionEvent
	errc   chan error
	done   chan struct{}
}

// Events returns a channel that receives subscription events.
func (s *Subscription) Events() <-chan SubscriptionEvent {
	return s.events
}

// Err returns a channel that receives the first read error (including clean shutdown).
func (s *Subscription) Err() <-chan error {
	return s.errc
}

// Close terminates the subscription and closes the underlying connection.
func (s *Subscription) Close() error {
	select {
	case <-s.done:
		return nil
	default:
		close(s.done)
		return s.conn.close()
	}
}

func (s *Subscription) readLoop() {
	defer close(s.events)
	for {
		select {
		case <-s.done:
			return
		default:
		}
		raw, err := s.conn.readFrame()
		if err != nil {
			select {
			case s.errc <- err:
			default:
			}
			return
		}
		arr, ok := raw.([]any)
		if !ok || len(arr) != 5 {
			continue
		}
		tag, _ := arr[0].(string)
		if tag != "event" {
			continue
		}
		evtType, _ := arr[1].(string)
		keyspace, _ := arr[2].(string)
		detail, _ := arr[3].(string)
		var ts int64
		switch v := arr[4].(type) {
		case int64:
			ts = v
		case string:
			ts, _ = strconv.ParseInt(v, 10, 64)
		}
		evt := SubscriptionEvent{
			EventType: evtType,
			Keyspace:  keyspace,
			Detail:    detail,
			Timestamp: ts,
		}
		select {
		case s.events <- evt:
		case <-s.done:
			return
		}
	}
}

// Subscribe opens a dedicated connection and subscribes to the given channel.
// The returned Subscription streams events until Close is called or an error occurs.
func (c *Client) Subscribe(channel string) (*Subscription, error) {
	conn, err := dial(c.host, c.port, c.tls)
	if err != nil {
		return nil, err
	}
	if c.auth != "" {
		if _, err := conn.execute("AUTH", c.auth); err != nil {
			conn.close()
			return nil, err
		}
	}
	resp, err := conn.execute("SUBSCRIBE", channel)
	if err != nil {
		conn.close()
		return nil, err
	}
	m, ok := resp.(map[string]any)
	if !ok {
		conn.close()
		return nil, fmt.Errorf("shroudb_courier: expected map response for SUBSCRIBE, got %T", resp)
	}
	if status, _ := m["status"].(string); status != "OK" {
		conn.close()
		return nil, fmt.Errorf("shroudb_courier: subscribe failed: %v", m)
	}
	sub := &Subscription{
		conn:   conn,
		events: make(chan SubscriptionEvent, 64),
		errc:   make(chan error, 1),
		done:   make(chan struct{}),
	}
	go sub.readLoop()
	return sub, nil
}

type uriConfig struct {
	host      string
	port      int
	tls       bool
	authToken string
	keyspace  string
}

func parseURI(uri string) (*uriConfig, error) {
	cfg := &uriConfig{port: defaultPort}

	switch {
	case strings.HasPrefix(uri, "shroudb-courier+tls://"):
		cfg.tls = true
		uri = "shroudb-courier://" + uri[len("shroudb-courier+tls://"):]
	case strings.HasPrefix(uri, "shroudb-courier://"):
		// ok
	default:
		return nil, fmt.Errorf("shroudb_courier: invalid URI scheme (expected shroudb-courier:// or shroudb-courier+tls://): %s", uri)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("shroudb_courier: invalid URI: %w", err)
	}

	cfg.host = u.Hostname()
	if cfg.host == "" {
		cfg.host = "localhost"
	}

	if p := u.Port(); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			cfg.port = n
		}
	}

	if u.User != nil {
		cfg.authToken = u.User.Username()
	}

	cfg.keyspace = strings.TrimPrefix(u.Path, "/")

	return cfg, nil
}

// Auth — Authenticate the connection
func (c *Client) Auth(token string) (*AuthResponse, error) {
	args := []string{
		"AUTH",
		token,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseAuthResponse(m), nil
}

// ChannelInfo — Get subscriber count for a WebSocket channel
func (c *Client) ChannelInfo(channel string) (*ChannelInfoResponse, error) {
	args := []string{
		"CHANNEL_INFO",
		channel,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseChannelInfoResponse(m), nil
}

// ChannelList — List all active WebSocket channels
func (c *Client) ChannelList() (*ChannelListResponse, error) {
	args := []string{
		"CHANNEL_LIST",
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseChannelListResponse(m), nil
}

// Connections — Get total WebSocket connections
func (c *Client) Connections() (*ConnectionsResponse, error) {
	args := []string{
		"CONNECTIONS",
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseConnectionsResponse(m), nil
}

// Deliver — Deliver a notification (decrypts recipient, renders template, sends via adapter)
func (c *Client) Deliver(json string) (*DeliverResponse, error) {
	args := []string{
		"DELIVER",
		json,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseDeliverResponse(m), nil
}

// Health — Check server health
func (c *Client) Health() error {
	args := []string{
		"HEALTH",
	}
	_, err := c.exec(args...)
	return err
}

// TemplateInfo — Get information about a specific template
func (c *Client) TemplateInfo(name string) (*TemplateInfoResponse, error) {
	args := []string{
		"TEMPLATE_INFO",
		name,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseTemplateInfoResponse(m), nil
}

// TemplateList — List all loaded templates
func (c *Client) TemplateList() (*TemplateListResponse, error) {
	args := []string{
		"TEMPLATE_LIST",
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseTemplateListResponse(m), nil
}

// TemplateReload — Reload templates from disk
func (c *Client) TemplateReload() (*TemplateReloadResponse, error) {
	args := []string{
		"TEMPLATE_RELOAD",
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseTemplateReloadResponse(m), nil
}
