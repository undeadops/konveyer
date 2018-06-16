package deployments

import (
	"fmt"

	restclient "k8s.io/client-go/rest"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

// NewClient - Return New Kubernetes Client
func NewClient(ke []byte) (*restclient.Config, error) {

	rawconfig, err := clientcmd.Load(ke)
	if err != nil {
		return nil, fmt.Errorf("Unable to create RestClient Config, %s", err.Error())
	}

	clientconfig := clientcmd.NewDefaultClientConfig(*rawconfig, nil)
	config, err := clientconfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Unable to create RestClient Config, %v", err.Error())
	}
	// if ke.Cluster.ServerUri == "" {
	// 	return nil, fmt.Errorf("no cluster endpoint provided")
	// }

	// tlsConfig := &restclient.TLSClientConfig{
	// 	Insecure: ke.Cluster.InSecure,
	// 	ServerName: ke.Cluster.Name,
	// 	CertFile: ke.UserAuth.ClientCertificate,
	// 	KeyFile: ke.UserAuth.ClientKey,
	// 	CAFile: ke.Cluster.CertificateAuthority
	// }

	// ca, err := load(ke.Cluster.CertificateAuthority, ke.Cluster.CertificateAuthorityData)
	// if err != nil {
	// 	return nil, fmt.Errorf("loading certificate authority: %v", err)
	// }

	// clientCert, err := load(ke.UserAuth.ClientCertificate, ke.UserAuth.ClientCertificateData)
	// if err != nil {
	// 	return nil, fmt.Errorf("loading client cert: %v", err)
	// }

	// clientKey, err := load(ke.UserAuth.ClientKey, ke.UserAuth.ClientKeyData)
	// if err != nil {
	// 	return nil, fmt.Errorf("loading client key: %v", err)
	// }

	// if len(clientCert) != 0 {
	// 	tlsConfig.CertData = clientCert
	// }
	// if len(clientKey) != 0 {
	// 	tlsConfig.KeyData = clientKey
	// }
	// if len(ca) != 0 {
	// 	tlsConfig.CAData = ca
	// }
	// if len(ca) != 0 {
	// 	tlsConfig.RootCAs = x509.NewCertPool()
	// 	if !tlsConfig.RootCAs.AppendCertsFromPEM(ca) {
	// 		return nil, errors.New("Certificate Authority doesn't contain any certificates")
	// 	}
	// }

	// if len(clientCert) != 0 {
	// 	cert, err := tls.X509KeyPair(clientCert, clientKey)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("Invalid Client Cert and Key Pair: %v", err)
	// 	}
	// 	tlsConfig.Certificates = []tls.Certificate{cert}
	// }

	// transport := &http.Transport{
	// 	Proxy: http.ProxyFromEnvironment,
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   30 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 	}).DialContext,
	// 	TLSClientConfig:       tlsConfig,
	// 	MaxIdleConns:          100,
	// 	IdleConnTimeout:       90 * time.Second,
	// 	TLSHandshakeTimeout:   10 * time.Second,
	// 	ExpectContinueTimeout: 1 * time.Second,
	// }
	// if err := http2.ConfigureTransport(transport); err != nil {
	// 	return nil, err
	// }

	// client := &k8s.Client{
	// 	Endpoint:  ke.Cluster.ServerUri,
	// 	Namespace: ke.Namespace,
	// 	Client: &http.Client{
	// 		Transport: transport,
	// 	},
	// }

	// if ke.UserAuth.Token != "" {
	// 	client.SetHeaders = func(h http.Header) error {
	// 		h.Set("Authorization", "Bearer "+ke.UserAuth.Token)
	// 		return nil
	// 	}
	// }

	// if ke.UserAuth.Username != "" && ke.UserAuth.Password != "" {
	// 	auth := ke.UserAuth.Username + ":" + ke.UserAuth.Password
	// 	auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	// 	client.SetHeaders = func(h http.Header) error {
	// 		h.Set("Authorization", auth)
	// 		return nil
	// 	}
	// }

	return config, nil
}
