# cobo-mpc-callback-server-v2-golang

## Overview

This is a Golang implementation of the TSS Node callback server. It provides a basic template for handling TSS Node requests and can be customized according to specific business requirements.

## Requirements

- Go 1.23.1

## Deployment Steps

### 1. Clone the Repository
```bash
git clone https://github.com/CoboGlobal/cobo-mpc-callback-server-v2-template.git
cd cobo-mpc-callback-server-v2-template/cobo-mpc-callback-server-v2-golang

```

### 2. Build
```bash 
go build -trimpath -o build/bin/tss-node-callback-server cmd/main.go
```
### 3. Configure Keys

Place the following key files in the project root directory:

- configs/tss-node-callback-pub.key (TSS Node's RSA public key)
- configs/callback-server-pri.pem (Callback server's RSA private key)

### 4. Start the Server
```bash 
./build/bin/tss-node-callback-server
```

The server will start on port 11020 by default.


## Testing

### 1. Health Check

```bash
curl http://127.0.0.1:11020/ping
```

### 2. Integration Testing

To test the complete workflow with TSS Node:

- Ensure your callback server is running
- Configure and start your TSS Node
- Send requests through TSS Node to the callback server

For detailed TSS Node setup, refer to the [Callback Server Overview](https://www.cobo.com/developers/v2/guides/mpc-wallets/server-co-signer/callback-server-overview).

## Important Notes

### Basic Implementation

This template implements only the basic server structure.
All requests are allowed by default.
Implement your own callback logic based on your business requirements.


### Dependencies

The `extra_info` risk control parameter structure is defined in [cobo-waas2-go-sdk](https://github.com/CoboGlobal/cobo-waas2-go-sdk)
Refer to the SDK documentation for detailed parameter definitions.
