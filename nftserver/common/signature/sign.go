package signature

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nftexchange/nftserver/ethhelper"
	"golang.org/x/crypto/sha3"
	"os"
	"strings"
	"time"
)

func TextAndHash(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}

func GetEthAddr(msg string, sigStr string) (common.Address, error){
	sigData := hexutil.MustDecode(sigStr)
	if len(sigData) != 65 {
		return common.Address{}, fmt.Errorf("signature must be 65 bytes long")
	}
	if sigData[64] != 27 && sigData[64] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sigData[64] -= 27
	hash, _ := TextAndHash([]byte(msg))
	fmt.Println("sigdebug hash=", hexutil.Encode(hash))
	rpk, err := crypto.SigToPub(hash, sigData)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*rpk), nil
}

func IsValidVerifyAddr(rawData string, sig string) (bool, error) {
	addrList, err := ethhelper.AdminList()
	if err != nil {
		return false, err
	}

	verificationAddr, err := GetEthAddr(rawData, sig)
	if err != nil {
		return false, err
	}
	verificationAddrS := verificationAddr.String()
	verificationAddrS = strings.ToLower(verificationAddrS)

	for _, addr := range addrList {
		if !strings.HasPrefix(addr, "0x") {
			addr = "0x" + addr
		}
		addr = strings.ToLower(addr)
		fmt.Printf("sigdebug verificationAddrS = [%s], Manager's addr = [%s]\n", verificationAddrS, addr)
		if verificationAddrS == addr {
			fmt.Println("sigdebug verify [Y]")
			return true, nil
		}
	}
	fmt.Println("sigdebug verify [N]")
	//return true, nil

	return false, errors.New("verification address is invalid")
}

func IsValidAddr(
	rawData string,
	sig string,
	addr string) (bool, error) {
	verificationAddr, err := GetEthAddr(rawData, sig)
	if err != nil {
		return false, err
	}
	verificationAddrS := verificationAddr.String()
	verificationAddrS = strings.ToLower(verificationAddrS)

	addr = strings.ToLower(addr)
	fmt.Printf("sigdebug verificationAddrS = [%s], approveAddr's addr = [%s]\n", verificationAddrS, addr)
	if verificationAddrS == addr {
		fmt.Println("sigdebug verify [Y]")
		return true, nil
	}
	fmt.Println("sigdebug verify [N]")
	//return true, nil

	return false, errors.New("address is invalid")
}

func RemoveSignData(jsonDataS string) string {
	lastIndex1 := strings.LastIndex(jsonDataS, ",")
	lastIndex2 := strings.LastIndex(jsonDataS, "}")
	newJsonDataS := string([]byte(jsonDataS)[:lastIndex1]) + string([]byte(jsonDataS)[lastIndex2:])
	fmt.Printf("jsonDataS=%s\n", jsonDataS)
	fmt.Printf("newJsonDataS=%s\n", newJsonDataS)
	return newJsonDataS
}

func VerifyAppconf(filePath string, oldAddr string) bool {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return false
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(bytesread)
	msg := string(buffer)
	index := strings.Index(msg, "app.conf.sig = ")
	sig := msg[index + len("app.conf.sig = "):]
	var message []byte = []byte(msg[:strings.Index(msg, "#签名数据")])
	addr, err := GetEthAddr(string(message), sig)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if addr.String() != oldAddr {
		return false
	}
	fmt.Println(addr)
	return true
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func SignAppconf(filePath string) {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bytesread)
	msg := string(buffer)
	tm := fmt.Sprintf(time.Now().String())
	msg = msg + "[time]\n" + "date = " + tm + "\n\n"

	var message []byte = []byte(msg)
	key, err := crypto.HexToECDSA("8c995fd78bddf528bd548cce025f62d4c3c0658362dbfd31b23414cf7ce2e8ed")
	if err != nil {
		fmt.Println(err)
	}
	sig, err := crypto.Sign(signHash(message), key)
	if err != nil {
		fmt.Println("signature error: %s", err)
	}
	sig[64] += 27
	sigstr := hexutil.Encode(sig)
	msg = msg + "#签名数据\n" + "[sig]\n" + "app.conf.sig = " + sigstr
	_, err = file.WriteAt([]byte(msg), 0)
	if err != nil {
		fmt.Println(err)
	}
}
