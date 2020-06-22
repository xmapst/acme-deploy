package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/certificate"
	"github.com/go-acme/lego/v3/challenge/dns01"
	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/providers/dns"
	"github.com/go-acme/lego/v3/registration"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func getEmail(ctx *cli.Context) string {
	email := ctx.GlobalString("email")
	if len(email) == 0 {
		log.Fatal("You have to pass an account (email address) to the program using --email or -m")
	}
	return email
}

// getKeyType the type from which private keys should be generated
func getKeyType(ctx *cli.Context) certcrypto.KeyType {
	keyType := ctx.GlobalString("key-type")
	switch strings.ToUpper(keyType) {
	case "RSA2048":
		return certcrypto.RSA2048
	case "RSA4096":
		return certcrypto.RSA4096
	case "RSA8192":
		return certcrypto.RSA8192
	case "EC256":
		return certcrypto.EC256
	case "EC384":
		return certcrypto.EC384
	}

	log.Fatalf("Unsupported KeyType: %s", keyType)
	return ""
}

func setup(ctx *cli.Context, accounts *Accounts) *lego.Client {
	config := lego.NewConfig(accounts)
	config.CADirURL = ctx.GlobalString("server")
	config.Certificate = lego.CertificateConfig{
		KeyType: getKeyType(ctx),
		Timeout: time.Duration(ctx.GlobalInt("cert.timeout")) * time.Second,
	}
	config.UserAgent = fmt.Sprintf("acme-cli/%s", ctx.App.Version)
	if ctx.GlobalIsSet("http-timeout") {
		config.HTTPClient.Timeout = time.Duration(ctx.GlobalInt("http-timeout")) * time.Second
	}
	
	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatalf("Could not create client: %v", err)
	}
	if client.GetExternalAccountRequired() && !ctx.GlobalIsSet("eab") {
		log.Fatal("Server requires External Account Binding. Use --eab with --kid and --hmac.")
	}
	return client
}

func setupChallenges(ctx *cli.Context, client *lego.Client) {
	if !ctx.GlobalIsSet("dns") {
		log.Fatal("No challenge selected. You must specify at least one challenge: `--dns`.")
	}
	provider, err := dns.NewDNSChallengeProviderByName(ctx.GlobalString("dns"))
	if err != nil {
		log.Fatal(err)
	}
	
	servers := ctx.GlobalStringSlice("dns.resolvers")
	err = client.Challenge.SetDNS01Provider(provider,
		dns01.CondOption(len(servers) > 0,
			dns01.AddRecursiveNameservers(dns01.ParseNameservers(ctx.GlobalStringSlice("dns.resolvers")))),
		dns01.CondOption(ctx.GlobalBool("dns.disable-cp"),
			dns01.DisableCompletePropagationRequirement()),
		dns01.CondOption(ctx.GlobalIsSet("dns-timeout"),
			dns01.AddDNSTimeout(time.Duration(ctx.GlobalInt("dns-timeout"))*time.Second)),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func register(ctx *cli.Context, client *lego.Client) (*registration.Resource, error) {
	if ctx.GlobalBool("eab") {
		kid := ctx.GlobalString("kid")
		hmacEncoded := ctx.GlobalString("hmac")
		
		if kid == "" || hmacEncoded == "" {
			log.Fatalf("Requires arguments --kid and --hmac.")
		}
		
		return client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
			TermsOfServiceAgreed: true,
			Kid:                  kid,
			HmacEncoded:          hmacEncoded,
		})
	}
	return client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
}

func obtainCertificate(ctx *cli.Context, client *lego.Client) (*certificate.Resource, error) {
	bundle := !ctx.Bool("no-bundle")
	domains := ctx.GlobalStringSlice("domains")
	if len(domains) > 0 {
		// obtain a certificate, generating a new private key
		request := certificate.ObtainRequest{
			Domains:    domains,
			Bundle:     bundle,
			MustStaple: ctx.Bool("must-staple"),
		}
		return client.Certificate.Obtain(request)
	}
	
	// read the CSR
	csr, err := readCSRFile(ctx.GlobalString("csr"))
	if err != nil {
		return nil, err
	}
	
	// obtain a certificate for this CSR
	return client.Certificate.ObtainForCSR(*csr, bundle)
}

func readCSRFile(filename string) (*x509.CertificateRequest, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	raw := bytes
	
	// see if we can find a PEM-encoded CSR
	var p *pem.Block
	rest := bytes
	for {
		// decode a PEM block
		p, rest = pem.Decode(rest)
		
		// did we fail?
		if p == nil {
			break
		}
		
		// did we get a CSR?
		if p.Type == "CERTIFICATE REQUEST" {
			raw = p.Bytes
		}
	}
	
	// no PEM-encoded CSR
	// assume we were given a DER-encoded ASN.1 CSR
	// (if this assumption is wrong, parsing these bytes will fail)
	return x509.ParseCertificateRequest(raw)
}