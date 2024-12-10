package grabber

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	com "github.com/rockiecn/ens-log/common"
	"github.com/rockiecn/ens-log/database"
	"github.com/rockiecn/ens-log/db"
	"github.com/rockiecn/ens-log/logs"
	"github.com/rockiecn/ens-log/utils"
)

var (
	// blockNumber = big.NewInt(0)
	logger = logs.Logger("dumper")
)

type Grabber struct {
	Endpoint     string
	ABI          abi.ABI                       // controller abi
	Address      common.Address                // controller address
	NameRegTopic string                        // nameRegister event topic
	EventNameMap map[common.Hash]string        // map for event names in contract
	IndexedMap   map[common.Hash]abi.Arguments // map for all topic args of an event
}

type RegisterEvent struct {
	Name     string
	Label    [32]byte
	Owner    common.Address
	BaseCost *big.Int
	Premium  *big.Int
	Expires  *big.Int
}

func NewGrabber(ep string, con_abi string, con_addr string, nameRegTopic string) (g *Grabber, err error) {
	// ABI
	a, err := abi.JSON(strings.NewReader(string(con_abi)))
	if err != nil {
		return nil, err
	}

	g = &Grabber{
		Endpoint:     ep,
		ABI:          a,
		Address:      common.HexToAddress(con_addr),
		NameRegTopic: nameRegTopic,
		EventNameMap: make(map[common.Hash]string),
		IndexedMap:   make(map[common.Hash]abi.Arguments),
	}

	// each event
	for name, event := range a.Events {
		// save event name
		g.EventNameMap[event.ID] = name

		var indexed abi.Arguments
		// each topic
		for _, arg := range g.ABI.Events[name].Inputs {
			if arg.Indexed {
				indexed = append(indexed, arg)
			}
		}
		// save topics for an event
		g.IndexedMap[event.ID] = indexed
	}

	return g, nil
}

// grab nameRegister logs
func (g *Grabber) GrabLogs(cli *ethclient.Client, fromBlock *big.Int, toBlock *big.Int, rDB *db.BadgerDB, sDB *db.BadgerDB) error {
	// from must less than to
	if fromBlock.Cmp(toBlock) >= 0 {
		return nil
	}

	// create a query
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{
			g.Address,
		},
	}

	// filter logs
	logs, err := cli.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("log number:", len(logs))

	// if no logs in these blocks, set local block to toBlock and return
	if len(logs) == 0 {
		fmt.Println("no logs in these blocks.")

		// update local block
		if toBlock.Uint64() < com.ChainBlock {
			k := []byte("localblock")
			v := utils.Uint64toBytes(toBlock.Uint64())
			sDB.Set(k, v)
		} else {
			k := []byte("localblock")
			v := utils.Uint64toBytes(com.ChainBlock)
			sDB.Set(k, v)
		}

		return nil
	}

	var count uint64

	// check each log
	for _, vLog := range logs {
		// NameRegistered event found
		if strings.Compare(vLog.Topics[0].Hex(), g.NameRegTopic) == 0 {
			//fmt.Println("NameRegistered found, log index: ", i)

			// unpack a log's data
			var out RegisterEvent
			err := g.unpack(vLog, &out)
			if err != nil {
				return err
			}

			// fmt.Println("name: ", out.Name)
			// fmt.Println("label: ", hex.EncodeToString(out.Label[:]))
			// fmt.Println("owner: ", out.Owner.String())
			// fmt.Println("cost: ", out.BaseCost.String())
			// fmt.Println("expires: ", out.Expires.String())

			if len(out.Name) > 64 {
				fmt.Println("name too long, skipped")
				continue
			}

			count++
			// store record
			database.Insert(out.Name, hex.EncodeToString(out.Label[:]), out.Owner.String(), out.BaseCost.String(), out.Expires.String())
		}
	}

	fmt.Printf("%d logs stored into database\n", count)

	// update local block
	if toBlock.Uint64() < com.ChainBlock {
		k := []byte("localblock")
		v := utils.Uint64toBytes(toBlock.Uint64())
		sDB.Set(k, v)
	} else {
		k := []byte("localblock")
		v := utils.Uint64toBytes(com.ChainBlock)
		sDB.Set(k, v)
	}

	// eventSignature := []byte("Transfer(address indexed, address indexed, uint256)")
	// hash := crypto.Keccak256Hash(eventSignature)

	return nil
}

// unpack a log
func (g *Grabber) unpack(log types.Log, out interface{}) error {
	// get event name from map with hash
	eventName := g.EventNameMap[log.Topics[0]]
	// get all topics
	indexed := g.IndexedMap[log.Topics[0]]

	logger.Debug("event name: ", eventName)

	// parse data
	err := g.ABI.UnpackIntoInterface(out, eventName, log.Data)
	if err != nil {
		return err
	}
	logger.Debug("unpack out(no topics):", out)

	// parse topic
	logger.Debug("parse topic")
	err = abi.ParseTopics(out, indexed, log.Topics[1:])
	if err != nil {
		return err
	}
	logger.Debug("unpack out(with topics):", out)

	return nil
}
