package btcadaptor

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/palletone/adaptor"
)

func TestCreateTransferTokenTx(t *testing.T) {
	//rpcParams := RPCParams{
	//	Host:      "localhost:18334",
	//	RPCUser:   "test",
	//	RPCPasswd: "123456",
	//	CertPath:  GCertPath,
	//}

	rpcParams := RPCParams{
		Host:      "123.126.106.86:28556",
		RPCUser:   "pallettest",
		RPCPasswd: "pallet123456",
		CertPath:  "./rpc-123.126.106.86-testnet.cert",
	}

	//
	var input adaptor.CreateTransferTokenTxInput
	input.FromAddress = "mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5"     //2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA
	input.ToAddress = "2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA"      //mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5
	input.Amount = adaptor.NewAmountAssetString("980000", "BTC") //dao, 0.0099 btc
	input.Fee = adaptor.NewAmountAssetString("10000", "BTC")     //dao,0.0001 btc
	//input.Extra

	input.ToAddress = "PalletOne"
	input.ToAddress = "P19z4r7G9MpZtaYMZcATWinTwXeGBj7fWTd"
	input.Amount = adaptor.NewAmountAssetString("0", "BTC") //op_return

	//idIndex, _ := hex.DecodeString("101d482b60cd3f74a61ce265d62e383456b9c21c84477931d207ea8f503d84cc01")
	//input.Extra = append(input.Extra, idIndex...)

	//{"transaction":"AQAAAAHMhD1Qj+oH0jF5R4QcwrlWNDgu1mXiHKZ0P81gK0gdEAEAAAAAAAAAAAJAQg8AAAAAABl2qRS93Jpi6bfDz9vhyBdSDiTjLDOfMoisQEIPAAAAAAAZdqkUDwjlW8/CB2MtLc/D1NtLbY2Rsi6IrAAAAAA=","extra":"EB1IK2DNP3SmHOJl1i44NFa5whyER3kx0gfqj1A9hMwB"}
	output, err := CreateTransferTokenTx(&input, &rpcParams, NETID_TEST)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		resultJSON, _ := json.Marshal(output)
		fmt.Println(string(resultJSON))
		fmt.Printf("Extra : %x\n", output.Extra)
		rawTxHex := fmt.Sprintf("%x", output.Transaction)
		fmt.Printf("rawTxHex : %s\n", rawTxHex)
		_, err = decodeRawTransaction(rawTxHex, NETID_TEST)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func TestCalcTxHash(t *testing.T) {
	tx, _ := hex.DecodeString("010000000144080d1b48b02a483b5341ba70b2660616d72d1aeb006518f3e261f7f651e48d00000000fdfd00004730440220475bc769dd39d5820131c3e3527f6710886c6852e66635b1f68aff0fdad88f5702206ce2ca0ac655f74ec5f01c0ca2ae9bc1228176bab1855a0a053c0c177bb5ae7c01483045022100dfa8811cee744502d183f67e8166480b0bcb1f0d4259f6fe3b2ee16aa0be568a022079aedff69da83654fe15bb0ab927510c737f772d1d048c219c25546d76af73eb014c69522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553ae000000000250c30000000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88acf07e0e000000000017a9147e037d8b8093a7cf3a6ec83aa8c852761a5d0cce8700000000")

	input := &adaptor.CalcTxHashInput{Transaction: tx}

	output, err := CalcTxHash(input)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		resultJSON, _ := json.Marshal(output)
		fmt.Println(string(resultJSON))
		fmt.Printf("%x\n", output.Hash)
	}
}

func TestGetBlockInfo(t *testing.T) {
	//rpcParams := RPCParams{
	//	Host:      "localhost:18334",
	//	RPCUser:   "test",
	//	RPCPasswd: "123456",
	//	CertPath:  GCertPath,
	//}

	rpcParams := RPCParams{
		Host:      "123.126.106.86:28556",
		RPCUser:   "pallettest",
		RPCPasswd: "pallet123456",
		CertPath:  "./rpc-123.126.106.86-testnet.cert",
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
	//rpcParams := RPCParams{
	//	Host:      "localhost:18334",
	//	RPCUser:   "test",
	//	RPCPasswd: "123456",
	//	CertPath:  GCertPath,
	//}

	rpcParams := RPCParams{
		Host:      "123.126.106.86:28556",
		RPCUser:   "pallettest",
		RPCPasswd: "pallet123456",
		CertPath:  "./rpc-123.126.106.86-testnet.cert",
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

	txIDHex := "6b2c4379b326757dd5b847f3c584170c5fe2649e6e33f962cf7e9826f77f07b6"

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
	//rpcParams := RPCParams{
	//	Host:      "localhost:18334",
	//	RPCUser:   "test",
	//	RPCPasswd: "123456",
	//	CertPath:  GCertPath,
	//}

	rpcParams := RPCParams{
		Host:      "123.126.106.86:28556",
		RPCUser:   "pallettest",
		RPCPasswd: "pallet123456",
		CertPath:  "./rpc-123.126.106.86-testnet.cert",
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
	//rpcParams := RPCParams{
	//	Host:      "localhost:18334",
	//	RPCUser:   "test",
	//	RPCPasswd: "123456",
	//	CertPath:  GCertPath,
	//}

	rpcParams := RPCParams{
		Host:      "123.126.106.86:28556",
		RPCUser:   "pallettest",
		RPCPasswd: "pallet123456",
		CertPath:  "./rpc-123.126.106.86-testnet.cert",
	}

	txIDHex := "52dab174e0d719704316c9301f146e1e90e7797ec8fe9f357a5fdfb0a62a1ab4"
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
		if 0 != len(output.Tx.AttachData) {
			fmt.Println("data : ", string(output.Tx.AttachData))
		}
	}
}

func TestGetTxBasicInfoHttp(t *testing.T) {
	txIDHex := "8b886fb5033d26c2bad728d73188e4eac46e2eb61260a2638b3330484498c576"
	txID, _ := hex.DecodeString(txIDHex)

	var input adaptor.GetTxBasicInfoInput
	input.TxID = txID
	fmt.Printf("%x\n", input.TxID)

	output, err := GetTxBasicInfoHttp(&input, NETID_TEST)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		resultJSON, _ := json.Marshal(output)
		fmt.Println(string(resultJSON))
	}
}
