package zeptomail

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// validate runs a validation on the incoming json payload
var validate = validator.New(validator.WithRequiredStructEnabled())

// Client is an object for the configs
type Client struct {
	Http    *http.Client
	BaseUrl string
	Token   string
}

// New initializes the ZeptoMail client configuration, with optional base URL support.
func New(httpClient *http.Client, token string, baseURL ...string) *Client {
	var url string
	if len(baseURL) > 0 {
		url = baseURL[0]
	} else {
		url = "https://api.zeptomail.com/v1.1/"
	}
	return &Client{
		BaseUrl: url,
		Http:    httpClient,
		Token:   token,
	}
}

// newRequest makes a http request to the zeptomail server and decodes the server response into the reqBody parameter passed into the newRequest method
func (c *Client) newRequest(method, reqURL string, reqBody, resp interface{}) error {
	newURL := c.BaseUrl + reqURL
	var body io.Reader

	if reqBody != nil {
		bb, err := json.Marshal(reqBody)
		if err != nil {
			return errors.Wrap(err, "http client ::: unable to marshal request struct")
		}
		body = bytes.NewReader(bb)
	}

	req, err := http.NewRequest(method, newURL, body)
	if err != nil {
		return errors.Wrap(err, "http client ::: unable to create request body")
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", c.Token)
	}

	res, err := c.Http.Do(req)
	if err != nil {
		return errors.Wrap(err, "http client ::: client failed to execute request")
	}
	defer res.Body.Close()

	bb, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "http client ::: client failed to read file")
	}

	// fmt.Println("Response body:", string(bb))

	if err := json.Unmarshal(bb, &resp); err != nil {
		return errors.Errorf("Error ::: %v", err)
	}

	return nil
}
