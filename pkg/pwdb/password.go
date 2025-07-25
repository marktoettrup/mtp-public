package pwdb

import (
	"context"
	"net/http"
	"path"
	"strconv"
)

const endpointPassword = "passwords"

type PasswordResponse struct {
	PasswordID int    `json:"PasswordID"`
	Username   string `json:"UserName"`
	Password   string `json:"Password"`
	Notes      string `json:"Notes"`
}

type PasswordRequest struct {
	PasswordListID string `json:"PasswordListID"`
	Title          string `json:"Title"`
	UserName       string `json:"UserName"`
	Password       string `json:"Password"`
	PasswordID     int    `json:"PasswordID"`
}

func (c *Client) GetPassword(ctx context.Context, id int) (*PasswordResponse, error) {
	endpoint := path.Join(endpointPassword, strconv.Itoa(id))
	req, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	r := make([]*PasswordResponse, 1)
	if err := c.do(req, &r); err != nil {
		return nil, err
	}

	return r[0], nil
}

func (c *Client) CreatePassword(ctx context.Context, r PasswordRequest) (*PasswordResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, endpointPassword, r)
	if err != nil {
		return nil, err
	}

	res := make([]*PasswordResponse, 1)
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	return res[0], nil
}

func (c *Client) UpdatePassword(ctx context.Context, r PasswordRequest) (*PasswordResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPut, endpointPassword, r)
	if err != nil {
		return nil, err
	}

	res := &PasswordResponse{}
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}
