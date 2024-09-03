# ETH Collector

This tool is designed to collect all ETH(or other EVM chain's native token) from multiple addresses to a single receiver address.

It's useful for consolidating funds from multiple wallets or for automated fund management.

## Features

- Connect to any EVM compatible network
- Transfer all ETH (minus gas costs) from multiple sender addresses to a single receiver address
- Configure via command-line flags or a JSON config file
- Automatic gas price estimation
- Detailed logging of transfer operations

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/huahuayu/eth-collector.git
   cd eth-collector
   ```

2. Build the project:

   ```bash
   go build -o eth-collector
   ```

## Usage

You can run the tool using command-line flags or a config file.

### Using Command-line Flags

```bash
./eth-collector -rpc="https://mainnet.infura.io/v3/YOUR-PROJECT-ID" -receiver="0x742d35Cc6634C0532925a3b844Bc454e4438f44e" -sender="private_key1" -sender="private_key2"
```

### Using a Config File

1. Create a JSON config file (e.g., `config.json`):

   ```json
   {
     "rpc": "https://mainnet.infura.io/v3/YOUR-PROJECT-ID",
     "receiverAddress": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
     "senderPrivateKeys": ["private_key1", "private_key2"]
   }
   ```

2. Run the tool with the config file:
   ```bash
   ./eth-collector -config="config.json"
   ```

## Configuration Options

- `-rpc`: RPC URL
- `-receiver`: Address of the receiver
- `-sender`: Private key of a sender (can be specified multiple times)
- `-config`: Path to the JSON config file

## Security Reminder

Never share your private keys or include them in public repositories.
