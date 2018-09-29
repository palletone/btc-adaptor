/*
   This file is part of go-palletone.
   go-palletone is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-palletone is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
*/
/*
 * @author PalletOne core developers <dev@pallet.one>
 * @date 2018
 */
package adaptorbtc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"

	"github.com/palletone/adaptor"
)

func SignTransaction(signTransactionParams *adaptor.SignTransactionParams, rpcParams *RPCParams, netID int) (string, error) {
	//	//convert params from json format
	//	var signTransactionParams SignTransactionParams
	//	err := json.Unmarshal([]byte(params), &signTransactionParams)
	//	if err != nil {
	//		return err.Error()
	//	}

	//check empty string
	if "" == signTransactionParams.TransactionHex {
		return "", errors.New("Params error : NO TransactionHex.")
	}

	//decode Transaction hexString to bytes
	rawTXBytes, err := hex.DecodeString(signTransactionParams.TransactionHex)
	if err != nil {
		return "", err
	}
	//deserialize to MsgTx
	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewReader(rawTXBytes))
	if err != nil {
		return "", err
	}

	//chainnet
	var realNet *chaincfg.Params
	if netID == NETID_MAIN {
		realNet = &chaincfg.MainNetParams
	} else {
		realNet = &chaincfg.TestNet3Params
	}

	//get private keys for sign
	var keys []string
	for _, key := range signTransactionParams.Privkeys {
		key = strings.TrimSpace(key) //Trim whitespace
		if len(key) == 0 {
			continue
		}
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return "", err
	}

	//sign the UTXO hash, must know RedeemHex which contains in RawTxInput
	var rawInputs []btcjson.RawTxInput
	for {
		if "" == signTransactionParams.RedeemHex {
			break
		}
		//decode redeem's hexString to bytes
		redeem, err := hex.DecodeString(signTransactionParams.RedeemHex)
		if err != nil {
			break
		}
		//get multisig payScript
		scriptAddr, err := btcutil.NewAddressScriptHash(redeem, realNet)
		scriptPkScript, err := txscript.PayToAddrScript(scriptAddr)
		//multisig transaction need redeem for sign
		for _, txinOne := range tx.TxIn {
			rawInput := btcjson.RawTxInput{
				txinOne.PreviousOutPoint.Hash.String(), //txid
				txinOne.PreviousOutPoint.Index,         //outindex
				hex.EncodeToString(scriptPkScript),     //multisig pay script
				signTransactionParams.RedeemHex}        //redeem
			rawInputs = append(rawInputs, rawInput)
		}
		break
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		return "", err
	}
	defer client.Shutdown()

	//if complete ruturn true
	resultTX, complete, err := client.SignRawTransaction3(&tx, rawInputs, keys)
	if err != nil {
		return "", err
	}

	//SerializeSize transaction to bytes
	bufTX := bytes.NewBuffer(make([]byte, 0, resultTX.SerializeSize()))
	if err := resultTX.Serialize(bufTX); err != nil {
		return "", err
	}

	//result for return
	var signTransactionResult adaptor.SignTransactionResult
	signTransactionResult.TransactionHex = hex.EncodeToString(bufTX.Bytes())
	signTransactionResult.Complete = complete

	jsonResult, err := json.Marshal(signTransactionResult)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

func SendTransaction(params string, rpcParams *RPCParams) string {
	//convert params from json format
	var sendTransactionParams adaptor.SendTransactionParams
	err := json.Unmarshal([]byte(params), &sendTransactionParams)
	if err != nil {
		return err.Error()
	}

	//check empty string
	if "" == sendTransactionParams.TransactionHex {
		return "Params error : NO TransactionHex."
	}

	//decode Transaction hexString to bytes
	rawTXBytes, err := hex.DecodeString(sendTransactionParams.TransactionHex)
	if err != nil {
		return err.Error()
	}
	//deserialize to MsgTx
	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewReader(rawTXBytes))
	if err != nil {
		return err.Error()
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		return err.Error()
	}
	defer client.Shutdown()

	//send to network
	hashTX, err := client.SendRawTransaction(&tx, false)
	if err != nil {
		return err.Error()
	}

	//result for return
	var sendTransactionResult adaptor.SendTransactionResult
	sendTransactionResult.TransactionHah = hashTX.String()

	jsonResult, err := json.Marshal(sendTransactionResult)
	if err != nil {
		return err.Error()
	}

	return string(jsonResult)
}

