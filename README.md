# github.com/shroudb/shroudb-go

Unified Go SDK for all ShrouDB engines. Provides namespaced, type-safe
access to every engine with built-in serialization, connection pooling, and
dual transport support (RESP3 for direct connections, HTTP for Moat gateway).

## Installation

```bash
go get github.com/shroudb/shroudb-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	shroudb "github.com/shroudb/shroudb-go"
)

func main() {
	ctx := context.Background()

	// Connect via Moat gateway (routes all engines through one endpoint)
	db, err := shroudb.New(shroudb.Options{
		Moat:  "https://moat.example.com",
		Token: "my-token",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Encrypt data
	result, err := db.Cipher.Encrypt(ctx, "my-keyring", "SGVsbG8=")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Ciphertext)
}
```

### Direct Engine Connections

```go
db, err := shroudb.New(shroudb.Options{
	Shroudb: "shroudb://token@localhost:6399",
	Cipher: "shroudb-cipher://token@localhost:6599",
	Sigil: "sigil://token@localhost:6499",
	Veil: "shroudb-veil://token@localhost:6799",
	Sentry: "shroudb-sentry://token@localhost:6799",
	Forge: "shroudb-forge://token@localhost:6699",
	Keep: "shroudb-keep://token@localhost:6899",
	Courier: "shroudb-courier://token@localhost:6999",
	Chronicle: "chronicle://token@localhost:7099",
})
```

## Connection Modes

### Moat Gateway (HTTP)

Routes all engine commands through a single Moat endpoint via HTTP POST.

```go
db, _ := shroudb.New(shroudb.Options{Moat: "https://moat.example.com", Token: "my-token"})
```

### Moat Gateway (RESP3)

Direct RESP3 connection to Moat with engine-prefixed commands.

```go
db, _ := shroudb.New(shroudb.Options{Moat: "shroudb-moat://my-token@moat.example.com:8201"})
```

### Direct Engine Connections

Connect to individual engines via RESP3. Only configure the engines you need.

```go
db, _ := shroudb.New(shroudb.Options{
	Cipher: "shroudb-cipher://token@cipher-host:6599",
	Sigil:  "shroudb-sigil://token@sigil-host:6499",
})
```

### Mixed Mode

Route most engines through Moat, but connect directly to specific engines.

```go
db, _ := shroudb.New(shroudb.Options{
	Moat:   "https://moat.example.com",
	Cipher: "shroudb-cipher://token@dedicated-cipher:6599", // direct
	Token:  "moat-token",
})
```

## Engines

### `db.Shroudb`

Encrypted key-value database

| Method | Description |
|--------|-------------|
| `Auth(ctx, token)` | Authenticate the connection with a token |
| `CommandList(ctx)` | List all supported commands |
| `ConfigGet(ctx, key)` | Read a runtime configuration value |
| `ConfigSet(ctx, key, value)` | Set a runtime configuration value (admin only) |
| `Delete(ctx, namespace, key)` | Delete a key by writing a tombstone |
| `Get(ctx, namespace, key, META, opts)` | Retrieve the value at a key |
| `Health(ctx)` | Check server health |
| `List(ctx, namespace, opts)` | List active keys in a namespace |
| `NamespaceAlter(ctx, name, opts)` | Update namespace configuration (enforce-on-write-only) |
| `NamespaceCreate(ctx, name, opts)` | Create a new namespace |
| `NamespaceDrop(ctx, name, FORCE)` | Drop a namespace |
| `NamespaceInfo(ctx, name)` | Get metadata about a namespace |
| `NamespaceList(ctx, opts)` | List namespaces (filtered by token grants) |
| `NamespaceValidate(ctx, name)` | Check existing entries against current MetaSchema |
| `Ping(ctx)` | Test connectivity |
| `Pipeline(ctx, count)` | Execute commands atomically (all succeed or all roll back) |
| `Put(ctx, namespace, key, value, opts)` | Store a value at the given key. Auto-increments version. |
| `Subscribe(ctx, namespace, opts)` | Subscribe to change events on a namespace |
| `Unsubscribe(ctx)` | End the current subscription |
| `Versions(ctx, namespace, key, opts)` | Retrieve version history for a key (most recent first) |

### `db.Cipher`

Encryption-as-a-service

| Method | Description |
|--------|-------------|
| `Auth(ctx, token)` | Authenticate the connection |
| `CommandList(ctx)` | List all supported commands |
| `Decrypt(ctx, keyring, ciphertext, opts)` | Decrypt ciphertext using the embedded key version |
| `Encrypt(ctx, keyring, plaintext, opts)` | Encrypt plaintext with the active key version |
| `GenerateDataKey(ctx, keyring, opts)` | Generate a data encryption key (envelope encryption pattern) |
| `Health(ctx)` | Check server health |
| `KeyInfo(ctx, keyring)` | Get keyring metadata and key version information |
| `KeyringCreate(ctx, name, algorithm, opts)` | Create a new keyring with its first active key |
| `KeyringList(ctx)` | List all keyring names |
| `Ping(ctx)` | Simple connectivity check — returns PONG |
| `Rewrap(ctx, keyring, ciphertext, opts)` | Re-encrypt ciphertext with the current active key version |
| `Rotate(ctx, keyring, opts)` | Rotate the keyring to a new key version |
| `Sign(ctx, keyring, data)` | Create a detached signature |
| `VerifySignature(ctx, keyring, data, signature)` | Verify a detached signature |

