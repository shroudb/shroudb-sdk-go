// Package shroudb_keep provides a client for the ShroudbKeep Secrets manager.
//
// Auto-generated from shroudb-keep protocol spec. Do not edit.
//
// Connect using a ShroudbKeep URI:
//
//	client, err := shroudb_keep.Connect("shroudb-keep://localhost")
//	result, err := client.Issue(ctx, "my-keyspace", &shroudb_keep.IssueOptions{TTL: 3600})
//	fmt.Println(result.CredentialID, result.Token)
//	client.Close()
package shroudb_keep

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Client is a ShroudbKeep client backed by a connection pool.
type Client struct {
	pool *pool
	host string
	port int
	tls  bool
	auth string
}

// Connect creates a new Client from a ShroudbKeep URI.
//
// Supported URI formats:
//
//	shroudb-keep://localhost
//	shroudb-keep://localhost:6899
//	shroudb-keep+tls://prod.example.com
//	shroudb-keep://mytoken@localhost:6899
//	shroudb-keep://mytoken@localhost/sessions
//	shroudb-keep+tls://tok@host:6899/keys
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
		return nil, fmt.Errorf("shroudb_keep: expected map response, got %T", result)
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
		return nil, fmt.Errorf("shroudb_keep: expected map response for SUBSCRIBE, got %T", resp)
	}
	if status, _ := m["status"].(string); status != "OK" {
		conn.close()
		return nil, fmt.Errorf("shroudb_keep: subscribe failed: %v", m)
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
	case strings.HasPrefix(uri, "shroudb-keep+tls://"):
		cfg.tls = true
		uri = "shroudb-keep://" + uri[len("shroudb-keep+tls://"):]
	case strings.HasPrefix(uri, "shroudb-keep://"):
		// ok
	default:
		return nil, fmt.Errorf("shroudb_keep: invalid URI scheme (expected shroudb-keep:// or shroudb-keep+tls://): %s", uri)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("shroudb_keep: invalid URI: %w", err)
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

// Delete — Soft-delete a secret
func (c *Client) Delete(path string) (*DeleteResponse, error) {
	args := []string{
		"DELETE",
		path,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseDeleteResponse(m), nil
}

// GetOptions are optional parameters for Get.
type GetOptions struct {
	Version string
}

// Get — Retrieve a secret (latest or specific version)
func (c *Client) Get(path string, opts *GetOptions) (*GetResponse, error) {
	args := []string{
		"GET",
		path,
	}
	if opts != nil {
		if opts.Version != "" {
			args = append(args, "VERSION", opts.Version)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseGetResponse(m), nil
}

// Health — Check server health
func (c *Client) Health(path string) error {
	args := []string{
		"HEALTH",
	}
	if path != "" {
		args = append(args, path)
	}
	_, err := c.exec(args...)
	return err
}

// List — List secret paths matching a prefix
func (c *Client) List(prefix string) (*ListResponse, error) {
	args := []string{
		"LIST",
	}
	if prefix != "" {
		args = append(args, prefix)
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseListResponse(m), nil
}

// PutOptions are optional parameters for Put.
type PutOptions struct {
	Meta string
}

// Put — Store a secret (creates a new version)
func (c *Client) Put(path string, value string, opts *PutOptions) (*PutResponse, error) {
	args := []string{
		"PUT",
		path,
		value,
	}
	if opts != nil {
		if opts.Meta != "" {
			args = append(args, "META", opts.Meta)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parsePutResponse(m), nil
}

// Rotate — Re-encrypt the latest version under the current key
func (c *Client) Rotate(path string) (*RotateResponse, error) {
	args := []string{
		"ROTATE",
		path,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseRotateResponse(m), nil
}

// Versions — Show version history for a secret
func (c *Client) Versions(path string) (*VersionsResponse, error) {
	args := []string{
		"VERSIONS",
		path,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseVersionsResponse(m), nil
}
