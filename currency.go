// package money - Currency definitions
// Shamelessly stolen from github.com/Rhymond/go-money/currency.go
// Changes from original:
// 		* Added support for "types" of currency, not just fiat.
// 		* Modified some internal names as well to reduce conflicts
package money

import (
	"strings"
)

// CurrType assigns a currency type to the given currency. This is an extension
// to the original Currency code from the original Rhymond/go-money/currency.go
type CurrType int

const (
	FIAT	CurrType	= iota	//	FIAT   	(regular currency)
	CRYPTO						//	CRYPTO	(the new kids on the block, like BTC, XEM, etc)
	LOYALTY						//	LOYALTY	(loyalty program points)
	REWARD						//	REWARD	(rewards points. Not sure if this is the same as loyalty, but something in my gut tells me I'll want to differentiate some day, so... meh)
	GAME						//	GAME	(Game credits)
	POINTS						//	POINTS	(Generic value store)
	
	UNKNOWN				= 9999  //  UNKNOWN    (Testing currencies. Should never see in production)
)

// UnknownCurrencyCode is used when creating Money objects where we don't yet know the currency.
const (
	UnknownCurrencyCode		= "???"
)

// Currency represents money currency information required for formatting
type Currency struct {
	Type	 CurrType
	Code     string
	Fraction int
	Grapheme string
	Template string
	DecPoint  string
	Thousand string
}

