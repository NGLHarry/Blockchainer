package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	//Data,trade data
	Data []byte
	//Nonce random number
	Nonce int
}

//Serialize a Block object to []byte
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encode := gob.NewEncoder(&result)

	err := encode.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

//Deserialize []byte into Blocks
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
