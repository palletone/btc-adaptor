package adaptorbtc

import (
	"fmt"
	"testing"

	"github.com/palletone/adaptor"
)

func TestNewPrivateKey(t *testing.T) {
	key := NewPrivateKey(NETID_MAIN)
	fmt.Println(key)
	keyTest := NewPrivateKey(NETID_TEST)
	fmt.Println(keyTest)
}

func TestGetPublicKey(t *testing.T) {
	key := "cUakDAWEeNeXTo3B93WBs9HRMfaFDegXcbEGooLz8BSxRBfmpYcX"
	testPubkey := "029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef5"
	pubkey := GetPublicKey(key, NETID_TEST)
	if testPubkey != pubkey {
		t.Errorf("unexpected pubkey bytes - got: %v, "+
			"want: %v", pubkey, testPubkey)
	}
}

func TestGetAddress(t *testing.T) {
	key := "cUakDAWEeNeXTo3B93WBs9HRMfaFDegXcbEGooLz8BSxRBfmpYcX"
	testAddr := "mxprH5bkXtn9tTTAxdQGPXrvruCUvsBNKt"
	addr := GetAddress(key, NETID_TEST)
	if testAddr != addr {
		t.Errorf("unexpected address - got: %v, "+
			"want: %v", addr, testAddr)
	}
}

func TestGetAddressByPubkey(t *testing.T) {
	pubkey := "029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef5"
	testAddr := "mxprH5bkXtn9tTTAxdQGPXrvruCUvsBNKt"
	addr, err := GetAddressByPubkey(pubkey, NETID_TEST)
	if err != nil {
		fmt.Println(err.Error())
	}
	if testAddr != addr {
		t.Errorf("unexpected address - got: %v, "+
			"want: %v", addr, testAddr)
	}
}

func TestCreateMultiSigAddress(t *testing.T) {
	//	//
	//	params := `{
	//    "publicKeys": ["03940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d5410514","020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb","029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef5"],
	//    "n": 3,
	//    "m": 2
	//  	}
	//	`
	pubkeyAlice := "03940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d5410514"
	pubkeyBob := "020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb"
	pubkeyPallet := "029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef5"
	var createMultiSigParams adaptor.CreateMultiSigParams
	createMultiSigParams.PublicKeys = append(createMultiSigParams.PublicKeys, pubkeyAlice)
	createMultiSigParams.PublicKeys = append(createMultiSigParams.PublicKeys, pubkeyBob)
	createMultiSigParams.PublicKeys = append(createMultiSigParams.PublicKeys, pubkeyPallet)
	createMultiSigParams.M = 2
	createMultiSigParams.N = 3

	//	resultMain := CreateMultiSigAddress(params, NETID_MAIN)
	resultMain, err := CreateMultiSigAddress(&createMultiSigParams, NETID_MAIN)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	testRedeem := "522103940ab29fbf214da2d8ec99c47db63879957311bd90d2f1c635828604d541051421020106ca23b4f28dbc83838ee4745accf90e5621fe70df5b1ee8f7e1b3b41b64cb21029d80ff37838e4989a6aa26af41149d4f671976329e9ddb9b78fdea9814ae6ef553ae"
	multiAddrMain := "3DBKFERmXBuy8btJ2x778Dyz8BQw7obDNn"
	if resultMain.RedeemScript != testRedeem {
		t.Errorf("unexpected address - got: %x, "+
			"want: %x", resultMain, testRedeem)
	}
	if resultMain.P2ShAddress != multiAddrMain {
		t.Errorf("unexpected address - got: %x, "+
			"want: %x", resultMain, multiAddrMain)
	}

	//
	//	resultTest := CreateMultiSigAddress(params, NETID_TEST)
	resultTest, err := CreateMultiSigAddress(&createMultiSigParams, NETID_TEST)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	multiAddrTest := "2N4jXJyMo8eRKLPWqi5iykAyFLXd6szehwA"
	if resultTest.RedeemScript != testRedeem {
		t.Errorf("unexpected address - got: %s, "+
			"want: %s", resultTest, testRedeem)
	}
	if resultTest.P2ShAddress != multiAddrTest {
		t.Errorf("unexpected address - got: %s, "+
			"want: %s", resultTest, multiAddrTest)
	}
}
