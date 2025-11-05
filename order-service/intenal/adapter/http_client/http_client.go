package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"order-service/config"
	"os"
	"time"

	"github.com/labstack/gommon/log"
)

type HttpClient interface {
	Connect()
	CallURL(method, url string, header map[string]string, rawData []byte) (*http.Response, error)
}

type Options struct {
	timeout int
	http    *http.Client
	logger  *log.Logger
}

type loggingTransport struct {
	logger *log.Logger
}

func NewHttpClient(cfg *config.Config) HttpClient {
	opt := new(Options)
	opt.timeout = cfg.App.ServerTimeOut
	return opt
}

func (o *Options) Connect() {
	now := time.Now().Format("2006-01-02")
	file, err := os.Create("./logs/" + now + ".log")
	if err != nil {
		log.Errorf("[FAILED] Create file logger : %v", err)
	}

	logger := log.New(file, "", log.LstdFlags)

	httpClient := &http.Client{
		Timeout:   time.Duration(o.timeout) * time.Second,
		Transport: &loggingTransport{logger: logger},
	}

	o.http = httpClient
	o.logger = logger
}

func (o *Options) CallURL(method, url string, header map[string]string, rawData []byte) (*http.Response, error) {
	o.Connect()
	req, err := http.NewRequest(method, url, bytes.NewBuffer(rawData))
	if err != nil {
		log.Errorf("[HttpClient-1] CallURL: %v", err)
		return nil, err
	}

	if len(header) > 0 {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	resp, err := o.http.Do(req)
	if err != nil {
		log.Errorf("[HttpClient-2] CallURL: %v", err)
		return nil, err
	}

	return resp, nil
}

func (lt *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Logging sebelum request
	lt.logger.Printf("Making request to: %s %s", req.Method, req.URL)
	lt.logger.Printf("Request Headers: %+v", req.Header)

	// Mengganti request body karena sudah dibaca dalam fungsi logging
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	lt.logger.Printf("Request Body: %s", reqBody)

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		lt.logger.Printf("Request failed: %v", err)
		return nil, err
	}

	// Logging setelah menerima respons
	lt.logger.Printf("Received response with status: %s", resp.Status)
	lt.logger.Printf("Response Headers: %+v", resp.Header)

	// Menampilkan Response Body (jika ada)
	respBody, err := io.ReadAll(resp.Body)
	if err == nil {
		lt.logger.Printf("Response Body: %s", respBody)
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

	return resp, nil
}
