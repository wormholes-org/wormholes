package nftexchangev2

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/nftexchange/nftserver/controllers"
	"math/big"
	"strconv"
	"sync"
	"time"
)

const (
	DeleteMapList = time.Second * 1
	PrintMapListLog = time.Minute * 10
	MapListItemsNum = 59
	ResultScope = 10000
	EXPIREDURATION = 60 * 60 * 10
)

var answerMapList AnswerMapList
var tokenMap TokenMap
var approveAddrsMap ApproveAddrsMap
func init() {
	answerMapList.MapList = make(map[int64]map[string]int64)
	go answerMapList.DeleteLoop()

	tokenMap.Tokens = make(map[string]TokenInfo)
	go tokenMap.DeleteTokensLoop()

	approveAddrsMap.ApproveAddrs = make(map[string]AddrInfo)
}

type AnswerMapList struct {
	Mux sync.Mutex
	MapList map[int64]map[string]int64
}

func (mapList *AnswerMapList) GetHashAndSecret() (string, string, int64, error) {
	var hashString string
	//hash := sha256.New()
	secret := GenerateString()
	result, err := rand.Int(rand.Reader, big.NewInt(ResultScope))
	if err == nil {
		rawData := secret + strconv.Itoa(int(result.Int64()))
		sum := md5.Sum([]byte(rawData))
		hashString = hex.EncodeToString(sum[:])
	} else {
		return "", "", 0, err
	}

	return hashString, secret, result.Int64(), nil
}

func (mapList *AnswerMapList) AddResult(userAddr string, result int64) int64 {
	var timeKey int64
	timeKey = time.Now().Unix()
	mapList.Mux.Lock()
	if len(mapList.MapList[timeKey]) == 0 {
		answerMap := make(map[string]int64)
		mapList.MapList[timeKey] = answerMap
	}

	mapList.MapList[timeKey][userAddr] = result
	mapList.Mux.Unlock()

	return timeKey
}

func (mapList *AnswerMapList) deleteAnswerMap(timeKey int64) {
	mapList.Mux.Lock()
	delete(mapList.MapList, timeKey)
	mapList.Mux.Unlock()
}

func (mapList *AnswerMapList) IsValid(timeStamp int64, userAddr string, result int64) bool {
	value, ok := mapList.MapList[timeStamp][userAddr]
	if ok && value == result {
		return true
	} else {
		return false
	}
}

func (mapList *AnswerMapList) DeleteLoop() {
	tick := time.Tick(DeleteMapList)
	tickLog := time.Tick(PrintMapListLog)
	for {
		select {
		case <- tick:
			deleteKey := time.Now().Unix() - MapListItemsNum
			mapList.Mux.Lock()
			for k, _ := range mapList.MapList {
				if k < deleteKey {
					delete(mapList.MapList, k)
				}
			}
			mapList.Mux.Unlock()
		case <- tickLog:
			fmt.Println("AnswerMapList len =", len(mapList.MapList))
			var total int
			for _, v := range mapList.MapList {
				total = total + len(v)
			}
			fmt.Println("AnswerMapList.MapList total items len =", total)
		}
	}
}

type TokenMap struct {
	Mux sync.Mutex
	Tokens map[string]TokenInfo
}
type TokenInfo struct {
	Token string
	TimeStamp int64
}

func (tokenMap *TokenMap) AddToken(userAddr string, token string) {
	tokenMap.Mux.Lock()
	tokenInfo := TokenInfo{
		Token: token,
		TimeStamp: time.Now().Unix(),
	}
	tokenMap.Tokens[userAddr] = tokenInfo
	tokenMap.Mux.Unlock()
}

func (tokenMap *TokenMap) GetToken(userAddr string) (string, error) {
	tokenInfo, ok := tokenMap.Tokens[userAddr]
	if ok {
		return tokenInfo.Token, nil
	}
	return "", errors.New("token is not exist!")
}

func (tokenMap *TokenMap) IsValidToken(tokenString string) bool {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&controllers.ExchangerCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(controllers.PRIVATEKEY), nil
		})
	if err != nil {
		return false
	}
	_, ok := token.Claims.(*controllers.ExchangerCustomClaims)
	if !ok || !token.Valid {
		return false
	}

	return true
}

func (tokenMap *TokenMap) DeleteExpireTokens() {
	deleteTokenKeys := make([]string, 0)
	nowStamp := time.Now().Unix()

	tokenMap.Mux.Lock()
	for tokenKey, tokenInfo := range tokenMap.Tokens {
		//validToken := tokenMap.IsValidToken(tokenInfo.Token)
		//if !validToken {
		//	deleteTokenKeys = append(deleteTokenKeys, tokenKey)
		//}
		if nowStamp - tokenInfo.TimeStamp > EXPIREDURATION {
			deleteTokenKeys = append(deleteTokenKeys, tokenKey)
		}
	}


	for _, v := range deleteTokenKeys {
		delete(tokenMap.Tokens, v)
	}
	tokenMap.Mux.Unlock()
}

func (tokenMap *TokenMap) DeleteTokensLoop() {
	tick := time.Tick(time.Second * 10)
	for {
		select {
		case <- tick:
			tokenMap.DeleteExpireTokens()
		}
	}
}


type ApproveAddrsMap struct {
	Mux sync.Mutex
	ApproveAddrs map[string]AddrInfo
}
type AddrInfo struct {
	Addr string
	TimeStamp int64
}

func (approveAddrsMap *ApproveAddrsMap) AddApproveAddr(userAddr string, approveAddr string) {
	approveAddrsMap.Mux.Lock()
	addrInfo := AddrInfo{
		Addr: approveAddr,
		TimeStamp: time.Now().Unix(),
	}
	approveAddrsMap.ApproveAddrs[userAddr] = addrInfo
	approveAddrsMap.Mux.Unlock()
}

func(approveAddrsMap *ApproveAddrsMap) GetApproveAddr(userAddr string) (string, error) {
	approveAddrInfo, ok := approveAddrsMap.ApproveAddrs[userAddr]
	if ok {
		return approveAddrInfo.Addr, nil
	}
	return "", errors.New("Not Exist!")
}


