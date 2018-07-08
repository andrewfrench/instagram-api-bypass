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
	ID                 string              `json:"id"`
	Shortcode          string              `json:"shortcode"`
	TakenAtTimestamp   int                 `json:"taken_at_timestamp"`
	Owner              string              `json:"owner"`
	Caption            string              `json:"caption"`
	CommentCount       int                 `json:"comment_count"`
	LikeCount          int                 `json:"like_count"`
	IsVideo            bool                `json:"is_video"`
	CommentsDisabled   bool                `json:"comments_disabled"`
	ThumbnailSrc       string              `json:"thumbnail_src"`
	ThumbnailResources []ThumbnailResource `json:"thumbnail_resources"`
	DisplayURL         string              `json:"display_url"`
}

type ThumbnailResource struct {
	Src    string `json:"src"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
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
						Edges []struct {
							Node struct {
								ID               string `json:"id"`
								Shortcode        string `json:"shortcode"`
								TakenAtTimestamp int    `json:"taken_at_timestamp"`
								Dimensions       struct {
									Height int `json:"height"`
									Width  int `json:"width"`
								} `json:"dimensions"`
								Owner struct {
									ID string `json:"id"`
								} `json:"owner"`
								EdgeMediaToCaption struct {
									Edges []struct {
										Node struct {
											Caption string `json:"text"`
										} `json:"node"`
									} `json:"edges"`
								} `json:"edge_media_to_caption"`
								Comments struct {
									Count int `json:"count"`
								} `json:"edge_media_to_comment"`
								Likes struct {
									Count int `json:"count"`
								} `json:"edge_liked_by"`
								IsVideo            bool   `json:"is_video"`
								CommentsDisabled   bool   `json:"comments_disabled"`
								ThumbnailSrc       string `json:"thumbnail_src"`
								ThumbnailResources []struct {
									Src    string `json:"src"`
									Width  int    `json:"config_width"`
									Height int    `json:"config_height"`
								} `json:"thumbnail_resources"`
								DisplayURL string `json:"display_url"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"edge_owner_to_timeline_media"`
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
	var recentMedia []RecentMedia
	for _, m := range resp.EntryData.ProfilePage[0].GraphQL.User.Media.Edges {
		media := RecentMedia{
			ID:               m.Node.ID,
			Shortcode:        m.Node.Shortcode,
			TakenAtTimestamp: m.Node.TakenAtTimestamp,
			Owner:            m.Node.Owner.ID,
			CommentCount:     m.Node.Comments.Count,
			LikeCount:        m.Node.Likes.Count,
			IsVideo:          m.Node.IsVideo,
			CommentsDisabled: m.Node.CommentsDisabled,
			ThumbnailSrc:     m.Node.ThumbnailSrc,
			DisplayURL:       m.Node.DisplayURL,
		}

		if len(m.Node.EdgeMediaToCaption.Edges) > 0 {
			media.Caption = m.Node.EdgeMediaToCaption.Edges[0].Node.Caption
		}

		for _, tr := range m.Node.ThumbnailResources {
			resource := ThumbnailResource{
				Src:    tr.Src,
				Width:  tr.Width,
				Height: tr.Height,
			}

			media.ThumbnailResources = append(media.ThumbnailResources, resource)
		}

		recentMedia = append(recentMedia, media)
	}

	return recentMedia
}
