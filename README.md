# Herd: A Secure NoSQL In-Memory Key-Value Store with Persistence

<div align="center">
  <img src="./banner.png" alt="Herd Logo">
  <br>
  <img src="https://img.shields.io/github/actions/workflow/status/defoeam/Herd/docker.yaml" alt="Build Status">
</div>

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Setup](#setup)
4. [Building and Running](#building-and-running)
5. [Usage](#usage)
6. [Contributing](#contributing)
7. [License](#license)


## Overview

Herd is an open-source, cloud-native, in-memory NoSQL key-value store (KVS) optimized for speed, security, and persistence. Designed as a research-friendly alternative to proprietary solutions like Redis or Memcached, Herd combines:

- **Write-Behind Logging** for data durability.
- **TLS-based security** for encrypted communication.
- **gRPC APIs** for seamless integration.

Herd provides:

- High throughput (>5000 operations/second with 1000 concurrent clients).
- Minimal latency overhead, even with persistence and encryption enabled.
- Modular design for extensibility and research applications.

Written in the memory-safe Go programming language, Herd also includes a Python client library for integration with machine learning and data pipelines.

## Features

- **Key-Value Cache Functionality:** Efficient retrieval and storage of key-value pairs.
- **Transaction Logging with Snapshotting:** Ensures data durability and faster recovery.
- **Secure Communication:** Encrypted client-server interactions using TLS.
- **gRPC API:** Enables easy interaction with support for extensibility.
- **Python Client Library:** Simplifies integration into Python workflows.
- **Dockerized Deployment:** Streamlined setup and portability via Docker Compose.



## Setup

### Prerequisites

- **Docker and Docker Compose:** Install [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/).
- **Make:** Ensure `make` is installed on your system.

### Generating TLS Certificates

To generate the necessary TLS certificates and key pairs, run:

```bash
make
```

The generated files will be located in the `./certs/` directory.



## Building and Running

### Building with Docker Compose

To build the Herd project:

```bash
docker compose build
```

### Running the Server Locally

To start Herd:

```bash
docker compose up -d
```

To stop the server:

```bash
docker compose down
```



## Usage

### Python Client Example

Here’s how you can interact with Herd using the Python client found here: https://github.com/defoeam/Herd-python.

```python
from herd_client import HerdClient

# Initialize the client
client = HerdClient("localhost:7878", secure=True)

# Set a key-value pair
client.set("key1", "value1")

# Get a value by key
value = client.get("key1")
print("Retrieved value:", value)

# Delete a key
client.delete("key1")
```

### gRPC API Operations

Herd supports the following gRPC operations:

- **SET:** Add or update a key-value pair.
- **GET:** Retrieve the value for a key.
- **DELETE:** Remove a key-value pair.
- **GETALL:** Retrieve all key-value pairs.
- **DELETEALL:** Clear the entire store.



## Architecture

Herd’s architecture is designed for modularity and performance:

1. **Write-Behind Logging (WBL):** Ensures persistence by asynchronously logging changes, minimizing I/O bottlenecks.
2. **Snapshotting:** Periodically saves the database state to optimize recovery processes.
3. **Transport Layer Security (TLS):** Provides encrypted client-server communication.
4. **gRPC-based API:** Allows low-latency and language-agnostic integration.
5. **Python Client Library:** Simplifies interaction with Python-based applications.


## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or fix.
3. Submit a pull request with a detailed description of your changes.

For major changes, please open an issue to discuss your proposal.


## License

Herd is open-source software licensed under the [MIT License](./LICENSE).


