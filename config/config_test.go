package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gnames/gnverify/config"
)

var url = "https://gnames.globalnames.org"

var _ = Describe("Config", func() {
	Describe("NewConfig", func() {
		It("Creates a default GNparser", func() {
			cnf := NewConfig()
			deflt := Config{
				Format:      CSV,
				VerifierURL: "https://:8888",
			}
			Expect(cnf).To(Equal(deflt))
		})
	})

	It("Takes options to update default settings", func() {
		opts := opts()
		cnf := NewConfig(opts...)
		updt := Config{
			Format:           PrettyJSON,
			PreferredOnly:    true,
			NameField:        3,
			PreferredSources: []int{1, 2, 3},
			VerifierURL:      url,
		}
		Expect(cnf).To(Equal(updt))
	})

	Describe("NewFormat", func() {
		It("Creates format out of string", func() {
			inputs := []formatTest{
				{String: "csv", Format: CSV},
				{String: "compact", Format: CompactJSON},
				{String: "pretty", Format: PrettyJSON},
				{String: "badstring", Format: InvalidFormat},
			}
			for _, v := range inputs {
				Expect(NewFormat(v.String)).To(Equal(v.Format))
			}
		})
	})
})

type formatTest struct {
	String string
	Format
}

func opts() []Option {
	return []Option{
		OptFormat(PrettyJSON),
		OptPreferredOnly(true),
		OptNameField(3),
		OptPreferredSources([]int{1, 2, 3}),
		OptVerifierURL(url),
	}
}
