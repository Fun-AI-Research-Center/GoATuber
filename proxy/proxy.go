package proxy

import (
	"net/http"
	"net/url"

	"GoATuber-2.0/engine"
)

func Client(e *engine.Engine) (http.Client, error) {
	config := e.Config.Proxy
	if config.UseProxy == false {
		return http.Client{}, nil
	}
	u, err := url.Parse(config.ProxyURL)
	if err != nil {
		return http.Client{}, err
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(u),
		},
	}
	return client, nil
}
