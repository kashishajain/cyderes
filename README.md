
## Overview

This service fetches posts from the JSONPlaceholder API, transforms each record by adding metadata (`ingested_at`, `source`), and stores them in Amazon DynamoDB. It is containerized with Docker, tested end-to-end (unit and integration), and can be extended with CI/CD pipelines.

---

## Architecture

1. **Fetcher**: HTTP client with retry and fetches data from `https://jsonplaceholder.typicode.com/posts`.
2. **Transformer**: Adds UTC timestamp `ingested_at` and static `source` field.
3. **Storage**: Writes each transformed item into DynamoDB table `cyderes_api_logs`.
4. **Containerization**: `Dockerfile` OR `docker-compose.yml` spins up the service.
5. **Testing**:

   * Unit tests for transformation logic.
   * Unit tests for fetch_data using httptest.
   * End to end Integration tests against DynamoDB Local via Docker.


## Storage Justification

* **Amazon DynamoDB**: A fully managed NoSQL database optimized for high-throughput writes and key-value/JSON data. It auto-scales, has a free tier, and the AWS SDK for Go v2 provides seamless integration. It fits moderately structured logs with simple primary-key lookups. Trade-offs: cost can rise under heavy read/write loads, but our moderate volume keeps it economical.

---

## Setup & Running Locally

### Prerequisites

* Docker & Docker Compose
* Go 1.20+
* Create IAM role configured with required DynamoDB policies (AmazonDynamoDBFullAcess). 
* AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
* AWS CLI (for `aws configure`, optional)

### 1. Clone & Build

```bash
git clone <repo-url>
cd cyderes
```

### 2. Export Env variables 

```bash
export ACCESS_KEY_ID=<access_key_id>
export SECRET_ACCESS_KEY=<secret_access_key>
```

### 3. Build docker image & Start service using docker conatiner

```bash
docker build -f build/docker/Dockerfile -t cyderes-app:1.1 .
```
### If AWS credentials are configured using env variable
```bash
docker run --env AWS_ACCESS_KEY_ID=$ACCESS_KEY_ID --env AWS_SECRET_ACCESS_KEY=$SECRET_ACCESS_KEY cyderes-app:1.1
```
### If AWS credentials are configured using AWS CLI
```bash
docker run -v ~/.aws:/root/.aws cyderes-app:1.1
```
### OR 
### Run service using docker-compose.yml

```bash
AWS_ACCESS_KEY_ID=$ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY=$SECRET_ACCESS_KEY docker compose -f ./build/docker/docker-compose.yml up
```

---


## Documentation

* **API Endpoint**: none (CLI application) by default
* **Transformation Logic**: see `/transformer/transform_data.go`
* **DB Schema**:

  * Partition key: `userID` (Number)
  * Sort key: `id`
  * Attributes: `userId`, `title`, `body`, `ingested_at`, `source`

---

## Trade-offs

* **DynamoDB vs S3**: Chose DynamoDB for per-item writes and querying; S3 would require batch files and ETL.
* **DynamoDB Local**: Used DynamoDB Local for end to end integration test.

---

## Hardest Parts

* Testing integration with a real database (solved via DynamoDB Local)
* Managing environment-specific configs in Docker Compose

---

## Future Improvements

* Add REST API to expose ingested data
* Add monitoring and alerting (CloudWatch, Prometheus)
* Updating github actions workflow to push the image on registry and deploy it to EC2 or lambda function.

---

End of README