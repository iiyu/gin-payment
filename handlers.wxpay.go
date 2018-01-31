package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cihub/seelog"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	K_WX_PAY_FAILED  = "0"
	K_WX_PAY_SUCCESS = "1"
)

type WxPayHandler struct {
	db *gorm.DB
}

func InitWxPayHandler(app *App) *WxPayHandler {
	defer seelog.Flush()
	h := &WxPayHandler{
		app.Db(),
	}
	v := app.engine.Group("")
	{
		//微信JSAPI要求域名下有这个文件
		v.StaticFile("MP_verify_retetryfdg.txt", basePath+"/src/gin-payment/MP_verify_retetryfdg.txt")
	}
	v1 := app.engine.Group("/v1")
	{
		v1.POST("/wxpay/notify", h.WxpayCallback)
		v1.POST("/wxpay", h.WxQueryOrde)
	}

	return h
}

func (h *WxPayHandler) WxQueryOrde(c *gin.Context) {
	//获取支付参数保存到订单表
	//........

	params["body"] = "测试" //显示标题
	params["out_trade_no"] = out_trade_no
	params["total_fee"] = ToDiscount(Money, Discount)
	params["ip"] = c.ClientIP()
	//params["attach"] = "abc" //自定义参数
	//params["openid"] = wxuser.Openid
	params["openid"] = c.MustGet("openid").(string)

	var modwx UnifyOrderReq
	res := modwx.CreateOrder(params)
	//ip := c.ClientIP()
	fmt.Println(res)
	jsUniforderResp := JSUnifyOrderResp{}
	jsUniforderResp.Appid = res.Appid
	times := fmt.Sprintf("%d", time.Now().Unix())
	jsUniforderResp.TimeStamp = times
	jsUniforderResp.Package = "prepay_id=" + res.Prepay_id
	jsUniforderResp.SignType = "MD5"
	Nonce_strs := randStr(32, "alphanum")
	jsUniforderResp.Nonce_str = Nonce_strs
	retur := make(map[string]interface{})
	retur["appId"] = res.Appid
	retur["timeStamp"] = times
	retur["package"] = "prepay_id=" + res.Prepay_id
	retur["signType"] = "MD5"
	retur["nonceStr"] = Nonce_strs
	jsUniforderResp.Sign = WxpayCalcSign(retur, WxAppKey)
	//seelog.Info("retur:", retur)
	h.responseHandler.JSON(c, StatusOK, jsUniforderResp)
}

func (h *WxPayHandler) GetWxId(c *gin.Context) (openid string) {

	if code, ok := c.GetQuery("code"); ok {
		if sn, err := getWxOpendidFromoauth2(code, WxAppId, WxMchId); err == nil {
			openid = sn.Openid
		}

	} else {
		encodUrl := url.QueryEscape("http://" + DomainUrl + "/wxpay/getwxid")
		//seelog.Info(encodUrl)
		urlStr := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" +
			WxAppId + "&redirect_uri=" + encodUrl +
			"&response_type=code&scope=snsapi_base&state=STATE#wechat_redirect"

		c.Redirect(302, urlStr)
	}

	return
}

func (h *WxPayHandler) WxpayCallback(c *gin.Context) {
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
		seelog.Info("succes:", "succes")
		//处理订单
		//、、、、trader_info

		if err != nil {
			seelog.Info("未获取到此订单信息:", reqMap["out_trade_no"].(string))
			//h.responseHandler.Error(c, "未获取到此订单信息", 467, "未获取到此订单信息")
			//return
		} else if trader_info.IsSuccess == K_WX_PAY_SUCCESS {
			seelog.Info("此订单信息已处理:", reqMap["out_trade_no"].(string))
			//h.responseHandler.Error(c, "此订单信息已处理", 468, "此订单信息已处理")
			//return
		} else if trader_info.TotalFee != strconv.Itoa(reqMap["total_fee"].(int)/100) {
			seelog.Info("订单价格不相等：", trader_info.TotalFee, "-", reqMap["total_fee"].(int)/100)
			//h.responseHandler.Error(c, "订单价格不相等", 469, "订单价格不相等：")
		} else {
			//修改定单状态
			//、、、、、
			//给用户加积分操作
			//、、、、
		}
		resp.Return_code = "SUCCESS"
		resp.Return_msg = "OK"
	} else {
		resp.Return_code = "FAIL"
		resp.Return_msg = "failed to verify sign, please retry!"
	}

	//结果返回，微信要求返回return_code "SUCCESS"
	bytes, _err := xml.Marshal(resp)
	strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)
	if _err != nil {
		fmt.Println("xml编码失败，原因：", _err)
		//return
	}
	seelog.Info("return：", strResp)
	c.String(200, strResp)
	return
}
