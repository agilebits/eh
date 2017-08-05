package secrets

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/agilebits/urlreader"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/pkg/errors"
)

const (
	typeLocal  = "local"
	typeAWSKMS = "awskms"
)

// Encrypt will generate a new key and encrypt the protected values.
func Encrypt(contents []byte) ([]byte, error) {
	tree, err := hcl.ParseBytes(contents)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Parse")
	}

	var wrapper Wrapper
	if err := hcl.DecodeObject(&wrapper, tree); err != nil {
		return nil, errors.Wrap(err, "failed to DecodeObject")
	}

	if wrapper.Header.Encrypted {
		return nil, errors.New("contents is already encrypted")
	}

	keyService, err := getKeyService(wrapper.Header.Service)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to obtain key service for parameters: %v", wrapper.Header.Service)
	}

	kid := "sm-" + time.Now().Format(time.RFC3339)
	encryptionKey, err := keyService.GenerateKey(kid)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate encryption key")
	}

	protect := make(map[string]bool)
	for _, s := range wrapper.Header.Protect {
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

// decryptWithHeader will access the key service and decrypt the protected values in the content. It returns unformatted AST file and 'eh' header found in the contents.
func decryptWithHeader(contents []byte, failIfNotEncrypted bool) (*ast.File, *Header, error) {
	tree, err := hcl.ParseBytes(contents)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to ParseBytes")
	}

	var wrapper Wrapper
	if err := hcl.DecodeObject(&wrapper, tree); err != nil {
		return nil, nil, errors.Wrap(err, "failed to DecodeObject")
	}

	if !wrapper.Header.Encrypted {
		if failIfNotEncrypted {
			return nil, nil, errors.New("contents is not encrypted")
		}

		return tree, &wrapper.Header, nil
	}

	keyBytes, err := base64.RawURLEncoding.DecodeString(wrapper.Header.Key)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to decode the encryption key")
	}

	encryptionKey := EncryptionKey{}
	if err := json.Unmarshal(keyBytes, &encryptionKey); err != nil {
		return nil, nil, errors.Wrap(err, "failed to Unmarshal the encryption key")
	}

	keyService, err := getKeyService(wrapper.Header.Service)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to obtain key service for parameters: %v", wrapper.Header.Service)
	}

	if err := keyService.DecryptKey(&encryptionKey); err != nil {
		return nil, nil, errors.Wrap(err, "failed to obtain decrypt key")
	}

	protect := make(map[string]bool)
	for _, s := range wrapper.Header.Protect {
		protect[s] = true
	}

	if err := processNode("", tree, opDecrypt, &encryptionKey, protect); err != nil {
		return nil, nil, errors.Wrap(err, "failed to process")
	}

	if err := removeEncryptionKey(tree); err != nil {
		return nil, nil, errors.Wrap(err, "failed to removeEncryptionKey")
	}

	return tree, &wrapper.Header, nil
}

// FormatASTFile returns formatted text representation of the file
func FormatASTFile(file *ast.File) ([]byte, error) {
	var c printer.Config
	var result bytes.Buffer
	if err := c.Fprint(&result, file); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

// Decrypt will access the key service and decrypt the protected values in the content.
func Decrypt(contents []byte) ([]byte, error) {
	result, _, err := decryptWithHeader(contents, true)
	if err != nil {
		return nil, err
	}

	return FormatASTFile(result)
}

// Read loads and decrypt the contents at the specifed URL. It also processes and merges all included files specified in the header.
func Read(url string) ([]byte, error) {
	reader, err := urlreader.Open(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open url %q", url)
	}

	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read from url %q", url)
	}

	if err := reader.Close(); err != nil {
		return nil, err
	}

	tree, header, err := decryptWithHeader(contents, false)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decrypt url %q", url)
	}

	removeHeader(tree)

	text, err := FormatASTFile(tree)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to format contents of url %q", url)
	}

	var result bytes.Buffer
	result.Write(text)

	for _, name := range header.Include {
		if strings.HasPrefix(name, "./") {
			dir, _ := path.Split(url)
			name = dir + "/" + name[2:]
		}

		fragment, err := Read(name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to include %q", name)
		}

		result.WriteString("\n\n// " + name + "\n")
		result.Write(fragment)
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
	case typeAWSKMS:
		keyService = NewAwsKeyService(service.Region, service.MasterKey)
	default:
		return nil, fmt.Errorf("unsupported service type: %+q", service.Type)
	}

	return keyService, nil
}
