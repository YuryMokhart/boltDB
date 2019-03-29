package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/wcharczuk/go-chart"
	"log"
	"math/big"
	// "os"
	"io/ioutil"
	"sort"
	"bytes"
)

func main() {
	db, err := bolt.Open("/Users/yurymokhart/Library/Ethereum/geth/chaindata", 0600, &bolt.Options{ReadOnly: true})
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
	check(err)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("AT"))
		b2 := tx.Bucket([]byte("secure-key-"))

		var balanceSlice []*big.Int
		var additionalAddress [][20]byte
		var address2 [20]byte

		b.ForEach(func(k, v []byte) error {
			address := b2.Get(k)
			input := v
			balance := new(big.Int)
			result, _ := rlpDecode(input)
			if len(result) == 2 || len(result) == 4 {
				balance = balance.SetBytes(result[1])
			} else {
				fmt.Errorf("len(result) isn't 2 or 4")
			}

			copy(address2[:], address)

			balanceSlice = append(balanceSlice, balance)
			additionalAddress = append(additionalAddress, address2)

			return nil
		})

		var accounts Accounts = Accounts{len(balanceSlice), additionalAddress, balanceSlice}

		sort.Sort(accounts)

		// file, err := os.Create("accounts.txt")
		// check(err)
		// defer file.Close()

		// for i := 0; i < 100; i++ {
		// 	//fmt.Printf("address = %x  balance = %d\n", accounts.keys[i], accounts.values[i])
		// 	//fmt.Fprintf(file, "address = %x  balance = %d\n", accounts.keys[i], accounts.values[i])
		// }

		count := 0
		var valueSlice []*big.Int
		var countSlice []int
		prevVal := accounts.values[0]

		for i , val := range accounts.values{
			if (prevVal).Cmp(val) != 0{
				valueSlice = append(valueSlice, prevVal)
				countSlice = append(countSlice, count)
				prevVal = val
			}
			count++
			if i == len(accounts.values)-1{ //because otherwise it doesn't include the last element not equal to the previous
				valueSlice = append(valueSlice, val)
				countSlice = append(countSlice, count)
			}
		}
		//fmt.Println(valueSlice)
		//fmt.Println(countSlice)
		chartAccountsBalances(countSlice, valueSlice)
		return nil
	})
	check(err)
}

func chartAccountsBalances(accI []int, balanceI []*big.Int) {
	var acc []float64
	var bal []float64
	for _, temporary := range accI{
		acc = append(acc, float64(temporary)/1e6)
	}
	for _, temporary := range balanceI{
		f := new(big.Float).SetInt(temporary)
		f64, _ := f.Float64()
		if f64 > 1e24{
			continue
		}
		bal = append(bal, f64)
	}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "wei",
			NameStyle: chart.StyleShow(),
			Style: chart.StyleShow(),
			Range: &chart.LogRange{
				Min: bal[0],
				Max: bal[len(bal)-1],
			},
			Ticks: []chart.Tick{
				{Value: 0.0, Label: "0"},
				{Value: 1e3, Label: "1e3"},
				{Value: 1e6, Label: "1e6"},
				{Value: 1e9, Label: "1e9"},
				{Value: 1e12, Label: "1e12"},
				{Value: 1e15, Label: "1e15"},
				{Value: 1e18, Label: "1e18"},
				{Value: 1e21, Label: "1e21"},
				{Value: 1e24, Label: "1e24"},
			},
		},
		YAxis: chart.YAxis{
			Name: "number of accounts, million",
			NameStyle: chart.StyleShow(),
			Style: chart.StyleShow(),
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.DefaultStrokeColor,
				StrokeWidth: 1.0,
			},
			GridMinorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.DefaultStrokeColor,
				StrokeWidth: 1.0,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.ColorRed,
					FillColor:   chart.ColorRed.WithAlpha(50),
				},
				XValues: bal,
				YValues: acc,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	check(err)
	err = ioutil.WriteFile("output.PNG", buffer.Bytes(), 0644)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

type Accounts struct {
	length int
	keys   [][20]byte
	values []*big.Int
}

func (ac Accounts) Len() int {
	return ac.length
}

func (ac Accounts) Less(i, j int) bool {
	return (ac.values[i]).Cmp(ac.values[j]) == -1
}

func (ac Accounts) Swap(i, j int) {
	ac.values[i], ac.values[j] = ac.values[j], ac.values[i]
	ac.keys[i], ac.keys[j] = ac.keys[j], ac.keys[i]
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
