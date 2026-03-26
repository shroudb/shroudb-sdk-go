// ShroudbSentry response types.
//
// Auto-generated from shroudb-sentry protocol spec. Do not edit.

package shroudb_sentry

// AuthResponse is the response from the AUTH command.
type AuthResponse struct {
	Status any // OK on success
}

func parseAuthResponse(m map[string]any) *AuthResponse {
	r := &AuthResponse{}
	if v, ok := m["status"].(any); ok { r.Status = v }
	return r
}

// EvaluateResponse is the response from the EVALUATE command.
type EvaluateResponse struct {
	Decision string // Authorization decision (allow or deny)
	Token string // Signed JWT encoding the decision
	Reasons any // Policy evaluation reasons
}

func parseEvaluateResponse(m map[string]any) *EvaluateResponse {
	r := &EvaluateResponse{}
	if v, ok := m["decision"].(string); ok { r.Decision = v }
	if v, ok := m["token"].(string); ok { r.Token = v }
	if v, ok := m["reasons"].(any); ok { r.Reasons = v }
	return r
}

// KeyInfoResponse is the response from the KEY_INFO command.
type KeyInfoResponse struct {
	KeyId any // Current signing key identifier
	Algorithm any // Signing algorithm
	CreatedAt any // Key creation time (RFC 3339)
}

func parseKeyInfoResponse(m map[string]any) *KeyInfoResponse {
	r := &KeyInfoResponse{}
	if v, ok := m["key_id"].(any); ok { r.KeyId = v }
	if v, ok := m["algorithm"].(any); ok { r.Algorithm = v }
	if v, ok := m["created_at"].(any); ok { r.CreatedAt = v }
	return r
}

// KeyRotateResponse is the response from the KEY_ROTATE command.
type KeyRotateResponse struct {
	KeyId any // New signing key identifier
	PreviousKeyId *any // Previous signing key identifier
}

func parseKeyRotateResponse(m map[string]any) *KeyRotateResponse {
	r := &KeyRotateResponse{}
	if v, ok := m["key_id"].(any); ok { r.KeyId = v }
	if v, ok := m["previous_key_id"].(any); ok { r.PreviousKeyId = &v }
	return r
}

// PolicyInfoResponse is the response from the POLICY_INFO command.
type PolicyInfoResponse struct {
	Name string // Policy name
	Version any // Policy version
	Rules any // Number of rules in the policy
	LoadedAt any // When the policy was loaded (RFC 3339)
}

func parsePolicyInfoResponse(m map[string]any) *PolicyInfoResponse {
	r := &PolicyInfoResponse{}
	if v, ok := m["name"].(string); ok { r.Name = v }
	if v, ok := m["version"].(any); ok { r.Version = v }
	if v, ok := m["rules"].(any); ok { r.Rules = v }
	if v, ok := m["loaded_at"].(any); ok { r.LoadedAt = v }
	return r
}

// PolicyListResponse is the response from the POLICY_LIST command.
type PolicyListResponse struct {
	Policies any // List of policy names
}

func parsePolicyListResponse(m map[string]any) *PolicyListResponse {
	r := &PolicyListResponse{}
	if v, ok := m["policies"].(any); ok { r.Policies = v }
	return r
}

// PolicyReloadResponse is the response from the POLICY_RELOAD command.
type PolicyReloadResponse struct {
	Count any // Number of policies loaded
}

func parsePolicyReloadResponse(m map[string]any) *PolicyReloadResponse {
	r := &PolicyReloadResponse{}
	if v, ok := m["count"].(any); ok { r.Count = v }
	return r
}

// SubscriptionEvent represents a real-time event from a SUBSCRIBE stream.
type SubscriptionEvent struct {
	EventType string
	Keyspace  string
	Detail    string
	Timestamp int64
}
