package anyslice

import (
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAnySlice(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	RegisterFailHandler(Fail)
	RunSpecs(t, "AnySlice Suite")
}

var _ = Describe("AnySlice", func() {
	It("[]string", func() {
		base := []string{"a", "b", "c"}

		By("To []any")
		anySlice := ToAny(base)
		Expect(anySlice).To(Equal([]interface{}{"a", "b", "c"}))

		By("To []string")
		stringSlice, err := FromAny[string](anySlice)
		Expect(err).NotTo(HaveOccurred())
		Expect(stringSlice).To(Equal(base))
	})

	It("[]int", func() {
		base := []int{1, 2, 3}

		By("To []any")
		anySlice := ToAny(base)
		Expect(anySlice).To(Equal([]interface{}{1, 2, 3}))

		By("To []int")
		intSlice, err := FromAny[int](anySlice)
		Expect(err).NotTo(HaveOccurred())
		Expect(intSlice).To(Equal(base))
	})
})
