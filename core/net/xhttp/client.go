package xhttp

import (
	"bytes"
	"context"
	"github.com/LuoHongLiang0921/kuaigo/core/encoding"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	nurl "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/uuid"
)

var (
	// HTTPNoKeepAliveClient is http client without keep alive
	HTTPNoKeepAliveClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	defaultHTTPClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 2048,
			IdleConnTimeout:     time.Minute * 5,
		},
	}
	defaultTimeout    = 500
	defaultRetryCount = 2
)

// PostRaw PostRaw
func PostRaw(ctx context.Context, client *http.Client, url string, header http.Header, reqBody interface{}, params ...int) ([]byte, error) {
	var (
		data []byte
		err  error
	)
	if com, ok := klog.FromContext(ctx); ok {
		header.Set(constant.HeaderFieldAi, strconv.Itoa(com.AppId))
		traceId := com.TraceId
		if traceId == "" {
			traceId = uuid.New()
		}
		header.Set(constant.HeaderFieldTi, traceId)
	}
	timeOut, retryCount := genDefaultParams(params...)
	for i := 0; i < retryCount; i++ {
		data, err = do(client, http.MethodPost, url, header, reqBody, timeOut)
		if err == nil {
			break
		}
	}
	if err != nil {
		//xlog.Error("PostRaw").Stack().Msgf("err:%s", err)
		klog.KuaigoLogger.Error(ctx, "PostRaw", klog.FieldCommon(struct{}{}), klog.FieldParams(map[string]interface{}{
			"name": "PostRaw",
		}), klog.FieldErr(err))
	}
	return data, err
}

// PostFormRaw PostRaw
func PostFormRaw(ctx context.Context, client *http.Client, url string, header http.Header, reqBody map[string]string, params ...int) ([]byte, error) {
	var (
		data []byte
		err  error
	)
	if com, ok := klog.FromContext(ctx); ok {
		header.Set(constant.HeaderFieldAi, strconv.Itoa(com.AppId))
		traceId := com.TraceId
		if traceId == "" {
			traceId = uuid.New()
		}
		header.Set(constant.HeaderFieldTi, traceId)
	}
	timeOut, retryCount := genDefaultParams(params...)
	for i := 0; i < retryCount; i++ {
		data, err = doForm(client, http.MethodPost, url, header, reqBody, timeOut)
		if err == nil {
			break
		}
	}
	if err != nil {
		//xlog.Error("PostRaw").Stack().Msgf("err:%s", err)
		klog.KuaigoLogger.Error(ctx, "PostRaw", klog.FieldCommon(struct{}{}), klog.FieldParams(map[string]interface{}{
			"name": "PostRaw",
		}), klog.FieldErr(err))
	}
	return data, err
}

// PostWithUnmarshal do http get with unmarshal
func PostWithUnmarshal(ctx context.Context, client *http.Client, url string, header http.Header, reqBody interface{}, resp interface{}, params ...int) error {
	data, err := PostRaw(ctx, client, url, header, reqBody, params...)
	if err != nil {
		return err
	}
	// for no resp needed request.
	if resp == nil {
		return nil
	}
	// for big int
	decoder := encoding.JSON.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	err = decoder.Decode(resp)
	if err != nil {
		klog.KuaigoLogger.Error(ctx, "PostWithUnmarshal.Decode", klog.FieldParams(map[string]interface{}{
			"url":      url,
			"respData": string(data),
		}), klog.FieldErr(err))
	}
	return err
}

// PostWithUnmarshal do http get with unmarshal
func PostFormWithUnmarshal(ctx context.Context, client *http.Client, url string, header http.Header, reqBody map[string]string, resp interface{}, params ...int) error {
	data, err := PostFormRaw(ctx, client, url, header, reqBody, params...)
	if err != nil {
		return err
	}
	// for no resp needed request.
	if resp == nil {
		return nil
	}
	// for big int
	decoder := encoding.JSON.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	err = decoder.Decode(resp)
	if err != nil {
		klog.KuaigoLogger.Error(ctx, "PostWithUnmarshal.Decode", klog.FieldParams(map[string]interface{}{
			"url":      url,
			"respData": string(data),
		}), klog.FieldErr(err))
	}
	return err
}

