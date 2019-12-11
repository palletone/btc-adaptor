package btcadaptor

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/palletone/adaptor"
)

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
		fmt.Printf(string(output.Signature)) //H7GJvdYPHxjeJWkAj5negCzcoW2JjvcY93UF5wMqW0D5etuRkK0qDr7KflmcZrIkflSzb4I6X8ZqSZHEXYG3jSA=
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
	addr := "mgtT62nq65DsPPAzPp6KhsWoHjNQUR9Bu5"
	keyHex := "ac18d6ffa4e006ee5c297e14962d213030910a045dccc9d393686eb1613a1477"
	key, _ := hex.DecodeString(keyHex)

	tx, _ := hex.DecodeString("0100000001cc843d508fea07d2317947841cc2b95634382ed665e21ca6743fcd602b481d100100000000000000000240420f000000000017a9147e037d8b8093a7cf3a6ec83aa8c852761a5d0cce8740420f00000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88ac00000000")

	input := &adaptor.SignTransactionInput{PrivateKey: key, Transaction: tx, Extra: []byte(addr)}

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

	tx, _ := hex.DecodeString("0100000001cc843d508fea07d2317947841cc2b95634382ed665e21ca6743fcd602b481d10010000006b483045022100aaeac1a6f3cf3386fcc27d726f385b0d0d5543abd5e69d7f81cfd52e72ccdc7602205eef66553eea6c35e19e48eb31baf24a9f77495c01c14071a225a346a165f5600121020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb000000000240420f000000000017a9147e037d8b8093a7cf3a6ec83aa8c852761a5d0cce8740420f00000000001976a9140f08e55bcfc207632d2dcfc3d4db4b6d8d91b22e88ac00000000")

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
//
//func TestMergeTransaction(t *testing.T) {
//	//raw tx
//	paramsMerge := &adaptor.MergeTransactionParams{
//		UserTransactionHex: "010000000236045404e65bd741109db92227ca0dc9274ef717a6612c96cd77b24a17d1bcd70000000000ffffffff7c1f7d5407b41abf29d41cf6f122ef2d40f76d956900d2c89314970951ef5b940000000000ffffffff014431d309000000001976a914bddc9a62e9b7c3cfdbe1c817520e24e32c339f3288ac00000000",
//		InputRedeemIndex:   []int{0},
//		RedeemHex:          []string{"522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553ae"},
//	}
//	//apped the tx signed by alice, for merge
//	paramsMerge.MergeTransactionHexs = append(paramsMerge.MergeTransactionHexs,
//		"010000000236045404e65bd741109db92227ca0dc9274ef717a6612c96cd77b24a17d1bcd700000000b400473044022024e6a6ca006f25ccd3ebf5dadf21397a6d7266536cd336061cd17cff189d95e402205af143f6726d75ac77bc8c80edcb6c56579053d2aa31601b23bc8da41385dd86014c69522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553aeffffffff7c1f7d5407b41abf29d41cf6f122ef2d40f76d956900d2c89314970951ef5b9400000000b40047304402206a1d7a2ae07840957bee708b6d3e1fbe7858760ac378b1e21209b348c1e2a5c402204255cd4cd4e5b5805d44bbebe7464aa021377dca5fc6bf4a5632eb2d8bc9f9e4014c69522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553aeffffffff014431d309000000001976a914bddc9a62e9b7c3cfdbe1c817520e24e32c339f3288ac00000000")
//	result, err := MergeTransaction(paramsMerge, NETID_TEST)
//	//not complete
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(result)
//	}
//
//	//apped the tx signed by bob, for merge
//	paramsMerge.MergeTransactionHexs = append(paramsMerge.MergeTransactionHexs,
//		"010000000236045404e65bd741109db92227ca0dc9274ef717a6612c96cd77b24a17d1bcd700000000b5004830450221009b4f02e07cab7f3c68125499d466d22086c02c3a921b7bebe32434106e6cf7f1022016894e5e1639420210d2b5ecff93c48f698e4bf4b189ed5cb2d275dbe91dcc4c014c69522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553aeffffffff7c1f7d5407b41abf29d41cf6f122ef2d40f76d956900d2c89314970951ef5b9400000000b500483045022100ee6f94e09447b2dd66cd10a5b6c6ba4e6b3215d5a50d042bacac3c895c369705022079b97833c85c91cd15efe99d945f98be75143c46f7b2fd9ff95d5daeb1e4d9f7014c69522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553aeffffffff014431d309000000001976a914bddc9a62e9b7c3cfdbe1c817520e24e32c339f3288ac00000000")
//	result2, err := MergeTransaction(paramsMerge, NETID_TEST)
//	//complete
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(result2)
//	}
//	return
//}
//

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
