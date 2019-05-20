package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/qizikd/vdswallet/routerHandler"
	"github.com/qizikd/vdswallet/transfer"
	"net/http"
)

func main() {
	port := flag.String("port", "80", "Listen port")
	pwd := flag.String("pwd", "123456", "wallet password")
	wallethost := flag.String("wallethost", "http://127.0.0.1:8091/", "wallet host")
	walletrpcuser := flag.String("walletrpcuser", "", "wallet rpc user")
	walletrpcpwd := flag.String("walletrpcpwd", "", "wallet rpc password")
	flag.Parse()

	transfer.WalletPassWord = *pwd
	transfer.WalletHost = *wallethost
	transfer.WalletRpcUser = *walletrpcuser
	transfer.WalletRpcPwd = *walletrpcpwd
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
	//获取充值地址
	router.GET("/wallet/vds/newaddress", routerHandler.GetNewAddress)
	//获取钱包余额
	router.GET("/wallet/vds/balance", routerHandler.GetBalance)
	//获取钱包交易记录
	router.GET("/wallet/vds/txs", routerHandler.GetTransactions)
	//提币 余额归拢
	router.POST("/wallet/vds/sendto", routerHandler.SendBtc)

	router.Run(":" + *port)
}
