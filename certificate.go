package apns

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"golang.org/x/crypto/pkcs12"
)

// Декодирует и проверяет валидность сертификата в формат p12.
func Certificate(p12 []byte, password string) (tlsCert tls.Certificate, err error) {
	privateKey, cert, err := pkcs12.Decode(p12, password) // декодируем сертификат
	if err != nil {
		return
	}
	// Проверяем валидность сертификата. Игнорируем ошибку UnknownAuthority, т.к. сертификаты
	// Apple не входят в основной пул сертификатов.
	if _, err = cert.Verify(x509.VerifyOptions{}); err != nil {
		if _, ok := err.(x509.UnknownAuthorityError); !ok {
			return
		}
	}
	tlsCert = tls.Certificate{ // формируем описание сертификата
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  privateKey,
		Leaf:        cert,
	}
	return
}

// LoadTLSCertificate загружает и разбирает сертификат в формате p12 из файла. Возвращает разобранный
// сертификат или ошибку.
func LoadCertificate(filename, password string) (tlsCert tls.Certificate, err error) {
	p12, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	return Certificate(p12, password)
}
