package main

import (
	"fmt"
	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/hbagdi/go-kong/kong"
	"github.com/urfave/cli"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"os"
	"time"
)

func checkExpirationTime(ctx *cli.Context) bool {
	deployType := ctx.GlobalString("deploy")
	domains := ctx.GlobalStringSlice("domains")
	days := ctx.GlobalInt("days")
	switch deployType {
	case "istio":
		client, err := k8sClient()
		if err != nil {
			log.Printf("[ERROR] kubernetes connection failed. %s\n", err)
			return false
		}
		CERT := os.Getenv("CERT")
		namespace := os.Getenv("NAMESPACE")
		secretName := os.Getenv("SECRETNAME")
		if len(namespace) == 0 || len(secretName) == 0 {
			fmt.Println("[ERROR] NAMESPACE OR CONFMAP is not set")
			return false
		}
		if len(CERT) == 0 {
			CERT = "cert"
		}
		secretsCli := client.CoreV1().Secrets(namespace)
		secret, err := secretsCli.Get(secretName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("[WARNING] istio certificate does not exist. %s\n", err)
			return true
		}
		if len(secret.Data[CERT]) == 0 {
			return true
		}
		return checkCA(string(secret.Data[CERT]), days)
	case "kong":
		// Timed tasks that depend on k8s
		kongUrl := os.Getenv("KONG_ADMIN_ADDR")
		if len(kongUrl) == 0 {
			log.Println("[ERROR] KONG_ADMIN_ADDR is not set")
			return false
		}
		client, err := kong.NewClient(kong.String(kongUrl), nil)
		if err != nil {
			log.Printf("[ERROR] %s\n", err)
			return false
		}
		listAll, err := client.Certificates.ListAll(nil)
		if err != nil {
			log.Printf("[ERROR] %s\n", err)
			return false
		}
		var certs []string
		for _, v := range listAll {
			sins := stringValueSlice(v.SNIs)
			doma := domainsContains(domains, sins)
			if len(doma) == 0 {
				continue
			}
			certs = append(certs, *v.Cert)
		}
		certs = sliceDedup(certs)
		if len(certs) == 0 {
			return true
		}
		
		var expiredSlice []bool
		for i:=0; i < len(certs); i++ {
			expiredSlice = append(expiredSlice, checkCA(certs[i], days))
		}
		for i:=0; i < len(expiredSlice); i ++ {
			if expiredSlice[i] {
				return true
			}
		}
		return false
	case "nginx":
		certPath := os.Getenv("CERT_PATH")
		keyPath := os.Getenv("KEY_PATH")
		if len(certPath) == 0 || len(keyPath) == 0 {
			log.Println("[ERROR] CERT_PATH or KEY_PATH location is not set")
			return false
		}
		cert, err := ioutil.ReadFile(certPath)
		if err != nil {
			log.Println("[WARNING] nginx certificate does not exist")
		}
		if len(cert) == 0 {
			return true
		}
		return checkCA(string(cert), days)
	default:
		log.Println("[INFO] No search deploy type")
		return true
	}
}

func checkCA(cert string, days int) bool {
	certificates, err := certcrypto.ParsePEMBundle([]byte(cert))
	if err != nil {
		log.Fatalf("[ERROR] Open certificates failed. %s\n", err)
	}
	
	certificate := certificates[0]
	if certificate.IsCA {
		log.Fatalf("[ERROR] Certificate bundle starts with a CA certificate")
	}
	notAfter := int(time.Until(certificate.NotAfter).Hours() / 24.0)
	if notAfter > days {
		log.Printf("[INFO] %s The certificate expires in %d days, the number of days defined to perform the renewal is %d: no renewal.",
			certificate.DNSNames, notAfter, days)
		return false
	}
	return true
}
