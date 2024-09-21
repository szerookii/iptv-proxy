package stb

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/szerookii/iptv-proxy/utils"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL     string // URL of the STB server
	macAddress  string // MAC address of the STB
	bearerToken string // retrieved when handshake is successful
}

type Response[T any] struct {
	Data T `json:"js"`
}

type HandshakeResponse struct {
	Token string `json:"token"`
}

func NewClient(baseURL, macAddress string) (*Client, error) {
	c := &Client{
		baseURL:    baseURL,
		macAddress: macAddress,
	}

	token, err := c.Handshake()
	if err != nil {
		return nil, err
	}

	c.bearerToken = token

	// TODO: check valid subscription

	return c, nil
}

func (c *Client) Handshake() (string, error) {
	bodyBytes, err := c.doRequest("/portal.php?", prepareURLValues("stb", "handshake", nil))
	if err != nil {
		return "", err
	}

	var handshake Response[HandshakeResponse]
	if err := sonic.Unmarshal(bodyBytes, &handshake); err != nil {
		return "", err
	}

	return handshake.Data.Token, nil
}

func prepareURLValues(t, action string, extra map[string]string) url.Values {
	urlValues := url.Values{}
	urlValues.Add("type", t)
	urlValues.Add("action", action)

	for k, v := range extra {
		urlValues.Add(k, v)
	}

	urlValues.Add("JsHttpRequest", "1-xml")

	return urlValues
}

func (c *Client) doRequest(endpoint string, urlValues url.Values) ([]byte, error) {
	u := c.baseURL + endpoint + urlValues.Encode()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	if c.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	}

	req.Header.Set("X-User-Agent", "Model: MAG520; Link: Ethernet")

	var cookiesMap = map[string]string{
		"stb_lang": "en",
		"timezone": "Europe/Germany",
	}

	if c.macAddress != "" {
		cookiesMap["mac"] = c.macAddress
	}

	req.Header.Set("Cookie", utils.Cookies2header(cookiesMap))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(bodyBytes) <= 0 {
		return nil, errors.New("empty response")
	}

	return bodyBytes, nil
}
