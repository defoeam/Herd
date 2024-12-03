<div align="center">

![logo image](./banner.png)
![build](https://img.shields.io/github/actions/workflow/status/defoeam/Herd/docker.yaml)

</div>

Todo: Summary

## Features
- Basic key:value cache functionality
- Transaction logging with Snapshotting
- Transport Layer Security (TLS)

## Setup
To setup the TLS certificates, run:
```bash
make
```
The certificates and key pairs will be in the ./certs/ directory.

## Building with Docker Compose
To build:
```bash
docker compose build
```

To start Herd locally, simply run:

```bash
docker compose up -d
```

and then to bring down the server, run:

```bash
docker compose down
```
