package client

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/isucon10-qualify/isucon10-qualify/bench/parameter"
)

func NewClient(isBot bool) *Client {
	var userAgent string
	if isBot {
		userAgent = GenerateBotUserAgent()
	} else {
		userAgent = GenerateUserAgent()
	}

	return &Client{
		userAgent: userAgent,
		isBot:     isBot,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// HTTPの時は無視されるだけ
					ServerName: ShareTargetURLs.TargetHost,
				},
			},
			Timeout: parameter.DefaultAPITimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return fmt.Errorf("redirect attempted")
			},
		},
	}
}

func NewClientForInitialize() *Client {
	return &Client{
		userAgent: "isucon-initialize",
		isBot:     false,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// HTTPのときには無視される
					ServerName: ShareTargetURLs.TargetHost,
				},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return fmt.Errorf("redirect attempted")
			},
		},
	}
}

func NewClientForVerify() *Client {
	return &Client{
		userAgent: "isucon-verify",
		isBot:     false,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// HTTPのときには無視される
					ServerName: ShareTargetURLs.TargetHost,
				},
			},
			Timeout: parameter.DefaultAPITimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return fmt.Errorf("redirect attempted")
			},
		},
	}
}

func NewClientForDraft() *Client {
	return &Client{
		userAgent: GenerateUserAgent(),
		isBot:     false,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// HTTPのときには無視される
					ServerName: ShareTargetURLs.TargetHost,
				},
			},
			Timeout: parameter.DraftTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return fmt.Errorf("redirect attempted")
			},
		},
	}
}
