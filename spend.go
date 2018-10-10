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
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"

	//"github.com/btcsuite/btcd/rpcclient/chain"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"

	"github.com/palletone/adaptor"
)

// decodeHexStr decodes the hex encoding of a string, possibly prepending a
// leading '0' character if there is an odd number of bytes in the hex string.
// This is to prevent an error for an invalid hex string when using an odd
// number of bytes when calling hex.Decode.
func decodeHexStr(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDecodeHexString,
			Message: "Hex string decode failed: " + err.Error(),
		}
	}
	return decoded, nil
}

// signRawTransaction handles the signrawtransaction command.
func signRawTransactionCmd(icmd interface{}, chainParams *chaincfg.Params,
	chainClient *rpcclient.Client) (interface{}, error) {
	cmd := icmd.(*btcjson.SignRawTransactionCmd)

	serializedTx, err := decodeHexStr(cmd.RawTx)
	if err != nil {
		return nil, err
	}
	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewBuffer(serializedTx))
	if err != nil {
		return nil, errors.New("TX decode failed")
	}

	var hashType txscript.SigHashType
	switch *cmd.Flags {
	case "ALL":
		hashType = txscript.SigHashAll
	case "NONE":
		hashType = txscript.SigHashNone
	case "SINGLE":
		hashType = txscript.SigHashSingle
	case "ALL|ANYONECANPAY":
		hashType = txscript.SigHashAll | txscript.SigHashAnyOneCanPay
	case "NONE|ANYONECANPAY":
		hashType = txscript.SigHashNone | txscript.SigHashAnyOneCanPay
	case "SINGLE|ANYONECANPAY":
		hashType = txscript.SigHashSingle | txscript.SigHashAnyOneCanPay
	default:
		return nil, errors.New("Invalid sighash parameter")
	}

	// TODO: really we probably should look these up with btcd anyway to
	// make sure that they match the blockchain if present.
	inputs := make(map[wire.OutPoint][]byte)
	scripts := make(map[string][]byte)
	var cmdInputs []btcjson.RawTxInput
	if cmd.Inputs != nil {
		cmdInputs = *cmd.Inputs
	}
	for _, rti := range cmdInputs {
		inputHash, err := chainhash.NewHashFromStr(rti.Txid)
		if err != nil {
			return nil, err
		}

		script, err := decodeHexStr(rti.ScriptPubKey)
		if err != nil {
			return nil, err
		}

		// redeemScript is only actually used iff the user provided
		// private keys. In which case, it is used to get the scripts
		// for signing. If the user did not provide keys then we always
		// get scripts from the wallet.
		// Empty strings are ok for this one and hex.DecodeString will
		// DTRT.
		if cmd.PrivKeys != nil && len(*cmd.PrivKeys) != 0 {
			redeemScript, err := decodeHexStr(rti.RedeemScript)
			if err != nil {
				return nil, err
			}

			addr, err := btcutil.NewAddressScriptHash(redeemScript,
				chainParams)
			if err != nil {
				return nil, err
			}
			scripts[addr.String()] = redeemScript
		}
		inputs[wire.OutPoint{
			Hash:  *inputHash,
			Index: rti.Vout,
		}] = script
	}

	// Now we go and look for any inputs that we were not provided by
	// querying btcd with getrawtransaction. We queue up a bunch of async
	// requests and will wait for replies after we have checked the rest of
	// the arguments.
	requested := make(map[wire.OutPoint]rpcclient.FutureGetTxOutResult)
	for _, txIn := range tx.TxIn {
		// Did we get this outpoint from the arguments?
		if _, ok := inputs[txIn.PreviousOutPoint]; ok {
			continue
		}

		// Asynchronously request the output script.
		requested[txIn.PreviousOutPoint] = chainClient.GetTxOutAsync(
			&txIn.PreviousOutPoint.Hash, txIn.PreviousOutPoint.Index,
			true)
	}

	// Parse list of private keys, if present. If there are any keys here
	// they are the keys that we may use for signing. If empty we will
	// use any keys known to us already.
	var keys map[string]*btcutil.WIF
	if cmd.PrivKeys != nil {
		keys = make(map[string]*btcutil.WIF)

		for _, key := range *cmd.PrivKeys {
			wif, err := btcutil.DecodeWIF(key)
			if err != nil {
				return nil, err
			}

			if !wif.IsForNet(chainParams) {
				s := "key network doesn't match wallet's"
				return nil, errors.New(s)
			}

			addr, err := btcutil.NewAddressPubKey(wif.SerializePubKey(),
				chainParams)
			if err != nil {
				return nil, err
			}
			keys[addr.EncodeAddress()] = wif
		}
	}

	// We have checked the rest of the args. now we can collect the async
	// txs. TODO: If we don't mind the possibility of wasting work we could
	// move waiting to the following loop and be slightly more asynchronous.
	for outPoint, resp := range requested {
		result, err := resp.Receive()
		if err != nil {
			return nil, err
		}
		if result == nil {
			s := fmt.Sprintf("the output %s - %d has been spent already", outPoint.Hash, outPoint.Index)
			return nil, errors.New(s)
		}
		script, err := hex.DecodeString(result.ScriptPubKey.Hex)
		if err != nil {
			return nil, err
		}
		inputs[outPoint] = script
	}

	// All args collected. Now we can sign all the inputs that we can.
	// `complete' denotes that we successfully signed all outputs and that
	// all scripts will run to completion. This is returned as part of the
	// reply.
	signErrs, err := SignTransactionReal(&tx, hashType, inputs, keys, scripts, chainParams)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Grow(tx.SerializeSize())

	// All returned errors (not OOM, which panics) encounted during
	// bytes.Buffer writes are unexpected.
	if err = tx.Serialize(&buf); err != nil {
		panic(err)
	}

	signErrors := make([]btcjson.SignRawTransactionError, 0, len(signErrs))
	for _, e := range signErrs {
		input := tx.TxIn[e.InputIndex]
		signErrors = append(signErrors, btcjson.SignRawTransactionError{
			TxID:      input.PreviousOutPoint.Hash.String(),
			Vout:      input.PreviousOutPoint.Index,
			ScriptSig: hex.EncodeToString(input.SignatureScript),
			Sequence:  input.Sequence,
			Error:     e.Error.Error(),
		})
	}

	return btcjson.SignRawTransactionResult{
		Hex:      hex.EncodeToString(buf.Bytes()),
		Complete: len(signErrors) == 0,
		Errors:   signErrors,
	}, nil
}

