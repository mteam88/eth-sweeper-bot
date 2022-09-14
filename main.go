package main

import (
	"fmt"
    "context"
    "log"
    "crypto/ecdsa"
    "math/big"
    "time"
    "os"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/params"
    "github.com/joho/godotenv"
)

func etherToWei(val *big.Int) *big.Int {
	return new(big.Int).Mul(val, big.NewInt(params.Ether))
}

func weiToEther(val *big.Int) *big.Int {
	return new(big.Int).Div(val, big.NewInt(params.Ether))
}

func checkerr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    godotenv.Load(".envrc", ".env")
    
    fmt.Println(os.Getenv("NODE_ENDPOINT"))
    
    client, err := ethclient.Dial(os.Getenv("NODE_ENDPOINT"))
    checkerr(err)

    privateKey, err := crypto.HexToECDSA(os.Getenv("TARGET_PRIVATE_KEY"))
    
    checkrate := time.Duration(15000000000)
    for {
        time.Sleep(checkrate)
        // Get the balance of an account
        publicKey := privateKey.Public()
        publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
        if !ok {
            log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
        }
    
        fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
        nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
        checkerr(err)
        
        balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
        checkerr(err)
        
        fmt.Println("Nonce: ", nonce)
        fmt.Println("Account balance:", balance)
    
        
        gasLimit := uint64(21000)          // in units
        gasPrice, err := client.SuggestGasPrice(context.Background())
        checkerr(err)
        
        gasLimitBigInt := new(big.Int).SetUint64(gasLimit)
        gasExpense := new(big.Int)
        gasExpense.Mul(gasLimitBigInt, gasPrice)
        fmt.Println("gasExpense: ", gasExpense)
    
        value := new(big.Int).Sub(balance, gasExpense)
        fmt.Println("value: ", value)

        if value.Int64() >= 0 {
            fmt.Println("Valid balance: Initialize transaction")
        
            toAddress := common.HexToAddress(os.Getenv("HQ_ADDRESS"))
            var data []byte
            tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
        
            chainID, err := client.NetworkID(context.Background())
            checkerr(err)
        
            signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
            checkerr(err)
        
            err = client.SendTransaction(context.Background(), signedTx)
            checkerr(err)
        
            balance, err = client.BalanceAt(context.Background(), fromAddress, nil)
            checkerr(err)
            
            fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
        
            txHash := signedTx.Hash()
            sleeptime := time.Duration(30000000000)
            for {
                time.Sleep(sleeptime)
                _, isPending, err := client.TransactionByHash(context.Background(), txHash)
                checkerr(err)
                fmt.Println("Waiting for transaction to finish...")
                if !isPending {
                    fmt.Println("TX done.")
                    break
                }
            }
        } else {
            fmt.Println("Balance too low. Waiting...")
        }
    }
}