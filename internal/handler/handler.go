package handler

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/gin-gonic/gin"
	"goTSVParser/config"
	"goTSVParser/internal/domains"
	"goTSVParser/internal/shema"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"
)

type Handler struct {
	service domains.Service
	engine  *gin.Engine
	config  config.Config
}

func NewHandler(service domains.Service, cnf config.Config) *Handler {
	router := gin.Default()
	h := &Handler{
		service: service,
		engine:  router,
		config:  cnf,
	}
	Route(router, h)
	return h
}

func (s *Handler) Start() {
	if s.config.TLS && (s.config.PrivateKey == "" || s.config.Certificate == "") {
		cert := &x509.Certificate{
			SerialNumber: big.NewInt(2024),
			Subject: pkix.Name{
				Organization: []string{"andogeek"},
				Country:      []string{"USA"},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(10, 0, 0),
			IsCA:                  true,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			BasicConstraintsValid: true,
		}

		privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			log.Fatal(err)
		}

		certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
		if err != nil {
			log.Fatal(err)
		}

		var certPEM bytes.Buffer
		pem.Encode(&certPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		})

		var privateKeyPEM bytes.Buffer
		pem.Encode(&privateKeyPEM, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})
		cerPemFile, err := os.Create("cerPem.crt")
		if err != nil {
			log.Fatalf("Failed to create file: %s", err.Error())
		}
		_, err = cerPemFile.Write(certPEM.Bytes())
		if err != nil {
			log.Fatalf("Failed to write file: %s", err.Error())
		}

		defer cerPemFile.Close()

		privateFile, err := os.Create("private.key")
		if err != nil {
			log.Fatalf("Failed to create file: %s", err.Error())
		}

		_, err = privateFile.Write(privateKeyPEM.Bytes())
		if err != nil {
			log.Fatalf("Failed to write file: %s", err.Error())
		}

		defer privateFile.Close()

		http.ListenAndServeTLS(s.config.Host, "cerPem.crt", "private.key", s.engine.Handler())
	} else if s.config.TLS && s.config.Certificate != "" && s.config.PrivateKey != "" {
		http.ListenAndServeTLS(s.config.Host, s.config.Certificate, s.config.PrivateKey, s.engine.Handler())
	}
	s.engine.Run(s.config.Host)
}

func (s *Handler) GetAll(c *gin.Context) {
	var r shema.Request
	err := c.ShouldBindJSON(&r)
	if err != nil {
		HandlerErr(c, err)
		return
	}
	ctx := c.Request.Context()
	result, err := s.service.GetAll(ctx, r)
	if err != nil {
		HandlerErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)

}
