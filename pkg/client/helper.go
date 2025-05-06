package client

import "net/url"

func SetQueryParams(url *url.URL, params map[string]string) {
	values := url.Query()
	for key, value := range params {
		values.Set(key, value)
	}
	url.RawQuery = values.Encode()
}
