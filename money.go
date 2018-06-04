// Package money implements an arbitrary precision fixed-point decimal.
//
// Functionality here, with respect to numerical functions is largely
//    a wrapper for github.com/shopspring/decimal.go functions.
//
// The best way to create a new Money is to use money.NewFromString, ex:
//
//     n, err := money.NewFromString("-123.4567")
//     n.String() // output: "-123.4567"
//
package money

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"math/big"
)

// Core Monetary construct which uses shopspring's decimal number and adds a
//  currency to the mix. All mathematical functions use the underlying decimal
//  package, which is designed to be "money safe".
// Note: No currency mixing is allowed. For that we'll create an exchange library.
//  Trying to perform operations (add/subtract/compare/etc) on mixed currency Moneys
//  will panic. YOU HAVE BEEN WARNED.
type Money struct {
	amount   decimal.Decimal
	currency *Currency
}

// DivisionPrecision is the number of decimal places in the result when it
// doesn't divide exactly. Overriding this from the original 16 to 20. Just because.
//
// Example:
//
//     d1 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(3)
//     d1.String() // output: "0.66666666666666666667"
//     d2 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(30000)
//     d2.String() // output: "0.00006666666666666667"
//     d3 := decimal.NewFromFloat(20000).Div(decimal.NewFromFloat(3)
//     d3.String() // output: "6666.66666666666666666667"
//     decimal.DivisionPrecision = 3
//     d4 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(3)
//     d4.String() // output: "0.667"
//
var DivisionPrecision = 20

// Zero constant, to make computations faster.
var ZeroMoney = Money{amount: decimal.Zero, currency: getUnknownCurrency()}

// New returns a new Money of type currency, with an amount of value * 10 ^ exp.
func New(curr string, value int64, exp int32) (Money, error) {

	c, ok := GetCurrency(curr)
	if !ok {
		return Money{amount: decimal.Zero, currency: getBadCurrency()}, fmt.Errorf("Currency [%s] not supported", curr)
	}
	return Money{
		amount:   decimal.New(value, exp),
		currency: c,
	}, nil

}

// NewFromBigInt returns a new Money from a big.Int, value * 10 ^ exp
func NewFromBigInt(curr string, value *big.Int, exp int32) (Money, error) {

	c, ok := GetCurrency(curr)
	if !ok {
		return Money{amount: decimal.Zero, currency: getBadCurrency()}, fmt.Errorf("Currency [%s] not supported", curr)
	}

	return Money{
		amount:   decimal.NewFromBigInt(value, exp),
		currency: c,
	}, nil

}

// NewFromString returns a new Decimal from a string representation.
//
// Example:
//
//     d, err := NewFromString("USD", "-123.45")
//     d2, err := NewFromString("AUD", ".0001")
//
func NewFromString(curr string, value string) (Money, error) {

	c, ok := GetCurrency(curr)
	if !ok {
		return Money{amount: decimal.Zero, currency: getBadCurrency()}, fmt.Errorf("Currency [%s] not supported", curr)
	}
	d, errr := decimal.NewFromString(value)
	if errr != nil {
		return Money{amount: decimal.Zero, currency: getBadCurrency()}, errr
	}
	return Money{
		amount:   d,
		currency: c,
	}, nil
}

// RequireFromString returns a new Money from a string representation
// or panics if NewFromString would have returned an error.
//
// Example:
//
//     d := RequireFromString("AUD", "-123.45")
//     d2 := RequireFromString("AUD", ".0001")
//
func RequireFromString(curr string, value string) Money {
	mon, err := NewFromString(curr, value)
	if err != nil {
		panic(err)
	}
	return mon
}

// NewFromFloat converts a float64 to Money.
//
// Example:
//
//     NewFromFloat("AUD", 123.45678901234567).String() // output: "$123.4567890123456"
//     NewFromFloat("AUD", .00000000000000001).String() // output: "$0.00000000000000001"
//
// NOTE: some float64 numbers can take up about 300 bytes of memory in decimal representation.
// Consider using NewFromFloatWithExponent if space is more important than precision.
//
// NOTE: this will panic on NaN, +/-inf
func NewFromFloat(curr string, value float64) (Money, error) {
	return NewFromFloatWithExponent(curr, value, math.MinInt32)
}

