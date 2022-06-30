package request

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"
)

// LogRedirects to sanitize "Location" URLs
type LogRedirects struct {
	Transport http.RoundTripper
}

func (l LogRedirects) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t := l.Transport
	if t == nil {
		t = http.DefaultTransport
	}
	resp, err = t.RoundTrip(req)
	if err != nil {
		return
	}
	switch resp.StatusCode {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
		LocationURL, _ := resp.Location()
		if !strings.ContainsAny(LocationURL.String(), " ") {
			return
		}
		resp.Header.Set("Location", strings.ReplaceAll(LocationURL.String(), " ", "%20"))
	}
	return
}

func DefaultClient() *http.Client {
	return &http.Client{
		Transport: LogRedirects{&http.Transport{
			DisableCompression:  true,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify:       true,
				PreferServerCipherSuites: false,
				CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.CurveP384, tls.CurveP521, tls.X25519},
			},
			IdleConnTimeout: 5 * time.Second,
		}},
		Timeout: 10 * time.Minute,
	}
}
