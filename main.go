package main

import (
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
)

func UserHomeDir() string {
	usr, err := user.Current()
	HandleError(err)
	return usr.HomeDir
}

func BlocksDir() string {
	return filepath.Join(UserHomeDir(), ".bitcoin", "blocks")
}

func BlockFiles() []string {
	files, _ := filepath.Glob(filepath.Join(BlocksDir(), "blk*.dat"))
	sort.Strings(files)
	return files
}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func BlockStream() Streamer {
	blkFiles := BlockFiles()
	readers := make([](*os.File), len(blkFiles))
	// todo: make a coroutine that opens files on demand,
	// and closes them asap.
	// => not all at once
	for i, blkFile := range blkFiles {
		reader, err := os.Open(blkFile)
		HandleError(err)
		readers[i] = reader
	}
	return Stream(readers...)
}

func main() {

	db := DBOpen()
	defer db.Close()

	// it := db.NewIterator(readOptions)
	// defer it.Close()
	// for it.Seek([]byte("foo")); it.Valid(); it.Next() {
	// 	fmt.Println(string(it.Key()))
	// }

	var processedOffset BlkProcessedOffset
	if !DBGet(db, &processedOffset) {
		processedOffset = BlkProcessedOffset(0)
	}

	reader := BlockStream()

	// continue where we left off
	reader.Skip(int64(processedOffset))

	for {
		block := Block{}
		err := block.Decode(reader)
		if err == io.EOF {
			break
		} else {
			HandleError(err)
		}
		//fmt.Println(block)
		IndexBlock(db, &block)
		// update current stream offset, s.t. we can resume
		DBPut(db, BlkProcessedOffset(reader.Position()))
	}
}
