package stb

import (
	"github.com/bytedance/sonic"
	"regexp"
	"strconv"
)

type ChannelsResponse struct {
	TotalItems   int           `json:"total_items"`
	MaxPageItems int           `json:"max_page_items"`
	SelectedItem int           `json:"selected_item"`
	CurPage      int           `json:"cur_page"`
	Data         []*ItvChannel `json:"data"`
}

type CreateLinkResponse struct {
	Id  string `json:"id"`
	Cmd string `json:"cmd"` // contains the link
}

type ItvGenre struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Alias     string `json:"alias"`
	ActiveSub bool   `json:"active_sub"`
	Censored  int    `json:"censored"`
}

type ItvChannel struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	TvGenreId string `json:"tv_genre_id"`
	Cmd       string `json:"cmd"`
	Logo      string `json:"logo"`
	Cmds      []struct {
		ID   string `json:"id"`
		ChId string `json:"ch_id"`
	} `json:"cmds"`
}

func (c *Client) GenreChannels(genreId string) ([]*ItvChannel, error) {
	pageNumber := 1
	var foundChannels []*ItvChannel

	for {
		bodyBytes, err := c.doRequest("/portal.php?", prepareURLValues("itv", "get_ordered_list", map[string]string{
			"genre": genreId,
			"p":     strconv.Itoa(pageNumber),
		}))

		if err != nil {
			return nil, err
		}

		var channels Response[ChannelsResponse]
		if err := sonic.Unmarshal(bodyBytes, &channels); err != nil {
			return nil, err
		}

		foundChannels = append(foundChannels, channels.Data.Data...)

		if len(foundChannels) >= channels.Data.TotalItems {
			break
		}

		pageNumber++
	}

	return foundChannels, nil
}

func (c *Client) ItvGenres() ([]*ItvGenre, error) {
	bodyBytes, err := c.doRequest("/portal.php?", prepareURLValues("itv", "get_genres", nil))
	if err != nil {
		return nil, err
	}

	var genres Response[[]*ItvGenre]
	if err := sonic.Unmarshal(bodyBytes, &genres); err != nil {
		return nil, err
	}

	//genres.Data = append(genres.Data[:0], genres.Data[1:]...) // remove "All" genre because it's not a real genre bro it's a placeholder for all channels when no genre is selected

	return genres.Data, nil
}

func (c *Client) Channels() ([]*ItvChannel, error) {
	bodyBytes, err := c.doRequest("/portal.php?", prepareURLValues("itv", "get_all_channels", nil))
	if err != nil {
		return nil, err
	}

	var channels Response[ChannelsResponse]
	if err := sonic.Unmarshal(bodyBytes, &channels); err != nil {
		return nil, err
	}

	return channels.Data.Data, nil
}

func (c *Client) CreateLink(t string, cmd string) (string, error) {
	bodyBytes, err := c.doRequest("/portal.php?", prepareURLValues(t, "create_link", map[string]string{
		"cmd": cmd,
		"mac": c.macAddress,
	}))
	if err != nil {
		return "", err
	}

	var response Response[CreateLinkResponse]
	if err := sonic.Unmarshal(bodyBytes, &response); err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(?m)http.*$`)
	link := re.FindString(response.Data.Cmd)

	return DecodeLink(link), nil
}

func DecodeLink(link string) string {
	re := regexp.MustCompile(`\\/`)
	link = re.ReplaceAllString(link, "/")

	return link
}

func EncodeLink(link string) string {
	re := regexp.MustCompile(`/`)
	link = re.ReplaceAllString(link, "\\/")

	return link
}
