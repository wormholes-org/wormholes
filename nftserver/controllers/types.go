package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nftexchange/nftserver/models"
)

// API访问时token所需要的设置
const (
	PRIVATEKEY = "NFTEXCHANGER.WORMHOLES.202110191729"
	DEFAULT_EXPIRE_TIME_SECONDS = 60 * 60 * 10
)
type User struct {
	UserAddr string
}
type ExchangerCustomClaims struct {
	User
	jwt.StandardClaims
}

//Server response
type HttpResponseData struct {
	Code string			`json:"code"`
	Msg string			`json:"msg"`
	Data interface{}	`json:"data"`
	TotalCount uint64	`json:"total_count"`
}

// 为方便解析filter字段，定义了此结构体来解析http请求数据，
// 没有使用通用的方式map[string]string
type HttpRequestFilter struct {
	Match string `json:"match"`
	Filter []models.StQueryField `json:"filter"`
	Sort []models.StSortField `json:"sort"`
	StartIndex string `json:"start_index"`
	Count string `json:"count"`
}


