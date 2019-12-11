package btcadaptor

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/palletone/adaptor"
)

func TestGetBalance(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18334",
		RPCUser:   "test",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	//testResult := `{"value":0.999}`
	testResult := 0.999
	testResultUint64 := uint64(testResult * 1e8)

	//	getBalanceParams := &adaptor.GetBalanceParams{Address:"mxprH5bkXtn9tTTAxdQGPXrvruCUvsBNKt"}
	getBalanceParams := &adaptor.GetBalanceInput{Address: "mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5"}
	//	getBalanceParams := &adaptor.GetBalanceParams{Address:"miZqthevf8LWguQmUR6EwynULqjKmYWxyY"}
	//	getBalanceParams := &adaptor.GetBalanceParams{Address:"2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA"}

	result, err := GetBalance(getBalanceParams, &rpcParams, NETID_TEST)
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
	rpcParams := RPCParams{
		Host:      "localhost:18334",
		RPCUser:   "test",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
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
