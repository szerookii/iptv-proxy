package stb

import (
	"github.com/bytedance/sonic"
)

type MainInfoResponse struct {
	Mac   string `json:"mac"`
	Phone string `json:"phone"` // date when subscription ends, idk why it's called phone but it's a date
}

func (c *Client) MainInfo() (*MainInfoResponse, error) {
	bodyBytes, err := c.doRequest("/portal.php?", prepareURLValues("account_info", "get_main_info", nil))
	if err != nil {
		return nil, err
	}

	var mainInfo Response[MainInfoResponse]
	if err := sonic.Unmarshal(bodyBytes, &mainInfo); err != nil {
		return nil, err
	}

	return &mainInfo.Data, nil
}
