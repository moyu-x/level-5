package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/moyu-x/level-5/pkg/config"
)

var (
	instance     *tls.Config
	instanceOnce sync.Once
)

func Tls(c *config.Bootstrap) *tls.Config {
	instanceOnce.Do(func() {
		instance = buildTlsConfig(c)
	})
	return instance
}

func buildTlsConfig(c *config.Bootstrap) *tls.Config {
	certFile := readFile(c.Kafka.CertFilePath, "kafka cert file")
	keyFile := readFile(c.Kafka.KeyFilePath, "kafka key file")
	keyPair, err := tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		log.Error().Err(err).Msg("failed to build x509 key pair")
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(certFile)

	return &tls.Config{
		RootCAs:            certPool,
		Certificates:       []tls.Certificate{keyPair},
		InsecureSkipVerify: true, //nolint:gosec
	}
}

func readFile(path, description string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Msgf("failed to read %s", description)
	}
	return data
}
