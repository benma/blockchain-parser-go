package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
)

// constants
const blockMagic = 0xd9b4bef9

// types

type Hash256 [32]byte
type SatoshiValue uint64
type Script []byte

type VarInt uint64

type Block struct {
	header       Header
	txCount      VarInt
	transactions []Transaction
}

type Header struct {
	raw [80]byte
}

type Transaction struct {
	version  uint32
	inputs   []Input
	outputs  []Output
	lockTime uint32
}

type Input struct {
	prevTxHash      Hash256
	prevOutputIndex uint32
	script          Script
	sequenceNumber  uint32
}

type Output struct {
	value  SatoshiValue
	script Script
}

// Decode

func (h *Hash256) Decode(r io.Reader) error {
	return binary.Read(r, binary.LittleEndian, h)
}

func (script *Script) Decode(r io.Reader) error {
	var scriptLength VarInt
	if err := scriptLength.Decode(r); err != nil {
		return err
	}
	*script = make([]byte, scriptLength)
	if err := binary.Read(r, binary.LittleEndian, script); err != nil {
		return err
	}
	return nil
}

func (input *Input) Decode(r io.Reader) error {
	if err := input.prevTxHash.Decode(r); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &input.prevOutputIndex); err != nil {
		return err
	}
	if err := input.script.Decode(r); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &input.sequenceNumber); err != nil {
		return err
	}
	return nil
}

func (output *Output) Decode(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &output.value); err != nil {
		return err
	}

	if err := output.script.Decode(r); err != nil {
		return err
	}

	return nil
}

func (vi *VarInt) Decode(r io.Reader) error {
	var w uint8
	if err := binary.Read(r, binary.LittleEndian, &w); err != nil {
		return err
	}
	switch {
	case w < 0xfd:
		// fmt.Println("VarInt: Case 1")
		*vi = VarInt(w)
	case w == 0xfd:
		// fmt.Println("VarInt: Case 2")
		var res uint16
		binary.Read(r, binary.LittleEndian, &res)
		*vi = VarInt(res)
	case w == 0xfe:
		// fmt.Println("VarInt: Case 3")
		var res uint32
		binary.Read(r, binary.LittleEndian, &res)
		*vi = VarInt(res)
	case w == 0xff:
		// fmt.Println("VarInt: Case 4")
		var res uint64
		binary.Read(r, binary.LittleEndian, &res)
		*vi = VarInt(res)
	default:
		return errors.New("Failed decoding VarInt")
	}
	return nil
}

func (header *Header) Decode(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &header.raw); err != nil {
		return err
	}
	return nil
}

func (tx *Transaction) Decode(r io.Reader) error {
	// transaction version, should be 1, but is 2/garbage for a total of 5 transactions
	if err := binary.Read(r, binary.LittleEndian, &tx.version); err != nil {
		return err
	}

	// inputs
	var inCount VarInt
	if err := inCount.Decode(r); err != nil {
		return err
	}
	tx.inputs = make([]Input, inCount)
	for i := 0; i < int(inCount); i++ {
		if err := tx.inputs[i].Decode(r); err != nil {
			return err
		}
	}

	// outputs
	var outCount VarInt
	if err := outCount.Decode(r); err != nil {
		return err
	}
	tx.outputs = make([]Output, outCount)
	for i := 0; i < int(outCount); i++ {
		if err := tx.outputs[i].Decode(r); err != nil {
			return err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &tx.lockTime); err != nil {
		return err
	}

	return nil
}

func (block *Block) Decode(r io.Reader) error {
	var magic uint32
	if err := binary.Read(r, binary.LittleEndian, &magic); err != nil {
		return err
	}
	if magic != blockMagic {
		return errors.New("Block magic constant mismatch")
	}

	var size uint32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return err
	}

	rawHeader := make([]byte, 80)
	//var rawHeader [80]byte
	if err := binary.Read(r, binary.LittleEndian, &rawHeader); err != nil {
		return err
	}
	if err := block.header.Decode(bytes.NewReader(rawHeader)); err != nil {
		return err
	}
	if err := block.txCount.Decode(r); err != nil {
		return err
	}

	block.transactions = make([]Transaction, block.txCount)
	for i := 0; i < int(block.txCount); i++ {
		if err := block.transactions[i].Decode(r); err != nil {
			return err
		}
	}

	return nil
}

func (block *Block) BlockHash() Hash256 {
	result := hash256(hash256(block.header.raw[:]))
	var hash [32]byte
	copy(hash[:], result)
	return Hash256(hash)
}

func (block *Block) PrevHash() Hash256 {
	var hash [32]byte
	copy(hash[:], block.header.raw[4:4+32])
	return Hash256(hash)
}

func hash256(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}
