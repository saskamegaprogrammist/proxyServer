package certificate

import (
	"crypto/x509"
	"time"
)

const (
	MAXCERTAGE = 5 * 365 * 24 * time.Hour
	MAXLEAFAGE = 24 * time.Hour
	USAGE = x509.KeyUsageDigitalSignature |
		x509.KeyUsageContentCommitment |
		x509.KeyUsageKeyEncipherment |
		x509.KeyUsageDataEncipherment |
		x509.KeyUsageKeyAgreement |
		x509.KeyUsageCertSign |
		x509.KeyUsageCRLSign
	CERTTYPE = "CERTIFICATE"
	KEYTYPE = "ECDSA PRIVATE KEY"
	ROOTCERTFILENAME = "rootCert.cert"
	ROOTKEYFILENAME = "rootKey.pem"
	PERMISSIONS = 0644

)
