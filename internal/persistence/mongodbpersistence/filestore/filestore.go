package filestore

import (
	"bytes"
	"io"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type MongoDBGridFSFileStore struct {
	dbStore mongodbpersistence.MongoDBStore
	bucket  *gridfs.Bucket
}

func NewMongoDBGridFSFileStore(dbStore mongodbpersistence.MongoDBStore) MongoDBGridFSFileStore {
	bucket, err := gridfs.NewBucket(dbStore.Database())
	if err != nil {
		panic(err)
	}
	return MongoDBGridFSFileStore{
		dbStore: dbStore,
		bucket:  bucket,
	}
}

func (f MongoDBGridFSFileStore) Get(fileID string) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	_, err := f.bucket.DownloadToStreamByName(fileID, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (f MongoDBGridFSFileStore) Save(file io.Reader) (string, error) {
	fileID := shared.NewID().String()
	_, err := f.bucket.UploadFromStream(fileID, file)
	if err != nil {
		return "", err
	}
	return fileID, nil
}
