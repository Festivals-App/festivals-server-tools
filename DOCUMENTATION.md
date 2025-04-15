# 🧰 FestivalsApp Server Tools

The `servertools` package provides a modular set of tools designed to support and streamline server-side functionality
in distributed service environments. It focuses on common operational needs such as service monitoring, HTTP handling,
structured logging, update automation, and utility functions.

<hr/>
<p align="center">
  <a href="#Servertools">Servertools</a> •
  <a href="#Loggertools">Loggertools</a> •
  <a href="#Heartbeattools">Heartbeattools</a> •
  <a href="#Responsetools">Responsetools</a> •
  <a href="#Updatetools">Updatetools</a>
</p>
<hr/>

This package is built to promote consistency, observability, and maintainability across microservices or backend systems.
It's especially valuable in production-grade deployments where reliability and traceability are key.

## Servertools

The servertools component of the `servertools` package provides basic file-related helper functions.

- `FileExists`: Checks whether a given file path exists and is not a directory. Useful for validating paths
  before attempting file operations.
- `ExpandTilde`: Expands a leading `~` in a file path to the current user's home directory. This allows support
  for user-friendly paths like `~/config.yaml`.

These utilities are simple but important for improving portability and robustness when working with file paths and user environments.

## Loggertools

The loggertools component of the `servertools` package provides structured logging and request tracing middleware
for HTTP servers, enabling robust observability and error handling in production environments.

At its core, the `Middleware` function wraps incoming HTTP requests to capture detailed telemetry,
including request duration, HTTP status, and response size. Each request is automatically enriched
with a unique request ID (if available), and is logged at different verbosity levels based on the outcome:

- **Successful requests** (2xx) are logged at the `trace` level to a dedicated trace logger.
- **Error responses** (3xx and above) are logged at the `debug` level using the global logger.

The middleware also includes built-in **panic recovery**, logging stack traces and sending a 500 error response
in case of unexpected runtime failures—ensuring system stability and post-mortem visibility.

The package supports structured, timestamped logging using [Zerolog](https://github.com/rs/zerolog),
and enables both **global** and **trace-specific** loggers:

- `InitializeGlobalLogger`: Configures the system-wide logger, supporting optional console output and rolling file logs.
- `TraceLogger`: Creates a logger for access tracing used in the middleware.
- `NewRollingFile`: Log files are rotated using [Lumberjack](https://github.com/natefinch/lumberjack) based on size (50MB),
  backups (10 files), and age (31 days).

All logs are timestamped in RFC3339 format.

### Common Format

Each log entry is a single-line JSON object, with the following common structure:

```json
{
  "level": "trace|debug|error|fatal",
  "time": "2025-04-15T12:34:56Z",
  "message": "Log message here",
  "request_id": "uuid-optional"
}
```

### Trace Logger Format

Used for successful HTTP requests (2xx responses), logged at trace level. Example structure:

```json
{
  "level": "trace",
  "time": "2025-04-15T12:34:56Z",
  "type": "access",
  "request_id": "123e4567-e89b-12d3-a456-426614174000",
  "url": "api.example.com/v1/resource",
  "method": "GET",
  "status": 200,
  "latency_ms": 42.37,
  "bytes_out": 512,
  "message": "Incoming request"
}
```

- type: Always "access" to distinguish trace logs.
- latency_ms: Request duration in milliseconds.
- bytes_out: Number of bytes written in the response.

### Global Logger Format

Used for:

- All other manual log events
- Failed HTTP requests (non-2xx)
- Errors, panics, or fatal issues

Example: debug log from a failed request:

```json
{
  "level": "debug",
  "time": "2025-04-15T12:35:01Z",
  "request_id": "123e4567-e89b-12d3-a456-426614174001",
  "url": "api.example.com/v1/resource",
  "method": "POST",
  "status": 500,
  "latency_ms": 128.09,
  "bytes_out": 0,
  "message": "Incoming request"
}
```

## Heartbeattools

The heartbeattools component of the `servertools` package provides essential functionality for service availability
monitoring in a distributed system. It enables individual services to securely report their operational status (heartbeats)
to the festivals-gateway using mutual TLS (mTLS) for authentication and encryption.

At the core of the heartbeattools is the Heartbeat struct, which contains key metadata about each service—such as
service name, host, port, and availability status. These heartbeat signals are sent via HTTP POST requests using
the SendHeartbeat function, which automatically attaches important headers including a unique request ID and service identifier.

To ensure secure communication, the HeartbeatClient function constructs a custom http.Client with full TLS configuration.
It loads client certificates and a trusted certificate authority (CA), allowing both the client and server
to mutually authenticate each other—preventing unauthorized access and ensuring data integrity.

The heartbeattools supports a predefined list of service identifiers, making it easy to integrate across a microservice architecture.

## Responsetools

The responsetools component of the `servertools` package provides utility functions to simplify and
standardize HTTP responses across the application.

It includes several response types:

- `RespondCode`: Sends a plain response with a status code only.
- `RespondString`: Sends a plain text response with a custom message.
- `RespondJSON`: Wraps a payload in a JSON structure under a `"data"` key and writes it to the response.
  Handles empty slices and marshaling errors gracefully.
- `RespondError`: Similar to `RespondJSON`, but wraps the message under an `"error"` key and
  is used for sending error responses.
- `UnauthorizedResponse`: Sends a `401 Unauthorized` response with the appropriate `WWW-Authenticate` header
  for HTTP Basic Auth.

All functions include error logging using Zerolog, ensuring failures in response writing are captured.
This package promotes consistency in API responses and improves maintainability in request handling logic.

## Updatetools

The updatetools component of the `servertools` package provides a simple mechanism for automatically updating a server
based on the latest GitHub release.

The RunUpdate function checks the current server version against the latest available tag from a specified GitHub repository.
If a newer version exists, it executes a provided update script. Development servers are excluded from this process
and must be updated manually. The version check is handled by the LatestVersion function, which uses
the GitHub API to retrieve the most recent tag from the given repository.

This utility helps streamline version management and ensures production environments stay up-to-date
with minimal manual intervention.
