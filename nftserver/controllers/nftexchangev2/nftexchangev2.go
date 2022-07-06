package nftexchangev2

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/nftexchange/nftserver/controllers"
	"math/big"
	"time"
)

const (
	PattenString = "^[0-9a-zA-Z_]+$"
	PattenNumber = "^[0-9]+$"
	PattenHex = "^[0-9a-fA-F]+$"
	PattenOperator = "^[<>=]+$"
	PattenEmail = "^[A-Za-z0-9]+([-_.][A-Za-z0-9]+)*@([A-Za-z0-9]+[-.])+[A-Za-z0-9]{2,4}$"
	PattenAddr = "^0x[0-9a-fA-F]{40}$"
	PattenImageBase64 = "^data:image/[a-zA-Z0-9]+;base64,[a-zA-Z0-9/+]+=?=?$"
)

var (
	ERRINPUTINVALID = errors.New("input data invalid")
	ERRTOKEN = errors.New("token invalid, please relogin!")
)


type NftExchangeControllerV2 struct {
	beego.Controller
}

//func (nft *NftExchangeControllerV2) GetImageFromIPFS() {
//	var data map[string]interface{}
//	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
//	defer nft.Ctx.Request.Body.Close()
//	json.Unmarshal(bytes, &data)
//	s, ok := data["hash"].(string)
//	fmt.Printf(">>>>>>>>s=%s\n", s)
//	if ok {
//
//	}
//}


//owner修改价格
//func (nft *NftExchangeControllerV2) ModifyNFT() {
//	fmt.Println("ModifyNFT()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
//	var httpResponseData HttpResponseData
//	nd := new(NftDb)
//	err := nd.ConnectDB(sqldsndb)
//	if err != nil {
//		fmt.Printf("connect database err = %s\n", err)
//	}
//	defer nd.Close()
//
//	var data map[string]string
//	defer nft.Ctx.Request.Body.Close()
//	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
//	//fmt.Printf("receive data = %s\n", string(bytes))
//	err = json.Unmarshal(bytes, &data)
//	if err != nil {
//		httpResponseData.Code = "500"
//		httpResponseData.Msg = err.Error()
//		httpResponseData.Data = []interface{}{}
//	} else {
//		rawData := nft.removeSignData(string(bytes))
//		_, err := nft.isValidAddr(rawData, data["sig"], data["user_addr"])
//		if err != nil {
//			httpResponseData.Code = "500"
//			httpResponseData.Msg = err.Error()
//			httpResponseData.Data = []interface{}{}
//		} else {
//			err = nd.InsertSigData(data["sig"], rawData)
//			if err != nil {
//				httpResponseData.Code = "500"
//				httpResponseData.Msg = err.Error()
//			} else {
//				//modify the database value if the verification address is valid.
//				price1, _ := strconv.ParseUint(data["price1"], 10, 64)
//				price2, _ := strconv.ParseUint(data["price2"], 10, 64)
//				//startTime, _ := time.Parse("2006-01-02 15:04:05", data["time1"])
//				//closeTime, _ := time.Parse("2006-01-02 15:04:05", data["time2"])
//				startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", data["time1"], time.Local)
//				closeTime, _ := time.ParseInLocation("2006-01-02 15:04:05", data["time2"], time.Local)
//				err = nd.ModifyPrice(data["user_addr"], data["nft_contract_addr"],
//					data["nft_token_id"], data["sell_type"], price1, price2,
//					data["currency"], data["hide"], data["trade_sig"], data["sig"], startTime, closeTime)
//				if err != nil {
//					httpResponseData.Code = "500"
//					httpResponseData.Msg = err.Error()
//					httpResponseData.Data = []interface{}{}
//				} else {
//					httpResponseData.Code = "200"
//					httpResponseData.Data = []interface{}{}
//				}
//			}
//		}
//	}
//
//	responseData, _ := json.Marshal(httpResponseData)
//	nft.Ctx.ResponseWriter.Write(responseData)
//	fmt.Println("ModifyNFT()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
//}

func RefreshToken(tokenString string) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString,
		&controllers.ExchangerCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(controllers.PRIVATEKEY), nil
		})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*controllers.ExchangerCustomClaims)
	if !ok || !token.Valid {
		return "", errors.New("token is invalid!")
	}
	signKey := []byte(controllers.PRIVATEKEY)
	expireAt := time.Now().Add(time.Second * time.Duration(controllers.DEFAULT_EXPIRE_TIME_SECONDS)).Unix()
	newClaims := controllers.ExchangerCustomClaims{
		claims.User,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			IssuedAt: time.Now().Unix(),
			Issuer: "Wormholes Exchanger",
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	newTokenString, err := newToken.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return newTokenString, nil
}

func ValidateToken(tokenString string) bool {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&controllers.ExchangerCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(controllers.PRIVATEKEY), nil
		})
	if err != nil {
		return false
	}
	claims, ok := token.Claims.(*controllers.ExchangerCustomClaims)
	if !ok || !token.Valid {
		return false
	}

	fmt.Println("claims = ", claims)
	return true
}

func GenerateToken(userAddr string) (string, error) {
	signKey := []byte(controllers.PRIVATEKEY)
	expireAt := time.Now().Add(time.Second * time.Duration(controllers.DEFAULT_EXPIRE_TIME_SECONDS)).Unix()
	claims := controllers.ExchangerCustomClaims{
		controllers.User{UserAddr: userAddr},
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			IssuedAt: time.Now().Unix(),
			Issuer: "Wormholes Exchanger",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CheckToken(ctx *context.Context) {
	var httpResponseData controllers.HttpResponseData
	token := ctx.Request.Header.Get("Token")
	//fmt.Println("X-Real-IP=", ctx.Request.Header.Get("X-Real-IP"))
	//fmt.Println("X-Forward-For=", ctx.Request.Header.Get("X-Forward-For"))
	//fmt.Println("RemoteAddr=", ctx.Request.RemoteAddr)
	//fmt.Println("ctx.Request.Host=", ctx.Request.Host)
	isValidToken := ValidateToken(token)
	if !isValidToken {
		httpResponseData.Code = "401"
		httpResponseData.Data = []interface{}{}
		httpResponseData.Msg = "token invalid!"
		responseData, _ := json.Marshal(httpResponseData)
		ctx.ResponseWriter.Write(responseData)
	}
}

func GenerateString() string {
	var strLength = 4
	var randStr string
	characterSet := []string {
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
		"u", "v", "w", "x", "y", "z",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
		"u", "v", "w", "x", "y", "z",
	}

	for i:=0; i<strLength; i++ {
		key, err := rand.Int(rand.Reader, big.NewInt(62))
		if err != nil {
			return ""
		}
		randStr = randStr + characterSet[key.Int64()]
	}

	return randStr
}





