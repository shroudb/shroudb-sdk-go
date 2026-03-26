// Package shroudb_mint provides a client for the ShroudbMint Lightweight internal Certificate Authority.
//
// Auto-generated from shroudb-mint protocol spec. Do not edit.
//
// Connect using a ShroudbMint URI:
//
//	client, err := shroudb_mint.Connect("shroudb-mint://localhost")
//	result, err := client.Issue(ctx, "my-keyspace", &shroudb_mint.IssueOptions{TTL: 3600})
//	fmt.Println(result.CredentialID, result.Token)
//	client.Close()
package shroudb_mint

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Client is a ShroudbMint client backed by a connection pool.
type Client struct {
	pool *pool
	host string
	port int
	tls  bool
	auth string
}

// Connect creates a new Client from a ShroudbMint URI.
//
// Supported URI formats:
//
//	shroudb-mint://localhost
//	shroudb-mint://localhost:6699
//	shroudb-mint+tls://prod.example.com
//	shroudb-mint://mytoken@localhost:6699
//	shroudb-mint://mytoken@localhost/sessions
//	shroudb-mint+tls://tok@host:6699/keys
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
		return nil, fmt.Errorf("shroudb_mint: expected map response, got %T", result)
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
		return nil, fmt.Errorf("shroudb_mint: expected map response for SUBSCRIBE, got %T", resp)
	}
	if status, _ := m["status"].(string); status != "OK" {
		conn.close()
		return nil, fmt.Errorf("shroudb_mint: subscribe failed: %v", m)
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
	case strings.HasPrefix(uri, "shroudb-mint+tls://"):
		cfg.tls = true
		uri = "shroudb-mint://" + uri[len("shroudb-mint+tls://"):]
	case strings.HasPrefix(uri, "shroudb-mint://"):
		// ok
	default:
		return nil, fmt.Errorf("shroudb_mint: invalid URI scheme (expected shroudb-mint:// or shroudb-mint+tls://): %s", uri)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("shroudb_mint: invalid URI: %w", err)
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

// CaCreateOptions are optional parameters for CaCreate.
type CaCreateOptions struct {
	Parent string
}

// CaCreate — Create a new Certificate Authority
func (c *Client) CaCreate(ca string, algorithm string, subject string, ttl_days string, opts *CaCreateOptions) (*CaCreateResponse, error) {
	args := []string{
		"CA_CREATE",
		ca,
		algorithm,
		subject,
		ttl_days,
	}
	if opts != nil {
		if opts.Parent != "" {
			args = append(args, "PARENT", opts.Parent)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseCaCreateResponse(m), nil
}

// CaExportOptions are optional parameters for CaExport.
type CaExportOptions struct {
	Format string
}

// CaExport — Export the CA's public certificate
func (c *Client) CaExport(ca string, opts *CaExportOptions) (*CaExportResponse, error) {
	args := []string{
		"CA_EXPORT",
		ca,
	}
	if opts != nil {
		if opts.Format != "" {
			args = append(args, "FORMAT", opts.Format)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseCaExportResponse(m), nil
}

// CaInfo — Get information about a CA
func (c *Client) CaInfo(ca string) (*CaInfoResponse, error) {
	args := []string{
		"CA_INFO",
		ca,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseCaInfoResponse(m), nil
}

// CaList — List all CAs
func (c *Client) CaList() (*CaListResponse, error) {
	args := []string{
		"CA_LIST",
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseCaListResponse(m), nil
}

// CaRotateOptions are optional parameters for CaRotate.
type CaRotateOptions struct {
	Force string
	Dryrun string
}

// CaRotate — Rotate the CA's signing key
func (c *Client) CaRotate(ca string, opts *CaRotateOptions) (*CaRotateResponse, error) {
	args := []string{
		"CA_ROTATE",
		ca,
	}
	if opts != nil {
		if opts.Force != "" {
			args = append(args, "FORCE", opts.Force)
		}
		if opts.Dryrun != "" {
			args = append(args, "DRYRUN", opts.Dryrun)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseCaRotateResponse(m), nil
}

// CrlInfo — Get CRL information for a CA
func (c *Client) CrlInfo(ca string) (*CrlInfoResponse, error) {
	args := []string{
		"CRL_INFO",
		ca,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseCrlInfoResponse(m), nil
}

// Health — Check server health
func (c *Client) Health(ca string) error {
	args := []string{
		"HEALTH",
	}
	if ca != "" {
		args = append(args, ca)
	}
	_, err := c.exec(args...)
	return err
}

// Inspect — Inspect a certificate
func (c *Client) Inspect(ca string, serial string) (*InspectResponse, error) {
	args := []string{
		"INSPECT",
		ca,
		serial,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseInspectResponse(m), nil
}

// IssueOptions are optional parameters for Issue.
type IssueOptions struct {
	Ttl string
	SanDns string
	SanIp string
}

// Issue — Issue a new certificate
func (c *Client) Issue(ca string, subject string, profile string, opts *IssueOptions) (*IssueResponse, error) {
	args := []string{
		"ISSUE",
		ca,
		subject,
		profile,
	}
	if opts != nil {
		if opts.Ttl != "" {
			args = append(args, "TTL", opts.Ttl)
		}
		if opts.SanDns != "" {
			args = append(args, "SAN_DNS", opts.SanDns)
		}
		if opts.SanIp != "" {
			args = append(args, "SAN_IP", opts.SanIp)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseIssueResponse(m), nil
}

// IssueFromCsrOptions are optional parameters for IssueFromCsr.
type IssueFromCsrOptions struct {
	Ttl string
}

// IssueFromCsr — Issue a certificate from a CSR
func (c *Client) IssueFromCsr(ca string, csr_pem string, profile string, opts *IssueFromCsrOptions) (*IssueFromCsrResponse, error) {
	args := []string{
		"ISSUE_FROM_CSR",
		ca,
		csr_pem,
		profile,
	}
	if opts != nil {
		if opts.Ttl != "" {
			args = append(args, "TTL", opts.Ttl)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseIssueFromCsrResponse(m), nil
}

// ListCertsOptions are optional parameters for ListCerts.
type ListCertsOptions struct {
	State string
	Limit int64
	Offset int64
}

// ListCerts — List certificates for a CA
func (c *Client) ListCerts(ca string, opts *ListCertsOptions) (*ListCertsResponse, error) {
	args := []string{
		"LIST_CERTS",
		ca,
	}
	if opts != nil {
		if opts.State != "" {
			args = append(args, "STATE", opts.State)
		}
		if opts.Limit != 0 {
			args = append(args, "LIMIT", strconv.FormatInt(opts.Limit, 10))
		}
		if opts.Offset != 0 {
			args = append(args, "OFFSET", strconv.FormatInt(opts.Offset, 10))
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseListCertsResponse(m), nil
}

// RenewOptions are optional parameters for Renew.
type RenewOptions struct {
	Ttl string
}

// Renew — Renew a certificate
func (c *Client) Renew(ca string, serial string, opts *RenewOptions) (*RenewResponse, error) {
	args := []string{
		"RENEW",
		ca,
		serial,
	}
	if opts != nil {
		if opts.Ttl != "" {
			args = append(args, "TTL", opts.Ttl)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseRenewResponse(m), nil
}

// RevokeOptions are optional parameters for Revoke.
type RevokeOptions struct {
	Reason string
}

// Revoke — Revoke a certificate
func (c *Client) Revoke(ca string, serial string, opts *RevokeOptions) (*RevokeResponse, error) {
	args := []string{
		"REVOKE",
		ca,
		serial,
	}
	if opts != nil {
		if opts.Reason != "" {
			args = append(args, "REASON", opts.Reason)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseRevokeResponse(m), nil
}
