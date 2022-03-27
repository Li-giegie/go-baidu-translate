package BaiDuFanYi

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const appid = `20220327001144941`

const key = `bRLXWsSHVinqfTkAmaS0`

//通用翻译API HTTPS 地址
const baiduUrl = `https://fanyi-api.baidu.com/api/trans/vip/translate`

var ErrorCodeInfo map[string]map[string]string

type Trans_result struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type Translate_obj struct {
	From         string         `json:"from"`
	To           string         `json:"to"`
	TransResults []Trans_result `json:"trans_result"`
	Error_code   string         `json:"error_Code"`
	Error_msg    string         `json:"error_Msg"`
}

type BaiDu_translate struct {
	appid      string
	secretKey  string
	BaiDuAPI   string `json:"baidu_util_url"`
	Query      string `json:"query"`
	sign       string
	rndNum     string
	RequestUrl string `json:"request_url"`
	srcLang    string
	toLang     string
}

//计算签名
//count Sign
func (BaiDu *BaiDu_translate) setSign() {
	//appid+q+salt(rand num )+密钥
	BaiDu.rndNum = fmt.Sprintf("%v", time.Now().Unix())

	BaiDu.sign = fmt.Sprintf("%x", md5.Sum([]byte(BaiDu.appid+BaiDu.Query+BaiDu.rndNum+BaiDu.secretKey)))

}

//设置请求URL
//set URL
func (BaiDu *BaiDu_translate) setRequestUrl() {
	fastUrl := BaiDu.BaiDuAPI
	fmt.Println("fastUrl ", fastUrl)
	if fastUrl == "" {
		fastUrl = baiduUrl
	}
	if BaiDu.srcLang == "" {
		BaiDu.srcLang = "auto"
	}
	url := fastUrl + "?q=" + url.QueryEscape(BaiDu.Query) + "&from=" + BaiDu.srcLang + "&to=" + BaiDu.toLang + "&appid=" + BaiDu.appid + "&salt=" + BaiDu.rndNum + "&sign=" + BaiDu.sign

	fmt.Println("准备请求的url ", url)

	BaiDu.RequestUrl = url
}

func New(appid string, secretKey string) BaiDu_translate {
	//设置错误map
	setErrorCodeInfo()
	return BaiDu_translate{appid: appid, secretKey: secretKey}
}

func (BaiDu *BaiDu_translate) Run(text string, srcLang string, toLang string) (Translate_obj, error) {
	var respTran Translate_obj
	BaiDu.Query = text

	BaiDu.srcLang = srcLang

	BaiDu.toLang = toLang

	BaiDu.setSign()

	BaiDu.setRequestUrl()

	resp, err := http.Get(BaiDu.RequestUrl)

	if err != nil {
		return respTran, fmt.Errorf("请求服务器错误：%v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(data, &respTran)
	if respTran.Error_code != "" {
		fmt.Println("resp error >>", GetErrorCodeInfo(respTran.Error_code))
	}
	return respTran, err

}

//func main() {
//
//	bdtran := New("20220327001144941", "bRLXWsSHVinqfTkAmaS0")
//	bdtran.BaiDuAPI = ""
//	str := []string{"苹果", "橘子", "香蕉"}
//	for i := 0; i < len(str); i++ {
//		fmt.Println(bdtran.Run(str[i], "auto", "en"))
//	}
//
//}
func setErrorCodeInfo() {
	ErrorCodeInfo = make(map[string]map[string]string)

	ErrorCodeInfo["52000"] = map[string]string{"成功": ""}
	ErrorCodeInfo["52001"] = map[string]string{"请求超时": "请重试 "}
	ErrorCodeInfo["52002"] = map[string]string{"系统错误": "请重试"}
	ErrorCodeInfo["52003"] = map[string]string{"未授权用户": "请检查appid是否正确或者服务是否开通"}

	ErrorCodeInfo["54000"] = map[string]string{"必填参数为空": "请检查是否少传参数"}
	ErrorCodeInfo["54001"] = map[string]string{"签名错误": "请检查您的签名生成方法"}
	ErrorCodeInfo["54003"] = map[string]string{"访问频率受限": "请降低您的调用频率，或进行身份认证后切换为高级版/尊享版"}
	ErrorCodeInfo["54004"] = map[string]string{"账户余额不足": "请前往管理控制台为账户充值"}
	ErrorCodeInfo["54005"] = map[string]string{"长query请求频繁": "请降低长query的发送频率，3s后再试"}

	ErrorCodeInfo["58000"] = map[string]string{"客户端IP非法": "检查个人资料里填写的IP地址是否正确，可前往开发者信息-基本信息修改"}
	ErrorCodeInfo["58001"] = map[string]string{"译文语言方向不支持": "检查译文语言是否在语言列表里"}
	ErrorCodeInfo["58002"] = map[string]string{"服务当前已关闭": "请前往管理控制台开启服务"}

	ErrorCodeInfo["90107 "] = map[string]string{"认证未通过或未生效 ": "请前往我的认证查看认证进度 "}
}
func GetErrorCodeInfo(mapkey string) map[string]string {
	return ErrorCodeInfo[mapkey]
}
