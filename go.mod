module github.com/gnames/gnverify

go 1.15

require (
	github.com/gnames/gnames v0.0.0-00010101000000-000000000000
	github.com/gnames/gnlib v0.1.2
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/nxadm/tail v1.4.5 // indirect
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.6.1
)

replace github.com/gnames/gnlib => ../gnlib

replace github.com/gnames/gnames => ../gnames
