# Going_Cache Documentation

## Overview

Going_Cache is a distributed caching system implemented in Go that provides:

-   Distributed key-value storage
-   Consistent hashing for data sharding
-   Data replication across nodes
-   LRU eviction policy
-   TTL-based cache expiration

## Prerequisites

-   Ubuntu Server/Desktop
-   Go 1.22.3 or later
-   Git

## Installation

1. Install Go on Ubuntu:

```bash
sudo apt update
sudo apt install golang-go
```

2. Clone the repository:

```bash
git clone https://github.com/StaphoneWizzoh/Going_Cache.git
cd Going_Cache
```

3. Build the project:

```bash
go build
```

## Configuration

The cache server accepts two command-line flags:

-   `--port`: HTTP server port (default ":8080")
-   `--peers`: Comma-separated list of peer addresses

## Running the Cache Cluster

1. Start the first node:

```bash
./Going_Cache --port=:8080
```

2. Start additional nodes with peer information:

```bash
./Going_Cache --port=:8081 --peers="http://localhost:8080"
./Going_Cache --port=:8082 --peers="http://localhost:8080,http://localhost:8081"
```

## API Usage

### Set a Value

```bash
curl -X POST http://localhost:8080/set \
  -H "Content-Type: application/json" \
  -d '{"key": "mykey", "value": "myvalue"}'
```

### Get a Value

```bash
curl http://localhost:8080/get?key=mykey
```

## Architecture

### Components

1. **Cache Server** (server.go)

-   Handles HTTP requests for get/set operations
-   Manages data replication
-   Routes requests to appropriate nodes

2. **Cache** (cache.go)

-   Implements LRU eviction
-   Handles TTL expiration
-   Thread-safe operations

3. **Hash Ring** (hashring.go)

-   Implements consistent hashing
-   Manages node distribution
-   Handles node addition/removal

### Key Features

1. **Data Sharding**

-   Uses consistent hashing to distribute keys
-   Minimizes key redistribution when nodes change
-   Ensures even data distribution

2. **Replication**

-   Automatically replicates data across nodes
-   Handles peer synchronization
-   Prevents replication loops

3. **Eviction Policies**

-   LRU (Least Recently Used) eviction
-   Time-based expiration (TTL)
-   Configurable cache capacity

## Monitoring

The server logs important events including:

-   Server start
-   Request forwarding
-   Replication failures
-   Connection issues

## Limitations

-   No persistence
-   Basic error handling
-   Simple replication strategy
-   Fixed TTL of 1 hour for cached items

## Best Practices

1. **Deployment**

-   Run nodes on separate machines for redundancy
-   Configure appropriate timeout values
-   Monitor system resources

2. **Scaling**

-   Add nodes gradually
-   Ensure network connectivity between nodes
-   Monitor cache hit rates

3. **Maintenance**

-   Regular monitoring of logs
-   Check for failed nodes
-   Monitor memory usage

## Health Checks

Monitor node health using:

```bash
curl http://localhost:8080/get?key=health
```

For more information, refer to the

README.md

or raise issues on the GitHub repository.
