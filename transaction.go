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

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"

	"github.com/palletone/adaptor"
)

func RawTransactionGen(rawTransactionGenParams *adaptor.RawTransactionGenParams, rpcParams *RPCParams, netID int) (string, error) {
	//	//convert params from json format
	//	var rawTransactionGenParams RawTransactionGenParams
	//	err := json.Unmarshal([]byte(params), &rawTransactionGenParams)
	//	if err != nil {
	//		return "", err
	//	}

	//transaction inputs
	var inputs []btcjson.TransactionInput
	for _, inputOne := range rawTransactionGenParams.Inputs {
		input := btcjson.TransactionInput{inputOne.Txid, inputOne.Vout}
		inputs = append(inputs, input)
	}
	if len(inputs) == 0 {
		return "", errors.New("Params error : NO Input.")
	}

	//chainnet
	var realNet *chaincfg.Params
	if netID == NETID_MAIN {
		realNet = &chaincfg.MainNetParams
	} else {
		realNet = &chaincfg.TestNet3Params
	}

	//transaction outputs
	amounts := map[btcutil.Address]btcutil.Amount{}
	for _, outOne := range rawTransactionGenParams.Outputs {
		if len(outOne.Address) == 0 || outOne.Amount <= 0 {
			continue
		}
		addr, err := btcutil.DecodeAddress(outOne.Address, realNet)
		if err != nil {
			return "", err
		}
		amounts[addr] = btcutil.Amount(outOne.Amount * 1e8)
	}
	if len(amounts) == 0 {
		return "", errors.New("Params error : NO Output.")
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		return "", err
	}
	defer client.Shutdown()

	//only inputs and outputs, no redeem
	msgTx, err := client.CreateRawTransaction(inputs, amounts, &rawTransactionGenParams.Locktime)
	if err != nil {
		return "", err
	}
	//SerializeSize transaction to bytes
	buf := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))
	if err := msgTx.Serialize(buf); err != nil {
		return "", err
	}
	//result for return
	var rawTransactionGenResult adaptor.RawTransactionGenResult
	rawTransactionGenResult.Rawtx = hex.EncodeToString(buf.Bytes())

	jsonResult, err := json.Marshal(rawTransactionGenResult)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

func DecodeRawTransaction(decodeRawTransactionParams *adaptor.DecodeRawTransactionParams, rpcParams *RPCParams) (string, error) {
	//	//convert params from json format
	//	var decodeRawTransactionParams DecodeRawTransactionParams
	//	err := json.Unmarshal([]byte(params), &decodeRawTransactionParams)
	//	if err != nil {
	//		return "", err
	//	}
	if "" == decodeRawTransactionParams.Rawtx {
		return "", errors.New("Params error : NO Rawtx.")
	}

	//covert rawtransaction hexString to bytes
	rawTXBytes, err := hex.DecodeString(decodeRawTransactionParams.Rawtx)
	if err != nil {
		return "", err
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		return "", err
	}
	defer client.Shutdown()

	//rpc DecodeRawTransaction
	resultTxRaw, err := client.DecodeRawTransaction(rawTXBytes)
	if err != nil {
		return "", err
	}

	//result for return
	var result adaptor.DecodeRawTransactionResult
	result.Locktime = resultTxRaw.LockTime
	for i, _ := range resultTxRaw.Vin {
		result.Inputs = append(result.Inputs, adaptor.Input{resultTxRaw.Vin[i].Txid, resultTxRaw.Vin[i].Vout})
	}
	for i, _ := range resultTxRaw.Vout {
		result.Outputs = append(result.Outputs, adaptor.Output{resultTxRaw.Vout[i].ScriptPubKey.Addresses[0], resultTxRaw.Vout[i].Value})
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

func GetTransactionByHash(getTransactionByHashParams *adaptor.GetTransactionByHashParams, rpcParams *RPCParams) (string, error) {
	//covert TxHash
	hash, err := chainhash.NewHashFromStr(getTransactionByHashParams.TxHash)
	if err != nil {
		return "", err
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		return "", err
	}
	defer client.Shutdown()

	//rpc GetRawTransactionVerbose
	txResult, err := client.GetRawTransactionVerbose(hash)
	if err != nil {
		return "", err
	}

	//result for return
	var getTransactionByHashResult adaptor.GetTransactionByHashResult
	for _, out := range txResult.Vout {
		getTransactionByHashResult.Outputs = append(getTransactionByHashResult.Outputs,
			adaptor.OutputIndex{out.N, out.ScriptPubKey.Addresses[0], out.Value})
	}
	for _, in := range txResult.Vin {
		getTransactionByHashResult.Inputs = append(getTransactionByHashResult.Inputs,
			adaptor.Input{in.Txid, in.Vout})
	}
	getTransactionByHashResult.Txid = txResult.Txid
	getTransactionByHashResult.Confirms = txResult.Confirmations

	jsonResult, err := json.Marshal(getTransactionByHashResult)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}
