package cmd

import (
	"fmt"
	"log"
	"math/big"

	"github.com/dgraph-io/badger"
	com "github.com/rockiecn/ens-log/common"
	"github.com/rockiecn/ens-log/db"
	"github.com/rockiecn/ens-log/grab"
	"github.com/rockiecn/ens-log/utils"
	"github.com/urfave/cli/v2"
)

// 100 blocks each time
const STEP uint64 = 10000

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

			// end
			if from.Cmp(to) >= 0 {
				return nil
			} else {
				fmt.Printf("======== grabbing logs for block: %d to %d\n", from, to)
				grab.GrabLogs(com.CLI, from, to, rDB, sDB)
			}
		}

		return nil
	},
}
