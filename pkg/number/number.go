package number

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

const nullString = "null"

// A Number is an abstraction of *big.Float
// Number should always copied, do not pass by reference because it just holds a *big.Float
// You can check whether the number is null by calling IsNil
// All operations on a number is safe even when the number is nil, nil number is treated as zero
// Copying a number doesn't copy the underlying value, use Copy to actually copy a number
type Number struct {
	// A value can be nil
	value *big.Float
}

func (num Number) IsNil() bool {
	return num.value == nil
}

// This will return true if number is nil or has zero value
func (num Number) IsZero() bool {
	if num.IsNil() {
		return true
	}
	if num.value.Cmp(new(big.Float)) == 0 {
		return true
	}
	return false
}

// All operations has to use the number.value should call this before using it
func (num *Number) checknil() {
	if num.value == nil {
		num.value = new(big.Float)
	}
}

func NewFromInt(i int64) Number {
	return Number{new(big.Float).SetInt64(i)}
}

func NewFromFloat(f float64) Number {
	return Number{new(big.Float).SetFloat64(f)}
}

func NewFromBigInt(b *big.Int) Number {
	return Number{new(big.Float).SetInt(b)}
}

func NewFromString(s string) (Number, bool) {
	num := Number{}
	ok := num.SetString(s)
	return num, ok
}

func (num *Number) SetString(s string) bool {
	v, ok := new(big.Float).SetString(s)
	if ok {
		num.value = v
		return true
	}
	return false
}

func (num Number) Copy() Number {
	if num.IsNil() {
		return Number{}
	}
	return Number{new(big.Float).Set(num.value)}
}

func (num Number) Add(num1 Number) Number {
	num.checknil()
	num1.checknil()
	return Number{new(big.Float).Add(num.value, num1.value)}
}

func (num Number) AddInt64(i int64) Number {
	return num.Add(NewFromInt(i))
}

func (num Number) AddFloat64(f float64) Number {
	return num.Add(NewFromFloat(f))
}

func (num Number) Sub(num1 Number) Number {
	num.checknil()
	num1.checknil()
	return Number{new(big.Float).Sub(num.value, num1.value)}
}

func (num Number) SubInt64(i int64) Number {
	return num.Sub(NewFromInt(i))
}

func (num Number) SubFloat64(f float64) Number {
	return num.Sub(NewFromFloat(f))
}

func (num Number) Mul(num1 Number) Number {
	num.checknil()
	num1.checknil()
	return Number{new(big.Float).Mul(num.value, num1.value)}
}

func (num Number) MulInt64(i int64) Number {
	return num.Mul(NewFromInt(i))
}

func (num Number) MulFloat64(f float64) Number {
	return num.Mul(NewFromFloat(f))
}

func (num Number) Div(num1 Number) Number {
	num.checknil()
	num1.checknil()
	return Number{new(big.Float).Quo(num.value, num1.value)}
}

func (num Number) DivInt64(i int64) Number {
	return num.Div(NewFromInt(i))
}

func (num Number) DivFloat64(f float64) Number {
	return num.Div(NewFromFloat(f))
}

func (num Number) Cmp(num1 Number) int {
	num.checknil()
	num1.checknil()
	return num.value.Cmp(num1.value)
}

func (num Number) Equals(num1 Number) bool {
	return num.Cmp(num1) == 0
}

func (num Number) LessThan(num1 Number) bool {
	return num.Cmp(num1) < 0
}

func (num Number) LessThanOrEqual(num1 Number) bool {
	return num.Cmp(num1) <= 0
}

func (num Number) GreaterThan(num1 Number) bool {
	return num.Cmp(num1) > 0
}

func (num Number) GreaterThanOrEqual(num1 Number) bool {
	return num.Cmp(num1) >= 0
}

// num1 must be an integer, it panics otherwise
// TODO: find a way for better Pow
func (num Number) Pow(num1 Number) Number {
	num.checknil()
	num1.checknil()
	if !num1.IsInteger() {
		panic("number: pow accepts only integer values")
	}
	if num.IsInteger() {
		return NewFromBigInt(new(big.Int).Exp(num.BigInt(), num1.BigInt(), nil))
	} else {
		result := new(big.Float).Set(num.value)
		for i := int64(0); i < num1.Int64()-1; i++ {
			result.Mul(result, result)
		}
		return Number{result}
	}
}

func (num Number) Sqrt() Number {
	num.checknil()
	return Number{new(big.Float).Sqrt(num.value)}
}

func (num Number) IsInteger() bool {
	num.checknil()
	return num.value.IsInt()
}

func (num Number) BigInt() *big.Int {
	num.checknil()
	i, _ := num.value.Int(nil)
	return i
}

func (num Number) BigFloat() *big.Float {
	num.checknil()
	return new(big.Float).Set(num.value)
}

func (num Number) Int64() int64 {
	num.checknil()
	i, _ := num.value.Int64()
	return i
}

func (num Number) Float64() float64 {
	num.checknil()
	f, _ := num.value.Float64()
	return f
}

// will return null if number is null
func (num Number) String() string {
	if num.IsNil() {
		return nullString
	}
	return num.StringFixed(-1)
}

// will return null if number is null
func (num Number) StringFixed(prec int) string {
	if num.IsNil() {
		return nullString
	}
	return num.value.Text('f', prec)
}

// will return null if number is null
func (num Number) StringPretty() string {
	if num.IsNil() {
		return nullString
	}
	fix := func(v Number, postfix string, places int) string {
		if v.IsInteger() {
			return v.String() + postfix
		}
		s := v.StringFixed(places)
		if s[len(s)-1] == '0' {
			// TODO: find a better and faster way
			s = strings.TrimRight(s, "0")
			s = strings.TrimRight(s, ".")
		}
		return s + postfix
	}
	e9 := NewFromInt(1e9)
	if num.GreaterThanOrEqual(e9) {
		return fix(num.Div(e9), "b", 1) // billion
	}
	e6 := NewFromInt(1e6)
	if num.GreaterThanOrEqual(e6) {
		return fix(num.Div(e6), "m", 1) // million
	}
	e5 := NewFromInt(1e5)
	e3 := NewFromInt(1e3)
	if num.GreaterThanOrEqual(e5) { // this is correct, 9999 -> 10k
		return fix(num.Div(e3), "k", 1) // thousand
	}
	// 1000 > num >= 100
	e2 := NewFromInt(1e2)
	if num.GreaterThanOrEqual(e2) { // 111.1
		return fix(num, "", 1)
	}
	// 100 > num >= 10
	e1 := NewFromInt(1e1)
	if num.GreaterThanOrEqual(e1) { // 11.11
		return fix(num, "", 2)
	}
	// 10 > num
	return fix(num, "", 3) // 1.111
}

func (num Number) MarshalText() ([]byte, error) {
	if num.IsNil() {
		return []byte(nullString), nil
	}
	return []byte(num.String()), nil
}

func (num *Number) UnmarshalText(data []byte) error {
	text := string(data)
	if text == nullString {
		return nil
	}
	ok := num.SetString(text)
	if !ok {
		return fmt.Errorf("can't unmarshal %s into a number", text)
	}
	return nil
}

func (num Number) MarshalJSON() ([]byte, error) {
	if num.IsNil() {
		return []byte(nullString), nil
	}
	return json.Marshal(num.String())
}

func (num *Number) UnmarshalJSON(data []byte) error {
	text := string(data)
	if text == nullString {
		return nil
	}
	ok := num.SetString(strings.Trim(text, `"`))
	if !ok {
		return fmt.Errorf("can't unmarshal %s into a number", text)
	}
	return nil
}
