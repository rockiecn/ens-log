package grab

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	com "github.com/rockiecn/ens-log/common"
	"github.com/rockiecn/ens-log/db"
	"github.com/rockiecn/ens-log/utils"
)

// abi of ETHRegistrarController.sol
const ABI string = `[{"inputs":[{"internalType":"contract BaseRegistrarImplementation","name":"_base","type":"address"},{"internalType":"contract IPriceOracle","name":"_prices","type":"address"},{"internalType":"uint256","name":"_minCommitmentAge","type":"uint256"},{"internalType":"uint256","name":"_maxCommitmentAge","type":"uint256"},{"internalType":"contract ReverseRegistrar","name":"_reverseRegistrar","type":"address"},{"internalType":"contract INameWrapper","name":"_nameWrapper","type":"address"},{"internalType":"contract ENS","name":"_ens","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[{"internalType":"bytes32","name":"commitment","type":"bytes32"}],"name":"CommitmentTooNew","type":"error"},{"inputs":[{"internalType":"bytes32","name":"commitment","type":"bytes32"}],"name":"CommitmentTooOld","type":"error"},{"inputs":[{"internalType":"uint256","name":"duration","type":"uint256"}],"name":"DurationTooShort","type":"error"},{"inputs":[],"name":"InsufficientValue","type":"error"},{"inputs":[],"name":"MaxCommitmentAgeTooHigh","type":"error"},{"inputs":[],"name":"MaxCommitmentAgeTooLow","type":"error"},{"inputs":[{"internalType":"string","name":"name","type":"string"}],"name":"NameNotAvailable","type":"error"},{"inputs":[],"name":"ResolverRequiredWhenDataSupplied","type":"error"},{"inputs":[{"internalType":"bytes32","name":"commitment","type":"bytes32"}],"name":"UnexpiredCommitmentExists","type":"error"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"name","type":"string"},{"indexed":true,"internalType":"bytes32","name":"label","type":"bytes32"},{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint256","name":"baseCost","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"premium","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"expires","type":"uint256"}],"name":"NameRegistered","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"name","type":"string"},{"indexed":true,"internalType":"bytes32","name":"label","type":"bytes32"},{"indexed":false,"internalType":"uint256","name":"cost","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"expires","type":"uint256"}],"name":"NameRenewed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[],"name":"MIN_REGISTRATION_DURATION","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"}],"name":"available","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"commitment","type":"bytes32"}],"name":"commit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"commitments","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"duration","type":"uint256"},{"internalType":"bytes32","name":"secret","type":"bytes32"},{"internalType":"address","name":"resolver","type":"address"},{"internalType":"bytes[]","name":"data","type":"bytes[]"},{"internalType":"bool","name":"reverseRecord","type":"bool"},{"internalType":"uint16","name":"ownerControlledFuses","type":"uint16"}],"name":"makeCommitment","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"maxCommitmentAge","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"minCommitmentAge","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"nameWrapper","outputs":[{"internalType":"contract INameWrapper","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"prices","outputs":[{"internalType":"contract IPriceOracle","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_token","type":"address"},{"internalType":"address","name":"_to","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"recoverFunds","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"duration","type":"uint256"},{"internalType":"bytes32","name":"secret","type":"bytes32"},{"internalType":"address","name":"resolver","type":"address"},{"internalType":"bytes[]","name":"data","type":"bytes[]"},{"internalType":"bool","name":"reverseRecord","type":"bool"},{"internalType":"uint16","name":"ownerControlledFuses","type":"uint16"}],"name":"register","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"duration","type":"uint256"}],"name":"renew","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"duration","type":"uint256"}],"name":"rentPrice","outputs":[{"components":[{"internalType":"uint256","name":"base","type":"uint256"},{"internalType":"uint256","name":"premium","type":"uint256"}],"internalType":"struct IPriceOracle.Price","name":"price","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"reverseRegistrar","outputs":[{"internalType":"contract ReverseRegistrar","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes4","name":"interfaceID","type":"bytes4"}],"name":"supportsInterface","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"}],"name":"valid","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

var (
	// mainnet
	//controllerAddress = common.HexToAddress("0x253553366Da8546fC250F225fe3d25d0C782303b")
	// goerli
	controllerAddress = common.HexToAddress("0x283Af0B28c62C092C9727F1Ee09c02CA627EB7F5")

	// nameRegister event topic
	// new version
	//nameRegTopic = "0x69e37f151eb98a09618ddaa80c8cfaf1ce5996867c489f45b555b412271ebf27"
	// old version
	nameRegTopic = "0xca6abbe9d7f11422cb6ca7629fbf6fe9efb1c621f71ce8f02b9f2a230097404f"
)

// grab nameRegister logs for certain blocks
func GrabLogs(
	cli *ethclient.Client,
	fromBlock *big.Int,
	toBlock *big.Int,
	rDB *db.BadgerDB,
	sDB *db.BadgerDB,
) error {
	// create a query
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{
			controllerAddress,
		},
	}

	// from must less than to
	if fromBlock.Cmp(toBlock) >= 0 {
		return nil
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

	// get contract abi
	jsonABI, err := abi.JSON(strings.NewReader(string(ABI)))
	if err != nil {
		log.Fatal(err)
	}

	// check each log
	for i, vLog := range logs {

		//fmt.Printf("log index: %d\n", i)

		// current log's block
		//logBlock := vLog.BlockNumber
		//fmt.Println("log's block number:", logBlock)

		//fmt.Println("txHash:", vLog.TxHash)

		// NameRegistered (string name, index_topic_1 bytes32 label, index_topic_2 address owner, uint256 cost, uint256 expires)
		// get topics from log
		// 3 topics in all, the first is the hash of the event
		topics := make([]string, 3)
		for i := range vLog.Topics {
			topics[i] = vLog.Topics[i].Hex()
			//fmt.Printf("topic%d: %s\n", i, topics[i])
		}
		// record item
		label := topics[1]
		owner := topics[2]

		// NameRegistered event found
		cmp := strings.Compare(topics[0], nameRegTopic)
		if cmp == 0 {
			fmt.Println("NameRegistered found, log index: ", i)

			// unpack log data
			logData, err := jsonABI.Unpack("NameRegistered", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			// 3 data in event: name, cost, expires
			// name string
			name, ok := logData[0].(string)
			if !ok {
				log.Fatal("log data assertion failed")
			}
			//fmt.Println("name:", name)
			// cost uint256
			cost, ok := logData[1].(*big.Int)
			if !ok {
				log.Fatal("log data assertion failed")
			}
			_ = cost
			//fmt.Println("cost:", cost.String())
			// expires uint256
			expires, ok := logData[2].(*big.Int)
			if !ok {
				log.Fatal("log data assertion failed")
			}
			//fmt.Println("expires:", expires.String())

			// make record
			r := utils.RECORD{
				Name:    name,
				Label:   label,
				Owner:   owner,
				Expires: expires,
			}

			// encode k
			k, err := utils.Encode(r.Name)
			if err != nil {
				log.Fatal(err)
			}

			// illegal key length, caused by long name
			if len(k) > 65000 {
				fmt.Println("name too long to be record, skip it")
				break
			} else { // record this name
				// encode v
				v, err := utils.Encode(r)
				if err != nil {
					log.Fatal(err)
				}
				// set k-v for record
				err = rDB.Set(k, v)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

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
