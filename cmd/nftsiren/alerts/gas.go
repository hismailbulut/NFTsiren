package alerts

import (
	"fmt"

	"nftsiren/pkg/number"
)

type GasAlertType int

const (
	GasAlertTypeLessThan GasAlertType = iota
	GasAlertTypeGreaterThan
)

func (t GasAlertType) String() string {
	switch t {
	case GasAlertTypeLessThan:
		return "Gas Less Than"
	case GasAlertTypeGreaterThan:
		return "Gas Greater Than"
	}
	return "UNKNOWN"
}

func (t GasAlertType) Label() string {
	return "Gas Price (GWEI)"
}

func (t GasAlertType) NeedsInterval() bool {
	return false
}

type GasAlert struct {
	Type GasAlertType  `json:"type" bson:"type"`
	Base number.Number `json:"base" bson:"base"`
	Loop bool          `json:"loop" bson:"loop"`
}

func (alert GasAlert) String() string {
	return fmt.Sprintf("%v:%v:%v", alert.Type, alert.Base, alert.Loop)
}

func (alert GasAlert) Description() string {
	switch alert.Type {
	case GasAlertTypeLessThan:
		return fmt.Sprintf("Checks for GAS < %v", alert.Base)
	case GasAlertTypeGreaterThan:
		return fmt.Sprintf("Checks for GAS > %v", alert.Base)
	}
	return "UNKNOWN"
}

func (alert GasAlert) NeedsInterval() bool {
	return false
}

func (alert GasAlert) Interval() int {
	return 0
}

func (alert GasAlert) Looping() bool {
	return alert.Loop
}

// returns true when check passed
func (alert GasAlert) Check(num number.Number) bool {
	switch alert.Type {
	case GasAlertTypeLessThan:
		return num.LessThan(alert.Base)
	case GasAlertTypeGreaterThan:
		return num.GreaterThan(alert.Base)
	}
	return false
}

func (alert GasAlert) NotificationText() string {
	switch alert.Type {
	case GasAlertTypeLessThan:
		return fmt.Sprintf("Gas less than %v", alert.Base)
	case GasAlertTypeGreaterThan:
		return fmt.Sprintf("Gas greater than %v", alert.Base)
	}
	return "UNKNOWN"
}
