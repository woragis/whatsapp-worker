# WhatsApp Worker v1.0.0

## Overview

The WhatsApp Worker is a production-ready, distributed WhatsApp messaging service that consumes message jobs from RabbitMQ and sends them via WhatsApp using the whatsmeow library. This initial stable release features QR code authentication, persistent session management, health checks, metrics, structured logging, and comprehensive message validation.

## Key Features

### Core Functionality
- **RabbitMQ Integration**: Consumes WhatsApp message jobs from configurable RabbitMQ queues
- **WhatsApp Messaging**: Sends text messages via WhatsApp using the whatsmeow library
- **QR Code Authentication**: Interactive QR code-based login for WhatsApp account connection
- **Session Persistence**: SQLite-based session storage for maintaining WhatsApp authentication
- **Connection Management**: Automatic connection handling with graceful reconnection
- **Message Validation**: Comprehensive validation including phone number format, content length, and security checks

### Production Features
- **Health Check Endpoint**: HTTP endpoint at `/healthz` (port 8080) for Kubernetes liveness/readiness probes
- **Prometheus Metrics**: Comprehensive metrics exposed at `/metrics` endpoint
  - Job processing counters and durations
  - Failed job tracking by error type
  - Queue depth monitoring
  - Dead letter queue size tracking
- **Structured Logging**: JSON logging in production, text logging in development with trace ID support
- **Graceful Shutdown**: Handles SIGTERM/SIGINT signals for clean shutdown
- **Connection State Tracking**: Real-time connection status monitoring

### WhatsApp Integration
- **whatsmeow Library**: Uses go.mau.fi/whatsmeow for WhatsApp Web API integration
- **End-to-End Encryption**: Native support for WhatsApp's end-to-end encryption
- **Event Handling**: Comprehensive event handling for connection state changes
- **QR Code Management**: QR code generation, storage, and automatic cleanup

### Configuration
- **Environment-based Configuration**: All settings via environment variables
- **RabbitMQ Configuration**: Flexible connection string or individual parameters
- **Queue Configuration**: Configurable queue names, exchanges, and routing keys
- **Session Configuration**: Configurable session storage path
- **Prefetch Control**: Configurable message prefetch count for optimal throughput

### Code Quality
- **Comprehensive Testing**: Unit tests and integration tests included
- **Test Coverage**: Coverage reporting with threshold checks (70% minimum)
- **Clean Architecture**: Separation of concerns with internal and pkg packages
- **Go 1.24**: Built with latest Go version

## Architecture

```
┌─────────────┐
│   Server    │
└──────┬──────┘
       │
       │ Publish WhatsAppEnvelope
       │
┌──────▼──────────┐
│   RabbitMQ      │
│   Queue         │
└──────┬──────────┘
       │
       │ WhatsAppEnvelope Message
       │
┌──────▼─────────────────────────┐
│   WhatsApp Worker              │
│                                │
│  ┌──────────────────┐          │
│  │  Queue Consumer  │          │
│  └────────┬─────────┘          │
│           │                    │
│  ┌────────▼─────────┐          │
│  │ Message          │          │
│  │ Validator        │          │
│  └────────┬─────────┘          │
│           │                    │
│  ┌────────▼─────────┐          │
│  │ WhatsApp         │          │
│  │ Notifier         │          │
│  │ (whatsmeow)      │          │
│  └────────┬─────────┘          │
└───────────┼────────────────────┘
            │
    ┌───────┴────────┐
    │                │
┌───▼────┐      ┌───▼─────┐
│WhatsApp│      │SQLite   │
│API     │      │Session  │
└────────┘      └─────────┘
```

## Dependencies

- **go.mau.fi/whatsmeow**: WhatsApp Web API client library
- **github.com/rabbitmq/amqp091-go**: RabbitMQ client library
- **github.com/prometheus/client_golang**: Prometheus metrics
- **github.com/mattn/go-sqlite3**: SQLite driver for session storage
- **github.com/google/uuid**: UUID generation
- **github.com/stretchr/testify**: Testing utilities

## Configuration Variables

### RabbitMQ Configuration
- `RABBITMQ_URL`: Full AMQP connection URL (overrides individual settings)
- `RABBITMQ_USER`: RabbitMQ username (default: "woragis")
- `RABBITMQ_PASSWORD`: RabbitMQ password (default: "woragis")
- `RABBITMQ_HOST`: RabbitMQ hostname (default: "rabbitmq")
- `RABBITMQ_PORT`: RabbitMQ port (default: "5672")
- `RABBITMQ_VHOST`: RabbitMQ virtual host (default: "woragis")

### Worker Configuration
- `WHATSAPP_QUEUE_NAME`: Queue name (default: "whatsapp.queue")
- `WHATSAPP_EXCHANGE`: Exchange name (default: "woragis.notifications")
- `WHATSAPP_ROUTING_KEY`: Routing key (default: "whatsapp.send")
- `WHATSAPP_PREFETCH_COUNT`: Prefetch count (default: 1)

### WhatsApp Configuration
- `WHATSAPP_SESSION_PATH`: Path to store WhatsApp session data (default: "/app/session")
  - Should be a persistent volume in production
  - Contains SQLite database with authentication tokens

### Environment
- `ENV`: Environment mode - "development" or "production" (affects logging)
- `LOG_TO_FILE`: Enable file logging in development (default: false)
- `LOG_DIR`: Directory for log files (default: "logs")

## Message Format

The worker expects `WhatsAppEnvelope` messages in the following JSON format:

