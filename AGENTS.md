# ShrouDB SDK — Agent Instructions

> Unified Go SDK for all ShrouDB security engines. Provides namespaced, type-safe access with built-in serialization.

## Quick Context

- **Module**: `github.com/shroudb/shroudb-go`
- **Transport**: RESP3 (direct engine connections) or HTTP (Moat gateway)
- **Pattern**: `db.Engine.Method(ctx, params)` — all methods take `context.Context`, return `(*Response, error)` or `error`
- **Serialization**: Handled internally — pass Go types, get typed structs back

## Connection

```go
import shroudb "github.com/shroudb/shroudb-go"

// Moat gateway (HTTP) — all engines through one endpoint
db, err := shroudb.New(shroudb.Options{Moat: "https://moat.example.com", Token: "my-token"})

// Direct — only the engines you need
db, err := shroudb.New(shroudb.Options{Cipher: "shroudb-cipher://token@host:6599"})

// Mixed — Moat default + direct overrides
db, err := shroudb.New(shroudb.Options{
	Moat:   "https://moat.example.com",
	Cipher: "shroudb-cipher://token@dedicated:6599",
	Token:  "moat-token",
})

// Always close when done
defer db.Close()
```

## `db.Shroudb` — Encrypted key-value database

| Method | Args | Returns | Description |
|--------|------|---------|-------------|
| `Auth` | `ctx, token` | `*ShroudbAuthResponse, error` | Authenticate the connection with a token |
| `CommandList` | `ctx` | `*ShroudbCommandListResponse, error` | List all supported commands |
| `ConfigGet` | `ctx, key` | `*ShroudbConfigGetResponse, error` | Read a runtime configuration value |
| `ConfigSet` | `ctx, key, value` | `error` | Set a runtime configuration value (admin only). Only registered config keys are accepted; unknown keys return an error. Values are type-checked against the key's schema (u64, bool, string). Valid keys: max_segment_bytes, max_segment_entries, snapshot_entry_threshold, snapshot_time_threshold_secs. |
| `Delete` | `ctx, namespace, key` | `*ShroudbDeleteResponse, error` | Delete a key by writing a tombstone |
| `Get` | `ctx, namespace, key, META?, opts` | `*ShroudbGetResponse, error` | Retrieve the value at a key |
| `Health` | `ctx` | `*ShroudbHealthResponse, error` | Check server health |
| `List` | `ctx, namespace, opts` | `*ShroudbListResponse, error` | List active keys in a namespace. Returns an error if the CURSOR value does not correspond to a key that exists in the namespace. |
| `NamespaceAlter` | `ctx, name, opts` | `error` | Update namespace configuration (enforce-on-write-only) |
| `NamespaceCreate` | `ctx, name, opts` | `error` | Create a new namespace |
| `NamespaceDrop` | `ctx, name, FORCE?` | `error` | Drop a namespace |
| `NamespaceInfo` | `ctx, name` | `*ShroudbNamespaceInfoResponse, error` | Get metadata about a namespace |
| `NamespaceList` | `ctx, opts` | `*ShroudbNamespaceListResponse, error` | List namespaces (filtered by token grants) |
| `NamespaceValidate` | `ctx, name` | `*ShroudbNamespaceValidateResponse, error` | Check existing entries against current MetaSchema |
| `Ping` | `ctx` | `*ShroudbPingResponse, error` | Test connectivity |
| `Pipeline` | `ctx, count` | `error` | Execute commands atomically (all succeed or all roll back) |
| `Put` | `ctx, namespace, key, value?, opts` | `*ShroudbPutResponse, error` | Store a value at the given key. Auto-increments version. |
| `Rekey` | `ctx` | `*ShroudbRekeyResponse, error` | Begin online rekey (zero-downtime master key rotation) |
| `RekeyStatus` | `ctx` | `*ShroudbRekeyStatusResponse, error` | Query progress of an in-flight rekey operation |
| `Subscribe` | `ctx, namespace, opts` | `error` | Subscribe to change events on a namespace |
| `Unsubscribe` | `ctx` | `error` | End the current subscription |
| `Versions` | `ctx, namespace, key, opts` | `*ShroudbVersionsResponse, error` | Retrieve version history for a key (most recent first) |

