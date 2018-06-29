package account

import (
	"fmt"
	"common/request"
	"encoding/json"
	"common/extract"
)

type Account struct {
	Id            string        `json:"id"`
	Username      string        `json:"username"`
	FullName      string        `json:"full_name"`
	Biography     string        `json:"biography"`
	ProfilePicUrl string        `json:"profile_pic_url"`
	Followers     int           `json:"followers"`
	Following     int           `json:"following"`
	IsPrivate     bool          `json:"is_private"`
	IsVerified    bool          `json:"is_verified"`
	RecentMedia   []RecentMedia `json:"recent_media"`
}

type RecentMedia struct {
	Id               string `json:"id"`
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
			User struct {
				Biography string `json:"biography"`
				FollowedBy struct {
					Count int `json:"count"`
				} `json:"followed_by"`
				Follows struct {
					Count int `json:"count"`
				}
				FullName string `json:"full_name"`
				Id string `json:"id"`
				IsPrivate bool `json:"is_private"`
				IsVerified bool `json:"is_verified"`
				Username string `json:"username"`
				ProfilePicUrl string `json:"profile_pic_url"`
				Media struct {
					Nodes []struct {
						Id string `json:"id"`
						Code string `json:"code"`
						Date int `json:"date"`
						Owner struct {
							Id string `json:"id"`
						} `json:"owner"`
						Caption string `json:"caption"`
						Comments struct {
							Count int `json:"count"`
						} `json:"comments"`
						Likes struct {
							Count int `json:"count"`
						} `json:"likes"`
						IsVideo bool `json:"is_video"`
						CommentsDisabled bool `json:"comments_disabled"`
						ThumbnailSrc string `json:"thumbnail_src"`
						DisplaySrc string `json:"display_src"`
					} `json:"nodes"`
				} `json:"media"`
			} `json:"user"`
		} `json:"ProfilePage"`
	} `json:"entry_data"`
}

func Get(username string) (*Account, error) {
	url := fmt.Sprintf("http://www.instagram.com/%s/", username)
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
	u := resp.EntryData.ProfilePage[0].User
	return &Account{
		Id:            u.Id,
		Username:      u.Username,
		FullName:      u.FullName,
		Biography:     u.Biography,
		ProfilePicUrl: u.ProfilePicUrl,
		Followers:     u.FollowedBy.Count,
		Following:     u.Follows.Count,
		IsPrivate:     u.IsPrivate,
		IsVerified:    u.IsVerified,
	}
}

func responseToRecentMediaSlice(resp *Response) []RecentMedia {
	recentMedia := []RecentMedia{}
	for _, m := range resp.EntryData.ProfilePage[0].User.Media.Nodes {
		media := RecentMedia{
			Id:               m.Id,
			Code:             m.Code,
			Date:             m.Date,
			Caption:          m.Caption,
			Owner:            m.Owner.Id,
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
