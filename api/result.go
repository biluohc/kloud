package api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
)

// int16 在-32768 到 +32767
// int 是当前操作系统位数(符合echo的api)，32/64都没有问题
// 可能需要整理一张 const 变量表，这样挺乱的

// 如果后面需要多语言支持可以自定义通过错误码在前端实现翻译？
// 比如方法是 4010 属于 401, 后面的10个可以是更详细的错误，echo那里判断大小也可以兼容标准的 http 状态码
// 注意如果给 echo 1000+ 作为状态码将得到 500..
type Result struct {
	//状态码
	Code int `json:"code"`
	//请求结果，通常会是出错时指明出错原因，或者修改权限等小操作时提示的`操作成功`
	Msg string `json:"msg"`
	//错误时为 null
	Data interface{} `json:"data"`
}

func (self *Result) Error() string {
	return self.Msg
}

func (self *Result) String() string {
	json, _ := json.Marshal(self)
	return string(json)
}

func Ok(code int, data interface{}) *Result {
	return &Result{
		Code: code,
		Data: data,
	}
}

func OkWithMsg(code int, message string, data interface{}) *Result {
	return &Result{
		Code: code,
		Msg:  message,
		Data: data,
	}
}

func Err(code int, message string) *Result {
	return &Result{
		Code: code,
		Msg:  message,
		// 注释与否一样的效果, 都是 null
		// Data:    nil,
	}
}

// 有点像走了魔道一般，但是 return Ok(...) 好方便，不用多套个 c.JSON 再写两个状态码了。
func HTTPErrorHandler(err error, c echo.Context) {
	// 注意这里都是指针。。
	if res, ok := err.(*Result); ok {
		httpCode := res.Code
		if res.Code > 1000 {
			httpCode = res.Code / 10
		}
		c.JSON(httpCode, res)
	} else if es, ok := err.(*echo.HTTPError); ok {
		c.JSON(es.Code, Err(es.Code, err.Error()))
	} else {
		code := http.StatusInternalServerError
		c.JSON(code, Err(code, err.Error()))
	}
}
