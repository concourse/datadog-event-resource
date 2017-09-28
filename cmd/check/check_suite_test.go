package main_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"
	"testing"

	"encoding/json"

	"bytes"

	"github.com/concourse/datadog-event-resource/cmd"
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

	if _, err = os.Stat("/opt/resource/check"); err == nil {
		binPath = "/opt/resource/check"
	} else {
		binPath, err = gexec.Build("github.com/concourse/datadog-event-resource/cmd/check")
		Expect(err).NotTo(HaveOccurred())
	}
})

var _ = AfterEach(func() {
	fakeDataDogServer.Close()
})

func RunCheckSuccessfully(id *string, filter string) *gexec.Session {
	sess := RunCheck(id, filter)
	Expect(sess).To(gexec.Exit(0))
	return sess
}

func RunCheck(id *string, filter string) *gexec.Session {
	var version *cmd.Version
	if id != nil {
		version = &cmd.Version{
			Id: *id,
		}
	}

	payload := cmd.CheckPayload{
		Source: cmd.Source{
			ApplicationKey: applicationKey,
			ApiKey:         apiKey,
			Filter:         filter,
		},
		Version: version,
	}

	b, err := json.Marshal(&payload)
	Expect(err).NotTo(HaveOccurred())

	c := exec.Command(binPath)
	c.Stdin = bytes.NewBuffer(b)
	c.Env = append(c.Env, "DATADOG_HOST=http://"+fakeDataDogServer.Addr())
	sess, err := gexec.Start(c, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	<-sess.Exited
	return sess
}
