package httpclient

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"nftsiren/pkg/log"
	"nftsiren/pkg/mutex"

	"github.com/beefsack/go-rate"
)

const (
	TIMEOUT = 30 * time.Second
)

type Client struct {
	handle    *http.Client
	baseUrl   *url.URL
	rateLimit *rate.RateLimiter
	headers   *mutex.Map[string, string]
	inQueue   mutex.Counter
}

func NewClient(baseUrl string) *Client {
	return &Client{
		handle: &http.Client{
			Timeout: TIMEOUT,
		},
		baseUrl: mustParseURL(baseUrl),
		headers: mutex.NewMap[string, string](),
	}
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

func NewClientWithLimit(baseUrl string, limit int, interval time.Duration) *Client {
	client := NewClient(baseUrl)
	client.rateLimit = rate.New(limit, interval)
	return client
}

// Default headers will be added to every request
func (client *Client) SetDefaultHeader(key, value string) {
	client.headers.Store(key, value)
}

func (client *Client) DeleteDefaultHeader(key string) {
	client.headers.Delete(key)
}

// This will also clear authorization
func (client *Client) ClearDefaultHeaders() {
	client.headers.Clear()
}

func (client *Client) SetBasicAuth(username, password string) {
	auth := []byte(username + ":" + password)
	token := base64.StdEncoding.EncodeToString(auth)
	client.SetDefaultHeader("Authorization", "Basic "+token)
}

func (client *Client) SetBearerAuth(token string) {
	client.SetDefaultHeader("Authorization", "Bearer "+token)
}

func (client *Client) DisableAuth() {
	client.DeleteDefaultHeader("Authorization")
}

func (client *Client) DoRequest(req *Request) (*http.Response, error) {
	if client.rateLimit != nil {
		client.inQueue.Increment()
		if client.inQueue.Value() >= 10 {
			log.Warn().Println("Too many requests waiting for rate limit")
		}
		client.rateLimit.Wait()
		client.inQueue.Decrement()
	}
	// Build url
	reqURL := req.URL(client.baseUrl)
	// log.Debug().Println(req.Method, reqURL)
	// Create request
	httpReq, err := http.NewRequest(req.Method, reqURL, req.Payload)
	if err != nil {
		return nil, err
	}
	// Set default headers before request headers because request headers
	// should be able to override defaults
	client.headers.Range(func(index int, key, value string) {
		httpReq.Header.Set(key, value)
	})
	// Set request headers
	for k, v := range req.Header {
		httpReq.Header.Set(k, v)
	}
	// Do request
	return client.handle.Do(httpReq)
}

// Makes a GET request and returns response
func (client *Client) Get(path ...string) (*http.Response, error) {
	return client.DoRequest(NewRequest(http.MethodGet, path))
}

// respObjRef should be reference to an object
// Status code may zero if there is a network error, also may return json error
func (client *Client) GetJson(path []string, params map[string]string, respObjRef any) (int, error) {
	req := NewRequest(http.MethodGet, path).SetAcceptJSON()
	req.Params = params
	resp, err := client.DoRequest(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(respObjRef)
	// This is for debugging
	/*
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		log.Debug().
			Field("body", string(body)).
			Println("Response of the request", req)
		return resp.StatusCode, json.Unmarshal(body, respObjRef)
	*/
}

// Post makes a POST request and returns response
func (client *Client) Post(path []string, payload []byte, contentType string) (*http.Response, error) {
	req := NewRequest(http.MethodPost, path).SetPayloadBytes(payload).SetContentType(contentType)
	return client.DoRequest(req)
}

// respObjRef should be reference of an object
// This function takes an object and posts it as json and also expects a json object from server
// Status code may zero if there is a network error, also may return json encoding or decoding error
func (client *Client) PostJson(path []string, params map[string]string, bodyObj any, respObjRef any) (int, error) {
	payload, err := json.Marshal(bodyObj)
	if err != nil {
		return 0, err
	}
	req := NewRequest(http.MethodPost, path).SetPayloadBytes(payload).SetContentTypeJSON().SetAcceptJSON()
	resp, err := client.DoRequest(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(respObjRef)
}
