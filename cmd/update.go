package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/dgraph-io/badger"
	"github.com/ethereum/go-ethereum/ethclient"
	com "github.com/rockiecn/ens-log/common"
	"github.com/rockiecn/ens-log/db"
	"github.com/rockiecn/ens-log/grabber"
	"github.com/rockiecn/ens-log/utils"
	"github.com/urfave/cli/v2"
)

// 100 blocks each time
const STEP uint64 = 5000

// update local db
var UpdateCmd = &cli.Command{
	Name:  "update",
	Usage: "update db",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		fmt.Println("updating")

		rDB := db.New("./badger_record")
		err := rDB.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer rDB.Close()

		sDB := db.New("./badger_sys")
		err = sDB.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer sDB.Close()

		// new grabber
		fmt.Println("new grabber")
		g, err := grabber.NewGrabber(com.Endpoint, com.ABI, com.Address, com.NameRegTopic)
		if err != nil {
			panic(err)
		}

		// connect chain
		cli, err := ethclient.Dial(g.Endpoint)
		if err != nil {
			log.Fatal(err)
		}

		// get current block number
		ctx := context.Background()
		com.ChainBlock, err = cli.BlockNumber(ctx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("current block number is: ", com.ChainBlock)

		// update local db
		for com.LocalBlock < com.ChainBlock {
			// get localblock
			k := []byte("localblock")
			b, err := sDB.Get(k)
			if err != nil {
				// get localblock in a view
				// no record for localblock, set to 0
				if err == badger.ErrKeyNotFound {
					// set localblock to 0
					err = sDB.Set(k, []byte("0"))
					if err != nil {
						log.Fatal(err)
					}
					com.LocalBlock = 0
					fmt.Println("LocalBlock is : ", com.LocalBlock)
				} else { // other err
					log.Fatal(err)
				}
			}
			com.LocalBlock = utils.BytestoUint64(b)

			// grab for 10000 blocks
			from := new(big.Int).SetUint64(com.LocalBlock)
			to := new(big.Int).SetUint64(com.LocalBlock + STEP)
			// limit to chain block
			if to.Uint64() > com.ChainBlock {
				to.SetUint64(com.ChainBlock)
			}

			// update end
			if from.Cmp(to) >= 0 {
				fmt.Println("grab end")
				return nil
			} else {
				fmt.Printf("======== grabbing logs for block: %d to %d\n", from, to)
				g.GrabLogs(cli, from, to, rDB, sDB)
			}
		}

		return nil
	},
}
