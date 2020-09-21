package gnverify

import "gitlab.com/gogna/gnverify/config"

type GNVerify struct {
	config.Config
}

func NewGNVerify(cnf config.Config) GNVerify {
	return GNVerify{
		Config: cnf,
	}
}
