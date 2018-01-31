package main

var (
	privateKey = []byte(`
-----BEGIN PRIVATE KEY-----
-----END PRIVATE KEY-----
`)
	publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
-----END PUBLIC KEY-----
`)
	partnerID = "" // fake partner ID
	sellerID  = "" // fake seller ID
	AppID     = ""
)

func NewAlipayTradeAppPayRequest(goods *Goods, out_trade_no string) string {
	var client = AlipayNew(AppID, partnerID, publicKey, privateKey, true)
	var p = AliPayTradeWapPay{}
	p.NotifyURL = "http://" + DomainUrl + "/v1/notify_url"
	p.ReturnURL = "http://" + DomainUrl + "/v1/return_url"
	p.Subject = "123"
	p.OutTradeNo = out_trade_no
	p.TotalAmount = ToDiscount(Money, Discount)
	p.ProductCode = "QUICK_WAP_WAY"
	pa, _ := client.PageExecute(p)
	a := p.ToHTML(pa)
	return a
}
