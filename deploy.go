package main

import (
	"fmt"
	"github.com/go-acme/lego/v3/certificate"
	"github.com/hbagdi/go-kong/kong"
	"github.com/urfave/cli"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"os"
	"os/exec"
	"strings"
)

func deploy(ctx *cli.Context, cert *certificate.Resource) error {
	domains := ctx.GlobalStringSlice("domains")
	deployType := ctx.GlobalString("deploy")
	log.Printf("[INFO] Deployment certificate to %s.\n", deployType)
	switch deployType {
	case "secret":
		return deployIsTio(cert)
	case "kong":
		return deployKong(domains, cert)
	case "nginx":
		return deployNginx(cert)
	default:
		return fmt.Errorf("%q is not yet supported", deployType)
	}
}

// deploy to SECRET
func deployIsTio(cert *certificate.Resource) error {
	namespace := os.Getenv("NAMESPACE")
	secretName := os.Getenv("SECRETNAME")
	CERT := os.Getenv("CERT")
	KEY := os.Getenv("KEY")
	if len(namespace) == 0 || len(secretName) == 0 {
		return fmt.Errorf("[INFO] NAMESPACE OR SECRETNAME is not set")
	}
	if len(CERT) == 0 {
		CERT = "cert"
	}
	if len(KEY) == 0 {
		KEY = "key"
	}

	client, err := k8sClient()
	if err != nil {
		return err
	}
	secretsCli := client.CoreV1().Secrets(namespace)
	secret, err := secretsCli.Get(secretName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	secret.Data = map[string][]byte{
		CERT: cert.Certificate,
		KEY:  cert.PrivateKey,
	}
	_, err = secretsCli.Update(secret)
	return err
}

// deploy to kong
func deployKong(domains []string, cert *certificate.Resource) error {
	kongUrl := os.Getenv("KONG_ADMIN_ADDR")
	if len(kongUrl) == 0 {
		return fmt.Errorf("KONG_URL is not set")
	}
	client, err := kong.NewClient(kong.String(kongUrl), nil)
	if err != nil {
		return err
	}

	kongRoot, err := client.Root(nil)
	if err != nil {
		return err
	}

	listAll, err := client.Certificates.ListAll(nil)
	if err != nil {
		return err
	}
	domaNew := domains
	for _, v := range listAll {
		sins := stringValueSlice(v.SNIs)
		doma := domainsContains(domains, sins)
		if len(doma) == 0 {
			continue
		}
		domaNew = append(domaNew, doma...)
		v.SNIs = kong.StringSlice()
		for _, sin := range sins {
			if !contains(sin, doma) {
				v.SNIs = append(v.SNIs, kong.String(sin))
			}
		}
		if len(v.SNIs) != 0 {
			_, err := client.Certificates.Update(nil, v)
			if err != nil {
				return err
			}
			continue
		}
		err := client.Certificates.Delete(nil, v.ID)
		if err != nil {
			return err
		}
	}
	
	certs := &kong.Certificate{
		Cert: kong.String(string(cert.Certificate)),
		Key:  kong.String(string(cert.PrivateKey)),
		SNIs: stringSlice(sliceDedup(domaNew)),
	}

	newCert, err := client.Certificates.Create(nil, certs)
	if fmt.Sprintf("%v", kongRoot["version"]) > "1.0.0" {
		return err
	}
	
	// kong version < 1.0.x
	domainSlice := strings.Split(cert.Domain, ".")
	defautDomain := strings.Join(domainSlice[1:], ".")
	for _, v := range listAll {
		sins := stringValueSlice(v.SNIs)
		for _, s := range sins {
			sniDomainSlice := strings.Split(s, ".")
			sniDefaultDomain := strings.Join(sniDomainSlice[1:], ".")
			if len(sniDomainSlice) == len(domainSlice) && sniDefaultDomain == defautDomain {
				_, err = SinUpdate(client, nil, &SNI{
					Name:             kong.String(s),
					SslCertificateId: newCert.ID,
				})
			}
		}
	}
	for i:=0; i < len(domains); i++ {
		sni := &SNI{
			Name:             kong.String(domains[i]),
			SslCertificateId: newCert.ID,
		}
		exist, err := SinGet(client, nil, kong.String(domains[i]))
		if err != nil && err.Error() != "Not found" {
			log.Printf("[INFO] %s\n", err.Error())
		}
		if exist {
			_, err = SinUpdate(client, nil, sni)
		} else {
			_, err = SinCreate(client, nil, sni)
		}
	}
	return err
}

// deploy to nginx
func deployNginx(cert *certificate.Resource) error {
	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")
	if len(certPath) == 0 || len(keyPath) == 0 {
		return fmt.Errorf("CERT_PATH or KEY_PATH location is not set")
	}

	log.Printf("[INFO] Write certificate to %s.\n", certPath)
	if err := ioutil.WriteFile(certPath, cert.Certificate, 0600); err != nil {
		return err
	}

	log.Printf("[INFO] Write private key to %s. \n", keyPath)
	if err := ioutil.WriteFile(keyPath, cert.PrivateKey, 0600); err != nil {
		return err
	}

	log.Println("[INFO] Reload nginx server")
	_, err := exec.LookPath("nginx")
	if err != nil {
		return fmt.Errorf("didn't find 'nginx' executable\n")
	}

	cmd := exec.Command("nginx", "-s", "reload")
	_, err = cmd.CombinedOutput()
	return err
}
