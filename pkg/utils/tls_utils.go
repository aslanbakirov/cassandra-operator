package utils

import(
	"crypto/tls"
	"io/ioutil"
	"os"
	"path/filepath"
	"fmt"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
	"net"

)

const (
	CliCertFile = "operator-client.crt"
	CliKeyFile  = "operator-client.key"
	CliCAFile   = "operator-client-ca.crt"
)

type TLSInfo struct {
	CertFile       string
	KeyFile        string
	CAFile         string
	TrustedCAFile  string
	ClientCertAuth bool

	// ServerName ensures the cert matches the given host in case of discovery / virtual hosting
	ServerName string

	selfCert bool

	parseFunc func([]byte, []byte) (tls.Certificate, error)
}

// ClientConfig generates a tls.Config object for use by an HTTP client.
func (info TLSInfo) ClientConfig() (*tls.Config, error) {
	var cfg *tls.Config
	var err error

	if !info.Empty() {
		cfg, err = info.baseConfig()
		if err != nil {
			return nil, err
		}
	} else {
		cfg = &tls.Config{ServerName: info.ServerName}
	}

	CAFiles := info.cafiles()
	if len(CAFiles) > 0 {
		cfg.RootCAs, err = NewCertPool(CAFiles)
		if err != nil {
			return nil, err
		}
		// if given a CA, trust any host with a cert signed by the CA
		//log.Println("warning: ignoring ServerName for user-provided CA for backwards compatibility is deprecated")
		cfg.ServerName = ""
	}

	if info.selfCert {
		cfg.InsecureSkipVerify = true
	}
	return cfg, nil
}

// NewCertPool creates x509 certPool with provided CA files.
func NewCertPool(CAFiles []string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	for _, CAFile := range CAFiles {
		pemByte, err := ioutil.ReadFile(CAFile)
		if err != nil {
			return nil, err
		}

		for {
			var block *pem.Block
			block, pemByte = pem.Decode(pemByte)
			if block == nil {
				break
			}
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			certPool.AddCert(cert)
		}
	}

	return certPool, nil
}


func NewTLSConfig(certData, keyData, caData []byte) (*tls.Config, error) {
	dir, err := ioutil.TempDir("", "cassandra-operator-cluster-tls")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	certFile, err := writeFile(dir, CliCertFile, certData)
	if err != nil {
		return nil, err
	}
	keyFile, err := writeFile(dir, CliKeyFile, keyData)
	if err != nil {
		return nil, err
	}
	caFile, err := writeFile(dir, CliCAFile, caData)
	if err != nil {
		return nil, err
	}

	tlsInfo := TLSInfo{
		CertFile:      certFile,
		KeyFile:       keyFile,
		TrustedCAFile: caFile,
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}
	return tlsConfig, nil
}


func (info TLSInfo) String() string {
	return fmt.Sprintf("cert = %s, key = %s, ca = %s, trusted-ca = %s, client-cert-auth = %v", info.CertFile, info.KeyFile, info.CAFile, info.TrustedCAFile, info.ClientCertAuth)
}

func (info TLSInfo) Empty() bool {
	return info.CertFile == "" && info.KeyFile == ""
}

func SelfCert(dirpath string, hosts []string) (info TLSInfo, err error) {
	if err = TouchDirAll(dirpath); err != nil {
		return
	}

	certPath := filepath.Join(dirpath, "cert.pem")
	keyPath := filepath.Join(dirpath, "key.pem")
	_, errcert := os.Stat(certPath)
	_, errkey := os.Stat(keyPath)
	if errcert == nil && errkey == nil {
		info.CertFile = certPath
		info.KeyFile = keyPath
		info.selfCert = true
		return
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return
	}

	tmpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{Organization: []string{"cassandra"}},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * (24 * time.Hour)),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, host := range hosts {
		h, _, _ := net.SplitHostPort(host)
		if ip := net.ParseIP(h); ip != nil {
			tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
		} else {
			tmpl.DNSNames = append(tmpl.DNSNames, h)
		}
	}

	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if err != nil {
		return
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		return
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return
	}
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	keyOut.Close()

	return SelfCert(dirpath, hosts)
}

