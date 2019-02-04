package main

import (
	"log"
	"fmt"
	//"bytes"
	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("/home/yury/myData/geth/chaindata", 0600, &bolt.Options{ReadOnly:true})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
        return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
            fmt.Println(string(name))
            return nil
        })
    })
    if err != nil {
        fmt.Println(err)
        return
    }

    err = db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("AT"))
        b2 := tx.Bucket([]byte("secure-key-"))

        b.ForEach(func(k, v []byte)error{
        	address := b2.Get(k)
        	input := v
        	length := len(input)
        	groupVar := v[0]
        	group := 0
        	if (groupVar < 0x80) {
        		group = 1
        	} else if (groupVar >= 0x80 && groupVar < 0xb8){
        		group = 2
        	} else if (groupVar >= 0xb8 && groupVar < 0xc0){
        		group = 3
        	} else if (groupVar >= 0xc0 && groupVar < 0xf8){
        		group = 4
        	} else if (groupVar >= 0xf8){
        		group = 5
        	}
        	fmt.Printf("key=%x, value=%x\n ||length = %d, group = %d, groupVar = %d \n", address, v, length, group, groupVar)
        	return nil
        	})
        return nil
    })
    if err != nil {
        fmt.Println(err)
        return
    }   
}

// func rlpDecode(input []byte)[][]byte{
// 	prefix := input[0]
// 	length := len(input)
// 	if (len(input) ==0){
// 		return 
// 	}
// 	output := ''
// 	if (prefix >= 0x00 && prefix <= 0x7f){
// 		offset := 0
// 		dataLen := 1
// 		output := instantiate_str(string(input, offset, dataLen))
// 		return output
// 	} else if (prefix >= 0x80 && prefix <= 0xb7 && length > prefix - 0x80){
// 		strlen := prefix - 0x80
// 		offset := 1
// 		dataLen := strlen
// 		output := instantiate_str(substr(input, offset, dataLen))
// 		return output
// 	} else if(prefix >= 0xb8 && prefix <= 0xbf && length > prefix - 0xb7 && length > prefix - 0xb7 + toInteger(substr(input,1,prefix-0xb7))){
// 		lenOfStrlen := prefix - 0xb7
// 		strlen := toInteger(substr(input,1,lenOfStrlen))
// 		offset := 1 + lenOfStrlen
// 		dataLen := strlen
// 		output := instantiate_str(substr(input, offset, dataLen))
// 		return output
// 	} else if (prefix >= 0xc0 && prefix <= 0xf7 && length > prefix - 0xc0){
// 		listLen := prefix - 0xc0
// 		offset := 1
// 		dataLen := listLen
// 		output := instantiate_list(substr(input, offset, dataLen))
//     output + rlp_decode(substr(input, offset + dataLen))
// 		return output
// 	} else if (prefix >= 0xf8 && prefix <= 0xff && length > prefix - 0xf7 && length > prefix - 0xf7 + toInteger(substr(input, 1, prefix - 0xf7))){
//         lenOfListLen := prefix - 0xf7
//         listLen := toInteger(substr(input, 1, lenOfListLen))
//         offset := 1 + lenOfListLen
// 		dataLen := listLen
//         output := instantiate_list(substr(input, offset, dataLen))
//     output + rlp_decode(substr(input, offset + dataLen))
//         return output
//     } else {
//     	fmt.Errorf("input doesn't conform RLP encoding form")
//     }

// }

// func toInteger(b [][]byte) int{
// 	length := len(b)
// 	if (length == 0){
// 		fmt.Errorf("Input is nil")
// 	} else if (length == 1){
// 		return b[0]
// 	} else {
// 		return substr(b,-1) + toInteger(substr(b,0,-1)) * 256
// 	}
// }