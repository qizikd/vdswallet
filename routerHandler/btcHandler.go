package routerHandler

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/qizikd/walletMiddleware/transfer"
	"github.com/tyler-smith/go-bip39"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

func GetNewAddress(c *gin.Context) {
	account := c.Query("account")
	address, err := transfer.GetNewBtcAddress(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "获取地址失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"address": address,
		},
	})
	return
}

func GetBtcAccountBalance(c *gin.Context) {
	account := c.Query("account")
	balance, err := transfer.GetBtcBalance(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "获取余额失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"balance": balance,
		},
	})
	return
}

func GetBtcAddressReceive(c *gin.Context) {
	address := c.Query("address")
	var addr btcutil.Address
	var err error
	if transfer.IsBTCTestNet3 {
		addr, err = btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	} else {
		addr, err = btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "地址非法",
			"errorinfo": fmt.Sprintf("地址非法(%s)", err.Error()),
		})
		return
	}
	balance, err := transfer.GetBtcAddressReceive(addr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "获取余额失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"recv": balance,
		},
	})
	return
}

func GetUsdtTxs(c *gin.Context) {
	address := c.Query("address")
	count := c.Query("count")
	skip := c.Query("skip")
	_count, err := strconv.Atoi(count)
	if err != nil {
		_count = 10
	}
	_skip, err := strconv.Atoi(skip)
	if err != nil {
		_skip = 0
	}
	if transfer.IsBTCTestNet3 {
		_, err = btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	} else {
		_, err = btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "地址非法",
			"errorinfo": fmt.Sprintf("地址非法(%s)", err.Error()),
		})
		return
	}
	balance, err := transfer.GetUsdtTransactions(address, _count, _skip)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "获取交易列表失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"recv": balance,
		},
	})
	return
}

func GetBtcAllAddressReceive(c *gin.Context) {
	result, err := transfer.GetBtcAllAddressReceive()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "获取余额失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})
	return
}

func GetUsdtBalance(c *gin.Context) {
	address := c.Query("address")
	var err error
	if transfer.IsBTCTestNet3 {
		_, err = btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	} else {
		_, err = btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "地址非法",
			"errorinfo": fmt.Sprintf("地址非法(%s)", err.Error()),
		})
		return
	}
	balance, err := transfer.GetUsdtBalance(address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "获取余额失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"balance": balance,
		},
	})
	return
}

func GetBtcTxInfo(c *gin.Context) {
	txid := c.Query("txid")
	txhash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "交易哈希非法",
			"errorinfo": err.Error(),
		})
		return
	}
	result, err := transfer.GetBtcTxInfo(txhash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "获取交易详情失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})
	return
}

func GetUsdtTxInfo(c *gin.Context) {
	txid := c.Query("txid")
	txhash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "交易哈希非法",
			"errorinfo": err.Error(),
		})
		return
	}
	result, err := transfer.GetUsdtTxInfo(txhash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "获取交易详情失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})
	return
}

func SendBtc(c *gin.Context) {
	account := c.PostForm("account")
	toaddress := c.PostForm("toaddress")
	var addr btcutil.Address
	var err error
	if transfer.IsBTCTestNet3 {
		addr, err = btcutil.DecodeAddress(toaddress, &chaincfg.TestNet3Params)
	} else {
		addr, err = btcutil.DecodeAddress(toaddress, &chaincfg.MainNetParams)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "地址非法",
			"errorinfo": fmt.Sprintf("地址非法(%s)", err.Error()),
		})
		return
	}
	_amount := c.PostForm("amount")
	//防止包含小数点，去掉小数点后的数值再转换类型
	if strings.Contains(_amount, ".") {
		a := strings.Split(_amount, ".")
		_amount = a[0]
	}
	amount, err := strconv.Atoi(_amount)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "金额非法",
			"errorinfo": "金额非法",
		})
		glog.Error("金额非法")
		return
	}
	if amount < 546 {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "转账金额不能低于546聪",
			"errorinfo": "转账金额不能低于546聪",
		})
		glog.Error("转账金额不能低于546聪")
		return
	}
	balance, _ := transfer.GetBtcBalance(account)
	if balance < int64(amount) {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "余额不足",
			"errorinfo": "余额不足",
		})
		glog.Error(fmt.Sprintf("余额不足(账号:%s,余额:%d,转账金额:%d)", account, balance, amount))
		return
	}
	txid, err := transfer.SendBtc(account, addr, btcutil.Amount(amount))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "交易失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"txid": txid,
		},
	})
	return
}

