package main

import (
	"fmt"

	"github.com/cihub/seelog"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AliPayHandler struct {
	db *gorm.DB
}

func InitAliPayHandler(app *App) *AliPayHandler {
	defer seelog.Flush()
	h := &AliPayHandler{
		app.Db(),
	}
	//支付宝异步和同步通知页面自己需要实现
	html := app.engine.Group("html")
	{
		html.GET("/alipayok", func(c *gin.Context) {
			c.HTML(200, "alipayok.html", nil)
		})
		html.GET("/alipayfail", func(c *gin.Context) {
			c.HTML(200, "alipayfail.html", nil)
		})
		html.GET("/alipayerror/:errmsg", func(c *gin.Context) {
			c.HTML(200, "alipayerror.html", gin.H{
				"errmsg": c.Param("errmsg"),
			})
		})
	}
	v1 := app.engine.Group("/v1")
	{
		//支付宝下单借口
		v1.POST("/alipay", h.Native)
		//同步地址
		v1.GET("/return_url", h.Return)
		//异步回调地址
		v1.POST("/notify_url", h.Notify)
	}
	return h
}

func (h *AliPayHandler) Native(c *gin.Context) {
	//获取支付参数保存到订单表
	//........

	//支付宝下单，返回支付html
	out_trade_no := GetOutTradeNo()
	goods := "商品详情"
	alipay_html := NewAlipayTradeAppPayRequest(goods, out_trade_no)

	h.responseHandler.Html(c, alipay_html)

}

//同步地址
func (h *AliPayHandler) Return(c *gin.Context) {
	// 实例化参数
	param := &AliPayNotifyResponse{}

	// 解析表单内容，失败
	if err := c.Bind(param); err != nil {
		h.responseHandler.Redirect(c, "/html/alipayfail")
		//h.responseHandler.MalformedJSON(c)
		//return
	}
	fmt.Println(param)
	if param.OutTradeNo == "" { //不存在交易号
		h.responseHandler.Redirect(c, "/html/alipayfail")
		//h.responseHandler.Error(c, "不存在交易号", 464, "不存在交易号")
		return
	}
	ok, err := verifySign(param.ToMap(), publicKey)
	if err != nil {
		h.responseHandler.Redirect(c, "/html/alipayfail")
		//h.responseHandler.ValidationErrors(c, err)
		return
	}
	if ok { //只有相同才说明该订单成功了
		h.responseHandler.Redirect(c, "/html/alipayok")
	} else { // 签名认证失败，返回错误代码-2
		//h.responseHandler.Error(c, "签名认证失败", 462, "签名认证失败")
		h.responseHandler.Redirect(c, "/html/alipayfail")
		return
	}
}

//异步回调
func (h *AliPayHandler) Notify(c *gin.Context) {
	// 实例化参数
	param := &AliPayNotifyResponse{}

	// 解析表单内容，失败
	if err := c.Bind(param); err != nil {
		seelog.Info("解析表单内容，失败")
		h.responseHandler.String(c, "fail")
	}
	fmt.Println(param)
	// 如果最基本的网站交易号为空
	if param.OutTradeNo == "" { //不存在交易号
		seelog.Info("请求数据中不存在交易号")
		h.responseHandler.String(c, "fail")
	}
	if param.SellerId != partnerID || param.AppId != AppID {
		seelog.Info("appid或sellerid不合法：", param.SellerId, param.AppId != AppID)
		h.responseHandler.String(c, "fail")
	}
	ok, err := verifySign(param.ToMap(), publicKey)
	//seelog.Info("issign:", err)
	if err != nil {
		seelog.Info("签名验证方法出错：", err)
		h.responseHandler.String(c, "fail")
	}
	if ok { //只有相同才说明该订单成功了
		// 判断订单是否已完成
		if param.TradeStatus == "TRADE_FINISHED" || param.TradeStatus == "TRADE_SUCCESS" { //交易成功
			//处理订单
			//、、、、trader_info

			if err != nil {
				seelog.Info("未获取到此订单信息:", param.OutTradeNo)
				h.responseHandler.String(c, "success")
				//h.responseHandler.Error(c, "未获取到此订单信息", 467, "未获取到此订单信息")
				//return
			} else if trader_info.IsSuccess == K_ALI_PAY_SUCCESS {
				seelog.Info("此订单信息已处理:", param.OutTradeNo)
				h.responseHandler.String(c, "success")
				//h.responseHandler.Error(c, "此订单信息已处理", 468, "此订单信息已处理")
				//return
			} else if trader_info.TotalFee+".00" != param.TotalAmount { //支付宝返回的是两位小数格式
				seelog.Info("订单价格不相等：", trader_info.TotalFee, param.TotalAmount)
				h.responseHandler.String(c, "success")
				//h.responseHandler.Error(c, "订单价格不相等", 469, "订单价格不相等：")
			} else {
				//修改定单状态
				//、、、、、
				//给用户加积分操作
				//、、、、
			}
			h.responseHandler.String(c, "success")
		} else { // 交易未完成，返回错误代码-4
			seelog.Info("交易未完成")
			h.responseHandler.String(c, "success")
			//h.responseHandler.Error(c, "交易未完成", 460, "交易未完成")
			//return
		}
	} else { // 签名认证失败，返回错误代码-2
		seelog.Info("签名认证失败")
		h.responseHandler.String(c, "fail")
		//h.responseHandler.Error(c, "签名认证失败", 462, "签名认证失败")
		return
	}
}
