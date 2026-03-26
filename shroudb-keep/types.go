// ShroudbKeep response types.
//
// Auto-generated from shroudb-keep protocol spec. Do not edit.

package shroudb_keep

// AuthResponse is the response from the AUTH command.
type AuthResponse struct {
	Status any // OK on success
}

func parseAuthResponse(m map[string]any) *AuthResponse {
	r := &AuthResponse{}
	if v, ok := m["status"].(any); ok { r.Status = v }
	return r
}

// DeleteResponse is the response from the DELETE command.
type DeleteResponse struct {
	Path string // Deleted secret path
	DeletedAt any // Deletion timestamp (RFC 3339)
}

func parseDeleteResponse(m map[string]any) *DeleteResponse {
	r := &DeleteResponse{}
	if v, ok := m["path"].(string); ok { r.Path = v }
	if v, ok := m["deleted_at"].(any); ok { r.DeletedAt = v }
	return r
}

// GetResponse is the response from the GET command.
type GetResponse struct {
	Path string // Secret path
	Value string // Base64-encoded secret value
	Version int32 // Version number
	Meta string // JSON metadata
	CreatedAt any // Version creation time (RFC 3339)
}

func parseGetResponse(m map[string]any) *GetResponse {
	r := &GetResponse{}
	if v, ok := m["path"].(string); ok { r.Path = v }
	if v, ok := m["value"].(string); ok { r.Value = v }
	if v, ok := m["version"].(int32); ok { r.Version = v }
	if v, ok := m["meta"].(string); ok { r.Meta = v }
	if v, ok := m["created_at"].(any); ok { r.CreatedAt = v }
	return r
}

// ListResponse is the response from the LIST command.
type ListResponse struct {
	Paths any // List of matching secret paths
}

func parseListResponse(m map[string]any) *ListResponse {
	r := &ListResponse{}
	if v, ok := m["paths"].(any); ok { r.Paths = v }
	return r
}

// PutResponse is the response from the PUT command.
type PutResponse struct {
	Path string // Secret path
	Version int32 // Created version number
}

func parsePutResponse(m map[string]any) *PutResponse {
	r := &PutResponse{}
	if v, ok := m["path"].(string); ok { r.Path = v }
	if v, ok := m["version"].(int32); ok { r.Version = v }
	return r
}

// RotateResponse is the response from the ROTATE command.
type RotateResponse struct {
	Path string // Secret path
	Version int32 // New version number
}

func parseRotateResponse(m map[string]any) *RotateResponse {
	r := &RotateResponse{}
	if v, ok := m["path"].(string); ok { r.Path = v }
	if v, ok := m["version"].(int32); ok { r.Version = v }
	return r
}

// VersionsResponse is the response from the VERSIONS command.
type VersionsResponse struct {
	Path string // Secret path
	Versions any // List of version summaries with timestamps
}

func parseVersionsResponse(m map[string]any) *VersionsResponse {
	r := &VersionsResponse{}
	if v, ok := m["path"].(string); ok { r.Path = v }
	if v, ok := m["versions"].(any); ok { r.Versions = v }
	return r
}

// SubscriptionEvent represents a real-time event from a SUBSCRIBE stream.
type SubscriptionEvent struct {
	EventType string
	Keyspace  string
	Detail    string
	Timestamp int64
}
