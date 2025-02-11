package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func Sign(value interface{}, privateKey string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("%s", "Failed to decode PEM block containing the private key")
	}
	block, _ := pem.Decode(bytes)
	if block == nil {
		return "", fmt.Errorf("%s", "Failed to decode PEM block containing the private key")
	}
	priKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("%s", fmt.Sprintf("Failed to parse private key: %v", err))
	}
	rsaPrivateKey, _ := priKey.(*rsa.PrivateKey)

	val := reflect.ValueOf(value).Elem()
	var fields []string
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fields = append(fields, field.Name)
	}
	sort.Strings(fields)
	var result []string
	for _, fieldName := range fields {
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			if field.Name == "Signature" {
				continue
			}
			if field.Name == fieldName {
				fieldValue := val.Field(i)
				var value string
				if field.Type.Kind() == reflect.Map {
					if fieldValue.IsNil() {
						continue
					}
					bytes, err := json.Marshal(fieldValue.Interface())
					if err != nil {
						return "", err
					}
					value = string(bytes)
				} else {
					value = fmt.Sprintf("%v", fieldValue.Interface())
				}
				tag := field.Tag.Get("json")
				tagParts := strings.Split(tag, ",")
				if len(tagParts) > 0 {
					tag = tagParts[0]
				}
				result = append(result, fmt.Sprintf("%s=%s", tag, value))
				break
			}
		}
	}
	dataString := strings.Join(result, "&")
	hash := sha256.New()
	hash.Write([]byte(dataString))
	hashed := hash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}