// NewFromFloatWithExponent converts a float64 to Decimal, with an arbitrary
// number of fractional digits.
//
// Example:
//
//     NewFromFloatWithExponent(123.456, -2).String() // output: "123.46"
//
func NewFromFloatWithExponent(curr string, value float64, exp int32) (Money, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		panic(fmt.Sprintf("Cannot create a Decimal from %v", value))
	}

	c, ok := GetCurrency(curr)
	if !ok {
		return Money{amount: decimal.Zero, currency: getBadCurrency()}, fmt.Errorf("Currency [%s] not supported", curr)
	}

	return Money{
		amount:   decimal.NewFromFloatWithExponent(value, exp),
		currency: c,
	}, nil
}

// UpdateCurrency(newCurr string)
// Allows you to update the currency to the correct code, but only if an UnknownCurrencyCode.
// Otherwise it returns an error (nil if ok)
func (m *Money) UpdateCurrency(newCurr string) error {

	if m.currency.Code != UnknownCurrencyCode {
		return fmt.Errorf("Cannot change currency to [%s]. Already set to [%s]!", newCurr, m.currency.Code)
	}

	c, ok := GetCurrency(newCurr)
	if !ok {
		return fmt.Errorf("Currency [%s] not supported", newCurr)
	}

	m.currency = c

	return nil

}

// Abs returns the absolute value of the decimal.
func (m Money) Abs() Money {

	m.ensureInitialized()

	return Money{
		amount:   m.amount.Abs(),
		currency: m.currency,
	}
}

// Add returns m + m2.
//
// NOTE: This will panic if you try to add Moneys of differing currencies.
// That functionality may come later
func (m Money) Add(m2 Money) Money {

	m.ensureInitialized()
	m2.ensureInitialized()

	if !m.currency.equals(m2.currency) {
		panic(fmt.Sprintf("Cannot add mismatched currencies m1[%s] m2[%s]", m.currency, m2.currency))
	}

	return Money{
		amount:   m.amount.Add(m2.amount),
		currency: m.currency,
	}
}

// Sub returns m - m2.
//
// NOTE: This will panic if you try to subtract Moneys of differing currencies.
// That functionality may come later
func (m Money) Sub(m2 Money) Money {

	m.ensureInitialized()
	m2.ensureInitialized()

	if !m.currency.equals(m2.currency) {
		panic(fmt.Sprintf("Cannot subtract mismatched currencies m1[%s] m2[%s]", m.currency, m2.currency))
	}

	m.ensureInitialized()

	return Money{
		amount:   m.amount.Sub(m2.amount),
		currency: m.currency,
	}
}

// Neg returns -d.
func (m Money) Neg() Money {

	m.ensureInitialized()

	return Money{
		amount:   m.amount.Neg(),
		currency: m.currency,
	}
}

// Mul returns d * d2.
//
// NOTE: This will panic if you try to multiply Moneys of differing currencies.
// That functionality may come later
//
// NOTE: This will also panic if you manage to overflow the amount
func (m Money) Mul(m2 Money) Money {

	m.ensureInitialized()
	m2.ensureInitialized()

	if !m.currency.equals(m2.currency) {
		panic(fmt.Sprintf("Cannot multiply mismatched currencies m1[%s] m2[%s]", m.currency, m2.currency))
	}

	return Money{
		amount:   m.amount.Mul(m2.amount),
		currency: m.currency,
	}
}

// Shift shifts the Money amount in base 10.
// It shifts left when shift is positive and right if shift is negative.
// In simpler terms, the given value for shift is added to the exponent
// of the decimal.
func (m Money) Shift(shift int32) Money {

	m.ensureInitialized()

	return Money{
		amount:   m.amount.Shift(shift),
		currency: m.currency,
	}
}

