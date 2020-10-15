package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/db"
	"github.com/weidonglian/golang-notes-app/logging"
	"testing"
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store Suite")
}

var (
	dbSessionPool *db.SessionPool
)

var _ = BeforeSuite(func() {
	config.SetTestMode()
	dbSessionPool = db.LoadSessionPool(logging.NewLogger(), config.GetConfig())
})

var _ = AfterSuite(func() {
	if dbSessionPool != nil {
		db.UnloadSessionPool()
		dbSessionPool = nil
	}
})
