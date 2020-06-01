package btcadaptor

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/palletone/adaptor"
)

func TestGetBalance(t *testing.T) {
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

	//return
	//testResult := `{"value":0.999}`
	testResult := 0.999
	testResultUint64 := uint64(uint64(decimal.NewFromFloat(testResult).Mul(decimal.New(1, 8)).IntPart()))

	//	input := &adaptor.GetBalanceInput{Address:"mxprH5bkXtn9tTTAxdQGPXrvruCUvsBNKt"}
	//	input := &adaptor.GetBalanceInput{Address:"miZqthevf8LWguQmUR6EwynULqjKmYWxyY"}
	input := &adaptor.GetBalanceInput{Address: "mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5"}
	//input := &adaptor.GetBalanceInput{Address: "2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA"}

	result, err := GetBalance(input, &rpcParams, NETID_TEST)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result.Balance.Amount.Uint64())
		if result.Balance.Amount.Uint64() != testResultUint64 {
			fmt.Printf("unexpected result - got: %v, "+"want: %v\n", result, testResult)
		}
	}
}

//func TestGetBalanceHttp(t *testing.T) {
//	//testResult := `{"value":16.0191}`
//	testResult := 16.0191
//
//	//getBalanceParams := &adaptor.GetBalanceParams{"mp6DHyYNuD28aiE1MQdKuRAdH7ZNydqUBC", 0}
//	getBalanceParams := &adaptor.GetBalanceHttpParams{"mp6DHyYNuD28aiE1MQdKuRAdH7ZNydqUBC", 6}
//	//getBalanceParams := &adaptor.GetBalanceParams{"tb1q73un52phlrsug2r35fgnrlme987tr3cgm88k8j", 0}//?
//	//getBalanceParams := &adaptor.GetBalanceParams{"1DEP8i3QJCsomS4BSMY2RpU1upv62aGvhD", 0}
//
//	result, err := GetBalanceHttp(getBalanceParams, NETID_TEST)
//	if err != nil {
//		fmt.Println(err.Error())
//	} else {
//		fmt.Println(result)
//		if result.Value != testResult {
//			fmt.Printf("unexpected result - got: %v, "+"want: %v\n", result, testResult)
//		}
//	}
//}

func TestGetTransactions(t *testing.T) {
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

	var input adaptor.GetAddrTxHistoryInput
	input.FromAddress = "mxprH5bkXtn9tTTAxdQGPXrvruCUvsBNKt"

	output, err := GetTransactions(&input, &rpcParams, NETID_TEST)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		resultJSON, _ := json.Marshal(output)
		fmt.Println(string(resultJSON))
		fmt.Println(output.Count)
	}
}
