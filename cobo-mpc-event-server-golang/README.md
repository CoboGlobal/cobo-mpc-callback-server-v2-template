# cobo-mpc-event-server-golang

## Overview

This is a Golang implementation of the TSS Node event server.
When `event` are configured in the TSS Node, the TSS Node will generate different events during its operation and send them to the event server.

It provides a basic template for handling TSS Node event and can be customized according to specific business requirements.

## Requirements

- Go 1.23.1

## Deployment Steps

### 1. Clone the Repository

```bash
git clone https://github.com/CoboGlobal/cobo-mpc-callback-server-v2-template.git
cd cobo-mpc-callback-server-v2-template/cobo-mpc-event-server-golang

```

### 2. Build

```bash
go build -trimpath -o build/bin/tss-node-event-server cmd/main.go
```

### 3. Configure Keys

Place the following key files in the project root directory:

- configs/tss-node-event-pub.key (TSS Node's RSA event public key)

### 4. Start the Server

```bash
./build/bin/tss-node-event-server
```

The server will start on port 11030 by default.

## Testing

### 1. Health Check

```bash
curl http://127.0.0.1:11030/ping
```

### 2. Integration Testing

To test the complete workflow with TSS Node:

- Ensure your event server is running
- Configure and start your TSS Node
- Send event through TSS Node to the event server

For detailed TSS Node setup, refer to the document.

## Important Notes

### Basic Implementation

This template implements only the basic server structure.
Implement your own handle logic based on your business requirements.

### Dependencies

Refer to the SDK documentation for detailed parameter definitions.