### Examples

```go
ctx := context.Background()
resp, err := db.Shroudb.ConfigGet(ctx, "key")
// resp.Key
err := db.Shroudb.ConfigSet(ctx, "key", "alice@example.com")
resp, err := db.Shroudb.Delete(ctx, "namespace", "key")
// resp.Version
```

## `db.Cipher` — Encryption-as-a-service

| Method | Args | Returns | Description |
|--------|------|---------|-------------|
| `Auth` | `ctx, token` | `*CipherAuthResponse, error` | Authenticate the connection |
| `CommandList` | `ctx` | `*CipherCommandListResponse, error` | List all supported commands |
| `Decrypt` | `ctx, keyring, ciphertext, opts` | `*CipherDecryptResponse, error` | Decrypt ciphertext using the embedded key version |
| `Encrypt` | `ctx, keyring, plaintext, opts` | `*CipherEncryptResponse, error` | Encrypt plaintext with the active key version |
| `GenerateDataKey` | `ctx, keyring, opts` | `*CipherGenerateDataKeyResponse, error` | Generate a data encryption key (envelope encryption pattern) |
| `Health` | `ctx` | `*CipherHealthResponse, error` | Check server health |
| `KeyInfo` | `ctx, keyring` | `*CipherKeyInfoResponse, error` | Get keyring metadata and key version information |
| `KeyringCreate` | `ctx, name, algorithm, opts` | `*CipherKeyringCreateResponse, error` | Create a new keyring with its first active key |
| `KeyringList` | `ctx` | `*CipherKeyringListResponse, error` | List all keyring names |
| `Ping` | `ctx` | `*CipherPingResponse, error` | Simple connectivity check — returns PONG |
| `Rewrap` | `ctx, keyring, ciphertext, opts` | `*CipherRewrapResponse, error` | Re-encrypt ciphertext with the current active key version |
| `Rotate` | `ctx, keyring, opts` | `*CipherRotateResponse, error` | Rotate the keyring to a new key version |
| `Sign` | `ctx, keyring, data` | `*CipherSignResponse, error` | Create a detached signature |
| `VerifySignature` | `ctx, keyring, data, signature` | `*CipherVerifySignatureResponse, error` | Verify a detached signature |

### Examples

```go
ctx := context.Background()
resp, err := db.Cipher.Decrypt(ctx, "my-keyring", "k3Xm:encrypted...")
// resp.Plaintext
resp, err := db.Cipher.Encrypt(ctx, "my-keyring", "SGVsbG8=")
// resp.Ciphertext
resp, err := db.Cipher.GenerateDataKey(ctx, "my-keyring")
// resp.PlaintextKey
```

## `db.Sigil` — Schema-driven credential envelope engine

