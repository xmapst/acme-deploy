package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"time"
	//"github.com/go-acme/lego/v3/certificate"
)

type resource struct {
	Domain            string `json:"domain"`
	CertURL           string `json:"cert_url"`
	CertStableURL     string `json:"cert_stable_url"`
	PrivateKey        string `json:"private_key"`
	Certificate       string `json:"certificate"`
	IssuerCertificate string `json:"issuer_certificate"`
	CSR               string `json:"csr"`
	Timestamp         int64  `json:"timestamp"`
}

func run(ctx *cli.Context) error {
	log.Println("[INFO] Check certificate expiration time")
	if !checkExpirationTime(ctx) {
		log.Println("[INFO] The certificate has not expired, No need to update")
		return nil
	}
	
	log.Println("[INFO] New accounts need an email and private key to start")
	accounts := NewAccounts(ctx)

	log.Println("[INFO] Client facilitates communication with the CA server")
	client := setup(ctx, &accounts)

	log.Println("[INFO] setup challenge")
	setupChallenges(ctx, client)

	log.Println("[INFO] New users will need to register")
	reg, err := register(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	accounts.Registration = reg

	log.Println("[INFO] Start issuing certificates")
	cert, err := obtainCertificate(ctx, client)
	if err != nil {
		// Make sure to return a non-zero exit code if ObtainSANCertificate returned at least one error.
		// Due to us not returning partial certificate we can just exit here instead of at the end.
		log.Fatalf("Could not obtain certificates:\n\t%v", err)
	}

	// save certificates to {domain}.json
	log.Printf("[INFO] Save certificates to %s.json.\n", cert.Domain)
	jsonBytes, err := json.Marshal(resource{
		Domain:            cert.Domain,
		CertURL:           cert.CertURL,
		CertStableURL:     cert.CertStableURL,
		PrivateKey:        string(cert.PrivateKey),
		Certificate:       string(cert.Certificate),
		IssuerCertificate: string(cert.IssuerCertificate),
		CSR:               string(cert.CSR),
		Timestamp:         time.Now().UTC().Unix(),
	})
	if err != nil {
		log.Fatalf("Unable to marshal CertResource for domain %s\n\t%v", cert.Domain, err)
	}
	if err := writeFile(fmt.Sprintf("%s.json", cert.Domain), jsonBytes); err != nil {
		log.Fatalf("[ERROR] %s.\n", err.Error())
	}

	log.Println("[INFO] Certificate issued successfully")

	if !ctx.GlobalIsSet("deploy") {
		log.Println("[INFO] No deployment selected. You can specify the deployment as needed: --deploy.")
		log.Printf("[INFO] Write to local file, %s.key, %s.cert /n", cert.Domain, cert.Domain)

		if err := writeFile(fmt.Sprintf("%s.key", cert.Domain), cert.PrivateKey); err != nil {
			log.Fatalf("[ERROR] %s.\n", err.Error())
		}
		if err := writeFile(fmt.Sprintf("%s.cert", cert.Domain), cert.Certificate); err != nil {
			log.Fatalf("[ERROR] %s.\n", err.Error())
		}
		return nil
	}
	if err := deploy(ctx, cert); err != nil {
		log.Fatalf("[ERROR] %s", err.Error())
	}
	log.Printf("[INFO] Successfully deployed to %s.\n", ctx.GlobalString("deploy"))
	return nil
}

func writeFile(name string, content []byte) error {
	if err := ioutil.WriteFile(name, content, 0600); err != nil {
		return fmt.Errorf("[ERROR] %s.\n", err.Error())
	}
	return nil
}
