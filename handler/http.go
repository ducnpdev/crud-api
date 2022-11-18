package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Reponse struct {
	Message string      `json:"message"`
	Err     string      `json:"err"`
	Data    interface{} `json:"data"`
}

func Set() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, Reponse{
			Message: string("set-data"),
		})
	}
}
func Health() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func call(ctx context.Context, mess, url string) Reponse {
	log.Default().Println("message:", mess)
	log.Default().Println("url:", url)
	clientHttp := NewClientHttp()
	h := map[string]string{
		"Content-Type": "application/json",
	}
	req := ClientHttpRequest{
		Body:   mess,
		Method: "GET",
		Url:    url,
		Header: h,
	}
	resp, err := clientHttp.Get(ctx, req)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	var status int
	if resp != nil {
		status = resp.StatusCode
	}
	return Reponse{
		Message: (mess),
		Err:     errStr,
		Data:    status,
	}
}
func Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, call(ctx, "get-aaaa-data", os.Getenv("url_get_a")))
	}
}

func GetB() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, call(ctx, "get-bbbbbb-data", os.Getenv("url_get_b")))
	}
}

type ClientHttpRequest struct {
	Body   interface{}
	Method string
	Url    string
	Header map[string]string
}

const (
	MimeJSON = "application/json"
)

type ClientHttp interface {
	Post(ctx context.Context, req ClientHttpRequest) (*http.Response, error)
	Get(ctx context.Context, req ClientHttpRequest) (*http.Response, error)
}

type clientHttp struct {
	client *http.Client
}

func NewClientHttp() ClientHttp {
	client := &http.Client{Transport: getTransport()}
	return &clientHttp{
		client: client,
	}
}

func getTransport() *http.Transport {
	tr := &http.Transport{
		// MaxIdleConns:       10,
		// IdleConnTimeout:    30 * time.Second,
		// DisableCompression: true,
	}
	return tr
}

// build common header
func buildHeader(mapHeader map[string]string) (header http.Header) {
	header = make(http.Header)
	for key, value := range mapHeader {
		header.Set(key, value)
	}
	return header
}

// build body of api
func buildBody(ctx context.Context, contentType string, bodyReq interface{}) (*bytes.Reader, error) {
	var (
		body     *bytes.Reader
		err      error
		bodyByte []byte
	)
	switch contentType {
	case MimeJSON:
		bodyByte, err = json.Marshal(bodyReq)
	default:
		return body, fmt.Errorf("content type of body only application/json")
	}
	if err != nil {
		return body, err
	}
	body = bytes.NewReader(bodyByte)
	return body, err
}

// build request data of http
func buildRequestHttp(ctx context.Context, req ClientHttpRequest) (*http.Request, error) {
	var (
		httpReq *http.Request
		err     error
	)
	body, err := buildBody(ctx, MimeJSON, req.Body)
	if err != nil {
		return httpReq, err
	}
	httpReq, err = http.NewRequestWithContext(ctx, req.Method, req.Url, body)
	if err != nil {
		return httpReq, err
	}
	httpReq.Header = buildHeader(req.Header)
	return httpReq, err
}

// post api
func (h *clientHttp) Post(ctx context.Context,
	req ClientHttpRequest) (httpResp *http.Response, err error) {
	req.Method = http.MethodPost
	reqhttp, err := buildRequestHttp(ctx, req)
	if err != nil {
		return httpResp, err
	}
	// client http send data
	httpResp, err = h.client.Do(reqhttp)
	return httpResp, err
}

// get api
func (h *clientHttp) Get(ctx context.Context,
	req ClientHttpRequest) (httpResp *http.Response, err error) {
	req.Method = http.MethodGet
	reqhttp, err := buildRequestHttp(ctx, req)
	if err != nil {
		return httpResp, err
	}
	// client http send data
	httpResp, err = h.client.Do(reqhttp)
	return httpResp, err
}