| Method | Args | Returns | Description |
|--------|------|---------|-------------|
| `CredentialChange` | `ctx, schema, id, field, old, new` | `*SigilCredentialChangeResponse, error` | Change a credential field (requires old value for verification) |
| `CredentialImport` | `ctx, schema, id, field, hash, opts` | `*SigilCredentialImportResponse, error` | Import a pre-hashed credential (bcrypt, scrypt, argon2). Transparently rehashed to Argon2id on next verify. |
| `CredentialReset` | `ctx, schema, id, field, new` | `*SigilCredentialResetResponse, error` | Force-reset a credential field without requiring old value (admin/reset token) |
| `EnvelopeCreate` | `ctx, schema, id, json` | `*SigilEnvelopeCreateResponse, error` | Create an envelope with field routing per schema annotations |
| `EnvelopeDelete` | `ctx, schema, id` | `*SigilEnvelopeDeleteResponse, error` | Delete an envelope and all associated data |
| `EnvelopeGet` | `ctx, schema, id` | `*SigilEnvelopeGetResponse, error` | Get an envelope record |
| `EnvelopeImport` | `ctx, schema, id, json` | `*SigilEnvelopeImportResponse, error` | Import an envelope with pre-hashed credential fields. Non-credential fields processed normally. |
| `EnvelopeLookup` | `ctx, schema, field, value` | `*SigilEnvelopeLookupResponse, error` | Look up an envelope by indexed or searchable field value |
| `EnvelopeUpdate` | `ctx, schema, id, json` | `*SigilEnvelopeUpdateResponse, error` | Update non-credential fields on an existing envelope |
| `EnvelopeVerify` | `ctx, schema, id, field, value` | `*SigilEnvelopeVerifyResponse, error` | Verify a credential field on an envelope by explicit field name |
| `Health` | `ctx` | `*SigilHealthResponse, error` | Health check |
| `Jwks` | `ctx, schema` | `*SigilJwksResponse, error` | Get the JSON Web Key Set for external token verification |
| `PasswordChange` | `ctx, schema, id, old, new` | `*SigilPasswordChangeResponse, error` | Sugar: change password. Infers credential field from schema. Equivalent to CREDENTIAL CHANGE with implicit field. |
| `PasswordImport` | `ctx, schema, id, hash, opts` | `*SigilPasswordImportResponse, error` | Sugar: import pre-hashed password. Infers credential field from schema. Equivalent to CREDENTIAL IMPORT with implicit field. |
| `PasswordReset` | `ctx, schema, id, new` | `*SigilPasswordResetResponse, error` | Sugar: force-reset password. Infers credential field from schema. Equivalent to CREDENTIAL RESET with implicit field. |
| `SchemaAlter` | `ctx, name, action, opts` | `*SigilSchemaAlterResponse, error` | Add or remove fields from a schema, producing a new version. Added fields are optional (required=false). Existing envelopes remain readable. |
| `SchemaGet` | `ctx, name` | `*SigilSchemaGetResponse, error` | Get a schema definition by name |
| `SchemaList` | `ctx` | `*SigilSchemaListResponse, error` | List all registered schema names |
| `SchemaRegister` | `ctx, name, json` | `*SigilSchemaRegisterResponse, error` | Register a credential envelope schema |
| `SessionCreate` | `ctx, schema, id, password, opts` | `*SigilSessionCreateResponse, error` | Verify credentials and issue access + refresh tokens |
| `SessionList` | `ctx, schema, id` | `*SigilSessionListResponse, error` | List active sessions for an entity |
| `SessionRefresh` | `ctx, schema, token` | `*SigilSessionRefreshResponse, error` | Rotate refresh token and issue new access token |
| `SessionRevoke` | `ctx, schema, token` | `*SigilSessionRevokeResponse, error` | Revoke a single refresh token (logout one session) |
| `SessionRevokeAll` | `ctx, schema, id` | `*SigilSessionRevokeAllResponse, error` | Revoke all sessions for an entity (logout everywhere) |
| `UserCreate` | `ctx, schema, id, json` | `*SigilUserCreateResponse, error` | Sugar: create an envelope. Equivalent to ENVELOPE CREATE. |
| `UserDelete` | `ctx, schema, id` | `*SigilUserDeleteResponse, error` | Sugar: delete an envelope. Equivalent to ENVELOPE DELETE. |
| `UserGet` | `ctx, schema, id` | `*SigilUserGetResponse, error` | Sugar: get an envelope. Equivalent to ENVELOPE GET. |
| `UserImport` | `ctx, schema, id, json` | `*SigilUserImportResponse, error` | Sugar: import an envelope with pre-hashed credentials. Equivalent to ENVELOPE IMPORT. |
| `UserUpdate` | `ctx, schema, id, json` | `*SigilUserUpdateResponse, error` | Sugar: update non-credential fields. Equivalent to ENVELOPE UPDATE. |
| `UserVerify` | `ctx, schema, id, password` | `*SigilUserVerifyResponse, error` | Sugar: verify credential. Infers the credential field from schema. Equivalent to ENVELOPE VERIFY with implicit field. |

### Examples

```go
ctx := context.Background()
resp, err := db.Sigil.CredentialChange(ctx, "myapp", "alice", "email", "old", "new")
// resp.Status
resp, err := db.Sigil.CredentialImport(ctx, "myapp", "alice", "email", "hash")
// resp.Algorithm
resp, err := db.Sigil.CredentialReset(ctx, "myapp", "alice", "email", "new")
// resp.Status
```

