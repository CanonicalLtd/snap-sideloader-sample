package snapd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const (
	socketFile     = "/run/snapd.socket"
	urlAssertions  = "/v2/assertions"
	urlSnaps       = "/v2/snaps"
	typeAssertions = "application/x.ubuntu.assertion"
	baseURL        = "http://localhost"
)

// Client is the abstract client interface
type Client interface {
	Ack(assertion []byte) error
	InstallPath(name, filePath string) error
	List() ([]byte, error)
	SideloadInstall(name, revision string) error
}

// Snapd service to access the snapd REST API
type Snapd struct {
	downloadPath string
	client       *http.Client
}

// NewClient returns a snapd API client
func NewClient(downloadPath string) *Snapd {
	return &Snapd{
		downloadPath: downloadPath,
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", socketFile)
				},
			},
		},
	}
}

func (snap *Snapd) call(method, url, contentType string, body io.Reader) (*http.Response, error) {
	u := baseURL + url

	switch method {
	case "POST":
		return snap.client.Post(u, contentType, body)
	case "GET":
		return snap.client.Get(u)
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
}

// Ack acknowledges a (snap) assertion
func (snap *Snapd) Ack(assertion []byte) error {
	_, err := snap.call("POST", urlAssertions, typeAssertions, bytes.NewReader(assertion))
	return err
}

// InstallPath installs a snap from a local file
func (snap *Snapd) InstallPath(name, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open: %q", filePath)
	}

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	go sendSnapFile(name, filePath, f, pw, mw)

	_, err = snap.call("POST", urlSnaps, mw.FormDataContentType(), pr)
	return err
}

// List the installed snaps
func (snap *Snapd) List() ([]byte, error) {
	resp, err := snap.call("GET", urlSnaps, "application/json; charset=UTF-8", nil)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

// SideloadInstall side loads a snap by acknowledging the assertion and installing the snap
func (snap *Snapd) SideloadInstall(name, revision string) error {
	assertsPath := path.Join(snap.downloadPath, fmt.Sprintf("%s_%s.assert", name, revision))
	snapPath := path.Join(snap.downloadPath, fmt.Sprintf("%s_%s.snap", name, revision))

	// acknowledge the snap assertion
	dataAssert, err := ioutil.ReadFile(assertsPath)
	if err != nil {
		return err
	}
	if err := snap.Ack(dataAssert); err != nil {
		return err
	}

	// install the snap
	return snap.InstallPath(name, snapPath)
}

func sendSnapFile(name, snapPath string, snapFile *os.File, pw *io.PipeWriter, mw *multipart.Writer) {
	defer snapFile.Close()

	fields := []struct {
		name  string
		value string
	}{
		{"action", "install"},
		{"name", name},
		{"snap-path", snapPath},
	}
	for _, s := range fields {
		if s.value == "" {
			continue
		}
		if err := mw.WriteField(s.name, s.value); err != nil {
			pw.CloseWithError(err)
			return
		}
	}

	fw, err := mw.CreateFormFile("snap", filepath.Base(snapPath))
	if err != nil {
		pw.CloseWithError(err)
		return
	}

	_, err = io.Copy(fw, snapFile)
	if err != nil {
		pw.CloseWithError(err)
		return
	}

	mw.Close()
	pw.Close()
}
