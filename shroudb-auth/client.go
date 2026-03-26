// Package shroudb_auth provides a Go HTTP client for the ShroudbAuth Authentication service.
//
// Auto-generated from shroudb-auth protocol spec. Do not edit.
package shroudb_auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Auth mode constants used by the do helper.
const (
	authNone    = 0
	authAccess  = 1
	authRefresh = 2
)

// Client is an HTTP client for the ShroudbAuth Authentication service.
type Client struct {
	// BaseURL is the root URL of the ShroudbAuth server (e.g. "http://localhost:4001").
	BaseURL string

	// Keyspace is the default keyspace used for endpoints that require one.
	Keyspace string

	// AccessToken is the current JWT access token. Updated automatically
	// after Signup, Login, and Refresh calls.
	AccessToken string

	// RefreshToken is the current opaque refresh token. Updated automatically
	// after Signup, Login, and Refresh calls.
	RefreshToken string

	// HTTPClient is the underlying HTTP client. Defaults to http.DefaultClient.
	HTTPClient *http.Client
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithKeyspace sets the default keyspace.
func WithKeyspace(ks string) ClientOption {
	return func(c *Client) { c.Keyspace = ks }
}

// WithAccessToken sets the initial access token.
func WithAccessToken(t string) ClientOption {
	return func(c *Client) { c.AccessToken = t }
}

// WithRefreshToken sets the initial refresh token.
func WithRefreshToken(t string) ClientOption {
	return func(c *Client) { c.RefreshToken = t }
}

// WithHTTPClient sets a custom *http.Client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) { c.HTTPClient = hc }
}