## `db.Veil` — veil

| Method | Args | Returns | Description |
|--------|------|---------|-------------|
| `Auth` | `ctx, token` | `*VeilAuthResponse, error` | Authenticate this connection |
| `CommandList` | `ctx` | `*VeilCommandListResponse, error` | List all supported commands |
| `Delete` | `ctx, index, id` | `*VeilDeleteResponse, error` | Remove an entry's blind tokens from the index |
| `Health` | `ctx` | `*VeilHealthResponse, error` | Health check |
| `IndexCreate` | `ctx, name` | `*VeilIndexCreateResponse, error` | Create a new blind index with a fresh HMAC key |
| `IndexDestroy` | `ctx, name` | `*VeilIndexDestroyResponse, error` | Crypto-shred an index: zeroize the HMAC key, delete all entries, and remove the index. After destruction, the index name can be reused. |
| `IndexInfo` | `ctx, name` | `*VeilIndexInfoResponse, error` | Get information about a blind index |
| `IndexList` | `ctx` | `*VeilIndexListResponse, error` | List all blind index names |
| `IndexReconcile` | `ctx, name, valid_ids` | `*VeilIndexReconcileResponse, error` | Remove orphaned entries from the index. Compares stored entry IDs against the provided valid set and deletes any entries not in the set. |
| `IndexReindex` | `ctx, name` | `*VeilIndexReindexResponse, error` | Clear all entries and update the tokenizer version to current. The HMAC key is preserved. After reindex, the application must re-submit all entries via PUT. Use this when the tokenizer algorithm has been upgraded. |
| `IndexRotate` | `ctx, name` | `*VeilIndexRotateResponse, error` | Rotate an index's HMAC key. Generates a new key, deletes all existing entries. The application must re-index all entries after rotation. |
| `Ping` | `ctx` | `*VeilPingResponse, error` | Ping-pong |
| `Put` | `ctx, index, id, data_b64, opts` | `*VeilPutResponse, error` | Store blind tokens for an entry. In standard mode, data_b64 is base64-encoded plaintext (server tokenizes). With BLIND flag, data_b64 is base64-encoded BlindTokenSet JSON (client pre-tokenized, for E2EE). |
| `Search` | `ctx, index, query, opts` | `*VeilSearchResponse, error` | Search a blind index. In standard mode, query is plain text (server tokenizes). With BLIND flag, query is base64-encoded BlindTokenSet JSON (client pre-tokenized, for E2EE). |
| `Tokenize` | `ctx, index, plaintext_b64, opts` | `*VeilTokenizeResponse, error` | Generate blind tokens from plaintext without storing. Returns HMAC-derived tokens for external use. |

### Examples

```go
ctx := context.Background()
resp, err := db.Veil.Delete(ctx, "index", "alice")
// resp.Status
resp, err := db.Veil.IndexCreate(ctx, "my-keyring")
// resp.Status
resp, err := db.Veil.IndexDestroy(ctx, "my-keyring")
// resp.Status
```

## `db.Sentry` — sentry

| Method | Args | Returns | Description |
|--------|------|---------|-------------|
| `Auth` | `ctx, token` | `*SentryAuthResponse, error` | Authenticate the connection with a token |
| `CommandList` | `ctx` | `*SentryCommandListResponse, error` | List all supported commands |
| `Evaluate` | `ctx, json` | `*SentryEvaluateResponse, error` | Evaluate an authorization request against policies and return a signed decision |
| `Health` | `ctx` | `*SentryHealthResponse, error` | Server health check |
| `Jwks` | `ctx` | `*SentryJwksResponse, error` | Get the JSON Web Key Set for verifying decision tokens |
| `KeyInfo` | `ctx` | `*SentryKeyInfoResponse, error` | Get signing key metadata |
| `KeyRotate` | `ctx, opts` | `*SentryKeyRotateResponse, error` | Rotate the signing key |
| `Ping` | `ctx` | `error` | Connectivity check |
| `PolicyCreate` | `ctx, name, json` | `*SentryPolicyCreateResponse, error` | Create a new authorization policy |
| `PolicyDelete` | `ctx, name` | `*SentryPolicyDeleteResponse, error` | Delete a policy |
| `PolicyGet` | `ctx, name` | `*SentryPolicyGetResponse, error` | Get a policy by name |
| `PolicyHistory` | `ctx, name` | `*SentryPolicyHistoryResponse, error` | Get version history of a policy (all past versions plus current) |
| `PolicyList` | `ctx` | `*SentryPolicyListResponse, error` | List all policy names |
| `PolicyUpdate` | `ctx, name, json` | `*SentryPolicyUpdateResponse, error` | Update an existing policy |

