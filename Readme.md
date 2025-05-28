# Go Proxy Server with LRU Cache

This project implements a multithreaded HTTP proxy server in Go, featuring an in-memory Least Recently Used (LRU) cache. The server acts as an intermediary between clients and target web servers, enhancing performance by caching frequently accessed content.

## Features

- **HTTP Proxying:** Forwards client HTTP requests to target web servers and returns their responses.
- **Multithreaded (Goroutines):** Handles multiple client requests concurrently using Go's goroutines, ensuring high throughput.
- **LRU Cache:** Stores responses from target servers in an in-memory cache. If a subsequent request for the same resource is received, the server checks its cache first.
  - **Cache Hit:** If the response is found in the cache, it's served directly, reducing latency and external server load.
  - **Cache Miss:** If the response is not in the cache, the proxy fetches it from the target server, stores it in the cache, and then serves it to the client.
- **Least Recently Used (LRU) Eviction:** When the cache reaches its maximum capacity, the least recently accessed items are automatically removed to make space for new ones.
- **Concurrency Safe Cache:** The LRU cache is protected with mutexes to ensure safe access from multiple goroutines.

## How it Works

1. A client sends an HTTP request to the proxy server.
2. The proxy server extracts the target URL from the client's request.
3. It generates a unique cache key based on the request URL.
4. The server first checks if the response for this key exists in its LRU cache:
   - **If found (cache hit):** The cached response is immediately sent back to the client.
   - **If not found (cache miss):**
     - The proxy server makes a new request to the actual target web server.
     - Upon receiving the target server's response, the proxy stores it in its LRU cache.
     - Finally, the response is forwarded to the client.
5. All these operations are handled concurrently for multiple clients using Go's lightweight goroutines.

## Getting Started

### Prerequisites

- Go (version 1.16 or higher recommended)

### Installation

Clone the repository (or create the project structure):

```sh
git clone https://github.com/your-username/go-proxy-server.git # Replace with your actual repo
cd go-proxy-server
```

(If you are building from scratch, simply create a directory and `go mod init` inside it.)

Initialize Go module (if not already done):

```sh
go mod init your_proxy_project # Use a meaningful name for your project
```

### Running the Server

Build the application:

```sh
go build -o proxy-server .
```

Run the server:

```sh
./proxy-server
```

By default, the server will listen on [http://localhost:8080](http://localhost:8080).

### Configuration (Future)

- **Proxy Port:** You can configure the listening port (e.g., via command-line flags or environment variables).
- **Cache Size:** The maximum number of items the LRU cache can hold will be configurable.

## Usage

Once the proxy server is running, you can configure your browser or application to use it as an HTTP proxy.

**Example (using curl):**

To proxy a request to `http://example.com` through your server running on `localhost:8080`:

```sh
curl -x http://localhost:8080 http://example.com
```

You will see the response from `example.com` served through your proxy. Subsequent requests to `http://example.com` should be faster as they will be served from the cache.

## Project Structure (Expected)

```
.
├── main.go               # Main application logic, server setup, and proxy handler
├── cache/                # Directory for cache implementation
│   └── lru.go            # LRU cache data structure and methods
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
└── README.md             # This file
```

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is open-source and available under the MIT License.
