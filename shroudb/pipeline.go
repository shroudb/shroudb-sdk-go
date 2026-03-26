// Shroudb pipeline for batching commands.
//
// Auto-generated from shroudb protocol spec. Do not edit.

package shroudb

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Pipeline batches multiple Shroudb commands and executes them in a single round-trip.
//
// Usage:
//
//	pipe := client.Pipeline()
//	pipe.Issue("keyspace", &IssueOptions{TTL: 3600})
//	pipe.Verify("keyspace", token, nil)
//	results, err := pipe.Execute()
type Pipeline struct {
	pool     *pool
	commands []pipelineCmd
}

type pipelineCmd struct {
	args   []string
	parser func(map[string]any) any // nil for simple responses
}

// Execute sends all queued commands and returns typed responses.
func (p *Pipeline) Execute() ([]any, error) {
	conn, err := p.pool.get()
	if err != nil {
		return nil, err
	}

	// Send all commands
	for _, cmd := range p.commands {
		if err := conn.sendCommand(cmd.args...); err != nil {
			conn.close()
			return nil, err
		}
	}

	// Read all responses
	results := make([]any, 0, len(p.commands))
	for _, cmd := range p.commands {
		raw, err := conn.readResponse()
		if err != nil {
			conn.close()
			return nil, err
		}
		if cmd.parser != nil {
			if m, ok := raw.(map[string]any); ok {
				results = append(results, cmd.parser(m))
			} else {
				results = append(results, raw)
			}
		} else {
			results = append(results, raw)
		}
	}

	p.pool.put(conn)
	p.commands = p.commands[:0]
	return results, nil
}

// Len returns the number of queued commands.
func (p *Pipeline) Len() int { return len(p.commands) }

// Clear discards all queued commands.
func (p *Pipeline) Clear() { p.commands = p.commands[:0] }

// Auth queues a AUTH command.
func (p *Pipeline) Auth(token string) *Pipeline {
	args := []string{
		"AUTH",
		token,
	}
	p.commands = append(p.commands, pipelineCmd{args: args})
	return p
}

// ConfigGet queues a CONFIG command.
func (p *Pipeline) ConfigGet(key string) *Pipeline {
	args := []string{
		"CONFIG", "GET",
		key,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseConfigGetResponse(m) }})
	return p
}

// ConfigSet queues a CONFIG command.
func (p *Pipeline) ConfigSet(key string, value string) *Pipeline {
	args := []string{
		"CONFIG", "SET",
		key,
		value,
	}
	p.commands = append(p.commands, pipelineCmd{args: args})
	return p
}

// Health queues a HEALTH command.
func (p *Pipeline) Health(keyspace string) *Pipeline {
	args := []string{
		"HEALTH",
	}
	if keyspace != "" {
		args = append(args, keyspace)
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseHealthResponse(m) }})
	return p
}

// Inspect queues a INSPECT command.
func (p *Pipeline) Inspect(keyspace string, credential_id string) *Pipeline {
	args := []string{
		"INSPECT",
		keyspace,
		credential_id,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseInspectResponse(m) }})
	return p
}

// Issue queues a ISSUE command.
func (p *Pipeline) Issue(keyspace string, opts *IssueOptions) *Pipeline {
	args := []string{
		"ISSUE",
		keyspace,
	}
	if opts != nil {
		if opts.Claims != nil {
			b, _ := json.Marshal(opts.Claims)
			args = append(args, "CLAIMS", string(b))
		}
		if opts.Metadata != nil {
			b, _ := json.Marshal(opts.Metadata)
			args = append(args, "META", string(b))
		}
		if opts.TtlSecs != 0 {
			args = append(args, "TTL", strconv.FormatInt(opts.TtlSecs, 10))
		}
		if opts.IdempotencyKey != "" {
			args = append(args, "IDEMPOTENCY_KEY", opts.IdempotencyKey)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseIssueResponse(m) }})
	return p
}

// Jwks queues a JWKS command.
func (p *Pipeline) Jwks(keyspace string) *Pipeline {
	args := []string{
		"JWKS",
		keyspace,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseJwksResponse(m) }})
	return p
}

// Keys queues a KEYS command.
func (p *Pipeline) Keys(keyspace string, opts *KeysOptions) *Pipeline {
	args := []string{
		"KEYS",
		keyspace,
	}
	if opts != nil {
		if opts.Cursor != "" {
			args = append(args, "CURSOR", opts.Cursor)
		}
		if opts.Pattern != "" {
			args = append(args, "MATCH", opts.Pattern)
		}
		if opts.StateFilter != "" {
			args = append(args, "STATE", opts.StateFilter)
		}
		if opts.Count != 0 {
			args = append(args, "COUNT", strconv.FormatInt(opts.Count, 10))
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseKeysResponse(m) }})
	return p
}

// Keystate queues a KEYSTATE command.
func (p *Pipeline) Keystate(keyspace string) *Pipeline {
	args := []string{
		"KEYSTATE",
		keyspace,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseKeystateResponse(m) }})
	return p
}

