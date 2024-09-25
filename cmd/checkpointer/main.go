package main

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/samber/lo"
)

func generateNewDump() {
	var maxItems int64 = 133000000
	minItemID := math.MaxInt64 - maxItems

	items := make(map[string]int64, maxItems)
	bytes := make([]byte, 8)
	for val := range maxItems {
		binary.LittleEndian.PutUint64(bytes, uint64(minItemID+val))
		itemID := fmt.Sprintf("%036x", minItemID+val)
		items[itemID] = 0
		if val%10000000 == 0 {
			fmt.Println("val", val)
		}
	}
	fmt.Println("Generated, write to dump")

	output := lo.Must(os.Create("output.gob"))
	lo.Must0(gob.NewEncoder(output).Encode(items))
	lo.Must0(output.Close())
	fmt.Println("File written")
}

func main() {
	output := lo.Must(os.Open("output.gob"))
	defer output.Close()
	var items map[string]int64
	fmt.Println("loading file", time.Now())
	lo.Must0(gob.NewDecoder(output).Decode(&items))
	fmt.Println("file loaded", time.Now(), len(items))
}
