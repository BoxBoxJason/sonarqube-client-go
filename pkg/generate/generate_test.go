package generate_test

import (
	. "github.com/boxboxjason/sonarqube-client-go/pkg/generate"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetCurrentModulePath", func() {
	Describe("Test module path retrieval", func() {
		It("Should read module path from go.mod", func() {
			modulePath, err := GetCurrentModulePath()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(modulePath).To(Equal("github.com/boxboxjason/sonarqube-client-go"))
		})
	})
})
