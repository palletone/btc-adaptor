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

	"github.com/btcsuite/btcd/txscript"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
)

//==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ===
type RawTransactionGenParams struct {
	Inputs   []Input  `json:"inputs"`
	Outputs  []Output `json:"outputs"`
	Locktime int64    `json:"locktime"`
}
type RawTransactionGenResult struct {
	Rawtx string `json:"rawtx"`
}

func RawTransactionGen(rawTransactionGenParams *RawTransactionGenParams, rpcParams *RPCParams, netID int) (string, error) {
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
	var rawTransactionGenResult RawTransactionGenResult
	rawTransactionGenResult.Rawtx = hex.EncodeToString(buf.Bytes())

	jsonResult, err := json.Marshal(rawTransactionGenResult)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

type DecodeRawTransactionParams struct {
	Rawtx string `json:"rawtx"`
}

type Input struct {
	Txid string `json:"txid"`
	Vout uint32 `json:"vout"`
}
type Output struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}
type DecodeRawTransactionResult struct {
	Inputs   []Input  `json:"inputs"`
	Outputs  []Output `json:"outputs"`
	Locktime uint32   `json:"locktime"`
}

func DecodeRawTransaction(decodeRawTransactionParams *DecodeRawTransactionParams, rpcParams *RPCParams) (string, error) {
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
	var result DecodeRawTransactionResult
	result.Locktime = resultTxRaw.LockTime
	for i, _ := range resultTxRaw.Vin {
		result.Inputs = append(result.Inputs, Input{resultTxRaw.Vin[i].Txid, resultTxRaw.Vin[i].Vout})
	}
	for i, _ := range resultTxRaw.Vout {
		result.Outputs = append(result.Outputs, Output{resultTxRaw.Vout[i].ScriptPubKey.Addresses[0], resultTxRaw.Vout[i].Value})
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

type GetTransactionByHashParams struct {
	TxHash string `json:"txhash"`
}

type GetTransactionByHashResult struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []output `json:"outputs"`
}

func GetTransactionByHash(getTransactionByHashParams *GetTransactionByHashParams, rpcParams *RPCParams) (string, error) {
	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		return "", err
	}
	defer client.Shutdown()

	//
	hash, err := chainhash.NewHashFromStr(getTransactionByHashParams.TxHash)
	if err != nil {
		return "", err
	}
	//
	txResult, err := client.GetRawTransaction(hash)
	if err != nil {
		return "", err
	}

	msgTx := txResult.MsgTx()
	fmt.Println(msgTx)

	//result for return
	var getTransactionByHashResult GetTransactionByHashResult
	for i := range msgTx.TxOut {
		_, addrs, _, err := txscript.ExtractPkScriptAddrs(
			msgTx.TxOut[i].PkScript, &chaincfg.TestNet3Params)
		if err != nil {
			return "", err
		}
		getTransactionByHashResult.Outputs = append(getTransactionByHashResult.Outputs,
			output{uint32(i), addrs[0].String(), msgTx.TxOut[i].Value})
	}
	for i := range msgTx.TxIn {
		getTransactionByHashResult.Inputs = append(getTransactionByHashResult.Inputs,
			Input{msgTx.TxIn[i].PreviousOutPoint.Hash.String(), msgTx.TxIn[i].PreviousOutPoint.Index})

	}

	jsonResult, err := json.Marshal(getTransactionByHashResult)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}
