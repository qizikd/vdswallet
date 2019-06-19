package transfer

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/qizikd/btcd/btcjson"
	"github.com/qizikd/vdswallet/rpc"
)

var WalletPassWord = "123456"
var WalletHost = "http://127.0.0.1:8332"
var WalletRpcUser = "rpcuser"
var WalletRpcPwd = "rpcpwd"

func GetNewAddress() (address string, err error) {
	fmt.Printf("WalletHost:%s,WalletRpcUser:%s,WalletRpcPwd:%s\n", WalletHost, WalletRpcUser, WalletRpcPwd)
	client, err := rpc.DialHTTP(WalletHost, &rpc.AuthCfg{User: WalletRpcUser, PassWord: WalletRpcPwd})
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
	client, err := rpc.DialHTTP(WalletHost, &rpc.AuthCfg{User: WalletRpcUser, PassWord: WalletRpcPwd})
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

func ListWalletTransactions(offset int, count int) (txs *[]btcjson.ListTransactionsResult, err error) {
	client, err := rpc.DialHTTP(WalletHost, &rpc.AuthCfg{User: WalletRpcUser, PassWord: WalletRpcPwd})
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return nil, errors.New("连接节点失败")
	}
	defer client.Close()
	//var amount float3
	err = client.Call(&txs, "listtransactions", count, offset)
	if err != nil {
		glog.Error(fmt.Sprintf("读取余额失败: %s", err.Error()))
		return nil, errors.New("读取余额失败")
	}
	return
}

func SendBtc(toAddress string, amount float64) (txid string, err error) {
	client, err := rpc.DialHTTP(WalletHost, &rpc.AuthCfg{User: WalletRpcUser, PassWord: WalletRpcPwd})
	if err != nil {
		glog.Error("连接节点失败: ", err)
		return "", errors.New("连接节点失败")
	}
	defer client.Close()
	////先解锁钱包
	//err = client.Call(&txid, "walletpassphrase", WalletPassWord, 30)
	//if err != nil {
	//	glog.Error("钱包解锁失败: ", err)
	//	return "", errors.New("钱包解锁失败")
	//}
	////锁定钱包
	//defer client.Call(&txid, "walletlock")
	//发送交易
	err = client.Call(&txid, "sendtoaddress", toAddress, amount)
	if err != nil {
		glog.Error(fmt.Sprintf("发送交易失败: %s", err.Error()))
		return
	}
	return
}
