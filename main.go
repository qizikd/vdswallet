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
	flag.Parse()

	transfer.WalletPassWord = *pwd
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
