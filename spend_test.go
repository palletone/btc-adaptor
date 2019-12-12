package btcadaptor

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/palletone/adaptor"
)

func TestHashMessage(t *testing.T) {
	var input adaptor.HashMessageInput
	input.Message = []byte("Hello, World!")
	output, err := HashMessage(&input)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%x\n", string(output.Hash))
	}
}

func TestSignMessage(t *testing.T) {
	keyHex := "d0e26e9189b9f047036ed21294c8f36d41df6b51852fc932595d849d727223d0"
	key, _ := hex.DecodeString(keyHex)

	var input adaptor.SignMessageInput
	input.Message = []byte("Hello, World!")
	input.PrivateKey = key
	output, err := SignMessage(&input)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%x\n", string(output.Signature))
	}
}

func TestVerifySignature(t *testing.T) {
	pubKeyHex := "029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef5"
	pubKey, _ := hex.DecodeString(pubKeyHex)

	var input adaptor.VerifySignatureInput
	input.Message = []byte("Hello, World!")
	input.PublicKey = pubKey
	input.Signature = []byte("H7GJvdYPHxjeJWkAj5negCzcoW2JjvcY93UF5wMqW0D5etuRkK0qDr7KflmcZrIkflSzb4I6X8ZqSZHEXYG3jSA=")
	output, err := VerifySignature(&input)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(output.Pass)
	}
}

func TestSignTransaction(t *testing.T) {
	keyHex := "ac18d6ffa4e006ee5c297e14962d213030910a045dccc9d393686eb1613a1477"
	//keyHex := "5102a03540efe05623c25fb35a2b250466d15b302caf04f9523401b96fae5cda"
	key, _ := hex.DecodeString(keyHex)

	tx, _ := hex.DecodeString("01000000016af6e1f77fbff03a3439d465c4ceb7871f4722dd1463e745129a21bf80670f49000000000000000000020000000000000000256a235031397a34723747394d705a7461594d5a63415457696e5477586547426a376657546420f40e00000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88ac00000000")

	addrOrRedeem := ""
	addrOrRedeem = "522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553ae"
	addrOrRedeem = "mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5"

	input := &adaptor.SignTransactionInput{PrivateKey: key, Transaction: tx, Extra: []byte(addrOrRedeem)}
	output, err := SignTransaction(input, NETID_TEST)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		resultJSON, _ := json.Marshal(output)
		fmt.Println(string(resultJSON))
		fmt.Printf("%x\n", output.SignedTx)
	}
}

func TestSendTransaction(t *testing.T) {
	rpcParams := RPCParams{
		Host:      "localhost:18334",
		RPCUser:   "test",
		RPCPasswd: "123456",
		CertPath:  GCertPath,
	}

	tx, _ := hex.DecodeString("01000000016af6e1f77fbff03a3439d465c4ceb7871f4722dd1463e745129a21bf80670f49000000006a473044022012f857974deeacf8cc255b31c89971ef280be070bcd1390f49f2dec083c600360220410ba0580b53166bf165c124dda35539dcf765c78232e0cd273287f0844985ff0121020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb00000000020000000000000000256a235031397a34723747394d705a7461594d5a63415457696e5477586547426a376657546420f40e00000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88ac00000000")

	input := &adaptor.SendTransactionInput{Transaction: tx}
	output, err := SendTransaction(input, &rpcParams)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		resultJSON, _ := json.Marshal(output)
		fmt.Println(string(resultJSON))
		fmt.Printf("%x\n", output.TxID)
	}
}

//func TestSendTransactionHttp(t *testing.T) {
//	paramsSend := adaptor.SendTransactionHttpParams{
//		TransactionHex: "010000000145796293a6d01a72d2b160dbaefe27acb659ad72d1a6d47b2e204c7be221998d010000006a473044022020bc4bb2aa9419c7ed6bd2355c6ffd1a8a845ee0b03e00bfe7725cc4fdf32a3b022069b3275dc9ae5e6c66b2ead2afb18b50e029e5e1015ff8ad4c9c90a5c922a4660121020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cbffffffff0200e1f5050000000017a9147e037d8b8093a7cf3a6ec83aa8c852761a5d0cce8730033c05000000001976a914bddc9a62e9b7c3cfdbe1c817520e24e32c339f3288ac00000000",
//	}
//
//	result, err := SendTransactionHttp(&paramsSend, NETID_TEST)
//	if err != nil {
//		fmt.Println(err.Error())
//	} else {
//		fmt.Println(result)
//	}
//}