### Examples

```go
ctx := context.Background()
resp, err := db.Sentry.Evaluate(ctx, "json")
// resp.CacheUntil
resp, err := db.Sentry.PolicyCreate(ctx, "name", "json")
// resp.Effect
resp, err := db.Sentry.PolicyDelete(ctx, "name")
// resp.Status
```

## `db.Forge` — Internal certificate authority engine

| Method | Args | Returns | Description |
|--------|------|---------|-------------|
| `CaCreate` | `ctx, name, algorithm, subject, opts` | `*ForgeCaCreateResponse, error` | Create a new Certificate Authority |
| `CaExport` | `ctx, name` | `*ForgeCaExportResponse, error` | Export the active CA certificate (PEM) |
| `CaInfo` | `ctx, name` | `*ForgeCaInfoResponse, error` | Get CA metadata and key version status |
| `CaList` | `ctx` | `*ForgeCaListResponse, error` | List all Certificate Authorities |
| `CaRotate` | `ctx, name, opts` | `*ForgeCaRotateResponse, error` | Rotate CA signing key |
| `ConfigGet` | `ctx, key` | `*ForgeConfigGetResponse, error` | Get a runtime configuration value |
| `ConfigSet` | `ctx, key, value` | `*ForgeConfigSetResponse, error` | Set a runtime configuration value (only scheduler_interval_secs is mutable) |
| `Inspect` | `ctx, ca, serial` | `*ForgeInspectResponse, error` | Get certificate details |
| `Issue` | `ctx, ca, subject, profile, opts` | `*ForgeIssueResponse, error` | Issue a new certificate. Returns cert + private key (private key never stored). |
| `IssueFromCsr` | `ctx, ca, csr_pem, profile, opts` | `*ForgeIssueFromCsrResponse, error` | Issue a certificate from a PEM-encoded CSR |
| `ListCerts` | `ctx, ca, opts` | `*ForgeListCertsResponse, error` | List certificates for a CA |
| `Renew` | `ctx, ca, serial, opts` | `*ForgeRenewResponse, error` | Renew a certificate (re-issue with same profile and SANs) |
| `Revoke` | `ctx, ca, serial, opts` | `*ForgeRevokeResponse, error` | Revoke a certificate |

### Examples

```go
ctx := context.Background()
resp, err := db.Forge.CaCreate(ctx, "name", "algorithm", "subject")
// resp.ActiveVersion
resp, err := db.Forge.CaExport(ctx, "name")
// resp.CertificatePem
resp, err := db.Forge.CaInfo(ctx, "name")
// resp.Algorithm
```

## `db.Keep` — Secrets manager with path-based access control and versioning

| Method | Args | Returns | Description |
|--------|------|---------|-------------|
| `Auth` | `ctx, token` | `*KeepAuthResponse, error` | Authenticate this connection with a token. |
| `CommandList` | `ctx` | `*KeepCommandListResponse, error` | List all supported commands. |
| `Delete` | `ctx, path` | `*KeepDeleteResponse, error` | Soft-delete a secret. Version history is preserved. |
| `Get` | `ctx, path, opts` | `*KeepGetResponse, error` | Retrieve a secret value. Returns the latest version by default. |
| `Health` | `ctx` | `*KeepHealthResponse, error` | Health check. |
| `List` | `ctx, prefix?` | `*KeepListResponse, error` | List secret paths, optionally filtered by prefix. Excludes deleted secrets. |
| `Ping` | `ctx` | `error` | Ping-pong. |
| `Purge` | `ctx, path` | `*KeepPurgeResponse, error` | Permanently remove a secret and all its versions. Irreversible — used for GDPR right-to-erasure compliance. After purge, GET returns not-found (not deleted). |
| `Put` | `ctx, path, value` | `*KeepPutResponse, error` | Store a new version of a secret. Creates the secret if it doesn't exist. Undeletes if soft-deleted. |
| `Rekey` | `ctx, new_key` | `*KeepRekeyResponse, error` | Re-encrypt all secrets with a new master key. Iterates all secrets (including deleted ones), decrypts every version with the current master key, re-encrypts with the new key, and switches to the new key for all future operations. |
| `Rotate` | `ctx, path` | `*KeepRotateResponse, error` | Re-encrypt the latest version with a new nonce. Creates a new version with the same plaintext. |
| `Versions` | `ctx, path` | `*KeepVersionsResponse, error` | Get version history for a secret. Includes deleted secrets. |

