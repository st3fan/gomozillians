// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package mozillians

import (
	"code.google.com/p/go.net/context"
	"encoding/json"
	"net/http"
	"strings"
)

type Client struct {
	app string
	key string
}

type Meta struct {
	Limit      int `json:"limit"`
	Next       int `json:"next"`
	Offset     int `json:"offset"`
	Previous   int `json:"previous"`
	TotalCount int `json:"total_count"`
}

type User struct {
	Username       string `json:"username"`
	Photo          string `json:"photo"`
	PhotoThumbnail string `json:"photo_thumbnail"`
	FullName       string `json:"full_name"`
}

type UsersResponse struct {
	Meta  Meta   `json:"meta"`
	Users []User `json:"objects"`
}

func NewClient(app, key string) *Client {
	return &Client{app: app, key: key}
}

type UsersOptions struct {
	Email   string
	Country string
	Region  string
	Groups  []string
	Skills  []string
}

func (c *Client) Users(ctx context.Context, options UsersOptions) ([]User, Meta, error) {
	req, err := http.NewRequest("GET", "https://mozillians.org/api/v1/users/", nil)
	if err != nil {
		return nil, Meta{}, err
	}

	q := req.URL.Query()
	q.Set("app_name", c.app)
	q.Set("app_key", c.key)

	if options.Email != "" {
		q.Set("email", options.Email)
	}

	if options.Country != "" {
		q.Set("country", options.Country)
	}

	if len(options.Groups) != 0 {
		q.Set("groups", strings.Join(options.Groups, ","))
	}

	if len(options.Skills) != 0 {
		q.Set("skills", strings.Join(options.Skills, ","))
	}

	req.URL.RawQuery = q.Encode()

	var users []User
	var meta Meta

	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var usersResponse UsersResponse
		if err := json.NewDecoder(resp.Body).Decode(&usersResponse); err != nil {
			return err
		}

		users = usersResponse.Users
		meta = usersResponse.Meta

		return nil
	})

	return users, meta, err
}

// Run the HTTP request in a goroutine and pass the response to f.
func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}
