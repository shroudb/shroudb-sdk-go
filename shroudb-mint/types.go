// ShroudbMint response types.
//
// Auto-generated from shroudb-mint protocol spec. Do not edit.

package shroudb_mint

// AuthResponse is the response from the AUTH command.
type AuthResponse struct {
	Status any // OK on success
}

func parseAuthResponse(m map[string]any) *AuthResponse {
	r := &AuthResponse{}
	if v, ok := m["status"].(any); ok { r.Status = v }
	return r
}

// CaCreateResponse is the response from the CA_CREATE command.
type CaCreateResponse struct {
	Ca string // Created CA name
	Serial string // CA certificate serial
	Certificate string // CA certificate in PEM format
}

func parseCaCreateResponse(m map[string]any) *CaCreateResponse {
	r := &CaCreateResponse{}
	if v, ok := m["ca"].(string); ok { r.Ca = v }
	if v, ok := m["serial"].(string); ok { r.Serial = v }
	if v, ok := m["certificate"].(string); ok { r.Certificate = v }
	return r
}

// CaExportResponse is the response from the CA_EXPORT command.
type CaExportResponse struct {
	Certificate string // CA certificate in requested format
}

func parseCaExportResponse(m map[string]any) *CaExportResponse {
	r := &CaExportResponse{}
	if v, ok := m["certificate"].(string); ok { r.Certificate = v }
	return r
}

// CaInfoResponse is the response from the CA_INFO command.
type CaInfoResponse struct {
	Ca string // CA name
	Algorithm string // Key algorithm
	Subject string // Subject DN
	Serial string // CA certificate serial
	NotBefore any // Validity start (RFC 3339)
	NotAfter any // Validity end (RFC 3339)
	IssuedCount any // Number of issued certificates
}

func parseCaInfoResponse(m map[string]any) *CaInfoResponse {
	r := &CaInfoResponse{}
	if v, ok := m["ca"].(string); ok { r.Ca = v }
	if v, ok := m["algorithm"].(string); ok { r.Algorithm = v }
	if v, ok := m["subject"].(string); ok { r.Subject = v }
	if v, ok := m["serial"].(string); ok { r.Serial = v }
	if v, ok := m["not_before"].(any); ok { r.NotBefore = v }
	if v, ok := m["not_after"].(any); ok { r.NotAfter = v }
	if v, ok := m["issued_count"].(any); ok { r.IssuedCount = v }
	return r
}

// CaListResponse is the response from the CA_LIST command.
type CaListResponse struct {
	Cas any // List of CA names
}

func parseCaListResponse(m map[string]any) *CaListResponse {
	r := &CaListResponse{}
	if v, ok := m["cas"].(any); ok { r.Cas = v }
	return r
}

// CaRotateResponse is the response from the CA_ROTATE command.
type CaRotateResponse struct {
	Serial string // New CA certificate serial
	PreviousSerial string // Previous CA certificate serial
}

func parseCaRotateResponse(m map[string]any) *CaRotateResponse {
	r := &CaRotateResponse{}
	if v, ok := m["serial"].(string); ok { r.Serial = v }
	if v, ok := m["previous_serial"].(string); ok { r.PreviousSerial = v }
	return r
}

// CrlInfoResponse is the response from the CRL_INFO command.
type CrlInfoResponse struct {
	Ca string // CA name
	CrlNumber any // Current CRL number
	LastUpdate any // Last CRL update (RFC 3339)
	NextUpdate any // Next CRL update (RFC 3339)
	RevokedCount any // Number of revoked certificates
}

func parseCrlInfoResponse(m map[string]any) *CrlInfoResponse {
	r := &CrlInfoResponse{}
	if v, ok := m["ca"].(string); ok { r.Ca = v }
	if v, ok := m["crl_number"].(any); ok { r.CrlNumber = v }
	if v, ok := m["last_update"].(any); ok { r.LastUpdate = v }
	if v, ok := m["next_update"].(any); ok { r.NextUpdate = v }
	if v, ok := m["revoked_count"].(any); ok { r.RevokedCount = v }
	return r
}

