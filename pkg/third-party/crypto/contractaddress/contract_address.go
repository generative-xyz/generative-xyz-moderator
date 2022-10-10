package contractaddress

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

type ContractAddress struct {
	// walletAddress is the address of the token contract, is also public key of the token contract.
	walletAddress string

	// privateKey is the private key of the token contract.
	privateKey string

	// privateKeyEncrypted is the encrypted private key of the token contract.
	privateKeyEncrypted string
}

// WithAddress set address
func (a *ContractAddress) WithAddress(address string) *ContractAddress {
	a.walletAddress = address
	return a
}

// WithPrivateKey set private key
func (a *ContractAddress) WithPrivateKey(privateKey string) *ContractAddress {
	a.privateKey = privateKey
	return a
}

// WithPrivateKeyEncrypted set encrypted private key
func (a *ContractAddress) WithPrivateKeyEncrypted(privateKeyEncrypted string) *ContractAddress {
	a.privateKeyEncrypted = privateKeyEncrypted
	return a
}

// GetWalletAddress get wallet address
func (a *ContractAddress) GetWalletAddress() string {
	return a.walletAddress
}

// GetPrivateKey get private key
func (a *ContractAddress) GetPrivateKey() string {
	return a.privateKey
}

// GetPrivateKeyEncrypted get encrypted private key
func (a *ContractAddress) GetPrivateKeyEncrypted() string {
	return a.privateKeyEncrypted
}

// EncryptPrivateKey encrypt private key
func (a *ContractAddress) EncryptPrivateKey(publicKey string) error {
	block, _ := pem.Decode([]byte(publicKey))
	rsaPublicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return err
	}

	encryptedBytes, err := rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		rsaPublicKey,
		[]byte(a.privateKey),
		nil)
	if err != nil {
		return err
	}

	a.privateKeyEncrypted = base64.StdEncoding.EncodeToString(encryptedBytes)
	a.privateKey = ""
	return nil
}

// DecryptPrivateKey decrypt private key
func (a *ContractAddress) DecryptPrivateKey(privateKey string) error {
	block, _ := pem.Decode([]byte(privateKey))
	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(a.privateKeyEncrypted)
	if err != nil {
		return err
	}
	decryptedBytes, err := rsa.DecryptOAEP(
		sha512.New(),
		rand.Reader,
		rsaPrivateKey,
		encryptedBytes,
		nil)
	if err != nil {
		return err
	}

	a.privateKey = string(decryptedBytes)
	return nil
}

// NewContractAddress create new empty contract address
func NewContractAddress() *ContractAddress {
	return &ContractAddress{}
}
