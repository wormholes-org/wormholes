package types

func (b *Block) SetExtra(data []byte) error {
	b.header.Extra = data
	return nil
}
