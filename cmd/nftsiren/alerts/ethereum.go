package alerts

import (
	"fmt"

	"nftsiren/pkg/number"
)

type EthereumAlertType int

const (
	EthereumAlertTypeLessThan EthereumAlertType = iota
	EthereumAlertTypeGreaterThan
)

func (t EthereumAlertType) String() string {
	switch t {
	case EthereumAlertTypeLessThan:
		return "Eth Less Than"
	case EthereumAlertTypeGreaterThan:
		return "Eth Greater Than"
	}
	return "UNKNOWN"
}

func (t EthereumAlertType) Label() string {
	return "Eth Price (USD)"
}

func (t EthereumAlertType) NeedsInterval() bool {
	return false
}

type EthereumAlert struct {
	Type EthereumAlertType `json:"type" bson:"type"`
	Base number.Number     `json:"base" bson:"base"`
	Loop bool              `json:"loop" bson:"loop"`
}

func (alert EthereumAlert) String() string {
	return fmt.Sprintf("%v:%v:%v", alert.Type, alert.Base, alert.Loop)
}

func (alert EthereumAlert) Description() string {
	switch alert.Type {
	case EthereumAlertTypeLessThan:
		return fmt.Sprintf("Checks for ETH$ < %v", alert.Base)
	case EthereumAlertTypeGreaterThan:
		return fmt.Sprintf("Checks for ETH$ > %v", alert.Base)
	}
	return "UNKNOWN"
}

func (alert EthereumAlert) NeedsInterval() bool {
	return false
}

func (alert EthereumAlert) Interval() int {
	return 0
}

func (alert EthereumAlert) Looping() bool {
	return alert.Loop
}

// returns true when check passed
func (alert EthereumAlert) Check(num number.Number) bool {
	switch alert.Type {
	case EthereumAlertTypeLessThan:
		return num.LessThan(alert.Base)
	case EthereumAlertTypeGreaterThan:
		return num.GreaterThan(alert.Base)
	}
	return false
}

func (alert EthereumAlert) NotificationText() string {
	switch alert.Type {
	case EthereumAlertTypeLessThan:
		return fmt.Sprintf("ETH less than %v", alert.Base)
	case EthereumAlertTypeGreaterThan:
		return fmt.Sprintf("ETH greater than %v", alert.Base)
	}
	return "UNKNOWN"
}
