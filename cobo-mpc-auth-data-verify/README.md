# Cobo MPC Auth Data Verify

A comprehensive Go library for verifying transaction approval data in Cobo's MPC (Multi-Party Computation) system. This library provides tools to validate transaction approval information, verify digital signatures, and ensure data integrity using WaaS2 SDK and Jinja2 templating.

## Overview

This project implements a complete verification pipeline for Cobo's MPC transaction approval system. It fetches transaction and approval data from WaaS2 API, builds approval messages using Jinja2 templates, and verifies cryptographic signatures to ensure the authenticity and integrity of transaction approvals.

## Architecture

### Core Components

1. **Validator** (`validator/`): Core message verification logic
   - `validator.go`: Main validation interface and AuthData structure
   - `statement.go`: Jinja2 template processing and message building
   - `signature.go`: Cryptographic signature verification

2. **WaaS2** (`waas2/`): WaaS2 SDK integration
   - `client.go`: WaaS2 API client implementation
   - `waas2.go`: Transaction and approval detail building logic
   - `validator.go`: Transaction approval detail validation

3. **Example** (`example/`): Complete usage example
   - `main.go`: Demonstrates the full verification workflow

## Workflow

The verification process follows these steps:

### 1. Configuration

Set up the required configuration parameters:

```go
var (
    pubkeyWhitelist = []string{
        "your_cobo_guard_public_key_1",
        "your_cobo_guard_public_key_2",
    }
    
    apiSecret = "your_waas2_api_secret"
    env       = coboWaas2.DevEnv  // or coboWaas2.ProdEnv
)
```

**Required Configuration:**
- **Cobo Guard Public Key Whitelist**: List of trusted public keys for signature verification
- **WaaS2 API Secret**: API credentials for accessing WaaS2 services
- **Environment**: Development or production environment
- **Transaction ID**: The transaction to be verified

### 2. Build Approval Information

Create a WaaS2 client and build transaction approval details:

```go
// Initialize WaaS2 client
client := waas2.NewClient(apiSecret, env)
waas2Client := waas2.NewWaas2(client)

// Build transaction and approval details
transactionIds := []string{"your_transaction_id"}
txApprovalDetails, err := waas2Client.Build(context.Background(), transactionIds)
if err != nil {
    panic(fmt.Errorf("failed to build transaction approval details: %w", err))
}
```

This step:
- Fetches detailed transaction information from WaaS2 API
- Retrieves transaction approval details and user information
- Downloads and caches all required approval templates
- Builds comprehensive `TxApprovalDetail` structures

### 3. Verify Approval Information

Validate the approval information using the validator:

```go
config := waas2.Config{
    PubkeyWhitelist: pubkeyWhitelist,
}

for _, txApprovalDetail := range txApprovalDetails {
    // Verify transaction approval detail
    validator := waas2.NewTxApprovalDetailValidator(txApprovalDetail, &config)
    err = validator.Verify(context.Background())
    if err != nil {
        panic(fmt.Errorf("failed to verify transaction approval detail: %w", err))
    }
}
```

This verification process includes:

#### 3.1 Build Business Data
- Merges transaction details with user approval information
- Creates comprehensive business data for template processing

#### 3.2 Build Approval Message
- Uses Jinja2 templates to generate approval messages
- Applies custom filters and methods for data transformation
- Ensures consistent message formatting

#### 3.3 Verify Approval Signature
- Validates cryptographic signatures using Cobo Guard public keys
- Verifies message integrity and authenticity
- Checks approval results and thresholds

### 4. Transaction Consistency Verification

After successful verification, compare transaction details with TSS callback data to ensure consistency:

```go
// Verify txApprovalDetail (transaction and approval detail)
// txApprovalDetail and tss callback data are matched
```

This final step ensures that the verified transaction data matches the actual TSS callback transaction data.

## Dependencies

- **Go 1.24.4+**: Required Go version
- **Cobo WaaS2 Go SDK**: For accessing WaaS2 API services
- **Gonja v2**: Jinja2 template engine for Go
- **Google Go-CMP**: For data comparison utilities
- **Testify**: Testing framework

## Installation

```bash
go mod tidy
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/CoboGlobal/cobo-mpc-callback-server-v2-template/cobo-mpc-auth-data-verify/waas2"
    coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

func main() {
    // Configuration
    pubkeyWhitelist := []string{"your_public_key"}
    apiSecret := "your_api_secret"
    env := coboWaas2.DevEnv
    transactionIds := []string{"your_transaction_id"}
    
    // Initialize client
    client := waas2.NewClient(apiSecret, env)
    waas2Client := waas2.NewWaas2(client)
    
    // Build approval details
    txApprovalDetails, err := waas2Client.Build(context.Background(), transactionIds)
    if err != nil {
        panic(err)
    }
    
    // Verify approval details
    config := waas2.Config{PubkeyWhitelist: pubkeyWhitelist}
    for _, txApprovalDetail := range txApprovalDetails {
        validator := waas2.NewTxApprovalDetailValidator(txApprovalDetail, &config)
        if err := validator.Verify(context.Background()); err != nil {
            panic(err)
        }
        fmt.Println("Transaction verification successful!")
    }
}
```

## Testing

Run the test suite:

```bash
go test ./...
```

## Security Considerations

- Always validate public keys against your whitelist
- Use secure storage for API secrets
- Verify transaction data consistency with TSS callbacks
- Monitor for signature verification failures

## License

This project is part of the Cobo MPC Callback Server template system.

## Support

For technical support and questions, please refer to the Cobo documentation or contact the development team.
