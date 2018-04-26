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

type CuratorRunConfig struct {
	RunDir string
	Args   []string
}

var _ = Describe("Growing", func() {

	var (
		client garden.Client

		tempDir          string
		curatorRunConfig *CuratorRunConfig
		stdout           string

		blueprint *bprint.Blueprint
	)

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		curatorRunConfig = &CuratorRunConfig{
			RunDir: tempDir,
			Args:   []string{"grow"},
		}

		client = gclient.New(gconn.New("tcp", "10.244.0.2:7777"))
	})

	JustBeforeEach(func() {
		bytes, err := yaml.Marshal(blueprint)
		Expect(err).NotTo(HaveOccurred())
		Expect(ioutil.WriteFile(filepath.Join(tempDir, "blueprint.yml"), bytes, 0755)).To(Succeed())

		session := execCurator(curatorRunConfig)
		session.Wait()
		stdout = string(session.Out.Contents())
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

		It("logs feedback that each container is being created", func() {
			Expect(stdout).To(ContainSubstring("growing 'base'..."))
			Expect(stdout).To(ContainSubstring("growing 'ace'..."))
			Expect(stdout).To(ContainSubstring("growing 'keith'..."))
		})
	})

	Context("when the blueprint location is passed as an argument", func() {
		BeforeEach(func() {
			curatorRunConfig.RunDir = ""
			curatorRunConfig.Args = append(curatorRunConfig.Args, "-b", filepath.Join(tempDir, "blueprint.yml"))
			blueprint = &bprint.Blueprint{Containers: []string{"distant-blueprint"}}
		})

		It("uses that file as the blueprint", func() {
			containers, err := client.Containers(garden.Properties{})

			Expect(err).NotTo(HaveOccurred())
			Expect(containers).To(HaveLen(1))
			Expect(containers[0].Handle()).To(Equal("distant-blueprint"))
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
