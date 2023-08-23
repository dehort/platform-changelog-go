package endpoints_test

import (
	"fmt"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	"github.com/redhatinsights/platform-changelog-go/internal/utils"
	"gorm.io/gorm"
)

var (
	testDB     *embeddedpostgres.EmbeddedPostgres
	testGormDB *gorm.DB
	testDBImpl *db.DBConnectorImpl
)

func TestEndpoints(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Endpoints Suite")
}

var _ = BeforeSuite(func() {
	cfg := config.Get()

	var err error
	testDB, testGormDB, err = utils.CreateTestDB(cfg, "file://../../migrations")
	if err != nil {
		testDB.Stop()
	}

	Expect(err).To(BeNil())

	testDBImpl = db.SetDBConnector(testGormDB)

	seed(cfg, testDBImpl)
})

var _ = AfterSuite(func() {
	err := testDB.Stop()
	Expect(err).To(BeNil())
	fmt.Println("TEST DB STOPPED")
})

func seed(cfg *config.Config, db *db.DBConnectorImpl) {
	// Add all services in the services config
	for key, service := range cfg.Services {
		_, err := db.CreateServiceTableEntry(key, service)
		Expect(err).To(BeNil())
	}
}
