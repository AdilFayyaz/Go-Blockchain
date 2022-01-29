package assignment02IBC

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

const miningReward = 100
const rootUser = "Satoshi"


type BlockData struct {
	Title    string
	Sender   string
	Receiver string
	Amount   int
}
type Block struct {
	Data        []BlockData
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}
type UserBalance struct {
	Name	string
	Balance int
	AmtToSend int
}


func CalculateBalance(userName string, chainHead *Block) int {
	balance := 0
	for chainHead != nil{
		for _,info := range chainHead.Data {
			// Subtract amount if sender
			if info.Sender == userName {
				balance -= info.Amount
			}
			// Add amount if receiver
			if info.Receiver == userName {
				balance += info.Amount
			}
		}
		chainHead = chainHead.PrevPointer
	}

	return balance
}
func CalculateHash(inputBlock *Block) string {
	val := sha256.New()
	val.Write([]byte(fmt.Sprintf("%v", inputBlock.Data)))
	return fmt.Sprintf("%x", val.Sum(nil))
}
func VerifyTransaction(transaction *BlockData, chainHead *Block) bool {
	tempBalance := CalculateBalance(transaction.Sender, chainHead)
	if transaction.Amount <= tempBalance{
		return true
	}
	return false
}

func getUserBalance(block BlockData, chainHead *Block, name string) (string,int){
	return name, CalculateBalance(block.Sender, chainHead)
}

func InsertBlock(blockData []BlockData, chainHead *Block) *Block {

	// Coin based transaction
	tempTitle2 := "Coinbase"
	tempSender2 := "System"
	tempReceiver2 := rootUser
	tempAmount2 := miningReward
	// Create a new block - current block
	var newBlock = new(Block)

	newBlock.Data = blockData

	indivFailing := false
	// Verify the transactions
	for _,val := range blockData{
		isValid := VerifyTransaction(&val, chainHead)
		if isValid == false {
			indivFailing = true
		}
	}

	// Individual Transactions are checked
	// Check collective transactions
	// Bonus
	var usrBal []UserBalance

	for _, val := range blockData{
		// Get the balance for Sender
		var name, bal = getUserBalance(val, chainHead, val.Sender)

		// Get the balance for the Receiver
		var name2, bal2 = getUserBalance(val, chainHead, val.Receiver)

		nameExists2 := false
		nameExists := false
		// Check if Name exists in the array
		for _, b := range usrBal{

			if b.Name == name{
				nameExists = true
			}
			if b.Name == name2{
				nameExists2 = true
			}
		}
		if nameExists == false{
			usrBal = append(usrBal, UserBalance{name, bal, 0})
		}
		if nameExists2 == false{
			usrBal = append(usrBal, UserBalance{name2,bal2, 0 })
		}

		// Now we have a list of users and their balances
		// Update the balances
		for i, x:= range usrBal{
			if x.Name == val.Sender{
				usrBal[i].AmtToSend += val.Amount
				// subtract from balance
				usrBal[i].Balance -= val.Amount
			}
			if x.Name == val.Receiver{
				usrBal[i].Balance += val.Amount
			}
		}
	}
	for _, val := range blockData {
		for i, j := range usrBal {
			if j.Balance < 0 {
				if indivFailing == true {
					fmt.Printf("ERROR: " + val.Sender + " has " + strconv.Itoa(CalculateBalance(j.Name,chainHead)) + " coins -> " + strconv.Itoa(val.Amount) + " required\n")
					return chainHead
				}
				fmt.Printf("ERROR: " + val.Sender + " has " + strconv.Itoa(CalculateBalance(j.Name,chainHead)) + " coins -> " + strconv.Itoa(usrBal[i].AmtToSend) + " required\n")
				return chainHead
			}
		}
	}
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
	// Append the coinbase transaction to the block
	newBlock.Data = append(newBlock.Data,BlockData{tempTitle2, tempSender2, tempReceiver2, tempAmount2})
	return newBlock

}
func ListBlocks(chainHead *Block) {
	fmt.Printf("******************** Printing Blocks: ********************\n")
	for chainHead != nil {
		for _,val := range chainHead.Data{
			fmt.Printf("Title: " + val.Title + " Sender:" + val.Sender + " Receiver:" + val.Receiver + " Amount:" + strconv.Itoa(val.Amount) + "\n")
		}
		chainHead = chainHead.PrevPointer
		if chainHead != nil {
			fmt.Printf("\t\t\t\t\t\t|\n")
			fmt.Printf("\t\t\t\t\t\tV\n")
		}else {
			fmt.Printf("\n\n")
		}

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


func PremineChain(chainHead *Block, numBlocks int) *Block {
	i := 0
	if chainHead == nil {

		for i < numBlocks {

			// Pre mined transaction
			tempTitle := "Premined Sender"
			tempSender := "nil"
			tempReceiver := "nil"
			tempAmount := 0

			// Coin based transaction
			tempTitle2 := "Coinbase"
			tempSender2 := "System"
			tempReceiver2 := rootUser
			tempAmount2 := miningReward

			preminedTransactions := []BlockData{{tempTitle, tempSender, tempReceiver, tempAmount}, {tempTitle2,
				tempSender2, tempReceiver2, tempAmount2}}

			var newBlock = new(Block)
			newBlock.Data = preminedTransactions
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
			chainHead = newBlock
			i++
		}
	}
	return chainHead
}