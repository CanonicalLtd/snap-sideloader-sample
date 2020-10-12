package service

import (
	"fmt"
	"github.com/slimjim777/snap-sideloader/snapd"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Web is the web service
type Web struct {
	host        string
	port        string
	snapdClient snapd.Client
}

// NewWebService creates a web service
func NewWebService(host, port string, client snapd.Client) *Web {
	return &Web{port: port, host: host, snapdClient: client}
}

// Start the web service
func (web *Web) Start() error {
	addr := fmt.Sprintf("%s:%s", web.host, web.port)
	log.Println("Starting web server at", addr)

	http.HandleFunc("/", web.handler)
	return http.ListenAndServe(addr, nil)
}

func (web *Web) handler(w http.ResponseWriter, r *http.Request) {
	first, second, err := parseURL(r.URL.Path)
	if err != nil {
		formatStandardResponse("error", err.Error(), w)
		return
	}

	switch {
	case first == "list" && second == "":
		resp, err := web.snapdClient.List()
		if err != nil {
			formatStandardResponse("error", err.Error(), w)
			return
		}
		fmt.Fprintf(w, string(resp))

	case first != "" && second != "":
		// sideload the snap from the predefined path
		if err := web.snapdClient.SideloadInstall(first, second); err != nil {
			formatStandardResponse("error", err.Error(), w)
			return
		}
		formatStandardResponse("", "Snap submitted", w)

	default:
		formatStandardResponse("error", "Incorrect URL format, use: /snap-name/revision", w)
	}

}

func parseURL(urlPath string) (string, string, error) {
	p := strings.Split(urlPath, "/")[1:]

	if len(p) == 1 {
		return p[0], "", nil
	}
	if len(p) != 2 {
		return "", "", fmt.Errorf("incorrect URL format, use: /snap-name/revision")
	}

	// check that the revision is an integer
	if _, err := strconv.Atoi(p[1]); err != nil {
		return "", "", fmt.Errorf("incorrect URL format, use: /snap-name/revision")
	}

	return p[0], p[1], nil
}