func SignTxSend(signTxSendParams *adaptor.SignTxSendParams, rpcParams *RPCParams, netID int) (string, error) {
	//check empty string
	if "" == signTxSendParams.TransactionHex {
		return "", errors.New("Params error : NO TransactionHex.")
	}

	//decode Transaction hexString to bytes
	rawTXBytes, err := hex.DecodeString(signTxSendParams.TransactionHex)
	if err != nil {
		return "", err
	}
	//deserialize to MsgTx
	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewReader(rawTXBytes))
	if err != nil {
		return "", err
	}

	//chainnet
	var realNet *chaincfg.Params
	if netID == NETID_MAIN {
		realNet = &chaincfg.MainNetParams
	} else {
		realNet = &chaincfg.TestNet3Params
	}

	//get private keys for sign
	var keys []string
	for _, key := range signTxSendParams.Privkeys {
		key = strings.TrimSpace(key) //Trim whitespace
		if len(key) == 0 {
			continue
		}
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return "", err
	}

	//sign the UTXO hash, must know RedeemHex which contains in RawTxInput
	var rawInputs []btcjson.RawTxInput
	for {
		if "" == signTxSendParams.RedeemHex {
			break
		}
		//decode redeem's hexString to bytes
		redeem, err := hex.DecodeString(signTxSendParams.RedeemHex)
		if err != nil {
			break
		}
		//get multisig payScript
		scriptAddr, err := btcutil.NewAddressScriptHash(redeem, realNet)
		scriptPkScript, err := txscript.PayToAddrScript(scriptAddr)
		//multisig transaction need redeem for sign
		for _, txinOne := range tx.TxIn {
			rawInput := btcjson.RawTxInput{
				txinOne.PreviousOutPoint.Hash.String(), //txid
				txinOne.PreviousOutPoint.Index,         //outindex
				hex.EncodeToString(scriptPkScript),     //multisig pay script
				signTxSendParams.RedeemHex}             //redeem
			rawInputs = append(rawInputs, rawInput)
		}
		break
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		return "", err
	}
	defer client.Shutdown()

	//if complete ruturn true
	resultTX, complete, err := client.SignRawTransaction3(&tx, rawInputs, keys)
	if err != nil {
		return "", err
	}

	//result for return
	var signTxSendResult adaptor.SignTxSendResult
	if complete {
		//send to network
		hashTX, err := client.SendRawTransaction(resultTX, false)
		if err != nil {
			return "", err
		}
		signTxSendResult.TransactionHah = hashTX.String()

		//SerializeSize transaction to bytes
		bufTX := bytes.NewBuffer(make([]byte, 0, resultTX.SerializeSize()))
		if err := resultTX.Serialize(bufTX); err != nil {
			return "", err
		}

		signTxSendResult.TransactionHex = hex.EncodeToString(bufTX.Bytes())
		signTxSendResult.Complete = complete
	}

	jsonResult, err := json.Marshal(signTxSendResult)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

//==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ===

type addressToKey struct {
	key        *btcec.PrivateKey
	compressed bool
}

//find the privatekey by address
func mkGetKey(keys map[string]addressToKey) txscript.KeyDB {
	if keys == nil {
		return txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
			return nil, false, errors.New("nope")
		})
	}
	return txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
		a2k, ok := keys[addr.EncodeAddress()]
		if !ok {
			return nil, false, errors.New("nope")
		}
		return a2k.key, a2k.compressed, nil
	})
}

