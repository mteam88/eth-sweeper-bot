# goethwithdrawbot
Constantly attempt to withdraw all ETH in provided wallet.

## Installation
Clone this repo with 
```
git clone https://github.com/mteam88/goethwithdrawbot
cd goethwithdrawbot
```
You must have `go` installed.

## Usage
Run
```go run main.go```
in `goethwithdrawbot` directory.
IMPORTANT: Complete the configuration steps below before running.

## Configuration
Set the Environment variables below manually or place them in a `.envrc` file.

NODE_ENDPOINT = <YOUR NODE PROVIDER HERE> // Usually a local node (ganache, geth) or a hosted node (infura, alchemy.io)
TARGET_PRIVATE_KEY = <YOUR TARGET PRIVATE KEY HERE>
HQ_ADDRESS = <YOUR PERSONAL ADDRESS HERE> // This is where the ETH will be sent.
