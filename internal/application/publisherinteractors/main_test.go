package publisherinteractors_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	siteTestSetup()
	os.Exit(m.Run())
}