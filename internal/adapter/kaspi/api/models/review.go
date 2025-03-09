package models

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

const (
	ReviewQualityPositive ReviewQuality = "POSITIVE"
	ReviewQualityNegative ReviewQuality = "NEGATIVE"
)

type (
	GetReviewsRequest struct {
		PageNumber int
		PageSize   int
		Filter     ReviewFilter
	}

	ReviewFilter struct {
		MerchantReviews MerchantReviews
	}

	MerchantReviews struct {
		ApprovedDateGe time.Time
		ApprovedDateLe time.Time
		Quality        ReviewQuality
	}

	ReviewsResponse struct {
		Data struct {
			Id            string           `json:"id"`
			Type          string           `json:"type"`
			Attributes    ReviewAttributes `json:"attributes"`
			Relationships struct {
				Order Relationship `json:"order"`
			} `json:"relationships"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
		Included []interface{} `json:"included"`
	}

	ReviewResponse struct {
		Data struct {
			Id            string           `json:"id"`
			Type          string           `json:"type"`
			Attributes    ReviewAttributes `json:"attributes"`
			Relationships struct {
				Order Relationship `json:"order"`
			} `json:"relationships"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
		Included []interface{} `json:"included"`
	}

	ReviewAttributes struct {
		Rating   int    `json:"rating"`
		Comment  string `json:"comment"`
		UserName string `json:"userName"`
	}

	ReviewQuality string
)

var (
	ErrApprovedDateGe = errors.New("filter merchant reviews approved date ge must be greater than 0")
	ErrApprovedDateLe = errors.New("filter merchant reviews approved date le must be greater than 0")
)

func (r *GetReviewsRequest) ToUrlValues() (url.Values, error) {
	if r.PageNumber <= 0 {
		return nil, ErrPageNumber
	}
	if r.PageSize <= 0 {
		return nil, ErrPageSize
	}
	if r.PageSize > 100 {
		return nil, ErrPageSizeLimit
	}
	if r.Filter.MerchantReviews.ApprovedDateGe.Unix() <= 0 {
		return nil, ErrApprovedDateGe
	}
	if r.Filter.MerchantReviews.ApprovedDateLe.Unix() <= 0 {
		return nil, ErrApprovedDateLe
	}

	params := make(url.Values)
	params.Add("page[number]", strconv.Itoa(r.PageNumber))
	params.Add("page[size]", strconv.Itoa(r.PageSize))
	if r.Filter.MerchantReviews.Quality != "" {
		params.Add("filter[merchantreviews][quality]", string(r.Filter.MerchantReviews.Quality))
	}
	params.Add("filter[merchantreviews][approvedDate][$ge]", strconv.FormatInt(r.Filter.MerchantReviews.ApprovedDateGe.Unix(), 10))
	params.Add("filter[merchantreviews][approvedDate][$le]", strconv.FormatInt(r.Filter.MerchantReviews.ApprovedDateLe.Unix(), 10))
	return params, nil
}
