// ShroudbMint pipeline for batching commands.
//
// Auto-generated from shroudb-mint protocol spec. Do not edit.

package shroudb_mint

import (
	"fmt"
	"strconv"
)

// Pipeline batches multiple ShroudbMint commands and executes them in a single round-trip.
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

// CaCreate queues a CA_CREATE command.
func (p *Pipeline) CaCreate(ca string, algorithm string, subject string, ttl_days string, opts *CaCreateOptions) *Pipeline {
	args := []string{
		"CA_CREATE",
		ca,
		algorithm,
		subject,
		ttl_days,
	}
	if opts != nil {
		if opts.Parent != "" {
			args = append(args, "PARENT", opts.Parent)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseCaCreateResponse(m) }})
	return p
}

// CaExport queues a CA_EXPORT command.
func (p *Pipeline) CaExport(ca string, opts *CaExportOptions) *Pipeline {
	args := []string{
		"CA_EXPORT",
		ca,
	}
	if opts != nil {
		if opts.Format != "" {
			args = append(args, "FORMAT", opts.Format)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseCaExportResponse(m) }})
	return p
}

// CaInfo queues a CA_INFO command.
func (p *Pipeline) CaInfo(ca string) *Pipeline {
	args := []string{
		"CA_INFO",
		ca,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseCaInfoResponse(m) }})
	return p
}

// CaList queues a CA_LIST command.
func (p *Pipeline) CaList() *Pipeline {
	args := []string{
		"CA_LIST",
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseCaListResponse(m) }})
	return p
}

// CaRotate queues a CA_ROTATE command.
func (p *Pipeline) CaRotate(ca string, opts *CaRotateOptions) *Pipeline {
	args := []string{
		"CA_ROTATE",
		ca,
	}
	if opts != nil {
		if opts.Force != "" {
			args = append(args, "FORCE", opts.Force)
		}
		if opts.Dryrun != "" {
			args = append(args, "DRYRUN", opts.Dryrun)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseCaRotateResponse(m) }})
	return p
}

// CrlInfo queues a CRL_INFO command.
func (p *Pipeline) CrlInfo(ca string) *Pipeline {
	args := []string{
		"CRL_INFO",
		ca,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseCrlInfoResponse(m) }})
	return p
}

// Health queues a HEALTH command.
func (p *Pipeline) Health(ca string) *Pipeline {
	args := []string{
		"HEALTH",
	}
	if ca != "" {
		args = append(args, ca)
	}
	p.commands = append(p.commands, pipelineCmd{args: args})
	return p
}

// Inspect queues a INSPECT command.
func (p *Pipeline) Inspect(ca string, serial string) *Pipeline {
	args := []string{
		"INSPECT",
		ca,
		serial,
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseInspectResponse(m) }})
	return p
}

// Issue queues a ISSUE command.
func (p *Pipeline) Issue(ca string, subject string, profile string, opts *IssueOptions) *Pipeline {
	args := []string{
		"ISSUE",
		ca,
		subject,
		profile,
	}
	if opts != nil {
		if opts.Ttl != "" {
			args = append(args, "TTL", opts.Ttl)
		}
		if opts.SanDns != "" {
			args = append(args, "SAN_DNS", opts.SanDns)
		}
		if opts.SanIp != "" {
			args = append(args, "SAN_IP", opts.SanIp)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseIssueResponse(m) }})
	return p
}

// IssueFromCsr queues a ISSUE_FROM_CSR command.
func (p *Pipeline) IssueFromCsr(ca string, csr_pem string, profile string, opts *IssueFromCsrOptions) *Pipeline {
	args := []string{
		"ISSUE_FROM_CSR",
		ca,
		csr_pem,
		profile,
	}
	if opts != nil {
		if opts.Ttl != "" {
			args = append(args, "TTL", opts.Ttl)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseIssueFromCsrResponse(m) }})
	return p
}

// ListCerts queues a LIST_CERTS command.
func (p *Pipeline) ListCerts(ca string, opts *ListCertsOptions) *Pipeline {
	args := []string{
		"LIST_CERTS",
		ca,
	}
	if opts != nil {
		if opts.State != "" {
			args = append(args, "STATE", opts.State)
		}
		if opts.Limit != 0 {
			args = append(args, "LIMIT", strconv.FormatInt(opts.Limit, 10))
		}
		if opts.Offset != 0 {
			args = append(args, "OFFSET", strconv.FormatInt(opts.Offset, 10))
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseListCertsResponse(m) }})
	return p
}

// Renew queues a RENEW command.
func (p *Pipeline) Renew(ca string, serial string, opts *RenewOptions) *Pipeline {
	args := []string{
		"RENEW",
		ca,
		serial,
	}
	if opts != nil {
		if opts.Ttl != "" {
			args = append(args, "TTL", opts.Ttl)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseRenewResponse(m) }})
	return p
}

// Revoke queues a REVOKE command.
func (p *Pipeline) Revoke(ca string, serial string, opts *RevokeOptions) *Pipeline {
	args := []string{
		"REVOKE",
		ca,
		serial,
	}
	if opts != nil {
		if opts.Reason != "" {
			args = append(args, "REASON", opts.Reason)
		}
	}
	p.commands = append(p.commands, pipelineCmd{args: args, parser: func(m map[string]any) any { return parseRevokeResponse(m) }})
	return p
}

// Ensure imports are used.
var _ = fmt.Sprintf
var _ = strconv.FormatInt
