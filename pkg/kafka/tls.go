package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"sync"

	"github.com/go-logr/logr"

	"github.com/moyu-x/level-5/pkg/config"
)

var (
	instance     *tls.Config
	instanceOnce sync.Once
)

func Tls(l *logr.Logger, c *config.Bootstrap) *tls.Config {
	instanceOnce.Do(func() {
		instance = buildTlsConfig(l, c)
	})
	return instance
}

func buildTlsConfig(l *logr.Logger, c *config.Bootstrap) *tls.Config {
	keyPair, err := tls.X509KeyPair(getCertFile(l, c), getKeyFile(l, c))
	if err != nil {
		l.Error(err, "build x509 key pair has error, reason: %v")
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(getCertFile(l, c))
	return &tls.Config{
		RootCAs:            certPool,
		Certificates:       []tls.Certificate{keyPair},
		InsecureSkipVerify: true, //nolint:gosec
	}
}

func getCertFile(l *logr.Logger, c *config.Bootstrap) []byte {
	certFile, err := os.ReadFile(c.Kafka.CertFilePath)
	if err != nil {
		l.Error(err, "read kafka cert file error, reason")
	}
	return certFile
}

// kafka.client.key
func getKeyFile(l *logr.Logger, c *config.Bootstrap) []byte {
	keyFile, err := os.ReadFile(c.Kafka.KeyFilePath)
	if err != nil {
		l.Error(err, "read kafka key file has error, reason")
	}
	return keyFile
}
