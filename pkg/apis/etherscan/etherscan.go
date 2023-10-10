package etherscan

import (
	"errors"
	"fmt"
	"net/http"

	"nftsiren/pkg/httpclient"
	"nftsiren/pkg/mutex"
)

var client = httpclient.NewClient("https://api.etherscan.io/api")
var apiKey mutex.Value[string]

func SetApiKey(key string) {
	apiKey.Store(key)
}

func get[T any](module, action string) (EtherscanCommonResponse[T], error) {
	params := map[string]string{
		"module": module,
		"action": action,
	}
	if apiKey.Load() != "" {
		params["apikey"] = apiKey.Load()
	}
	var resp EtherscanCommonResponse[T]
	status, err := client.GetJson(nil, params, &resp)
	if err != nil {
		if status >= 400 {
			return EtherscanCommonResponse[T]{}, errors.New(http.StatusText(status))
		}
		return EtherscanCommonResponse[T]{}, err
	}
	if resp.Error != nil {
		return EtherscanCommonResponse[T]{}, resp.Error
	}
	return resp, nil
}

func FetchEthPrice() (EthPrice, error) {
	resp, err := get[EthPrice]("stats", "ethprice")
	if err != nil {
		return EthPrice{}, err
	}
	if resp.Status != "1" {
		return EthPrice{}, fmt.Errorf(resp.Message)
	}
	return resp.Result, nil
}

func FetchGasPrice() (GasPrice, error) {
	resp, err := get[GasPrice]("gastracker", "gasoracle")
	if err != nil {
		return GasPrice{}, err
	}
	if resp.Status != "1" {
		return GasPrice{}, fmt.Errorf(resp.Message)
	}
	return resp.Result, nil
}