// DivRound divides and rounds to a given precision
// i.e. to an integer multiple of 10^(-precision)
//   for a positive quotient digit 5 is rounded up, away from 0
//   if the quotient is negative then digit 5 is rounded down, away from 0
// Note that precision<0 is allowed as input.
//
// NOTE: This will panic if you try to divide Moneys of differing currencies.
// That functionality may come later
func (m Money) DivRound(m2 Money, precision int32) Money {

	m.ensureInitialized()
	m2.ensureInitialized()

	if !m.currency.equals(m2.currency) {
		panic(fmt.Sprintf("Cannot divide amounts with mismatched currencies m1[%s] m2[%s]", m.currency, m2.currency))
	}

	return Money{
		amount:   m.amount.DivRound(m2.amount, precision),
		currency: m.currency,
	}
}

// Div returns d / d2. If it doesn't divide exactly, the result will have
// DivisionPrecision digits after the decimal point.
//
// NOTE: This will panic (thrown eventually from DivRound) if you try to
// divide Moneys of differing currencies.
// That functionality may come later
func (m Money) Div(m2 Money) Money {
	return m.DivRound(m2, int32(DivisionPrecision))
}

// QuoRem does divsion with remainder
// d.QuoRem(d2,precision) returns quotient q and remainder r such that
//   d = d2 * q + r, q an integer multiple of 10^(-precision)
//   0 <= r < abs(d2) * 10 ^(-precision) if d>=0
//   0 >= r > -abs(d2) * 10 ^(-precision) if d<0
// Note that precision<0 is allowed as input.
func (m Money) QuoRem(m2 Money, precision int32) (Money, Money) {
	m.ensureInitialized()
	m2.ensureInitialized()

	if !m.currency.equals(m2.currency) {
		panic(fmt.Sprintf("Cannot divide amounts with mismatched currencies m1[%s] m2[%s]", m.currency, m2.currency))
	}

	d1, d2 := m.amount.QuoRem(m2.amount, precision)

	return Money{
			amount:   d1,
			currency: m.currency,
		},
		Money{
			amount:   d2,
			currency: m2.currency,
		}
}

// Mod returns d % d2.
func (m Money) Mod(m2 Money) Money {
	m.ensureInitialized()
	m2.ensureInitialized()

	if !m.currency.equals(m2.currency) {
		panic(fmt.Sprintf("Cannot modulo amounts with mismatched currencies m1[%s] m2[%s]", m.currency, m2.currency))
	}

	return Money{
		amount:   m.amount.Mod(m2.amount),
		currency: m.currency,
	}
}

// Pow returns d to the power d2
func (m Money) Pow(m2 Money) Money {

	m.ensureInitialized()
	m2.ensureInitialized()

	if !m.currency.equals(m2.currency) {
		panic(fmt.Sprintf("Cannot take power of amounts with mismatched currencies m1[%s] m2[%s]", m.currency, m2.currency))
	}

	return Money{
		amount:   m.amount.Pow(m2.amount),
		currency: m.currency,
	}
}

// Cmp compares the numbers represented by d and d2 and returns:
//
//     -1 if d <  d2
//      0 if d == d2
//     +1 if d >  d2
//
// NOTE: This will panic if you try to compare Moneys of differing currencies.
// That functionality may come later
func (m Money) Cmp(m2 Money) int {

	m.ensureInitialized()
	m2.ensureInitialized()

	if !m.currency.equals(m2.currency) {
		panic(fmt.Sprintf("Cannot compare amounts with mismatched currencies m1[%s] m2[%s]", m.currency, m2.currency))
	}

	return m.amount.Cmp(m2.amount)
}

// Equal returns whether the numbers represented by d and d2 are equal.
func (m Money) Equal(m2 Money) bool {
	return m.Cmp(m2) == 0
}

// Equals is deprecated, please use Equal method instead
func (m Money) Equals(m2 Money) bool {
	return m.Equal(m2)
}

// GreaterThan (GT) returns true when d is greater than d2.
func (m Money) GreaterThan(m2 Money) bool {
	return m.Cmp(m2) == 1
}

