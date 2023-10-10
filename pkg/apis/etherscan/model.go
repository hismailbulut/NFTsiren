package etherscan

import (
	"encoding/json"
	"fmt"
	"nftsiren/pkg/number"
)

type EtherscanCommonResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  T      `json:"result"`
	Error   error  `json:"-"`
}

// Etherscan api may return different types for same field in json when error
// we need to check for errors before unmarshalling actual data
func (resp *EtherscanCommonResponse[T]) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		return nil
	}
	var pre map[string]json.RawMessage
	if err := json.Unmarshal(data, &pre); err != nil {
		return err
	}
	if status, ok := pre["status"]; ok {
		if err := json.Unmarshal(status, &resp.Status); err != nil {
			return err
		}
	}
	if message, ok := pre["message"]; ok {
		if err := json.Unmarshal(message, &resp.Message); err != nil {
			return err
		}
	}
	// An API call that encounters an error will return 0 as its status code and display the cause of the error under the result field.
	if resp.Status == "0" {
		errstr := "unknown error"
		if result, ok := pre["result"]; ok {
			if err := json.Unmarshal(result, &errstr); err != nil {
				return err
			}
		}
		resp.Error = fmt.Errorf("%s: %s", resp.Message, errstr)
	} else { // No error
		if result, ok := pre["result"]; ok {
			if err := json.Unmarshal(result, &resp.Result); err != nil {
				return err
			}
		}
	}
	return nil
}

type EthPrice struct {
	Ethbtc          number.Number `json:"ethbtc"`
	EthbtcTimestamp number.Number `json:"ethbtc_timestamp"`
	Ethusd          number.Number `json:"ethusd"`
	EthusdTimestamp number.Number `json:"ethusd_timestamp"`
}

type GasPrice struct {
	LastBlock       number.Number `json:"LastBlock"`
	SafeGasPrice    number.Number `json:"SafeGasPrice"`
	ProposeGasPrice number.Number `json:"ProposeGasPrice"`
	FastGasPrice    number.Number `json:"FastGasPrice"`
	SuggestBaseFee  number.Number `json:"suggestBaseFee"`
	GasUsedRatio    string        `json:"gasUsedRatio"`
}
