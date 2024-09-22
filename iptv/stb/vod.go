package stb

import (
	"fmt"
	"github.com/bytedance/sonic"
	"strconv"
)

type VodCategory struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Alias    string `json:"alias"`
	Censored int    `json:"censored"`
}

type Vod struct {
	Id                   string        `json:"id"`
	Owner                string        `json:"owner"`
	Name                 string        `json:"name"`
	OldName              string        `json:"old_name"`
	OName                string        `json:"o_name"`
	Fname                string        `json:"fname"`
	Description          string        `json:"description"`
	Pic                  string        `json:"pic"`
	Cost                 int           `json:"cost"`
	Time                 any           `json:"time"` // TODO: find a fix for time sometimes being "N/A" for no reason
	File                 string        `json:"file"`
	Path                 string        `json:"path"`
	Protocol             string        `json:"protocol"`
	RtspUrl              string        `json:"rtsp_url"`
	Censored             int           `json:"censored"`
	Series               []interface{} `json:"series"`
	VolumeCorrection     int           `json:"volume_correction"`
	CategoryId           string        `json:"category_id"`
	GenreId              int           `json:"genre_id"`
	GenreId1             int           `json:"genre_id_1"`
	GenreId2             int           `json:"genre_id_2"`
	GenreId3             int           `json:"genre_id_3"`
	Hd                   int           `json:"hd"`
	GenreId4             int           `json:"genre_id_4"`
	CatGenreId1          string        `json:"cat_genre_id_1"`
	CatGenreId2          int           `json:"cat_genre_id_2"`
	CatGenreId3          int           `json:"cat_genre_id_3"`
	CatGenreId4          int           `json:"cat_genre_id_4"`
	Director             string        `json:"director"`
	Actors               string        `json:"actors"`
	Year                 string        `json:"year"`
	Accessed             int           `json:"accessed"`
	Status               int           `json:"status"`
	DisableForHdDevices  int           `json:"disable_for_hd_devices"`
	Added                string        `json:"added"`
	Count                int           `json:"count"`
	CountFirst05         int           `json:"count_first_0_5"`
	CountSecond05        int           `json:"count_second_0_5"`
	VoteSoundGood        int           `json:"vote_sound_good"`
	VoteSoundBad         int           `json:"vote_sound_bad"`
	VoteVideoGood        int           `json:"vote_video_good"`
	VoteVideoBad         int           `json:"vote_video_bad"`
	Rate                 string        `json:"rate"`
	LastRateUpdate       string        `json:"last_rate_update"`
	LastPlayed           string        `json:"last_played"`
	ForSdStb             int           `json:"for_sd_stb"`
	RatingImdb           string        `json:"rating_imdb"`
	RatingCountImdb      string        `json:"rating_count_imdb"`
	RatingLastUpdate     string        `json:"rating_last_update"`
	Age                  string        `json:"age"`
	HighQuality          int           `json:"high_quality"`
	RatingKinopoisk      string        `json:"rating_kinopoisk"`
	Comments             string        `json:"comments"`
	LowQuality           int           `json:"low_quality"`
	IsSeries             int           `json:"is_series"`
	YearEnd              int           `json:"year_end"`
	AutocompleteProvider string        `json:"autocomplete_provider"`
	Screenshots          string        `json:"screenshots"`
	IsMovie              int           `json:"is_movie"`
	Lock                 int           `json:"lock"`
	Fav                  int           `json:"fav"`
	ForRent              int           `json:"for_rent"`
	ScreenshotUri        string        `json:"screenshot_uri"`
	GenresStr            string        `json:"genres_str"`
	Cmd                  string        `json:"cmd"`
	Yesterday            string        `json:"yesterday"`
	WeekAndMore          string        `json:"week_and_more"`
	HasFiles             int           `json:"has_files"`
}

type VodsResponse struct {
	TotalItems   int    `json:"total_items"`
	MaxPageItems int    `json:"max_page_items"`
	SelectedItem int    `json:"selected_item"`
	CurPage      int    `json:"cur_page"`
	Data         []*Vod `json:"data"`
}

func (c *Client) VodCategories() ([]*VodCategory, error) {
	bodyBytes, err := c.doRequest("/portal.php?", prepareURLValues("vod", "get_categories", nil))
	if err != nil {
		return nil, err
	}

	var genres Response[[]*VodCategory]
	if err := sonic.Unmarshal(bodyBytes, &genres); err != nil {
		return nil, err
	}

	//genres.Data = append(genres.Data[:0], genres.Data[1:]...) // remove "All" genre because it's not a real genre bro it's a placeholder for all channels when no genre is selected

	return genres.Data, nil
}

func (c *Client) CategoryVods(categoryId string) ([]*Vod, error) {
	pageNumber := 1
	var foundVods []*Vod

	for {
		bodyBytes, err := c.doRequest("/portal.php?", prepareURLValues("vod", "get_ordered_list", map[string]string{
			"category": categoryId,
			"p":        strconv.Itoa(pageNumber),
		}))

		if err != nil {
			return nil, err
		}

		var channels Response[VodsResponse]
		if err := sonic.Unmarshal(bodyBytes, &channels); err != nil {
			fmt.Println(err)
			return nil, err
		}

		foundVods = append(foundVods, channels.Data.Data...)

		if len(foundVods) >= channels.Data.TotalItems {
			break
		}

		pageNumber++
	}

	return foundVods, nil
}
