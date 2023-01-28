package gin

import "C"
import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"mime/multipart"
	"strconv"
)

type IRequest interface {
	//查询URL中的参数
	DefaultQueryInt(key string, def int) (int, bool)
	DefaultQueryInt64(key string, def int64) (int64, bool)
	DefaultQueryFloat32(key string, def float32) (float32, bool)
	DefaultQueryFloat64(key string, def float64) (float64, bool)
	DefaultQueryBool(key string, def bool) (bool, bool)
	DefaultQueryString(key string, def string) (string, bool)
	DefaultQueryStringSlice(key string, def []string) ([]string, bool)

	//路由匹配中的参数
	DefaultParamInt(key string, def int) (int, bool)
	DefaultParamInt64(key string, def int64) (int64, bool)
	DefaultParamFloat32(key string, def float32) (float32, bool)
	DefaultParamFloat64(key string, def float64) (float64, bool)
	DefaultParamBool(key string, def bool) (bool, bool)
	DefaultParamString(key string, def string) (string, bool)
	DefaultParam(key string) interface{}

	//form表单中带的参数
	DefaultFormInt(key string, def int) (int, bool)
	DefaultFormInt64(key string, def int64) (int64, bool)
	DefaultFormFloat32(key string, def float32) (float32, bool)
	DefaultFormFloat64(key string, def float64) (float64, bool)
	DefaultFormBool(key string, def bool) (bool, bool)
	DefaultFormString(key string, def string) (string, bool)
	DefaultFormStringSlice(key string, def []string) ([]string, bool)
	DefaultFormFile(key string) (*multipart.FileHeader, error)
	DefaultForm(key string) interface{}

	//绑定JSON
	BindJson(obj interface{}) error
	//绑定XML
	BindXml(obj interface{}) error

	//获取原始数据
	GetRawData() ([]byte, error)

	//基本信息
	Uri() string
	Method() string
	Host() string
	ClientIp() string

	//头信息
	Headers() map[string][]string
	Header(key string) (string, bool)

	//cookie
	Cookies() map[string]string
	Cookie(key string) (string, bool)
}

func (c *Context) QueryAll() map[string][]string {
	c.initQueryCache()
	return map[string][]string(c.queryCache)
}

func (c *Context) DefaultQueryInt(key string, def int) (int, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.Atoi(val[len-1])
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) DefaultQueryInt64(key string, def int64) (int64, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseInt(val[len-1], 10, 64)
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) DefaultQueryFloat32(key string, def float32) (float32, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseFloat(val[len-1], 32)
			if err != nil {
				return def, false
			}
			return float32(v), true
		}
	}
	return def, false
}

func (c *Context) DefaultQueryFloat64(key string, def float64) (float64, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseFloat(val[len-1], 64)
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) DefaultQueryBool(key string, def bool) (bool, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseBool(val[len-1])
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) DefaultQueryString(key string, def string) (string, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			return val[len-1], true
		}
	}
	return def, false
}

func (c *Context) DefaultQueryStringSlice(key string, def []string) ([]string, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			return val, true
		}
	}
	return def, false
}

func (c *Context) GetExtParam(key string) string {
	if val, ok := c.Params.Get(key); ok {
		return val
	}
	return ""
}

func (c *Context) DefaultParamInt(key string, def int) (int, bool) {
	val := c.GetExtParam(key)
	if val != "" {
		v, err := strconv.Atoi(val)
		if err != nil {
			return def, false
		}
		return v, true
	}
	return def, false
}

func (c *Context) DefaultParamInt64(key string, def int64) (int64, bool) {
	val := c.GetExtParam(key)
	if val != "" {
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return def, false
		}
		return v, true
	}
	return def, false
}

func (c *Context) DefaultParamFloat32(key string, def float32) (float32, bool) {
	val := c.GetExtParam(key)
	if val != "" {
		v, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return def, false
		}
		return float32(v), true
	}
	return def, false
}

func (c *Context) DefaultParamFloat64(key string, def float64) (float64, bool) {
	val := c.GetExtParam(key)
	if val != "" {
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return def, false
		}
		return v, true
	}
	return def, false
}

func (c *Context) DefaultParamBool(key string, def bool) (bool, bool) {
	val := c.GetExtParam(key)
	if val != "" {
		v, err := strconv.ParseBool(val)
		if err != nil {
			return def, false
		}
		return v, true
	}
	return def, false
}

func (c *Context) DefaultParamString(key string, def string) (string, bool) {
	val := c.GetExtParam(key)
	if val != "" {
		return val, true
	}
	return def, false
}

func (c *Context) DefaultParam(key string) interface{} {
	return c.GetExtParam(key)
}

func (c *Context) FormAll() map[string][]string {
	c.initFormCache()
	return map[string][]string(c.formCache)
}

func (c *Context) DefaultFormInt(key string, def int) (int, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.Atoi(val[len-1])
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) DefaultFormInt64(key string, def int64) (int64, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseInt(val[len-1], 10, 64)
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) DefaultFormFloat32(key string, def float32) (float32, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseFloat(val[len-1], 32)
			if err != nil {
				return def, false
			}
			return float32(v), true
		}
	}
	return def, false
}

func (c *Context) DefaultFormFloat64(key string, def float64) (float64, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseFloat(val[len-1], 64)
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) DefaultFormBool(key string, def bool) (bool, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseBool(val[len-1])
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) DefaultFormString(key string, def string) (string, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			return val[len-1], true
		}
	}
	return def, false
}

func (c *Context) DefaultFormStringSlice(key string, def []string) ([]string, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			return val, true
		}
	}
	return def, false
}

func (c *Context) DefaultFormFile(key string) (*multipart.FileHeader, error) {
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(defaultMultipartMemory); err != nil {
			return nil, err
		}
	}
	file, header, err := c.Request.FormFile(key)
	if err != nil {
		return nil, err
	}
	file.Close()
	return header, nil
}

func (c *Context) DefaultForm(key string) interface{} {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			return val
		}
	}
	return nil
}

func (c *Context) BindJson(obj interface{}) error {
	if c.Request != nil {
		all, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return err
		}
		//body只能读一次，读出来后需要重置下body
		c.Request.Body = io.NopCloser(bytes.NewBuffer(all))

		err = json.Unmarshal(all, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("request empty")
	}
	return nil
}

func (c *Context) BindXml(obj interface{}) error {
	if c.Request != nil {
		all, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return err
		}
		//body只能读一次，读出来后需要重置下body
		c.Request.Body = io.NopCloser(bytes.NewBuffer(all))

		err = xml.Unmarshal(all, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("request empty")
	}
	return nil
}

func (c *Context) Uri() string {
	return c.Request.RequestURI
}

func (c *Context) Method() string {
	return c.Request.Method
}

func (c *Context) Host() string {
	return c.Request.URL.Host
}

func (c *Context) ClientIp() string {
	ip := c.Request.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = c.Request.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = c.Request.RemoteAddr
	}
	return ip
}

func (c *Context) Headers() map[string][]string {
	return c.Request.Header
}

func (c *Context) Cookies() map[string]string {
	cookies := c.Request.Cookies()
	ret := map[string]string{}
	for _, cookie := range cookies {
		ret[cookie.Name] = cookie.Value
	}
	return ret
}
