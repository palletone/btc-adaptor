package btcadaptor

import (
	"encoding/hex"
	"fmt"
	"testing"
	//"github.com/btcsuite/btcutil"
	//"github.com/palletone/adaptor"
	"github.com/palletone/adaptor"
)

func TestNewPrivateKey(t *testing.T) {
	//wifKey, _ := btcutil.DecodeWIF("cUakDAWEeNeXTo3B93WBs9HRMfaFDegXcbEGooLz8BSxRBfmpYcX")
	//fmt.Printf("%x\n", wifKey.PrivKey.Serialize())
	key, _ := NewPrivateKey(NETID_TEST)
	fmt.Printf("%x\n", key)
}

func TestGetPublicKey(t *testing.T) {
	keyHex := "d0e26e9189b9f047036ed21294c8f36d41df6b51852fc932595d849d727223d0"
	testPubkey := "029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef5"

	key, _ := hex.DecodeString(keyHex)
	pubkey, _ := GetPublicKey(key, NETID_TEST)
	pubkeyHex := hex.EncodeToString(pubkey)
	fmt.Println(pubkeyHex)
	if testPubkey != pubkeyHex {
		t.Errorf("unexpected pubkey bytes - got: %s, "+
			"want: %s", pubkeyHex, testPubkey)
	}
}

func TestPubKeyToAddress(t *testing.T) {
	pubKeyHex := "029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef5"
	testAddr := "mxprH5bkXtn9tTTAxdQGPXrvruCUvsBNKt"

	pubKey, _ := hex.DecodeString(pubKeyHex)

	addr, err := PubKeyToAddress(pubKey, NETID_TEST)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(addr)
	if testAddr != addr {
		t.Errorf("unexpected address - got: %v, "+
			"want: %v", addr, testAddr)
	}
}

func TestCreateMultiSigAddress(t *testing.T) {
	pubkeyAliceHex := "03940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d5410514"
	pubkeyBobHex := "020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb"
	pubkeyPalletHex := "029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef5"

	pubkeyAlice, _ := hex.DecodeString(pubkeyAliceHex)
	pubkeyBob, _ := hex.DecodeString(pubkeyBobHex)
	pubkeyPallet, _ := hex.DecodeString(pubkeyPalletHex)

	var input adaptor.CreateMultiSigAddressInput
	input.Keys = append(input.Keys, pubkeyAlice)
	input.Keys = append(input.Keys, pubkeyBob)
	input.Keys = append(input.Keys, pubkeyPallet)
	input.SignCount = 2

	resultMain, err := CreateMultiSigAddress(&input, NETID_MAIN)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	testRedeem := "522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553ae"
	multiAddrMain := "3DBKFERmXBuy8btJ2x778Dyz8BQw7obDNn"
	redeemMain := hex.EncodeToString(resultMain.Extra)
	fmt.Println(redeemMain)
	if redeemMain != testRedeem {
		t.Errorf("unexpected address - got: %x, "+
			"want: %x", resultMain, testRedeem)
	}
	fmt.Println(resultMain.Address)
	if resultMain.Address != multiAddrMain {
		t.Errorf("unexpected address - got: %x, "+
			"want: %x", resultMain, multiAddrMain)
	}

	resultTest, err := CreateMultiSigAddress(&input, NETID_TEST)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	redeemTest := hex.EncodeToString(resultMain.Extra)
	fmt.Println(redeemTest)
	if redeemTest != testRedeem {
		t.Errorf("unexpected address - got: %s, "+
			"want: %s", redeemTest, testRedeem)
	}
	multiAddrTest := "2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA"
	fmt.Println(resultTest.Address)
	if resultTest.Address != multiAddrTest {
		t.Errorf("unexpected address - got: %s, "+
			"want: %s", resultTest, multiAddrTest)
	}
}
