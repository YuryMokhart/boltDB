package main

import (
	"log"
	"fmt"
	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key := []byte("firstKey")
	value := []byte("creativeFirstValue")

	db.Update(func(tx *bolt.Tx) error{
		b, err := tx.CreateBucket([]byte("FirstBucket"))
		if err != nil{
			return fmt.Errorf("Create Bucket %s\n", err)
		}

		err = b.Put(key, value)
		return nil
		})

	db.Update(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte("FirstBucket"))
		err := b.Put([]byte("someData"), []byte("111"))
		return err
		})

	db.View(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte("FirstBucket"))
		v := b.Get([]byte("someData"))
		value1 := b.Get(key) //just trying to use another way
		fmt.Printf("Value of someData is %s\n", v)
		fmt.Println(string(value1)) //just another way
		return nil
	})

	db.Update(func(tx *bolt.Tx) error{
		b, err := tx.CreateBucket([]byte("DelBucket"))
		if err != nil{
			return fmt.Errorf("Create Bucket %s\n", err)
		}
		err = b.Put([]byte("DelKey"), []byte("DelVal"))
		return nil
		})

	db.View(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte("DelBucket"))
		v := b.Get([]byte("DelKey"))
		fmt.Println(string(v))
		return nil
	})

	db.View(func (tx *bolt.Tx) error{
		b := tx.Bucket([]byte("FirstBucket"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next(){
			fmt.Printf("key = %s, value = %s\n", k, v)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error{
		err := tx.DeleteBucket([]byte("DelBucket"))
		if err != nil{
			return fmt.Errorf("Bucket wasn't deleted.")
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error{
		err := tx.DeleteBucket([]byte("FirstBucket"))
		if err != nil{
			return fmt.Errorf("Bucket wasn't deleted.")
		}
		return nil
	})
}