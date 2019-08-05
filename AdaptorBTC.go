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
	"github.com/palletone/adaptor"
)

type RPCParams struct {
	Host      string `json:"host"`
	RPCUser   string `json:"rpcUser"`
	RPCPasswd string `json:"rpcPasswd"`
	CertPath  string `json:"certPath"`
}

type AdaptorBTC struct {
	NetID int
	RPCParams
}

const (
	NETID_MAIN = iota
	NETID_TEST
)

func (abtc AdaptorBTC) NewPrivateKey() (wifPriKey string) {
	return NewPrivateKey(abtc.NetID)
}
func (abtc AdaptorBTC) GetPublicKey(wifPriKey string) (pubKey string) {
	return GetPublicKey(wifPriKey, abtc.NetID)
}
func (abtc AdaptorBTC) GetAddress(wifPriKey string) (address string) {
	return GetAddress(wifPriKey, abtc.NetID)
}
func (abtc AdaptorBTC) GetAddressByPubkey(pubKeyHex string) (string, error) {
	return GetAddressByPubkey(pubKeyHex, abtc.NetID)
}
func (abtc AdaptorBTC) CreateMultiSigAddress(params *adaptor.CreateMultiSigParams) (*adaptor.CreateMultiSigResult, error) {
	return CreateMultiSigAddress(params, abtc.NetID)
}

func (abtc AdaptorBTC) GetUTXO(params *adaptor.GetUTXOParams) (*adaptor.GetUTXOResult, error) {
	return GetUTXO(params, &abtc.RPCParams, abtc.NetID)
}
func (abtc AdaptorBTC) GetUTXOHttp(params *adaptor.GetUTXOHttpParams) (*adaptor.GetUTXOHttpResult, error) {
	return GetUTXOHttp(params, abtc.NetID)
}

func (abtc AdaptorBTC) RawTransactionGen(params *adaptor.RawTransactionGenParams) (*adaptor.RawTransactionGenResult, error) {
	return RawTransactionGen(params, abtc.NetID)
}
func (abtc AdaptorBTC) DecodeRawTransaction(params *adaptor.DecodeRawTransactionParams) (*adaptor.DecodeRawTransactionResult, error) {
	return DecodeRawTransaction(params, abtc.NetID)
}
func (abtc AdaptorBTC) GetTransactionByHash(params *adaptor.GetTransactionByHashParams) (*adaptor.GetTransactionByHashResult, error) {
	return GetTransactionByHash(params, &abtc.RPCParams)
}
func (abtc AdaptorBTC) GetTransactionHttp(params *adaptor.GetTransactionHttpParams) (*adaptor.GetTransactionHttpResult, error) {
	return GetTransactionHttp(params, abtc.NetID)
}

func (abtc AdaptorBTC) SignTransaction(params *adaptor.SignTransactionParams) (*adaptor.SignTransactionResult, error) {
	return SignTransaction(params, abtc.NetID)
}
func (abtc AdaptorBTC) SignTxSend(params *adaptor.SignTxSendParams) (*adaptor.SignTxSendResult, error) {
	return SignTxSend(params, &abtc.RPCParams, abtc.NetID)
}
func (abtc AdaptorBTC) GetBalance(params *adaptor.GetBalanceParams) (*adaptor.GetBalanceResult, error) {
	return GetBalance(params, &abtc.RPCParams, abtc.NetID)
}
func (abtc AdaptorBTC) GetBalanceHttp(params *adaptor.GetBalanceHttpParams) (*adaptor.GetBalanceHttpResult, error) {
	return GetBalanceHttp(params, abtc.NetID)
}
func (abtc AdaptorBTC) GetTransactions(params *adaptor.GetTransactionsParams) (*adaptor.TransactionsResult, error) {
	return GetTransactions(params, &abtc.RPCParams, abtc.NetID)
}

func (abtc AdaptorBTC) SendTransaction(params *adaptor.SendTransactionParams) (*adaptor.SendTransactionResult, error) {
	return SendTransaction(params, &abtc.RPCParams)
}
func (abtc AdaptorBTC) SendTransactionHttp(params *adaptor.SendTransactionHttpParams) (*adaptor.SendTransactionHttpResult, error) {
	return SendTransactionHttp(params, abtc.NetID)
}

func (abtc AdaptorBTC) MergeTransaction(params *adaptor.MergeTransactionParams) (*adaptor.MergeTransactionResult, error) {
	return MergeTransaction(params, abtc.NetID)
}

func (abtc AdaptorBTC) SignMessage(params *adaptor.SignMessageParams) (*adaptor.SignMessageResult, error) {
	return SignMessage(params)
}

func (abtc AdaptorBTC) VerifyMessage(params *adaptor.VerifyMessageParams) (*adaptor.VerifyMessageResult, error) {
	return VerifyMessage(params, abtc.NetID)
}
