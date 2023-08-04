package siteverification

import (
	"os"
	"testing"
)


func TestCanCheckIfRecordIsInAdsTxt(t *testing.T) {
	buf, err := os.ReadFile("testdata/tutspointads.txt")
	if err != nil {
		t.Fatal(err.Error())
	}
	testAdsTxt := string(buf)
	record := AdsTxtRecord{
		AdExchangeDomain: "consumable.com",
		PublisherId:      "2000970",
		Relation:         "DIRECT",
	}

	if !hasAdsTxtRecord(testAdsTxt, record) {
		t.FailNow()
	}

}

