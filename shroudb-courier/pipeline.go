// ShroudbCourier pipeline for batching commands.
//
// Auto-generated from shroudb-courier protocol spec. Do not edit.

package shroudb_courier

import (
	"fmt"
	"strconv"
)

// Pipeline batches multiple ShroudbCourier commands and executes them in a single round-trip.
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

// ChannelInfo queues a CHANNEL_INFO command.
func (p *Pipeline) ChannelInfo(channel string) *Pipeline {
	args := []string{
		"CHANNEL_INFO",
		channel,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseChannelInfoResponse(m) }})
	return p
}

// ChannelList queues a CHANNEL_LIST command.
func (p *Pipeline) ChannelList() *Pipeline {
	args := []string{
		"CHANNEL_LIST",
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseChannelListResponse(m) }})
	return p
}

// Connections queues a CONNECTIONS command.
func (p *Pipeline) Connections() *Pipeline {
	args := []string{
		"CONNECTIONS",
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseConnectionsResponse(m) }})
	return p
}

// Deliver queues a DELIVER command.
func (p *Pipeline) Deliver(json string) *Pipeline {
	args := []string{
		"DELIVER",
		json,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseDeliverResponse(m) }})
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

// TemplateInfo queues a TEMPLATE_INFO command.
func (p *Pipeline) TemplateInfo(name string) *Pipeline {
	args := []string{
		"TEMPLATE_INFO",
		name,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseTemplateInfoResponse(m) }})
	return p
}

// TemplateList queues a TEMPLATE_LIST command.
func (p *Pipeline) TemplateList() *Pipeline {
	args := []string{
		"TEMPLATE_LIST",
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseTemplateListResponse(m) }})
	return p
}

// TemplateReload queues a TEMPLATE_RELOAD command.
func (p *Pipeline) TemplateReload() *Pipeline {
	args := []string{
		"TEMPLATE_RELOAD",
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseTemplateReloadResponse(m) }})
	return p
}

// Ensure imports are used.
var _ = fmt.Sprintf
var _ = strconv.FormatInt
