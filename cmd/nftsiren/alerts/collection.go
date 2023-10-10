package alerts

import (
	"fmt"

	"nftsiren/pkg/number"
)

type CollectionAlertType int

const (
	CollectionAlertTypeFloorLessThan CollectionAlertType = iota
	CollectionAlertTypeFloorGreaterThan
	CollectionAlertTypeSalesGreaterThan
)

func (t CollectionAlertType) String() string {
	switch t {
	case CollectionAlertTypeFloorLessThan:
		return "Floor Less Than"
	case CollectionAlertTypeFloorGreaterThan:
		return "Floor Greater Than"
	case CollectionAlertTypeSalesGreaterThan:
		return "Sales Greater Than"
	}
	return "UNKNOWN"
}

// Returns required parameters to create this type of alert
// inner array contains two values, first one is entry label and seconds one is parameter name
func (t CollectionAlertType) Label() string {
	switch t {
	case CollectionAlertTypeFloorLessThan:
		return "Floor Price (ETH)"
	case CollectionAlertTypeFloorGreaterThan:
		return "Floor Price (ETH)"
	case CollectionAlertTypeSalesGreaterThan:
		return "Total Sales In Interval"
	}
	return "UNKNOWN"
}

func (t CollectionAlertType) NeedsInterval() bool {
	switch t {
	case CollectionAlertTypeFloorLessThan:
		return false
	case CollectionAlertTypeFloorGreaterThan:
		return false
	case CollectionAlertTypeSalesGreaterThan:
		return true
	}
	return false
}

type CollectionAlert struct {
	Type CollectionAlertType `json:"type"     bson:"type"`     // Type of the alert
	Base number.Number       `json:"base"     bson:"base"`     // A number to check against something specified by type
	Intv int                 `json:"interval" bson:"interval"` // The time between each update in seconds
	Loop bool                `json:"loop"     bson:"loop"`     // Is it needs to be checked continiously?
}

func (alert CollectionAlert) String() string {
	return fmt.Sprintf("%v:%v:%d:%v", alert.Type, alert.Base, alert.Intv, alert.Loop)
}

func (alert CollectionAlert) Description() string {
	switch alert.Type {
	case CollectionAlertTypeFloorLessThan:
		return fmt.Sprintf("Checks for Floor < %v", alert.Base)
	case CollectionAlertTypeFloorGreaterThan:
		return fmt.Sprintf("Checks for Floor > %v", alert.Base)
	case CollectionAlertTypeSalesGreaterThan:
		return fmt.Sprintf("Checks for Sales > %v every %d seconds", alert.Base, alert.Intv)
	}
	return "UNKNOWN"
}

func (alert CollectionAlert) NeedsInterval() bool {
	return alert.Type.NeedsInterval()
}

func (alert CollectionAlert) Interval() int {
	return alert.Intv
}

func (alert CollectionAlert) Looping() bool {
	return alert.Loop
}

// returns true when check passed
func (alert CollectionAlert) Check(num number.Number) bool {
	switch alert.Type {
	case CollectionAlertTypeFloorLessThan:
		return num.LessThan(alert.Base)
	case CollectionAlertTypeFloorGreaterThan:
		return num.GreaterThan(alert.Base)
	case CollectionAlertTypeSalesGreaterThan:
		return num.LessThan(alert.Base)
	}
	return false
}

func (alert CollectionAlert) NotificationText() string {
	switch alert.Type {
	case CollectionAlertTypeFloorLessThan:
		return fmt.Sprintf("Floor is less than %v", alert.Base)
	case CollectionAlertTypeFloorGreaterThan:
		return fmt.Sprintf("Floor is greater than %v", alert.Base)
	case CollectionAlertTypeSalesGreaterThan:
		return fmt.Sprintf("Sales passed %v in past %d seconds", alert.Base, alert.Intv)
	}
	return "UNKNOWN"
}
