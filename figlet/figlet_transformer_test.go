package figlet_test

import (
	"github.com/calebamiles/example-figlet-service/figlet"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("figlet.Transformer", func() {
	Describe("Figletize()", func() {
		It("applies a figlet transformation", func() {
			txt := "cool"
			figletedTxt, err := figlet.NewTransformer().Figletize(txt)

			Expect(err).ToNot(HaveOccurred(), "expected no error ")
			Expect(figletedTxt).ToNot(BeEmpty(), "expected figleted text to be non empty")

			Expect(len(figletedTxt)).To(BeNumerically(">", len(txt)), "expected figetlet transformation to be strictly larger than input")
		})
	})
})
