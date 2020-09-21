package config_test

import (
	log "github.com/sirupsen/logrus"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	log.SetLevel(log.PanicLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}
