package transfer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bcext/cashutil/base58"
	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"github.com/qizikd/btcd/btcjson"
	"github.com/qizikd/btcd/rpcclient"
	"github.com/qizikd/walletMiddleware/token"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var IsBTCTestNet3 = false
var WalletPassWord = "123456"

func newOmniClient() (client *rpcclient.Client, err error) {
	if IsBTCTestNet3 {
		client, err = rpcclient.New(&rpcclient.ConnConfig{
			HTTPPostMode: true,
			DisableTLS:   true,
			//rpc.blockchain.info
			Host: "47.92.148.83:8888",
			User: "omnicorerpc",
			Pass: "abcd1234",
		}, nil)
	} else {
		client, err = rpcclient.New(&rpcclient.ConnConfig{
			HTTPPostMode: true,
			DisableTLS:   true,
			//rpc.blockchain.info
			Host: "127.0.0.1:8332",
			User: "omnicorerpc",
			Pass: "abcd1234",
		}, nil)
	}
	return
}

func newGobcy(coin string) (api gobcy.API) {
	//随机获取一个GobcyApi key，防止用一个有请求限制
	gobcyApis := []string{"9184cf751ace44f090769b52643ade0b", "269d9eb40f3048a6875b45e5aee017e9", "64cb8fe59b934d8d9df104fa8d59a85b",
		"3dc9ad5c6d8449a499103de610ab12d8", "c2c26b546bf04f049ea06e7e539d868a", "dd4b3e08a28347dfa222469472ad1a73"}
	gobcyApiFile := "/home/gopath/bin/gobcyApiToken.json"
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		gobcyApiFile = dir + "/gobcyApiToken.json"
	}
	content, err := ioutil.ReadFile(gobcyApiFile)
	if err == nil {
		json.Unmarshal(content, &gobcyApis)
	} else {
		glog.Error("加载gobcyApiToken.json出错", err)
	}
	rand.Seed(time.Now().Unix())
	apiIndex := rand.Intn(len(gobcyApis))
	var gobcyApi = gobcyApis[apiIndex]
	if IsBTCTestNet3 {
		api = gobcy.API{gobcyApi, coin, "test3"}
	} else {
		api = gobcy.API{gobcyApi, coin, "main"}
	}
	return
}

func GetNewBtcAddress(account string) (address string, err error) {
	client, err := newOmniClient()
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return "", errors.New("连接节点失败")
	}
	defer client.Disconnect()
	addr, err := client.GetNewAddress(account)
	if err != nil {
		glog.Error(fmt.Sprintf("从节点生成地址失败(%s): %s", account, err.Error()))
		return "", errors.New("生成地址失败")
	}
	return addr.EncodeAddress(), nil
}

func GetBtcBalance(account string) (balance int64, err error) {
	client, err := newOmniClient()
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return 0, errors.New("连接节点失败")
	}
	defer client.Disconnect()
	amount, err := client.GetBalance(account)
	if err != nil {
		glog.Error(fmt.Sprintf("读取余额失败(%s): %s", account, err.Error()))
		return 0, errors.New("读取余额失败")
	}
	return int64(amount), nil
}

func GetUsdtBalance(address string) (balance int64, err error) {
	client, err := newOmniClient()
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return 0, errors.New("连接节点失败")
	}
	defer client.Disconnect()
	var amount int
	if IsBTCTestNet3 {
		amount, err = client.GetOmniBalance(address, 1)
	} else {
		amount, err = client.GetOmniBalance(address, 31)
	}
	if err != nil {
		glog.Error(fmt.Sprintf("读取余额失败(%s): %s", address, err.Error()))
		return 0, errors.New("读取余额失败")
	}
	return int64(amount), nil
}

func GetBtcAddressReceive(addr btcutil.Address) (balance int64, err error) {
	client, err := newOmniClient()
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return 0, errors.New("连接节点失败")
	}
	defer client.Disconnect()
	amount, err := client.GetReceivedByAddress(addr)
	if err != nil {
		glog.Error(fmt.Sprintf("读取接收金额失败(%s): %s", addr.EncodeAddress(), err.Error()))
		return 0, errors.New("读取接收金额失败")
	}
	return int64(amount), nil
}

func GetUsdtTransactions(address string, count int, skip int) (result []rpcclient.Omni_ListtransactionResult, err error) {
	client, err := newOmniClient()
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return []rpcclient.Omni_ListtransactionResult{}, errors.New("连接节点失败")
	}
	defer client.Disconnect()
	result, err = client.Omni_Listtransactions(address, count, skip)
	if err != nil {
		glog.Error(fmt.Sprintf("获取地址交易列表失败: %s", err.Error()))
		return []rpcclient.Omni_ListtransactionResult{}, errors.New("获取地址交易列表失败")
	}
	return
}

