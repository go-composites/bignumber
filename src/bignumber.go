package BigNumber

import (
	"fmt"
	"math/big"

	Error "github.com/go-composites/error/src"
	MethodNotImplementedError "github.com/go-composites/error/src/method_not_implemented"
	Result "github.com/go-composites/result/src"
)

/*
BigNumber is an arbitrary-precision integer composite over a math/big.Int.

It mirrors Ruby's unbounded Integer: there is no overflow, only memory. Its
fallible operations (notably Div and Mod) return a Result.Interface so that
failures — such as a division by zero — are values rather than panics, and
they never return a bare nil.
*/
type Interface interface {
	ToGoString() string
	ToInt64() int64
	IsNull() bool
	Add(Interface) Result.Interface
	Sub(Interface) Result.Interface
	Mul(Interface) Result.Interface
	Div(Interface) Result.Interface
	Mod(Interface) Result.Interface
	Abs() Result.Interface
	Neg() Result.Interface
	Equal(Interface) bool
	LessThan(Interface) bool
	GreaterThan(Interface) bool
	Inspect() String
}

// String is the lightweight inspection representation of a BigNumber.
type String = string

type data struct {
	value *big.Int
}

/*
FromInt64 is the BigNumber constructor from a Go int64.

	n := BigNumber.FromInt64(42) // 42
*/
func FromInt64(v int64) Interface {
	return &data{value: big.NewInt(v)}
}

/*
FromString parses a base-10 decimal string into a BigNumber.

It returns a Result whose payload is the parsed BigNumber. When the input is
not a valid base-10 integer the Result carries an Error instead of a payload —
the parse never panics and never returns nil. This is how arbitrary-precision
values that overflow int64 are constructed.

	r := BigNumber.FromString("123456789012345678901234567890")
	if !r.HasError() {
	    n := r.Payload().(BigNumber.Interface)
	}
*/
func FromString(s string) Result.Interface {
	value, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return Result.New(
			Result.WithError(
				Error.New("invalid integer: " + s),
			),
		)
	}
	return Result.New(
		Result.WithPayload(
			&data{value: value},
		),
	)
}

/*
Null returns the Null-Object variant of BigNumber.

It is defined in src/null; this thin re-export keeps a Null next to the
concrete constructors. The returned value satisfies Interface and reports
IsNull() == true.
*/
func Null() Interface {
	return newNull()
}

/*
ToGoString returns the base-10 decimal representation of the value.
*/
func (d data) ToGoString() string {
	return d.value.String()
}

/*
ToInt64 returns the value as a Go int64.

When the value does not fit in an int64 the result is undefined in the same way
as math/big.Int.Int64 (the low 64 bits), so callers handling arbitrary
precision should prefer ToGoString.
*/
func (d data) ToInt64() int64 {
	return d.value.Int64()
}

/*
IsNull reports whether the BigNumber is the Null-Object variant.

A concrete BigNumber is never null.
*/
func (d data) IsNull() bool {
	return false
}

/*
Add returns a Result whose payload is the sum of the receiver and other.

A fresh big.Int backs the payload; the operands are never mutated.
*/
func (d data) Add(other Interface) Result.Interface {
	return payload(
		new(big.Int).Add(d.value, fromInterface(other)),
	)
}

/*
Sub returns a Result whose payload is the difference of the receiver and other.
*/
func (d data) Sub(other Interface) Result.Interface {
	return payload(
		new(big.Int).Sub(d.value, fromInterface(other)),
	)
}

/*
Mul returns a Result whose payload is the product of the receiver and other.
*/
func (d data) Mul(other Interface) Result.Interface {
	return payload(
		new(big.Int).Mul(d.value, fromInterface(other)),
	)
}

/*
Div returns a Result whose payload is the quotient of the receiver and other.

When other is zero the Result carries an Error ("division by zero") instead of
a payload — the division never panics and never returns nil.
*/
func (d data) Div(other Interface) Result.Interface {
	rhs := fromInterface(other)
	if rhs.Sign() == 0 {
		return Result.New(
			Result.WithError(
				Error.New("division by zero"),
			),
		)
	}
	return payload(
		new(big.Int).Quo(d.value, rhs),
	)
}

