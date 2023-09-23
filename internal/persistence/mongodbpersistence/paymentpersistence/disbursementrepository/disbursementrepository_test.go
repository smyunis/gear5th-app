//go:build db
package disbursementrepository_test

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/disbursement"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/mongotestdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/paymentpersistence/disbursementrepository"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var disbuRepo disbursement.DisbursementRepository

func setup() {
	configProvider := mongotestdoubles.NewTestEnvConfigurationProvider()
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	logger := testdoubles.ConsoleLogger{}
	disbuRepo = disbursementrepository.NewMongoDBDisbursementRepository(dbStore, logger)
}

func TestSaveDisbursement(t *testing.T) {
	ph, _ := user.NewPhoneNumber("0978564312")
	pf := disbursement.PaymentProfile{
		PaymentMethod: disbursement.CommercialBankOfEthiopia,
		Account:       "1000456666678",
		FullName:      "Son Goku",
		PhoneNumber:   ph,
	}
	d := disbursement.NewDisbursement(shared.NewID(), pf, 4677, time.Now(), time.Now())

	err := disbuRepo.Save(context.TODO(), d)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSavedDisbursement(t *testing.T) {
	ph, _ := user.NewPhoneNumber("0978564312")
	pf := disbursement.PaymentProfile{
		PaymentMethod: disbursement.CommercialBankOfEthiopia,
		Account:       "1000456666678",
		FullName:      "Son Goku",
		PhoneNumber:   ph,
	}
	d := disbursement.NewDisbursement(shared.NewID(), pf, 4677, time.Now(), time.Now())

	err := disbuRepo.Save(context.TODO(), d)
	if err != nil {
		t.Fatal(err)
	}

	fd, err := disbuRepo.Get(context.TODO(), d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fd.PaymentProfile, d.PaymentProfile) {
		t.FailNow()
	}
}

func TestPublisherDisbursements(t *testing.T) {
	ph, _ := user.NewPhoneNumber("0978564312")
	pf := disbursement.PaymentProfile{
		PaymentMethod: disbursement.CommercialBankOfEthiopia,
		Account:       "1000456666678",
		FullName:      "Son Goku",
		PhoneNumber:   ph,
	}
	d := disbursement.NewDisbursement(shared.NewID(), pf, 4677, time.Now(), time.Now())

	err := disbuRepo.Save(context.TODO(), d)
	if err != nil {
		t.Fatal(err)
	}

	ds, err := disbuRepo.DisbursementsForPublisher(d.PublisherID, disbursement.Requested)
	if err != nil {
		t.Fatal(err)
	}

	if len(ds) == 0 {
		t.FailNow()
	}
}
