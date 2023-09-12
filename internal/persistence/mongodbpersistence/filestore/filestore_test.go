//go:build db

package filestore_test

import (
	"io"
	"os"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/filestore"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/mongotestdoubles"
)

func TestMain(m *testing.M) {
	setup()

	os.Exit(m.Run())
}

var store application.FileStore

func setup() {
	configProvider := mongotestdoubles.NewTestEnvConfigurationProvider()
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	store = filestore.NewMongoDBGridFSFileStore(dbStore)
}

func TestSaveFile(t *testing.T) {
	testImg, err := os.Open("testdata/luffy.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer testImg.Close()

	id, err := store.Save(testImg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
}

func TestDownloadSavedFile(t *testing.T) {
	testImg, err := os.Open("testdata/luffy.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer testImg.Close()

	id, err := store.Save(testImg)
	if err != nil {
		t.Fatal(err)
	}

	file, err := store.Get(id)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := io.ReadAll(file)
	os.WriteFile("testdata/luffy_dn.jpg", b, 0644)

}
