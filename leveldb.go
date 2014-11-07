package main

import "github.com/jmhodges/levigo"

var readOptions *levigo.ReadOptions = ReadOptions()
var writeOptions *levigo.WriteOptions = WriteOptions()

func ReadOptions() *levigo.ReadOptions {
	readOptions := levigo.NewReadOptions()
	// readOptions.SetVerifyChecksums(true)
	// readOptions.SetFillCache(false)
	return readOptions
}

func WriteOptions() *levigo.WriteOptions {
	writeOptions := levigo.NewWriteOptions()
	//writeOptions.SetSync(true)
	return writeOptions
}

func DBOpen() *levigo.DB {
	dbOptions := levigo.NewOptions()
	dbOptions.SetCreateIfMissing(true)
	db, err := levigo.Open("db", dbOptions)
	HandleError(err)
	return db
}

type ValueEncoder interface {
	ValueEncode() []byte
}
type KeyEncoder interface {
	KeyEncode() []byte
}
type ValueDecoder interface {
	ValueDecode([]byte)
}

func DBGet(db *levigo.DB, kv interface {
	KeyEncoder
	ValueDecoder
}) bool {
	value, err := db.Get(readOptions, kv.KeyEncode())
	HandleError(err)
	if value == nil {
		return false
	}
	kv.ValueDecode(value)
	return true
}

func DBPut(db *levigo.DB, kv interface {
	KeyEncoder
	ValueEncoder
}) {
	db.Put(writeOptions, kv.KeyEncode(), kv.ValueEncode())
}

func DBWPut(wb *levigo.WriteBatch, kv interface {
	KeyEncoder
	ValueEncoder
}) {
	wb.Put(kv.KeyEncode(), kv.ValueEncode())
}

func DBWrite(db *levigo.DB, wb *levigo.WriteBatch) {
	db.Write(writeOptions, wb)
}
