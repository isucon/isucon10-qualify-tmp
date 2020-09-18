package client

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/morikuni/failure"
	// "github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
)

type Client struct {
	userAgent  string
	isBot      bool
	httpClient *http.Client
}

type TargetURLs struct {
	AppURL     url.URL
	TargetHost string
}

var (
	ShareTargetURLs *TargetURLs
)

func SetShareTargetURLs(appURL, targetHost string) error {
	var err error
	ShareTargetURLs, err = newTargetURLs(appURL, targetHost)
	if err != nil {
		return err
	}

	return nil
}

func newTargetURLs(appURL, targetHost string) (*TargetURLs, error) {
	if len(appURL) == 0 {
		return nil, fmt.Errorf("client: missing url")
	}

	appParsedURL, err := urlParse(appURL)
	if err != nil {
		return nil, failure.Wrap(err, failure.Messagef("failed to parse url: %s", appURL))
	}

	return &TargetURLs{
		AppURL:     *appParsedURL,
		TargetHost: targetHost,
	}, nil
}

func urlParse(ref string) (*url.URL, error) {
	u, err := url.Parse(ref)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, fmt.Errorf("host is empty")
	}

	return &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
	}, nil
}

func (c *Client) newGetRequest(u url.URL, spath string) (*http.Request, error) {
	if len(spath) > 0 {
		u.Path = spath
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Host = ShareTargetURLs.TargetHost
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

func (c *Client) newGetRequestWithQuery(u url.URL, spath string, q url.Values) (*http.Request, error) {
	if len(spath) > 0 {
		u.Path = spath
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Host = ShareTargetURLs.TargetHost
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

func (c *Client) newPostRequest(u url.URL, spath string, body io.Reader) (*http.Request, error) {
	u.Path = spath

	req, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Host = ShareTargetURLs.TargetHost
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

func checkStatusCode(res *http.Response, expectedStatusCodes []int) error {
	for _, expectedStatusCode := range expectedStatusCodes {
		if res.StatusCode == expectedStatusCode {
			return nil
		}
	}

	return failure.New(fails.ErrApplication)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		if nerr, ok := err.(net.Error); ok {
			if nerr.Timeout() {
				return nil, failure.Translate(err, fails.ErrTimeout)
			} else if nerr.Temporary() {
				return nil, failure.Translate(err, fails.ErrTemporary)
			}
		}

		return nil, err
	}

	if !c.isBot && res.StatusCode == http.StatusServiceUnavailable {
		return nil, failure.New(fails.ErrTemporary)
	}

	return res, nil
}

func (c *Client) GetEmail() string {
	return fmt.Sprintf("%s@isucon.com", c.userAgent)
}
