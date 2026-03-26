// ShroudbTransit pipeline for batching commands.
//
// Auto-generated from shroudb-transit protocol spec. Do not edit.

package shroudb_transit

import (
	"fmt"
	"strconv"
)

// Pipeline batches multiple ShroudbTransit commands and executes them in a single round-trip.
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

// Decrypt queues a DECRYPT command.
func (p *Pipeline) Decrypt(keyring string, ciphertext string, opts *DecryptOptions) *Pipeline {
	args := []string{
		"DECRYPT",
		keyring,
		ciphertext,
	}
	if opts != nil {
		if opts.Context != "" {
			args = append(args, "CONTEXT", opts.Context)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseDecryptResponse(m) }})
	return p
}

// Encrypt queues a ENCRYPT command.
func (p *Pipeline) Encrypt(keyring string, plaintext string, opts *EncryptOptions) *Pipeline {
	args := []string{
		"ENCRYPT",
		keyring,
		plaintext,
	}
	if opts != nil {
		if opts.Context != "" {
			args = append(args, "CONTEXT", opts.Context)
		}
		if opts.KeyVersion != "" {
			args = append(args, "KEY_VERSION", opts.KeyVersion)
		}
		if opts.Convergent != "" {
			args = append(args, "CONVERGENT", opts.Convergent)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseEncryptResponse(m) }})
	return p
}

// GenerateDataKey queues a GENERATE_DATA_KEY command.
func (p *Pipeline) GenerateDataKey(keyring string, opts *GenerateDataKeyOptions) *Pipeline {
	args := []string{
		"GENERATE_DATA_KEY",
		keyring,
	}
	if opts != nil {
		if opts.Bits != "" {
			args = append(args, "BITS", opts.Bits)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseGenerateDataKeyResponse(m) }})
	return p
}

// Health queues a HEALTH command.
func (p *Pipeline) Health(keyring string) *Pipeline {
	args := []string{
		"HEALTH",
	}
	if keyring != "" {
		args = append(args, keyring)
	}
	p.commands = append(p.commands, pipelineCmd{args: args})
	return p
}

// KeyInfo queues a KEY_INFO command.
func (p *Pipeline) KeyInfo(keyring string) *Pipeline {
	args := []string{
		"KEY_INFO",
		keyring,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseKeyInfoResponse(m) }})
	return p
}

// Rewrap queues a REWRAP command.
func (p *Pipeline) Rewrap(keyring string, ciphertext string, opts *RewrapOptions) *Pipeline {
	args := []string{
		"REWRAP",
		keyring,
		ciphertext,
	}
	if opts != nil {
		if opts.Context != "" {
			args = append(args, "CONTEXT", opts.Context)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseRewrapResponse(m) }})
	return p
}

// Rotate queues a ROTATE command.
func (p *Pipeline) Rotate(keyring string, opts *RotateOptions) *Pipeline {
	args := []string{
		"ROTATE",
		keyring,
	}
	if opts != nil {
		if opts.Force != "" {
			args = append(args, "FORCE", opts.Force)
		}
		if opts.Dryrun != "" {
			args = append(args, "DRYRUN", opts.Dryrun)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseRotateResponse(m) }})
	return p
}

// Sign queues a SIGN command.
func (p *Pipeline) Sign(keyring string, data string, opts *SignOptions) *Pipeline {
	args := []string{
		"SIGN",
		keyring,
		data,
	}
	if opts != nil {
		if opts.Algorithm != "" {
			args = append(args, "ALGORITHM", opts.Algorithm)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseSignResponse(m) }})
	return p
}

// VerifySignature queues a VERIFY_SIGNATURE command.
func (p *Pipeline) VerifySignature(keyring string, data string, signature string) *Pipeline {
	args := []string{
		"VERIFY_SIGNATURE",
		keyring,
		data,
		signature,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseVerifySignatureResponse(m) }})
	return p
}

// Ensure imports are used.
var _ = fmt.Sprintf
var _ = strconv.FormatInt
