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
	"fmt"
	"path/filepath"
	//"fmt"
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"

	"github.com/palletone/adaptor"
)

//==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ==== ===

var GHomeDir = btcutil.AppDataDir("btcwallet", false)
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
	Addresses    []string `json:"addresses"`
	Minconf      int      `json:"minconf"`
	Maxconf      int      `json:"maxconf"`
	MaximumCount int      `json:"maximumCount"`
}

func GetUnspendUTXO(params string, rpcParams *RPCParams, netID adaptor.NetID) string {
	//convert params from json format
	var getUTXOParams GetUTXOParams
	err := json.Unmarshal([]byte(params), &getUTXOParams)
	if err != nil {
		log.Fatal(err)
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
	var addrs []btcutil.Address
	for _, address := range getUTXOParams.Addresses {
		address = strings.TrimSpace(address) //Trim whitespace
		if len(address) == 0 {
			continue
		}
		addr, err := btcutil.DecodeAddress(address, realNet)
		if err != nil {
			log.Fatal(err)
			return err.Error()
		}
		addrs = append(addrs, addr)
	}
	if len(addrs) == 0 {
		return "Params error : NO addresss."
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	defer client.Shutdown()

	//return [{utxo},{utxo}]
	results, err := client.ListUnspentMinMaxAddresses(
		getUTXOParams.Minconf, getUTXOParams.Maxconf, addrs)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}

	//
	var resultSplit []btcjson.ListUnspentResult
	if getUTXOParams.MaximumCount > 0 {
		if getUTXOParams.MaximumCount < len(results) {
			resultSplit = results[:getUTXOParams.MaximumCount]
		} else {
			resultSplit = results[:]
		}
	}
	jsonResult, err := json.Marshal(resultSplit)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}

	return string(jsonResult)
}

func GetBalance(getBalanceParams *adaptor.GetBalanceParams, rpcParams *RPCParams, netID adaptor.NetID) (string, error) {
	//	//convert params from json format
	//	var getBalanceParams GetBalanceParams
	//	err := json.Unmarshal([]byte(params), &getBalanceParams)
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

	//return [{utxo},{utxo}]
	var results []btcjson.ListUnspentResult
	results, err = client.ListUnspentMinMaxAddresses(
		getBalanceParams.Minconf, 999999, addrs)
	if err != nil {
		return "", err
	}

	//compute total Amount for balance
	var result adaptor.GetBalanceResult
	for _, resultOne := range results {
		result.Value += resultOne.Amount
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

func ImportMultisig(importMultisigParams *adaptor.ImportMultisigParams, rpcParams *RPCParams, netID adaptor.NetID) (string, error) {
	//	//convert params from json format
	//	var importMultisigParams ImportMultisigParams
	//	err := json.Unmarshal([]byte(params), &importMultisigParams)
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
	var addrs []btcutil.Address
	for _, publicKey := range importMultisigParams.PublicKeys {
		publicKey = strings.TrimSpace(publicKey) //Trim whitespace
		if len(publicKey) == 0 {
			continue
		}
		addr, err := btcutil.DecodeAddress(publicKey, realNet)
		if err != nil {
			return "", err
		}
		addrs = append(addrs, addr)
	}
	if len(addrs) == 0 {
		return "", errors.New("Params error : Must one address.")
	}
	if len(addrs) < importMultisigParams.MRequires {
		return "", errors.New("Params error : Need more publickeys.")
	}

	//get rpc client
	client, err := GetClient(rpcParams)
	if err != nil {
		return "", err
	}
	defer client.Shutdown()

	//unlock wallet 3 seconds
	err = client.WalletPassphrase(importMultisigParams.WalletPasswd, 3)
	if err != nil {
		return "", err
	}
	//add to wallet, return multsig address
	//'imported' is the account btcwallet required,
	//'imported' is used in ListUTXO too
	_, err = client.AddMultisigAddress(
		importMultisigParams.MRequires, addrs, "imported")
	if err != nil {
		return "", err
	}

	//
	var result adaptor.ImportMultisigResult
	result.Import = true

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
			fmt.Println(addrs)
			return addrs[0].String(), msgTx.TxOut[index].Value
		}
	}
	//return empty if error
	return "", 0
}

func GetTransactions(getTransactionsParams *adaptor.GetTransactionsParams, rpcParams *RPCParams, netID adaptor.NetID) (string, error) {
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
