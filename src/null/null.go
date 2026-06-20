package NullBigNumber

import (
	BigNumber "github.com/go-composites/bignumber/src"
	MethodNotImplementedError "github.com/go-composites/error/src/method_not_implemented"
	Result "github.com/go-composites/result/src"
)

/*
NullBigNumber is the Null-Object variant of BigNumber.

It satisfies BigNumber.Interface so callers never have to test for a bare nil:
its value is zero, its arithmetic yields a Result carrying a
"method not implemented" Error, its comparisons are false (except Equal against
another null), and IsNull() returns true.
*/
type Interface interface {
	BigNumber.Interface
}

type data struct{}

/*
New returns a NullBigNumber.
*/
func New() Interface {
	return &data{}
}

func (d data) ToGoString() string {
	return ``
}

func (d data) ToInt64() int64 {
	return 0
}

func (d data) IsNull() bool {
	return true
}

func notImplemented(methodName string) Result.Interface {
	return Result.New(
		Result.WithError(
			MethodNotImplementedError.New(methodName),
		),
	)
}

func (d data) Add(BigNumber.Interface) Result.Interface {
	return notImplemented(`Add`)
}

func (d data) Sub(BigNumber.Interface) Result.Interface {
	return notImplemented(`Sub`)
}

func (d data) Mul(BigNumber.Interface) Result.Interface {
	return notImplemented(`Mul`)
}

func (d data) Div(BigNumber.Interface) Result.Interface {
	return notImplemented(`Div`)
}

func (d data) Mod(BigNumber.Interface) Result.Interface {
	return notImplemented(`Mod`)
}

func (d data) Abs() Result.Interface {
	return notImplemented(`Abs`)
}

func (d data) Neg() Result.Interface {
	return notImplemented(`Neg`)
}

func (d data) Equal(other BigNumber.Interface) bool {
	return other.IsNull()
}

func (d data) LessThan(BigNumber.Interface) bool {
	return false
}

func (d data) GreaterThan(BigNumber.Interface) bool {
	return false
}

func (d data) Inspect() BigNumber.String {
	return `<NullBigNumber>`
}