// currencies represents a collection of currency
// If adding to this list, you should only use 3 ASCII chars as the code. 
// If this changes, we'll need to fix the (Un)MarshallBinary functions as they'll break badly. 
var currencies = map[string]*Currency{
	// Fiat Currencies
	"AED": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "AED", Fraction: 2, Grapheme: ".\u062f.\u0625", Template: "1 $"},
	"AFN": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "AFN", Fraction: 2, Grapheme: "\u060b", Template: "1 $"},
	"ALL": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "ALL", Fraction: 2, Grapheme: "L", Template: "$1"},
	"AMD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "AMD", Fraction: 2, Grapheme: "\u0564\u0580.", Template: "1 $"},
	"ANG": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "ANG", Fraction: 2, Grapheme: "\u0192", Template: "$1"},
	"ARS": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "ARS", Fraction: 2, Grapheme: "$", Template: "$1"},
	"AUD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "AUD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"AWG": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "AWG", Fraction: 2, Grapheme: "\u0192", Template: "$1"},
	"AZN": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "AZN", Fraction: 2, Grapheme: "\u20bc", Template: "$1"},
	"BAM": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BAM", Fraction: 2, Grapheme: "KM", Template: "$1"},
	"BBD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BBD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"BGN": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BGN", Fraction: 2, Grapheme: "\u043b\u0432", Template: "$1"},
	"BHD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BHD", Fraction: 3, Grapheme: ".\u062f.\u0628", Template: "1 $"},
	"BMD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BMD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"BND": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BND", Fraction: 2, Grapheme: "$", Template: "$1"},
	"BOB": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BOB", Fraction: 2, Grapheme: "Bs.", Template: "$1"},
	"BRL": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BRL", Fraction: 2, Grapheme: "R$", Template: "$1"},
	"BSD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BSD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"BWP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BWP", Fraction: 2, Grapheme: "P", Template: "$1"},
	"BYN": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BYN", Fraction: 2, Grapheme: "p.", Template: "1 $"},
	"BYR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BYR", Fraction: 0, Grapheme: "p.", Template: "1 $"},
	"BZD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "BZD", Fraction: 2, Grapheme: "BZ$", Template: "$1"},
	"CAD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "CAD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"CLP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "CLP", Fraction: 0, Grapheme: "$", Template: "$1"},
	"CNY": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "CNY", Fraction: 2, Grapheme: "\u5143", Template: "1 $"},
	"COP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "COP", Fraction: 0, Grapheme: "$", Template: "$1"},
	"CRC": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "CRC", Fraction: 2, Grapheme: "\u20a1", Template: "$1"},
	"CUP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "CUP", Fraction: 2, Grapheme: "$MN", Template: "$1"},
	"CZK": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "CZK", Fraction: 2, Grapheme: "K\u010d", Template: "1 $"},
	"DKK": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "DKK", Fraction: 2, Grapheme: "kr", Template: "1 $"},
	"DOP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "DOP", Fraction: 2, Grapheme: "RD$", Template: "$1"},
	"DZD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "DZD", Fraction: 2, Grapheme: ".\u062f.\u062c", Template: "1 $"},
	"EEK": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "EEK", Fraction: 2, Grapheme: "kr", Template: "$1"},
	"EGP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "EGP", Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"EUR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "EUR", Fraction: 2, Grapheme: "\u20ac", Template: "$1"},
	"FJD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "FJD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"FKP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "FKP", Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"GBP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "GBP", Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"GGP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "GGP", Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"GHC": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "GHC", Fraction: 2, Grapheme: "\u00a2", Template: "$1"},
	"GIP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "GIP", Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"GTQ": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "GTQ", Fraction: 2, Grapheme: "Q", Template: "$1"},
	"GYD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "GYD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"HKD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "HKD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"HNL": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "HNL", Fraction: 2, Grapheme: "L", Template: "$1"},
	"HRK": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "HRK", Fraction: 2, Grapheme: "kn", Template: "$1"},
	"HUF": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "HUF", Fraction: 0, Grapheme: "Ft", Template: "$1"},
	"IDR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "IDR", Fraction: 2, Grapheme: "Rp", Template: "$1"},
	"ILS": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "ILS", Fraction: 2, Grapheme: "\u20aa", Template: "$1"},
	"IMP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "IMP", Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"INR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "INR", Fraction: 2, Grapheme: "\u20b9", Template: "$1"},
	"IQD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "IQD", Fraction: 3, Grapheme: ".\u062f.\u0639", Template: "1 $"},
	"IRR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "IRR", Fraction: 2, Grapheme: "\ufdfc", Template: "1 $"},
	"ISK": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "ISK", Fraction: 2, Grapheme: "kr", Template: "$1"},
	"JEP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "JEP", Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"JMD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "JMD", Fraction: 2, Grapheme: "J$", Template: "$1"},
	"JOD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "JOD", Fraction: 3, Grapheme: ".\u062f.\u0625", Template: "1 $"},
	"JPY": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "JPY", Fraction: 0, Grapheme: "\u00a5", Template: "$1"},
	"KES": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "KES", Fraction: 2, Grapheme: "KSh", Template: "$1"},
	"KGS": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "KGS", Fraction: 2, Grapheme: "\u0441\u043e\u043c", Template: "$1"},
	"KHR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "KHR", Fraction: 2, Grapheme: "\u17db", Template: "$1"},
	"KPW": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "KPW", Fraction: 0, Grapheme: "\u20a9", Template: "$1"},
	"KRW": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "KRW", Fraction: 0, Grapheme: "\u20a9", Template: "$1"},
	"KWD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "KWD", Fraction: 3, Grapheme: ".\u062f.\u0643", Template: "1 $"},
	"KYD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "KYD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"KZT": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "KZT", Fraction: 2, Grapheme: "\u20b8", Template: "$1"},
	"LAK": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "LAK", Fraction: 2, Grapheme: "\u20ad", Template: "$1"},
	"LBP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "LBP", Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"LKR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "LKR", Fraction: 2, Grapheme: "\u20a8", Template: "$1"},
	"LRD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "LRD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"LTL": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "LTL", Fraction: 2, Grapheme: "Lt", Template: "$1"},
	"LVL": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "LVL", Fraction: 2, Grapheme: "Ls", Template: "1 $"},
	"LYD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "LYD", Fraction: 3, Grapheme: ".\u062f.\u0644", Template: "1 $"},
	"MAD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "MAD", Fraction: 2, Grapheme: ".\u062f.\u0645", Template: "1 $"},
	"MKD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "MKD", Fraction: 2, Grapheme: "\u0434\u0435\u043d", Template: "$1"},
	"MNT": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "MNT", Fraction: 2, Grapheme: "\u20ae", Template: "$1"},
	"MUR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "MUR", Fraction: 2, Grapheme: "\u20a8", Template: "$1"},
	"MXN": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "MXN", Fraction: 2, Grapheme: "$", Template: "$1"},
	"MWK": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "MWK", Fraction: 2, Grapheme: "MK", Template: "$1"},
	"MYR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "MYR", Fraction: 2, Grapheme: "RM", Template: "$1"},
	"MZN": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "MZN", Fraction: 2, Grapheme: "MT", Template: "$1"},
	"NAD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "NAD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"NGN": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "NGN", Fraction: 2, Grapheme: "\u20a6", Template: "$1"},
	"NIO": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "NIO", Fraction: 2, Grapheme: "C$", Template: "$1"},
	"NOK": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "NOK", Fraction: 2, Grapheme: "kr", Template: "1 $"},
	"NPR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "NPR", Fraction: 2, Grapheme: "\u20a8", Template: "$1"},
	"NZD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "NZD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"OMR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "OMR", Fraction: 3, Grapheme: "\ufdfc", Template: "1 $"},
	"PAB": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "PAB", Fraction: 2, Grapheme: "B/.", Template: "$1"},
	"PEN": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "PEN", Fraction: 2, Grapheme: "S/", Template: "$1"},
	"PHP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "PHP", Fraction: 2, Grapheme: "\u20b1", Template: "$1"},
	"PKR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "PKR", Fraction: 2, Grapheme: "\u20a8", Template: "$1"},
	"PLN": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "PLN", Fraction: 2, Grapheme: "z\u0142", Template: "1 $"},
	"PYG": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "PYG", Fraction: 0, Grapheme: "Gs", Template: "1$"},
	"QAR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "QAR", Fraction: 2, Grapheme: "\ufdfc", Template: "1 $"},
	"RON": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "RON", Fraction: 2, Grapheme: "lei", Template: "$1"},
	"RSD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "RSD", Fraction: 2, Grapheme: "\u0414\u0438\u043d.", Template: "$1"},
	"RUB": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "RUB", Fraction: 2, Grapheme: "\u20bd", Template: "1 $"},
	"RUR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "RUR", Fraction: 2, Grapheme: "\u20bd", Template: "1 $"},
	"SAR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "SAR", Fraction: 2, Grapheme: "\ufdfc", Template: "1 $"},
	"SBD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "SBD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"SCR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "SCR", Fraction: 2, Grapheme: "\u20a8", Template: "$1"},
	"SEK": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "SEK", Fraction: 2, Grapheme: "kr", Template: "1 $"},
	"SGD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "SGD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"SHP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "SHP", Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"SOS": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "SOS", Fraction: 2, Grapheme: "S", Template: "$1"},
	"SRD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "SRD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"SVC": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "SVC", Fraction: 2, Grapheme: "$", Template: "$1"},
	"SYP": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "SYP", Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"THB": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "THB", Fraction: 2, Grapheme: "\u0e3f", Template: "$1"},
	"TND": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "TND", Fraction: 3, Grapheme: ".\u062f.\u062a", Template: "1 $"},
	"TRL": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "TRL", Fraction: 2, Grapheme: "\u20a4", Template: "$1"},
	"TRY": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "TRY", Fraction: 2, Grapheme: "\u20ba", Template: "$1"},
	"TTD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "TTD", Fraction: 2, Grapheme: "TT$", Template: "$1"},
	"TWD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "TWD", Fraction: 0, Grapheme: "NT$", Template: "$1"},
	"TZS": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "TZS", Fraction: 0, Grapheme: "TSh", Template: "$1"},
	"UAH": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "UAH", Fraction: 2, Grapheme: "\u20b4", Template: "$1"},
	"UGX": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "UGX", Fraction: 0, Grapheme: "USh", Template: "$1"},
	"USD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "USD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"UYU": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "UYU", Fraction: 0, Grapheme: "$U", Template: "$1"},
	"UZS": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "UZS", Fraction: 2, Grapheme: "so\u2019m", Template: "$1"},
	"VEF": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "VEF", Fraction: 2, Grapheme: "Bs", Template: "$1"},
	"VND": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "VND", Fraction: 0, Grapheme: "\u20ab", Template: "1 $"},
	"XCD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "XCD", Fraction: 2, Grapheme: "$", Template: "$1"},
	"YER": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "YER", Fraction: 2, Grapheme: "\ufdfc", Template: "1 $"},
	"ZAR": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "ZAR", Fraction: 2, Grapheme: "R", Template: "$1"},
	"ZMW": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "ZMW", Fraction: 2, Grapheme: "ZK", Template: "$1"},
	"ZWD": {Type: FIAT, DecPoint: ".", Thousand: ",", Code: "ZWD", Fraction: 2, Grapheme: "Z$", Template: "$1"},

	// Cryptocurrencies
	// Bitcoin has 2 accepted codes as of now. ISO 4217 standard is moving to XBT at some point
	"BTC": {Type: CRYPTO, DecPoint: ".", Thousand: ",", Code: "BTC", Fraction: 8, Grapheme: "\u20bf", Template: "$1"},
	"XBT": {Type: CRYPTO, DecPoint: ".", Thousand: ",", Code: "XBT", Fraction: 8, Grapheme: "\u20bf", Template: "$1"},
	
	// Unknown currency.
	// Only to be used in Test code.
	// Seriously, don't be a dick with them.
	"???": {Type: UNKNOWN, DecPoint: ".", Thousand: ",", Code: "???", Fraction: 2, Grapheme: "$", Template: "$1"},

}

