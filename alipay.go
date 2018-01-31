package main

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/smartwalle/alipay/encoding"
)

const (
	K_TIME_FORMAT = "2006-01-02 15:04:05"

	K_ALI_PAY_TRADE_STATUS_WAIT_BUYER_PAY = "WAIT_BUYER_PAY" // 交易创建，等待买家付款
	K_ALI_PAY_TRADE_STATUS_TRADE_CLOSED   = "TRADE_CLOSED"   // 未付款交易超时关闭，或支付完成后全额退款
	K_ALI_PAY_TRADE_STATUS_TRADE_SUCCESS  = "TRADE_SUCCESS"  // 交易支付成功
	K_ALI_PAY_TRADE_STATUS_TRADE_FINISHED = "TRADE_FINISHED" // 交易结束，不可退款

	K_ALI_PAY_SANDBOX_API_URL    = "https://openapi.alipaydev.com/gateway.do"
	K_ALI_PAY_PRODUCTION_API_URL = "https://openapi.alipay.com/gateway.do"

	K_ALI_PAY_SUCCESS = "1"
	K_ALI_PAY_FAILED  = "0"

	K_FORMAT       = "JSON"
	K_CHARSET      = "utf-8"
	K_VERSION      = "1.0"
	K_PRODUCT_CODE = "QUICK_WAP_WAY"

	// https://doc.open.alipay.com/docs/doc.htm?treeId=291&articleId=105806&docType=1
	K_SUCCESS_CODE = "10000"

	k_RESPONSE_SUFFIX = "_response"
	k_ERROR_RESPONSE  = "error_response"
	k_SIGN_NODE_NAME  = "sign"

	K_SIGN_TYPE_RSA2 = "RSA2"
	K_SIGN_TYPE_RSA  = "RSA"
)

type AliPay struct {
	appId           string
	apiDomain       string
	partnerId       string
	publicKey       []byte
	privateKey      []byte
	AliPayPublicKey []byte
	client          *http.Client
	SignType        string
}
type AliPayParam interface {
	// 用于提供访问的 method
	APIName() string

	// 返回参数列表
	Params() map[string]string

	// 返回扩展 JSON 参数的字段名称
	ExtJSONParamName() string

	// 返回扩展 JSON 参数的字段值
	ExtJSONParamValue() string
}

type Result struct {
	// 状态
	Status int
	// 本网站订单号
	OrderNo string
	// 支付宝交易号
	TradeNo string
	// 买家支付宝账号
	BuyerEmail string
	// 错误提示
	Message string

	GmtPayment string

	TotalFee string

	Extra_common_param string
}

func AlipayNew(appId, partnerId string, publicKey, privateKey []byte, isProduction bool) (client *AliPay) {
	client = &AliPay{}
	client.appId = appId
	client.partnerId = partnerId
	client.privateKey = privateKey
	client.publicKey = publicKey
	client.client = http.DefaultClient
	if isProduction {
		client.apiDomain = K_ALI_PAY_PRODUCTION_API_URL
	} else {
		client.apiDomain = K_ALI_PAY_SANDBOX_API_URL
	}
	client.SignType = K_SIGN_TYPE_RSA
	return client
}

func signRSA2(keys []string, param map[string]string, privateKey []byte) (s string, err error) {
	if param == nil {
		param = make(map[string]string, 0)
	}

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(param[key])
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var src = strings.Join(pList, "&")
	sig, err := encoding.SignPKCS1v15([]byte(src), privateKey, crypto.SHA256)
	if err != nil {
		return "", err
	}
	s = base64.StdEncoding.EncodeToString(sig)
	return s, nil
}

func signRSA(keys []string, param map[string]string, privateKey []byte) (s string, err error) {
	if param == nil {
		param = make(map[string]string, 0)
	}

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(param[key])
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var src = strings.Join(pList, "&")
	//seelog.Info(src)
	sig, err := encoding.SignPKCS1v15([]byte(src), privateKey, crypto.SHA1)
	if err != nil {
		return "", err
	}
	s = base64.StdEncoding.EncodeToString(sig)
	return s, nil
}

func verifySign(param map[string]string, publicKey []byte) (ok bool, err error) {
	sign, err := base64.StdEncoding.DecodeString(param["sign"])
	fmt.Println(sign, param["sign"])
	signType := param["sign_type"]
	if err != nil {
		return false, err
	}

	var keys = make([]string, 0, 0)
	for key, value := range param {
		if key == "sign" || key == "sign_type" {
			continue
		}
		if len(value) > 0 {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(param[key])
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var s = strings.Join(pList, "&")
	fmt.Println(s)
	//seelog.Info("parm-url:", s)
	//seelog.Info("sing:", sign)
	if signType == K_SIGN_TYPE_RSA {
		err = encoding.VerifyPKCS1v15([]byte(s), sign, publicKey, crypto.SHA1)
	} else {
		err = encoding.VerifyPKCS1v15([]byte(s), sign, publicKey, crypto.SHA256)
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func verifyResponseData(data []byte, signType, sign string, key []byte) (ok bool, err error) {
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	if signType == K_SIGN_TYPE_RSA {
		err = encoding.VerifyPKCS1v15(data, signBytes, key, crypto.SHA1)
	} else {
		err = encoding.VerifyPKCS1v15(data, signBytes, key, crypto.SHA256)
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
