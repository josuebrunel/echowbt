package echowbt

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
)

func concat(ss ...string) string {
	var ns strings.Builder
	for _, s := range ss {
		ns.WriteString(s)
	}
	return ns.String()
}

// Dict represents a dict object
type Dict map[string]interface{}

// DictString represents a dict of string key and values
type DictString map[string]string

// URL is a type representing an URL with :
// * Path
// * Params
// * Values
type URL struct {
	Path        string
	Params      []string
	Values      []string
	QueryString map[string]string
}

func (u URL) String() string {
	if u.QueryString == nil {
		return u.Path
	}
	u.Path = concat(u.Path, "?")
	for k, v := range u.QueryString {
		if u.Path[len(u.Path)-1:] == "?" {
			u.Path = concat(u.Path, k, "=", v)
		} else {
			u.Path = concat(u.Path, "&", k, "=", v)
		}
	}
	return u.Path
}

// NewURL is a function returning an URL struct
// p as Path
// uv as Url Values ( e.g /:username/ )
// qs as Url Query String ( e.g /?debug=true)
func NewURL(p string, uv DictString, qs DictString) (u URL) {
	u = URL{Path: p, QueryString: qs}
	// set url values
	for k, v := range uv {
		u.Params = append(u.Params, k)
		u.Values = append(u.Values, v)
	}
	return
}

// JSONDecode returns an interface from a json formatted
// http.ResponseRecorder.Body
func JSONDecode(b *bytes.Buffer) Dict {
	v := Dict{}
	err := json.Unmarshal(b.Bytes(), &v)
	if err != nil {
		panic(err)
	}
	return v
}

// JSONEncode transforms interface into a []byte
func JSONEncode(s interface{}) []byte {
	out, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return out
}

// MultiPartForm is the struct returned by FormData func
type MultiPartForm struct {
	Data        []byte
	ContentType string
}

// FormData helps create a form data payload
func FormData(fields DictString, files DictString) (MultiPartForm, error) {
	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)
	// write fields
	for k, v := range fields {
		writer.WriteField(k, v)
	}
	// handle file fields
	for fname, fpath := range files {
		file, err := os.Open(fpath)
		if err != nil {
			return MultiPartForm{}, err
		}
		defer file.Close()
		part, err := writer.CreateFormFile(fname, filepath.Base(fpath))
		if err != nil {
			return MultiPartForm{}, err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return MultiPartForm{}, err
		}
	}
	if err := writer.Close(); err != nil {
		return MultiPartForm{}, err
	}
	return MultiPartForm{body.Bytes(), writer.FormDataContentType()}, nil
}

// Client represents the client instance
type Client struct {
	E *echo.Echo
	H DictString
}

// New returns a client instance
func New() (c Client) {
	c = Client{E: echo.New(), H: DictString{"Content-Type": "application/json"}}
	return
}

// SetHeaders allow you define some headers
func (c *Client) SetHeaders(headers DictString) {
	for k, v := range headers {
		c.H[k] = v
	}
}

// Request is the method performing the request
func (c Client) Request(method string, url URL, handler echo.HandlerFunc, data []byte, headers DictString) *httptest.ResponseRecorder {
	methods := map[string]string{
		"get":    http.MethodGet,
		"post":   http.MethodPost,
		"put":    http.MethodPut,
		"delete": http.MethodDelete,
		"patch":  http.MethodPatch,
	}
	path := url.String()
	req := httptest.NewRequest(methods[method], path, bytes.NewReader(data))
	// set client default headers
	for k, v := range c.H {
		req.Header.Set(k, v)
	}
	// set call headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	ctx := c.E.NewContext(req, rec)
	ctx.SetPath(path)
	ctx.SetParamNames(url.Params...)
	ctx.SetParamValues(url.Values...)
	handler(ctx)
	return rec
}

// Get represents a Get Request
func (c Client) Get(url URL, handler echo.HandlerFunc, data []byte, headers DictString) *httptest.ResponseRecorder {
	return c.Request("get", url, handler, data, headers)
}

// Post represents a Post Request
func (c Client) Post(url URL, handler echo.HandlerFunc, data []byte, headers DictString) *httptest.ResponseRecorder {
	return c.Request("post", url, handler, data, headers)
}

// Put represents a Put Request
func (c Client) Put(url URL, handler echo.HandlerFunc, data []byte, headers DictString) *httptest.ResponseRecorder {
	return c.Request("put", url, handler, data, headers)
}

// Patch represents a Patch Request
func (c Client) Patch(url URL, handler echo.HandlerFunc, data []byte, headers DictString) *httptest.ResponseRecorder {
	return c.Request("patch", url, handler, data, headers)
}

// Delete represents a Delete Request
func (c Client) Delete(url URL, handler echo.HandlerFunc, data []byte, headers DictString) *httptest.ResponseRecorder {
	return c.Request("delete", url, handler, data, headers)
}
