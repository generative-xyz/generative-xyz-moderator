package contractaddress

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestContractAddress_EncryptPrivateKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatal(err)
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)
	privateKeyStr := string(privateKeyPem)
	// t.Log(privateKeyStr)

	// The public key is a part of the *rsa.privateKey struct
	publicKey := privateKey.PublicKey
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&publicKey)
	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)
	publicKeyStr := string(publicKeyPem)
	// t.Log(publicKeyStr)

	type fields struct {
		Address             string
		PrivateKey          string
		PrivateKeyEncrypted string
	}
	type args struct {
		publicKey  string
		privateKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				Address:    "0x0",
				PrivateKey: "319717766d44c592d9971d8c595c95caea59017755933deb850b37f0b7941dc4",
				// PrivateKey: "vh5DT8RYMy65kJL5Eq7aHZJfLa6RUjqKeYWsdWqQd5AqFziUzi2HtQedZxYUbt8ARPkkbrWuWJR4db1TfJ7Xni7",
			},
			wantErr: false,
			args: args{
				publicKey:  publicKeyStr,
				privateKey: privateKeyStr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewContractAddress().WithAddress(tt.fields.Address).
				WithPrivateKey(tt.fields.PrivateKey)
			if err := a.EncryptPrivateKey(tt.args.publicKey); (err != nil) != tt.wantErr {
				t.Errorf("EncryptPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := a.DecryptPrivateKey(tt.args.privateKey); (err != nil) != tt.wantErr {
				t.Errorf("DecryptPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
			}

			if a.privateKey != tt.fields.PrivateKey {
				t.Errorf("DecryptPrivateKey() error = %v, wantErr %v", a.privateKey, tt.fields.PrivateKey)
			}
		})
	}
}