func SendUsdt(c *gin.Context) {
	fromaddress := c.PostForm("fromaddress")
	toaddress := c.PostForm("toaddress")
	var addr btcutil.Address
	var err error
	if transfer.IsBTCTestNet3 {
		addr, err = btcutil.DecodeAddress(fromaddress, &chaincfg.TestNet3Params)
	} else {
		addr, err = btcutil.DecodeAddress(fromaddress, &chaincfg.MainNetParams)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "转出地址非法",
			"errorinfo": fmt.Sprintf("转出地址非法(%s)", err.Error()),
		})
		return
	}
	if transfer.IsBTCTestNet3 {
		addr, err = btcutil.DecodeAddress(toaddress, &chaincfg.TestNet3Params)
	} else {
		addr, err = btcutil.DecodeAddress(toaddress, &chaincfg.MainNetParams)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "转入地址非法",
			"errorinfo": fmt.Sprintf("转入地址非法(%s)", err.Error()),
		})
		return
	}
	_amount := c.PostForm("amount")
	//防止包含小数点，去掉小数点后的数值再转换类型
	if strings.Contains(_amount, ".") {
		a := strings.Split(_amount, ".")
		_amount = a[0]
	}
	amount, err := strconv.Atoi(_amount)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "金额非法",
			"errorinfo": "金额非法",
		})
		glog.Error("金额非法")
		return
	}
	balance, _ := transfer.GetUsdtBalance(fromaddress)
	if balance < int64(amount) {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "余额不足",
			"errorinfo": "余额不足",
		})
		glog.Error(fmt.Sprintf("余额不足(账号:%s,余额:%d,转账金额:%d)", fromaddress, balance, amount))
		return
	}
	txid, err := transfer.SendUsdt(fromaddress, addr, strconv.FormatFloat(float64(float32(amount)/float32(1e8)), 'f', 6, 64))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "交易失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"txid": txid,
		},
	})
	return
}

func TransferFee(c *gin.Context) {
	fromAddress := c.PostForm("fromaddress")
	privKey := c.PostForm("privkey")
	toAddress := c.PostForm("toaddress")
	_amount := c.PostForm("amount")
	send(c, fromAddress, privKey, toAddress, _amount, "btc")
}

func send(c *gin.Context, fromAddress string, privKey string, toAddress string, amount string, coin string) {
	glog.Infof("from:%s,to:%s,prikey:%s,amount:%s,coin:%s\n", fromAddress, toAddress, privKey, amount, coin)
	_, err := btcutil.DecodeAddress(fromAddress, &chaincfg.MainNetParams)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "转出地址非法",
			"errorinfo": fmt.Sprintf("转出地址非法(%s)", err.Error()),
		})
		glog.Error(fmt.Sprintf("转出地址非法:%s", fromAddress))
		return
	}
	_, err = btcutil.DecodeAddress(toAddress, &chaincfg.MainNetParams)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "转入地址非法",
			"errorinfo": fmt.Sprintf("转入地址非法(%s)", err.Error()),
		})
		glog.Error(fmt.Sprintf("转入地址非法:%s", toAddress))
		return
	}
	_amount := amount
	//防止包含小数点，去掉小数点后的数值再转换类型
	if strings.Contains(_amount, ".") {
		a := strings.Split(_amount, ".")
		_amount = a[0]
	}
	a, err := strconv.Atoi(_amount)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "金额非法",
			"errorinfo": fmt.Sprintf("金额非法(%s)", err.Error()),
		})
		glog.Error("金额非法")
		return
	}
	if a < 546 {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "转账金额不能低于546聪",
			"errorinfo": "转账金额不能低于546聪",
		})
		glog.Error("转账金额不能低于546聪")
		return
	}
	balance, _ := transfer.GetBalance(fromAddress, coin)
	if balance < (a) {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "余额不足",
			"errorinfo": "余额不足",
		})
		glog.Error(fmt.Sprintf("余额不足(地址:%s,余额:%d,转账金额:%d)", fromAddress, balance, a))
		return
	}
	txHex, err := transfer.Transaction(fromAddress, toAddress, privKey, a, coin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       fmt.Sprintf("交易失败:%s", err.Error()),
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"txHex": txHex,
		"txUrl": fmt.Sprintf("https://live.blockcypher.com/%s/tx/%s/", coin, txHex),
		"data": gin.H{
			"txHex": txHex,
			"txUrl": fmt.Sprintf("https://live.blockcypher.com/btc/tx/%s/", txHex),
		},
	})
	glog.Info(fmt.Sprintf("txUrl:https://live.blockcypher.com/%s/tx/%s/", coin, txHex))
	return
}

