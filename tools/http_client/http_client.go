package http_client

import (
	"app/tools/conv"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	domain  string
	data    map[string]any
	cli     *http.Client
	headers map[string]string
}

func (c *Client) Headers() map[string]string {
	return c.headers
}

func (c *Client) SetHeaders(headers map[string]string) {
	c.headers = headers
}

func (c *Client) ClearHeader() {
	c.headers = map[string]string{}
}

func (c *Client) C() *http.Client {
	return c.cli
}

func (c *Client) SetC(cli *http.Client) {
	c.cli = cli
}

func (c *Client) Data() map[string]any {
	return c.data
}

func (c *Client) SetData(data map[string]any) {
	c.data = data
}

func (c *Client) ClearData() {
	c.data = map[string]any{}
}

func (c *Client) Domain() string {
	return c.domain
}

func (c *Client) SetDomain(domain string) {
	c.domain = domain
}

func (c *Client) Get(path string) (error, []byte, int) {
	var query string
	if len(c.Data()) > 0 {
		query = "?"
		for k, v := range c.Data() {
			val, _ := conv.Conv[string](v)
			if query != "?" {
				query = query + "&" + k + "=" + val
			} else {
				query = query + k + "=" + val
			}
		}
		path += query
	}

	req, _ := http.NewRequest("GET", c.domain+path, nil)
	if len(c.Headers()) > 0 {
		for k, v := range c.Headers() {
			req.Header.Add(k, v)
		}
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return err, nil, 0
	}
	defer func(Body io.ReadCloser) {
		c.ClearData()
		err := Body.Close()
		if err != nil {
			fmt.Println("http client close fail")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read body failed:", err)
		return err, nil, 0
	}
	return nil, body, resp.StatusCode
}

func (c *Client) Post(path string) (error, []byte, int) {
	d := c.formatData()
	req, _ := http.NewRequest("POST", c.domain+path, d)
	if len(c.Headers()) > 0 {
		for k, v := range c.Headers() {
			req.Header.Add(k, v)
		}
	}
	resp, err := c.cli.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return err, nil, 0
	}
	defer func(Body io.ReadCloser) {
		c.ClearData()
		err := Body.Close()
		if err != nil {
			fmt.Println("http client close fail")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read body failed:", err)
		return err, nil, 0
	}
	return nil, body, resp.StatusCode
}

func (c *Client) formatData() *bytes.Buffer {
	contentType, ok := c.Headers()["Content-Type"]
	if !ok {
		panic("empty content type")
	}
	if contentType == "application/json" {
		str, err := json.Marshal(c.Data())
		if err != nil {
			panic("json.Marshal failed:" + err.Error())
		}
		return bytes.NewBuffer(str)
	} else {
		// form-data
		var query string
		if len(c.Data()) > 0 {
			query = ""
			for k, v := range c.Data() {
				val, _ := conv.Conv[string](v)
				if query != "" {
					query = query + "&" + k + "=" + val
				} else {
					query = query + k + "=" + val
				}
			}
		}
		return bytes.NewBuffer([]byte(query))
	}

}
