// ShroudbCourier response types.
//
// Auto-generated from shroudb-courier protocol spec. Do not edit.

package shroudb_courier

// AuthResponse is the response from the AUTH command.
type AuthResponse struct {
	Status any // OK on success
}

func parseAuthResponse(m map[string]any) *AuthResponse {
	r := &AuthResponse{}
	if v, ok := m["status"].(any); ok { r.Status = v }
	return r
}

// ChannelInfoResponse is the response from the CHANNEL_INFO command.
type ChannelInfoResponse struct {
	Channel any // Channel name
	Subscribers any // Number of active subscribers
}

func parseChannelInfoResponse(m map[string]any) *ChannelInfoResponse {
	r := &ChannelInfoResponse{}
	if v, ok := m["channel"].(any); ok { r.Channel = v }
	if v, ok := m["subscribers"].(any); ok { r.Subscribers = v }
	return r
}

// ChannelListResponse is the response from the CHANNEL_LIST command.
type ChannelListResponse struct {
	Channels any // List of active channels with subscriber counts
}

func parseChannelListResponse(m map[string]any) *ChannelListResponse {
	r := &ChannelListResponse{}
	if v, ok := m["channels"].(any); ok { r.Channels = v }
	return r
}

// ConnectionsResponse is the response from the CONNECTIONS command.
type ConnectionsResponse struct {
	Connections any // Total active WebSocket connections
}

func parseConnectionsResponse(m map[string]any) *ConnectionsResponse {
	r := &ConnectionsResponse{}
	if v, ok := m["connections"].(any); ok { r.Connections = v }
	return r
}

// DeliverResponse is the response from the DELIVER command.
type DeliverResponse struct {
	DeliveryId string // Unique delivery identifier
	Channel string // Channel used for delivery
	Status any // Delivery status (sent, queued)
}

func parseDeliverResponse(m map[string]any) *DeliverResponse {
	r := &DeliverResponse{}
	if v, ok := m["delivery_id"].(string); ok { r.DeliveryId = v }
	if v, ok := m["channel"].(string); ok { r.Channel = v }
	if v, ok := m["status"].(any); ok { r.Status = v }
	return r
}

// TemplateInfoResponse is the response from the TEMPLATE_INFO command.
type TemplateInfoResponse struct {
	Name string // Template name
	Channels any // Supported delivery channels
	Variables any // Template variables
	LoadedAt any // When the template was loaded (RFC 3339)
}

func parseTemplateInfoResponse(m map[string]any) *TemplateInfoResponse {
	r := &TemplateInfoResponse{}
	if v, ok := m["name"].(string); ok { r.Name = v }
	if v, ok := m["channels"].(any); ok { r.Channels = v }
	if v, ok := m["variables"].(any); ok { r.Variables = v }
	if v, ok := m["loaded_at"].(any); ok { r.LoadedAt = v }
	return r
}

// TemplateListResponse is the response from the TEMPLATE_LIST command.
type TemplateListResponse struct {
	Templates any // List of template names
}

func parseTemplateListResponse(m map[string]any) *TemplateListResponse {
	r := &TemplateListResponse{}
	if v, ok := m["templates"].(any); ok { r.Templates = v }
	return r
}

// TemplateReloadResponse is the response from the TEMPLATE_RELOAD command.
type TemplateReloadResponse struct {
	Count any // Number of templates loaded
}

func parseTemplateReloadResponse(m map[string]any) *TemplateReloadResponse {
	r := &TemplateReloadResponse{}
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