// NewClient creates a new ShroudbAuth client.
//
// baseURL should include the scheme and host, e.g. "http://localhost:4001".
func NewClient(baseURL string, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// do performs an HTTP request and decodes the JSON response.
func (c *Client) do(ctx context.Context, method, path string, body any, result any, authMode int, expectedStatus int) error {
	url := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("shroudb_auth: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("shroudb_auth: create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	switch authMode {
	case authAccess:
		if c.AccessToken != "" {
			req.Header.Set("Authorization", "Bearer "+c.AccessToken)
		}
	case authRefresh:
		if c.RefreshToken != "" {
			req.Header.Set("Authorization", "Bearer "+c.RefreshToken)
		}
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("shroudb_auth: request %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("shroudb_auth: read response: %w", err)
	}

	if resp.StatusCode != expectedStatus {
		var apiErr ShroudbAuthError
		if json.Unmarshal(respBody, &apiErr) == nil && apiErr.Code != "" {
			return &apiErr
		}
		return &ShroudbAuthError{
			Code:    http.StatusText(resp.StatusCode),
			Message: string(respBody),
		}
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("shroudb_auth: decode response: %w", err)
		}
	}

	return nil
}

// ChangePassword — Change password for the currently authenticated user
func (c *Client) ChangePassword(ctx context.Context, newPassword string, oldPassword string) (*ChangePasswordResponse, error) {
	body := map[string]any{
		"new_password": newPassword,
		"old_password": oldPassword,
	}
	result := &ChangePasswordResponse{}
	err := c.do(ctx, "POST", fmt.Sprintf("/auth/%s/change-password", c.Keyspace), body, result, authAccess, 200)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ForgotPassword — Request a password reset token (always returns 200 to prevent enumeration)
func (c *Client) ForgotPassword(ctx context.Context, userId string) (*ForgotPasswordResponse, error) {
	body := map[string]any{
		"user_id": userId,
	}
	result := &ForgotPasswordResponse{}
	err := c.do(ctx, "POST", fmt.Sprintf("/auth/%s/forgot-password", c.Keyspace), body, result, authNone, 200)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Health — Health check endpoint
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	result := &HealthResponse{}
	err := c.do(ctx, "GET", "/auth/health", nil, result, authNone, 200)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Jwks — Public JSON Web Key Set for verifying access tokens
func (c *Client) Jwks(ctx context.Context) (*JwksResponse, error) {
	result := &JwksResponse{}
	err := c.do(ctx, "GET", fmt.Sprintf("/auth/%s/.well-known/jwks.json", c.Keyspace), nil, result, authNone, 200)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Login — Authenticate a user and receive access + refresh tokens
func (c *Client) Login(ctx context.Context, password string, userId string) (*LoginResponse, error) {
	body := map[string]any{
		"password": password,
		"user_id": userId,
	}
	result := &LoginResponse{}
	err := c.do(ctx, "POST", fmt.Sprintf("/auth/%s/login", c.Keyspace), body, result, authNone, 200)
	if err != nil {
		return nil, err
	}
	c.AccessToken = result.AccessToken
	c.RefreshToken = result.RefreshToken
	return result, nil
}

// Logout — Revoke the current refresh token family and clear cookies
func (c *Client) Logout(ctx context.Context) (*LogoutResponse, error) {
	result := &LogoutResponse{}
	err := c.do(ctx, "POST", fmt.Sprintf("/auth/%s/logout", c.Keyspace), nil, result, authRefresh, 200)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// LogoutAll — Revoke all refresh token families for a user
func (c *Client) LogoutAll(ctx context.Context, userId string) (*LogoutAllResponse, error) {
	body := map[string]any{
		"user_id": userId,
	}
	result := &LogoutAllResponse{}
	err := c.do(ctx, "POST", fmt.Sprintf("/auth/%s/logout-all", c.Keyspace), body, result, authAccess, 200)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Refresh — Exchange a refresh token for new access + refresh tokens
func (c *Client) Refresh(ctx context.Context) (*RefreshResponse, error) {
	result := &RefreshResponse{}
	err := c.do(ctx, "POST", fmt.Sprintf("/auth/%s/refresh", c.Keyspace), nil, result, authRefresh, 200)
	if err != nil {
		return nil, err
	}
	c.AccessToken = result.AccessToken
	c.RefreshToken = result.RefreshToken
	return result, nil
}

// ResetPassword — Reset password using a single-use reset token (revoked after use)
func (c *Client) ResetPassword(ctx context.Context, newPassword string, token string) (*ResetPasswordResponse, error) {
	body := map[string]any{
		"new_password": newPassword,
		"token": token,
	}
	result := &ResetPasswordResponse{}
	err := c.do(ctx, "POST", fmt.Sprintf("/auth/%s/reset-password", c.Keyspace), body, result, authNone, 200)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Session — Validate current session and return user info
func (c *Client) Session(ctx context.Context) (*SessionResponse, error) {
	result := &SessionResponse{}
	err := c.do(ctx, "GET", fmt.Sprintf("/auth/%s/session", c.Keyspace), nil, result, authAccess, 200)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Sessions — List active sessions (refresh token families) for the authenticated user
func (c *Client) Sessions(ctx context.Context) (*SessionsResponse, error) {
	result := &SessionsResponse{}
	err := c.do(ctx, "GET", fmt.Sprintf("/auth/%s/sessions", c.Keyspace), nil, result, authAccess, 200)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Signup — Register a new user and receive access + refresh tokens
func (c *Client) Signup(ctx context.Context, password string, userId string, opts *SignupOptions) (*SignupResponse, error) {
	body := map[string]any{
		"password": password,
		"user_id": userId,
	}
	if opts != nil {
		if opts.Metadata != nil {
			body["metadata"] = opts.Metadata
		}
	}
	result := &SignupResponse{}
	err := c.do(ctx, "POST", fmt.Sprintf("/auth/%s/signup", c.Keyspace), body, result, authNone, 201)
	if err != nil {
		return nil, err
	}
	c.AccessToken = result.AccessToken
	c.RefreshToken = result.RefreshToken
	return result, nil
}
