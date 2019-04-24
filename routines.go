package main

import (
	"log"
	"fmt"
	"github.com/boltdb/bolt"
	"encoding/binary"
	"time"
)

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	check(err)
	defer db.Close()
	
	go input(db)
	go output(db)

	var n string
	fmt.Scanln(&n)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
		return
	}
}

func input(db *bolt.DB) {
	for val := 1; val < 2000; val++ {
		fmt.Printf("Writing at %s\n", time.Now())
		var valBytes [4]byte
		binary.BigEndian.PutUint32(valBytes[:], uint32(val))
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Sequence"))
			if b != nil{
				err := tx.DeleteBucket([]byte("Sequence"))
				check(err)
			}
			b, err := tx.CreateBucket([]byte("Sequence"))
			check(err)

			for i := 1; i < 1001; i++ {
				var sl [4]byte
				binary.BigEndian.PutUint32(sl[:], uint32(i))
				err = b.Put(sl[:], valBytes[:])
				time.Sleep(time.Millisecond)
			}
			return nil
		})
		fmt.Printf("Finished writing %x at %s\n", valBytes, time.Now())
	}
}

func output(db *bolt.DB) {
	for {
		fmt.Printf("Reading at %s\n", time.Now())
		var valBytes [4]byte
		db.View(func(tx *bolt.Tx) error {
			c := tx.Bucket([]byte("Sequence")).Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				copy(valBytes[:], v)
			}
			return nil
		})
		fmt.Printf("Last value: %x\n", valBytes)
		time.Sleep(time.Second)
	}
}