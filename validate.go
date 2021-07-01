package main

import (
	"errors"
	"net/http"
	"nyh00-product/common"
	"nyh00-product/encrypt"
)

func Check(w http.ResponseWriter, r *http.Request) {
	// 执行正常业务逻辑
	println("执行checck")
}

// 统一验证拦截器，每个接口都需要提前验证
func Auth(rw http.ResponseWriter, r *http.Request) error {
	println("执行验证")
	// 添加基于cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

// 身份校验函数
func CheckUserInfo(r *http.Request) error {
	// 1.获取cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户id Cookie获取失败！")
	}
	// 获取用户加密串
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户加密串 Cookie获取失败！")
	}
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串已被篡改！")
	}
	println("结果比对")
	println("用户ID：" + uidCookie.Value)
	println("解密后用户ID：" + string(signByte))
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	return errors.New("身份校验失败")
}

func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}

func main() {
	// 过滤器
	filter := common.NewFilter()
	// 注册拦截器
	filter.RegisterFilterUri("/check", Auth)
	// 启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	// 启动端口
	http.ListenAndServe("8083", nil)
}
