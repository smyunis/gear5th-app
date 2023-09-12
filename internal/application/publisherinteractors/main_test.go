package publisherinteractors_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	adSlotTestSetup()
	publisherTestSetup()
	siteTestSetup()
	os.Exit(m.Run())
}