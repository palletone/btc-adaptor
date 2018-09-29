package adaptorbtc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/palletone/adaptor"
)

func TestRawTransactionGen(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18332",
		RPCUser:   "zxl",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

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
	rawTransactionGenParams.Inputs = append(rawTransactionGenParams.Inputs, adaptor.Input{txid1, uint32(vout1)})
	rawTransactionGenParams.Inputs = append(rawTransactionGenParams.Inputs, adaptor.Input{txid2, uint32(vout2)})
	rawTransactionGenParams.Outputs = append(rawTransactionGenParams.Outputs, adaptor.Output{address, amount})

	testResult := "0100000002897b0b90db2e7eebdd82bfd18a01bdf18a9c8cc7c10df8e19ff131c04fa696160000000000ffffffffea60d9d11cc3e40e17f068ca1191ff5d3ec11d016393addf63305001ce813c990000000000ffffffff010c8aa30d000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88ac00000000"
	result, err := RawTransactionGen(&rawTransactionGenParams, &rpcParams, NETID_TEST)
	if !strings.Contains(result, testResult) {
		t.Errorf("unexpected result - got: %v, "+"want: %v", result, testResult)
	}
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}

func TestDecodeRawTransaction(t *testing.T) {

	rpcParams := RPCParams{
		Host:      "localhost:18332",
		RPCUser:   "zxl",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

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
	rawtx := "010000000145796293a6d01a72d2b160dbaefe27acb659ad72d1a6d47b2e204c7be221998d010000006a473044022020bc4bb2aa9419c7ed6bd2355c6ffd1a8a845ee0b03e00bfe7725cc4fdf32a3b022069b3275dc9ae5e6c66b2ead2afb18b50e029e5e1015ff8ad4c9c90a5c922a4660121020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cbffffffff0200e1f5050000000017a9147e037d8b8093a7cf3a6ec83aa8c852761a5d0cce8730033c05000000001976a914bddc9a62e9b7c3cfdbe1c817520e24e32c339f3288ac00000000"
	var decodeRawTransactionParams adaptor.DecodeRawTransactionParams
	decodeRawTransactionParams.Rawtx = rawtx
	result, err := DecodeRawTransaction(&decodeRawTransactionParams, &rpcParams)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}

func TestGetTransactionByHash(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18332",
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