// GreaterThanOrEqual (GTE) returns true when d is greater than or equal to d2.
func (m Money) GreaterThanOrEqual(m2 Money) bool {
	cmp := m.Cmp(m2)
	return cmp == 1 || cmp == 0
}

// LessThan (LT) returns true when d is less than d2.
func (m Money) LessThan(m2 Money) bool {
	return m.Cmp(m2) == -1
}

// LessThanOrEqual (LTE) returns true when d is less than or equal to d2.
func (m Money) LessThanOrEqual(m2 Money) bool {
	cmp := m.Cmp(m2)
	return cmp == -1 || cmp == 0
}

// Sign returns:
//
//	-1 if d <  0
//	 0 if d == 0
//	+1 if d >  0
//
func (m Money) Sign() int {
	m.ensureInitialized()
	return m.amount.Sign()
}

// Exponent returns the exponent, or scale component of the decimal.
func (m Money) Exponent() int32 {
	m.ensureInitialized()
	return m.amount.Exponent()
}

// Coefficient returns the coefficient of the decimal.  It is scaled by 10^Exponent()
func (m Money) Coefficient() *big.Int {
	// we copy the coefficient so that mutating the result does not mutate the
	// Decimal.
	m.ensureInitialized()

	return m.amount.Coefficient()
}

// IntPart returns the integer component of the decimal.
func (m Money) IntPart() int64 {
	m.ensureInitialized()
	return m.amount.IntPart()
}

// Rat returns a rational number representation of the decimal.
func (m Money) Rat() *big.Rat {
	m.ensureInitialized()
	return m.amount.Rat()
}

// Float64 returns the nearest float64 value for d and a bool indicating
// whether f represents d exactly.
// For more details, see the documentation for big.Rat.Float64
func (m Money) Float64() (f float64, exact bool) {
	m.ensureInitialized()
	return m.amount.Float64()
}

// String returns a simple string representation of the decimal
// with the fixed point.
// Note: It does not pretty-print the amount with currency symbols or
//   thousand separators. Use the CurrencyStringX functions if that's what you need.
// Example:
//
//     d := New("AUD",-12345, -3)
//     println(d.String())
//
// Output:
//
//     -12.345
//
//TODO Fix this.
func (m Money) String() string {
	m.ensureInitialized()
	return m.amount.String()
}

// StringFixed returns a rounded fixed-point string with places digits after
// the decimal point.
//
// Example:
//
// 	   NewFromFloat(0).StringFixed(2) // output: "0.00"
// 	   NewFromFloat(0).StringFixed(0) // output: "0"
// 	   NewFromFloat(5.45).StringFixed(0) // output: "5"
// 	   NewFromFloat(5.45).StringFixed(1) // output: "5.5"
// 	   NewFromFloat(5.45).StringFixed(2) // output: "5.45"
// 	   NewFromFloat(5.45).StringFixed(3) // output: "5.450"
// 	   NewFromFloat(545).StringFixed(-1) // output: "550"
//
//TODO Fix this.
func (m Money) StringFixed(places int32) string {
	m.ensureInitialized()

	return m.amount.StringFixed(places)
}

// StringFixedBank returns a banker rounded fixed-point string with places digits
// after the decimal point.
//
// Example:
//
// 	   NewFromFloat(0).StringFixed(2) // output: "0.00"
// 	   NewFromFloat(0).StringFixed(0) // output: "0"
// 	   NewFromFloat(5.45).StringFixed(0) // output: "5"
// 	   NewFromFloat(5.45).StringFixed(1) // output: "5.4"
// 	   NewFromFloat(5.45).StringFixed(2) // output: "5.45"
// 	   NewFromFloat(5.45).StringFixed(3) // output: "5.450"
// 	   NewFromFloat(545).StringFixed(-1) // output: "550"
//
func (m Money) StringFixedBank(places int32) string {
	m.ensureInitialized()

	return m.amount.StringFixedBank(places)
}

func (m Money) StringFixedCash(interval uint8) string {
	m.ensureInitialized()

	return m.amount.StringFixedCash(interval)
}

