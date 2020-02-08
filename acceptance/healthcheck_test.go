package acceptance_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Healthchecking", func() {
	var (
		session *gexec.Session
	)

	BeforeEach(func() {
		cmd := exec.Command(serverPath)
		session = execBin(cmd)
	})

	AfterEach(func() {
		session.Kill().Wait()
	})

	It("eventually responds 200 OK", func() {
		Eventually(healthcheck).Should(Succeed())
	})
})

func healthcheck() error {
	resp, err := http.Get("http://localhost:3000/healthcheck")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200 but got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != "OK" {
		return fmt.Errorf("expected body to say `OK` but got `%s`", string(body))
	}

	return nil
}
