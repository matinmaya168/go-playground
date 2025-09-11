package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

// Address represents a wallet address with private and public keys and balance
type Address struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Balance    float64
}

// Wallet represents a collection of addresses
type Wallet struct {
	Addresses map[string]*Address
}

// Transaction represents a transfer operation
type Transaction struct {
	Sender    string
	Receiver  string
	Amount    float64
	Signature []byte
}

// Block represents a single block in the blockchain
type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	PreviousHash string
	Hash         string
}

// Blockchain represents the blockchain structure
type Blockchain struct {
	Chain        []Block
	PendingTxs   []Transaction
	HashFunction func([]byte) string
}

// NewWallet initializes a new wallet
func NewWallet() *Wallet {
	return &Wallet{
		Addresses: make(map[string]*Address),
	}
}

// CreateAddress generates a new address with public/private key pair
func (w *Wallet) CreateAddress() string {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.PublicKey
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	hash := sha256.Sum256(publicKeyBytes)
	address := hex.EncodeToString(hash[:])
	w.Addresses[address] = &Address{
		PrivateKey: privateKey,
		PublicKey:  &publicKey,
		Balance:    0.0,
	}
	return address
}

// Transfer handles the transfer of found between address
func (w *Wallet) Transfer(from, to string, amount float64) (*Transaction, error) {
	fromAddr, exists := w.Addresses[from]
	if !exists {
		return nil, fmt.Errorf("sender address %s does not exist", from)
	}
	toAddr, exists := w.Addresses[to]
	if !exists {
		return nil, fmt.Errorf("receiver address %s does not exist", to)
	}
	if fromAddr.Balance < amount {
		return nil, fmt.Errorf("insufficient balance in sender address %s", from)
	}

	// Create transaction
	tx := Transaction{
		Sender:   from,
		Receiver: to,
		Amount:   amount,
	}

	// Sign transaction
	txData := []byte(from + to + fmt.Sprintf("%f", amount))
	hash := sha256.Sum256(txData)
	signature, err := ecdsa.SignASN1(rand.Reader, fromAddr.PrivateKey, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}
	tx.Signature = signature

	fromAddr.Balance -= amount
	toAddr.Balance += amount

	return &tx, nil
}

// VerifyTransaction verifies the transaction signature
func VerifyTransaction(tx Transaction) bool {
	addr, err := hex.DecodeString(tx.Sender)
	if err != nil {
		return false
	}
	publicKeyBytes := addr
	hash := sha256.Sum256([]byte(tx.Sender + tx.Receiver + fmt.Sprintf("%f", tx.Amount)))
	var publicKey ecdsa.PublicKey
	// Assuming public key store in Address; in practice, fetch from wallet
	curve := elliptic.P256()
	x, y := elliptic.Unmarshal(curve, publicKeyBytes)
	if x == nil {
		return false
	}
	publicKey = ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	return ecdsa.VerifyASN1(&publicKey, hash[:], tx.Signature)
}

// NewBlockChain initializes a new blockchain
func NewBlockChain() *Blockchain {
	return &Blockchain{
		Chain:        []Block{},
		PendingTxs:   []Transaction{},
		HashFunction: calculateHash,
	}
}

// calculateHash computes the hash of a block
func calculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// AddTransaction adds a transaction to the pending pool
func (bc *Blockchain) AddTransaction(tx Transaction) bool {
	if !VerifyTransaction(tx) {
		return false
	}
	bc.PendingTxs = append(bc.PendingTxs, tx)
	return true
}

// Main function to demonstrate transfer operations
func main() {
	// Create a new wallet
	wallet := NewWallet()

	// Create two addresses
	addr1 := wallet.CreateAddress()
	addr2 := wallet.CreateAddress()

	// Initialize balance for addr1
	wallet.Addresses[addr1].Balance = 100.0

	// Create a blockchain
	bc := NewBlockChain()

	// Perform a transfer
	tx, err := wallet.Transfer(addr1, addr2, 50.0)
	if err != nil {
		log.Fatal(err)
	}

	// Add transaction to blockchain
	if bc.AddTransaction(*tx) {
		fmt.Printf("Transaction added: %v\n", tx)
	} else {
		fmt.Println("Transaction verification failed")
	}

	// Print balances
	fmt.Printf("Balance of %s: %.2f\n", addr1, wallet.Addresses[addr1].Balance)
	fmt.Printf("Balance of %s: %.2f\n", addr2, wallet.Addresses[addr2].Balance)
}