### Examples

```go
ctx := context.Background()
resp, err := db.Keep.Delete(ctx, "path")
// resp.Status
resp, err := db.Keep.Get(ctx, "path")
// resp.Status
resp, err := db.Keep.List(ctx, "prefix")
// resp.Status
```

## `db.Courier` — Just-in-time decryption delivery engine

| Method | Args | Returns | Description |
|--------|------|---------|-------------|
| `Auth` | `ctx, token` | `*CourierAuthResponse, error` | Authenticate the connection with a token |
| `ChannelCreate` | `ctx, name, type, config_json` | `*CourierChannelCreateResponse, error` | Create a delivery channel |
| `ChannelDelete` | `ctx, name` | `*CourierChannelDeleteResponse, error` | Delete a channel |
| `ChannelGet` | `ctx, name` | `*CourierChannelGetResponse, error` | Get channel configuration |
| `ChannelList` | `ctx` | `*CourierChannelListResponse, error` | List all channels |
| `CommandList` | `ctx` | `*CourierCommandListResponse, error` | List available commands |
| `Deliver` | `ctx, json` | `*CourierDeliverResponse, error` | Decrypt recipient and deliver a message |
| `DeliveryGet` | `ctx, id` | `*CourierDeliveryGetResponse, error` | Get a delivery receipt by ID |
| `DeliveryList` | `ctx, opts` | `*CourierDeliveryListResponse, error` | List delivery receipts, optionally filtered by channel |
| `Health` | `ctx` | `*CourierHealthResponse, error` | Server health check |
| `Metrics` | `ctx` | `*CourierMetricsResponse, error` | Get delivery metrics (total, success, failure counts, per-channel breakdown) |
| `NotifyEvent` | `ctx, channel, subject, body` | `*CourierNotifyEventResponse, error` | Trigger a notification on a pre-configured channel (e.g. rotation/expiry alerts) |
| `Ping` | `ctx` | `error` | Connectivity check |

### Examples

```go
ctx := context.Background()
resp, err := db.Courier.ChannelCreate(ctx, "name", "type", "config_json")
// resp.ChannelType
resp, err := db.Courier.ChannelDelete(ctx, "name")
// resp.Name
resp, err := db.Courier.ChannelGet(ctx, "name")
// resp.ChannelType
```

## `db.Chronicle` — Structured audit event engine

| Method | Args | Returns | Description |
|--------|------|---------|-------------|
| `Actors` | `ctx, opts` | `*ChronicleActorsResponse, error` | Active actors in time window |
| `Auth` | `ctx, token` | `*ChronicleAuthResponse, error` | Authenticate this connection |
| `Count` | `ctx, opts` | `*ChronicleCountResponse, error` | Count events matching filter predicates |
| `Errors` | `ctx, opts` | `*ChronicleErrorsResponse, error` | Error rates by action |
| `Health` | `ctx` | `*ChronicleHealthResponse, error` | Health check |
| `Hotspots` | `ctx, opts` | `*ChronicleHotspotsResponse, error` | Top actors by event volume |
| `Ingest` | `ctx, event_json` | `*ChronicleIngestResponse, error` | Ingest a single structured audit event |
| `IngestBatch` | `ctx, events_json` | `*ChronicleIngestBatchResponse, error` | Ingest multiple events in a single call |
| `Ping` | `ctx` | `error` | Keepalive |
| `Query` | `ctx, opts` | `*ChronicleQueryResponse, error` | Query events with filter predicates |
| `Verify` | `ctx` | `*ChronicleVerifyResponse, error` | Verify the cryptographic hash chain integrity of all events. Returns the number of verified events or an error if tampering is detected. |

