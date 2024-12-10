package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rockiecn/ens-log/database"
	"github.com/rockiecn/ens-log/db"
	"github.com/rockiecn/ens-log/utils"
)

// public
var (
	CLI        *ethclient.Client // endpoint
	ChainBlock uint64            // current block number on chain
	LocalBlock uint64            // current updated block in local db

	DB_RECORD *db.BadgerDB
	DB_SYS    *db.BadgerDB
)

const (
	// endpoint
	Endpoint string = "https://rpc.ankr.com/eth" // mainnet
	//Endpoint string= "https://rpc.ankr.com/eth_sepolia" // sepolia

	// controller addr
	Address string = "0x253553366Da8546fC250F225fe3d25d0C782303b" // mainnet
	//Address string= "0xFED6a969AaA60E4961FCD3EBF1A2e8913ac65B72" // sepolia

	// abi of ETHRegistrarController.sol
	ABI string = `[{"inputs":[{"internalType":"contract BaseRegistrarImplementation","name":"_base","type":"address"},{"internalType":"contract IPriceOracle","name":"_prices","type":"address"},{"internalType":"uint256","name":"_minCommitmentAge","type":"uint256"},{"internalType":"uint256","name":"_maxCommitmentAge","type":"uint256"},{"internalType":"contract ReverseRegistrar","name":"_reverseRegistrar","type":"address"},{"internalType":"contract INameWrapper","name":"_nameWrapper","type":"address"},{"internalType":"contract ENS","name":"_ens","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[{"internalType":"bytes32","name":"commitment","type":"bytes32"}],"name":"CommitmentTooNew","type":"error"},{"inputs":[{"internalType":"bytes32","name":"commitment","type":"bytes32"}],"name":"CommitmentTooOld","type":"error"},{"inputs":[{"internalType":"uint256","name":"duration","type":"uint256"}],"name":"DurationTooShort","type":"error"},{"inputs":[],"name":"InsufficientValue","type":"error"},{"inputs":[],"name":"MaxCommitmentAgeTooHigh","type":"error"},{"inputs":[],"name":"MaxCommitmentAgeTooLow","type":"error"},{"inputs":[{"internalType":"string","name":"name","type":"string"}],"name":"NameNotAvailable","type":"error"},{"inputs":[],"name":"ResolverRequiredWhenDataSupplied","type":"error"},{"inputs":[{"internalType":"bytes32","name":"commitment","type":"bytes32"}],"name":"UnexpiredCommitmentExists","type":"error"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"name","type":"string"},{"indexed":true,"internalType":"bytes32","name":"label","type":"bytes32"},{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint256","name":"baseCost","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"premium","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"expires","type":"uint256"}],"name":"NameRegistered","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"name","type":"string"},{"indexed":true,"internalType":"bytes32","name":"label","type":"bytes32"},{"indexed":false,"internalType":"uint256","name":"cost","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"expires","type":"uint256"}],"name":"NameRenewed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[],"name":"MIN_REGISTRATION_DURATION","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"}],"name":"available","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"commitment","type":"bytes32"}],"name":"commit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"commitments","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"duration","type":"uint256"},{"internalType":"bytes32","name":"secret","type":"bytes32"},{"internalType":"address","name":"resolver","type":"address"},{"internalType":"bytes[]","name":"data","type":"bytes[]"},{"internalType":"bool","name":"reverseRecord","type":"bool"},{"internalType":"uint16","name":"ownerControlledFuses","type":"uint16"}],"name":"makeCommitment","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"maxCommitmentAge","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"minCommitmentAge","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"nameWrapper","outputs":[{"internalType":"contract INameWrapper","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"prices","outputs":[{"internalType":"contract IPriceOracle","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_token","type":"address"},{"internalType":"address","name":"_to","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"recoverFunds","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"duration","type":"uint256"},{"internalType":"bytes32","name":"secret","type":"bytes32"},{"internalType":"address","name":"resolver","type":"address"},{"internalType":"bytes[]","name":"data","type":"bytes[]"},{"internalType":"bool","name":"reverseRecord","type":"bool"},{"internalType":"uint16","name":"ownerControlledFuses","type":"uint16"}],"name":"register","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"duration","type":"uint256"}],"name":"renew","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"duration","type":"uint256"}],"name":"rentPrice","outputs":[{"components":[{"internalType":"uint256","name":"base","type":"uint256"},{"internalType":"uint256","name":"premium","type":"uint256"}],"internalType":"struct IPriceOracle.Price","name":"price","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"reverseRegistrar","outputs":[{"internalType":"contract ReverseRegistrar","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes4","name":"interfaceID","type":"bytes4"}],"name":"supportsInterface","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"}],"name":"valid","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

	// nameRegister event topic
	NameRegTopic string = "0x69e37f151eb98a09618ddaa80c8cfaf1ce5996867c489f45b555b412271ebf27" // new version
	//NameRegTopic string= "0xca6abbe9d7f11422cb6ca7629fbf6fe9efb1c621f71ce8f02b9f2a230097404f" // old version
)

// new db and eth connection
func init() {
	fmt.Println("init db and connect eth")

	database.Open()

	// new dbs
	DB_RECORD = db.New("./badger_record")
	DB_SYS = db.New("./badger_sys")

	// open sys db
	DB_SYS.Open()
	defer DB_SYS.Close()

	var LocalBlock uint64

	// get localblock
	b, err := DB_SYS.Get([]byte("localblock"))
	if err != nil {
		fmt.Println("localblock not found in db, set it to 0")
		LocalBlock = 0
		SetLocalBlock(LocalBlock)
	} else {
		LocalBlock = utils.BytestoUint64(b)
	}

	fmt.Println("LocalBlock in db is : ", LocalBlock)

	// // connect endpoint
	// fmt.Println("connecting eth chain")
	// CLI, err = ethclient.Dial(Endpoint)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // get current block number
	// ctx := context.Background()
	// ChainBlock, err = CLI.BlockNumber(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("current block number is: ", ChainBlock)
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

	fmt.Println("setting localblock")

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

	fmt.Println("getting localblock")

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
