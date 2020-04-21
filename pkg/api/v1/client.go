package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
)

const (
	defaultEndpoint = "http://127.0.0.1:8008"
)

type Client struct {
	Alerts    *Alerts
	Logs      *Logs
	Receivers *Receivers
	Snoozes   *Snoozes
	Templates *Templates

	endpoint string
	client   *http.Client
	router   *mux.Router
}

type Response struct {
	*http.Response
}

func NewClient(httpClient *http.Client, endpoint string) *Client {
	return newClient(httpClient, endpoint)
}

func (c *Client) WithRouter(router *mux.Router) *Client {
	c.router = router
	return c
}

func newClient(httpClient *http.Client, endpoint string) *Client {
	if len(endpoint) == 0 {
		endpoint = defaultEndpoint
	}

	cli := httpClient
	if cli == nil {
		cli = http.DefaultClient
	}

	c := &Client{
		endpoint: endpoint,
		client:   cli,
	}

	c.Alerts = &Alerts{c}
	c.Logs = &Logs{c}
	c.Receivers = &Receivers{c}
	c.Snoozes = &Snoozes{c}
	c.Templates = &Templates{c}

	return c
}

func (c *Client) newRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u := joinURLPath(c.endpoint, urlStr)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

func (c *Client) doRequest(req *http.Request, v interface{}) (*Response, error) {
	var (
		resp *http.Response
		err  error
	)

	if c.router != nil {
		rr := httptest.NewRecorder()
		c.router.ServeHTTP(rr, req)
		resp = rr.Result()
	} else {
		resp, err = c.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		defer io.Copy(ioutil.Discard, resp.Body)
	}

	response := newResponse(resp)

	err = checkResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	serverError := &ServerError{}
	errDec := json.Unmarshal(data, serverError)
	if errDec == nil {
		return serverError
	}

	return fmt.Errorf("failed: %s", string(data))
}

func joinURLPath(endpoint string, path string) string {
	return strings.TrimRight(endpoint, "/") + "/" + strings.TrimLeft(path, "/")
}
