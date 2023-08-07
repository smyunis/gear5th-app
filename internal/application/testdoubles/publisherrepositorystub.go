package testdoubles

import (
	"reflect"
	"unsafe"

	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

type PublisherRepositoryStub struct{}

func (p PublisherRepositoryStub) Get(id shared.Id) (publisher.Publisher, error) {

	if id != shared.Id("stub-id-xxx") {
		return publisher.Publisher{}, shared.NewEntityNotFoundError(id.String(), "publisher")
	}

	pub := &publisher.Publisher{}
	uVal := reflect.ValueOf(pub).Elem()

	idField := uVal.FieldByName("publisherUserId")
	idField = reflect.NewAt(idField.Type(), unsafe.Pointer(idField.UnsafeAddr())).Elem()
	idField.Set(reflect.ValueOf(id))

	return *pub, nil
}

func (p PublisherRepositoryStub) Save(pub publisher.Publisher) error {
	return nil
}