// GetRaw get http raw
func GetRaw(ctx context.Context, client *http.Client, url string, header http.Header, reqBody interface{}, params ...int) ([]byte, error) {
	var (
		data []byte
		err  error
	)
	if com, ok := klog.FromContext(ctx); ok {
		header.Set(constant.HeaderFieldAi, strconv.Itoa(com.AppId))
		traceId := com.TraceId
		if traceId == "" {
			traceId = uuid.New()
		}
		header.Set(constant.HeaderFieldTi, traceId)
	}
	timeOut, retryCount := genDefaultParams(params...)
	for i := 0; i < retryCount; i++ {
		data, err = do(client, http.MethodGet, url, header, reqBody, timeOut)
		if err == nil {
			break
		}
	}
	if err != nil {
		//xlog.TabbyLogger.Error("GetRaw").Stack().Msgf("err:%s", err)
		klog.KuaigoLogger.Error(ctx, "GetRaw", klog.FieldCommon(struct{}{}), klog.FieldParams(map[string]interface{}{
			"name": "GetRaw",
		}), klog.FieldErr(err))
	}
	return data, err
}

// GetWithUnmarshal do http get with unmarshal
func GetWithUnmarshal(ctx context.Context, client *http.Client, url string, header http.Header, reqBody interface{}, resp interface{}, params ...int) error {
	data, err := GetRaw(ctx, client, url, header, reqBody, params...)
	if err != nil {
		return err
	}
	// for no resp needed request.
	if resp == nil {
		return nil
	}
	// for big int
	decoder := encoding.JSON.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	err = decoder.Decode(resp)
	if err != nil {
		//log.Error("GetWithUnmarshal.Decode").Stack().Msgf("err:%s, url:%s, respData:%s", err, url, string(data))
		klog.KuaigoLogger.Error(ctx, "GetWithUnmarshal.Decode", klog.FieldParams(map[string]interface{}{
			"name":     "GetWithUnmarshal.Decode",
			"url":      url,
			"respData": string(data),
		}), klog.FieldErr(err))
	}
	return err
}

func genDefaultParams(params ...int) (int, int) {
	timeOut, retryCount := defaultTimeout, defaultRetryCount
	switch {
	case len(params) >= 2:
		timeOut, retryCount = params[0], params[1]
	case len(params) >= 1:
		timeOut = params[0]
	}
	return timeOut, retryCount
}

func do(client *http.Client, method string, url string, header http.Header, reqBody interface{}, timeOut int) ([]byte, error) {
	if client == nil {
		client = defaultHTTPClient
	}
	var reader io.Reader
	switch v := reqBody.(type) {
	case nurl.Values:
		reader = strings.NewReader(v.Encode())
	case []byte:
		reader = bytes.NewBuffer(v)
	case string:
		reader = strings.NewReader(v)
	case io.Reader:
		reader = v
	default:
		buff := &bytes.Buffer{}
		err := encoding.JSON.NewEncoder(buff).Encode(v)
		if err != nil {
			return nil, err
		}
		reader = buff
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	if header != nil {
		req.Header = header
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeOut))
	defer cancelFunc()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err // TODO maybe should define ctx timeout in package errs
	}
	defer resp.Body.Close()
	// TODO maybe should handle status not equal 200
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func doForm(client *http.Client, method string, url string, header http.Header, reqBody map[string]string, timeOut int) ([]byte, error) {
	if client == nil {
		client = defaultHTTPClient
	}
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range reqBody {
		w.WriteField(k, v)
	}
	w.Close()
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if err != nil {
		return nil, err
	}
	//if header != nil {
	//	req.Header = header
	//}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeOut))
	defer cancelFunc()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err // TODO maybe should define ctx timeout in package errs
	}
	defer resp.Body.Close()
	// TODO maybe should handle status not equal 200
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
