package tokens_test

import (
	"os"
	"testing"

	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure/identity/tokens"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var hservice identityinteractors.DigitalSignatureValidationService

func setup() {
	hservice = tokens.NewHS256HMACValidationService()
}

func TestCanGenerateSignature(t *testing.T) {
	msg := "mymesagethatIwanttosign"
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
