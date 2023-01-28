package gin

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

type IResponse interface {
	IJson(obj interface{}) IResponse

	IJsonp(obj interface{}) IResponse

	IXml(obj interface{}) IResponse

	IHtml(file string, obj interface{}) IResponse

	IText(format string, values ...interface{}) IResponse

	IRedirect(path string) IResponse

	ISetHeader(key string, val string) IResponse

	ISetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	ISetStatus(code int) IResponse

	ISetOkStatus() IResponse
}

func (c *Context) IJson(obj interface{}) IResponse {
	json, err := json.Marshal(obj)
	if err != nil {
		return c.ISetStatus(http.StatusInternalServerError)
	}
	c.ISetHeader("Content-Type", "application/json")
	c.Writer.Write(json)
	return c
}

func (c *Context) IJsonp(obj interface{}) IResponse {
	callbackFunc, _ := c.DefaultQueryString("callback", "callback_func")
	c.ISetHeader("Content-Type", "application/javascript")
	callbackFunc = template.JSEscapeString(callbackFunc)

	_, err := c.Writer.Write([]byte(callbackFunc))
	if err != nil {
		return c
	}
	_, err = c.Writer.Write([]byte("("))
	if err != nil {
		return c
	}
	json, err := json.Marshal(obj)
	if err != nil {
		return c
	}
	_, err = c.Writer.Write(json)
	if err != nil {
		return c
	}
	_, err = c.Writer.Write([]byte(")"))
	if err != nil {
		return c
	}
	return c
}

func (c *Context) IXml(obj interface{}) IResponse {
	xml, err := xml.Marshal(obj)
	if err != nil {
		return c.ISetStatus(http.StatusInternalServerError)
	}
	c.ISetHeader("Content-Type", "application/xml")
	c.Writer.Write(xml)
	return c
}

func (c *Context) IHtml(file string, obj interface{}) IResponse {
	files, err := template.New("output").ParseFiles(file)
	if err != nil {
		return c
	}
	if err := files.Execute(c.Writer, obj); err != nil {
		return c
	}
	c.ISetHeader("Content-Type", "application/html")
	return c
}

func (c *Context) IText(format string, values ...interface{}) IResponse {
	text := fmt.Sprintf(format, values...)
	c.ISetHeader("Content-Type", "application/text")
	c.Writer.Write([]byte(text))
	return c
}

func (c *Context) IRedirect(path string) IResponse {
	http.Redirect(c.Writer, c.Request, path, http.StatusMovedPermanently)
	return c
}

func (c *Context) ISetHeader(key string, val string) IResponse {
	c.Writer.Header().Add(key, val)
	return c
}

func (c *Context) ISetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return c
}

func (c *Context) ISetStatus(code int) IResponse {
	c.Writer.WriteHeader(code)
	return c
}

func (c *Context) ISetOkStatus() IResponse {
	c.Writer.WriteHeader(http.StatusOK)
	return c
}
