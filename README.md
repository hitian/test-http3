# HTTP/3 Checker

A command-line tool to test HTTP/3 support for domains. This tool first checks if a server advertises HTTP/3 support via the Alt-Svc header using a standard HTTP connection, then attempts to establish an actual HTTP/3 connection for verification.

## Features

- Two-step verification process:
  1. Check Alt-Svc header via standard HTTP/1.1 or HTTP/2 connection
  2. Attempt actual HTTP/3/QUIC connection
- Detailed error diagnostics with specific troubleshooting suggestions
- Configurable timeout settings
- Cross-platform support (Linux and macOS, both amd64 and arm64)

## Installation

### From Source

```bash
go install github.com/hitian/test-http3@latest
```

### Build from Source

```bash
git clone https://github.com/hitian/test-http3
cd test-http3
go build -o http3check
```

### Using Makefile

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

## Usage

```bash
./http3check [options] <domain>
```

### Options

- `-timeout duration`: Set timeout for connection attempts (default: 10s)

### Examples

```bash
# Basic usage
./http3check google.com

# With custom timeout
./http3check -timeout 15s example.com

# Test with HTTPS URL
./http3check https://cloudflare.com
```

## Sample Output

### Successful HTTP/3 Connection
```
Testing https://google.com for HTTP/3 support (timeout: 10s)...

Step 1: Checking https://google.com with standard HTTP client...
Alt-Svc header: h3=":443"; ma=2592000
Standard connection protocol: HTTP/2.0

✅ https://google.com advertises HTTP/3 support via Alt-Svc header

Step 2: Attempting HTTP/3 connection to https://google.com...
✅ HTTP/3 connection successful!
Protocol: HTTP/3.0
Status: 200 OK
Server: gws
```

### Failed HTTP/3 Connection with Diagnostics
```
Testing https://example.com for HTTP/3 support (timeout: 10s)...

Step 1: Checking https://example.com with standard HTTP client...
Alt-Svc header: h3=":443"; ma=86400
Standard connection protocol: HTTP/1.1

✅ https://example.com advertises HTTP/3 support via Alt-Svc header

Step 2: Attempting HTTP/3 connection to https://example.com...

❌ HTTP/3 connection failed: context deadline exceeded

Possible issues:
• Connection timeout - server may not support HTTP/3 or QUIC is blocked
• Try increasing timeout with -timeout flag
• Check if UDP port 443 is accessible
• Server advertises HTTP/3 but implementation may be incomplete
• Try testing with other HTTP/3 tools like curl --http3
```

## How It Works

1. **Standard HTTP Check**: First connects using HTTP/1.1 or HTTP/2 to check the `Alt-Svc` header
2. **HTTP/3 Verification**: If HTTP/3 is advertised, attempts an actual HTTP/3 connection using QUIC
3. **Error Diagnosis**: Provides specific troubleshooting suggestions based on the type of failure

## Common Issues

- **UDP Blocking**: Many networks and firewalls block UDP traffic, which HTTP/3 requires
- **Incomplete Implementations**: Some servers advertise HTTP/3 but have buggy implementations
- **Certificate Issues**: HTTP/3 requires valid TLS certificates that support QUIC
- **Timeout Issues**: HTTP/3 handshakes can be slower than HTTP/2, try increasing timeout

## Requirements

- Go 1.19 or later
- Internet connection
- UDP port 443 access for HTTP/3 connections

## Dependencies

- [quic-go](https://github.com/quic-go/quic-go): Go implementation of QUIC and HTTP/3

## License

MIT License - see LICENSE file for details

## Contributing

Contributions welcome! Please feel free to submit a Pull Request.