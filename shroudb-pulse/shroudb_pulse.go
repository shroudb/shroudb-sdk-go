// Package shroudb_pulse provides a client for the ShroudbPulse Observability plane for unified audit event streaming.
//
// Auto-generated from shroudb-pulse protocol spec. Do not edit.
//
// Connect using a ShroudbPulse URI:
//
//	client, err := shroudb_pulse.Connect("shroudb-pulse://localhost")
//	result, err := client.Issue(ctx, "my-keyspace", &shroudb_pulse.IssueOptions{TTL: 3600})
//	fmt.Println(result.CredentialID, result.Token)
//	client.Close()
package shroudb_pulse

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Client is a ShroudbPulse client backed by a connection pool.
type Client struct {
	pool *pool
	host string
	port int
	tls  bool
	auth string
}

// Connect creates a new Client from a ShroudbPulse URI.
//
// Supported URI formats:
//
//	shroudb-pulse://localhost
//	shroudb-pulse://localhost:7099
//	shroudb-pulse+tls://prod.example.com
//	shroudb-pulse://mytoken@localhost:7099
//	shroudb-pulse://mytoken@localhost/sessions
//	shroudb-pulse+tls://tok@host:7099/keys
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
		return nil, fmt.Errorf("shroudb_pulse: expected map response, got %T", result)
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
		return nil, fmt.Errorf("shroudb_pulse: expected map response for SUBSCRIBE, got %T", resp)
	}
	if status, _ := m["status"].(string); status != "OK" {
		conn.close()
		return nil, fmt.Errorf("shroudb_pulse: subscribe failed: %v", m)
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
	case strings.HasPrefix(uri, "shroudb-pulse+tls://"):
		cfg.tls = true
		uri = "shroudb-pulse://" + uri[len("shroudb-pulse+tls://"):]
	case strings.HasPrefix(uri, "shroudb-pulse://"):
		// ok
	default:
		return nil, fmt.Errorf("shroudb_pulse: invalid URI scheme (expected shroudb-pulse:// or shroudb-pulse+tls://): %s", uri)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("shroudb_pulse: invalid URI: %w", err)
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

// ActorsOptions are optional parameters for Actors.
type ActorsOptions struct {
	Window string
}

// Actors — Most active actors
func (c *Client) Actors(opts *ActorsOptions) (*ActorsResponse, error) {
	args := []string{
		"ACTORS",
	}
	if opts != nil {
		if opts.Window != "" {
			args = append(args, "WINDOW", opts.Window)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseActorsResponse(m), nil
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

// Count — Count events matching filter arguments
func (c *Client) Count() (*CountResponse, error) {
	args := []string{
		"COUNT",
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseCountResponse(m), nil
}

// ErrorsOptions are optional parameters for Errors.
type ErrorsOptions struct {
	Engine string
	Window string
}

// Errors — Per-operation error rates
func (c *Client) Errors(opts *ErrorsOptions) (*ErrorsResponse, error) {
	args := []string{
		"ERRORS",
	}
	if opts != nil {
		if opts.Engine != "" {
			args = append(args, "ENGINE", opts.Engine)
		}
		if opts.Window != "" {
			args = append(args, "WINDOW", opts.Window)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseErrorsResponse(m), nil
}

// Health — Check server health
func (c *Client) Health() error {
	args := []string{
		"HEALTH",
	}
	_, err := c.exec(args...)
	return err
}

// HotspotsOptions are optional parameters for Hotspots.
type HotspotsOptions struct {
	Engine string
	Window string
}

// Hotspots — Find hotspot resources with highest activity
func (c *Client) Hotspots(opts *HotspotsOptions) (*HotspotsResponse, error) {
	args := []string{
		"HOTSPOTS",
	}
	if opts != nil {
		if opts.Engine != "" {
			args = append(args, "ENGINE", opts.Engine)
		}
		if opts.Window != "" {
			args = append(args, "WINDOW", opts.Window)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseHotspotsResponse(m), nil
}

// Ingest — Ingest a single audit event
func (c *Client) Ingest(json string) (*IngestResponse, error) {
	args := []string{
		"INGEST",
		json,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseIngestResponse(m), nil
}

// IngestBatch — Ingest a batch of audit events
func (c *Client) IngestBatch(json string) (*IngestBatchResponse, error) {
	args := []string{
		"INGEST_BATCH",
		json,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseIngestBatchResponse(m), nil
}

// Query — Query events with filter arguments
func (c *Client) Query() (*QueryResponse, error) {
	args := []string{
		"QUERY",
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseQueryResponse(m), nil
}

// SourceList — List configured event sources
func (c *Client) SourceList() (*SourceListResponse, error) {
	args := []string{
		"SOURCE_LIST",
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseSourceListResponse(m), nil
}

// SourceStatus — Show per-source ingestion statistics
func (c *Client) SourceStatus() (*SourceStatusResponse, error) {
	args := []string{
		"SOURCE_STATUS",
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseSourceStatusResponse(m), nil
}
