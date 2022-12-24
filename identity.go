package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/42LoCo42/go-zeolite"
)

type Identity struct {
	ident *zeolite.Identity
}

func (id *Identity) Generate() error {
	newIdentity, err := zeolite.NewIdentity()

	if err != nil {
		return err
	}

	id.ident = &newIdentity

	return nil
}

func (id *Identity) LoadIdentity(ident []byte) error {
	if len(ident) == 0 {
		return errors.New("Invalid identity file")
	}

	identity := IdentityFile{}
	err := json.Unmarshal(ident, &identity)

	if err != nil {
		return err
	}

	public, err := zeolite.Base64Dec(identity.Public)

	if err != nil {
		return errors.New("Unable to decode public key")
	}

	private, err := zeolite.Base64Dec(identity.Private)

	if err != nil {
		return errors.New("Unable to decode private key")
	}

	var newIdentity zeolite.Identity
	copy(newIdentity.Public[:], public)
	copy(newIdentity.Secret[:], private)

	id.ident = &newIdentity

	return nil
}

func (id *Identity) NewStream(conn io.ReadWriter, cb zeolite.TrustCB) (ret *zeolite.Stream, err error) {
	if id.ident == nil {
		return nil, errors.New("No identity set up yet")
	}

	return id.ident.NewStream(conn, cb)
}

func (id *Identity) PublicKey() string {
	if id.ident == nil {
		return ""
	}

	return zeolite.Base64Enc(id.ident.Public[:])
}

func (id *Identity) PrivateKey() string {
	if id.ident == nil {
		return ""
	}

	return zeolite.Base64Enc(id.ident.Secret[:])
}

func (id *Identity) Print() {
	if id.ident == nil {
		return
	}

	fmt.Printf("Public: %s\nPrivate: %s\n", id.PublicKey(), id.PrivateKey())
}

type IdentityFile struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

func SaveIdentity(path string, ident *Identity) error {
	identity := IdentityFile{
		Private: ident.PrivateKey(),
		Public:  ident.PublicKey(),
	}

	jsonString, err := json.Marshal(identity)

	if err != nil {
		return err
	}

	if err = WriteToFile(path, []byte(jsonString)); err != nil {
		return err
	}

	return nil
}
