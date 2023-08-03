package publisher_test

import (
	"testing"

	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

func TestCreateNewPublisher(t *testing.T) {
	_ = publisher.NewPublisher(shared.Id("xxxx-yyyy"))

}