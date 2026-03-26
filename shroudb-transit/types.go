// ShroudbTransit response types.
//
// Auto-generated from shroudb-transit protocol spec. Do not edit.

package shroudb_transit

// DecryptResponse is the response from the DECRYPT command.
type DecryptResponse struct {
	Plaintext string // Decrypted data
}

func parseDecryptResponse(m map[string]any) *DecryptResponse {
	r := &DecryptResponse{}
	if v, ok := m["plaintext"].(string); ok { r.Plaintext = v }
	return r
}

// EncryptResponse is the response from the ENCRYPT command.
type EncryptResponse struct {
	Ciphertext string // Encrypted data with embedded key version
	KeyVersion int32 // Key version used for encryption
}

func parseEncryptResponse(m map[string]any) *EncryptResponse {
	r := &EncryptResponse{}
	if v, ok := m["ciphertext"].(string); ok { r.Ciphertext = v }
	if v, ok := m["key_version"].(int32); ok { r.KeyVersion = v }
	return r
}

// GenerateDataKeyResponse is the response from the GENERATE_DATA_KEY command.
type GenerateDataKeyResponse struct {
	PlaintextKey string // Plaintext DEK (use for local encryption, then discard)
	WrappedKey string // Wrapped DEK (store alongside ciphertext, unwrap via DECRYPT)
	KeyVersion int32 // Key version used to wrap
}

func parseGenerateDataKeyResponse(m map[string]any) *GenerateDataKeyResponse {
	r := &GenerateDataKeyResponse{}
	if v, ok := m["plaintext_key"].(string); ok { r.PlaintextKey = v }
	if v, ok := m["wrapped_key"].(string); ok { r.WrappedKey = v }
	if v, ok := m["key_version"].(int32); ok { r.KeyVersion = v }
	return r
}

// KeyInfoResponse is the response from the KEY_INFO command.
type KeyInfoResponse struct {
	Keyring string // Keyring name
	Type any // Keyring type (aes256-gcm)
	ActiveVersion int32 // Currently active key version
	Versions any // List of all key versions with state
}

func parseKeyInfoResponse(m map[string]any) *KeyInfoResponse {
	r := &KeyInfoResponse{}
	if v, ok := m["keyring"].(string); ok { r.Keyring = v }
	if v, ok := m["type"].(any); ok { r.Type = v }
	if v, ok := m["active_version"].(int32); ok { r.ActiveVersion = v }
	if v, ok := m["versions"].(any); ok { r.Versions = v }
	return r
}

// RewrapResponse is the response from the REWRAP command.
type RewrapResponse struct {
	Ciphertext string // Re-encrypted ciphertext with new key version
	KeyVersion int32 // New key version used
}

func parseRewrapResponse(m map[string]any) *RewrapResponse {
	r := &RewrapResponse{}
	if v, ok := m["ciphertext"].(string); ok { r.Ciphertext = v }
	if v, ok := m["key_version"].(int32); ok { r.KeyVersion = v }
	return r
}

// RotateResponse is the response from the ROTATE command.
type RotateResponse struct {
	KeyVersion int32 // New active key version
	PreviousVersion *int32 // Previous active key version
}

func parseRotateResponse(m map[string]any) *RotateResponse {
	r := &RotateResponse{}
	if v, ok := m["key_version"].(int32); ok { r.KeyVersion = v }
	if v, ok := m["previous_version"].(int32); ok { r.PreviousVersion = &v }
	return r
}

// SignResponse is the response from the SIGN command.
type SignResponse struct {
	Signature string // Detached signature
	KeyVersion int32 // Key version used
}

func parseSignResponse(m map[string]any) *SignResponse {
	r := &SignResponse{}
	if v, ok := m["signature"].(string); ok { r.Signature = v }
	if v, ok := m["key_version"].(int32); ok { r.KeyVersion = v }
	return r
}

// VerifySignatureResponse is the response from the VERIFY_SIGNATURE command.
type VerifySignatureResponse struct {
	Valid any // Whether the signature is valid
}

func parseVerifySignatureResponse(m map[string]any) *VerifySignatureResponse {
	r := &VerifySignatureResponse{}
	if v, ok := m["valid"].(any); ok { r.Valid = v }
	return r
}

// SubscriptionEvent represents a real-time event from a SUBSCRIBE stream.
type SubscriptionEvent struct {
	EventType string
	Keyspace  string
	Detail    string
	Timestamp int64
}
