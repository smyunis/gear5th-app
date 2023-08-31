package adslot_test

import (
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

func TestChangeAdSlotName(t *testing.T) {
	siteID := shared.NewID()
	slot := adslot.NewAdSlot(siteID, "name-one", adslot.Box)
	slot.SetName("name-two")

	if slot.Name() != "name-two" {
		t.FailNow()
	}
}

