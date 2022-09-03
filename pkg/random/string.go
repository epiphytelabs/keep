package random

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func String(length int) (string, error) {
	data := make([]byte, 8192)

	if _, err := rand.Read(data); err != nil {
		return "", err
	}

	enc := base64.StdEncoding.EncodeToString(data)

	if len(enc) < length {
		return "", fmt.Errorf("could not generate random string of length %d", length)
	}

	return enc[0:length], nil
}
