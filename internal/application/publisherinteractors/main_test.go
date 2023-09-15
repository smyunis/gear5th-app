package publisherinteractors_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	publisherTestSetup()
	siteTestSetup()
	os.Exit(m.Run())
}