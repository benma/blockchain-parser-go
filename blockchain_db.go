package main

import (
	"bytes"
	"encoding/binary"
)

func deserialize(data []byte, value interface{}) {
	binary.Read(bytes.NewReader(data), binary.LittleEndian, value)
}

func serialize(data interface{}) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, data)
	HandleError(err)
	return buf.Bytes()
}

type BlkProcessedOffset int64

func (kv BlkProcessedOffset) KeyEncode() []byte {
	return []byte("processed-offset")
}

func (kv BlkProcessedOffset) ValueEncode() []byte {
	return serialize(kv)
}

func (kv *BlkProcessedOffset) ValueDecode(data []byte) {
	deserialize(data, kv)
}

type BlkKeyBlkToHeight struct {
	blockHash Hash256
	height    int32
}

func (kv BlkKeyBlkToHeight) KeyEncode() []byte {
	buf := new(bytes.Buffer)
	buf.Write([]byte("idx-"))
	buf.Write(kv.blockHash[:])
	return buf.Bytes()
}

func (kv BlkKeyBlkToHeight) ValueEncode() []byte {
	return serialize(kv.height)
}

func (kv *BlkKeyBlkToHeight) ValueDecode(data []byte) {
	deserialize(data, &kv.height)
}

type BlkKeyTip Hash256

func (kv BlkKeyTip) KeyEncode() []byte {
	return []byte("tip")
}

func (kv BlkKeyTip) ValueEncode() []byte {
	return serialize(kv)
}

func (kv *BlkKeyTip) ValueDecode(data []byte) {
	deserialize(data, kv)
}
