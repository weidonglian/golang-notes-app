package handlers_test

import (
	"github.com/weidonglian/notes-app/internal/db"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = AfterSuite(func() {
	db.UnloadSessionPool()
})
