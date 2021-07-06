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
	ListService interface {
		CreateList(input *CreateUpdateListInput) (*ListOutput, error)
		GetAllLists(input *PaginatorInput) (*AllListOutput, error)
		GetList(id string) (*ListOutput, error)
		UpdateList(id string, input *CreateUpdateListInput) (*ListOutput, error)
		ArchiveList(id string) error
		SubscribeList(id string, input *SubscribeListInput) (*SubscribeListOutput, error)
		UnsubscribeList(id, uid string) error
	}

	Lists service

	ListOutput struct {
		Id              string    `json:"id"`
		Title           string    `json:"title"`
		Description     string    `json:"description"`
		SubscriberCount int       `json:"subscriber_count"`
		BroadcastCount  int       `json:"broadcast_count"`
		DoubleOptIn     bool      `json:"double_optin"`
		RedirectUrl     string    `json:"redirect_url"`
		CreatedAt       time.Time `json:"created_at"`
	}
	AllListOutput struct {
		Data       []*ListOutput `json:"data"`
		NextCursor string        `json:"next_cursor"`
		PrevCursor string        `json:"prev_cursor"`
	}
	SubscribeListOutput struct {
		Uid string `json:"uid"`
	}
)

// CreateList creates a new list using the provided input - Documentation Link: https://engage.so/docs/api/lists#create-a-list
func (l *Lists) CreateList(input *CreateUpdateListInput) (*ListOutput, error) {
	if input.Title == nil {
		return nil, errors.New("goengage: title is required")
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := l.client.newRequest(http.MethodPost, "/lists", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	var output ListOutput
	err = l.client.makeRequest(req, &output)
	if err != nil {
		return nil, err
	}

	return &output, err
}

// GetAllLists returns as array of lists - Documentation Link: https://engage.so/docs/api/lists#get-all-list-data
func (l *Lists) GetAllLists(input *PaginatorInput) (*AllListOutput, error) {

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

	req, err := l.client.newRequest(http.MethodGet, fmt.Sprintf("/lists?%v", params.Encode()), nil)
	if err != nil {
		return nil, err
	}

	var output AllListOutput
	err = l.client.makeRequest(req, &output)
	if err != nil {
		return nil, err
	}

	return &output, err
}

// GetList retrieves the details of a list using it's ID - Documentation Link: https://engage.so/docs/api/lists#get-all-list-data
func (l *Lists) GetList(id string) (*ListOutput, error) {
	if id == "" {
		return nil, errors.New("goengage: id is required")
	}

	req, err := l.client.newRequest(http.MethodGet, fmt.Sprintf("/lists/%v", id), nil)
	if err != nil {
		return nil, err
	}

	var output ListOutput
	err = l.client.makeRequest(req, &output)
	if err != nil {
		return nil, err
	}

	return &output, err
}

// UpdateList updates properties of the list - Documentation Link: https://engage.so/docs/api/lists#update-a-list
func (l *Lists) UpdateList(id string, input *CreateUpdateListInput) (*ListOutput, error) {
	if id == "" {
		return nil, errors.New("goengage: id is required")
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := l.client.newRequest(http.MethodPut, fmt.Sprintf("/lists/%v", id), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	var output ListOutput
	err = l.client.makeRequest(req, &output)
	if err != nil {
		return nil, err
	}

	return &output, err
}

// ArchiveList archives list with provided ID - Documentation Link: https://engage.so/docs/api/lists#archive-a-list
func (l *Lists) ArchiveList(id string) error {
	if id == "" {
		return errors.New("goengage: id is required")
	}

	req, err := l.client.newRequest(http.MethodDelete, fmt.Sprintf("/lists/%v", id), nil)
	if err != nil {
		return err
	}

	var output map[string]string
	return l.client.makeRequest(req, &output)
}

// SubscribeList creates a user and subscribes to a list - Documentation Link: https://engage.so/docs/api/lists#subscribe-to-a-list
func (l *Lists) SubscribeList(id string, input *SubscribeListInput) (*SubscribeListOutput, error) {
	if input.Email == nil && input.Number == nil {
		return nil, errors.New("goengage: Email or Number is required")
	}

	if id == "" {
		return nil, errors.New("goengage: id is required")
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := l.client.newRequest(http.MethodPost, fmt.Sprintf("/lists/%v/subscribers", id), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	var output SubscribeListOutput
	err = l.client.makeRequest(req, &output)
	if err != nil {
		return nil, err
	}

	return &output, err
}

// UnsubscribeList Remove subscribers from list. - Documentation Link: https://engage.so/docs/api/lists#unsubscribe-from-a-list
func (l *Lists) UnsubscribeList(id, uid string) error {
	if id == "" {
		return errors.New("goengage: id is required")
	}

	if id == "" {
		return errors.New("goengage: uid is required")
	}

	req, err := l.client.newRequest(http.MethodDelete, fmt.Sprintf("/lists/%v/subscribers/%v", id, uid), nil)
	if err != nil {
		return err
	}

	var output map[string]string
	return l.client.makeRequest(req, &output)
}
