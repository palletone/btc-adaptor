package btcadaptor

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/palletone/adaptor"
)

func TestCreateTransferTokenTx(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18334",
		RPCUser:   "test",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	//
	var input adaptor.CreateTransferTokenTxInput
	input.FromAddress = "mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5"
	input.ToAddress = "2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA"
	input.Amount = adaptor.NewAmountAssetString("1000000", "BTC") //dao, 0.1 btc
	input.Fee = adaptor.NewAmountAssetString("10000", "BTC")      //dao,0.0001 btc

	//idIndex, _ := hex.DecodeString("101d482b60cd3f74a61ce265d62e383456b9c21c84477931d207ea8f503d84cc01")
	//input.Extra = append(input.Extra, idIndex...)

	//{"transaction":"AQAAAAHMhD1Qj+oH0jF5R4QcwrlWNDgu1mXiHKZ0P81gK0gdEAEAAAAAAAAAAAJAQg8AAAAAABl2qRS93Jpi6bfDz9vhyBdSDiTjLDOfMoisQEIPAAAAAAAZdqkUDwjlW8/CB2MtLc/D1NtLbY2Rsi6IrAAAAAA=","extra":"EB1IK2DNP3SmHOJl1i44NFa5whyER3kx0gfqj1A9hMwB"}
	output, err := CreateTransferTokenTx(&input, &rpcParams, NETID_TEST)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		resultJSON, _ := json.Marshal(output)
		fmt.Println(string(resultJSON))
		fmt.Printf("%x\n", output.Transaction)
	}
}

func TestGetBlockInfo(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18334",
		RPCUser:   "test",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	var input adaptor.GetBlockInfoInput

	input.Latest = true

	blkIDHex := "0000000094305b3c6173d333a9c9732a7b3af7534fade0fbdff95d932c3b1a94"
	blkID, _ := hex.DecodeString(blkIDHex)
	input.BlockID = blkID
	fmt.Printf("%x\n", input.BlockID)

	input.Height = 1

	output, err := GetBlockInfo(&input, &rpcParams)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		resultJSON, _ := json.Marshal(output)
		fmt.Println(string(resultJSON))
	}
}

func TestGetPalletOneMappingAddress(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18334",
		RPCUser:   "test",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	//txIDStr := "8b886fb5033d26c2bad728d73188e4eac46e2eb61260a2638b3330484498c576"
	//txIDByte, _ := hex.DecodeString(txIDStr)
	//
	//index := 128
	//str := fmt.Sprintf("%x", index)
	//fmt.Println(str)
	//indexByte, _ := hex.DecodeString(str)
	//fmt.Println(indexByte)
	//
	//fmt.Printf("%x\n", txIDByte)
	//txIDByte = append(txIDByte, indexByte[0])
	//fmt.Printf("%x\n", txIDByte)
	//
	//n := uint64(txIDByte[len(txIDByte)-1])
	//fmt.Println(n)
	//return

	txIDHex := "8b886fb5033d26c2bad728d73188e4eac46e2eb61260a2638b3330484498c576"

	var input adaptor.GetPalletOneMappingAddressInput
	input.MappingDataSource = txIDHex

	output, err := GetPalletOneMappingAddress(&input, &rpcParams)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		resultJSON, _ := json.Marshal(output)
		fmt.Println(string(resultJSON))
	}
}

func TestGetTxBasicInfo(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18334",
		RPCUser:   "test",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	txIDHex := "8b886fb5033d26c2bad728d73188e4eac46e2eb61260a2638b3330484498c576"
	txID, _ := hex.DecodeString(txIDHex)

	var input adaptor.GetTxBasicInfoInput
	input.TxID = txID
	fmt.Printf("%x\n", input.TxID)

	output, err := GetTxBasicInfo(&input, &rpcParams)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		resultJSON, _ := json.Marshal(output)
		fmt.Println(string(resultJSON))
	}
}

func TestGetTransferTx(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18334",
		RPCUser:   "test",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	txIDHex := "8b886fb5033d26c2bad728d73188e4eac46e2eb61260a2638b3330484498c576"
	txID, _ := hex.DecodeString(txIDHex)

	var input adaptor.GetTransferTxInput
	input.TxID = txID
	fmt.Printf("%x\n", input.TxID)

	output, err := GetTransferTx(&input, &rpcParams)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		resultJSON, _ := json.Marshal(output)
		fmt.Println(string(resultJSON))
	}
}

//func TestGetTransactionHttp(t *testing.T) {
//	var getTransactionParams adaptor.GetTransactionHttpParams
//	getTransactionParams.TxHash = "39176cc119da6a491472ef598335305839b9ccb6b7d8cd635c863456f8b09917"
//
//	//{"txid":"39176cc119da6a491472ef598335305839b9ccb6b7d8cd635c863456f8b09917","confirms":75542,"inputs":[{"txid":"6c4d7711f71dc5d075d9e10583351bb7ee530b44c8fd581cb97da91ea31d88cf","vout":0}],"outputs":[{"index":0,"addr":"2NGDzMbWC7Q1tv3bHc9B8FytBbKEwXJSgkg","value":0.60085864}]}
//
//	result, err := GetTransactionHttp(&getTransactionParams, NETID_TEST)
//	if err != nil {
//		fmt.Println(err.Error())
//	} else {
//		fmt.Println(result)
//	}
//
//}
