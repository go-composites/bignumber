<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/go-composites.png" alt="go-composites/bignumber" width="720"></p>

# bignumber

[![ci](https://github.com/go-composites/bignumber/actions/workflows/ci.yml/badge.svg)](https://github.com/go-composites/bignumber/actions/workflows/ci.yml)

An **arbitrary-precision integer** composite for **Composition-Oriented
Programming**. A `BigNumber` wraps Go's `math/big.Int` and mirrors Ruby's
unbounded `Integer`: there is no overflow, only memory. Its arithmetic is
exposed as **fallible operations that return a `Result`** — so failures (the
canonical example being a division by zero) are *values*, never panics and
never `nil`.

```golang
quotient := numerator.Div(denominator)
if quotient.HasError() {
    fmt.Println(quotient.Error().Message()) // "division by zero"
} else {
    fmt.Println(quotient.Payload().(BigNumber.Interface).ToGoString())
}
```

`BigNumber` follows the org's Null-Object / never-nil invariant (enforced by the
`nonnil` CI analyzer): the `NullBigNumber` variant in `src/null` satisfies the
same `Interface` and reports `IsNull() == true`.

## Install

```bash
export GOPRIVATE=github.com/go-composites GOPROXY=direct GOSUMDB=off
go get github.com/go-composites/bignumber@main
```

## Usage

> [!NOTE] main.go

```golang
package main

import (
    "fmt"

    BigNumber "github.com/go-composites/bignumber/src"
)

func main() {
    six := BigNumber.FromInt64(6)
    two := BigNumber.FromInt64(2)
    zero := BigNumber.FromInt64(0)

    // Arithmetic returns a Result.
    sum := six.Add(two)
    fmt.Println(sum.Payload().(BigNumber.Interface).ToGoString()) // 8

    // Division by zero is a value, not a panic.
    div := six.Div(zero)
    fmt.Println("has error:", div.HasError())      // true
    fmt.Println(div.Error().Message())             // division by zero

    // Arbitrary precision: a 40-digit number squared overflows int64,
    // yet BigNumber computes it exactly.
    huge := BigNumber.FromString("9999999999999999999999999999999999999999")
    n := huge.Payload().(BigNumber.Interface)
    fmt.Println(n.Mul(n).Payload().(BigNumber.Interface).ToGoString())

    fmt.Println(six.GreaterThan(two)) // true
    fmt.Println(six.Inspect())        // <BigNumber:0x... value=6>
}
```

```bash
$ go run .
```

## API

Constructors

- `FromInt64(v int64) Interface` — build from a Go int64.
- `FromString(s string) Result.Interface` — parse a base-10 decimal string; a
  `Result` carrying `Error.New(...)` when the input is not a valid integer. This
  is how values that overflow int64 are constructed.
- `Null() Interface` — the `NullBigNumber` Null-Object (`IsNull() == true`).
- `null.New() Interface` — the importable `NullBigNumber` Null-Object.

Conversions

- `ToGoString() string` (base-10 decimal), `ToInt64() int64`, `IsNull() bool`.

Arithmetic (each returns `Result.Interface`)

- `Add(other)` / `Sub(other)` / `Mul(other)` — sum, difference, product.
- `Div(other)` — quotient; a `Result` carrying `Error.New("division by zero")`
  when `other` is zero.
- `Mod(other)` — remainder; a `Result` carrying `Error.New("modulo by zero")`
  when `other` is zero.
- `Abs()` / `Neg()` — absolute value and negation.

Every operation works on a fresh `big.Int`, so operands are never mutated.

Comparisons (each returns `bool`)

- `Equal(other)` / `LessThan(other)` / `GreaterThan(other)`.

Inspection

- `Inspect() string` — `<BigNumber:0x... value=...>`.

## License

BSD-3-Clause — see [LICENSE](./LICENSE).