func (info TLSInfo) baseConfig() (*tls.Config, error) {
	if info.KeyFile == "" || info.CertFile == "" {
		return nil, fmt.Errorf("KeyFile and CertFile must both be present[key: %v, cert: %v]", info.KeyFile, info.CertFile)
	}

	tlsCert, err := NewCert(info.CertFile, info.KeyFile, info.parseFunc)
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{*tlsCert},
		MinVersion:   tls.VersionTLS12,
		ServerName:   info.ServerName,
	}
	return cfg, nil
}

// cafiles returns a list of CA file paths.
func (info TLSInfo) cafiles() []string {
	cs := make([]string, 0)
	if info.CAFile != "" {
		cs = append(cs, info.CAFile)
	}
	if info.TrustedCAFile != "" {
		cs = append(cs, info.TrustedCAFile)
	}
	return cs
}

// ServerConfig generates a tls.Config object for use by an HTTP server.
func (info TLSInfo) ServerConfig() (*tls.Config, error) {
	cfg, err := info.baseConfig()
	if err != nil {
		return nil, err
	}

	cfg.ClientAuth = tls.NoClientCert
	if info.CAFile != "" || info.ClientCertAuth {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	CAFiles := info.cafiles()
	if len(CAFiles) > 0 {
		cp, err := NewCertPool(CAFiles)
		if err != nil {
			return nil, err
		}
		cfg.ClientCAs = cp
	}

	// "h2" NextProtos is necessary for enabling HTTP2 for go's HTTP server
	cfg.NextProtos = []string{"h2"}

	return cfg, nil
}

// ClientConfig generates a tls.Config object for use by an HTTP client.
// ShallowCopyTLSConfig copies *tls.Config. This is only
// work-around for go-vet tests, which complains
//
//   assignment copies lock value to p: crypto/tls.Config contains sync.Once contains sync.Mutex
//
// Keep up-to-date with 'go/src/crypto/tls/common.go'
func ShallowCopyTLSConfig(cfg *tls.Config) *tls.Config {
	ncfg := tls.Config{
		Time:                     cfg.Time,
		Certificates:             cfg.Certificates,
		NameToCertificate:        cfg.NameToCertificate,
		GetCertificate:           cfg.GetCertificate,
		RootCAs:                  cfg.RootCAs,
		NextProtos:               cfg.NextProtos,
		ServerName:               cfg.ServerName,
		ClientAuth:               cfg.ClientAuth,
		ClientCAs:                cfg.ClientCAs,
		InsecureSkipVerify:       cfg.InsecureSkipVerify,
		CipherSuites:             cfg.CipherSuites,
		PreferServerCipherSuites: cfg.PreferServerCipherSuites,
		SessionTicketKey:         cfg.SessionTicketKey,
		ClientSessionCache:       cfg.ClientSessionCache,
		MinVersion:               cfg.MinVersion,
		MaxVersion:               cfg.MaxVersion,
		CurvePreferences:         cfg.CurvePreferences,
	}
	return &ncfg
}

func NewCert(certfile, keyfile string, parseFunc func([]byte, []byte) (tls.Certificate, error)) (*tls.Certificate, error) {
	cert, err := ioutil.ReadFile(certfile)
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	if parseFunc == nil {
		parseFunc = tls.X509KeyPair
	}

	tlsCert, err := parseFunc(cert, key)
	if err != nil {
		return nil, err
	}
	return &tlsCert, nil
}


func writeFile(dir, file string, data []byte) (string, error) {
	p := filepath.Join(dir, file)
	return p, ioutil.WriteFile(p, data, 0600)
}