func GetBtcAllAddressReceive() (result []btcjson.ListReceivedByAddressResult, err error) {
	client, err := newOmniClient()
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return []btcjson.ListReceivedByAddressResult{}, errors.New("连接节点失败")
	}
	defer client.Disconnect()
	result, err = client.ListReceivedByAddress()
	if err != nil {
		glog.Error(fmt.Sprintf("获取地址接收金额列表失败: %s", err.Error()))
		return []btcjson.ListReceivedByAddressResult{}, errors.New("获取地址接收金额列表失败")
	}
	return
}

func GetBtcTxInfo(txhash *chainhash.Hash) (result *btcjson.GetTransactionResult, err error) {
	client, err := newOmniClient()
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return nil, errors.New("连接节点失败")
	}
	defer client.Disconnect()
	result, err = client.GetTransaction(txhash)
	if err != nil {
		glog.Error(fmt.Sprintf("读取交易详情失败: %s", err.Error()))
		return nil, errors.New("读取交易详情失败")
	}
	return
}

func GetUsdtTxInfo(txhash *chainhash.Hash) (result rpcclient.Omni_ListtransactionResult, err error) {
	client, err := newOmniClient()
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return rpcclient.Omni_ListtransactionResult{}, errors.New("连接节点失败")
	}
	defer client.Disconnect()
	result, err = client.Omni_Gettransaction(txhash.String())
	if err != nil {
		glog.Error(fmt.Sprintf("读取交易详情失败: %s", err.Error()))
		return rpcclient.Omni_ListtransactionResult{}, errors.New("读取交易详情失败")
	}
	return
}

func SendBtc(account string, toAddress btcutil.Address, amount btcutil.Amount) (txid string, err error) {
	client, err := newOmniClient()
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return "", errors.New("连接节点失败")
	}
	defer client.Disconnect()
	//先解锁钱包
	err = client.WalletPassphrase(WalletPassWord, 30)
	if err != nil {
		glog.Error("钱包解锁失败: ", err)
		return "", errors.New("钱包解锁失败")
	}
	//发送交易
	txHash, err := client.SendFrom(account, toAddress, amount)
	if err != nil {
		glog.Error(fmt.Sprintf("发送交易失败(%s): %s", account, err.Error()))
		return
	}
	//锁定钱包
	client.WalletLock()
	return txHash.String(), nil
}

func SendUsdt(fromAddress string, toAddress btcutil.Address, amount string) (txid string, err error) {
	client, err := newOmniClient()
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return "", errors.New("连接节点失败")
	}
	defer client.Disconnect()
	//先解锁钱包
	err = client.WalletPassphrase(WalletPassWord, 30)
	if err != nil {
		glog.Error("钱包解锁失败: ", err)
		return "", errors.New("钱包解锁失败")
	}
	//正式链31，测试链1
	var propertyid = 31
	if IsBTCTestNet3 {
		propertyid = 1
	}
	//发送交易
	txHash, err := client.OmniSend(fromAddress, toAddress.EncodeAddress(), amount, propertyid)
	if err != nil {
		glog.Error(fmt.Sprintf("发送交易失败(%s): %s", fromAddress, err.Error()))
		return
	}
	//锁定钱包
	client.WalletLock()
	return txHash, nil
}

func GetBalance(address string, coin string) (balance int, err error) {
	bcy := newGobcy(coin)
	addr, err := bcy.GetAddrBal(address, nil)
	if err != nil {
		glog.Error(fmt.Sprintf("请求获取余额失败(%s): %s", bcy.Token, err))
		return 0, errors.New(fmt.Sprintf("请求获取余额失败: %s", err))
	}
	return addr.Balance, nil
}

func Transaction(fromAddress string, toAddress string, privateKey string, amount int, coin string) (tx string, err error) {
	bcy := newGobcy(coin)
	//讲私匙从wif格式转换为原始格式
	privwif := privateKey
	privb, _, err := base58.CheckDecode(privwif)
	if err != nil {
		glog.Error(err)
		return "", err
	}
	err = VerifyBtcAddressByPriv(privateKey, fromAddress)
	if err != nil {
		glog.Error("私钥错误:", err)
		return "", errors.New("私钥错误")
	}
	privstr := hex.EncodeToString(privb)
	privstr = privstr[0 : len(privstr)-2]
	//Post New TXSkeleton
	trans := gobcy.TempNewTX(fromAddress, toAddress, amount)
	skel, err := bcy.NewTX(trans, false)
	//Sign it locally
	var priv []string
	for i := 0; i < len(skel.ToSign); i++ {
		priv = append(priv, privstr)
	}
	err = skel.Sign(priv)
	if err != nil {
		glog.Error(err)
		return "", err
	}
	//Send TXSkeleton
	skel, err = bcy.SendTX(skel)
	if err != nil {
		glog.Error(err)
		return "", err
	}
	return skel.Trans.Hash, nil
}

