package routerHandler

import (
	"fmt"
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
	offset := c.Query("offset")
	count := c.Query("count")
	_offset, err := strconv.Atoi(offset)
	if err != nil {
		_offset = 0
	}
	_count, err := strconv.Atoi(count)
	if err != nil {
		_count = 100
	}
	result, err := transfer.ListWalletTransactions(_offset, _count)
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

func SendBtc(c *gin.Context) {
	toaddress := c.PostForm("toaddress")
	_amount := c.PostForm("amount")
	amount, err := strconv.ParseFloat(_amount, 64)
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