// AddCurrency lets you insert or update currency in currencies list
func AddCurrency(Type CurrType, Code, Grapheme, Template, DecPoint, Thousand string, Fraction int) *Currency {
	currencies[Code] = &Currency{
		Type:		Type,
		Code:     	Code,
		Grapheme: 	Grapheme,
		Template: 	Template,
		DecPoint:  	DecPoint,
		Thousand: 	Thousand,
		Fraction: 	Fraction,
	}

	return currencies[Code]
}

func newCurrency(code string) *Currency {
	return &Currency{Code: strings.ToUpper(code)}
}

// GetCurrency returns the currency given the code.
func GetCurrency(code string) (*Currency, bool) {
	c, err := currencies[code]
	return c, err 
}

// Formatter returns currency formatter representing
// used currency structure
func (c *Currency) Formatter() *Formatter {
	return &Formatter{
		Fraction: c.Fraction,
		DecPoint:  c.DecPoint,
		Thousand: c.Thousand,
		Grapheme: c.Grapheme,
		Template: c.Template,
	}
}

// getDefault represent default currency if currency is not found in currencies list.
// Grapheme and Code fields will be changed by currency code
func (c *Currency) getDefault() *Currency {
	return &Currency{Type: FIAT, DecPoint: ".", Thousand: ",", Code: c.Code, Fraction: 2, Grapheme: c.Code, Template: "1$"}
}

// getDefault represent default currency if currency is not found in currencies list.
// Grapheme and Code fields will be changed by currency code
func getUnknownCurrency() *Currency {
	return &Currency{Type: FIAT, DecPoint: ".", Thousand: ",", Code: UnknownCurrencyCode, Fraction: 2, Grapheme: "$", Template: "1$"}
}

// get extended currency using currencies list
func (c *Currency) get() *Currency {
	if curr, ok := currencies[c.Code]; ok {
		return curr
	}

	return c.getDefault()
}

func (c *Currency) equals(oc *Currency) bool {
	return c.Code == oc.Code
}

func (c *Currency) String() string {
	return c.Code
}