// StringFixedBank returns a banker rounded fixed-point string with places digits
// after the decimal point.
//
// Example:
//
// 	   NewFromFloat(0).StringFixed(2) // output: "0.00"
// 	   NewFromFloat(0).StringFixed(0) // output: "0"
// 	   NewFromFloat(5.45).StringFixed(0) // output: "5"
// 	   NewFromFloat(5.45).StringFixed(1) // output: "5.4"
// 	   NewFromFloat(5.45).StringFixed(2) // output: "5.45"
// 	   NewFromFloat(5.45).StringFixed(3) // output: "5.450"
// 	   NewFromFloat(545).StringFixed(-1) // output: "550"
//
func (m Money) FormattedStringBank() string {
	m.ensureInitialized()

	return m.currency.Formatter().FormatCurrency(m.amount)
}

// StringFixedBank returns a banker rounded fixed-point string with places digits
// after the decimal point.
//
// Example:
//
// 	   NewFromFloat(0).StringFixed(2) // output: "0.00"
// 	   NewFromFloat(0).StringFixed(0) // output: "0"
// 	   NewFromFloat(5.45).StringFixed(0) // output: "5"
// 	   NewFromFloat(5.45).StringFixed(1) // output: "5.4"
// 	   NewFromFloat(5.45).StringFixed(2) // output: "5.45"
// 	   NewFromFloat(5.45).StringFixed(3) // output: "5.450"
// 	   NewFromFloat(545).StringFixed(-1) // output: "550"
//
//TODO Fix this.
func (m Money) FormattedStringAccounting() string {
	m.ensureInitialized()

	return m.currency.Formatter().FormatAccounting(m.amount)
}

// StringFixedCash returns a Swedish/Cash rounded fixed-point string. For
// more details see the documentation at function RoundCash.
//TODO Fix this.
func (m Money) FormattedStringFixedCash(interval uint8) string {
	m.ensureInitialized()

	return m.currency.Formatter().FormatCurrency(m.RoundCash(interval).amount)
}

// Round rounds the decimal to places decimal places.
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
// Example:
//
// 	   NewFromFloat(5.45).Round(1).String() // output: "5.5"
// 	   NewFromFloat(545).Round(-1).String() // output: "550"
//
func (m Money) Round(places int32) Money {
	m.ensureInitialized()

	return Money{
		amount:   m.amount.Round(places),
		currency: m.currency,
	}
}

// RoundBank rounds the decimal to places decimal places.
// If the final digit to round is equidistant from the nearest two integers the
// rounded value is taken as the even number
//
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
// Examples:
//
// 	   NewFromFloat(5.45).Round(1).String() // output: "5.4"
// 	   NewFromFloat(545).Round(-1).String() // output: "540"
// 	   NewFromFloat(5.46).Round(1).String() // output: "5.5"
// 	   NewFromFloat(546).Round(-1).String() // output: "550"
// 	   NewFromFloat(5.55).Round(1).String() // output: "5.6"
// 	   NewFromFloat(555).Round(-1).String() // output: "560"
//
func (m Money) RoundBank(places int32) Money {
	m.ensureInitialized()

	return Money{
		amount:   m.amount.RoundBank(places),
		currency: m.currency,
	}

}

// RoundCash aka Cash/Penny/öre rounding rounds decimal to a specific
// interval. The amount payable for a cash transaction is rounded to the nearest
// multiple of the minimum currency unit available. The following intervals are
// available: 5, 10, 15, 25, 50 and 100; any other number throws a panic.
//	    5:   5 cent rounding 3.43 => 3.45
// 	   10:  10 cent rounding 3.45 => 3.50 (5 gets rounded up)
// 	   15:  10 cent rounding 3.45 => 3.40 (5 gets rounded down)
// 	   25:  25 cent rounding 3.41 => 3.50
// 	   50:  50 cent rounding 3.75 => 4.00
// 	  100: 100 cent rounding 3.50 => 4.00
// For more details: https://en.wikipedia.org/wiki/Cash_rounding
func (m Money) RoundCash(interval uint8) Money {
	m.ensureInitialized()

	return Money{
		amount:   m.amount.RoundCash(interval),
		currency: m.currency,
	}
	// TODO: optimize those calculations to reduce the high allocations (~29 allocs).
	// return d.Mul(dVal).Round(0).Div(dVal).Truncate(2)
}

