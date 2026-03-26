// Package shroudb_transit provides a client for the ShroudbTransit Encryption-as-a-service.
//
// Auto-generated from shroudb-transit protocol spec. Do not edit.
//
// Connect using a ShroudbTransit URI:
//
//	client, err := shroudb_transit.Connect("shroudb-transit://localhost")
//	result, err := client.Issue(ctx, "my-keyspace", &shroudb_transit.IssueOptions{TTL: 3600})
//	fmt.Println(result.CredentialID, result.Token)
//	client.Close()
package shroudb_transit

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Client is a ShroudbTransit client backed by a connection pool.
type Client struct {
	pool *pool
	host string
	port int
	tls  bool
	auth string
}

// Connect creates a new Client from a ShroudbTransit URI.
//
// Supported URI formats:
//
//	shroudb-transit://localhost
//	shroudb-transit://localhost:6499
//	shroudb-transit+tls://prod.example.com
//	shroudb-transit://mytoken@localhost:6499
//	shroudb-transit://mytoken@localhost/sessions
//	shroudb-transit+tls://tok@host:6499/keys
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
		return nil, fmt.Errorf("shroudb_transit: expected map response, got %T", result)
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
		return nil, fmt.Errorf("shroudb_transit: expected map response for SUBSCRIBE, got %T", resp)
	}
	if status, _ := m["status"].(string); status != "OK" {
		conn.close()
		return nil, fmt.Errorf("shroudb_transit: subscribe failed: %v", m)
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
	case strings.HasPrefix(uri, "shroudb-transit+tls://"):
		cfg.tls = true
		uri = "shroudb-transit://" + uri[len("shroudb-transit+tls://"):]
	case strings.HasPrefix(uri, "shroudb-transit://"):
		// ok
	default:
		return nil, fmt.Errorf("shroudb_transit: invalid URI scheme (expected shroudb-transit:// or shroudb-transit+tls://): %s", uri)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("shroudb_transit: invalid URI: %w", err)
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

// DecryptOptions are optional parameters for Decrypt.
type DecryptOptions struct {
	Context string
}

// Decrypt — Decrypt ciphertext using the embedded key version
func (c *Client) Decrypt(keyring string, ciphertext string, opts *DecryptOptions) (*DecryptResponse, error) {
	args := []string{
		"DECRYPT",
		keyring,
		ciphertext,
	}
	if opts != nil {
		if opts.Context != "" {
			args = append(args, "CONTEXT", opts.Context)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseDecryptResponse(m), nil
}

// EncryptOptions are optional parameters for Encrypt.
type EncryptOptions struct {
	Context string
	KeyVersion string
	Convergent string
}

// Encrypt — Encrypt plaintext with the active key version
func (c *Client) Encrypt(keyring string, plaintext string, opts *EncryptOptions) (*EncryptResponse, error) {
	args := []string{
		"ENCRYPT",
		keyring,
		plaintext,
	}
	if opts != nil {
		if opts.Context != "" {
			args = append(args, "CONTEXT", opts.Context)
		}
		if opts.KeyVersion != "" {
			args = append(args, "KEY_VERSION", opts.KeyVersion)
		}
		if opts.Convergent != "" {
			args = append(args, "CONVERGENT", opts.Convergent)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseEncryptResponse(m), nil
}

// GenerateDataKeyOptions are optional parameters for GenerateDataKey.
type GenerateDataKeyOptions struct {
	Bits string
}

// GenerateDataKey — Generate a data encryption key (envelope encryption pattern)
func (c *Client) GenerateDataKey(keyring string, opts *GenerateDataKeyOptions) (*GenerateDataKeyResponse, error) {
	args := []string{
		"GENERATE_DATA_KEY",
		keyring,
	}
	if opts != nil {
		if opts.Bits != "" {
			args = append(args, "BITS", opts.Bits)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseGenerateDataKeyResponse(m), nil
}

// Health — Check server health
func (c *Client) Health(keyring string) error {
	args := []string{
		"HEALTH",
	}
	if keyring != "" {
		args = append(args, keyring)
	}
	_, err := c.exec(args...)
	return err
}

// KeyInfo — Get keyring metadata and key version information
func (c *Client) KeyInfo(keyring string) (*KeyInfoResponse, error) {
	args := []string{
		"KEY_INFO",
		keyring,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseKeyInfoResponse(m), nil
}

// RewrapOptions are optional parameters for Rewrap.
type RewrapOptions struct {
	Context string
}

// Rewrap — Re-encrypt ciphertext with the current active key version
func (c *Client) Rewrap(keyring string, ciphertext string, opts *RewrapOptions) (*RewrapResponse, error) {
	args := []string{
		"REWRAP",
		keyring,
		ciphertext,
	}
	if opts != nil {
		if opts.Context != "" {
			args = append(args, "CONTEXT", opts.Context)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseRewrapResponse(m), nil
}

// RotateOptions are optional parameters for Rotate.
type RotateOptions struct {
	Force string
	Dryrun string
}

// Rotate — Rotate the keyring to a new key version
func (c *Client) Rotate(keyring string, opts *RotateOptions) (*RotateResponse, error) {
	args := []string{
		"ROTATE",
		keyring,
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
	return parseRotateResponse(m), nil
}

// SignOptions are optional parameters for Sign.
type SignOptions struct {
	Algorithm string
}

// Sign — Create a detached signature
func (c *Client) Sign(keyring string, data string, opts *SignOptions) (*SignResponse, error) {
	args := []string{
		"SIGN",
		keyring,
		data,
	}
	if opts != nil {
		if opts.Algorithm != "" {
			args = append(args, "ALGORITHM", opts.Algorithm)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseSignResponse(m), nil
}

// VerifySignature — Verify a detached signature
func (c *Client) VerifySignature(keyring string, data string, signature string) (*VerifySignatureResponse, error) {
	args := []string{
		"VERIFY_SIGNATURE",
		keyring,
		data,
		signature,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseVerifySignatureResponse(m), nil
}
