package main

import (
	"./code"
	"fmt"
	"io/ioutil"
	"os"
)

const BLOCK_SIZE = 2 ^ 14 // 16 KB

func main() {
	readAtorrentFile()
}

// kelines Experiment zum lesen einer bestehenden Torrentfile
func readAtorrentFile() {
	filename := "2018-03-13-raspbian-stretch-lite.zip.torrent"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	bencodeText, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := bencode.Decode(bencodeText)
	datamap := data.(map[string]interface{})

	//fmt.Println(datamap["encoding"])
	fmt.Println(datamap["announce"])
	//fmt.Println(datamap["created by"])
	//fmt.Println(datamap["creation date"])

	info := datamap["info"].(map[string]interface{})

	// sigle file
	fmt.Println(info["name"])
	fmt.Println(info["piece length"])
	fmt.Println(info["private"])
	fmt.Println(info["length"])
	fmt.Printf("%x\n", splitPieces(info["pieces"].(string)))
}

func splitPieces(s string) []string {
	var liste []string
	for i := 0; i < len(s); i += 20 {
		liste = append(liste, s[i:i+20])
	}
	return liste
}