/*
Mod returns a Result whose payload is the remainder of the receiver divided by
other.

When other is zero the Result carries an Error ("modulo by zero") instead of a
payload — the operation never panics and never returns nil.
*/
func (d data) Mod(other Interface) Result.Interface {
	rhs := fromInterface(other)
	if rhs.Sign() == 0 {
		return Result.New(
			Result.WithError(
				Error.New("modulo by zero"),
			),
		)
	}
	return payload(
		new(big.Int).Rem(d.value, rhs),
	)
}

/*
Abs returns a Result whose payload is the absolute value of the receiver.
*/
func (d data) Abs() Result.Interface {
	return payload(
		new(big.Int).Abs(d.value),
	)
}

/*
Neg returns a Result whose payload is the negation of the receiver.
*/
func (d data) Neg() Result.Interface {
	return payload(
		new(big.Int).Neg(d.value),
	)
}

/*
Equal reports whether the receiver and other hold the same integer value.
*/
func (d data) Equal(other Interface) bool {
	return d.value.Cmp(fromInterface(other)) == 0
}

/*
LessThan reports whether the receiver is strictly less than other.
*/
func (d data) LessThan(other Interface) bool {
	return d.value.Cmp(fromInterface(other)) < 0
}

/*
GreaterThan reports whether the receiver is strictly greater than other.
*/
func (d data) GreaterThan(other Interface) bool {
	return d.value.Cmp(fromInterface(other)) > 0
}

/*
Inspect returns a one-line representation of the BigNumber with its address and
value — mirroring the style of the other composites.
*/
func (d data) Inspect() String {
	return fmt.Sprintf(
		"<BigNumber:%p value=%s>",
		&d, d.value.String(),
	)
}

// nullData is the Null-Object variant returned by Null(). The importable
// NullBigNumber package in src/null mirrors it; this copy keeps a Null next to
// the concrete constructors without creating an import cycle.
type nullData struct{}

func newNull() Interface {
	return &nullData{}
}

func nullNotImplemented(methodName string) Result.Interface {
	return Result.New(
		Result.WithError(
			MethodNotImplementedError.New(methodName),
		),
	)
}

func (nullData) ToGoString() string             { return `` }
func (nullData) ToInt64() int64                 { return 0 }
func (nullData) IsNull() bool                   { return true }
func (nullData) Add(Interface) Result.Interface { return nullNotImplemented(`Add`) }
func (nullData) Sub(Interface) Result.Interface { return nullNotImplemented(`Sub`) }
func (nullData) Mul(Interface) Result.Interface { return nullNotImplemented(`Mul`) }
func (nullData) Div(Interface) Result.Interface { return nullNotImplemented(`Div`) }
func (nullData) Mod(Interface) Result.Interface { return nullNotImplemented(`Mod`) }
func (nullData) Abs() Result.Interface          { return nullNotImplemented(`Abs`) }
func (nullData) Neg() Result.Interface          { return nullNotImplemented(`Neg`) }
func (nullData) Equal(other Interface) bool     { return other.IsNull() }
func (nullData) LessThan(Interface) bool        { return false }
func (nullData) GreaterThan(Interface) bool     { return false }
func (nullData) Inspect() String                { return `<NullBigNumber>` }

// payload wraps a fresh big.Int in a success Result.
func payload(value *big.Int) Result.Interface {
	return Result.New(
		Result.WithPayload(
			&data{value: value},
		),
	)
}

// fromInterface extracts a *big.Int from any BigNumber.Interface, parsing its
// decimal string when the concrete type is unknown (e.g. the Null-Object).
// The returned big.Int is always a fresh copy, so operands are never shared.
func fromInterface(other Interface) *big.Int {
	if d, ok := other.(*data); ok {
		return new(big.Int).Set(d.value)
	}
	value, ok := new(big.Int).SetString(other.ToGoString(), 10)
	if !ok {
		return new(big.Int)
	}
	return value
}
