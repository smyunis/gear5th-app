package testdoubles

import (
	"context"
	"reflect"
	"unsafe"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type PublisherRepositoryStub struct{}

func (p PublisherRepositoryStub) Get(ctx context.Context, id shared.ID) (publisher.Publisher, error) {

	if id != shared.ID("stub-id-xxx") {
		return publisher.Publisher{}, application.ErrEntityNotFound
	}

	pub := &publisher.Publisher{}
	uVal := reflect.ValueOf(pub).Elem()

	idField := uVal.FieldByName("publisherUserId")
	idField = reflect.NewAt(idField.Type(), unsafe.Pointer(idField.UnsafeAddr())).Elem()
	idField.Set(reflect.ValueOf(id))

	return *pub, nil
}

func (p PublisherRepositoryStub) Save(ctx context.Context, pub publisher.Publisher) error {
	return nil
}
