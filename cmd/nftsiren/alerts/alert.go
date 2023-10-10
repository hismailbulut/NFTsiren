package alerts

import "nftsiren/pkg/number"

type Condition interface {
	// This can be shown to user
	String() string
	// Label of the value needed for this condition
	Label() string
	// Whether this condition requires a specific interval
	NeedsInterval() bool
}

var _ Condition = CollectionAlertType(0)
var _ Condition = EthereumAlertType(0)
var _ Condition = GasAlertType(0)

// Alert must be immutable
type Alert interface {
	// This is for debugging and for checking equality
	String() string
	// This will return a description of what this alert does, can be showed to user
	Description() string
	// Whether this alert needs an interval input from user
	NeedsInterval() bool
	// This is the time between each check, will be 0 if alert requires no interval
	// Only be trusted when alert needs interval
	Interval() int
	// Reports true when this alert needs to be checked constantly or false when one time only
	Looping() bool
	// Check will check alert Base number against num and returns true if alert condition is met, does not check for interval
	Check(num2 number.Number) bool
	// This text can be shown to user in a notification
	NotificationText() string
}

var _ Alert = &CollectionAlert{}
var _ Alert = &EthereumAlert{}
var _ Alert = &GasAlert{}
