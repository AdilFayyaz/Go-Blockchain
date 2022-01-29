package assignment01IBC

import (
	"crypto/sha256"
	"fmt"
)

type BlockData struct {
	Transactions []string
}
type Block struct {
	Data        BlockData
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}
func CalculateHash(inputBlock *Block) string {
	val := sha256.New()
	val.Write([]byte(fmt.Sprintf("%v", inputBlock.Data)))
	return fmt.Sprintf("%x", val.Sum(nil))
}
func InsertBlock(dataToInsert BlockData, chainHead *Block) *Block {
	// Create a new block - current block
	var newBlock = new(Block)
	newBlock.Data = dataToInsert

	// Genesis block
	if chainHead == nil {
		// Genesis block prev pointer points to nil
		newBlock.PrevPointer = nil
		newBlock.PrevHash = ""
		blockHash := CalculateHash(newBlock)

		// Write the hash value to the new Block
		newBlock.CurrentHash = blockHash
	} else{
		// Point the previous pointer to the chainHead
		// Update Prev Hash
		newBlock.PrevPointer = chainHead
		newBlock.PrevHash = chainHead.CurrentHash
		blockHash := CalculateHash(newBlock)

		// Write the hash value to the new Block
		newBlock.CurrentHash = blockHash
	}
	return newBlock
}

func ChangeBlock(oldTrans string, newTrans string, chainHead *Block) {
	currPtr := chainHead
	for currPtr != nil{
		// Iterate over the transactions
		for index, val := range currPtr.Data.Transactions{
			// If transaction found, update the data
			if oldTrans == val{
				currPtr.Data.Transactions[index] = newTrans
				return
			}
		}
		currPtr = currPtr.PrevPointer
	}
}
func ListBlocks(chainHead *Block) {
	fmt.Printf("******************** Printing Blocks: ********************\n")
	for chainHead != nil {
		fmt.Printf("Data: %s\n", chainHead.Data.Transactions)
		fmt.Printf("Current Hash: %s\n", chainHead.CurrentHash)
		fmt.Printf("Pervious Hash %s\n", chainHead.PrevHash)

		chainHead = chainHead.PrevPointer
		fmt.Printf( "\t\t\t\t\t\t|\n")
		fmt.Printf("\t\t\t\t\t\tV\n")
	}
}
func VerifyChain(chainHead *Block) {
	fmt.Printf( "******************** Verifying Blockchain... ********************\n")
	// While chain Head
	for chainHead != nil {
		// Fetch the Prev hash and point the head to
		// prev pointer
		var prevHashValue string
		prevHashValue = chainHead.PrevHash
		// Update the pointer to point to the prev block
		chainHead = chainHead.PrevPointer
		// Check if chainHead is null and if prevHashVal is not equal to
		// the current hash of the head
		if (chainHead != nil) && (prevHashValue != CalculateHash(chainHead)){
			fmt.Printf("Blockchain Compromised! :( ")
			return
		}
	}
	fmt.Printf("Blockchain Verified! :)")
}

