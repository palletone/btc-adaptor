package adaptorbtc

import (
	"fmt"
	"testing"

	"github.com/palletone/adaptor"
)

func TestRawTransactionGen(t *testing.T) {
	//1.64835652 + 0.64085864 = 2.28921516 - fee = 2.28821516
	//	rawTransactionGenParams := `{
	//	    "inputs": [
	//			{
	//	           "txid": "1696a64fc031f19fe1f80dc1c78c9c8af1bd018ad1bf82ddeb7e2edb900b7b89",
	//	           "vout": 0
	//			},
	//			{
	//	           "txid": "993c81ce01503063dfad9363011dc13e5dff9111ca68f0170ee4c31cd1d960ea",
	//	           "vout": 0
	//			}
	//	    ],
	//	    "outputs": [
	//			{
	//	           "address": "mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5",
	//	           "amount": 2.28821516
	//			}
	//	    ],
	//	    "locktime": 0
	//		}`
	//	result, err := RawTransactionGen(rawTransactionGenParams, &rpcParams, NETID_TEST)

	txid1 := "1696a64fc031f19fe1f80dc1c78c9c8af1bd018ad1bf82ddeb7e2edb900b7b89"
	vout1 := 0
	txid2 := "993c81ce01503063dfad9363011dc13e5dff9111ca68f0170ee4c31cd1d960ea"
	vout2 := 0
	address := "mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5"
	amount := 2.28821516
	//
	var rawTransactionGenParams adaptor.RawTransactionGenParams
	rawTransactionGenParams.Inputs = append(rawTransactionGenParams.Inputs, adaptor.Input{Txid: txid1, Vout: uint32(vout1)})
	rawTransactionGenParams.Inputs = append(rawTransactionGenParams.Inputs, adaptor.Input{Txid: txid2, Vout: uint32(vout2)})
	rawTransactionGenParams.Outputs = append(rawTransactionGenParams.Outputs, adaptor.Output{address, amount})

	testResult := "0100000002897b0b90db2e7eebdd82bfd18a01bdf18a9c8cc7c10df8e19ff131c04fa69616000000000000000000ea60d9d11cc3e40e17f068ca1191ff5d3ec11d016393addf63305001ce813c99000000000000000000010c8aa30d000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88ac00000000"
	result, err := RawTransactionGen(&rawTransactionGenParams, NETID_TEST)
	if result.Rawtx != testResult {
		t.Errorf("unexpected result - got: %v, "+"want: %v", result, testResult)
	}
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}

func TestDecodeRawTransaction(t *testing.T) {
	//	testResult := `{
	//		"hex":"",
	//		"txid":"0bf2bbdabd7561fe035eb383d14e376f04690c62301cc78d89dd189f7e6c3a72",
	//		"version":1,
	//		"locktime":0,
	//		"vin":[{
	//			"txid":"132154398e312b69b62973f8f6a91797bba9996bc60dc1d7b1f8697df196088d",
	//			"vout":0,
	//			"scriptSig":{"asm":"","hex":""},
	//			"sequence":4294967295
	//		}],
	//		"vout":[{
	//			"value":0.98811339,
	//			"n":0,
	//			"scriptPubKey":{
	//			"asm":"OP_DUP OP_HASH160 bddc9a62e9b7c3cfdbe1c817520e24e32c339f32 OP_EQUALVERIFY OP_CHECKSIG",
	//			"hex":"76a914bddc9a62e9b7c3cfdbe1c817520e24e32c339f3288ac",
	//			"reqSigs":1,
	//			"type":"pubkeyhash",
	//			"addresses":["mxprH5bkXtn9tTTAxdQGPXrvruCUvsBNKt"]
	//			}
	//		}]
	//	}
	//	`

	//	rawtx := "0100000002897b0b90db2e7eebdd82bfd18a01bdf18a9c8cc7c10df8e19ff131c04fa69616000000006b48304502210085b216c64e2c5311dfd4ee03038c175832f9531a8cadd5342cdc1896079c891e022046e0097c0c9b552e128851d264699e35b26b54c07e421d1e1a9bc148f2818b070121029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef5ffffffffea60d9d11cc3e40e17f068ca1191ff5d3ec11d016393addf63305001ce813c99000000006a47304402207d9ab909748b2a7e869e575a6ebf6814ea40b906356992b063f191cd3ae0b10102204593e60f7728dd5ae6c380d7bb8e04682e13f95e8a51182b2d65863386fe43bb0121029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef5ffffffff01f85b0c09000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88ac00000000"
	rawtx := "01000000012ab6af95c3edfc093f0c6e80b5080d5de31df4d46f775b8d5515e8c1ac25f9b8010000006c493046022100fa192625f84cdf22692338388e632157d7c3f7c24b596e2024926c2e5035d532022100af81d97b3d630dcdf40cf9eca7f41612297585cf8c2e579705306baca403a64001210313c403e04becbcb83e93fbdd9eb9d1b04d9479bfc0864ef46c49a6ca266b6f1fffffffff015113f204000000001976a914deb6c6367b4631fd3eeddcabc500dd5af64ce20888ac00000000"
	var decodeRawTransactionParams adaptor.DecodeRawTransactionParams
	decodeRawTransactionParams.Rawtx = rawtx
	result, err := DecodeRawTransaction(&decodeRawTransactionParams, NETID_TEST)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}

func TestGetTransactionByHash(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18334",
		RPCUser:   "zxl",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	var getTransactionParams adaptor.GetTransactionByHashParams
	getTransactionParams.TxHash = "39176cc119da6a491472ef598335305839b9ccb6b7d8cd635c863456f8b09917"

	result, err := GetTransactionByHash(&getTransactionParams, &rpcParams)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}

}

func TestGetTransactionHttp(t *testing.T) {
	var getTransactionParams adaptor.GetTransactionHttpParams
	getTransactionParams.TxHash = "39176cc119da6a491472ef598335305839b9ccb6b7d8cd635c863456f8b09917"

	//{"txid":"39176cc119da6a491472ef598335305839b9ccb6b7d8cd635c863456f8b09917","confirms":75542,"inputs":[{"txid":"6c4d7711f71dc5d075d9e10583351bb7ee530b44c8fd581cb97da91ea31d88cf","vout":0}],"outputs":[{"index":0,"addr":"2NGDzMbWC7Q1tv3bHc9B8FytBbKEwXJSgkg","value":0.60085864}]}

	result, err := GetTransactionHttp(&getTransactionParams, NETID_TEST)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}

}
