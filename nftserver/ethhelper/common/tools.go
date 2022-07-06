package common

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

const DEFAULT_MAX_WIDTH uint = 500
const DEFAULT_MAX_HEIGHT uint = 500

func GetImgBase64(imgUrl string) (string, error) {
	//imgPath := "/home/user1/chengdu"
	imgPath := "D:\\"
	fileName := path.Base(imgUrl)
	res, err := http.Get(imgUrl)
	if err != nil {
		fmt.Println("A error occurred!", err)
		return "", errors.New("url not correct:" + imgUrl)
	}

	// 获得get请求响应的reader对象

	//noEnd := false
	//if fileName[len(fileName)-1:] != "g" {
	//	fileName+=".png"
	//	//noEnd = true
	//}
	reader := bufio.NewReaderSize(res.Body, 64*1024)
	file, err := os.Create(imgPath + fileName)
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	_, err = io.Copy(writer, reader)

	defer func() {
		os.Remove(imgPath + fileName)
		os.Remove(imgPath + "comp_" + fileName)
		res.Body.Close()
	}()
	return makeThumbnail(imgPath+fileName, imgPath+"comp_"+fileName)

}

//// 计算图片缩放后的尺寸
//func calculateRatioFit(srcWidth, srcHeight int) (int, int) {
//	ratio := math.Min(DEFAULT_MAX_WIDTH/float64(srcWidth), DEFAULT_MAX_HEIGHT/float64(srcHeight))
//	return int(math.Ceil(float64(srcWidth) * ratio)), int(math.Ceil(float64(srcHeight) * ratio))
//}

// 生成缩略图
func makeThumbnail(imagePath, savePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", errors.New("makeThumbnail:" + err.Error())
	}
	defer file.Close()
	var img image.Image
	var base string
	if imagePath[len(imagePath)-2:] == "NG" || imagePath[len(imagePath)-2:] == "ng" {
		base = "data:image/png;base64,"
		img, err = png.Decode(file)
	} else {
		img, err = jpeg.Decode(file)
		base = "data:image/jpg;base64,"
	}
	if err != nil {
		return "", errors.New("makeThumbnail:" + err.Error())
	}

	if err != nil {
		return "", errors.New("makeThumbnail file invalid :" + err.Error())
	}

	//b := img.Bounds()
	//width := b.Max.X
	//height := b.Max.Y

	//w, h := calculateRatioFit(width, height)
	m := resize.Resize(DEFAULT_MAX_WIDTH, DEFAULT_MAX_HEIGHT, img, resize.Lanczos3)

	imgfile, _ := os.Create(savePath)

	defer func() {
		imgfile.Close()
	}()

	err = png.Encode(imgfile, m)
	data, err := ToBase64(savePath)
	return base + data, err
}
func ToBase64(path string) (string, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	defer imgFile.Close()
	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)

	return imgBase64Str, nil
}
