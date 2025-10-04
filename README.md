# Go Blockchain (Educational)

A minimal blockchain implementation in Go demonstrating blocks, transactions, hashing, proof-of-work, mining rewards, and simple balance calculation.

## Requirements
- Go 1.25+
- OS: Linux/Mac/Windows

## Run
- From project root:
  - go run blockchain.go

You will see block JSON (used for hashing), proof-of-work logs, mining logs, balances, and a readable chain printout.

## How the Blockchain Works

Core concepts:
- Transaction: a value transfer record between two addresses.
- Block: groups transactions, carries a nonce, timestamp, and the previous block’s hash.
- Blockchain: the ordered list of blocks plus a pool of pending transactions and a node (miner) address.
- Proof of Work: finds a nonce such that the block hash starts with a number of leading zeros equal to the mining difficulty.
- Mining: includes a special reward transaction and appends a new block upon successful proof-of-work.

Constants:
- MINING_DIFFICULTY = 3 (hash must begin with "000").
- MINING_SENDER = "THE BLOCKCHAIN" (issuer of mining reward).
- MINING_REWARD = 1.0 (amount rewarded to the miner per mined block).

High-level flow:
1) NewBlockchain initializes the chain with a genesis block and stores the miner’s blockchainAddress.
2) AddTransaction queues a transaction into the transaction pool.
3) ProofOfWork searches for a nonce that makes the next block’s hash satisfy the difficulty (using a stable copy of the pool).
4) Mining adds a reward transaction for the miner, runs ProofOfWork, then CreateBlock to append the new block and clear the pool.
5) Hash serializes a block to JSON and returns its SHA-256 hash (used for linking and PoW).
6) CalculateTotalAmount scans the chain to compute an address balance from sent/received transactions.
7) Print methods display blocks, transactions, and the entire chain.

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
    - blockchainAddress: string (address to receive mining rewards)
  - NewBlockchain(blockchainAddress string) -> *Blockchain
    - Initializes the chain with the miner’s address and creates the genesis block.
  - (bc *Blockchain) AddTransaction(sender, recipient, value)
    - Enqueues a transaction into transactionPool.
  - (bc *Blockchain) CopyTransactionPool() -> []*Transaction
    - Deep-copies the transaction pool for a stable proof-of-work input set.
  - (bc *Blockchain) ValidProof(nonce, previousHash, transactions, difficulty) -> bool
    - Builds a candidate block with these inputs and returns true if its hash has the required leading zeros.
  - (bc *Blockchain) ProofOfWork() -> int
    - Iteratively increments nonce until ValidProof is satisfied for MINING_DIFFICULTY.
  - (bc *Blockchain) Mining() -> bool
    - Adds a mining reward transaction (from MINING_SENDER to bc.blockchainAddress),
      runs ProofOfWork, creates the new block, logs success, and returns true.
  - (bc *Blockchain) CreateBlock(nonce, previousHash) -> *Block
    - Creates a block from the current transaction pool, appends to chain, clears the pool, and returns the new block.
  - (bc *Blockchain) LastBlock() -> *Block
    - Returns the most recent block in the chain.
  - (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) -> float32
    - Computes the balance by summing received minus sent amounts across all blocks.
  - (bc *Blockchain) Print()
    - Prints all blocks with separators for readability.

- init()
  - Sets a log prefix for application messages.

- main()
  - Demonstrates two flows (a basic flow is commented; mining flow is active):
    - Mining flow:
      1) Initialize blockchain with a miner address and print.
      2) Add a user transaction, call Mining (adds reward, runs PoW, creates block), print.
      3) Add more transactions, call Mining again, print.
      4) Print balances via CalculateTotalAmount for selected addresses.

## Notes
- Educational example: no network, no signatures, no validation rules.
- Hashing is performed on JSON-serialized block data.
- Proof-of-work target is defined by MINING_DIFFICULTY leading zeros in the hex hash.
- Mining rewards are simple coinbase-style transactions credited to the miner’s address.