// SignatureError records the underlying error when validating a transaction
// input signature.
type SignatureError struct {
	InputIndex uint32
	Error      error
}

func SignTransactionReal(tx *wire.MsgTx, hashType txscript.SigHashType,
	additionalPrevScripts map[wire.OutPoint][]byte,
	additionalKeysByAddress map[string]*btcutil.WIF,
	p2shRedeemScriptsByAddress map[string][]byte, chainParams *chaincfg.Params) ([]SignatureError, error) {

	signErrors := []SignatureError{}
	//var signErrors []SignatureErroerr := walletdb.View(w.db, func(dbtx walletdb.ReadTx) error {
	//addrmgrNs := dbtx.ReadBucket(waddrmgrNamespaceKey)
	//txmgrNs := dbtx.ReadBucZXL/ket(wtxmgrNamespaceKey)
	var err error
	for i, txIn := range tx.TxIn {
		prevOutScript, ok := additionalPrevScripts[txIn.PreviousOutPoint]
		if !ok {
			/*prevHash := &txIn.PreviousOutPoint.Hash
			prevIndex := txIn.PreviousOutPoint.Index
			txDetails, err := w.TxStore.TxDetails(txmgrNs, prevHash)
			if err != nil {
				return fmt.Errorf("cannot query previous transaction "+
					"details for %v: %v", txIn.PreviousOutPoint, err)
			}
			if txDetails == nil {
				return fmt.Errorf("%v not found",
					txIn.PreviousOutPoint)
			}
			prevOutScript = txDetails.MsgTx.TxOut[prevIndex].PkScript*/
		}

		// Set up our callbacks that we pass to txscript so it can
		// look up the appropriate keys and scripts by address.
		getKey := txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
			if len(additionalKeysByAddress) != 0 {
				addrStr := addr.EncodeAddress()
				wif, ok := additionalKeysByAddress[addrStr]
				if !ok {
					return nil, false, errors.New("no key for address")
				}
				return wif.PrivKey, wif.CompressPubKey, nil
			}
			return nil, false, errors.New("no key for address")
		})
		getScript := txscript.ScriptClosure(func(addr btcutil.Address) ([]byte, error) {
			// If keys were provided then we can only use the
			// redeem scripts provided with our inputs, too.
			if len(additionalKeysByAddress) != 0 {
				addrStr := addr.EncodeAddress()
				script, ok := p2shRedeemScriptsByAddress[addrStr]
				if !ok {
					return nil, errors.New("no script for address")
				}
				return script, nil
			}
			return nil, errors.New("no script for address")
		})

		// SigHashSingle inputs can only be signed if there's a
		// corresponding output. However this could be already signed,
		// so we always verify the output.
		if (hashType&txscript.SigHashSingle) !=
			txscript.SigHashSingle || i < len(tx.TxOut) {

			script, err := txscript.SignTxOutput(chainParams,
				tx, i, prevOutScript, hashType, getKey,
				getScript, txIn.SignatureScript)
			// Failure to sign isn't an error, it just means that
			// the tx isn't complete.
			if err != nil {
				signErrors = append(signErrors, SignatureError{
					InputIndex: uint32(i),
					Error:      err,
				})
				continue
			}
			txIn.SignatureScript = script
		}

		// Either it was already signed or we just signed it.
		// Find out if it is completely satisfied or still needs more.
		vm, err := txscript.NewEngine(prevOutScript, tx, i,
			txscript.StandardVerifyFlags, nil, nil, 0)
		if err == nil {
			err = vm.Execute()
		}
		if err != nil {
			signErrors = append(signErrors, SignatureError{
				InputIndex: uint32(i),
				Error:      err,
			})
		}
	}
	//return nil
	//})
	return signErrors, err
}

