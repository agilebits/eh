package secrets

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/pkg/errors"
)

const (
	typeLocal = "local"
	typeAWS   = "aws"
)

// Encrypt will generate a new key and encrypt the protected values.
func Encrypt(contents []byte) ([]byte, error) {
	tree, err := hcl.ParseBytes(contents)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Parse")
	}

	var header Header
	if err := hcl.DecodeObject(&header, tree); err != nil {
		return nil, errors.Wrap(err, "failed to DecodeObject")
	}

	if header.EHCL.Encrypted {
		return nil, errors.New("contents is already encrypted")
	}

	keyService, err := getKeyService(header.EHCL.Service)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to obtain key service for parameters: %v", header.EHCL.Service)
	}

	kid := "sm-" + time.Now().Format(time.RFC3339)
	encryptionKey, err := keyService.GenerateKey(kid)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate encryption key")
	}

	protect := make(map[string]bool)
	for _, s := range header.EHCL.Protect {
		protect[s] = true
	}

	if err := processNode("", tree, opEncrypt, encryptionKey, protect); err != nil {
		return nil, errors.Wrap(err, "failed to process")
	}

	if err := addEncryptionKey(tree, encryptionKey); err != nil {
		return nil, errors.Wrap(err, "failed to addEncryptionKey")
	}

	var c printer.Config
	var result bytes.Buffer
	if err := c.Fprint(&result, tree); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

// Decrypt will access the key service and decrypt the protected values in the content.
func Decrypt(contents []byte) ([]byte, error) {
	tree, err := hcl.ParseBytes(contents)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ParseBytes")
	}

	var header Header
	if err := hcl.DecodeObject(&header, tree); err != nil {
		return nil, errors.Wrap(err, "failed to DecodeObject")
	}

	if !header.EHCL.Encrypted {
		return nil, errors.New("contents is not encrypted")
	}

	keyBytes, err := base64.RawURLEncoding.DecodeString(header.EHCL.Key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode the encryption key")
	}

	encryptionKey := EncryptionKey{}
	if err := json.Unmarshal(keyBytes, &encryptionKey); err != nil {
		return nil, errors.Wrap(err, "failed to Unmarshal the encryption key")
	}

	keyService, err := getKeyService(header.EHCL.Service)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to obtain key service for parameters: %v", header.EHCL.Service)
	}

	if err := keyService.DecryptKey(&encryptionKey); err != nil {
		return nil, errors.Wrap(err, "failed to obtain decrypt key")
	}

	protect := make(map[string]bool)
	for _, s := range header.EHCL.Protect {
		protect[s] = true
	}

	if err := processNode("", tree, opDecrypt, &encryptionKey, protect); err != nil {
		return nil, errors.Wrap(err, "failed to process")
	}

	if err := removeEncryptionKey(tree); err != nil {
		return nil, errors.Wrap(err, "failed to removeEncryptionKey")
	}

	var c printer.Config
	var result bytes.Buffer
	if err := c.Fprint(&result, tree); err != nil {
		return nil, err
	}

	return result.Bytes(), nil

}

func getKeyService(service ServiceParams) (KeyService, error) {
	if service.Type == "" {
		return nil, errors.New("missing service type")
	}

	var keyService KeyService
	switch service.Type {
	case typeLocal:
		keyService = NewDevKeyService()
	case typeAWS:
		keyService = NewAwsKeyService(service.Region, service.MasterKey)
	default:
		return nil, fmt.Errorf("unsupported service type: %+q", service.Type)
	}

	return keyService, nil
}
