package main

import (
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("./example/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 读写事务
	err = db.Update(func(txn *badger.Txn) error {
		txn.Set([]byte("answer"), []byte("42"))
		txn.Get([]byte("answer"))
		return nil
	})

	// 只读事务
	err = db.View(func(txn *badger.Txn) error {
		txn.Get([]byte("answer_v1"))
		return nil
	})

	// 遍历keys
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)

		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(val []byte) error {
				fmt.Printf("key=%s, value=%s\n", k, val)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	err = db.RunValueLogGC(0.7)
	_ = err
}
