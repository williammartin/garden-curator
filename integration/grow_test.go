package integration

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	bprint "github.com/williammartin/garden-curator/blueprint"
	yaml "gopkg.in/yaml.v2"

	"code.cloudfoundry.org/garden"
	gclient "code.cloudfoundry.org/garden/client"
	gconn "code.cloudfoundry.org/garden/client/connection"
)

var _ = Describe("Growing", func() {

	var (
		tempDir string
		client  garden.Client

		blueprint *bprint.Blueprint
	)

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		client = gclient.New(gconn.New("tcp", "10.244.0.2:7777"))
	})

	JustBeforeEach(func() {
		bytes, err := yaml.Marshal(blueprint)
		Expect(err).NotTo(HaveOccurred())
		Expect(ioutil.WriteFile(filepath.Join(tempDir, "blueprint.yml"), bytes, 0755)).To(Succeed())

		session := execCurator(tempDir, "grow")
		session.Wait()
	})

	AfterEach(func() {
		destroyAllContainers(client)
		Expect(os.RemoveAll(tempDir)).To(Succeed())
	})

	Context("when the blueprint contains a single container", func() {
		BeforeEach(func() {
			blueprint = &bprint.Blueprint{Containers: []string{"interview"}}
		})

		It("grows that container", func() {
			containers, err := client.Containers(garden.Properties{})

			Expect(err).NotTo(HaveOccurred())
			Expect(containers).To(HaveLen(1))
			Expect(containers[0].Handle()).To(Equal("interview"))
		})
	})

	Context("when the blueprint contains multiple containers", func() {
		BeforeEach(func() {
			blueprint = &bprint.Blueprint{Containers: []string{"base", "ace", "keith"}}
		})

		It("grows all the containers", func() {
			containers, err := client.Containers(garden.Properties{})

			Expect(err).NotTo(HaveOccurred())
			Expect(containers).To(HaveLen(3))
			Expect(containersToHandles(containers)).To(ConsistOf("base", "ace", "keith"))
		})
	})
})

func destroyAllContainers(client garden.Client) {
	containers, err := client.Containers(garden.Properties{})
	Expect(err).NotTo(HaveOccurred())

	for _, container := range containers {
		Expect(client.Destroy(container.Handle())).To(Succeed())
	}
}

func containersToHandles(containers []garden.Container) []string {
	handles := []string{}
	for _, container := range containers {
		handles = append(handles, container.Handle())
	}
	return handles
}
