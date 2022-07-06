package models

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

func (nft NftDb) LoginNew(userAddr, sigData string) error {
	userAddr = strings.ToLower(userAddr)

	user := Users{}
	db := nft.db.Model(&user).Where("useraddr = ?", userAddr).First(&user)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			user.Useraddr = userAddr
			user.Signdata = sigData
			user.Userlogin = time.Now().Unix()
			user.Userlogout = time.Now().Unix()
			user.Username = ""
			user.Userregd = time.Now().Unix()
			db = nft.db.Model(&user).Create(&user)
			if db.Error != nil {
				fmt.Println("loging()->create() err=", db.Error)
				return db.Error
			}
		}
	} else {
		db = nft.db.Model(&Users{}).Where("useraddr = ?", userAddr).Update("userlogin", time.Now().Unix())
		if db.Error != nil {
			fmt.Printf("login()->UPdate() users err=%s\n", db.Error)
		}
	}
	return db.Error
}

func (nft NftDb) Login(userAddr, sigData string) error {
	userAddr = strings.ToLower(userAddr)

	userOld := Users{}
	db := nft.db.Model(&Users{}).Last(&userOld)
	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		fmt.Println("login() look last err=", db.Error)
		return db.Error
	}
	fmt.Println("login()", "userOld.id= ", userOld.ID)
	user := Users{}
	db = nft.db.Model(&user).Where("useraddr = ?", userAddr).First(&user)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			fmt.Println("log()", "userOld.id= ", userOld.ID)
			user.ID = userOld.ID + 1
			user.Useraddr = userAddr
			user.Signdata = sigData
			user.Userlogin = time.Now().Unix()
			user.Userlogout = time.Now().Unix()
			user.Username = ""
			user.Userregd = time.Now().Unix()
			db = nft.db.Model(&user).Create(&user)
			if db.Error != nil {
				fmt.Println("\"log num=\", num, loging()->create() err=", db.Error)
				return db.Error
			}
			fmt.Println("log()", "user.id= ", user.ID)
			fmt.Println("log()", "userOld.id= ", userOld.ID)
		}
	} else {
		fmt.Println("log()", "find user.id= ", user.ID)
		db = nft.db.Model(&Users{}).Where("useraddr = ?", userAddr).Update("userlogin", time.Now().Unix())
		if db.Error != nil {
			fmt.Printf("login()->UPdate() users err=%s\n", db.Error)
		}
	}
	return db.Error
}