### `db.Sigil`

Schema-driven credential envelope engine

| Method | Description |
|--------|-------------|
| `CredentialChange(ctx, schema, id, field, old, new)` | Change a credential field (requires old value for verification) |
| `CredentialImport(ctx, schema, id, field, hash, opts)` | Import a pre-hashed credential (bcrypt, scrypt, argon2). Transparently rehashed to Argon2id on next verify. |
| `CredentialReset(ctx, schema, id, field, new)` | Force-reset a credential field without requiring old value (admin/reset token) |
| `EnvelopeCreate(ctx, schema, id, json)` | Create an envelope with field routing per schema annotations |
| `EnvelopeDelete(ctx, schema, id)` | Delete an envelope and all associated data |
| `EnvelopeGet(ctx, schema, id)` | Get an envelope record |
| `EnvelopeImport(ctx, schema, id, json)` | Import an envelope with pre-hashed credential fields. Non-credential fields processed normally. |
| `EnvelopeLookup(ctx, schema, field, value)` | Look up an envelope by indexed or searchable field value |
| `EnvelopeUpdate(ctx, schema, id, json)` | Update non-credential fields on an existing envelope |
| `EnvelopeVerify(ctx, schema, id, field, value)` | Verify a credential field on an envelope by explicit field name |
| `Health(ctx)` | Health check |
| `Jwks(ctx, schema)` | Get the JSON Web Key Set for external token verification |
| `PasswordChange(ctx, schema, id, old, new)` | Sugar: change password. Infers credential field from schema. Equivalent to CREDENTIAL CHANGE with implicit field. |
| `PasswordImport(ctx, schema, id, hash, opts)` | Sugar: import pre-hashed password. Infers credential field from schema. Equivalent to CREDENTIAL IMPORT with implicit field. |
| `PasswordReset(ctx, schema, id, new)` | Sugar: force-reset password. Infers credential field from schema. Equivalent to CREDENTIAL RESET with implicit field. |
| `SchemaGet(ctx, name)` | Get a schema definition by name |
| `SchemaList(ctx)` | List all registered schema names |
| `SchemaRegister(ctx, name, json)` | Register a credential envelope schema |
| `SessionCreate(ctx, schema, id, password, opts)` | Verify credentials and issue access + refresh tokens |
| `SessionList(ctx, schema, id)` | List active sessions for an entity |
| `SessionRefresh(ctx, schema, token)` | Rotate refresh token and issue new access token |
| `SessionRevoke(ctx, schema, token)` | Revoke a single refresh token (logout one session) |
| `SessionRevokeAll(ctx, schema, id)` | Revoke all sessions for an entity (logout everywhere) |
| `UserCreate(ctx, schema, id, json)` | Sugar: create an envelope. Equivalent to ENVELOPE CREATE. |
| `UserDelete(ctx, schema, id)` | Sugar: delete an envelope. Equivalent to ENVELOPE DELETE. |
| `UserGet(ctx, schema, id)` | Sugar: get an envelope. Equivalent to ENVELOPE GET. |
| `UserImport(ctx, schema, id, json)` | Sugar: import an envelope with pre-hashed credentials. Equivalent to ENVELOPE IMPORT. |
| `UserUpdate(ctx, schema, id, json)` | Sugar: update non-credential fields. Equivalent to ENVELOPE UPDATE. |
| `UserVerify(ctx, schema, id, password)` | Sugar: verify credential. Infers the credential field from schema. Equivalent to ENVELOPE VERIFY with implicit field. |

### `db.Veil`

veil

| Method | Description |
|--------|-------------|
| `Auth(ctx, token)` | Authenticate this connection |
| `CommandList(ctx)` | List all supported commands |
| `Delete(ctx, index, id)` | Remove an entry's blind tokens from the index |
| `Health(ctx)` | Health check |
| `IndexCreate(ctx, name)` | Create a new blind index with a fresh HMAC key |
| `IndexInfo(ctx, name)` | Get information about a blind index |
| `IndexList(ctx)` | List all blind index names |
| `Ping(ctx)` | Ping-pong |
| `Put(ctx, index, id, plaintext_b64, opts)` | Tokenize plaintext and store the blind tokens under the given entry ID |
| `Search(ctx, index, query, opts)` | Search a blind index. Tokenizes the query, generates blind tokens, and compares against stored entries. |
| `Tokenize(ctx, index, plaintext_b64, opts)` | Generate blind tokens from plaintext without storing. Returns HMAC-derived tokens for external use. |

