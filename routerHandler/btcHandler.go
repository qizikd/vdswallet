package routerHandler

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/qizikd/vdswallet/transfer"
	"net/http"
	"strconv"
)

func GetNewAddress(c *gin.Context) {
	address, err := transfer.GetNewAddress()
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

func GetBalance(c *gin.Context) {
	balance, err := transfer.GetWalletBalance()
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

func GetTransactions(c *gin.Context) {
	start := c.Query("start")
	end := c.PostForm("end")
	_start, err := strconv.Atoi(start)
	if err != nil{
		_start = 0
	}
	_end, err := strconv.Atoi(end)
	if err != nil{
		_end = 100
	}
	result, err := transfer.ListWalletTransactions(_start,_end)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
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

func GetBtcAddressReceive(c *gin.Context) {
	address := c.Query("address")
	var addr btcutil.Address
	var err error
	addr, err = btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
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

func SendBtc(c *gin.Context) {
	toaddress := c.PostForm("toaddress")
	_amount := c.PostForm("amount")
	amount, err := strconv.ParseFloat(_amount,64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      -1,
			"msg":       "金额非法",
			"errorinfo": "金额非法",
		})
		glog.Error("金额非法")
		return
	}
	balance, _ := transfer.GetWalletBalance()
	if balance < amount {
		c.JSON(http.StatusOK, gin.H{
			"code":      -2,
			"msg":       "余额不足",
			"errorinfo": "余额不足",
		})
		glog.Error(fmt.Sprintf("余额不足(账号:%s,余额:%d,转账金额:%d)", balance, amount))
		return
	}
	txid, err := transfer.SendBtc(toaddress, amount)
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

