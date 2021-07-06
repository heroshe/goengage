package goengage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type (
	UserService interface {
		Create(input *CreateUserInput) (*UserOutput, error)
		Get(uid string) (*UserOutput, error)
		List(input *PaginatorInput) (*ListUserOutput, error)
		UpdateAttributes(uid string, input *UpdateUserAttributesInput) (*UserOutput, error)
		AddEvent(uid string, event *AddUserEvent) error
	}

	Users service

	UserOutput struct {
		Id           string                 `json:"id"`
		Uid          string                 `json:"uid"`
		UidUpdatable bool                   `json:"uid_updateable"`
		FirstName    string                 `json:"first_name"`
		LastName     string                 `json:"last_name"`
		Number       string                 `json:"number"`
		Email        string                 `json:"email"`
		Devices      []UserDevice           `json:"devices"`
		Lists        []UserList             `json:"lists"`
		Segments     []UserSegment          `json:"segments"`
		Meta         map[string]interface{} `json:"meta"`
		CreatedAt    time.Time              `json:"created_at"`
	}

	UserDevice struct {
		Token    string `json:"token"`
		Platform string `json:"platform"`
	}

	UserList struct {
		Id         string `json:"id"`
		Subscribed bool   `json:"subscribed"`
	}

	UserSegment struct {
		Id         string `json:"id"`
		Suppressed bool   `json:"suppressed"`
	}

	ListUserOutput struct {
		Data       []*UserOutput `json:"data"`
		NextCursor string        `json:"next_cursor"`
		PrevCursor string        `json:"prev_cursor"`
	}
)

// Create create a new user - Documentation Link: https://engage.so/docs/api/users#create-a-user
func (u *Users) Create(input *CreateUserInput) (*UserOutput, error) {
	if input.Id == "" {
		return nil, errors.New("goengage: id is required")
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := u.client.newRequest(http.MethodPost, "/users", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	var output UserOutput
	err = u.client.makeRequest(req, &output)
	if err != nil {
		return nil, err
	}

	return &output, err
}

// Get fetches and returns a user's profile - Documentation Link: https://engage.so/docs/api/users#retrieve-a-user
func (u *Users) Get(uid string) (*UserOutput, error) {
	if uid == "" {
		return nil, errors.New("goengage: uid is required")
	}

	req, err := u.client.newRequest(http.MethodGet, fmt.Sprintf("/users/%v", uid), nil)
	if err != nil {
		return nil, err
	}

	var output UserOutput
	err = u.client.makeRequest(req, &output)
	if err != nil {
		return nil, err
	}

	return &output, err
}

// List returns a list of users. - Documentation Link: https://engage.so/docs/api/users#list-users
func (u *Users) List(input *PaginatorInput) (*ListUserOutput, error) {

	if input.NextCursor != nil && input.PrevCursor != nil {
		return nil, errors.New("goengage: Cannot use Next and Prev cursor at the same time")
	}

	params := url.Values{}
	if input.Limit != nil {
		params.Add("limit", strconv.Itoa(*input.Limit))
	}

	if input.NextCursor != nil && input.PrevCursor == nil {
		params.Add("next_cursor", *input.NextCursor)
	}

	if input.PrevCursor != nil && input.NextCursor == nil {
		params.Add("prev_cursor", *input.PrevCursor)
	}

	req, err := u.client.newRequest(http.MethodGet, fmt.Sprintf("/users?%v", params.Encode()), nil)
	if err != nil {
		return nil, err
	}

	var output ListUserOutput
	err = u.client.makeRequest(req, &output)
	if err != nil {
		return nil, err
	}

	return &output, err
}

// UpdateAttributes updates user data and attributes. - Documentation Link: https://engage.so/docs/api/users#update-user-attributes
func (u *Users) UpdateAttributes(uid string, input *UpdateUserAttributesInput) (*UserOutput, error) {
	if uid == "" {
		return nil, errors.New("goengage: uid is required")
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := u.client.newRequest(http.MethodPut, fmt.Sprintf("/users/%v", uid), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	var output UserOutput
	err = u.client.makeRequest(req, &output)
	if err != nil {
		return nil, err
	}

	return &output, err
}

// AddEvent Add user events. It returns an error if any or nil if operation successful. Successful == 200 status code
// Documentation Link: https://engage.so/docs/api/users#add-user-events
func (u *Users) AddEvent(uid string, event *AddUserEvent) error {
	if uid == "" {
		return errors.New("goengage: uid is required")
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err

	}

	req, err := u.client.newRequest(http.MethodPut, fmt.Sprintf("/users/%v/events", uid), bytes.NewReader(payload))
	if err != nil {
		return err
	}

	var output map[string]string
	return u.client.makeRequest(req, &output)
}
