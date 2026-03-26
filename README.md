# ShrouDB Go SDK

Typed Go clients for all ShrouDB engines. Auto-generated from protocol specs.

Each engine is a separate Go module within this repository.

## Install

```bash
go get github.com/shroudb/shroudb-sdk-go/shroudb@latest
go get github.com/shroudb/shroudb-sdk-go/shroudb-transit@latest
go get github.com/shroudb/shroudb-sdk-go/shroudb-keep@latest
```

## Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/shroudb/shroudb-sdk-go/shroudb"
    "github.com/shroudb/shroudb-sdk-go/shroudb-transit"
)

func main() {
    ctx := context.Background()

    // Vault — credential management
    vault, _ := shroudb.Connect("shroudb://localhost")
    defer vault.Close()
    cred, _ := vault.Issue(ctx, "my-keyspace", &shroudb.IssueOptions{TTL: 3600})
    fmt.Println(cred.CredentialID)

    // Transit — encryption-as-a-service
    transit, _ := shroudb_transit.Connect("shroudb-transit://localhost")
    defer transit.Close()
    encrypted, _ := transit.Encrypt(ctx, "my-key", plaintext)
    fmt.Println(encrypted)
}
```

## Engines

| Module | Package | Description |
|--------|---------|-------------|
| `shroudb` | `shroudb` | Credential management |
| `shroudb-transit` | `shroudb_transit` | Encryption-as-a-service |
| `shroudb-auth` | `shroudb_auth` | Authentication |
| `shroudb-mint` | `shroudb_mint` | Token minting |
| `shroudb-sentry` | `shroudb_sentry` | Access control |
| `shroudb-keep` | `shroudb_keep` | Secret storage |
| `shroudb-courier` | `shroudb_courier` | Secure delivery |
| `shroudb-pulse` | `shroudb_pulse` | Telemetry & metrics |

## License

MIT OR Apache-2.0