func VerifyBtcAddressByPriv(privwif string, address string) (err error) {
	private_wif, err := btcutil.DecodeWIF(privwif)
	if err != nil {
		glog.Error("私钥格式不正确:", err)
		return errors.New("私钥格式不正确")
	}
	addr, err := btcutil.NewAddressPubKey(private_wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil {
		glog.Error("生成地址失败:", err)
		return errors.New("生成地址失败")
	}
	if strings.ToUpper(addr.EncodeAddress()) != strings.ToUpper(address) {
		return errors.New("验证失败")
	}
	return nil
}

func GetEthBalance(address string) (amount big.Int, err error) {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		glog.Error("连接infura节点失败", err)
		return amount, errors.New("连接infura节点失败")
	}
	defer client.Close()
	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		glog.Error("获取Eth余额失败", err)
		return amount, fmt.Errorf("获取Eth余额失败:%s", err)
	}
	return *balance, nil
}

func GetTokenBalance(tokenAddress string, address string) (amount big.Int, err error) {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		glog.Error("连接infura节点失败", err)
		return amount, errors.New("连接infura节点失败")
	}
	defer client.Close()
	instance, err := token.NewToken(common.HexToAddress(tokenAddress), client)
	if err != nil {
		glog.Error(err)
		return
	}
	balance, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(address))
	if err != nil {
		glog.Error("获取token balance失败", err)
		return amount, fmt.Errorf("获取token balance余额失败:%s", err)
	}
	return *balance, nil
}

func TransactionEth(fromAddress string, toAddress string, privateKey string, amount big.Int) (tx string, err error) {
	tx, err = sendRawTransaction(fromAddress, privateKey, toAddress, nil, amount)
	if err != nil {
		glog.Error("sendRawTransaction:", err)
		return
	}
	return tx, nil
}

func TransactionToken(fromAddress string, toAddress string, tokenAddress string, privateKey string, tokenAmount big.Int) (tx string, err error) {
	data := buildTransfer(toAddress, tokenAmount)
	tx, err = sendRawTransaction(fromAddress, privateKey, tokenAddress, data, *big.NewInt(0))
	if err != nil {
		glog.Error("sendRawTransaction:", err)
		return
	}
	return
}

func buildTransfer(toAddressHex string, tokenAmount big.Int) (data []byte) {
	toAddress := common.HexToAddress(toAddressHex)
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	//生成data methodABI+参数
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(tokenAmount.Bytes(), 32)
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	return
}

func sendRawTransaction(fromAddressHex string, privatekey string, toAddressHex string, data []byte, value big.Int) (tx string, err error) {
	//发送方私匙
	if privatekey[0:2] == "0x" {
		privatekey = privatekey[2:]
	}
	privateKey, err := crypto.HexToECDSA(privatekey)
	if err != nil {
		glog.Error("私钥格式不正确", err)
		return "", errors.New("私钥格式不正确")
	}
	err = VerifyEthAddressByPriv(privatekey, fromAddressHex)
	if err != nil {
		glog.Error("私钥错误:", err)
		return "", errors.New("私钥错误")
	}
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		glog.Error("连接infura节点失败", err)
		return "", errors.New("连接infura节点失败")
	}
	defer client.Close()
	//获取noce值
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(fromAddressHex))
	if err != nil {
		glog.Error("获取noce值失败:", err)
		return "", errors.New("获取noce值失败")
	}
	//获取gasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		glog.Error("获取gasPrice失败:", err)
		return "", errors.New("获取手续费价格失败")
	}
	//fmt.Println("gasPrice:", gasPrice)
	toAddress := common.HexToAddress(toAddressHex)

	//获取gasLimit
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     common.HexToAddress(fromAddressHex),
		To:       &toAddress,
		GasPrice: gasPrice,
		Data:     data,
	})
	if err != nil {
		glog.Error("获取gasLimit失败:", err)
		return "", errors.New("获取手续费失败")
	}

	//gasLimit := uint64(300000)
	//fmt.Printf("gasLimit:%d\n", gasLimit) // 23256
	//创建一个交易对象
	transaction := types.NewTransaction(nonce, toAddress, &value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		glog.Error("创建chainId失败:", err)
		return "", errors.New("创建chainID失败")
	}

	signedTx, err := types.SignTx(transaction, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		glog.Error("签名交易失败:", err)
		return "", errors.New("签名交易失败")
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		glog.Error("发送交易失败:", err)
		return "", errors.New("发送交易失败")
	}
	return signedTx.Hash().Hex(), nil
}

func VerifyEthAddressByPriv(privwif string, address string) (err error) {
	privatekey := privwif
	if privatekey[0:2] == "0x" {
		privatekey = privatekey[2:]
	}
	var privkey = new(big.Int)
	privkey.SetString(privatekey, 16)
	pk, err := crypto.ToECDSA(privkey.Bytes())
	if err != nil {
		glog.Error("私钥格式不正确:", err)
		return errors.New("私钥格式不正确")
	}
	addr := crypto.PubkeyToAddress(pk.PublicKey)
	if strings.ToUpper(addr.Hex()) != strings.ToUpper(address) {
		glog.Error(address, " != ", addr.Hex())
		return errors.New("验证失败")
	}
	return nil
}
