package main

import (
	"fmt"

	BigNumber "github.com/go-composites/bignumber/src"
	Error "github.com/go-composites/error/src"
	Result "github.com/go-composites/result/src"
)

func report(label string, result Result.Interface) {
	if result.HasError() {
		fmt.Printf("%s -> error: %s\n", label, result.Error().Message())
		return
	}
	fmt.Printf("%s -> %s\n", label, result.Payload().(BigNumber.Interface).ToGoString())
}

func main() {
	six := BigNumber.FromInt64(6)
	two := BigNumber.FromInt64(2)
	zero := BigNumber.FromInt64(0)

	report("6 + 2", six.Add(two))
	report("6 - 2", six.Sub(two))
	report("6 * 2", six.Mul(two))
	report("6 / 2", six.Div(two))

	// The canonical Result use-case: division by zero is a value, not a panic.
	divByZero := six.Div(zero)
	fmt.Println("6 / 0 has error:", divByZero.HasError())
	report("6 / 0", divByZero)

	// Errors are first-class values.
	var _ Error.Interface = divByZero.Error()

	// Arbitrary precision: a 40-digit number squared overflows int64 by far,
	// yet BigNumber computes it exactly.
	huge := BigNumber.FromString("9999999999999999999999999999999999999999")
	if !huge.HasError() {
		n := huge.Payload().(BigNumber.Interface)
		report("40-digit squared", n.Mul(n))
	}

	fmt.Println("6 == 2 :", six.Equal(two))
	fmt.Println("6 < 2  :", six.LessThan(two))
	fmt.Println("6 > 2  :", six.GreaterThan(two))
	fmt.Println(six.Inspect())
}
