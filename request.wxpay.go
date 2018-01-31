package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"strings"

	"github.com/astaxie/beego"
	"github.com/gin-gonic/gin"
)

type OrderqueryReq struct {
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"mch_id"`
	Transaction_id string `xml:"transaction_id"`
	Nonce_str      string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
}
type JSUnifyOrderResp struct {
	Appid     string `from:"appid"`
	Nonce_str string `from:"nonce_str"`
	Sign      string `from:"sign"`
	TimeStamp string `from:"time_stamp"`
	Package   string `from:"package"`
	SignType  string `from:"sign_type"`
}

type OrderqueryResp struct {
	Return_code      string `xml:"return_code"`
	Return_msg       string `xml:"return_msg"`
	Appid            string `xml:"appid"`
	Mch_id           string `xml:"mch_id"`
	Nonce_str        string `xml:"nonce_str"`
	Sign             string `xml:"sign"`
	Result_code      string `xml:"result_code"`
	Openid           string `xml:"prepay_id"`
	Trade_type       string `xml:"trade_type"`
	Trade_state      string `xml:"trade_state"`
	Bank_type        string `xml:"bank_type"`
	Total_fee        string `xml:"total_fee"`
	Cash_fee         int    `xml:"cash_fee"`
	Transaction_id   string `xml:"transaction_id"`
	Out_trade_no     string `xml:"out_trade_no"`
	Time_end         string `xml:"time_end"`
	Trade_state_desc string `xml:"trade_state_desc"`
}

type WXPayNotifyReq struct {
	Return_code    string `xml:"return_code"`
	Return_msg     string `xml:"return_msg"`
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"mch_id"`
	Nonce          string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Result_code    string `xml:"result_code"`
	Openid         string `xml:"openid"`
	Is_subscribe   string `xml:"is_subscribe"`
	Trade_type     string `xml:"trade_type"`
	Bank_type      string `xml:"bank_type"`
	Total_fee      int    `xml:"total_fee"`
	Fee_type       string `xml:"fee_type"`
	Cash_fee       int    `xml:"cash_fee"`
	Cash_fee_Type  string `xml:"cash_fee_type"`
	Transaction_id string `xml:"transaction_id"`
	Out_trade_no   string `xml:"out_trade_no"`
	Attach         string `xml:"attach"`
	Time_end       string `xml:"time_end"`
}
type UnifyOrderReq struct {
	Attach           string `xml:"attach"`
	Appid            string `xml:"appid"`
	Body             string `xml:"body"`
	Detail           string `xml:"detail"`
	Fee_type         string `xml:"fee_type"`
	Goods_tag        string `xml:"goods_tag"`
	Mch_id           string `xml:"mch_id"`
	Nonce_str        string `xml:"nonce_str"`
	Notify_url       string `xml:"notify_url"`
	Product_id       string `xml:"product_id"`
	Time_start       string `xml:"time_start"`
	Time_expire      string `xml:"time_expire"`
	Trade_type       string `xml:"trade_type"`
	Open_id          string `xml:"openid"`
	Spbill_create_ip string `xml:"spbill_create_ip"`
	Total_fee        int    `xml:"total_fee"`
	Out_trade_no     string `xml:"out_trade_no"`
	Sign             string `xml:"sign"`
}

type UnifyOrderResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
	Attach      string `xml:"attach"`
	Appid       string `xml:"appid"`
	Mch_id      string `xml:"mch_id"`
	Nonce_str   string `xml:"nonce_str"`
	Sign        string `xml:"sign"`
	Result_code string `xml:"result_code"`
	Prepay_id   string `xml:"prepay_id"`
	Trade_type  string `xml:"trade_type"`
	Code_url    string `xml:"code_url"`
}

type WXPayNotifyResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
}

