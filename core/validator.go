package core

import "fmt"

type Validator interface {
	ValidateBlock(block *Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

func (v *BlockValidator) ValidateBlock(b *Block) error {
	if v.bc.HasBlock(b.Height) {
		return fmt.Errorf("block at height %d with hash %s already exists", b.Height, b.Hash(BlockHasher{}))
	}

	if b.Height != v.bc.Height()+1 {
		return fmt.Errorf("invalid block height %d, expected %d", b.Height, v.bc.Height()+1)
	}

	prevHeader, err := v.bc.GetHeader(b.Height - 1)
	if err != nil {
		return err
	}

	hash := BlockHasher{}.Hash(prevHeader)
	if b.PrevBlockHash != hash {
		return fmt.Errorf("invalid prev block hash %s, expected %s", b.PrevBlockHash, hash)
	}

	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
