package app_test

import (
	"bytes"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/buildpack/pack/logging"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpack/pack/app"
	h "github.com/buildpack/pack/testhelpers"
)

func TestRun(t *testing.T) {
	color.NoColor = true
	rand.Seed(time.Now().UTC().UnixNano())
	spec.Run(t, "run", testRun, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testRun(t *testing.T, when spec.G, it spec.S) {
	when("#Run", func() {
		var (
			subject *app.Image
			docker  *client.Client
			err error
			outBuf, errBuf bytes.Buffer
		)

		it.Before(func() {
			logger := logging.NewLogger(&outBuf, &errBuf, true, false)
			repo := "some-org/" + h.RandString(10)
			docker, err = client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.38"))
			h.CreateImageOnLocal(t, docker, repo, `
FROM hashicorp/http-echo
CMD ["-text=\"hello world\""]
`)
			subject = &app.Image{
				RepoName: repo,
				Logger: logger,
			}
		})

		it("runs an image", func() {
			go func() {
				h.Eventually(t, func() bool {
					h.AssertContains(outBuf.String(), "Server is listening")
				})
				ctx.Canel
			}
			subject.Run(ctx)
		})

		when("the process is terminated", func() {
			it("stops the running container and cleans up", func() {
			})
		})
		when("the port is not specified", func() {
			it.Before(func() {
			})

			it("gets exposed ports from the image", func() {
			})
		})
		when("custom ports bindings are defined", func() {
			it("binds simple ports from localhost to the container on the same port", func() {
			})

			it("binds each port to the container", func() {
			})
		})
	})
}
