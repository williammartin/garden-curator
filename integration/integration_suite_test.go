package integration

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var curatorBin string

var _ = BeforeSuite(func() {
	var err error
	curatorBin, err = gexec.Build("github.com/williammartin/garden-curator")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func execCurator(cwd string, args ...string) *gexec.Session {
	cmd := exec.Command(curatorBin, args...)
	cmd.Dir = cwd
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}

func TestGardenCurator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GardenCurator Suite")
}