// InspectResponse is the response from the INSPECT command.
type InspectResponse struct {
	Serial string // Certificate serial
	Subject string // Subject DN
	NotBefore any // Validity start (RFC 3339)
	NotAfter any // Validity end (RFC 3339)
	State string // Certificate state
	Certificate string // Certificate PEM
}

func parseInspectResponse(m map[string]any) *InspectResponse {
	r := &InspectResponse{}
	if v, ok := m["serial"].(string); ok { r.Serial = v }
	if v, ok := m["subject"].(string); ok { r.Subject = v }
	if v, ok := m["not_before"].(any); ok { r.NotBefore = v }
	if v, ok := m["not_after"].(any); ok { r.NotAfter = v }
	if v, ok := m["state"].(string); ok { r.State = v }
	if v, ok := m["certificate"].(string); ok { r.Certificate = v }
	return r
}

// IssueResponse is the response from the ISSUE command.
type IssueResponse struct {
	Serial string // Certificate serial
	Certificate string // Issued certificate PEM
	PrivateKey string // Private key PEM
	Chain string // Full certificate chain PEM
	NotAfter any // Expiry (RFC 3339)
}

func parseIssueResponse(m map[string]any) *IssueResponse {
	r := &IssueResponse{}
	if v, ok := m["serial"].(string); ok { r.Serial = v }
	if v, ok := m["certificate"].(string); ok { r.Certificate = v }
	if v, ok := m["private_key"].(string); ok { r.PrivateKey = v }
	if v, ok := m["chain"].(string); ok { r.Chain = v }
	if v, ok := m["not_after"].(any); ok { r.NotAfter = v }
	return r
}

// IssueFromCsrResponse is the response from the ISSUE_FROM_CSR command.
type IssueFromCsrResponse struct {
	Serial string // Certificate serial
	Certificate string // Issued certificate PEM
	Chain string // Full certificate chain PEM
	NotAfter any // Expiry (RFC 3339)
}

func parseIssueFromCsrResponse(m map[string]any) *IssueFromCsrResponse {
	r := &IssueFromCsrResponse{}
	if v, ok := m["serial"].(string); ok { r.Serial = v }
	if v, ok := m["certificate"].(string); ok { r.Certificate = v }
	if v, ok := m["chain"].(string); ok { r.Chain = v }
	if v, ok := m["not_after"].(any); ok { r.NotAfter = v }
	return r
}

// ListCertsResponse is the response from the LIST_CERTS command.
type ListCertsResponse struct {
	Certificates any // List of certificate summaries
}

func parseListCertsResponse(m map[string]any) *ListCertsResponse {
	r := &ListCertsResponse{}
	if v, ok := m["certificates"].(any); ok { r.Certificates = v }
	return r
}

// RenewResponse is the response from the RENEW command.
type RenewResponse struct {
	Serial string // New certificate serial
	Certificate string // Renewed certificate PEM
	PrivateKey string // New private key PEM
	NotAfter any // New expiry (RFC 3339)
}

func parseRenewResponse(m map[string]any) *RenewResponse {
	r := &RenewResponse{}
	if v, ok := m["serial"].(string); ok { r.Serial = v }
	if v, ok := m["certificate"].(string); ok { r.Certificate = v }
	if v, ok := m["private_key"].(string); ok { r.PrivateKey = v }
	if v, ok := m["not_after"].(any); ok { r.NotAfter = v }
	return r
}

// RevokeResponse is the response from the REVOKE command.
type RevokeResponse struct {
	Serial string // Revoked certificate serial
	RevokedAt any // Revocation timestamp (RFC 3339)
}

func parseRevokeResponse(m map[string]any) *RevokeResponse {
	r := &RevokeResponse{}
	if v, ok := m["serial"].(string); ok { r.Serial = v }
	if v, ok := m["revoked_at"].(any); ok { r.RevokedAt = v }
	return r
}

// SubscriptionEvent represents a real-time event from a SUBSCRIBE stream.
type SubscriptionEvent struct {
	EventType string
	Keyspace  string
	Detail    string
	Timestamp int64
}
