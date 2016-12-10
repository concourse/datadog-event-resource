package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"
	"testing"

	"encoding/json"

	"github.com/concourse/datadog-resource/cmd"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

func TestCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Check Suite")
}

var (
	binPath           string
	err               error
	fakeDataDogServer *ghttp.Server

	applicationKey string
	apiKey         string
)

var _ = BeforeEach(func() {
	applicationKey = "some-application-key"
	apiKey = "some-api-key"

	fakeDataDogServer = ghttp.NewServer()

	binPath, err = gexec.Build("github.com/concourse/datadog-resource/cmd/check")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterEach(func() {
	fakeDataDogServer.Close()
})

func RunCheck(id *string) *gexec.Session {
	var version *cmd.Version
	if id != nil {
		version = &cmd.Version{
			Id: *id,
		}
	}

	payload := cmd.CheckPayload{
		Source: cmd.Source{
			ApplicationKey: "foobar",
			ApiKey:         "barbaz",
		},
		Version: version,
	}

	b, err := json.Marshal(&payload)
	Expect(err).NotTo(HaveOccurred())

	c := exec.Command(binPath, string(b))
	c.Env = append(c.Env, "DATADOG_HOST=http://"+fakeDataDogServer.Addr())
	sess, err := gexec.Start(c, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	<-sess.Exited
	Expect(sess).To(gexec.Exit(0))
	return sess
}
