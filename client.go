package api

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/dghubble/sling"
)

// Header to send along the HTTP Requests
type Header struct {
	Key   string
	Value string
}

// ClientOptions to setup client
type ClientOptions struct {
	BaseURL     string
	BaseHeaders []Header
}

// ClientRequestOptions ...
type ClientRequestOptions struct {
	Path       string
	AddHeaders []Header
}

// Client to make API requests.
type Client struct {
	SlingBase *sling.Sling
}

// newClientWithHeaders ..
func (c Client) newClientWithHeaders(options ClientRequestOptions) *sling.Sling {

	client := c.SlingBase.New() //.Post("").Add("X-HASURA-ROLE", "")
	for _, header := range options.AddHeaders {
		client.Add(header.Key, header.Value)
	}

	return client

}

// ErrorPayload is used to handle graphql errors
type ErrorPayload struct {
	Errors []struct {
		Path  string `json:"path"`
		Error string `json:"error"`
		Code  string `json:"code"`
	} `json:"errors"`
}

// Get Makes a POST request
func (c Client) Get(options ClientRequestOptions, output interface{}) error {

	var errorPayload ErrorPayload

	_, err := c.newClientWithHeaders(options).Get(options.Path).Receive(&output, &errorPayload)
	if err != nil {

		log.WithFields(log.Fields{
			"err": err,
		}).Error("ERROR MAKING REQUEST")

		return err

	}

	if len(errorPayload.Errors) > 0 {
		log.WithFields(log.Fields{
			"errorPayload": errorPayload,
		}).Error("CHECK ERROR PAYLOAD")

		return errors.New(errorPayload.Errors[0].Error)
	}

	return nil
}

// NewClient returns a new Client.
func NewClient(options *ClientOptions) *Client {
	var httpClient *http.Client

	log.Error(">>>>>>>>>>>>>>>>>>>>>>>>>> NewClient")

	base := sling.New().Client(httpClient).Base(options.BaseURL)

	for _, header := range options.BaseHeaders {

		base.Add(header.Key, header.Value)

		log.WithFields(log.Fields{
			"header.Key":   header.Key,
			"header.Value": header.Value,
		}).Error("header")
	}

	return &Client{
		SlingBase: base,
	}
}
