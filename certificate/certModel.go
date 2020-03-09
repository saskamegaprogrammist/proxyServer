package certificate

import "crypto/tls"

type Cert struct {
	tlsCert tls.Certificate
}