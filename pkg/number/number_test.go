package number

import (
	"encoding/json"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNumber(t *testing.T) {
	const testInt int64 = 1234567890
	const testFloatStr string = "1234567890.0123456789"
	numFromInt64 := NewFromInt(testInt)
	numFromBigInt := NewFromBigInt(new(big.Int).SetInt64(testInt))
	numFromString, ok := NewFromString(testFloatStr)
	assert.True(t, ok)
	assert.Equal(t, testInt, numFromInt64.Int64())
	assert.Equal(t, testInt, numFromBigInt.Int64())
	assert.Equal(t, numFromInt64.String(), numFromBigInt.String())
	assert.Equal(t, testFloatStr, numFromString.String())
}

func TestArithmeticOperations(t *testing.T) {
	f, ok := NewFromString("1453.1071")
	assert.True(t, ok)
	assert.Equal(t, "1454.2071", f.AddFloat64(1.1).StringFixed(4))
	assert.Equal(t, "2154.1071", f.AddInt64(701).String())
	assert.Equal(t, "1449.9671", f.SubFloat64(3.14).StringFixed(4))
	assert.Equal(t, "0.1071", f.SubInt64(1453).StringFixed(4))
	assert.Equal(t, "1", f.Div(f).String())
	assert.Equal(t, f.String(), f.MulInt64(1).String())
	assert.Equal(t, f.String(), f.Div(f).Mul(f).String())
	assert.Equal(t, "100", NewFromInt(10).Pow(NewFromInt(2)).String())
	f2, ok := NewFromString("3.13")
	assert.True(t, ok)
	assert.Equal(t, "9.7969", f2.Pow(NewFromInt(2)).StringFixed(4))
	f3, ok := NewFromString("4")
	assert.True(t, ok)
	assert.Equal(t, "2", f3.Sqrt().String())
	assert.Equal(t, "4", f3.String()) // Original number must still same
}

func TestComparisons(t *testing.T) {
	s, ok1 := NewFromString("1.10003")
	b, ok2 := NewFromString("1.10004")
	assert.True(t, ok1 && ok2)
	assert.True(t, s.LessThan(b))
	assert.False(t, s.GreaterThan(b))
	assert.False(t, b.LessThan(s))
	assert.True(t, b.GreaterThan(s))
	assert.True(t, s.Equals(s))
	assert.True(t, b.Equals(b))
}

func TestPrettyDecimal(t *testing.T) {
	assert.Equal(t, "193.6k", NewFromFloat(193586.9876).StringPretty())
	assert.Equal(t, "193", NewFromFloat(193.0004).StringPretty())
	assert.Equal(t, "193m", NewFromFloat(193000000.98).StringPretty())
	assert.Equal(t, "193.1m", NewFromFloat(193091456.98).StringPretty())
	assert.Equal(t, "19365.1", NewFromFloat(19365.1234).StringPretty())
	assert.Equal(t, "193.1", NewFromFloat(193.1234).StringPretty())
	assert.Equal(t, "19.12", NewFromFloat(19.12345).StringPretty())
	assert.Equal(t, "1.123", NewFromFloat(1.123456).StringPretty())
}

type testStruct struct {
	N Number `json:"num"`
}

func TestJSON(t *testing.T) {
	{
		// Marshal
		num, ok := NewFromString("9.998")
		assert.True(t, ok)
		ts1 := testStruct{num}
		data, err := json.Marshal(ts1)
		assert.Nil(t, err)
		assert.Equal(t, `{"num":"9.998"}`, string(data))
		// Unmarshal
		var ts2 testStruct
		err = json.Unmarshal(data, &ts2)
		assert.Nil(t, err)
		assert.False(t, ts2.N.IsNil())
		assert.True(t, ts1.N.Equals(ts2.N))
	}
	{
		// Nil marshal
		var ts1 testStruct
		data, err := json.Marshal(ts1)
		assert.Nil(t, err)
		assert.Equal(t, `{"num":null}`, string(data))
		// Nil unmarshal
		var ts2 testStruct
		err = json.Unmarshal(data, &ts2)
		assert.Nil(t, err)
		assert.True(t, ts2.N.IsNil())
	}
	{
		// Unmarshal int
		var ts1 testStruct
		data := []byte(`{"num":1.618}`)
		err := json.Unmarshal(data, &ts1)
		assert.Nil(t, err)
		assert.False(t, ts1.N.IsNil())
		ns, ok := NewFromString("1.618")
		assert.True(t, ok)
		assert.True(t, ts1.N.Equals(ns))
	}
}

func BenchmarkAdd(b *testing.B) {
	num1 := NewFromInt(rand.Int63())
	num2 := NewFromFloat(rand.NormFloat64())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		num1.Add(num2)
	}
}

func BenchmarkSub(b *testing.B) {
	num1 := NewFromInt(rand.Int63())
	num2 := NewFromFloat(rand.NormFloat64())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		num1.Sub(num2)
	}
}

func BenchmarkMul(b *testing.B) {
	num1 := NewFromInt(rand.Int63())
	num2 := NewFromFloat(rand.NormFloat64())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		num1.Mul(num2)
	}
}

func BenchmarkDiv(b *testing.B) {
	num1 := NewFromInt(rand.Int63())
	num2 := NewFromFloat(rand.NormFloat64())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		num1.Div(num2)
	}
}

func BenchmarkPowInt(b *testing.B) {
	num1 := NewFromInt(rand.Int63n(100))
	num2 := NewFromInt(rand.Int63n(100))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		num1.Pow(num2)
	}
}

func BenchmarkPowFloat(b *testing.B) {
	num1 := NewFromFloat(rand.Float64() * float64(rand.Int31n(100)))
	num2 := NewFromInt(rand.Int63n(100))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		num1.Pow(num2)
	}
}

func BenchmarkString(b *testing.B) {
	num := NewFromFloat(rand.NormFloat64())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = num.String()
	}
}

func BenchmarkStringPretty(b *testing.B) {
	num := NewFromFloat(rand.NormFloat64())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = num.StringPretty()
	}
}
