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
func (abtc AdaptorBTC) CreateMultiSigAddress(params *CreateMultiSigParams) (string, error) {
	return CreateMultiSigAddress(params, abtc.NetID)
}

func (abtc AdaptorBTC) GetUnspendUTXO(params string) string {
	return GetUnspendUTXO(params, &abtc.RPCParams, abtc.NetID)
}

func (abtc AdaptorBTC) RawTransactionGen(params *RawTransactionGenParams) (string, error) {
	return RawTransactionGen(params, &abtc.RPCParams, abtc.NetID)
}
func (abtc AdaptorBTC) DecodeRawTransaction(params *DecodeRawTransactionParams) (string, error) {
	return DecodeRawTransaction(params, &abtc.RPCParams)
}
func (abtc AdaptorBTC) GetTransactionByHash(params *GetTransactionByHashParams) (string, error) {
	return GetTransactionByHash(params, &abtc.RPCParams)
}
func (abtc AdaptorBTC) SignTransaction(params *SignTransactionParams) (string, error) {
	return SignTransaction(params, &abtc.RPCParams, abtc.NetID)
}
func (abtc AdaptorBTC) SignTxSend(params *SignTxSendParams) (string, error) {
	return SignTxSend(params, &abtc.RPCParams, abtc.NetID)
}
func (abtc AdaptorBTC) GetBalance(params *GetBalanceParams) (string, error) {
	return GetBalance(params, &abtc.RPCParams, abtc.NetID)
}
func (abtc AdaptorBTC) GetTransactions(params *GetTransactionsParams) (string, error) {
	return GetTransactions(params, &abtc.RPCParams, abtc.NetID)
}

func (abtc AdaptorBTC) ImportMultisig(params *ImportMultisigParams) (string, error) {
	return ImportMultisig(params, &abtc.RPCParams, abtc.NetID)
}

func (abtc AdaptorBTC) SendTransaction(params string) string {
	return SendTransaction(params, &abtc.RPCParams)
}