// Floor returns the nearest integer value less than or equal to d.
func (m Money) Floor() Money {
	m.ensureInitialized()

	return Money{
		amount:   m.amount.Floor(),
		currency: m.currency,
	}
}

// Ceil returns the nearest integer value greater than or equal to d.
func (m Money) Ceil() Money {
	m.ensureInitialized()

	return Money{
		amount:   m.amount.Ceil(),
		currency: m.currency,
	}
}

// Truncate truncates off digits from the number, without rounding.
//
// NOTE: precision is the last digit that will not be truncated (must be >= 0).
//
// Example:
//
//     decimal.NewFromString("123.456").Truncate(2).String() // "123.45"
//
func (m Money) Truncate(precision int32) Money {
	m.ensureInitialized()

	return Money{
		amount:   m.amount.Truncate(precision),
		currency: m.currency,
	}
}

// TODO
// UnmarshalJSON implements the json.Unmarshaler interface.
//func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
//	if string(decimalBytes) == "null" {
//		return nil
//	}
//
//	str, err := unquoteIfQuoted(decimalBytes)
//	if err != nil {
//		return fmt.Errorf("Error decoding string '%s': %s", decimalBytes, err)
//	}
//
//	decimal, err := NewFromString(str)
//	*d = decimal
//	if err != nil {
//		return fmt.Errorf("Error decoding string '%s': %s", str, err)
//	}
//	return nil
//}

