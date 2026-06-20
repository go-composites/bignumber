package BigNumber_test

import (
	BigNumber "github.com/go-composites/bignumber/src"
	Result "github.com/go-composites/result/src"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// payloadOf unwraps a success Result into a BigNumber.Interface.
func payloadOf(r interface {
	HasError() bool
	Payload() interface{}
}) BigNumber.Interface {
	gomega.ExpectWithOffset(1, r.HasError()).To(gomega.BeFalse())
	return r.Payload().(BigNumber.Interface)
}

// foreign is a BigNumber.Interface implementation that is NOT the package's
// own concrete type. It is used to exercise the string-bridging path of
// fromInterface with a value that DOES parse as base-10 (the success branch).
type foreign struct{ s string }

func (f foreign) ToGoString() string                     { return f.s }
func (foreign) ToInt64() int64                           { return 0 }
func (foreign) IsNull() bool                             { return false }
func (foreign) Add(BigNumber.Interface) Result.Interface { return nil }
func (foreign) Sub(BigNumber.Interface) Result.Interface { return nil }
func (foreign) Mul(BigNumber.Interface) Result.Interface { return nil }
func (foreign) Div(BigNumber.Interface) Result.Interface { return nil }
func (foreign) Mod(BigNumber.Interface) Result.Interface { return nil }
func (foreign) Abs() Result.Interface                    { return nil }
func (foreign) Neg() Result.Interface                    { return nil }
func (foreign) Equal(BigNumber.Interface) bool           { return false }
func (foreign) LessThan(BigNumber.Interface) bool        { return false }
func (foreign) GreaterThan(BigNumber.Interface) bool     { return false }
func (foreign) Inspect() BigNumber.String                { return `` }

var _ = ginkgo.Describe("BigNumber", func() {

	ginkgo.Describe("constructors", func() {
		ginkgo.It("builds from a Go int64", func() {
			n := BigNumber.FromInt64(42)
			gomega.Expect(n.ToInt64()).To(gomega.BeEquivalentTo(42))
			gomega.Expect(n.ToGoString()).To(gomega.Equal("42"))
			gomega.Expect(n.IsNull()).To(gomega.BeFalse())
		})
		ginkgo.It("parses a valid base-10 string", func() {
			r := BigNumber.FromString("12345678901234567890")
			gomega.Expect(r.HasError()).To(gomega.BeFalse())
			gomega.Expect(r.Payload().(BigNumber.Interface).ToGoString()).
				To(gomega.Equal("12345678901234567890"))
		})
		ginkgo.It("returns an error Result on a bad string", func() {
			r := BigNumber.FromString("not-a-number")
			gomega.Expect(r.HasError()).To(gomega.BeTrue())
			gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("invalid integer"))
		})
		ginkgo.It("exposes a Null-Object", func() {
			n := BigNumber.Null()
			gomega.Expect(n.IsNull()).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("arbitrary precision", func() {
		ginkgo.It("multiplies values that overflow int64", func() {
			big := payloadOf(BigNumber.FromString("9999999999999999999999999999999999999999"))
			r := big.Mul(big)
			gomega.Expect(r.HasError()).To(gomega.BeFalse())
			gomega.Expect(r.Payload().(BigNumber.Interface).ToGoString()).To(gomega.Equal(
				"99999999999999999999999999999999999999980000000000000000000000000000000000000001"))
		})
	})

	ginkgo.Describe("arithmetic", func() {
		var six = BigNumber.FromInt64(6)
		var two = BigNumber.FromInt64(2)

		ginkgo.It("adds", func() {
			gomega.Expect(payloadOf(six.Add(two)).ToInt64()).To(gomega.BeEquivalentTo(8))
		})
		ginkgo.It("subtracts", func() {
			gomega.Expect(payloadOf(six.Sub(two)).ToInt64()).To(gomega.BeEquivalentTo(4))
		})
		ginkgo.It("multiplies", func() {
			gomega.Expect(payloadOf(six.Mul(two)).ToInt64()).To(gomega.BeEquivalentTo(12))
		})
		ginkgo.It("divides", func() {
			gomega.Expect(payloadOf(six.Div(two)).ToInt64()).To(gomega.BeEquivalentTo(3))
		})
		ginkgo.It("computes the remainder", func() {
			gomega.Expect(payloadOf(six.Mod(BigNumber.FromInt64(4))).ToInt64()).
				To(gomega.BeEquivalentTo(2))
		})
		ginkgo.It("does not mutate its operands", func() {
			_ = six.Add(two)
			gomega.Expect(six.ToInt64()).To(gomega.BeEquivalentTo(6))
			gomega.Expect(two.ToInt64()).To(gomega.BeEquivalentTo(2))
		})

		ginkgo.Describe("division by zero", func() {
			ginkgo.It("returns a Result carrying an error instead of panicking", func() {
				r := six.Div(BigNumber.FromInt64(0))
				gomega.Expect(r.HasError()).To(gomega.BeTrue())
				gomega.Expect(r.Error().Message()).To(gomega.Equal("division by zero"))
			})
		})

		ginkgo.Describe("modulo by zero", func() {
			ginkgo.It("returns a Result carrying an error instead of panicking", func() {
				r := six.Mod(BigNumber.FromInt64(0))
				gomega.Expect(r.HasError()).To(gomega.BeTrue())
				gomega.Expect(r.Error().Message()).To(gomega.Equal("modulo by zero"))
			})
		})

		ginkgo.Describe("absolute value", func() {
			ginkgo.It("makes a negative number positive", func() {
				gomega.Expect(payloadOf(BigNumber.FromInt64(-7).Abs()).ToInt64()).
					To(gomega.BeEquivalentTo(7))
			})
			ginkgo.It("leaves a positive number unchanged", func() {
				gomega.Expect(payloadOf(BigNumber.FromInt64(7).Abs()).ToInt64()).
					To(gomega.BeEquivalentTo(7))
			})
		})

		ginkgo.Describe("negation", func() {
			ginkgo.It("negates a positive number", func() {
				gomega.Expect(payloadOf(BigNumber.FromInt64(7).Neg()).ToInt64()).
					To(gomega.BeEquivalentTo(-7))
			})
			ginkgo.It("negates a negative number", func() {
				gomega.Expect(payloadOf(BigNumber.FromInt64(-7).Neg()).ToInt64()).
					To(gomega.BeEquivalentTo(7))
			})
		})
	})

	ginkgo.Describe("operations against a Null operand", func() {
		var six = BigNumber.FromInt64(6)
		var null = BigNumber.Null()

		ginkgo.It("treats a Null operand as zero in addition", func() {
			gomega.Expect(payloadOf(six.Add(null)).ToInt64()).To(gomega.BeEquivalentTo(6))
		})
		ginkgo.It("guards division by a Null operand (zero)", func() {
			gomega.Expect(six.Div(null).HasError()).To(gomega.BeTrue())
		})
		ginkgo.It("bridges a foreign Interface through its decimal string", func() {
			gomega.Expect(payloadOf(six.Add(foreign{s: "4"})).ToInt64()).
				To(gomega.BeEquivalentTo(10))
		})
		ginkgo.It("treats an unparsable foreign operand as zero", func() {
			gomega.Expect(payloadOf(six.Add(foreign{s: "xx"})).ToInt64()).
				To(gomega.BeEquivalentTo(6))
		})
	})

	ginkgo.Describe("comparisons", func() {
		var six = BigNumber.FromInt64(6)
		var two = BigNumber.FromInt64(2)

		ginkgo.It("reports equality both ways", func() {
			gomega.Expect(six.Equal(six)).To(gomega.BeTrue())
			gomega.Expect(six.Equal(two)).To(gomega.BeFalse())
		})
		ginkgo.It("reports less-than both ways", func() {
			gomega.Expect(two.LessThan(six)).To(gomega.BeTrue())
			gomega.Expect(six.LessThan(two)).To(gomega.BeFalse())
		})
		ginkgo.It("reports greater-than both ways", func() {
			gomega.Expect(six.GreaterThan(two)).To(gomega.BeTrue())
			gomega.Expect(two.GreaterThan(six)).To(gomega.BeFalse())
		})
	})

	ginkgo.Describe("inspection", func() {
		ginkgo.It("renders a BigNumber", func() {
			gomega.Expect(BigNumber.FromInt64(6).Inspect()).
				To(gomega.ContainSubstring("value=6"))
		})
	})

	ginkgo.Describe("the package-local Null-Object", func() {
		var n = BigNumber.Null()

		ginkgo.It("converts to zero values", func() {
			gomega.Expect(n.ToInt64()).To(gomega.BeEquivalentTo(0))
			gomega.Expect(n.ToGoString()).To(gomega.Equal(``))
		})
		ginkgo.It("returns error Results for every arithmetic method", func() {
			gomega.Expect(n.Add(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Sub(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Mul(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Div(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Mod(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Abs().HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Neg().HasError()).To(gomega.BeTrue())
		})
		ginkgo.It("compares as a Null-Object", func() {
			gomega.Expect(n.Equal(BigNumber.Null())).To(gomega.BeTrue())
			gomega.Expect(n.Equal(BigNumber.FromInt64(0))).To(gomega.BeFalse())
			gomega.Expect(n.LessThan(BigNumber.FromInt64(1))).To(gomega.BeFalse())
			gomega.Expect(n.GreaterThan(BigNumber.FromInt64(-1))).To(gomega.BeFalse())
		})
		ginkgo.It("inspects as the null marker", func() {
			gomega.Expect(n.Inspect()).To(gomega.Equal(`<NullBigNumber>`))
		})
	})
})
