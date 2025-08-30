package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateRequestID creates a readable request ID with format: req_<random_string>.
// This must be used in each controller to generate a unique request ID for logging and tracing
func GenerateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	randomString := base64.URLEncoding.EncodeToString(b)[:10]
	return fmt.Sprintf("req_%s", randomString)
}