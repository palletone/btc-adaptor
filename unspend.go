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
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/big"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"

	"github.com/palletone/adaptor"
)

//==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ===

//var GHomeDir = btcutil.AppDataDir("btcwallet", false)
var GHomeDir = btcutil.AppDataDir("btcd", false)
var GCertPath = filepath.Join(GHomeDir, "rpc.cert")

func GetClient(rpcParams *RPCParams) (*rpcclient.Client, error) {
	//read cert from file
	var connCfg *rpcclient.ConnConfig
	if rpcParams.CertPath == "" {
		rpcParams.CertPath = GCertPath
	}
	if rpcParams.CertPath != "" {
		certs, err := ioutil.ReadFile(rpcParams.CertPath)
		if err != nil {
			return nil, err
		}

		// Connect to local bitcoin core RPC server using HTTP POST mode.
		connCfg = &rpcclient.ConnConfig{
			Host:         rpcParams.Host,
			Endpoint:     "ws",
			User:         rpcParams.RPCUser,
			Pass:         rpcParams.RPCPasswd,
			HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
			//DisableTLS:   true,  // Bitcoin core does not provide TLS by default
			Certificates: certs, // btcwallet provide TLS by default
		}
	} else {
		// Connect to local bitcoin core RPC server using HTTP POST mode.
		connCfg = &rpcclient.ConnConfig{
			Host:         rpcParams.Host,
			Endpoint:     "ws",
			User:         rpcParams.RPCUser,
			Pass:         rpcParams.RPCPasswd,
			HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
			DisableTLS:   true, // Bitcoin core does not provide TLS by default
			//Certificates: certs, // btcwallet provide TLS by default
		}
	}

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type GetUTXOParams struct {
	Address      string `json:"address"`
	Minconf      int    `json:"minconf"`
	Maxconf      int    `json:"maxconf"`
	MaximumCount int    `json:"maximumCount"`
}

type UTXO struct {
	TxID   string  `json:"txid"`
	Vout   uint32  `json:"vout"`
	Amount float64 `json:"amount"`
}
type GetUTXOResult struct {
	UTXOs []UTXO `json:"utxos"`
}

func GetUnspendUTXO(params string, rpcParams *RPCParams, netID int) string {
	//convert params from json format
	var getUTXOParams GetUTXOParams
	err := json.Unmarshal([]byte(params), &getUTXOParams)
	if err != nil {
		return err.Error()
	}

	//chainnet
	var realNet *chaincfg.Params
	if netID == NETID_MAIN {
		realNet = &chaincfg.MainNetParams
	} else {
		realNet = &chaincfg.TestNet3Params
	}

	//convert address from string
	address := strings.TrimSpace(getUTXOParams.Address) //Trim whitespace
	if len(address) == 0 {
		return "Params error : NO addresss."
	}
	addr, err := btcutil.DecodeAddress(address, realNet)
	if err != nil {
		return err.Error()
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	defer client.Shutdown()

	//get all raw transaction
	var strs []string
	accout := addr.String()
	count := 999999
	msgTxs, err := client.SearchRawTransactions(addr, 0, count, true, strs)
	if err != nil {
		return "Search : " + err.Error()
	}

	//save utxo to map, check next one transanction is spend or not
	outputIndex := map[string]int64{}
	sep := "-"

	//the result for return
	for _, msgTx := range msgTxs {
		//transaction inputs
		for _, in := range msgTx.TxIn {
			//check is spend or not
			_, exist := outputIndex[in.PreviousOutPoint.Hash.String()+sep+
				strconv.Itoa(int(in.PreviousOutPoint.Index))]
			if exist { //spend
				delete(outputIndex, in.PreviousOutPoint.Hash.String()+sep+
					strconv.Itoa(int(in.PreviousOutPoint.Index)))
			}
		}

		//transaction outputs
		for outIndex, out := range msgTx.TxOut {
			_, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript,
				&chaincfg.TestNet3Params)
			if err != nil {
				continue
			} else {
				if addrs[0].String() == accout {
					//fmt.Println("lock: ", hex.EncodeToString(out.PkScript))
					outputIndex[msgTx.TxHash().String()+sep+strconv.Itoa(outIndex)] = out.Value
				}
			}
		}
	}

	//compute total Amount for balance
	var result GetUTXOResult
	for oneOut, value := range outputIndex {
		keys := strings.Split(oneOut, sep)
		if len(keys) == 2 {
			vout, _ := strconv.Atoi(keys[1])
			//remove e+8
			bigFloat := new(big.Float)
			bigFloat.SetInt64(value)
			bigFloat.Mul(bigFloat, big.NewFloat(1e-8))
			amount, _ := bigFloat.Float64()
			oneUTXO := UTXO{keys[0], uint32(vout), amount}
			result.UTXOs = append(result.UTXOs, oneUTXO)
		} else {
			return "Process fatal error : key invalid."
		}

	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return err.Error()
	}

	return string(jsonResult)
}

func GetBalance(getBalanceParams *adaptor.GetBalanceParams, rpcParams *RPCParams, netID int) (string, error) {
	//chainnet
	var realNet *chaincfg.Params
	if netID == NETID_MAIN {
		realNet = &chaincfg.MainNetParams
	} else {
		realNet = &chaincfg.TestNet3Params
	}

	//convert address from string
	var addrs []btcutil.Address
	if len(getBalanceParams.Address) != 0 {
		addr, err := btcutil.DecodeAddress(getBalanceParams.Address, realNet)
		if err != nil {
			return "", err
		}
		addrs = append(addrs, addr)
	}
	if len(addrs) != 1 {
		return "", errors.New("Params error : Must one address.")
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		return "", err
	}
	defer client.Shutdown()

	//get all raw transaction
	var strs []string
	account := addrs[0].String()
	count := 999999
	msgTxs, err := client.SearchRawTransactions(addrs[0], 0, count, true, strs)
	if err != nil {
		return "", err
	}

	//save utxo to map, check next one transanction is spend or not
	msgIndex := map[string]int{}

	//the result for return
	var transAll adaptor.TransactionsResult
	for index, msgTx := range msgTxs {
		//one transaction result
		var transOne adaptor.Transaction
		transOne.TxHash = msgTx.TxHash().String()

		//		jsonTX, _ := json.Marshal(msgTx)
		//		fmt.Println(string(jsonTX))

		//transaction inputs
		isSpend := false
		for _, in := range msgTx.TxIn {
			//check is spend or not
			index, exist := msgIndex[in.PreviousOutPoint.Hash.String()+
				strconv.Itoa(int(in.PreviousOutPoint.Index))]
			if exist { //spend
				isSpend = true
				transOne.Inputs = append(transOne.Inputs,
					adaptor.InputIndex{in.PreviousOutPoint.Hash.String(),
						in.PreviousOutPoint.Index,
						transAll.Transactions[index].Outputs[in.PreviousOutPoint.Index].Addr,
						transAll.Transactions[index].Outputs[in.PreviousOutPoint.Index].Value})
			} else { //recv
				//to get addr and value
				addr, value := getAddrValue(client, realNet,
					&in.PreviousOutPoint.Hash,
					int(in.PreviousOutPoint.Index))
				if 0 == value {
					continue
				}
				transOne.Inputs = append(transOne.Inputs,
					adaptor.InputIndex{in.PreviousOutPoint.Hash.String(),
						in.PreviousOutPoint.Index,
						addr, value})
			}
		}

		//transaction outputs
		for outIndex, out := range msgTx.TxOut {
			_, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript,
				&chaincfg.TestNet3Params)
			if err != nil {
				continue
			} else {
				transOne.Outputs = append(transOne.Outputs,
					adaptor.OutputIndex{uint32(outIndex), addrs[0].String(), out.Value})
				if addrs[0].String() == account {
					msgIndex[msgTx.TxHash().String()+strconv.Itoa(outIndex)] = index
				}
			}
		}

		//calculate blancechanged
		if isSpend {
			totalInput := int64(0)
			for _, in := range transOne.Inputs {
				if account == in.Addr {
					totalInput += in.Value
				}
			}
			totalOutput := int64(0)
			for _, out := range transOne.Outputs {
				if account == out.Addr {
					totalOutput += out.Value
				}
			}
			//spend return detract from total input
			transOne.BlanceChanged = totalOutput - totalInput
		} else {
			totalRecv := int64(0)
			for _, out := range transOne.Outputs {
				if account == out.Addr {
					totalRecv += out.Value
				}
			}
			transOne.BlanceChanged = totalRecv
		}

		//add to result for return
		transAll.Transactions = append(transAll.Transactions, transOne)
	}

	//compute total Amount for balance
	var result adaptor.GetBalanceResult
	var allAmount int64
	for _, resultOne := range transAll.Transactions {
		allAmount += resultOne.BlanceChanged
	}

	//remove e+8
	bigFloat := new(big.Float)
	bigFloat.SetInt64(allAmount)
	bigFloat.Mul(bigFloat, big.NewFloat(1e-8))
	result.Value, _ = bigFloat.Float64()
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

func getAddrValue(client *rpcclient.Client, chainParams *chaincfg.Params,
	txHash *chainhash.Hash, index int) (addr string, value int64) {
	//get raw transaction by txHash
	tx, err := client.GetRawTransaction(txHash)
	if err != nil {
		log.Fatal(err)
		return "", 0
	}

	//get addr and value by index
	msgTx := tx.MsgTx()
	if index < len(msgTx.TxOut) {
		_, addrs, _, err := txscript.ExtractPkScriptAddrs(
			msgTx.TxOut[index].PkScript, chainParams)
		if err != nil {
			log.Fatal(err)
			return "", 0
		} else {
			//fmt.Println(addrs)
			return addrs[0].String(), msgTx.TxOut[index].Value
		}
	}
	//return empty if error
	return "", 0
}

func GetTransactions(getTransactionsParams *adaptor.GetTransactionsParams, rpcParams *RPCParams, netID int) (string, error) {
	//	//convert params from json format
	//	var getTransactionsParams GetTransactionsParams
	//	err := json.Unmarshal([]byte(params), &getTransactionsParams)
	//	if err != nil {
	//		log.Fatal(err)
	//		return err.Error()
	//	}

	//chainnet
	var realNet *chaincfg.Params
	if netID == NETID_MAIN {
		realNet = &chaincfg.MainNetParams
	} else {
		realNet = &chaincfg.TestNet3Params
	}

	//convert address from string
	addr, err := btcutil.DecodeAddress(getTransactionsParams.Account, realNet)
	if err != nil {
		return "", err
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		return "", err
	}
	defer client.Shutdown()

	//get all raw transaction
	var strs []string
	msgTxs, err := client.SearchRawTransactions(addr, 0, getTransactionsParams.Count, true, strs)
	if err != nil {
		return "", err
	}

	//save utxo to map, check next one transanction is spend or not
	msgIndex := map[string]int{}

	//the result for return
	var transAll adaptor.TransactionsResult
	for index, msgTx := range msgTxs {
		//one transaction result
		var transOne adaptor.Transaction
		transOne.TxHash = msgTx.TxHash().String()

		//		jsonTX, _ := json.Marshal(msgTx)
		//		fmt.Println(string(jsonTX))

		//transaction inputs
		isSpend := false
		for _, in := range msgTx.TxIn {
			//check is spend or not
			index, exist := msgIndex[in.PreviousOutPoint.Hash.String()+
				strconv.Itoa(int(in.PreviousOutPoint.Index))]
			if exist { //spend
				isSpend = true
				transOne.Inputs = append(transOne.Inputs,
					adaptor.InputIndex{in.PreviousOutPoint.Hash.String(),
						in.PreviousOutPoint.Index,
						transAll.Transactions[index].Outputs[in.PreviousOutPoint.Index].Addr,
						transAll.Transactions[index].Outputs[in.PreviousOutPoint.Index].Value})
			} else { //recv
				//to get addr and value
				addr, value := getAddrValue(client, realNet,
					&in.PreviousOutPoint.Hash,
					int(in.PreviousOutPoint.Index))
				if 0 == value {
					continue
				}
				transOne.Inputs = append(transOne.Inputs,
					adaptor.InputIndex{in.PreviousOutPoint.Hash.String(),
						in.PreviousOutPoint.Index,
						addr, value})
			}
		}

		//transaction outputs
		for outIndex, out := range msgTx.TxOut {
			_, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript,
				&chaincfg.TestNet3Params)
			if err != nil {
				continue
			} else {
				transOne.Outputs = append(transOne.Outputs,
					adaptor.OutputIndex{uint32(outIndex), addrs[0].String(), out.Value})
				if addrs[0].String() == getTransactionsParams.Account {
					msgIndex[msgTx.TxHash().String()+strconv.Itoa(outIndex)] = index
				}
			}
		}

		//calculate blancechanged
		if isSpend {
			totalInput := int64(0)
			for _, in := range transOne.Inputs {
				if getTransactionsParams.Account == in.Addr {
					totalInput += in.Value
				}
			}
			totalOutput := int64(0)
			for _, out := range transOne.Outputs {
				if getTransactionsParams.Account == out.Addr {
					totalOutput += out.Value
				}
			}
			//spend return detract from total input
			transOne.BlanceChanged = totalOutput - totalInput
		} else {
			totalRecv := int64(0)
			for _, out := range transOne.Outputs {
				if getTransactionsParams.Account == out.Addr {
					totalRecv += out.Value
				}
			}
			transOne.BlanceChanged = totalRecv
		}

		//add to result for return
		transAll.Transactions = append(transAll.Transactions, transOne)
	}

	jsonResult, err := json.Marshal(transAll)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}