func (o *UnifyOrderReq) CreateOrder(param map[string]interface{}) UnifyOrderResp {
	xmlResp := UnifyOrderResp{}
	unify_order_req := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	var pokerReq UnifyOrderReq
	//pokerReq.Attach = param["attach"].(string)
	pokerReq.Appid = WxAppId //微信开放平台我们创建出来的app的app id
	pokerReq.Body = param["body"].(string)
	//pokerReq.Detail = "123"
	//pokerReq.Fee_type = "CNY"
	//pokerReq.Goods_tag = "WXG"
	pokerReq.Mch_id = WxMchId
	pokerReq.Nonce_str = randStr(32, "alphanum")
	pokerReq.Notify_url = "http://" + DomainUrl + "/v1/wxpay/notify" //异步返回的地址
	//pokerReq.Product_id = param["product_id"].(string)
	pokerReq.Time_start = TimeConvert(1)
	pokerReq.Time_expire = TimeConvert(2)
	pokerReq.Trade_type = "JSAPI"
	pokerReq.Spbill_create_ip = param["ip"].(string)
	pokerReq.Open_id = param["openid"].(string)

	totalFee, _ := strconv.ParseFloat(param["total_fee"].(string), 64)
	totalfeeint := totalFee * 100
	pokerReq.Total_fee = int(totalfeeint) //单位是分，这里是1毛钱
	pokerReq.Out_trade_no = param["out_trade_no"].(string)

	//beego.Debug("pokerReq",pokerReq)
	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["attach"] = pokerReq.Attach
	m["appid"] = pokerReq.Appid
	m["body"] = pokerReq.Body
	m["detail"] = pokerReq.Detail
	m["fee_type"] = pokerReq.Fee_type
	m["goods_tag"] = pokerReq.Goods_tag
	m["mch_id"] = pokerReq.Mch_id
	m["nonce_str"] = pokerReq.Nonce_str
	m["notify_url"] = pokerReq.Notify_url
	//m["product_id"] = pokerReq.Product_id
	m["time_start"] = pokerReq.Time_start
	m["time_expire"] = pokerReq.Time_expire
	m["trade_type"] = pokerReq.Trade_type
	m["spbill_create_ip"] = pokerReq.Spbill_create_ip
	m["total_fee"] = pokerReq.Total_fee
	m["out_trade_no"] = pokerReq.Out_trade_no
	m["openid"] = pokerReq.Open_id
	pokerReq.Sign = WxpayCalcSign(m, WxAppKey) //这个是计算wxpay签名的函数上面已贴出

	//seelog.Info("pokerReq", pokerReq)
	bytes_req, err := xml.Marshal(pokerReq)
	if err != nil {
		fmt.Println("以xml形式编码发送错误, 原因:", err)
		return xmlResp
	}

	str_req := string(bytes_req)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	str_req = strings.Replace(str_req, "UnifyOrderReq", "xml", -1)
	bytes_req = []byte(str_req)

	//发送unified order请求.
	req, err := http.NewRequest("POST", unify_order_req, bytes.NewReader(bytes_req))
	if err != nil {
		fmt.Println("New Http Request发生错误，原因:", err)
		return xmlResp
	}
	req.Header.Set("Accept", "application/xml")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	c := http.Client{}
	resp, _err := c.Do(req)
	if _err != nil {
		fmt.Println("请求微信支付统一下单接口发送错误, 原因:", _err)
		return xmlResp
	}

	//xmlResp :=UnifyOrderResp{}
	body, _ := ioutil.ReadAll(resp.Body)
	_err = xml.Unmarshal(body, &xmlResp)
	if xmlResp.Return_code == "FAIL" {
		fmt.Println("微信支付统一下单不成功，原因:", xmlResp.Return_msg)
		return xmlResp
	}
	//seelog.Info("xmlResp", xmlResp)
	//这里已经得到微信支付的prepay id，需要返给客户端，由客户端继续完成支付流程
	fmt.Println("微信支付统一下单成功，预支付单号:", xmlResp.Prepay_id)
	return xmlResp

}

