package pwdb

import (
	"bytes"
	"context"
	"net/http"
	"path"
	"strconv"
)

const endpointDocument = "document/password"

type DocumentResponse struct{}

func (c *Client) GetDocument(ctx context.Context, id int) ([]byte, error) {
	endpoint := path.Join(endpointDocument, strconv.Itoa(id))
	req, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	bs, err := c.doRaw(req)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func (c *Client) CreateDocument(
	ctx context.Context,
	passwordID int,
	name string,
	description string,
	data []byte,
) (*DocumentResponse, error) {
	endpoint := path.Join(endpointDocument, strconv.Itoa(passwordID))
	req, err := c.newRawRequest(ctx, http.MethodPost, endpoint, map[string]string{
		"DocumentName":        name,
		"DocumentDescription": description,
	}, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var res DocumentResponse
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