func Balance(c *gin.Context) {
	address := c.Query("address")
	_, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "地址非法",
			"errorinfo": fmt.Sprintf("转出地址非法(%s)", err.Error()),
		})
		return
	}
	balance, err := transfer.GetBalance(address, "btc")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "获取余额失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"coin":    "btc",
			"address": address,
			"balance": balance,
		},
	})
	return
}

func BalanceEth(c *gin.Context) {
	address := c.Query("address")
	if !common.IsHexAddress(address) {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "地址非法",
			"errorinfo": "地址非法",
		})
		return
	}
	balance, err := transfer.GetEthBalance(address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "获取余额失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"coin":    "ETH",
			"address": address,
			"balance": balance.String(),
		},
	})
	return
}

func BalanceEthToken(c *gin.Context) {
	address := c.Query("address")
	tokenAddress := c.Query("tokenaddress")
	if !common.IsHexAddress(address) {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "地址非法",
			"errorinfo": "地址非法",
		})
		return
	}
	if !common.IsHexAddress(tokenAddress) {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "代币地址非法",
			"errorinfo": "代币地址非法",
		})
		return
	}
	balance, err := transfer.GetTokenBalance(tokenAddress, address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "获取余额失败",
			"errorinfo": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"coin":    "ETH",
			"address": address,
			"balance": balance.String(),
		},
	})
	return
}

func SendEthTo(c *gin.Context) {
	fromAddress := c.PostForm("fromaddress")
	privKey := c.PostForm("privkey")
	toAddress := c.PostForm("toaddress")
	amount := c.PostForm("amount")
	sendEthTo(c, fromAddress, privKey, toAddress, amount)
}

func SendEthTokenTo(c *gin.Context) {
	fromAddress := c.PostForm("fromaddress")
	privKey := c.PostForm("privkey")
	toAddress := c.PostForm("toaddress")
	tokenAddress := c.PostForm("tokenaddress")
	amount := c.PostForm("amount")
	sendEthTokenTo(c, fromAddress, privKey, toAddress, tokenAddress, amount)
}

func sendEthTo(c *gin.Context, fromAddress string, privKey string, toAddress string, amount string) {
	glog.Infof("from:%s,to:%s,prikey:%s,amount:%s\n", fromAddress, toAddress, privKey, amount)
	if !common.IsHexAddress(fromAddress) {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "转出地址非法",
			"errorinfo": "转出地址非法",
		})
		glog.Error(fmt.Sprintf("转出地址非法:%s", fromAddress))
		return
	}
	if !common.IsHexAddress(toAddress) {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "转入地址非法",
			"errorinfo": "转入地址非法",
		})
		glog.Error(fmt.Sprintf("转入地址非法:%s", toAddress))
		return
	}
	_amount := amount
	//防止包含小数点，去掉小数点后的数值再转换类型
	if strings.Contains(_amount, ".") {
		a := strings.Split(_amount, ".")
		_amount = a[0]
	}
	var a big.Int
	_, b := a.SetString(_amount, 10)
	if !b {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "金额非法",
			"errorinfo": "金额非法",
		})
		glog.Error("金额非法")
		return
	}
	balance, _ := transfer.GetEthBalance(fromAddress)
	if balance.Cmp(&a) == -1 {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "余额不足",
			"errorinfo": "余额不足",
		})
		glog.Error(fmt.Sprintf("余额不足(地址:%s,余额:%d,转账金额:%s)", fromAddress, balance.Int64(), a.String()))
		return
	}
	txHex, err := transfer.TransactionEth(fromAddress, toAddress, privKey, a)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "交易失败",
			"errorinfo": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"txHex": txHex,
		"txUrl": fmt.Sprintf("https://etherscan.io/tx/%s/", txHex),
		"data": gin.H{
			"txHex": txHex,
			"txUrl": fmt.Sprintf("https://etherscan.io/tx/%s/", txHex),
		},
	})
	glog.Info(fmt.Sprintf("txUrl:https://etherscan.io/tx/%s/", txHex))
	return
}

