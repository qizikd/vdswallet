package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/qizikd/walletMiddleware/routerHandler"
	"github.com/qizikd/walletMiddleware/transfer"
	"net/http"
)

func main() {
	port := flag.String("port", "80", "Listen port")
	isTest := flag.String("net", "normal", "network")
	pwd := flag.String("pwd", "123456", "wallet password")
	flag.Parse()
	if *isTest == "test" {
		transfer.IsBTCTestNet3 = true
	}
	transfer.WalletPassWord = *pwd
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
	//根据账户名获取btc充值地址
	router.GET("/walletmiddleware/btc/newaddress", routerHandler.GetNewAddress)
	//获取账户btc余额
	router.GET("/walletmiddleware/btc/accountbalance", routerHandler.GetBtcAccountBalance)
	//获取指定地址接收的金额
	router.GET("/walletmiddleware/btc/addressreceive", routerHandler.GetBtcAddressReceive)
	//获取钱包所有地址接收的金额
	router.GET("/walletmiddleware/btc/alladdressreceive", routerHandler.GetBtcAllAddressReceive)
	//获取交易详情
	router.GET("/walletmiddleware/btc/tx", routerHandler.GetBtcTxInfo)
	//提币 余额归拢
	router.POST("/walletmiddleware/btc/send", routerHandler.SendBtc)

	//usdt归拢流程，先请求账户btc余额，是否足够手续费 不够的话请先打入手续费
	//根据账户名获取usdt充值地址  同btc
	router.GET("/walletmiddleware/usdt/newaddress", routerHandler.GetNewAddress)
	//获取地址usdt余额
	router.GET("/walletmiddleware/usdt/balance", routerHandler.GetUsdtBalance)
	//获取指定地址usdt交易记录
	router.GET("/walletmiddleware/usdt/txs", routerHandler.GetUsdtTxs)
	//获取交易详情
	router.GET("/walletmiddleware/usdt/tx", routerHandler.GetUsdtTxInfo)
	//ustd提币 余额归拢
	router.POST("/walletmiddleware/usdt/send", routerHandler.SendUsdt)

	//打入手续费
	router.POST("/walletmiddleware/btc/transferfee", routerHandler.TransferFee)
	//查询地址btc余额
	router.GET("/walletmiddleware/btc/balance", routerHandler.Balance)

	//eth余额查询
	router.GET("/walletmiddleware/eth/balance", routerHandler.BalanceEth)
	//eth代币余额查询
	router.GET("/walletmiddleware/eth/tokenbalance", routerHandler.BalanceEthToken)
	//eth转账
	router.POST("/walletmiddleware/eth/sendto", routerHandler.SendEthTo)
	//eth代币转账
	router.POST("/walletmiddleware/eth/tokensendto", routerHandler.SendEthTokenTo)
	//生成eth地址
	router.POST("/walletmiddleware/eth/newaddress", routerHandler.GenerateAddress)
	router.Run(":" + *port)
}
