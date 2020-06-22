package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/go-acme/lego/v3/registration"
	"github.com/urfave/cli"
	"log"
)

type Accounts struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *Accounts) GetEmail() string {
	return u.Email
}
func (u Accounts) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *Accounts) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func NewAccounts(ctx *cli.Context) Accounts {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	return Accounts{
		Email: getEmail(ctx),
		key:   privateKey,
	}
}
