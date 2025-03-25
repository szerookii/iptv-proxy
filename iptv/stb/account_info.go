package stb

import (
	"errors"
	"github.com/bytedance/sonic"
	"strings"
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

func (c *Client) ConvertToXtream() (string, string, error) {
	vodCatogeries, err := c.VodCategories()
	if err != nil {
		return "", "", err
	}

	if len(vodCatogeries) == 0 {
		return "", "", errors.New("no vod categories found")
	}

	vod, err := c.CategoryVods(vodCatogeries[0].Id)
	if err != nil {
		return "", "", err
	}

	if len(vod) == 0 {
		return "", "", errors.New("no vods found")
	}

	link, err := c.CreateLink("vod", vod[0].Cmd)
	if err != nil {
		return "", "", err
	}

	args := strings.Split(link, "/")

	if len(args) < 6 {
		return "", "", errors.New("cannot be converted to xtream")
	}

	username := strings.Split(link, "/")[4]
	password := strings.Split(link, "/")[5]

	return username, password, nil
}
