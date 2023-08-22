package publisher_test

import (
	"reflect"
	"testing"

	"slices"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

func TestCreateNewPublisher(t *testing.T) {
	_ = publisher.NewPublisher(shared.ID("xxxx-yyyy"))

}

func TestNotifyPublisher(t *testing.T) {
	p := publisher.NewPublisher(shared.ID("xxxx-yyyy"))
	n := shared.NewNotification("message")

	p.Notify(n)

	if !slices.Contains(p.UnacknowledgedNotifications(), n) {
		t.FailNow()
	}
}

func TestAcknowledgeNotifications(t *testing.T) {
	p := publisher.NewPublisher(shared.ID("xxxx-yyyy"))
	n1 := shared.NewNotification("message 1")
	n2 := shared.NewNotification("message 2")

	p.Notify(n1)
	p.Notify(n2)

	p.AcknowledgeNotifications()

	if len(p.UnacknowledgedNotifications()) != 0 {
		t.FailNow()
	}

}

func TestReconsitutePublisher(t *testing.T) {
	id := shared.NewID()
	np := publisher.NewPublisher(id)
	np.Notify(shared.NewNotification("Alert!"))

	rp := publisher.ReconstitutePublisher(id, np.UnacknowledgedNotifications())

	if !reflect.DeepEqual(rp, np) {
		t.Fatal("reconstituted not equal")
	}
}
