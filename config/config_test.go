package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gnames/gnlib/format"
	. "github.com/gnames/gnverify/config"
)

var url = "https://gnames.globalnames.org"

var _ = Describe("Config", func() {
	Describe("NewConfig", func() {
		It("Creates a default GNparser", func() {
			cnf := NewConfig()
			deflt := Config{
				Format:      format.CSV,
				VerifierURL: "http://:8888",
			}
			Expect(cnf.Format).To(Equal(deflt.Format))
			Expect(cnf.VerifierURL).To(Equal(deflt.VerifierURL))
		})
	})

	It("Takes options to update default settings", func() {
		opts := opts()
		cnf := NewConfig(opts...)
		updt := Config{
			Format:           format.PrettyJSON,
			PreferredOnly:    true,
			NameField:        3,
			PreferredSources: []int{1, 2, 3},
			VerifierURL:      url,
		}
		Expect(cnf.Format).To(Equal(updt.Format))
		Expect(cnf.PreferredOnly).To(Equal(updt.PreferredOnly))
		Expect(cnf.NameField).To(Equal(updt.NameField))
		Expect(cnf.PreferredSources).To(Equal(updt.PreferredSources))
		Expect(cnf.VerifierURL).To(Equal(updt.VerifierURL))
	})

	Describe("NewFormat", func() {
		It("Creates format out of string", func() {
			inputs := []formatTest{
				{String: "csv", Format: format.CSV},
				{String: "compact", Format: format.CompactJSON},
				{String: "pretty", Format: format.PrettyJSON},
				// {String: "badstring", Format: format.FormatNone},
			}
			for _, v := range inputs {
				f, _ := format.NewFormat(v.String)
				Expect(f).To(Equal(v.Format))
			}
		})
	})
})

type formatTest struct {
	String string
	format.Format
}

func opts() []Option {
	return []Option{
		OptFormat(format.PrettyJSON),
		OptPreferredOnly(true),
		OptNameField(3),
		OptPreferredSources([]int{1, 2, 3}),
		OptVerifierURL(url),
	}
}
