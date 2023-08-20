package tokens

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type HS256HMACValidationService struct {
	key       []byte
	separator string
}

func NewHS256HMACValidationService() HS256HMACValidationService {
	return HS256HMACValidationService{
		key: []byte("F1WjYXPWgir9hU5HGm3YDZhwgzApGBcp"),
		separator: "+",
	}
}

func (s HS256HMACValidationService) Generate(message string) (string, error) {
	h := hmac.New(sha256.New, s.key)
	_, err := h.Write([]byte(message))
	if err != nil {
		return "", nil
	}
	generated := h.Sum(nil)

	generatedHMACStr := hex.EncodeToString(generated)
	return strings.Join([]string{message, generatedHMACStr}, s.separator), nil
}

func (s HS256HMACValidationService) Validate(hashedMessage string) bool {
	messageHash := strings.Split(hashedMessage, s.separator) //[message,hash]

	if len(messageHash) != 2 {
		return false
	}

	h := hmac.New(sha256.New, s.key)
	_, err := h.Write([]byte(messageHash[0]))
	if err != nil {
		return false
	}
	msgHMAC := h.Sum(nil)

	decodedHash, err := hex.DecodeString(messageHash[1])
	if err != nil {
		return false
	}

	return hmac.Equal(msgHMAC, decodedHash)
}
