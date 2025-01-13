# Cobo MPC Callback Server V2 Template

## Overview

This project serves as a template for implementing TSS Node callback servers in multiple programming languages.
It is a crucial component for using MPC wallets and deploying server co-signers with the TSS Node callback mechanism.
Clients can use this template as a reference to build their own TSS Node callback servers.

For detailed information about TSS Node callback mechanism, please visit  [Cobo Developer Documentation](https://www.cobo.com/developers/v2/guides/mpc-wallets/server-co-signer/callback-server-overview).

## Key Features

The template demonstrates the implementation of essential functionalities required for a TSS Node callback server:

1. HTTP endpoints for receiving task requests
2. Processing of JWT-signed messages from the TSS Node
3. Custom risk control logic implementation
4. Signed response handling back to the TSS Node

## Available Implementations

The template is available in three programming languages:

### [Go Implementation](./cobo-mpc-callback-server-v2-golang/README.md)
- Located in: `cobo-mpc-callback-server-v2-golang`
- Provides Go-specific implementation and setup instructions

### [Java Implementation](./cobo-mpc-callback-server-v2-java/README.md)
- Located in: `cobo-mpc-callback-server-v2-java`
- Includes Java-specific implementation and configuration details

### [Python Implementation](./cobo-mpc-callback-server-v2-python/README.md)
- Located in: `cobo-mpc-callback-server-v2-python`
- Contains Python-specific implementation and usage guidelines

## Getting Started

Please refer to the README in each language-specific directory for detailed setup and implementation instructions. Each implementation includes:

- Installation requirements
- Configuration steps
- Usage examples
- Implementation details

## Additional Resources

For comprehensive understanding and implementation details, refer to:
- [Official Documentation](https://www.cobo.com/developers/v2/guides/mpc-wallets/server-co-signer/callback-server-overview)
- Language-specific README files in respective directories


## Comparison with cobo-mpc-callback-server-examples

There are two different template repositories available:

- [cobo-mpc-callback-server-examples](https://github.com/CoboGlobal/cobo-mpc-callback-server-examples)
    - Uses WaaS SDK (v1)
 
- This repository (cobo-mpc-callback-server-v2-template)
    - Uses WaaS2 SDK

Please choose the appropriate template based on your WaaS SDK version.

## Contributing

We welcome contributions to improve the template. Please feel free to submit issues and pull requests.
