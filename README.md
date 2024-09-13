# top-k-system-go

Example **GoLang** implementation of a system that will allow querying top k most popular items for a given time period (1 hour, 1 day, 1 month or all time). The implementation will use a precise calculation rather than doing a probabilistic/approximate calculation (e.g using [count-min sketch](https://en.wikipedia.org/wiki/Count–min_sketch)).

The design is inspired by the [Top-K Youtube Videos](https://www.hellointerview.com/learn/system-design/answer-keys/top-k).

This is a work in progress

## System Design

Core Functional Requirements:
* It should be possible to query top K items (max 1000) for a given time window.
* Time windows are: last hour, last day, last month and all time.

Core Non Functional Requirements:
* Eventual consistency on querying top K items (up to a minute tolerance)
* The system should support large volume of events that are incrementing popularity of the item (e.g view events or similar)
* The system should support large volume of items
* The topK query should be fast (100ms or less)

## Project Setup

Please have the following tools installed: 
* [direnv](https://github.com/direnv/direnv) 
* [gobrew](https://github.com/kevincobain2000/gobrew#install-or-update)

Install/Update dependencies: 
```sh
# Install
go mod download
make tools

# Update:
go get -u ./... && go mod tidy
```

### Lint and Tests

Run all lint and tests:
```bash
make lint
make test
```

Run specific tests:
```bash
# Run once
go test -v ./service/pkg/api/http/v1controllers/ --run TestHealthCheckController

# Run same test multiple times
# This is useful for tests that are flaky
go test -v -count=5 ./service/pkg/api/http/v1controllers/ --run TestHealthCheckController

# Run and watch
gow test -v ./service/pkg/api/http/v1controllers/ --run TestHealthCheckController
```
### Run local API server:

```bash
# Regular mode
go run ./cmd/service/

# Watch mode (double ^C to stop)
gow run ./cmd/service/
```