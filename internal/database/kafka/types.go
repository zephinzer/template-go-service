package kafka

type ConnectionOpts struct {
	TLS *TlsOpts
}

type TlsOpts struct {
	CaCertificatePath     string
	ClientCertificatePath string
	ClientKeyPath         string
}
