
## Overview

This service fetches posts from the JSONPlaceholder API, transforms each record by adding metadata (`ingested_at`, `source`), and stores them in Amazon DynamoDB. It is containerized with Docker, tested end-to-end (unit and integration), and can be extended with CI/CD pipelines.

---

## Architecture

1. **Fetcher**: HTTP client with retry and fetches data from `https://jsonplaceholder.typicode.com/posts`.
2. **Transformer**: Adds UTC timestamp `ingested_at` and static `source` field.
3. **Storage**: Writes each transformed item into DynamoDB table `cyderes_api_logs`.
4. **Containerization**: Dockerfile builds the Go binary; `docker-compose.yml` spins up the service.
5. **Testing**:

   * Unit tests for transformation logic.
   * Unit tests for fetch_data using httptest
   * End to end Integration tests against DynamoDB Local via Docker.


## Storage Justification

* **Amazon DynamoDB**: a fully managed NoSQL database optimized for high-throughput writes and key-value/JSON data. It auto-scales, has a free tier, and the AWS SDK for Go v2 provides seamless integration. It fits moderately structured logs with simple primary-key lookups. Trade-offs: cost can rise under heavy read/write loads, but our moderate volume keeps it economical.

---

## Setup & Running Locally

### Prerequisites

* Docker & Docker Compose
* Go 1.20+
* AWS CLI (for `aws configure`, optional)
* IAM role configured with requir

### 1. Clone & Build

```bash
git clone <repo-url>
cd cyderes
```

### 2. Build docker image

```bash
docker build -f build/docker/Dockerfile -t cyderes-app:1.1 .
```

### 3. Start service using docker conatiner

```bash
docker run --env AWS_ACCESS_KEY_ID=<access_key_id> --env AWS_SECRET_ACCESS_KEY=<secret_access_key> cyderes-app:1.1
```

### 4. OR run service using docker-compose.yml

```bash
AWS_ACCESS_KEY_ID=<access_key_id> AWS_SECRET_ACCESS_KEY=<secret_access_key> docker compose -f ./build/docker/docker-compose.yml up
```

### 5. Observe Logs

```bash
docker-compose logs -f service
```

---

## Dockerfile

```dockerfile
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ingest main.go

FROM alpine:3.17
WORKDIR /root/
COPY --from=builder /app/ingest .
CMD ["./ingest"]
```

---

## docker-compose.yml

```yaml
version: '3.8'
services:
  dynamodb:
    image: amazon/dynamodb-local
    container_name: dynamodb_local
    ports:
      - "8000:8000"
  service:
    build: .
    environment:
      - AWS_REGION=us-east-1
      - AWS_ENDPOINT_URL=http://dynamodb:8000
    depends_on:
      - dynamodb
```

---

## Configuration

* **AWS\_ENDPOINT\_URL**: override endpoint for local DynamoDB
* **AWS\_REGION**: AWS region
* **TABLE\_NAME**: DynamoDB table (default `Logs`)

---

## Testing

### Unit Tests

File: `transformer/transformer_test.go`

```go
package transformer
import (
  "testing"
  "time"
)
func TestTransformData(t *testing.T) {
  raw := []byte(`[{"userId":1,"id":2,"title":"t","body":"b"}]`)
  out, err := TransformData(raw)
  if err != nil {
    t.Fatalf("unexpected error: %v", err)
  }
  if len(out) != 1 {
    t.Errorf("expected 1 record, got %d", len(out))
  }
  if out[0].Source != "placeholder_api" {
    t.Errorf("unexpected source: %s", out[0].Source)
  }
}
```

### Integration Tests

* Use Docker Compose to bring up DynamoDB Local
* In `tests/integration_test.go`, point AWS\_ENDPOINT\_URL at localhost:8000, call `StoreToDynamoDB`, then use AWS SDK to `Scan` the table and assert record count.

---

## Documentation

* **API Endpoint**: none (CLI application) by default
* **Transformation Logic**: see `internal/transformer/transformer.go`
* **DB Schema**:

  * Partition key: `id` (Number)
  * Attributes: `userId`, `title`, `body`, `ingested_at`, `source`

---

## Trade-offs

* **DynamoDB vs S3**: Chose DynamoDB for per-item writes and querying; S3 would require batch files and ETL.
* **Localstack vs DynamoDB Local**: Used DynamoDB Local for simplicity; localstack supports more AWS services.
* **Exponential backoff**: simple linear backoff implemented; a library (e.g., backoff) could provide jitter.

---

## Hardest Parts

* Handling AWS SigV4 signatures manually (avoided by using AWS SDK)
* Testing integration with a real database (solved via DynamoDB Local)
* Managing environment-specific configs in Docker Compose

---

## Future Improvements

* Batch writes with `BatchWriteItem` for performance
* Add REST API (Gin or Echo) to expose ingested data
* Implement CI/CD with GitHub Actions for build, test, and deploy
* Add monitoring and alerting (CloudWatch, Prometheus)

---

## Bonus: CI/CD with GitHub Actions

```yaml
name: CI
on: [push]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.20'
      - name: Install dependencies
        run: go mod download
      - name: Run tests
        run: go test ./...
      - name: Build Docker image
        run: docker build . -t data-ingest:latest
```

---

*End of README*