```json
{
  "user_id": "uuid",
  "text_message": "Hello, this is a WhatsApp message",
  "destination": "+5511999999999"
}
```

### Message Validation

Messages are validated before processing:
- **user_id**: Must be a valid UUID
- **destination**: Must be a valid international phone number (format: +[country code][number], 8-15 digits)
- **text_message**: 
  - Required, 1-4096 characters
  - Validated against SQL injection patterns
  - Validated against XSS patterns

## Phone Number Format

Phone numbers must be in international format:
- **Format**: `+[country code][number]`
- **Example**: `+5511999999999` (Brazil), `+1234567890` (USA)
- **Rules**:
  - Must start with `+`
  - 8-15 digits after the country code
  - Only digits allowed after `+`
  - Common separators (spaces, dashes, parentheses) are automatically cleaned

## WhatsApp Authentication

The worker uses QR code-based authentication:

1. **First Run**: 
   - No session found → QR code is generated
   - QR code is logged to stdout and logger
   - Scan QR code with WhatsApp mobile app to authenticate

2. **Subsequent Runs**:
   - Session stored in SQLite database
   - Automatically reconnects using stored session
   - No QR code needed unless session expires or is invalidated

3. **Session Storage**:
   - Stored in `WHATSAPP_SESSION_PATH/whatsapp.db`
   - Must be a persistent volume in production
   - Contains authentication tokens and device information

## Health Check

The service exposes a health check endpoint at `GET /healthz`:

- **Healthy Response (200 OK)**: All checks pass
- **Unhealthy Response (503)**: RabbitMQ connection is down

The health check includes:
- RabbitMQ connection status check (2s timeout)
- Result caching (5s TTL) for performance

**Note**: WhatsApp connection status is not checked in health endpoint (as it may require QR code scan initially).

## Metrics

Prometheus metrics are exposed at `GET /metrics`:

- `worker_jobs_processed_total{worker, status}`: Total jobs processed
- `worker_job_duration_seconds{worker}`: Job processing duration histogram
- `worker_jobs_failed_total{worker, error_type}`: Failed job counter (by error type)
- `worker_jobs_retried_total{worker}`: Retried job counter
- `queue_depth{queue_name}`: Current queue depth
- `queue_dlq_size{queue_name}`: Dead letter queue size

## Deployment

### Docker
The service can be containerized and deployed as a standalone worker or in Kubernetes.

**Important**: The session directory (`WHATSAPP_SESSION_PATH`) must be a persistent volume to maintain WhatsApp authentication across restarts.

### Kubernetes
- **Liveness Probe**: `GET /healthz` every 10s, timeout 5s
- **Readiness Probe**: `GET /healthz` every 10s, timeout 5s
- **Port**: 8080 (health and metrics)
- **Persistent Volume**: Required for `WHATSAPP_SESSION_PATH`

### Scaling
- **Single Instance**: Recommended for WhatsApp workers (WhatsApp Web allows one active session per account)
- **Multiple Accounts**: Run separate instances with different session paths for multiple WhatsApp accounts

## Development

### Testing
```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run with coverage
make test-cov

# Check coverage threshold (70%)
make test-cov-check

# Clean test artifacts
make clean
```

### Building
```bash
# Install dependencies
make install

# Build binary (manual - no make target, use go build)
go build -o bin/whatsapp-worker ./cmd/whatsapp-worker
```

### Running Locally
```bash
# Set environment variables
export RABBITMQ_URL="amqp://woragis:woragis@localhost:5672/woragis"
export WHATSAPP_SESSION_PATH="./session"
export ENV=development

# Run worker
./bin/whatsapp-worker

# On first run, scan QR code from logs with WhatsApp mobile app
```

## Error Handling

- **Invalid Messages**: Rejected without requeue (sent to DLQ)
  - Invalid phone number format
  - Missing required fields
  - Content validation failures
- **Send Failures**: Retried (requeued)
  - WhatsApp not connected (requires QR code scan)
  - Network errors
  - WhatsApp API errors
- **Connection Errors**: 
  - Logged with warnings
  - Messages are requeued for retry when connection is restored

## Important Notes

### WhatsApp Web Limitations
- **Single Session**: WhatsApp Web allows only one active session per account
- **Session Persistence**: Session must be stored in a persistent volume
- **QR Code Expiry**: QR codes expire after a timeout period; restart worker to generate new one
- **Account Restrictions**: Follow WhatsApp's Terms of Service and API usage policies

### Security Considerations
- **Session Security**: Session files contain authentication tokens - secure appropriately
- **Phone Number Privacy**: Never log full phone numbers in production
- **Content Validation**: Messages are validated for SQL injection and XSS patterns
- **Rate Limiting**: Be aware of WhatsApp's rate limits to avoid account restrictions

## Documentation

- **HEALTH_CHECK.md**: Detailed health check documentation
- **LOGGING.md**: Structured logging guidelines and usage
- **tests/README.md**: Testing documentation
- **internal/integration/README.md**: Integration test documentation

## Breaking Changes

None - This is the initial release.

## Future Enhancements

Potential improvements for future versions:
- Media message support (images, documents, videos)
- Message delivery status tracking
- Read receipts handling
- Group message support
- Message templates
- WhatsApp Business API integration option
- Connection retry logic with exponential backoff
- Session backup and restore
- Multi-account support with session routing
- Webhook notifications for message status
- Message scheduling
- Rate limiting and throttling

## Contributors

Initial release by the Woragis team.

## License

Part of the Woragis backend infrastructure.

