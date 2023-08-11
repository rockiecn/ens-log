package common

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rockiecn/ens-log/db"
	"github.com/rockiecn/ens-log/utils"
)

// public
var (
	CLI        *ethclient.Client // endpoint
	ChainBlock uint64            // block number on chain
	LocalBlock uint64            // current updated block in local db

	DB_RECORD *db.BadgerDB
	DB_SYS    *db.BadgerDB
)

// db and connection
func init() {
	// new dbs
	DB_RECORD = db.New("./badger_record")
	DB_SYS = db.New("./badger_sys")

	// open
	DB_SYS.Open()
	defer DB_SYS.Close()
	// get localblock
	b, err := DB_SYS.Get([]byte("localblock"))
	if err != nil {
		log.Fatal(err)
	}
	LocalBlock := utils.BytestoUint64(b)
	fmt.Println("LocalBlock is : ", LocalBlock)

	// connect endpoint
	// mainnet
	//client, err := ethclient.Dial("https://mainnet.infura.io/v3/c7bbf3c290ac408694435532b6838308")
	// goerli
	CLI, err = ethclient.Dial("https://goerli.infura.io/v3/c7bbf3c290ac408694435532b6838308")
	if err != nil {
		log.Fatal(err)
	}
	// get block number
	ctx := context.Background()
	ChainBlock, err = CLI.BlockNumber(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Chain block is: ", ChainBlock)
}

// set local block to specified number
func SetLocalBlock(l uint64) error {
	db := db.New("./badger_sys")
	err := db.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	k := []byte("localblock")
	v := utils.Uint64toBytes(l)

	err = db.Set(k, v)
	if err != nil {
		return err
	}

	fmt.Println("local block is set to ", l)

	return nil
}

// get local block
func GetLocalBlock() (uint64, error) {
	db := db.New("./badger_sys")
	err := db.Open()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	k := []byte("localblock")
	b, err := db.Get(k)
	if err != nil {
		return 0, err
	}

	i := utils.BytestoUint64(b)

	return i, nil
}

// get local block
func Iter() error {
	db := db.New("./badger_record")
	err := db.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	db.IteratorKeysAndValues()

	return nil
}