### Examples

```go
ctx := context.Background()
resp, err := db.Chronicle.Ingest(ctx, map[string]any{})
// resp.Status
resp, err := db.Chronicle.IngestBatch(ctx, map[string]any{})
// resp.Ingested
```

## `db.Stash` — Encrypted blob storage with S3 backend and envelope encryption

| Method | Args | Returns | Description |
|--------|------|---------|-------------|
| `Command` | `ctx` | `error` | List supported commands |
| `Fingerprint` | `ctx, id, viewer_id, opts` | `*StashFingerprintResponse, error` | Create a viewer-specific encrypted copy of a blob for leak tracing |
| `Health` | `ctx` | `error` | Health check |
| `Inspect` | `ctx, id` | `*StashInspectResponse, error` | Read blob metadata without downloading or decrypting |
| `List` | `ctx, opts` | `*StashListResponse, error` | List blobs for the current tenant |
| `Ping` | `ctx` | `error` | Ping-pong |
| `Retrieve` | `ctx, id` | `error` | Retrieve and decrypt a blob |
| `Revoke` | `ctx, id, opts` | `*StashRevokeResponse, error` | Revoke a blob (hard crypto-shred by default, SOFT for soft revoke) |
| `Rewrap` | `ctx, id` | `*StashRewrapResponse, error` | Re-wrap a blob's DEK under the current Cipher key version. The blob ciphertext is not re-encrypted — only the key wrapping changes. |
| `Store` | `ctx, id, data_b64, opts` | `*StashStoreResponse, error` | Store an encrypted blob |
| `Trace` | `ctx, id` | `*StashTraceResponse, error` | Return the viewer map (who has copies) for a blob |

### Examples

```go
ctx := context.Background()
resp, err := db.Stash.Fingerprint(ctx, "alice", "viewer_id")
// resp.CreatedAt
resp, err := db.Stash.Inspect(ctx, "alice")
// resp.BlobStatus
err := db.Stash.Retrieve(ctx, "alice")
```

## Error Handling

All methods return `error` (or `(*Response, error)`). Errors from the server are `*ShrouDBError` with a `Code` field matching the server error code (e.g., `NOTFOUND`, `DENIED`, `BADARG`).

```go
result, err := db.Cipher.Encrypt(ctx, "kr", data)
if err != nil {
	if shroudb.IsCode(err, shroudb.ErrNOTFOUND) {
		// handle not found
	}
}
```

## Error Codes

