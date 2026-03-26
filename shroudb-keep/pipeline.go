// ShroudbKeep pipeline for batching commands.
//
// Auto-generated from shroudb-keep protocol spec. Do not edit.

package shroudb_keep

import (
	"fmt"
	"strconv"
)

// Pipeline batches multiple ShroudbKeep commands and executes them in a single round-trip.
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

// Delete queues a DELETE command.
func (p *Pipeline) Delete(path string) *Pipeline {
	args := []string{
		"DELETE",
		path,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseDeleteResponse(m) }})
	return p
}

// Get queues a GET command.
func (p *Pipeline) Get(path string, opts *GetOptions) *Pipeline {
	args := []string{
		"GET",
		path,
	}
	if opts != nil {
		if opts.Version != "" {
			args = append(args, "VERSION", opts.Version)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseGetResponse(m) }})
	return p
}

// Health queues a HEALTH command.
func (p *Pipeline) Health(path string) *Pipeline {
	args := []string{
		"HEALTH",
	}
	if path != "" {
		args = append(args, path)
	}
	p.commands = append(p.commands, pipelineCmd{args: args})
	return p
}

// List queues a LIST command.
func (p *Pipeline) List(prefix string) *Pipeline {
	args := []string{
		"LIST",
	}
	if prefix != "" {
		args = append(args, prefix)
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseListResponse(m) }})
	return p
}

// Put queues a PUT command.
func (p *Pipeline) Put(path string, value string, opts *PutOptions) *Pipeline {
	args := []string{
		"PUT",
		path,
		value,
	}
	if opts != nil {
		if opts.Meta != "" {
			args = append(args, "META", opts.Meta)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parsePutResponse(m) }})
	return p
}

// Rotate queues a ROTATE command.
func (p *Pipeline) Rotate(path string) *Pipeline {
	args := []string{
		"ROTATE",
		path,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseRotateResponse(m) }})
	return p
}

// Versions queues a VERSIONS command.
func (p *Pipeline) Versions(path string) *Pipeline {
	args := []string{
		"VERSIONS",
		path,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseVersionsResponse(m) }})
	return p
}

// Ensure imports are used.
var _ = fmt.Sprintf
var _ = strconv.FormatInt