//find the redeemhex by address
func mkGetScript(scripts map[string][]byte) txscript.ScriptDB {
	if scripts == nil {
		return txscript.ScriptClosure(func(addr btcutil.Address) ([]byte, error) {
			return nil, errors.New("nope")
		})
	}
	return txscript.ScriptClosure(func(addr btcutil.Address) ([]byte, error) {
		script, ok := scripts[addr.EncodeAddress()]
		if !ok {
			return nil, errors.New("nope")
		}
		return script, nil
	})
}

//if complete, ruturn nil
func checkScripts(tx *wire.MsgTx, idx int, inputAmt int64,
	sigScript, scriptPkScript []byte) error {
	vm, err := txscript.NewEngine(scriptPkScript, tx, idx,
		txscript.ScriptBip16|txscript.ScriptVerifyDERSignatures,
		nil, nil, inputAmt)
	if err != nil {
		return err
	}

	err = vm.Execute()
	if err != nil {
		return err
	}

	return nil
}

// one input, one output, signed one by one.
// if not complete, return signedTransaction partSigedScript and false,
// if complete, ruturn signedTransaction lastSigedScript and true.
func MultisignOneByOne(prevTxHash string, index uint,
	amount int64, fee int64, recvAddress string,
	redeem string, partSigedScript string,
	wifKey string, netID int) (signedTransaction, newSigedScript string, complete bool) {
	//chainnet
	var realNet *chaincfg.Params
	if netID == NETID_MAIN {
		realNet = &chaincfg.MainNetParams
	} else {
		realNet = &chaincfg.TestNet3Params
	}

	//
	hash, _ := chainhash.NewHashFromStr(prevTxHash)
	outPoint := wire.NewOutPoint(hash, uint32(index))
	//
	txIn := wire.NewTxIn(outPoint, nil, nil)
	inputs := []*wire.TxIn{txIn}

	//
	var recvAmount = amount - fee
	addr, err := btcutil.DecodeAddress(recvAddress, realNet)
	if err != nil {
		return "", "", false
	}
	pubkeyScript, _ := txscript.PayToAddrScript(addr)
	//
	outputs := []*wire.TxOut{}
	outputs = append(outputs, wire.NewTxOut(recvAmount, pubkeyScript))

	//
	tx := &wire.MsgTx{
		Version:  1,
		TxIn:     inputs,
		TxOut:    outputs,
		LockTime: 0,
	}

	//
	var sigOldBytes []byte
	if partSigedScript == "" {
		sigOldBytes = nil
	} else {
		sigOldBytes, err = hex.DecodeString(partSigedScript)
		if err != nil {
			return "", "", false
		}
	}

	//
	key, err := btcutil.DecodeWIF(wifKey)
	//
	pub, err := btcutil.NewAddressPubKey(key.SerializePubKey(), realNet)
	if err != nil {
		return "", "", false
	}

	//
	pkScript, err := hex.DecodeString(redeem)
	//
	scriptAddr, err := btcutil.NewAddressScriptHash(pkScript, realNet)
	if err != nil {
		return "", "", false
	}

	//
	scriptPkScript, err := txscript.PayToAddrScript(scriptAddr)
	if err != nil {
		return "", "", false
	}

	// Two part multisig, sign with one key then the other.
	// Sign with the other key and merge
	sigScript, err := txscript.SignTxOutput(realNet,
		tx, 0, scriptPkScript, txscript.SigHashAll,
		mkGetKey(map[string]addressToKey{
			pub.EncodeAddress(): {key.PrivKey, true},
		}), mkGetScript(map[string][]byte{
			scriptAddr.EncodeAddress(): pkScript,
		}), sigOldBytes)
	if err != nil {
		return "", "", false
	}
	tx.TxIn[0].SignatureScript = sigScript

	//
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", "", false
	}

	//
	err = checkScripts(tx, 0, amount, sigScript, scriptPkScript)
	if err != nil {
		complete = false
	} else {
		complete = true
	}

	return hex.EncodeToString(buf.Bytes()), hex.EncodeToString(sigScript), complete
}
