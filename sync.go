package main

import (
	"fmt"
	"log"

	"github.com/jmhodges/levigo"
)

func Reverse(b []byte) {
	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func prettyHash(hash Hash256) string {
	b := hash // copy
	Reverse(b[:])
	return fmt.Sprintf("%x", b)
}

func GetTip(db *levigo.DB) (int32, *Hash256) {
	var tip BlkKeyTip
	if !DBGet(db, &tip) {
		return -1, nil
	}
	var kvHeight BlkKeyBlkToHeight
	tipHash := Hash256(tip)
	kvHeight.blockHash = tipHash
	DBGet(db, &kvHeight)

	return kvHeight.height, &tipHash
}

func IndexBlock(db *levigo.DB, block *Block) {
	tipHeight, tip := GetTip(db)

	if tip == nil || *tip == block.PrevHash() {
		// index this block

		curHash := block.BlockHash()
		// if tipHeight%10000 == 0 {
		// 	fmt.Println(tipHeight+1, curHash)
		// }

		wb := levigo.NewWriteBatch()
		DBWPut(wb, BlkKeyBlkToHeight{curHash, tipHeight + 1})
		DBWPut(wb, BlkKeyTip(curHash))
		DBWrite(db, wb)
	} else {
		// handle out of order
		log.Fatal("needs orphanizing!")
	}
}
