package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"

	badger "github.com/dgraph-io/badger/v3"
)

// record struct
type RECORD struct {
	Name    string
	Label   string
	Owner   string
	Expires *big.Int
}

// 用gob进行数据编码
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// -------------------
// Decode
// 用gob进行数据解码
func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

// test db and encoder
func Test() {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions("./badger")
	opts.Logger = nil // close log
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// save k-v
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("test"), []byte("42"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	// get data in a view
	err = db.View(func(txn *badger.Txn) error {
		// get item
		item, err := txn.Get([]byte("answer"))
		if err != nil {
			log.Fatal(err)
		}

		// copy data from item
		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("The answer is: %s\n", valCopy)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	r := RECORD{
		Label:   "label",
		Owner:   "wy",
		Expires: big.NewInt(123123123),
	}
	// encode
	fmt.Println("encoding")
	b, err := Encode(r)
	if err != nil {
		fmt.Println(err)
	}

	// decode
	res := RECORD{}
	fmt.Println("decoding")
	Decode(b, &res)
	if err != nil {
		fmt.Println("Error decoding GOB data:", err)
		log.Fatal(err)
	}

	fmt.Printf("%s %s %s\n", res.Label, res.Owner, res.Expires.String())
}

// read and decode
func Read(name string) {
	opts := badger.DefaultOptions("./badger")
	opts.Logger = nil // close log
	// open db
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// read data
	err = db.View(func(txn *badger.Txn) error {
		// get item
		k, err := Encode(name)
		if err != nil {
			log.Fatal(err)
		}

		// read
		fmt.Println("reading")
		item, err := txn.Get(k)
		if err != nil {
			log.Fatal(err)
		}

		// copy data from item
		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			log.Fatal(err)
		}

		// decode
		res := RECORD{}
		fmt.Println("decoding")
		Decode(valCopy, &res)

		fmt.Println("Record info:")
		fmt.Println("Name", res.Name)
		fmt.Println("Label", res.Label)
		fmt.Println("Owner", res.Owner)
		fmt.Println("Expires", res.Expires)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// uint64 to []byte
func Uint64toBytes(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

// [8]byte is needed
func BytestoUint64(b []byte) uint64 {
	if len(b) > 8 {
		log.Fatal("bytes too big")
	}

	// 8 bytes needed
	bytes8 := make([]byte, 8)
	copy(bytes8, b)

	i := binary.LittleEndian.Uint64(bytes8)
	return i
}