func TestBindTxAndSignature(t *testing.T) {
	//raw tx
	rawTx, _ := hex.DecodeString("010000000144080d1b48b02a483b5341ba70b2660616d72d1aeb006518f3e261f7f651e48d00000000000000000001301b0f00000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88ac00000000")
	signenTx1, _ := hex.DecodeString("010000000144080d1b48b02a483b5341ba70b2660616d72d1aeb006518f3e261f7f651e48d00000000b500483045022100f2e159bbcbe28829d75f68ac3cdda9851df5478a23c090b58e9a6189d3f7e1bc0220664fde594d470527f70a10542384a241487fcd209291819093a5e6d7c3f83c94014c69522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553ae0000000001301b0f00000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88ac00000000")
	//signenTx2, _ := hex.DecodeString("010000000144080d1b48b02a483b5341ba70b2660616d72d1aeb006518f3e261f7f651e48d00000000b500483045022100dfa8811cee744502d183f67e8166480b0bcb1f0d4259f6fe3b2ee16aa0be568a022079aedff69da83654fe15bb0ab927510c737f772d1d048c219c25546d76af73eb014c69522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553ae000000000250c30000000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88acf07e0e000000000017a9147e037d8b8093a7cf3a6ec83aa8c852761a5d0cce8700000000")
	signenTx2, _ := hex.DecodeString("010000000144080d1b48b02a483b5341ba70b2660616d72d1aeb006518f3e261f7f651e48d00000000b500483045022100af88bce0ce4b1327c2b761423872be62cbb6a3efc1fcaa37f20799bd2ec32f95022007eca388424e7e7a053e50edfa5eb311e4a343437a195c9445503f2e62693c6a014c69522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553ae0000000001301b0f00000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88ac00000000")

	input := &adaptor.BindTxAndSignatureInput{
		Transaction: rawTx,
		SignedTxs:   [][]byte{signenTx1, signenTx2},
		Extra:       []byte("522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553ae"),
	}
	output, err := BindTxAndSignature(input, NETID_TEST)
	//not complete
	if err != nil {
		fmt.Println(err)
	} else {
		resultJSON, _ := json.Marshal(output)
		fmt.Println(string(resultJSON))
		fmt.Printf("%x\n", output.SignedTx)
	}
}

//func TestMultisignOneByOne(t *testing.T) {
//	sigTransaction, partSigedScript, complete := MultisignOneByOne("101d482b60cd3f74a61ce265d62e383456b9c21c84477931d207ea8f503d84cc", 0,
//		100000000, 100000, "mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5",
//		"522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553ae",
//		"", "cUakDAWEeNeXTo3B93WBs9HRMfaFDegXcbEGooLz8BSxRBfmpYcX", NETID_TEST)
//	fmt.Println(sigTransaction)
//	fmt.Println(partSigedScript)
//
//	sigTransaction, lastSigedScript, complete := MultisignOneByOne("101d482b60cd3f74a61ce265d62e383456b9c21c84477931d207ea8f503d84cc", 0,
//		100000000, 100000, "mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5",
//		"522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553ae",
//		partSigedScript, "cQJB6w8SxVNoprVwp2xyxUFxvExMbpR2qj3banXYYXmhtTc1WxC8", NETID_TEST)
//	fmt.Println(sigTransaction)
//	fmt.Println(lastSigedScript)
//
//	sigTransactionTest := "0100000001cc843d508fea07d2317947841cc2b95634382ed665e21ca6743fcd602b481d1000000000fc00473044022016f92f6342c945132779c12009cedc67f8eff461fdb7327153aa28aa42b1a37d02205fb1d6eefe06fdc1d003ed17f65d367178202b32dd619d7f2cc98cbf75ae84b50147304402200b552fc38bdefbd85c069df29ceed32e3d3199f2ef7dbbdb0c9a72de6612b9f2022041018a9ec1a9fae45b1902ba92a99f161c545598fb188835754f4c809817aaa4014c69522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553aeffffffff01605af405000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88ac00000000"
//	lastSigedScriptTest := "00473044022016f92f6342c945132779c12009cedc67f8eff461fdb7327153aa28aa42b1a37d02205fb1d6eefe06fdc1d003ed17f65d367178202b32dd619d7f2cc98cbf75ae84b50147304402200b552fc38bdefbd85c069df29ceed32e3d3199f2ef7dbbdb0c9a72de6612b9f2022041018a9ec1a9fae45b1902ba92a99f161c545598fb188835754f4c809817aaa4014c69522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553ae"
//
//	if !strings.Contains(sigTransaction, sigTransactionTest) {
//		t.Errorf("unexpected address - got: %x, "+
//			"want: %x", sigTransaction, sigTransactionTest)
//	}
//	if !strings.Contains(lastSigedScript, lastSigedScriptTest) {
//		t.Errorf("unexpected address - got: %x, "+
//			"want: %x", lastSigedScript, lastSigedScriptTest)
//	}
//	if !complete {
//		t.Errorf("complete - got: false, want: true")
//	}
//}
