package integration

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/williammartin/garden-curator/blueprint"
	yaml "gopkg.in/yaml.v2"

	"code.cloudfoundry.org/garden"
	gclient "code.cloudfoundry.org/garden/client"
	gconn "code.cloudfoundry.org/garden/client/connection"
)

var _ = Describe("Growing", func() {

	It("can create a single container from a blueprint", func() {
		tempDir, err := ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		// Create a blueprint with the right contents
		blueprint := &blueprint.Blueprint{Containers: []string{"interview"}}
		bytes, err := yaml.Marshal(blueprint)
		Expect(err).NotTo(HaveOccurred())

		Expect(ioutil.WriteFile(filepath.Join(tempDir, "blueprint.yml"), bytes, 0755)).To(Succeed())

		// Run garden-curator grow
		session := execCurator(tempDir, "grow")
		session.Wait()

		// Check whether a container exists via the garden client library
		client := gclient.New(gconn.New("tcp", "10.244.0.2:7777"))
		containers, err := client.Containers(garden.Properties{})
		Expect(err).NotTo(HaveOccurred())

		Expect(containers).To(HaveLen(1))
		Expect(containers[0].Handle()).To(Equal("interview"))

		Expect(os.RemoveAll(tempDir)).To(Succeed())

		Expect(client.Destroy("interview")).To(Succeed())
	})
})
