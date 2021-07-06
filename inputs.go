package goengage

import "time"

type (
	CreateUserInput struct {
		Id             string                 `json:"id,omitempty"`
		FirstName      *string                `json:"first_name,omitempty"`
		LastName       *string                `json:"last_name,omitempty"`
		Email          *string                `json:"email,omitempty"`
		Number         *string                `json:"number,omitempty"`
		Lists          []string               `json:"lists,omitempty"`
		DeviceToken    *string                `json:"device_token,omitempty"`
		DevicePlatform *string                `json:"device_platform,omitempty"`
		CreatedAt      *time.Time             `json:"created_at,omitempty"`
		Meta           map[string]interface{} `json:"meta,omitempty"`
	}

	UpdateUserAttributesInput struct {
		FirstName      *string                `json:"first_name,omitempty"`
		LastName       *string                `json:"last_name,omitempty"`
		Email          *string                `json:"email,omitempty"`
		Number         *string                `json:"number,omitempty"`
		Lists          []string               `json:"lists,omitempty"`
		DeviceToken    *string                `json:"device_token,omitempty"`
		DevicePlatform *string                `json:"device_platform,omitempty"`
		CreatedAt      *time.Time             `json:"created_at,omitempty"`
		Meta           map[string]interface{} `json:"meta,omitempty"`
	}

	PaginatorInput struct {
		Limit      *int    `json:"limit"`
		NextCursor *string `json:"next_cursor"`
		PrevCursor *string `json:"prev_cursor"`
	}

	AddUserEvent struct {
		Event      string                 `json:"event,omitempty"`
		Value      interface{}            `json:"value,omitempty"`
		Properties map[string]interface{} `json:"properties,omitempty"`
		Timestamp  *time.Time             `json:"timestamp,omitempty"`
	}

	CreateUpdateListInput struct {
		Title       *string `json:"title,omitempty"`
		Description *string `json:"description,omitempty"`
		RedirectUrl *string `json:"redirect_url,omitempty"`
		DoubleOptIn *bool   `json:"double_optin,omitempty"`
	}

	SubscribeListInput struct {
		FirstName *string                `json:"first_name,omitempty"`
		LastName  *string                `json:"last_name,omitempty"`
		Email     *string                `json:"email,omitempty"`
		Number    *string                `json:"number,omitempty"`
		CreatedAt *time.Time             `json:"created_at,omitempty"`
		Meta      map[string]interface{} `json:"meta,omitempty"`
	}
)

func String(input string) *string {
	return &input
}

func Int(input int) *int {
	return &input
}

func Bool(input bool) *bool {
	return &input
}

func Time(input time.Time) *time.Time {
	return &input
}