func SignTransaction(signTransactionParams *adaptor.SignTransactionParams, rpcParams *RPCParams, netID int) (string, error) {
	//check empty string
	if "" == signTransactionParams.TransactionHex {
		return "", errors.New("Params error : NO TransactionHex.")
	}

	//chainnet
	var realNet *chaincfg.Params
	if netID == NETID_MAIN {
		realNet = &chaincfg.MainNetParams
	} else {
		realNet = &chaincfg.TestNet3Params
	}

	var err error
	//sign the UTXO hash, must know RedeemHex which contains in RawTxInput
	var rawInputs []btcjson.RawTxInput
	for {
		if "" == signTransactionParams.RedeemHex {
			break
		}

		//decode Transaction hexString to bytes
		rawTXBytes, err := hex.DecodeString(signTransactionParams.TransactionHex)
		if err != nil {
			break
		}
		//deserialize to MsgTx
		var tx wire.MsgTx
		err = tx.Deserialize(bytes.NewReader(rawTXBytes))
		if err != nil {
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
	if err != nil {
		return "", err
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		return "", err
	}
	defer client.Shutdown()

	//
	var cmd btcjson.SignRawTransactionCmd
	cmd.RawTx = signTransactionParams.TransactionHex
	cmd.Inputs = &rawInputs
	cmd.PrivKeys = &signTransactionParams.Privkeys
	flags := "ALL"
	cmd.Flags = &flags

	//if complete ruturn true
	result, err := signRawTransactionCmd(&cmd, realNet, client)
	if err != nil {
		return "", err
	}

	//result for return
	signRawResult := result.(btcjson.SignRawTransactionResult)
	var signTransactionResult adaptor.SignTransactionResult
	signTransactionResult.TransactionHex = signRawResult.Hex
	signTransactionResult.Complete = signRawResult.Complete

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

	//chainnet
	var realNet *chaincfg.Params
	if netID == NETID_MAIN {
		realNet = &chaincfg.MainNetParams
	} else {
		realNet = &chaincfg.TestNet3Params
	}

	var err error
	//sign the UTXO hash, must know RedeemHex which contains in RawTxInput
	var rawInputs []btcjson.RawTxInput
	for {
		if "" == signTxSendParams.RedeemHex {
			break
		}

		//decode Transaction hexString to bytes
		rawTXBytes, err := hex.DecodeString(signTxSendParams.TransactionHex)
		if err != nil {
			break
		}
		//deserialize to MsgTx
		var tx wire.MsgTx
		err = tx.Deserialize(bytes.NewReader(rawTXBytes))
		if err != nil {
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
				signTxSendParams.RedeemHex}        //redeem
			rawInputs = append(rawInputs, rawInput)
		}

		break
	}
	if err != nil {
		return "", err
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		return "", err
	}
	defer client.Shutdown()

	//
	var cmd btcjson.SignRawTransactionCmd
	cmd.RawTx = signTxSendParams.TransactionHex
	cmd.Inputs = &rawInputs
	cmd.PrivKeys = &signTxSendParams.Privkeys
	flags := "ALL"
	cmd.Flags = &flags

	//if complete ruturn true
	result, err := signRawTransactionCmd(&cmd, realNet, client)
	if err != nil {
		return "", err
	}

	//result for return
	signRawResult := result.(btcjson.SignRawTransactionResult)
	var signTxSendResult adaptor.SignTxSendResult
	if signRawResult.Complete {
		//decode Transaction hexString to bytes
		rawTXBytes, err := hex.DecodeString(signRawResult.Hex)
		if err != nil {
			return "", err
		}
		//deserialize to MsgTx
		var resultTX wire.MsgTx
		err = resultTX.Deserialize(bytes.NewReader(rawTXBytes))
		if err != nil {
			return "", err
		}

		//send to network
		hashTX, err := client.SendRawTransaction(&resultTX, false)
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
		signTxSendResult.Complete = true
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
