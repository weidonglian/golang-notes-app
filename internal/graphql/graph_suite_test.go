package graphql_test

import (
	"github.com/weidonglian/notes-app/internal/db"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGraph(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Graph Suite")
}

var _ = AfterSuite(func() {
	db.UnloadSessionPool()
})
