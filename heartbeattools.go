package servertools

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Service string

const (
	Gateway    Service = "festivals-gateway"
	Identity   Service = "festivals-identity-server"
	Server     Service = "festivals-server"
	Fileserver Service = "festivals-fileserver"
	Database   Service = "festivals-database-node"
	Website    Service = "festivals-website-node"
)

type Heartbeat struct {
	Service   string `json:"service"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Available bool   `json:"available"`
}

func HeartbeatClient(clientCert string, clientKey string, serverCA string) (*http.Client, error) {

	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	certContent, err := os.ReadFile(serverCA)
	if err != nil {
		return nil, err
	}
	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM(certContent); !ok {
		return nil, errors.New("failed to append certificate to certificate pool")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      rootCertPool,
			},
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 20 * time.Second,
	}

	return client, nil
}

func SendHeartbeat(client *http.Client, endpoint string, serviceKey Service, beat *Heartbeat) error {

	heartbeatwave, err := json.Marshal(beat)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(heartbeatwave))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("X-Request-ID", uuid.New().String())
	request.Header.Set("Service-Key", string(serviceKey))

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
