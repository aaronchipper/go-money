// package money - Currency definitions
// Shamelessly stolen from github.com/Rhymond/go-money/formatter.go
// Changes from original:
// 		* Changed to support DecPoint type, rather than int64

package money

import (
	"github.com/shopspring/decimal"
	"strings"
)

// Formatter stores Money formatting information
type Formatter struct {
	Fraction int
	DecPoint string
	Thousand string
	Grapheme string
	Template string
}

// NewFormatter creates new Formatter instance
func NewFormatter(fraction int, decpoint, thousand, grapheme, template string) *Formatter {
	return &Formatter{
		Fraction: fraction,
		DecPoint: decpoint,
		Thousand: thousand,
		Grapheme: grapheme,
		Template: template,
	}
}

// FormatWithOptions returns string of formatted integer using given currency template
//		amount: The amount to be displayed
//		noThousands : Boolean - If true, don't bother with the thousands separator.
//		noCurrencyGrapheme: Boolean - If true, we'll hide the $ (or whatever) symbol
//		negsInBrackets: Boolean - If true, we'll display negative numbers as "($1,000.00)" as opposed to "-$100.00"
func (f *Formatter) formatWithOptions(amount decimal.Decimal, noThousands, noCurrencyGrapheme, negsInBrackets bool) string {

	// Work with absolute amount value
	// Then print as a Bank Rounded number to the display amount based on the currency
	// Then split into int and fractional parts for correct formatting
	numBits := strings.Split(amount.Abs().StringFixedBank(int32(f.Fraction)), ".")

	fractionalPart := ""
	intPart := numBits[0]
	if len(numBits) > 1 {
		fractionalPart = numBits[1]
	}

	if !noThousands {
		if f.Thousand != "" {
			for i := len(intPart) - 3; i > 0; i -= 3 {
				intPart = intPart[:i] + f.Thousand + intPart[i:]
			}
		}
	}

	// Time to combine. Hijacking intPart because renaming is pointless now.
	if len(fractionalPart) > 0 {
		intPart += f.DecPoint + fractionalPart
	}

	// Got the number looking nice, now for the trimmings
	intPart = strings.Replace(f.Template, "1", intPart, 1)

	// Add (or hide) the currency grapheme
	if noCurrencyGrapheme {
		intPart = strings.TrimSpace(strings.Replace(intPart, "$", "", 1))
	} else {
		intPart = strings.Replace(intPart, "$", f.Grapheme, 1)
	}

	// Add minus sign for negative amount
	if amount.Sign() < 0 {
		if negsInBrackets {
			intPart = "(" + intPart + ")"
		} else {
			intPart = "-" + intPart
		}
	}

	return intPart
}

// Format returns string of formatted integer using given currency template
//		amount: The amount to be displayed
func (f *Formatter) FormatAccounting(amount decimal.Decimal) string {
	return f.formatWithOptions(amount, true, true, true)
}

// Format returns string of formatted integer using given currency template
//		amount: The amount to be displayed
func (f *Formatter) FormatCurrency(amount decimal.Decimal) string {
	return f.formatWithOptions(amount, false, false, false)
}
