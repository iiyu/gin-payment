package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// AliPayTradeWapPay https://doc.open.alipay.com/doc2/detail.htm?treeId=203&articleId=105463&docType=1
type AliPayTradeWapPay struct {
	NotifyURL string `json:"-"`
	ReturnURL string `json:"-"`

	// biz content，这四个参数是必须的
	Subject     string `json:"subject"`      // 商品的标题/交易标题/订单标题/订单关键字等。
	OutTradeNo  string `json:"out_trade_no"` // 商户网站唯一订单号
	TotalAmount string `json:"total_amount"` // 订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	ProductCode string `json:"product_code"` // 销售产品码，商家和支付宝签约的产品码。该产品请填写固定值：QUICK_WAP_WAY

	Body               string `json:"body,omitempty"`                 // 对一笔交易的具体描述信息。如果是多种商品，请将商品描述字符串累加传给body。
	TimeoutExpress     string `json:"timeout_express,omitempty"`      // 该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。
	SellerId           string `json:"seller_id,omitempty"`            // 收款支付宝用户ID。 如果该值为空，则默认为商户签约账号对应的支付宝用户ID
	AuthToken          string `json:"auth_token,omitempty"`           // 针对用户授权接口，获取用户相关数据时，用于标识用户授权关系
	GoodsType          string `json:"goods_type,omitempty"`           // 商品主类型：0—虚拟类商品，1—实物类商品 注：虚拟类商品不支持使用花呗渠道
	PassbackParams     string `json:"passback_params,omitempty"`      // 公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝会在异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝
	PromoParams        string `json:"promo_params,omitempty"`         // 优惠参数 注：仅与支付宝协商后可用
	ExtendParams       string `json:"extend_params,omitempty"`        // 业务扩展参数，详见下面的“业务扩展参数说明”
	EnablePayChannels  string `json:"enable_pay_channels,omitempty"`  // 可用渠道，用户只能在指定渠道范围内支付 当有多个渠道时用“,”分隔 注：与disable_pay_channels互斥
	DisablePayChannels string `json:"disable_pay_channels,omitempty"` // 禁用渠道，用户不可用指定渠道支付 当有多个渠道时用“,”分隔 注：与enable_pay_channels互斥
	StoreId            string `json:"store_id,omitempty"`             // 商户门店编号。该参数用于请求参数中以区分各门店，非必传项。
}

func (this AliPayTradeWapPay) APIName() string {
	return "alipay.trade.wap.pay"
}

func (this AliPayTradeWapPay) Params() map[string]string {
	var m = make(map[string]string)
	m["notify_url"] = this.NotifyURL
	m["return_url"] = this.ReturnURL
	return m
}

func (this AliPayTradeWapPay) ExtJSONParamName() string {
	return "biz_content"
}

func (this AliPayTradeWapPay) ExtJSONParamValue() string {
	var bytes, err = json.Marshal(this)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// ToHTML 转换form表单
func (this AliPayTradeWapPay) ToHTML(m map[string]string) string {
	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(fmt.Sprintf(`<input type='hidden' name='%s' value='%s'>`, k, strings.Replace(v, "'", "&apos;", -1)))
	}
	formatStr :=
		`<html>
    <meta http-equiv=Content-Type content="text/html;charset=utf-8">
    <body>
        <form id='paysubmit' name='paysubmit' action='%s' method = 'POST'>
        %s
        <input type='submit' value='ok' style='display:none'>
        </form>
        <script>
         (function(){
             document.forms['paysubmit'].submit();
         })();
        </script>
    </body>
</html>`
	//url=charset=utf-8必须加上这个编码并且要和formdata中的charset保持一致
	return fmt.Sprintf(formatStr, "https://openapi.alipay.com/gateway.do?charset=utf-8", buf.String())
}

func (this *AliPay) PageExecute(param AliPayParam) (p map[string]string, err error) {
	p = make(map[string]string)
	p["app_id"] = this.appId
	p["method"] = param.APIName()
	p["format"] = K_FORMAT
	p["charset"] = K_CHARSET
	p["sign_type"] = this.SignType
	p["timestamp"] = time.Now().Format(K_TIME_FORMAT)
	p["version"] = K_VERSION

	if len(param.ExtJSONParamName()) > 0 {
		p[param.ExtJSONParamName()] = param.ExtJSONParamValue()
	}

	var ps = param.Params()
	if ps != nil {
		for key, value := range ps {
			p[key] = value
		}
	}

	var keys = make([]string, 0, 0)
	for key, _ := range p {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var sign string
	if this.SignType == K_SIGN_TYPE_RSA {
		sign, err = signRSA(keys, p, this.privateKey)
	} else {
		sign, err = signRSA2(keys, p, this.privateKey)
	}
	p["sign"] = sign

	if err != nil {
		return nil, err
	}
	return p, nil
}
