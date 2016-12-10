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
)

var _ = BeforeEach(func() {
	fakeDataDogServer = ghttp.NewServer()

	binPath, err = gexec.Build("github.com/concourse/datadog-resource/cmd/check")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterEach(func() {
	fakeDataDogServer.Close()
})

func RunCheck(payload cmd.CheckPayload) *gexec.Session {
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
