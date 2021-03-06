package model

import (
	"crypto/md5"
	"db"
	"encoding/hex"
	"errors"
	"fmt"
)

func Login(account string, password string) {

	//定义一个登录的状态.
	//这里定义true，表示登录成功.
	status := true

	//验证登录的帐号.
	ac := CheckAccount(account)

	//如果验证用户名的方法传入过来的错误信息不为空，就表示用户名存在错误.
	//就需要返回错误状态.
	if ac != nil {
		status = false
	}

	//验证登录的密码.
	pw := CheckPassword(password)

	//验证登录密码是否返回的有错误信息.
	if pw != nil {
		status = false
	}

	//验证帐号密码的有效性.
	//根据account 来 查询 存放在数据库中加密过后的密码.
	//将登录页面上填写的密码进行加密与account所对应的密码进行对比.
	//如果对比成功则表示用户是真实的,否则表示登录失败.
	status = CheckLogin(account, password)

	fmt.Println(status)
	//	return status, msg

	//根据account来查询密码然后对比加密后的密码与传入过来的密码是否相同.

	//判断帐号和密码是否为空.

}

//判断帐号的长度是否为空等等.
//由于登录是可以使用用户名，邮箱，和手机号，所以不能验证用户名的长度，因为用户名可能为一个字符.
func CheckAccount(account string) error {

	err := errors.New("用户名不能为空")

	if len(account) != 0 {
		err = nil
	}
	return err
}

//判断登录的密码长度等等.
func CheckPassword(password string) error {

	err := errors.New("密码长度不够")

	if len(password) > 6 {
		err = nil
	}
	return err
}

//验证帐号密码的有效性.
//根据account 来 查询 存放在数据库中加密过后的密码.
//将登录页面上填写的密码进行加密与account所对应的密码进行对比.
//如果对比成功则表示用户是真实的,否则表示登录失败.
func CheckLogin(account string, password string) bool {

	//定义一个状态用来存储默认登录成功的状态 true.
	var loginStatus bool = true

	type Manage struct {
		id       int
		username string
		password string
		email    string
		mobile   int
	}
	//定义sql.
	sql := "SELECT `id`,`username`,`password`,`email`,`mobile` FROM `manage` WHERE `email` = ? OR `mobile` = ? OR `username` = ?"

	//获取db对象.
	obj := db.GetInstance()

	//	var username string
	//	var id int
	manage := &Manage{}

	row := obj.QueryRow(sql, account, account, account)

	//Scan方法是将查询出来的数据填充到各个指定的各个值中.
	//这里我们填充到一个结构体中.
	err := row.Scan(&manage.id, &manage.username, &manage.password, &manage.email, &manage.mobile)

	if err != nil {

		//如果有错误直接返回登录失败.
		return loginStatus
	}

	//对密码进行加密.
	//返回一个新的哈希散列来计算md5.
	encryption := md5.New()
	encryption.Write([]byte(password))
	password = hex.EncodeToString(encryption.Sum(nil))

	//判断查询出来的密码与登陆时候填写的密码是否一致.
	//这里登录页面传递过来的密码需要进行加密处理以后才能与查询出来的密码进行对比.
	if manage.password != password {

		//密码不相同就修改前面定义的默认登录状态.
		//将状态修改成false.
		loginStatus = false
	}

	return loginStatus
}
