// ShroudbAuth request/response types.
//
// Auto-generated from shroudb-auth protocol spec. Do not edit.

package shroudb_auth

// ChangePasswordResponse is the response from the ChangePassword endpoint.
type ChangePasswordResponse struct {
	// Always "OK"
	Status string `json:"status"`
}

// ForgotPasswordResponse is the response from the ForgotPassword endpoint.
type ForgotPasswordResponse struct {
	// Reset token TTL in seconds
	ExpiresIn *int64 `json:"expires_in,omitempty"`
	// Password reset token (only present if user exists)
	ResetToken *string `json:"reset_token,omitempty"`
	// Always "OK"
	Status string `json:"status"`
}

// HealthResponse is the response from the Health endpoint.
type HealthResponse struct {
	// "healthy" or "unhealthy"
	Status string `json:"status"`
}

// JwksResponse is the response from the Jwks endpoint.
type JwksResponse struct {
	// Array of JWK objects
	Keys []any `json:"keys"`
}

// LoginResponse is the response from the Login endpoint.
type LoginResponse struct {
	// JWT access token
	AccessToken string `json:"access_token"`
	// Access token TTL in seconds
	ExpiresIn int64 `json:"expires_in"`
	// Opaque refresh token
	RefreshToken string `json:"refresh_token"`
	// Authenticated user's ID
	UserId string `json:"user_id"`
}

// LogoutResponse is the response from the Logout endpoint.
type LogoutResponse struct {
	// Always "OK"
	Status string `json:"status"`
}

// LogoutAllResponse is the response from the LogoutAll endpoint.
type LogoutAllResponse struct {
	// Number of refresh token families revoked
	RevokedFamilies int64 `json:"revoked_families"`
	// Always "OK"
	Status string `json:"status"`
}

// RefreshResponse is the response from the Refresh endpoint.
type RefreshResponse struct {
	// New JWT access token
	AccessToken string `json:"access_token"`
	// Access token TTL in seconds
	ExpiresIn int64 `json:"expires_in"`
	// New opaque refresh token
	RefreshToken string `json:"refresh_token"`
}

// ResetPasswordResponse is the response from the ResetPassword endpoint.
type ResetPasswordResponse struct {
	// Always "OK"
	Status string `json:"status"`
}

// SessionResponse is the response from the Session endpoint.
type SessionResponse struct {
	// JWT claims from the access token
	Claims map[string]any `json:"claims"`
	// Token expiration as Unix timestamp
	ExpiresAt *int64 `json:"expires_at,omitempty"`
	// Authenticated user's ID
	UserId *string `json:"user_id,omitempty"`
}

// SessionsResponse is the response from the Sessions endpoint.
type SessionsResponse struct {
	// Array of active session objects
	ActiveSessions []any `json:"active_sessions"`
	// Authenticated user's ID
	UserId string `json:"user_id"`
}

// SignupOptions are optional parameters for Signup.
type SignupOptions struct {
	// Optional user metadata
	Metadata map[string]any `json:"metadata,omitempty"`
}

// SignupResponse is the response from the Signup endpoint.
type SignupResponse struct {
	// JWT access token
	AccessToken string `json:"access_token"`
	// Access token TTL in seconds
	ExpiresIn int64 `json:"expires_in"`
	// Opaque refresh token
	RefreshToken string `json:"refresh_token"`
	// The registered user's ID
	UserId string `json:"user_id"`
}

