package handlers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {
	var (
		testApp = NewTestAppAndServe()
	)

	It("POST /session", func() {
		testApp.App.Serve()
		Expect(1).To(BeZero())
	})
})