| Code | Constant | Description |
|------|----------|-------------|
| `BAD_ARG` | `ErrBAD_ARG` | Missing or malformed command argument |
| `DENIED` | `ErrDENIED` | Authentication required or insufficient permissions |
| `NAMESPACE_EXISTS` | `ErrNAMESPACE_EXISTS` | Namespace already exists |
| `NAMESPACE_NOT_EMPTY` | `ErrNAMESPACE_NOT_EMPTY` | Namespace is not empty (use FORCE to override) |
| `NAMESPACE_NOT_FOUND` | `ErrNAMESPACE_NOT_FOUND` | Namespace does not exist |
| `NOT_AUTHENTICATED` | `ErrNOT_AUTHENTICATED` | No auth token provided on this connection |
| `NOT_FOUND` | `ErrNOT_FOUND` | Key or resource does not exist |
| `NOT_READY` | `ErrNOT_READY` | Server is not in READY state |
| `PIPELINE_ABORTED` | `ErrPIPELINE_ABORTED` | Pipeline command failed, all commands rolled back |
| `VALIDATION_FAILED` | `ErrVALIDATION_FAILED` | Metadata validation failed against namespace schema |
| `VERSION_NOT_FOUND` | `ErrVERSION_NOT_FOUND` | Requested version does not exist |
| `BADARG` | `ErrBADARG` | Missing or invalid argument |
| `DISABLED` | `ErrDISABLED` | Keyring is disabled |
| `EXISTS` | `ErrEXISTS` | Keyring already exists |
| `INTERNAL` | `ErrINTERNAL` | Unexpected server error |
| `NOTFOUND` | `ErrNOTFOUND` | Keyring or key version not found |
| `POLICY` | `ErrPOLICY` | Operation denied by keyring policy |
| `RETIRED` | `ErrRETIRED` | Key version is retired — use REWRAP |
| `WRONGTYPE` | `ErrWRONGTYPE` | Operation not supported for this keyring type |
| `ACCOUNT_LOCKED` | `ErrACCOUNT_LOCKED` | Account locked after too many failed attempts |
| `CAPABILITY_MISSING` | `ErrCAPABILITY_MISSING` | Required engine capability not available (e.g., Cipher for PII fields) |
| `ENTITY_EXISTS` | `ErrENTITY_EXISTS` | Entity already exists |
| `ENTITY_NOT_FOUND` | `ErrENTITY_NOT_FOUND` | Entity does not exist |
| `IMPORT_FAILED` | `ErrIMPORT_FAILED` | Password import failed (invalid hash format) |
| `INVALID_FIELD` | `ErrINVALID_FIELD` | Field value is invalid or field cannot be updated via this path |
| `INVALID_TOKEN` | `ErrINVALID_TOKEN` | Token is invalid, expired, or revoked |
| `MISSING_FIELD` | `ErrMISSING_FIELD` | Required field missing from request |
| `SCHEMA_EXISTS` | `ErrSCHEMA_EXISTS` | Schema already exists |
| `SCHEMA_NOT_FOUND` | `ErrSCHEMA_NOT_FOUND` | Schema does not exist |
| `SCHEMA_VALIDATION` | `ErrSCHEMA_VALIDATION` | Schema definition is invalid |
| `TOKEN_REUSE` | `ErrTOKEN_REUSE` | Refresh token reuse detected — entire family revoked |
| `VERIFICATION_FAILED` | `ErrVERIFICATION_FAILED` | Credential verification failed (wrong password) |
| `AUTH_REQUIRED` | `ErrAUTH_REQUIRED` | Authentication required |
| `NOKEY` | `ErrNOKEY` | No active signing key available |
| `SIGNING` | `ErrSIGNING` | Failed to sign decision |
| `STORAGE` | `ErrSTORAGE` | Backend storage error |
| `DELETED` | `ErrDELETED` | Secret has been soft-deleted |
| `ENCRYPTION` | `ErrENCRYPTION` | Encryption or decryption failed |
| `VERSION_NOTFOUND` | `ErrVERSION_NOTFOUND` | Requested version does not exist |
| `ADAPTER` | `ErrADAPTER` | Delivery adapter failure |
| `DECRYPT` | `ErrDECRYPT` | Cipher decryption failed |
| `CIPHER_UNAVAILABLE` | `ErrCIPHER_UNAVAILABLE` | Cipher engine not available for envelope encryption |
| `CLIENT_ENCRYPTED` | `ErrCLIENT_ENCRYPTED` | Cannot fingerprint a client-encrypted blob (client manages encryption) |
| `CRYPTO` | `ErrCRYPTO` | Encryption or decryption failed |
| `DUPLICATE_VIEWER` | `ErrDUPLICATE_VIEWER` | Viewer already has a fingerprinted copy of this blob |
| `INVALID_ARGUMENT` | `ErrINVALID_ARGUMENT` | Invalid argument |
| `OBJECT_STORE` | `ErrOBJECT_STORE` | S3 object store operation failed |
| `REVOKED` | `ErrREVOKED` | Blob has been soft-revoked |
| `SHREDDED` | `ErrSHREDDED` | Blob has been crypto-shredded (unrecoverable) |

## Common Mistakes

- Always `defer db.Close()` to release connection pool resources
- Every method requires a `context.Context` as the first argument
- Engine methods handle serialization — pass Go maps for JSON params, not `json.Marshal()`
- Accessing a nil engine namespace panics — check your `Options` configuration
- Optional keyword params use pointer fields in the Options struct — use `&value` to set them