// TODO
// MarshalJSON implements the json.Marshaler interface.
//func (d Decimal) MarshalJSON() ([]byte, error) {
//	var str string
//	if MarshalJSONWithoutQuotes {
//		str = d.String()
//	} else {
//		str = "\"" + d.String() + "\""
//	}
//	return []byte(str), nil
//}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface. As a string representation
// is already used when encoding to text, this method stores that string as []byte
// NOTE: This is going to break really badlyif we have non ASCII
//    chars in the currency code. Should probably add a length byte at the start
//    but cannot be arsed right now.
func (m *Money) UnmarshalBinary(data []byte) error {

	var err error
	var mo Money

	if ld := len(data); ld < 8 {
		err = fmt.Errorf("Not enough data - only found [%v] bytes", ld)
	} else {
		// Extract the exponent
		curr := string(data[:3])

		// Extract the exponent
		exp := int32(binary.BigEndian.Uint32(data[3:7]))

		// Extract the value
		v := new(big.Int)

		if err = v.GobDecode(data[7:]); err == nil {
			mo, _ = NewFromBigInt(curr, v, exp)
			*m = mo
		} else {
		}
	}

	return err
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
// NOTE: This is going to break really badlyif we have non ASCII
//    chars in the currency code. Should probably add a length byte at the start
//    but cannot be arsed right now.
func (m Money) MarshalBinary() (data []byte, err error) {
	// Write currency first as it's meant to be a fixed size (3 bytes)
	b1 := []byte(m.currency.Code)

	// Write the exponent next since it's a fixed size
	b2 := make([]byte, 4)
	binary.BigEndian.PutUint32(b2, uint32(m.Exponent()))

	b1 = append(b1, b2...)

	// Add the value
	var b3 []byte
	var mCo = m.Coefficient()
	if b3, err = mCo.GobEncode(); err != nil {
		return
	}

	// Return the byte array
	data = append(b1, b3...)

	return
}

// Scan implements the sql.Scanner interface for database deserialization.
func (m *Money) Scan(value interface{}) error {
	// first try to see if the data is stored in database as a Numeric datatype
	switch v := value.(type) {

	case float32:
		*m, _ = NewFromFloat(UnknownCurrencyCode, float64(v))
		return nil

	case float64:
		// numeric in sqlite3 sends us float64
		*m, _ = NewFromFloat(UnknownCurrencyCode, v)
		return nil

	case int64:
		// at least in sqlite3 when the value is 0 in db, the data is sent
		// to us as an int64 instead of a float64 ...
		*m, _ = New(UnknownCurrencyCode, v, 0)
		return nil

	default:
		// default is trying to interpret value stored as string
		str, err := unquoteIfQuoted(v)
		if err != nil {
			return err
		}
		*m, err = NewFromString(UnknownCurrencyCode, str)
		return err
	}
}

// Value implements the driver.Valuer interface for database serialization.
func (m Money) Value() (driver.Value, error) {
	return m.String(), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization.
func (d *Money) UnmarshalText(text []byte) error {
	str := string(text)

	dec, err := NewFromString(UnknownCurrencyCode, str)
	*d = dec
	if err != nil {
		return fmt.Errorf("Error decoding string '%s': %s", str, err)
	}

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization.
func (d Money) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (m Money) GobEncode() ([]byte, error) {
	return m.MarshalBinary()
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (m *Money) GobDecode(data []byte) error {
	return m.UnmarshalBinary(data)
}

// Checks to see if we actually have a proper Money object.
// If not, create a valid Zero so we can at least not crash things too badly.
func (m *Money) ensureInitialized() {
	if m.currency == nil {
		m.currency = getUnknownCurrency()
	}
}

// Min returns the smallest Decimal that was passed in the arguments.
//
// To call this function with an array, you must do:
//
//     Min(arr[0], arr[1:]...)
//
// This makes it harder to accidentally call Min with 0 arguments.
func Min(first Money, rest ...Money) Money {
	ans := first
	for _, item := range rest {
		if item.Cmp(ans) < 0 {
			ans = item
		}
	}
	return ans
}

// Max returns the largest Decimal that was passed in the arguments.
//
// To call this function with an array, you must do:
//
//     Max(arr[0], arr[1:]...)
//
// This makes it harder to accidentally call Max with 0 arguments.
func Max(first Money, rest ...Money) Money {
	ans := first
	for _, item := range rest {
		if item.Cmp(ans) > 0 {
			ans = item
		}
	}
	return ans
}

// Sum returns the combined total of the provided first and rest Decimals
func Sum(first Money, rest ...Money) Money {
	total := first
	for _, item := range rest {
		total = total.Add(item)
	}

	return total
}

// Avg returns the average value of the provided first and rest Decimals
func Avg(first Money, rest ...Money) Money {
	count, _ := New(first.currency.Code, int64(len(rest)+1), 0)
	sum := Sum(first, rest...)
	return sum.Div(count)
}

func min(x, y int32) int32 {
	if x >= y {
		return y
	}
	return x
}

func unquoteIfQuoted(value interface{}) (string, error) {
	var bytes []byte

	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return "", fmt.Errorf("Could not convert value '%+v' to byte array of type '%T'",
			value, value)
	}

	// If the amount is quoted, strip the quotes
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes), nil
}

// NullMoney represents a nullable decimal with compatibility for
// scanning null values from the database.
type NullMoney struct {
	Money Money
	Valid bool
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *NullMoney) Scan(value interface{}) error {
	if value == nil {
		d.Valid = false
		return nil
	}
	d.Valid = true
	return d.Money.Scan(value)
}

// Value implements the driver.Valuer interface for database serialization.
func (d NullMoney) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Money.Value()
}

//// UnmarshalJSON implements the json.Unmarshaler interface.
//func (d *NullMoney) UnmarshalJSON(decimalBytes []byte) error {
//	if string(decimalBytes) == "null" {
//		d.Valid = false
//		return nil
//	}
//	d.Valid = true
//	return d.Money.UnmarshalJSON(decimalBytes)
//}
//
//// MarshalJSON implements the json.Marshaler interface.
//func (d NullMoney) MarshalJSON() ([]byte, error) {
//	if !d.Valid {
//		return []byte("null"), nil
//	}
//	return d.Money.MarshalJSON()
//}