func sendEthTokenTo(c *gin.Context, fromAddress string, privKey string, toAddress string, tokenAddress string, amount string) {
	glog.Infof("from:%s,to:%s,prikey:%s,amount:%s,tokenAddress:%s\n", fromAddress, toAddress, privKey, amount, tokenAddress)
	if !common.IsHexAddress(fromAddress) {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "转出地址非法",
			"errorinfo": "转出地址非法",
		})
		glog.Error(fmt.Sprintf("转出地址非法:%s", fromAddress))
		return
	}
	if !common.IsHexAddress(toAddress) {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "转入地址非法",
			"errorinfo": "转入地址非法",
		})
		glog.Error(fmt.Sprintf("转入地址非法:%s", toAddress))
		return
	}
	if !common.IsHexAddress(tokenAddress) {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "代币地址非法",
			"errorinfo": "代币地址非法",
		})
		glog.Error(fmt.Sprintf("代币地址非法:%s", tokenAddress))
		return
	}
	_amount := amount
	//防止包含小数点，去掉小数点后的数值再转换类型
	if strings.Contains(_amount, ".") {
		a := strings.Split(_amount, ".")
		_amount = a[0]
	}
	var a big.Int
	_, b := a.SetString(_amount, 10)
	if !b {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "金额非法",
			"errorinfo": "金额非法",
		})
		glog.Error("金额非法")
		return
	}
	balance, _ := transfer.GetTokenBalance(tokenAddress, fromAddress)
	if balance.Cmp(&a) == -1 {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "余额不足",
			"errorinfo": "余额不足",
		})
		glog.Error(fmt.Sprintf("余额不足(地址:%s,余额:%d,转账金额:%s)", fromAddress, balance.Int64(), a.String()))
		return
	}
	txHex, err := transfer.TransactionToken(fromAddress, toAddress, tokenAddress, privKey, a)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "交易失败",
			"errorinfo": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"txHex": txHex,
			"txUrl": fmt.Sprintf("https://etherscan.io/tx/%s/", txHex),
		},
	})
	glog.Info(fmt.Sprintf("txUrl:https://etherscan.io/tx/%s/", txHex))
	return
}

//用一套助记词，根据秘钥(项目名称)生成seed 根据userid生成私钥和地址
func GenerateAddress(c *gin.Context) {
	_userid := c.PostForm("userid")
	userid, err := strconv.Atoi(_userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "userid非法",
			"errorinfo": err.Error(),
		})
		glog.Error("userid非法 ", _userid)
		return
	}
	mnemonic := "alien ghost shock story abuse wish garlic chuckle trophy army silver category"
	p := fmt.Sprintf("m/44'/60'/%d'/0/0", userid)
	seed := bip39.NewSeed(mnemonic, "buzhidao")
	/*---------- 开始生成Eth地址 ----------*/
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "生成以太坊钱包失败",
			"errorinfo": err.Error(),
		})
		return
	}
	path := hdwallet.MustParseDerivationPath(p)
	account, err := wallet.Derive(path, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -3,
			"msg":       "生成以太坊钱包失败",
			"errorinfo": err.Error(),
		})
		return
	}
	privatekeyHex, err := wallet.PrivateKeyHex(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -4,
			"msg":       "生成以太坊钱包失败",
			"errorinfo": err.Error(),
		})
		return
	}
	/*---------- Eth地址生成完毕 ----------*/

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"eth_private": privatekeyHex,
			"eth_address": account.Address.Hex(),
		},
	})
	return
}
