package main

// 列举全部传参
type AliPayNotifyResponse struct {
	AppId      string `form:"app_id" json:"app_id"`             //支付宝分配给开发者的应用Id
	Subject    string `form:"subject" json:"subject"`           //商品的标题/交易标题/订单标题/订单关键字等，是请求时对应的参数，原样通知回来
	Body       string `form:"body" json:"body"`                 //该订单的备注、描述、明细等。对应请求时的body参数，原样通知回来
	TradeNo    string `form:"trade_no" json:"trade_no"`         //支付宝交易凭证号
	OutTradeNo string `form:"out_trade_no" json:"out_trade_no"` //原支付请求的商户订单号
	OutBizNo   string `form:"out_biz_no" json:"out_biz_no"`     //商户业务ID，主要是退款通知中返回退款申请的流水号
	BuyerId    string `form:"buyer_id" json:"buyer_id"`         //买家支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字

	AuthAppId         string `form:"auth_app_id" json:"auth_app_id"`                 //授权商户的AppIdString是授权商户的AppId
	BuyerLogonId      string `form:"buyer_logon_id" json:"buyer_logon_id"`           //买家支付宝账号
	SellerId          string `form:"seller_id" json:"seller_id"`                     //卖家支付宝用户号
	SellerEmail       string `form:"seller_email" json:"seller_email"`               //卖家支付宝账号
	TradeStatus       string `form:"trade_status" json:"trade_status"`               //交易目前所处的状态
	TotalAmount       string `form:"total_amount" json:"total_amount"`               //本次交易支付的订单金额，单位为人民币（元）
	ReceiptAmount     string `form:"receipt_amount" json:"receipt_amount"`           //商家在交易中实际收到的款项，单位为元
	InvoiceAmount     string `form:"invoice_amount" json:"invoice_amount"`           //用户在交易中支付的可开发票的金额
	BuyerPayAmount    string `form:"buyer_pay_amount" json:"buyer_pay_amount"`       //用户在交易中支付的金额
	PointAmount       string `form:"point_amount" json:"point_amount"`               //使用集分宝支付的金额
	RefundFee         string `form:"refund_fee" json:"refund_fee"`                   //退款通知中，返回总退款金额，单位为元，支持两位小数
	FundBillList      string `form:"fund_bill_list" json:"fund_bill_list"`           //支付成功的各个渠道金额信息
	Timestamp         string `form:"timestamp" json:"timestamp"`                     //同步时间戳
	Method            string `form:"method" json:"method"`                           //接口方法
	VoucherDetailList string `form:"voucher_detail_list" json:"voucher_detail_list"` //本交易支付时所使用的所有优惠券信息
	PassbackParams    string `form:"passback_params" json:"passback_params"`         //公共回传参数，如果请求时传递了该参数，则返回给商户时会在异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝
	Charset           string `form:"charset" json:"charset"`                         //编码格式，如utf-8、gbk、gb2312等
	Sign              string `form:"sign" json:"sign"`                               //601510b7970e52cc63db0f44997cf70e
	SignType          string `form:"sign_type" json:"sign_type"`                     //商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	NotifyId          string `form:"notify_id" json:"notify_id"`                     //通知校验ID
	NotifyType        string `form:"notify_type" json:"notify_type"`                 //通知的类型
	NotifyTime        string `form:"notify_time" json:"notify_time"`                 //通知的发送时间。格式为yyyy-MM-dd HH:mm:ss
	GmtCreate         string `form:"gmt_create" json:"gmt_create"`                   //该笔交易创建的时间。格式为yyyy-MM-dd HH:mm:ss
	GmtPayment        string `form:"gmt_payment" json:"gmt_payment"`                 //该笔交易的买家付款时间。格式为yyyy-MM-dd HH:mm:ss
	GmtRefund         string `form:"gmt_refund" json:"gmt_refund"`                   //该笔交易的退款时间。格式为yyyy-MM-dd HH:mm:ss.S
	GmtClose          string `form:"gmt_close" json:"gmt_close"`                     //该笔交易结束时间。格式为yyyy-MM-dd HH:mm:ss
	Version           string `form:"version" json:"version"`                         //调用的接口版本，固定为：1.0
}

//对象转成字典
func (s *AliPayNotifyResponse) ToMap() map[string]string {
	mapData := map[string]string{
		"app_id":              s.AppId,
		"subject":             s.Subject,
		"body":                s.Body,
		"trade_no":            s.TradeNo,
		"auth_app_id":         s.AuthAppId,
		"out_trade_no":        s.OutTradeNo,
		"out_biz_no":          s.OutBizNo,
		"buyer_id":            s.BuyerId,
		"buyer_logon_id":      s.BuyerLogonId,
		"seller_id":           s.SellerId,
		"seller_email":        s.SellerEmail,
		"trade_status":        s.TradeStatus,
		"total_amount":        s.TotalAmount,
		"receipt_amount":      s.ReceiptAmount,
		"invoice_amount":      s.InvoiceAmount,
		"buyer_pay_amount":    s.BuyerPayAmount,
		"point_amount":        s.PointAmount,
		"refund_fee":          s.RefundFee,
		"timestamp":           s.Timestamp,
		"method":              s.Method,
		"fund_bill_list":      s.FundBillList,
		"voucher_detail_list": s.VoucherDetailList,
		"passback_params":     s.PassbackParams,
		"charset":             s.Charset,
		"sign":                s.Sign,
		"sign_type":           s.SignType,
		"notify_id":           s.NotifyId,
		"notify_type":         s.NotifyType,
		"notify_time":         s.NotifyTime,
		"gmt_create":          s.GmtCreate,
		"gmt_payment":         s.GmtPayment,
		"gmt_refund":          s.GmtRefund,
		"gmt_close":           s.GmtClose,
		"version":             s.Version,
	}

	return mapData
}
