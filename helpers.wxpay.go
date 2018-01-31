package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

type AccessTokenErrorResponse struct {
	Errcode float64
	Errmsg  string
}

type SnsapiBase struct {
	AccessToken  string
	ExpiresIn    int
	RefreshToken string
	Openid       string
	Scope        string
	Unionid      string
}

//微信支付签名验证函数
func WxpayVerifySign(needVerifyM map[string]interface{}, sign string) bool {

	signCalc := WxpayCalcSign(needVerifyM, WxAppKey)

	seelog.Info("计算出来的sign: ", signCalc)
	seelog.Info("微信异步通知sign: ", sign)

	if sign == signCalc {
		fmt.Println("签名校验通过!")
		return true
	}

	fmt.Println("签名校验失败!")
	return false
}

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func WxpayCalcSign(mReq map[string]interface{}, key string) (sign string) {
	fmt.Println("微信支付签名计算, API KEY:", key)
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		fmt.Printf("k=%v, v=%v\n", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}

	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}

//和国内时间转换一致
func GetCurrentTime() time.Time {
	location := time.FixedZone("Asia/Shanghai", +8*60*60)
	abnow := time.Now().UTC()
	_now := abnow.In(location)
	return _now
}
func TimeConvert(span int) string {
	_now := GetCurrentTime()
	if span == 2 {
		_now = _now.Add(time.Hour * 2)
	}
	return _now.Format("20060102150405")
}

//随机订单号
func GetOutTradeNo() string {

	return time.Now().Format("0102150405") + fmt.Sprintf("%05d", rand.Intn(100000))
}
func GetWxId(code string) (openid string) {

	if sn, err := getWxOpendidFromoauth2(code, WxAppId, WxAppSecret); err == nil {
		openid = sn.Openid
	}

	return
}

func GetWxOrId(c *gin.Context, suburl string) (openid string) {

	if code, ok := c.GetQuery("code"); ok {
		if sn, err := getWxOpendidFromoauth2(code, WxAppId, WxAppSecret); err == nil {
			openid = sn.Openid
		}

	} else {
		encodUrl := url.QueryEscape("http://" + DomainUrl + suburl)
		//seelog.Info(encodUrl)
		urlStr := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" +
			WxAppId + "&redirect_uri=" + encodUrl +
			"&response_type=code&scope=snsapi_base&state=STATE#wechat_redirect"

		c.Redirect(302, urlStr)
	}

	return
}

func getWxOpendidFromoauth2(code, appid, secret string) (*SnsapiBase, error) {

	requestLine := strings.Join([]string{"https://api.weixin.qq.com/sns/oauth2/access_token",
		"?appid=", appid,
		"&secret=", secret,
		"&code=", code,
		"&grant_type=authorization_code"}, "")
	//seelog.Info(requestLine)
	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("发送get请求获取 openid 错误", err)
		seelog.Error("发送get请求获取 openid 错误", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("发送get请求获取 openid 读取返回body错误", err)
		seelog.Error("发送get请求获取 openid 读取返回body错误", err)
		return nil, err
	}
	if bytes.Contains(body, []byte("errcode")) {

		ater := AccessTokenErrorResponse{}
		err = json.Unmarshal(body, &ater)

		if err != nil {
			fmt.Printf("发送get请求获取 openid 的错误信息 %+v\n", ater)
			seelog.Error("发送get请求获取 openid 的错误信息 %+v\n", ater)
			return nil, err
		}
		return nil, fmt.Errorf("%s", ater.Errmsg)

	} else {
		atr := SnsapiBase{}
		//seelog.Info("body:", body)
		err = json.Unmarshal(body, &atr)
		if err != nil {
			fmt.Println("发送get请求获取 openid 返回数据json解析错误", err)
			seelog.Error("发送get请求获取 openid 返回数据json解析错误", err)
			return nil, err
		}
		return &atr, nil
	}
}
