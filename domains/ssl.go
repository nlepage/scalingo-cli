package domains

import (
	"os"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func DisableSSL(app string, domain string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	d, err := findDomain(c, app, domain)
	if err != nil {
		return errgo.Notef(err, "fail to find the matching domain to disable SSL")
	}

	_, err = c.DomainsUpdate(app, d.ID, "", "")
	if err != nil {
		return errgo.Notef(err, "fail to unset the domain certificate")
	}
	io.Status("SSL of " + domain + " has been disabled.")
	return nil
}

func EnableSSL(app, domain, certPath, keyPath string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	d, err := findDomain(c, app, domain)
	if err != nil {
		return errgo.Notef(err, "fail to find the matching domain to enable SSL")
	}

	certContent, keyContent, err := validateSSL(certPath, keyPath)
	if err != nil {
		return errgo.Notef(err, "fail to validate the given certificate and key")
	}

	d, err = c.DomainsUpdate(app, d.ID, certContent, keyContent)
	if err != nil {
		return errgo.Notef(err, "fail to set the domain certificate")
	}

	io.Status("The certificate and key have been installed for " + d.Name + " (Validity: " + d.Validity.UTC().String() + ")")
	return nil
}

func validateSSL(cert, key string) (string, string, error) {
	if cert == "" && key == "" {
		return "", "", nil
	}

	if cert == "" && key != "" {
		return "", "", errgo.New("--cert <certificate path> should be defined")
	}

	if key == "" && cert != "" {
		return "", "", errgo.New("--key <key path> should be defined")
	}

	certContent, err := os.ReadFile(cert)
	if err != nil {
		return "", "", errgo.Notef(err, "fail to read the TLS certificate")
	}
	keyContent, err := os.ReadFile(key)
	if err != nil {
		return "", "", errgo.Notef(err, "fail to read the private key")
	}
	return string(certContent), string(keyContent), nil
}
