package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
)

func main() {
	db, err := bolt.Open("/home/yury/myData/geth/chaindata", 0600, &bolt.Options{ReadOnly: true})
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

		b.ForEach(func(k, v []byte) error {
			address := b2.Get(k)
			input := v
			length := len(input)
			groupVar := v[0]
			var group int

			if groupVar < 0x80 {
				group = 1
			} else if groupVar >= 0x80 && groupVar < 0xb8 {
				group = 2
			} else if groupVar >= 0xb8 && groupVar < 0xc0 {
				group = 3
			} else if groupVar >= 0xc0 && groupVar < 0xf8 {
				group = 4
			} else if groupVar >= 0xf8 {
				group = 5
			}
			fmt.Printf("key=%x, value=%x\n~length = %d, group = %d\n", address, v, length, group)
			result, _ := rlpDecode(input)
			fmt.Printf("result = %x\nlen = %d\n", result, len(result))
			if len(result) == 2 || len(result) == 4 {
				balance := new(big.Int)
				balance = balance.SetBytes(result[1])
				fmt.Printf("Balance = %d\n", balance)
			}else {
				fmt.Errorf("len(result) isn't 2 or 4")
			}
			return nil
		})
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func rlpDecode(input []byte) (result [][]byte, err error) {
	totalLength := len(input)
	startPosition := 0
	for startPosition < totalLength-1 {
		prefix := input[0]
		if prefix <= 0x7f {
			additionalSlice := []byte{input[startPosition]}
			dataLength := 1
			result = append(result, additionalSlice)
			startPosition += int(dataLength)
		} else if prefix >= 0x80 && prefix <= 0xb7 {
			dataLength := prefix - 0x80
			if dataLength > 0 {
				for i := startPosition + 1; i < int(dataLength)+1; i++ {
					additionalSlice := []byte{input[i]}
					result = append(result, additionalSlice)
				}
			}
			startPosition += int(dataLength)
		} else if prefix >= 0xb8 && prefix <= 0xbf {
			dataLength := input[startPosition+1] - 0xb8
			if dataLength > 0 {
				for i := startPosition + 2; i < int(dataLength)+1; i++ {
					additionalSlice := []byte{input[i]}
					result = append(result, additionalSlice)
				}
			}
			startPosition += int(dataLength) + 1
		} else if prefix >= 0xc0 && prefix <= 0xf7 {
			totalPayload := prefix - 0xc0
			var additionalSlice []byte
			if totalPayload > 0 {
				for i := startPosition + 1; i < int(totalPayload)+1; i++ {
					additionalSlice = append(additionalSlice, input[i])
				}
			}
			preResult, consumed, _ := rlpDecode2(additionalSlice)
			result = append(result, preResult)
			startPosition += consumed
		} else if prefix == 0xf8 && startPosition < totalLength-2 {
			totalPayload := prefix - 0xf7
			var additionalSlice []byte
			if totalPayload > 0 {
				for i := startPosition + 2; i < int(totalLength); i++ {
					additionalSlice = append(additionalSlice, input[i])
				}
			}
			preResult, consumed, _ := rlpDecode2(additionalSlice)
			result = append(result, preResult)
			startPosition += consumed
			if startPosition == totalLength-2 {
				break
			}
		} else {
			return nil, fmt.Errorf("Prefix value is bigger than 0xf8. It isn't supported")
		}
	}
	// fmt.Printf("result from rlpDecode = %x\n", result)
	// fmt.Printf("cap = %d\n", cap(result))
	// if cap(result)==2 || cap(result)==4{
	// 	fmt.Printf("val2 = %x\n",result[1])
	// 	balance := new(big.Int)
	// 	balance = balance.SetBytes(result[1])
	// 	fmt.Println(balance)
	// }
	return
}

func rlpDecode2(input []byte) (result []byte, consumed int, err error) {
	startPosition := 0
	prefix := input[startPosition]
	if prefix <= 0x7f {
		dataLength := 1
		result = append(result, prefix)
		consumed = dataLength
		startPosition += int(dataLength)
	} else if prefix >= 0x80 && prefix <= 0xb7 {
		dataLength := prefix - 0x80
		if dataLength > 0 {
			for i := startPosition + 1; i < int(dataLength)+1; i++ {
				result = append(result, input[i])
			}
		}
		consumed = int(dataLength) + 1
		startPosition += int(dataLength)
	} else if prefix >= 0xb8 && prefix <= 0xbf {
		dataLength := input[startPosition+1] - 0xb8
		if dataLength > 0 {
			for i := startPosition + 2; i < int(dataLength)+1; i++ {
				result = append(result, input[i])
			}
		}
		consumed = int(dataLength) + 2
		startPosition += int(dataLength) + 1
	} else if prefix >= 0xc0 {
		return nil, 0, fmt.Errorf("Groups 4 & 5 in function rlpDecode2 aren't supported")
	}
	return
}
