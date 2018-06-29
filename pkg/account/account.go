package account

import (
	"encoding/json"
	"fmt"

	"github.com/andrewfrench/instagram-api-bypass/pkg/common/extract"
	"github.com/andrewfrench/instagram-api-bypass/pkg/common/request"
)

type Account struct {
	ID            string        `json:"id"`
	Username      string        `json:"username"`
	FullName      string        `json:"full_name"`
	ExternalURL   string        `json:"external_url"`
	Biography     string        `json:"biography"`
	ProfilePicURL string        `json:"profile_pic_url_hd"`
	Followers     int           `json:"followers"`
	Following     int           `json:"following"`
	IsPrivate     bool          `json:"is_private"`
	IsVerified    bool          `json:"is_verified"`
	RecentMedia   []RecentMedia `json:"recent_media"`
}

type RecentMedia struct {
	ID               string `json:"id"`
	Code             string `json:"code"`
	Date             int    `json:"date"`
	Owner            string `json:"owner"`
	Caption          string `json:"caption"`
	Comments         int    `json:"comments"`
	Likes            int    `json:"likes"`
	IsVideo          bool   `json:"is_video"`
	CommentsDisabled bool   `json:"comments_disabled"`
	ThumbnalSrc      string `json:"thumbnail_src"`
	DisplaySrc       string `json:"display_src"`
}

type Response struct {
	EntryData struct {
		ProfilePage []struct {
			GraphQL struct {
				User struct {
					Biography      string `json:"biography"`
					ExternalURL    string `json:"external_url"`
					EdgeFollowedBy struct {
						Count int `json:"count"`
					} `json:"edge_followed_by"`
					EdgeFollow struct {
						Count int `json:"count"`
					} `json:"edge_follow"`
					FullName      string `json:"full_name"`
					ID            string `json:"id"`
					IsPrivate     bool   `json:"is_private"`
					IsVerified    bool   `json:"is_verified"`
					Username      string `json:"username"`
					ProfilePicURL string `json:"profile_pic_url"`
					CountryBlock  bool   `json:"country_block"`
					Media         struct {
						Nodes []struct {
							ID    string `json:"id"`
							Code  string `json:"code"`
							Date  int    `json:"date"`
							Owner struct {
								ID string `json:"id"`
							} `json:"owner"`
							Caption  string `json:"caption"`
							Comments struct {
								Count int `json:"count"`
							} `json:"comments"`
							Likes struct {
								Count int `json:"count"`
							} `json:"likes"`
							IsVideo          bool   `json:"is_video"`
							CommentsDisabled bool   `json:"comments_disabled"`
							ThumbnailSrc     string `json:"thumbnail_src"`
							DisplaySrc       string `json:"display_src"`
						} `json:"nodes"`
					} `json:"media"`
				} `json:"user"`
			} `json:"graphql"`
		} `json:"ProfilePage"`
	} `json:"entry_data"`
}

func Get(username string) (*Account, error) {
	url := fmt.Sprintf("https://www.instagram.com/%s/", username)
	response, err := request.Get(url)
	if err != nil {
		return &Account{}, err
	}

	return parseAccount(response)
}

func parseAccount(input []byte) (*Account, error) {
	jsonBytes, err := extract.ExtractJson(input)
	if err != nil {
		return &Account{}, err
	}

	response := &Response{}
	err = json.Unmarshal(jsonBytes, response)
	if err != nil {
		return &Account{}, err
	}

	acc := responseToAccount(response)
	acc.RecentMedia = responseToRecentMediaSlice(response)

	return acc, err
}

func responseToAccount(resp *Response) *Account {
	u := resp.EntryData.ProfilePage[0].GraphQL.User
	return &Account{
		ID:            u.ID,
		Username:      u.Username,
		FullName:      u.FullName,
		Biography:     u.Biography,
		ProfilePicURL: u.ProfilePicURL,
		Followers:     u.EdgeFollowedBy.Count,
		Following:     u.EdgeFollow.Count,
		IsPrivate:     u.IsPrivate,
		IsVerified:    u.IsVerified,
	}
}

func responseToRecentMediaSlice(resp *Response) []RecentMedia {
	recentMedia := []RecentMedia{}
	for _, m := range resp.EntryData.ProfilePage[0].GraphQL.User.Media.Nodes {
		media := RecentMedia{
			ID:               m.ID,
			Code:             m.Code,
			Date:             m.Date,
			Caption:          m.Caption,
			Owner:            m.Owner.ID,
			Comments:         m.Comments.Count,
			Likes:            m.Likes.Count,
			IsVideo:          m.IsVideo,
			CommentsDisabled: m.CommentsDisabled,
			ThumbnalSrc:      m.ThumbnailSrc,
			DisplaySrc:       m.DisplaySrc,
		}

		recentMedia = append(recentMedia, media)
	}

	return recentMedia
}