package adaptorbtc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/palletone/adaptor"
)

func TestGetUnspendUTXO(t *testing.T) {
	//	parms := `{"address": "mxprH5bkXtn9tTTAxdQGPXrvruCUvsBNKt","minconf": 0,"maxconf": 999999,"maximumCount": 10}`
	//	parms := `{"address": "miZqthevf8LWguQmUR6EwynULqjKmYWxyY","minconf": 0,"maxconf": 999999,"maximumCount": 10}`
	//	parms := `{"address": "2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA","minconf": 0,"maxconf": 999999,"maximumCount": 10}`
	parms := `{"address": "mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5","minconf": 0,"maxconf": 999999,"maximumCount": 100}`

	rpcParams := RPCParams{
		Host:      "localhost:18334",
		RPCUser:   "zxl",
		RPCPasswd: "123456",
		//CertPath:  GCertPath,
	}

	//	testResult := "101d482b60cd3f74a61ce265d62e383456b9c21c84477931d207ea8f503d84cc"
	//	testResult := "cdc28467435bb3060333777e289adb200c033eee72c96c68cb9790534516f6eb"

	result := GetUnspendUTXO(parms, &rpcParams, NETID_TEST)
	//	if !strings.Contains(result, testResult) {
	//		t.Errorf("unexpected result - got: %v, "+"want: %v", result, testResult)
	//	}
	fmt.Println(result)
}

func TestGetBalance(t *testing.T) {
	//	parms := `{"address": "2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA","minconf": 1}`

	rpcParams := RPCParams{
		Host:      "localhost:18334",
		RPCUser:   "zxl",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	testResult := `{"value":0}`

	//	getBalanceParams := &adaptor.GetBalanceParams{"mxprH5bkXtn9tTTAxdQGPXrvruCUvsBNKt", 1}
	getBalanceParams := &adaptor.GetBalanceParams{"mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5", 0}
	//	getBalanceParams := &adaptor.GetBalanceParams{"miZqthevf8LWguQmUR6EwynULqjKmYWxyY", 1}
	//	getBalanceParams := &adaptor.GetBalanceParams{"2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA", 1}

	result, err := GetBalance(getBalanceParams, &rpcParams, NETID_TEST)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
		if !strings.Contains(result, testResult) {
			fmt.Printf("unexpected result - got: %v, "+"want: %v\n", result, testResult)
		}
	}
}

func TestGetTransactions(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18334",
		RPCUser:   "zxl",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	//	parms := `{
	//	    "account": "2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA",
	//	    "count": 10,
	//	    "skip": 0
	//  	}`

	var getTransactionsParams adaptor.GetTransactionsParams
	//	getTransactionsParams.Account = "2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA"
	//	getTransactionsParams.Account = "2NGDzMbWC7Q1tv3bHc9B8FytBbKEwXJSgkg"
	//	getTransactionsParams.Account = "2N2ApYikZS6mVUeWLVqVpDVtLWuE1ufwam2"
	getTransactionsParams.Account = "mxprH5bkXtn9tTTAxdQGPXrvruCUvsBNKt"
	getTransactionsParams.Count = 100

	//	testResult := "1696a64fc031f19fe1f80dc1c78c9c8af1bd018ad1bf82ddeb7e2edb900b7b89"

	result, err := GetTransactions(&getTransactionsParams, &rpcParams, NETID_TEST)
	//	if !strings.Contains(result, testResult) {
	//		t.Errorf("unexpected result - got: %v, "+"want: %v", result, testResult)
	//	}
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}
