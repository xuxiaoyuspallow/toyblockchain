package main


type BlockChain struct {
	blocks []*Block
}

func (bc *BlockChain) AddBlock(data string)  {
	preBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data,preBlock.Hash)
	bc.blocks = append(bc.blocks,newBlock)
}

func NewGenesisBlock() *Block{
	return NewBlock("Genesis Block", []byte{})   //you can write anything on genesis block
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}