### `db.Sentry`

sentry

| Method | Description |
|--------|-------------|
| `Auth(ctx, token)` | Authenticate the connection with a token |
| `CommandList(ctx)` | List all supported commands |
| `Evaluate(ctx, json)` | Evaluate an authorization request against policies and return a signed decision |
| `Health(ctx)` | Server health check |
| `Jwks(ctx)` | Get the JSON Web Key Set for verifying decision tokens |
| `KeyInfo(ctx)` | Get signing key metadata |
| `KeyRotate(ctx, opts)` | Rotate the signing key |
| `Ping(ctx)` | Connectivity check |
| `PolicyCreate(ctx, name, json)` | Create a new authorization policy |
| `PolicyDelete(ctx, name)` | Delete a policy |
| `PolicyGet(ctx, name)` | Get a policy by name |
| `PolicyList(ctx)` | List all policy names |
| `PolicyUpdate(ctx, name, json)` | Update an existing policy |

### `db.Forge`

Internal certificate authority engine

| Method | Description |
|--------|-------------|
| `CaCreate(ctx, name, algorithm, subject, opts)` | Create a new Certificate Authority |
| `CaExport(ctx, name)` | Export the active CA certificate (PEM) |
| `CaInfo(ctx, name)` | Get CA metadata and key version status |
| `CaList(ctx)` | List all Certificate Authorities |
| `CaRotate(ctx, name, opts)` | Rotate CA signing key |
| `Inspect(ctx, ca, serial)` | Get certificate details |
| `Issue(ctx, ca, subject, profile, opts)` | Issue a new certificate. Returns cert + private key (private key never stored). |
| `IssueFromCsr(ctx, ca, csr_pem, profile, opts)` | Issue a certificate from a PEM-encoded CSR |
| `ListCerts(ctx, ca, opts)` | List certificates for a CA |
| `Renew(ctx, ca, serial, opts)` | Renew a certificate (re-issue with same profile and SANs) |
| `Revoke(ctx, ca, serial, opts)` | Revoke a certificate |

### `db.Keep`

Secrets manager with path-based access control and versioning

| Method | Description |
|--------|-------------|
| `Auth(ctx, token)` | Authenticate this connection with a token. |
| `CommandList(ctx)` | List all supported commands. |
| `Delete(ctx, path)` | Soft-delete a secret. Version history is preserved. |
| `Get(ctx, path, opts)` | Retrieve a secret value. Returns the latest version by default. |
| `Health(ctx)` | Health check. |
| `List(ctx, prefix)` | List secret paths, optionally filtered by prefix. Excludes deleted secrets. |
| `Ping(ctx)` | Ping-pong. |
| `Put(ctx, path, value)` | Store a new version of a secret. Creates the secret if it doesn't exist. Undeletes if soft-deleted. |
| `Rotate(ctx, path)` | Re-encrypt the latest version with a new nonce. Creates a new version with the same plaintext. |
| `Versions(ctx, path)` | Get version history for a secret. Includes deleted secrets. |

### `db.Courier`

Just-in-time decryption delivery engine

| Method | Description |
|--------|-------------|
| `Auth(ctx, token)` | Authenticate the connection with a token |
| `ChannelCreate(ctx, name, type, config_json)` | Create a delivery channel |
| `ChannelDelete(ctx, name)` | Delete a channel |
| `ChannelGet(ctx, name)` | Get channel configuration |
| `ChannelList(ctx)` | List all channels |
| `CommandList(ctx)` | List available commands |
| `Deliver(ctx, json)` | Decrypt recipient and deliver a message |
| `Health(ctx)` | Server health check |
| `Ping(ctx)` | Connectivity check |

### `db.Chronicle`

Structured audit event engine

| Method | Description |
|--------|-------------|
| `Actors(ctx, opts)` | Active actors in time window |
| `Auth(ctx, token)` | Authenticate this connection |
| `Count(ctx, opts)` | Count events matching filter predicates |
| `Errors(ctx, opts)` | Error rates by action |
| `Health(ctx)` | Health check |
| `Hotspots(ctx, opts)` | Top actors by event volume |
| `Ingest(ctx, event_json)` | Ingest a single structured audit event |
| `IngestBatch(ctx, events_json)` | Ingest multiple events in a single call |
| `Ping(ctx)` | Keepalive |
| `Query(ctx, opts)` | Query events with filter predicates |

## Error Handling

```go
result, err := db.Cipher.Encrypt(ctx, "missing-keyring", data)
if err != nil {
	if shroudb.IsCode(err, shroudb.ErrNOTFOUND) {
		fmt.Println("Keyring not found")
	}
}
```
