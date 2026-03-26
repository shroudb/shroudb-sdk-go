// ShroudbPulse response types.
//
// Auto-generated from shroudb-pulse protocol spec. Do not edit.

package shroudb_pulse

// ActorsResponse is the response from the ACTORS command.
type ActorsResponse struct {
	Actors any // Actors ranked by activity
}

func parseActorsResponse(m map[string]any) *ActorsResponse {
	r := &ActorsResponse{}
	if v, ok := m["actors"].(any); ok { r.Actors = v }
	return r
}

// AuthResponse is the response from the AUTH command.
type AuthResponse struct {
	Status any // OK on success
}

func parseAuthResponse(m map[string]any) *AuthResponse {
	r := &AuthResponse{}
	if v, ok := m["status"].(any); ok { r.Status = v }
	return r
}

// CountResponse is the response from the COUNT command.
type CountResponse struct {
	Count any // Number of matching events
}

func parseCountResponse(m map[string]any) *CountResponse {
	r := &CountResponse{}
	if v, ok := m["count"].(any); ok { r.Count = v }
	return r
}

// ErrorsResponse is the response from the ERRORS command.
type ErrorsResponse struct {
	ErrorRates any // Operations with error counts and rates
}

func parseErrorsResponse(m map[string]any) *ErrorsResponse {
	r := &ErrorsResponse{}
	if v, ok := m["error_rates"].(any); ok { r.ErrorRates = v }
	return r
}

// HotspotsResponse is the response from the HOTSPOTS command.
type HotspotsResponse struct {
	Hotspots any // Resources ranked by activity
}

func parseHotspotsResponse(m map[string]any) *HotspotsResponse {
	r := &HotspotsResponse{}
	if v, ok := m["hotspots"].(any); ok { r.Hotspots = v }
	return r
}

// IngestResponse is the response from the INGEST command.
type IngestResponse struct {
	Id any // Assigned event identifier
}

func parseIngestResponse(m map[string]any) *IngestResponse {
	r := &IngestResponse{}
	if v, ok := m["id"].(any); ok { r.Id = v }
	return r
}

// IngestBatchResponse is the response from the INGEST_BATCH command.
type IngestBatchResponse struct {
	Count any // Number of events ingested
	Ids any // Assigned event identifiers
}

func parseIngestBatchResponse(m map[string]any) *IngestBatchResponse {
	r := &IngestBatchResponse{}
	if v, ok := m["count"].(any); ok { r.Count = v }
	if v, ok := m["ids"].(any); ok { r.Ids = v }
	return r
}

// QueryResponse is the response from the QUERY command.
type QueryResponse struct {
	Events any // Matching audit events
}

func parseQueryResponse(m map[string]any) *QueryResponse {
	r := &QueryResponse{}
	if v, ok := m["events"].(any); ok { r.Events = v }
	return r
}

// SourceListResponse is the response from the SOURCE_LIST command.
type SourceListResponse struct {
	Sources any // List of configured source names
}

func parseSourceListResponse(m map[string]any) *SourceListResponse {
	r := &SourceListResponse{}
	if v, ok := m["sources"].(any); ok { r.Sources = v }
	return r
}

// SourceStatusResponse is the response from the SOURCE_STATUS command.
type SourceStatusResponse struct {
	Sources any // Per-source ingestion stats (name, count, last_seen, lag)
}

func parseSourceStatusResponse(m map[string]any) *SourceStatusResponse {
	r := &SourceStatusResponse{}
	if v, ok := m["sources"].(any); ok { r.Sources = v }
	return r
}

// SubscriptionEvent represents a real-time event from a SUBSCRIBE stream.
type SubscriptionEvent struct {
	EventType string
	Keyspace  string
	Detail    string
	Timestamp int64
}
