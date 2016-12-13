package main_test

import (
	"encoding/json"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"io/ioutil"
	"testing"

	"github.com/concourse/datadog-resource/cmd"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

func TestOut(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Out Suite")
}

var (
	binPath           string
	err               error
	fakeDataDogServer *ghttp.Server

	applicationKey string
	apiKey         string

	tmpDir string
)

var _ = BeforeEach(func() {
	applicationKey = "some-application-key"
	apiKey = "some-api-key"

	fakeDataDogServer = ghttp.NewServer()

	binPath, err = gexec.Build("github.com/concourse/datadog-resource/cmd/out")
	Expect(err).NotTo(HaveOccurred())

	tmpDir, err = ioutil.TempDir("", "datadog_event_resource_out")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterEach(func() {
	fakeDataDogServer.Close()
})

func RunOut(p cmd.OutParams) *gexec.Session {
	payload := cmd.OutPayload{
		Source: cmd.Source{
			ApplicationKey: applicationKey,
			ApiKey:         apiKey,
		},
		Params: p,
	}

	b, err := json.Marshal(&payload)
	Expect(err).NotTo(HaveOccurred())

	c := exec.Command(binPath, tmpDir)
	c.Stdin = bytes.NewBuffer(b)
	c.Env = append(c.Env, "DATADOG_HOST=http://"+fakeDataDogServer.Addr())
	sess, err := gexec.Start(c, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	<-sess.Exited
	Expect(sess).To(gexec.Exit(0))
	return sess
}
