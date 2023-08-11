package db

import (
	"fmt"
	"os"
	"strings"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

type BadgerDB struct {
	Path string
	DB   *badger.DB
}

func New(path string) *BadgerDB {
	db := &BadgerDB{
		Path: path,
		DB:   nil,
	}

	return db
}

func (b *BadgerDB) Open() error {
	if _, err := os.Stat(b.Path); os.IsNotExist(err) {
		os.MkdirAll(b.Path, 0755)
	}
	opts := badger.DefaultOptions(b.Path)
	opts.Logger = nil
	opts.Dir = b.Path
	opts.ValueDir = b.Path
	opts.SyncWrites = false
	opts.ValueThreshold = 256
	opts.CompactL0OnClose = true
	db, err := badger.Open(opts)
	if err != nil {
		return fmt.Errorf("badger open failed. path %s, error %s", b.Path, err)
	}

	b.DB = db

	return nil
}

func (b *BadgerDB) Close() error {
	err := b.DB.Close()
	return err
}

func (b *BadgerDB) Set(key []byte, value []byte) error {
	wb := b.DB.NewWriteBatch()
	defer wb.Cancel()

	err := wb.SetEntry(badger.NewEntry(key, value).WithMeta(0))
	if err != nil {
		return fmt.Errorf("failed to write data to cache. key %s, value %s, error %s", string(key), string(value), err)
	}

	err = wb.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush data to cache. key %s, value %s, error %s", string(key), string(value), err)
	}

	return nil
}

func (b *BadgerDB) SetWithTTL(key []byte, value []byte, ttl int64) error {
	wb := b.DB.NewWriteBatch()
	defer wb.Cancel()

	err := wb.SetEntry(badger.NewEntry(key, value).WithMeta(0).WithTTL(time.Duration(ttl * time.Second.Nanoseconds())))
	if err != nil {
		return fmt.Errorf("failed to write data to cache. key %s, value %s, error %s", string(key), string(value), err)
	}

	err = wb.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush data to cache. key %s, error %s", string(key), err)
	}

	return nil
}

func (b *BadgerDB) Get(key []byte) ([]byte, error) {
	var ival []byte
	err := b.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		ival, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read data to cache. key %s, error %s", string(key), err)
	}
	return ival, nil
}

func (b *BadgerDB) Has(key []byte) (bool, error) {
	var exist bool = false
	err := b.DB.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil {
			return err
		} else {
			exist = true
		}
		return err
	})
	// align with leveldb, if the key doesn't exist, leveldb returns nil
	if strings.HasSuffix(err.Error(), "not found") {
		err = nil
	}
	return exist, err
}

func (b *BadgerDB) Delete(key []byte) error {
	wb := b.DB.NewWriteBatch()
	defer wb.Cancel()
	return wb.Delete(key)
}

func (b *BadgerDB) IteratorKeysAndValues() error {
	err := b.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%x, value=%x\n\n", k, v)
				return nil
			})
			if err != nil {
				return err
			}

			time.Sleep(100 * time.Millisecond)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to iterator keys and values from the cache. error %s", err)
	}

	return nil
}

func (b *BadgerDB) IteratorKeys() error {
	err := b.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			fmt.Printf("key=%s\n", k)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to iterator keys from the cache. error %s", err)
	}

	return nil
}

func (b *BadgerDB) SeekWithPrefix(prefixStr string) error {
	err := b.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(prefixStr)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%s, value=%s\n", k, v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to seek prefix from the cache. prefix %s, error %s", prefixStr, err)
	}

	return nil
}
