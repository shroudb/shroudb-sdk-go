// ShroudbSentry pipeline for batching commands.
//
// Auto-generated from shroudb-sentry protocol spec. Do not edit.

package shroudb_sentry

import (
	"fmt"
	"strconv"
)

// Pipeline batches multiple ShroudbSentry commands and executes them in a single round-trip.
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
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseAuthResponse(m) }})
	return p
}

// Evaluate queues a EVALUATE command.
func (p *Pipeline) Evaluate(json string) *Pipeline {
	args := []string{
		"EVALUATE",
		json,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseEvaluateResponse(m) }})
	return p
}

// Health queues a HEALTH command.
func (p *Pipeline) Health() *Pipeline {
	args := []string{
		"HEALTH",
	}
	p.commands = append(p.commands, pipelineCmd{args: args})
	return p
}

// KeyInfo queues a KEY_INFO command.
func (p *Pipeline) KeyInfo() *Pipeline {
	args := []string{
		"KEY_INFO",
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseKeyInfoResponse(m) }})
	return p
}

// KeyRotate queues a KEY_ROTATE command.
func (p *Pipeline) KeyRotate(opts *KeyRotateOptions) *Pipeline {
	args := []string{
		"KEY_ROTATE",
	}
	if opts != nil {
		if opts.Force != "" {
			args = append(args, "FORCE", opts.Force)
		}
		if opts.Dryrun != "" {
			args = append(args, "DRYRUN", opts.Dryrun)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseKeyRotateResponse(m) }})
	return p
}

// PolicyInfo queues a POLICY_INFO command.
func (p *Pipeline) PolicyInfo(name string) *Pipeline {
	args := []string{
		"POLICY_INFO",
		name,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parsePolicyInfoResponse(m) }})
	return p
}

// PolicyList queues a POLICY_LIST command.
func (p *Pipeline) PolicyList() *Pipeline {
	args := []string{
		"POLICY_LIST",
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parsePolicyListResponse(m) }})
	return p
}

// PolicyReload queues a POLICY_RELOAD command.
func (p *Pipeline) PolicyReload() *Pipeline {
	args := []string{
		"POLICY_RELOAD",
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parsePolicyReloadResponse(m) }})
	return p
}

// Ensure imports are used.
var _ = fmt.Sprintf
var _ = strconv.FormatInt
