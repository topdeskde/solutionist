package main

import (
	"fmt"
	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestDownloadBuildGradle(t *testing.T) {
	g := Goblin(t)

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Trying to download a build.gradle with authentication", func() {
		var ts *httptest.Server
		targetFolder := "test"
		username := "myusername"
		password := "mypassword"
		//username := ""
		//password := ""

		g.Before(func() {
			ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/download_with_basic_auth" {
					auth_array := r.Header["Authorization"]
					if len(auth_array) > 0 {
						auth := strings.TrimSpace(auth_array[0])
						w.WriteHeader(200)
						fmt.Fprint(w, auth)
					} else {
						w.WriteHeader(401)
						fmt.Fprint(w, "private")
					}
				}
			}))

			err := os.MkdirAll(targetFolder, 0777)
			if err != nil {
				fmt.Errorf("Target directory could not be created: %s", err)
			}
		})
		g.It("Should download a build.gradle", func() {
			url := ts.URL + "/download_with_basic_auth"
			downloadFromUrl(url, targetFolder, "build.gradle", username, password)
			_, err := os.Stat(targetFolder + "/build.gradle")
			Expect(err).Should(BeNil())

		})
		g.After(func() {
			err := os.RemoveAll(targetFolder)
			if err != nil {
				fmt.Errorf("Target directory could not be removed: %s", err)
			}
			ts.Close()
		})
	})

}
