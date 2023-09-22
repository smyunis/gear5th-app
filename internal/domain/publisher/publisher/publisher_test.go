package publisher_test

import (
	"testing"


	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

func TestCreateNewPublisher(t *testing.T) {
	_ = publisher.NewPublisher(shared.ID("xxxx-yyyy"))

}

