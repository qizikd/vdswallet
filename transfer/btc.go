package transfer

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/golang/glog"
	"github.com/qizikd/btcd/btcjson"
	"github.com/qizikd/btcd/rpcclient"
	"github.com/qizikd/vdswallet/rpc"
)

var WalletPassWord = "123456"
var WalletHost string = "http://47.244.57.162:13982"
var WalletRpcUser string = "coninbkb"
var WalletRpcPwd string = "coinbkb2019!"

func newOmniClient() (client *rpcclient.Client, err error) {
	client, err = rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host: WalletHost,
		User: WalletRpcUser,
		Pass: WalletRpcPwd,
	}, nil)
	return
}

func GetNewAddress() (address string, err error) {
	client, err := rpc.DialHTTP(WalletHost,&rpc.AuthCfg{User:WalletRpcUser,PassWord:WalletRpcPwd})
	if err != nil {
		glog.Error("连接节点失败", err.Error())
		return
	}
	defer client.Close()
	err = client.Call(&address, "getnewaddress")
	if err != nil {
		glog.Error("生成地址失败", err.Error())
		return
	}
	return
}

func GetWalletBalance() (balance float64, err error) {
	client, err := rpc.DialHTTP(WalletHost,&rpc.AuthCfg{User:WalletRpcUser,PassWord:WalletRpcPwd})
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return 0, errors.New("连接节点失败")
	}
	defer client.Close()
	//var amount float32
	err = client.Call(&balance, "getbalance")
	if err != nil {
		glog.Error(fmt.Sprintf("读取余额失败: %s", err.Error()))
		return 0, errors.New("读取余额失败")
	}
	return
}

func ListWalletTransactions(start int,end int) (txs *[]btcjson.ListTransactionsResult,err error) {
	client, err := rpc.DialHTTP(WalletHost,&rpc.AuthCfg{User:WalletRpcUser,PassWord:WalletRpcPwd})
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return nil,errors.New("连接节点失败")
	}
	defer client.Close()
	//var amount float3
	err = client.Call(&txs, "listtransactions",end,start)
	if err != nil {
		glog.Error(fmt.Sprintf("读取余额失败: %s", err.Error()))
		return nil,errors.New("读取余额失败")
	}
	return
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

func SendBtc(toAddress string, amount float64) (txid string, err error) {
	client, err := rpc.DialHTTP(WalletHost,&rpc.AuthCfg{User:WalletRpcUser,PassWord:WalletRpcPwd})
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return "",errors.New("连接节点失败")
	}
	defer client.Close()
	////先解锁钱包
	//err = client.WalletPassphrase(WalletPassWord, 30)
	//if err != nil {
	//	glog.Error("钱包解锁失败: ", err)
	//	return "", errors.New("钱包解锁失败")
	//}
	//发送交易
	err = client.Call(&txid, "sendtoaddress",toAddress,amount)
	if err != nil {
		glog.Error(fmt.Sprintf("发送交易失败: %s", err.Error()))
		return
	}
	//锁定钱包
	//client.WalletLock()
	return
}
