package messages

//go:generate msgp

type MessageBody []byte

func (b *MessageBody) ExtensionType() int8 { return 98 }

func (b *MessageBody) Len() int { return len(*b) }

func (b *MessageBody) MarshalBinaryTo(b1 []byte) error {
	b1 = *b
	return nil
}

func (b *MessageBody) UnmarshalBinary(b1 []byte) error {
	*b = b1
	return nil
}
