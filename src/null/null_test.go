package NullBigNumber_test

import (
	BigNumber "github.com/go-composites/bignumber/src"
	NullBigNumber "github.com/go-composites/bignumber/src/null"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("NullBigNumber", func() {
	var n NullBigNumber.Interface
	ginkgo.BeforeEach(func() {
		n = NullBigNumber.New()
	})

	ginkgo.It("satisfies the BigNumber interface", func() {
		var _ BigNumber.Interface = n
	})
	ginkgo.It("reports IsNull() true", func() {
		gomega.Expect(n.IsNull()).To(gomega.BeTrue())
	})
	ginkgo.It("converts to zero values", func() {
		gomega.Expect(n.ToInt64()).To(gomega.BeEquivalentTo(0))
		gomega.Expect(n.ToGoString()).To(gomega.Equal(``))
	})

	ginkgo.It("Add returns an error result", func() {
		r := n.Add(BigNumber.FromInt64(0))
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Add"))
	})
	ginkgo.It("Sub returns an error result", func() {
		r := n.Sub(BigNumber.FromInt64(0))
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Sub"))
	})
	ginkgo.It("Mul returns an error result", func() {
		r := n.Mul(BigNumber.FromInt64(0))
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Mul"))
	})
	ginkgo.It("Div returns an error result", func() {
		r := n.Div(BigNumber.FromInt64(0))
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Div"))
	})
	ginkgo.It("Mod returns an error result", func() {
		r := n.Mod(BigNumber.FromInt64(0))
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Mod"))
	})
	ginkgo.It("Abs returns an error result", func() {
		r := n.Abs()
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Abs"))
	})
	ginkgo.It("Neg returns an error result", func() {
		r := n.Neg()
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Neg"))
	})
	ginkgo.It("Equal is true only against another null", func() {
		gomega.Expect(n.Equal(NullBigNumber.New())).To(gomega.BeTrue())
		gomega.Expect(n.Equal(BigNumber.FromInt64(0))).To(gomega.BeFalse())
	})
	ginkgo.It("LessThan is always false", func() {
		gomega.Expect(n.LessThan(BigNumber.FromInt64(0))).To(gomega.BeFalse())
	})
	ginkgo.It("GreaterThan is always false", func() {
		gomega.Expect(n.GreaterThan(BigNumber.FromInt64(0))).To(gomega.BeFalse())
	})
	ginkgo.It("Inspect renders the null marker", func() {
		gomega.Expect(n.Inspect()).To(gomega.Equal(`<NullBigNumber>`))
	})
})
