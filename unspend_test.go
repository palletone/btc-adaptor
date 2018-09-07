package adaptorbtc

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetUnspendUTXO(t *testing.T) {
	//	parms := `{"addresses": ["mxprH5bkXtn9tTTAxdQGPXrvruCUvsBNKt"],"minconf": 0,"maxconf": 999999,"maximumCount": 10}`
	parms := `{"addresses": ["mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5"],"minconf": 0,"maxconf": 999999,"maximumCount": 10}`
	//	parms := `{"addresses": ["miZqthevf8LWguQmUR6EwynULqjKmYWxyY"],"minconf": 0,"maxconf": 999999,"maximumCount": 10}`
	//	parms := `{"addresses": ["2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA"],"minconf": 0,"maxconf": 999999,"maximumCount": 10}`

	rpcParams := RPCParams{
		Host:      "localhost:18332",
		RPCUser:   "zxl",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
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
		Host:      "localhost:18332",
		RPCUser:   "zxl",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	testResult := `{"value":1}`

	//	getBalanceParams := &GetBalanceParams{"mxprH5bkXtn9tTTAxdQGPXrvruCUvsBNKt", 1}
	getBalanceParams := &GetBalanceParams{"mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5", 1}
	//	getBalanceParams := &GetBalanceParams{"miZqthevf8LWguQmUR6EwynULqjKmYWxyY", 1}
	//	getBalanceParams := &GetBalanceParams{"2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA", 1}

	result, err := GetBalance(getBalanceParams, &rpcParams, NETID_TEST)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
		if !strings.Contains(result, testResult) {
			t.Errorf("unexpected result - got: %v, "+"want: %v", result, testResult)
		}
	}
}

func TestImportMultisig(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18332",
		RPCUser:   "zxl",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	//	params := `{
	//    "publicKeys": ["03940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d5410514","020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb","029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef5"],
	//    "mRequires": 2,
	//	"walletPasswd":"1"
	//  	}
	//	`

	pubkeyAlice := "03940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d5410514"
	pubkeyBob := "020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb"
	pubkeyPallet := "029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef5"
	var importMultisigParams ImportMultisigParams
	importMultisigParams.PublicKeys = append(importMultisigParams.PublicKeys, pubkeyAlice)
	importMultisigParams.PublicKeys = append(importMultisigParams.PublicKeys, pubkeyBob)
	importMultisigParams.PublicKeys = append(importMultisigParams.PublicKeys, pubkeyPallet)
	importMultisigParams.MRequires = 2
	importMultisigParams.WalletPasswd = "1"

	//btcwallet can't be closed when rescanning, otherwise failed to get UTXO
	result, err := ImportMultisig(&importMultisigParams, &rpcParams, NETID_TEST)
	if !strings.Contains(result, "true") {
		t.Errorf("unexpected result - got: %v, want: true", result)
	}
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}

func TestGetTransactions(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18332",
		RPCUser:   "zxl",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	//	parms := `{
	//	    "account": "2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA",
	//	    "count": 10,
	//	    "skip": 0
	//  	}`

	var getTransactionsParams GetTransactionsParams
	//	getTransactionsParams.Account = "2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA"
	getTransactionsParams.Account = "2NGDzMbWC7Q1tv3bHc9B8FytBbKEwXJSgkg"
	getTransactionsParams.Count = 10

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
