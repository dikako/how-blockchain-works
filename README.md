# Go Blockchain (Educational)

A minimal blockchain implementation in Go demonstrating blocks, transactions, hashing, and proof-of-work.

## Requirements
- Go 1.25+
- OS: Linux/Mac/Windows

## Run
- From project root:
    - go run blockchain.go

You will see block JSON (used for hashing), proof-of-work logs, and a readable chain printout.

## How the Blockchain Works

Core concepts:
- Transaction: a value transfer record between two addresses.
- Block: groups transactions, carries a nonce, timestamp, and the previous block’s hash.
- Blockchain: the ordered list of blocks plus a pool of pending transactions.
- Proof of Work: finds a nonce such that the block hash starts with a number of leading zeros equal to the mining difficulty.

Mining difficulty:
- MINING_DIFFICULTY = 3 (hash must begin with "000").

High-level flow:
1) NewBlockchain creates a genesis block that references the hash of an empty block holder.
2) AddTransaction places transactions into the transaction pool.
3) ProofOfWork searches for a nonce that makes the next block’s hash satisfy the difficulty.
4) CreateBlock builds the block with that nonce and links it via previousHash to the tip of the chain; the transaction pool is then cleared.
5) Hash serializes a block to JSON and returns its SHA-256 hash.
6) Print methods display the chain and transactions for inspection.

## Flow Code Run (per struct and function)

- Struct: Transaction
    - Fields:
        - senderBlockchainAddress: string
        - recipientBlockchainAddress: string
        - value: float32
    - NewTransaction(sender, recipient, value) -> *Transaction
        - Creates a transaction used in the pool and inside blocks.
    - (t *Transaction) Print()
        - Prints sender, recipient, and value for debugging.
    - (t *Transaction) MarshalJSON() -> []byte, error
        - Custom JSON for stable, readable output and consistent hashing.

- Struct: Block
    - Fields:
        - timestamp: int64 (nanoseconds)
        - nonce: int
        - previousHash: [32]byte
        - transactions: []*Transaction
    - NewBlock(nonce, previousHash, transactions) -> *Block
        - Creates a block with the current time and provided transactions.
    - (b *Block) Print()
        - Prints timestamp, nonce, previous hash, and all transactions.
    - (b *Block) Hash() -> [32]byte
        - Marshals the block to JSON (also printed) and returns the SHA-256 hash.
        - Used to link blocks (as the next block’s previousHash) and to verify proof-of-work.
    - (b *Block) MarshalJSON() -> []byte, error
        - Custom JSON to ensure predictable hashing layout.

- Struct: Blockchain
    - Fields:
        - transactionPool: []*Transaction (pending transactions to be mined)
        - chain: []*Block (ordered list of blocks)
    - NewBlockchain() -> *Blockchain
        - Initializes the chain by creating the genesis block using the hash of an empty block value holder.
    - (bc *Blockchain) AddTransaction(sender, recipient, value)
        - Enqueues a transaction into transactionPool.
    - (bc *Blockchain) CopyTransactionPool() -> []*Transaction
        - Deep-copies the transaction pool for a stable proof-of-work input set.
    - (bc *Blockchain) ValidProof(nonce, previousHash, transactions, difficulty) -> bool
        - Builds a candidate block with these inputs and returns true if its hash has the required leading zeros.
    - (bc *Blockchain) ProofOfWork() -> int
        - Iteratively increments nonce until ValidProof is satisfied for MINING_DIFFICULTY.
    - (bc *Blockchain) CreateBlock(nonce, previousHash) -> *Block
        - Creates a block from the current transaction pool, appends to chain, clears the pool, and returns the new block.
    - (bc *Blockchain) LastBlock() -> *Block
        - Returns the most recent block in the chain.
    - (bc *Blockchain) Print()
        - Prints all blocks with separators for readability.

- init()
    - Sets a log prefix for application messages.

- main()
    - Demonstrates the full flow:
        1) Initialize blockchain and print.
        2) Add a transaction, compute previousHash (tip’s hash), run ProofOfWork for a valid nonce, create the block, print.
        3) Add more transactions, repeat proof-of-work and block creation, print again.

## Notes
- Educational example: no consensus, validation, or networking.
- Hashing is performed on JSON-serialized block data.
- Proof-of-work target is defined by MINING_DIFFICULTY leading zeros in the hex hash.