// PasswordChange queues a PASSWORD command.
func (p *Pipeline) PasswordChange(keyspace string, user_id string, old_password string, new_password string) *Pipeline {
	args := []string{
		"PASSWORD", "CHANGE",
		keyspace,
		user_id,
		old_password,
		new_password,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parsePasswordChangeResponse(m) }})
	return p
}

// PasswordImport queues a PASSWORD command.
func (p *Pipeline) PasswordImport(keyspace string, user_id string, hash string, opts *PasswordImportOptions) *Pipeline {
	args := []string{
		"PASSWORD", "IMPORT",
		keyspace,
		user_id,
		hash,
	}
	if opts != nil {
		if opts.Metadata != nil {
			b, _ := json.Marshal(opts.Metadata)
			args = append(args, "META", string(b))
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parsePasswordImportResponse(m) }})
	return p
}

// PasswordSet queues a PASSWORD command.
func (p *Pipeline) PasswordSet(keyspace string, user_id string, password string, opts *PasswordSetOptions) *Pipeline {
	args := []string{
		"PASSWORD", "SET",
		keyspace,
		user_id,
		password,
	}
	if opts != nil {
		if opts.Metadata != nil {
			b, _ := json.Marshal(opts.Metadata)
			args = append(args, "META", string(b))
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parsePasswordSetResponse(m) }})
	return p
}

// PasswordVerify queues a PASSWORD command.
func (p *Pipeline) PasswordVerify(keyspace string, user_id string, password string) *Pipeline {
	args := []string{
		"PASSWORD", "VERIFY",
		keyspace,
		user_id,
		password,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parsePasswordVerifyResponse(m) }})
	return p
}

// Refresh queues a REFRESH command.
func (p *Pipeline) Refresh(keyspace string, token string) *Pipeline {
	args := []string{
		"REFRESH",
		keyspace,
		token,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseRefreshResponse(m) }})
	return p
}

// Revoke queues a REVOKE command.
func (p *Pipeline) Revoke(keyspace string, credential_id string) *Pipeline {
	args := []string{
		"REVOKE",
		keyspace,
		credential_id,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseRevokeResponse(m) }})
	return p
}

// RevokeBulk queues a REVOKE command.
func (p *Pipeline) RevokeBulk(keyspace string, opts *RevokeBulkOptions) *Pipeline {
	args := []string{
		"REVOKE",
		keyspace,
	}
	if opts != nil {
		if len(opts.Ids) > 0 {
			args = append(args, "BULK")
			args = append(args, opts.Ids...)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseRevokeBulkResponse(m) }})
	return p
}

// RevokeFamily queues a REVOKE command.
func (p *Pipeline) RevokeFamily(keyspace string, opts *RevokeFamilyOptions) *Pipeline {
	args := []string{
		"REVOKE",
		keyspace,
	}
	if opts != nil {
		if opts.FamilyId != "" {
			args = append(args, "FAMILY", opts.FamilyId)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseRevokeFamilyResponse(m) }})
	return p
}

// Rotate queues a ROTATE command.
func (p *Pipeline) Rotate(keyspace string, opts *RotateOptions) *Pipeline {
	args := []string{
		"ROTATE",
		keyspace,
	}
	if opts != nil {
		if opts.Force {
			args = append(args, "FORCE")
		}
		if opts.Nowait {
			args = append(args, "NOWAIT")
		}
		if opts.Dryrun {
			args = append(args, "DRYRUN")
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseRotateResponse(m) }})
	return p
}

// Schema queues a SCHEMA command.
func (p *Pipeline) Schema(keyspace string) *Pipeline {
	args := []string{
		"SCHEMA",
		keyspace,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseSchemaResponse(m) }})
	return p
}


// Suspend queues a SUSPEND command.
func (p *Pipeline) Suspend(keyspace string, credential_id string) *Pipeline {
	args := []string{
		"SUSPEND",
		keyspace,
		credential_id,
	}
	p.commands = append(p.commands, pipelineCmd{args: args})
	return p
}

// Unsuspend queues a UNSUSPEND command.
func (p *Pipeline) Unsuspend(keyspace string, credential_id string) *Pipeline {
	args := []string{
		"UNSUSPEND",
		keyspace,
		credential_id,
	}
	p.commands = append(p.commands, pipelineCmd{args: args})
	return p
}

// Update queues a UPDATE command.
func (p *Pipeline) Update(keyspace string, credential_id string, opts *UpdateOptions) *Pipeline {
	args := []string{
		"UPDATE",
		keyspace,
		credential_id,
	}
	if opts != nil {
		if opts.Metadata != nil {
			b, _ := json.Marshal(opts.Metadata)
			args = append(args, "META", string(b))
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args})
	return p
}

// Verify queues a VERIFY command.
func (p *Pipeline) Verify(keyspace string, token string, opts *VerifyOptions) *Pipeline {
	args := []string{
		"VERIFY",
		keyspace,
		token,
	}
	if opts != nil {
		if opts.Payload != "" {
			args = append(args, "PAYLOAD", opts.Payload)
		}
		if opts.CheckRevoked {
			args = append(args, "CHECKREV")
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseVerifyResponse(m) }})
	return p
}

// Ensure imports are used.
var _ = fmt.Sprintf
var _ = strconv.FormatInt
var _ = json.Marshal