func (o *OrderqueryReq) WxQueryOrder(transId string) OrderqueryResp {
	xmlResp := OrderqueryResp{}

	query_order_req := "https://api.mch.weixin.qq.com/pay/orderquery"

	var qeuryReq OrderqueryReq
	qeuryReq.Appid = "wx74a6f50de79ccd40" //微信开放平台我们创建出来的app的app id
	qeuryReq.Mch_id = "1319106901"
	qeuryReq.Transaction_id = transId
	qeuryReq.Nonce_str = randStr(32, "alphanum")

	//beego.Debug("qeuryReqEntity",qeuryReq)
	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = qeuryReq.Appid
	m["mch_id"] = qeuryReq.Mch_id
	m["transaction_id"] = qeuryReq.Transaction_id
	m["nonce_str"] = qeuryReq.Nonce_str
	appkey := beego.AppConfig.String("wxappkey")
	qeuryReq.Sign = WxpayCalcSign(m, appkey) //这个是计算wxpay签名的函数上面已贴出

	bytes_req, err := xml.Marshal(qeuryReq)
	if err != nil {
		fmt.Println("以xml形式编码发送错误, 原因:", err)
		return xmlResp
	}

	str_req := string(bytes_req)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	str_req = strings.Replace(str_req, "OrderqueryReq", "xml", -1)
	bytes_req = []byte(str_req)

	//发送unified order请求.
	req, err := http.NewRequest("POST", query_order_req, bytes.NewReader(bytes_req))
	if err != nil {
		fmt.Println("New Http Request发生错误，原因:", err)
		return xmlResp
	}
	req.Header.Set("Accept", "application/xml")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	c := http.Client{}
	resp, _err := c.Do(req)
	if _err != nil {
		fmt.Println("请求微信支付查询发送错误, 原因:", _err)
		return xmlResp
	}

	//xmlResp :=UnifyOrderResp{}
	body, _ := ioutil.ReadAll(resp.Body)
	_err = xml.Unmarshal(body, &xmlResp)
	if xmlResp.Return_code == "FAIL" {
		fmt.Println("微信支付统查询不成功，原因:", xmlResp.Return_msg)
		return xmlResp
	}
	return xmlResp

}
func (o *WXPayNotifyReq) WxpayCallback(c *gin.Context) map[string]interface{} {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("读取http body失败，原因!", err)
	}
	fmt.Println("微信支付异步通知，HTTP Body:", string(body))
	var mr WXPayNotifyReq
	err = xml.Unmarshal(body, &mr)
	if err != nil {
		fmt.Println("解析HTTP Body格式到xml失败，原因!", err)
	}

	//seelog.Info("body", body)

	var reqMap map[string]interface{}
	reqMap = make(map[string]interface{}, 0)

	reqMap["return_code"] = mr.Return_code
	reqMap["return_msg"] = mr.Return_msg
	reqMap["appid"] = mr.Appid
	reqMap["mch_id"] = mr.Mch_id
	reqMap["nonce_str"] = mr.Nonce
	reqMap["result_code"] = mr.Result_code
	reqMap["openid"] = mr.Openid
	reqMap["is_subscribe"] = mr.Is_subscribe
	reqMap["trade_type"] = mr.Trade_type
	reqMap["bank_type"] = mr.Bank_type
	reqMap["total_fee"] = mr.Total_fee
	reqMap["fee_type"] = mr.Fee_type
	reqMap["cash_fee"] = mr.Cash_fee
	reqMap["cash_fee_type"] = mr.Cash_fee_Type
	reqMap["transaction_id"] = mr.Transaction_id
	reqMap["out_trade_no"] = mr.Out_trade_no
	reqMap["attach"] = mr.Attach
	reqMap["time_end"] = mr.Time_end

	var resp WXPayNotifyResp

	//进行签名校验
	if WxpayVerifySign(reqMap, mr.Sign) {
		//这里就可以更新我们的后台数据库了，其他业务逻辑同理。
		//seelog.Info("succes", "succes")
		resp.Return_code = "SUCCESS"
		resp.Return_msg = "OK"
		//ctx.WriteString("SUCCESS")
		//return reqMap
	} else {
		resp.Return_code = "FAIL"
		resp.Return_msg = "failed to verify sign, please retry!"
	}

	//结果返回，微信要求如果成功需要返回return_code "SUCCESS"
	bytes, _err := xml.Marshal(resp)
	strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)
	if _err != nil {
		fmt.Println("xml编码失败，原因：", _err)
		//return
	}
	//seelog.Info("return", strResp)
	c.String(200, strResp)
	//ctx.WriteString(strResp)

	//c.Ctx.ResponseWriter.WriteHeader(200)
	//fmt.Fprint(c.Ctx.ResponseWriter, strResp)

	return reqMap
}
