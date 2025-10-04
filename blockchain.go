package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

// Block represents a single block in the blockchain containing metadata and a list of transactions.
type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

// NewBlock constructs a new Block with the given nonce, previous hash, and transactions.
func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		nonce:        nonce,
		previousHash: previousHash,
		timestamp:    time.Now().UnixNano(),
		transactions: transactions,
	}
}

// Blockchain holds the chain of blocks and a pool of pending transactions.
type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

// NewBlockchain initializes a new Blockchain with a genesis block.
func NewBlockchain(blockchainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	return bc
}

// CreateBlock creates a new block from the current transaction pool and appends it to the chain.
func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

// LastBlock returns the most recently added block in the chain.
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// Print outputs the block details and all contained transactions to stdout.
func (b *Block) Print() {
	fmt.Printf("timestamp: %d\n", b.timestamp)
	fmt.Printf("nonce: %d\n", b.nonce)
	fmt.Printf("previousHash: %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

// Print outputs the entire blockchain to stdout in a readable format.
func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

// AddTransaction creates a new transaction and adds it to the transaction pool.
func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

// CopyTransactionPool creates a deep copy of the current transaction pool and returns it as a slice of transactions.
func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(
				t.senderBlockchainAddress,
				t.recipientBlockchainAddress,
				t.value))
	}
	return transactions
}

// ValidProof checks if the hash of a block with the given nonce, previousHash, and transactions meets the difficulty target.
func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	fmt.Printf("guessHashStr: %s\n", guessHashStr)
	return guessHashStr[:difficulty] == zeros
}

// ProofOfWork computes a valid nonce for a new block by iteratively searching for a hash that meets the mining difficulty.
func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

// Mining executes the mining process, rewards the miner, and adds a new block to the blockchain. Returns true on success.
func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

// CalculateTotalAmount computes the total balance for a specific blockchain address by summing all received and sent transactions.
func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}

			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

// Transaction represents a transfer of value between two blockchain addresses.
type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

// NewTransaction constructs a new Transaction with sender, recipient, and value.
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

// Print outputs the transaction details to stdout.
func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("_", 40))
	fmt.Printf(" sender_blockchain_address: %s\n", t.senderBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address: %s\n", t.recipientBlockchainAddress)
	fmt.Printf(" value: %.1f\n", t.value)
}

// MarshalJSON provides a custom JSON representation for Transaction fields.
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})

}

// Hash computes and returns the SHA-256 hash of the block's JSON representation.
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	fmt.Println(string(m))
	return sha256.Sum256(m)
}

// MarshalJSON provides a custom JSON representation for Block fields.
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

// init configures the logger prefix for the application.
func init() {
	log.SetPrefix("Blockchain: ")
}

// main is the entry point of the application demonstrating basic blockchain operations.
func main() {
	// Demo blockchain operations
	//myBlockchainAddress := "my_address"
	//bc := NewBlockchain(myBlockchainAddress)
	//bc.Print()
	//
	//bc.AddTransaction("Dika", "Bejo", 1.0)
	//previousHash := bc.LastBlock().Hash()
	//nonce := bc.ProofOfWork()
	//bc.CreateBlock(nonce, previousHash)
	//bc.Print()
	//
	//bc.AddTransaction("Batman", "Superman", 2.0)
	//bc.AddTransaction("Tukimin", "Tukiplus", 3.0)
	//previousHash = bc.LastBlock().Hash()
	//nonce = bc.ProofOfWork()
	//bc.CreateBlock(nonce, previousHash)
	//bc.Print()

	// Mining demo
	myBlockchainAddress := "my_address"
	bc := NewBlockchain(myBlockchainAddress)
	bc.Print()

	bc.AddTransaction("Dika", "Bejo", 1.0)
	bc.Mining()
	bc.Print()

	bc.AddTransaction("Batman", "Superman", 2.0)
	bc.AddTransaction("Tukimin", "Tukiplus", 3.0)
	bc.Mining()
	bc.Print()

	fmt.Printf("my_address %.1f\n", bc.CalculateTotalAmount("my_address"))
	fmt.Printf("Batman %.1f\n", bc.CalculateTotalAmount("Batman"))
	fmt.Printf("Superman %.1f\n", bc.CalculateTotalAmount("Superman"))
}
