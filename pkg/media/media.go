package media

import (
	"fmt"
	"common/request"
	"encoding/json"
	"common/extract"
)

type Media struct {
	Id               string `json:"id"`
	Code             string `json:"code"`
	Date             int `json:"date"`
	Caption          string `json:"caption"`
	IsVideo          bool `json:"is_video"`
	IsAd             bool `json:"is_ad"`
	CommentsDisabled bool `json:"comments_disabled"`
	CommentCount     int `json:"comment_count"`
	Comments         []MediaComment `json:"comments"`
	DisplayUrl       string `json:"display"`
	Likes            int `json:"likes"`
	OwnerId          string `json:"owner_id"`
	OwnerUsername	 string `json:"owner_username"`
	Width            int `json:"width"`
	Height           int `json:"height"`
}

type MediaComment struct {
	Id string `json:"id"`
	Text string `json:"text"`
	CreatedAt int `json:"created_at"`
	OwnerId string `json:"owner_id"`
	OwnerUsername string `json:"owner_username"`
	OwnerProfilePicUrl string `json:"owner_profile_pic_url"`
}

type Response struct {
	EntryData struct {
		PostPage []struct {
			GraphQL struct {
				ShortcodeMedia struct {
					Id string `json:"id"`
					ShortCode string `json:"shortcode"`
					Dimensions struct {
						Width int `json:"width"`
						Height int `json:"height"`
					} `json:"dimensions"`
					DisplayUrl string `json:"display_url"`
					IsVideo bool `json:"is_video"`
					EdgeMediaPreviewLike struct {
						Count int `json:"count"`
					} `json:"edge_media_preview_like"`
					EdgeMediaToComment struct {
						Count int `json:"count"`
						Edges []struct {
							Node struct {
								Id string `json:"id"`
								Text string `json:"text"`
								CreatedAt int `json:"created_at"`
								Owner struct {
									Id string `json:"id"`
									Username string `json:"username"`
									ProfilePicUrl string `json:"profile_pic_url"`
								} `json:"owner"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"edge_media_to_comment"`
					Owner struct {
						Id string `json:"id"`
						Username string `json:"username"`
					} `json:"owner"`
					IsAd bool `json:"is_ad"`
					TakenAtTimestamp int `json:"taken_at_timestamp"`
				} `json:"shortcode_media"`
			} `json:"graphql"`
		} `json:"PostPage"`
	} `json:"entry_data"`
}

func Get(code string) (*Media, error) {
	url := fmt.Sprintf("https://www.instagram.com/p/%s/", code)
	resp, err := request.Get(url)
	if err != nil {
		return &Media{}, err
	}

	return parseMedia(resp)
}

func parseMedia(input []byte) (*Media, error) {
	jsonBytes, err := extract.ExtractJson(input)
	if err != nil {
		return &Media{}, err
	}

	response := &Response{}
	err = json.Unmarshal(jsonBytes, response)
	if err != nil {
		return &Media{}, err
	}

	med := responseToMedia(response)
	med.Comments = responseToCommentSlice(response)

	return med, err
}

func responseToMedia(resp *Response) *Media {
	m := resp.EntryData.PostPage[0].GraphQL.ShortcodeMedia
	return &Media{
		Id: m.Id,
		Code: m.ShortCode,
		Date: m.TakenAtTimestamp,
		CommentCount: m.EdgeMediaToComment.Count,
		IsVideo: m.IsVideo,
		IsAd: m.IsAd,
		DisplayUrl: m.DisplayUrl,
		Likes: m.EdgeMediaPreviewLike.Count,
		Width: m.Dimensions.Width,
		Height: m.Dimensions.Height,
		OwnerId: m.Owner.Id,
		OwnerUsername: m.Owner.Username,
	}
}

func responseToCommentSlice(resp *Response) []MediaComment {
	commentSlice := []MediaComment{}
	for _, c := range resp.EntryData.PostPage[0].GraphQL.ShortcodeMedia.EdgeMediaToComment.Edges {
		comment := MediaComment{
			Id: c.Node.Id,
			Text: c.Node.Text,
			CreatedAt: c.Node.CreatedAt,
			OwnerId: c.Node.Owner.Id,
			OwnerUsername: c.Node.Owner.Username,
			OwnerProfilePicUrl: c.Node.Owner.ProfilePicUrl,
		}

		commentSlice = append(commentSlice, comment)
	}

	return commentSlice
}
