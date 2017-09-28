package main_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"
	"testing"

	"encoding/json"

	"bytes"

	"io/ioutil"

	"strconv"

	"github.com/concourse/datadog-event-resource/cmd"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

func TestIn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "In Suite")
}

var (
	binPath           string
	err               error
	fakeDataDogServer *ghttp.Server

	applicationKey string
	apiKey         string

	id     int
	tmpDir string
)

var _ = BeforeEach(func() {
	applicationKey = "some-application-key"
	apiKey = "some-api-key"
	id = 1234

	fakeDataDogServer = ghttp.NewServer()

	if _, err = os.Stat("/opt/resource/in"); err == nil {
		binPath = "/opt/resource/in"
	} else {
		binPath, err = gexec.Build("github.com/concourse/datadog-event-resource/cmd/in")
		Expect(err).NotTo(HaveOccurred())
	}

	tmpDir, err = ioutil.TempDir("", "datadog_event_resource_in")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterEach(func() {
	fakeDataDogServer.Close()
})

func RunIn() *gexec.Session {
	payload := cmd.InPayload{
		Source: cmd.Source{
			ApplicationKey: applicationKey,
			ApiKey:         apiKey,
		},
		Version: cmd.Version{
			Id: strconv.Itoa(id),
		},
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
