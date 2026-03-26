// ShroudbPulse pipeline for batching commands.
//
// Auto-generated from shroudb-pulse protocol spec. Do not edit.

package shroudb_pulse

import (
	"fmt"
	"strconv"
)

// Pipeline batches multiple ShroudbPulse commands and executes them in a single round-trip.
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

// Actors queues a ACTORS command.
func (p *Pipeline) Actors(opts *ActorsOptions) *Pipeline {
	args := []string{
		"ACTORS",
	}
	if opts != nil {
		if opts.Window != "" {
			args = append(args, "WINDOW", opts.Window)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseActorsResponse(m) }})
	return p
}

// Auth queues a AUTH command.
func (p *Pipeline) Auth(token string) *Pipeline {
	args := []string{
		"AUTH",
		token,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseAuthResponse(m) }})
	return p
}

// Count queues a COUNT command.
func (p *Pipeline) Count() *Pipeline {
	args := []string{
		"COUNT",
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseCountResponse(m) }})
	return p
}

// Errors queues a ERRORS command.
func (p *Pipeline) Errors(opts *ErrorsOptions) *Pipeline {
	args := []string{
		"ERRORS",
	}
	if opts != nil {
		if opts.Engine != "" {
			args = append(args, "ENGINE", opts.Engine)
		}
		if opts.Window != "" {
			args = append(args, "WINDOW", opts.Window)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseErrorsResponse(m) }})
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

// Hotspots queues a HOTSPOTS command.
func (p *Pipeline) Hotspots(opts *HotspotsOptions) *Pipeline {
	args := []string{
		"HOTSPOTS",
	}
	if opts != nil {
		if opts.Engine != "" {
			args = append(args, "ENGINE", opts.Engine)
		}
		if opts.Window != "" {
			args = append(args, "WINDOW", opts.Window)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseHotspotsResponse(m) }})
	return p
}

// Ingest queues a INGEST command.
func (p *Pipeline) Ingest(json string) *Pipeline {
	args := []string{
		"INGEST",
		json,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseIngestResponse(m) }})
	return p
}

// IngestBatch queues a INGEST_BATCH command.
func (p *Pipeline) IngestBatch(json string) *Pipeline {
	args := []string{
		"INGEST_BATCH",
		json,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseIngestBatchResponse(m) }})
	return p
}

// Query queues a QUERY command.
func (p *Pipeline) Query() *Pipeline {
	args := []string{
		"QUERY",
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseQueryResponse(m) }})
	return p
}

// SourceList queues a SOURCE_LIST command.
func (p *Pipeline) SourceList() *Pipeline {
	args := []string{
		"SOURCE_LIST",
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseSourceListResponse(m) }})
	return p
}

// SourceStatus queues a SOURCE_STATUS command.
func (p *Pipeline) SourceStatus() *Pipeline {
	args := []string{
		"SOURCE_STATUS",
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseSourceStatusResponse(m) }})
	return p
}

// Ensure imports are used.
var _ = fmt.Sprintf
var _ = strconv.FormatInt
