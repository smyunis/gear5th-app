package tokens_test

import (
	"os"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/identity/tokens"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var hservice application.DigitalSignatureService

func setup() {
	hservice = tokens.NewHS256HMACValidationService()
}

func TestCanGenerateSignature(t *testing.T) {
	msg := "01HA4GNAXY6DKN506Z4HXM8DTY"
	sign, err := hservice.Generate(msg)
	if err != nil {
		t.FailNow()
	}
	if sign == "" {
		t.FailNow()
	}
	t.Log(sign)
}

func TestCanVerifyGeneratedSigniture(t *testing.T) {
	msg := "mymesagethatIwanttosign"
	sign, err := hservice.Generate(msg)
	if err != nil {
		t.FailNow()
	}
	if sign == "" {
		t.FailNow()
	}

	if !hservice.Validate(sign) {
		t.FailNow()
	}
}

func TestCanGetMessageFromGeneratedHash(t *testing.T) {
	msg := "his.awsomemail@gmail.com"
	hash, err := hservice.Generate(msg)
	if err != nil {
		t.FailNow()
	}

	m, err := hservice.GetMessage(hash)
	if err != nil {
		t.FailNow()
	}
	if m != msg {
		t.FailNow()
	}
}
