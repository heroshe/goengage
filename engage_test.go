package goengage

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	client      *Client
	fakeService *httptest.Server
	fakeUser    = UserOutput{
		Id:           "5fc6477241fcec31a9548e98",
		Uid:          "123456789",
		UidUpdatable: false,
		FirstName:    "Heroshe",
		LastName:     "Engineering",
		Number:       "987456321",
		Email:        "someone@heroshe.com",
		Devices: []UserDevice{
			{
				Token:    "QWERTYUIOP",
				Platform: "ANDROID",
			},
			{
				Token:    "ASDFGHJKL",
				Platform: "APPLE",
			},
		},
		Lists: []UserList{
			{
				Id:         "ALL_USERS",
				Subscribed: true,
			},
			{
				Id:         "MOBILE_USERS",
				Subscribed: false,
			},
		},
		Segments: []UserSegment{
			{
				Id:         "ACTIVE_USERS",
				Suppressed: true,
			},
			{
				Id:         "CHURNED_USERS",
				Suppressed: true,
			},
		},
		Meta: map[string]interface{}{
			"property_a": "value_a",
			"property_b": "value_b",
		},
		CreatedAt: time.Now().UTC(),
	}
	fakeList = ListOutput{
		Id:              "ASDEWSDEWQQASZXDSED",
		Title:           "Waiting List",
		Description:     "Waiting List",
		SubscriberCount: 10,
		BroadcastCount:  10,
		DoubleOptIn:     true,
		RedirectUrl:     "https://test.com/landing",
		CreatedAt:       time.Now().UTC(),
	}
)

func TestMain(m *testing.M) {
	fakeService = fakeServer()
	defer fakeService.Close()

	cfg := NewConfig().WithCredentials(NewStaticCredentials("my_public_key", "my_private_key"))
	client, _ = New(cfg)
	client.BaseUrl = fakeService.URL

	os.Exit(m.Run())
}

//User Tests

func TestUsers_Create(t *testing.T) {
	assert.NotPanics(t, func() {
		user, err := client.Users.Create(&CreateUserInput{
			Id:             "1234567",
			FirstName:      nil,
			LastName:       nil,
			Email:          nil,
			Number:         nil,
			Lists:          nil,
			DeviceToken:    nil,
			DevicePlatform: nil,
			CreatedAt:      nil,
			Meta:           nil,
		})

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, user.FirstName)
		assert.NotEmpty(t, user.Email)
		assert.NotEmpty(t, user.Uid)
	})
}

func TestUsers_Get(t *testing.T) {
	assert.NotPanics(t, func() {
		user, err := client.Users.Get("123456789")
		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, fakeUser.FirstName, user.FirstName)
		assert.Equal(t, fakeUser.LastName, user.LastName)
		assert.Equal(t, fakeUser.Email, user.Email)
		assert.Equal(t, 2, len(user.Devices))
		assert.Equal(t, 2, len(user.Segments))
		assert.Equal(t, 2, len(user.Lists))
	})
}

func TestUsers_List(t *testing.T) {
	assert.NotPanics(t, func() {
		users, err := client.Users.List(&PaginatorInput{
			Limit:      nil,
			NextCursor: nil,
			PrevCursor: nil,
		})

		assert.Nil(t, err)
		assert.NotNil(t, users)
		assert.Equal(t, 3, len(users.Data))
	})
}

func TestUsers_UpdateAttributes(t *testing.T) {
	assert.NotPanics(t, func() {
		user, err := client.Users.UpdateAttributes(fakeUser.Uid, &UpdateUserAttributesInput{
			Number: String("1234567890"),
		})

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "1234567890", user.Number)
	})
}

func TestUsers_AddEvent(t *testing.T) {
	assert.NotPanics(t, func() {
		err := client.Users.AddEvent(fakeUser.Uid, &AddUserEvent{
			Event: "login",
			Value: "yes",
			Properties: map[string]interface{}{
				"user":     1,
				"browser":  "Chrome",
				"isMobile": true,
			},
			Timestamp: Time(time.Now().UTC()),
		})

		assert.Nil(t, err)
	})
}

// List Tests

func TestLists_CreateList(t *testing.T) {
	assert.NotPanics(t, func() {
		list, err := client.Lists.CreateList(&CreateUpdateListInput{
			Title:       String(fakeList.Title),
			Description: String(fakeList.Description),
			RedirectUrl: String(fakeList.RedirectUrl),
			DoubleOptIn: Bool(fakeList.DoubleOptIn),
		})

		assert.Nil(t, err)
		assert.NotNil(t, list)
		assert.NotEmpty(t, list.Title)
		assert.NotEmpty(t, list.Description)
		assert.NotEmpty(t, list.RedirectUrl)
		assert.Equal(t, fakeList.Title, list.Title)
		assert.Equal(t, fakeList.SubscriberCount, list.SubscriberCount)
		assert.Equal(t, fakeList.BroadcastCount, list.BroadcastCount)
	})
}

func TestLists_GetAllLists(t *testing.T) {
	assert.NotPanics(t, func() {
		lists, err := client.Lists.GetAllLists(&PaginatorInput{
			Limit:      nil,
			NextCursor: nil,
			PrevCursor: nil,
		})

		assert.Nil(t, err)
		assert.NotNil(t, lists)
		assert.Equal(t, 3, len(lists.Data))
	})
}

func TestLists_GetList(t *testing.T) {
	assert.NotPanics(t, func() {
		list, err := client.Lists.GetList(fakeList.Id)

		assert.Nil(t, err)
		assert.NotNil(t, list)
		assert.Equal(t, fakeList.Title, list.Title)
		assert.Equal(t, fakeList.SubscriberCount, list.SubscriberCount)
		assert.Equal(t, fakeList.BroadcastCount, list.BroadcastCount)
	})
}

func TestLists_UpdateList(t *testing.T) {
	assert.NotPanics(t, func() {
		list, err := client.Lists.UpdateList(fakeList.Id, &CreateUpdateListInput{
			Title: String("New List Name"),
		})

		assert.Nil(t, err)
		assert.NotNil(t, list)
		assert.Equal(t, "New List Name", list.Title)
	})
}

func TestLists_ArchiveList(t *testing.T) {
	assert.NotPanics(t, func() {
		err := client.Lists.ArchiveList(fakeList.Id)

		assert.Nil(t, err)
	})
}

func TestLists_UnsubscribeList(t *testing.T) {
	assert.NotPanics(t, func() {
		err := client.Lists.UnsubscribeList(fakeList.Id, fakeUser.Uid)
		assert.Nil(t, err)
	})
}

// Config Tests

func TestConfig_WithCredentials(t *testing.T) {
	cfg := NewConfig().WithCredentials(NewStaticCredentials("my_pub_key", "my_priv_key"))
	assert.Equal(t, "my_pub_key", cfg.Credentials.publicKey)
	assert.Equal(t, "my_priv_key", cfg.Credentials.privateKey)

	os.Setenv("ENGAGE_SO_PUBLIC_KEY", "my_env_pub_key")
	os.Setenv("ENGAGE_SO_PRIVATE_KEY", "my_env_priv_key")
	cfg2 := NewConfig().WithCredentials(NewEnvCredentials())

	assert.Equal(t, "my_env_pub_key", cfg2.Credentials.publicKey)
	assert.Equal(t, "my_env_priv_key", cfg2.Credentials.privateKey)
}

func TestConfig_WithHttpClient(t *testing.T) {
	d := 10 * time.Second

	cfg := NewConfig().WithHttpClient(&http.Client{
		Timeout: d,
	})
	assert.NotNil(t, cfg.HTTPClient)
	assert.Equal(t, d, cfg.HTTPClient.Timeout)
}

//StartServer initializes a test HTTP server useful for request mocking
func fakeServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		user, _ := json.Marshal(fakeUser)
		list, _ := json.Marshal(fakeList)

		users, _ := json.Marshal(ListUserOutput{
			Data: []*UserOutput{
				&fakeUser, &fakeUser, &fakeUser,
			},
		})

		lists, _ := json.Marshal(AllListOutput{
			Data: []*ListOutput{
				&fakeList, &fakeList, &fakeList,
			},
		})

		switch r.URL.Path {
		case "/users":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				fmt.Fprintf(w, string(users))

			case http.MethodPost:
				w.WriteHeader(200)
				fmt.Fprintf(w, string(user))
			}

		case fmt.Sprintf("/users/%v", fakeUser.Uid):
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				fmt.Fprintf(w, string(user))
			case http.MethodPut:
				updatedUser := fakeUser
				updatedUser.Number = "1234567890"
				updatedUserJson, _ := json.Marshal(updatedUser)
				w.WriteHeader(200)
				fmt.Fprintf(w, string(updatedUserJson))
			}

		case fmt.Sprintf("/users/%v/events", fakeUser.Uid):
			w.WriteHeader(200)
			fmt.Fprintf(w, `{"status":"ok"}`)

		case "/lists":
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				fmt.Fprintf(w, string(lists))
			case http.MethodPost:
				w.WriteHeader(200)
				fmt.Fprintf(w, string(list))
			}

		case fmt.Sprintf("/lists/%v", fakeList.Id):
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(200)
				fmt.Fprintf(w, string(list))

			case http.MethodPut:
				updatedList := fakeList
				updatedList.Title = "New List Name"
				updatedListJson, _ := json.Marshal(updatedList)
				w.WriteHeader(200)
				fmt.Fprintf(w, string(updatedListJson))

			case http.MethodDelete:
				w.WriteHeader(200)
				fmt.Fprintf(w, `{"status": "ok"}`)

			}

		case fmt.Sprintf("/lists/%v/subscribers/%v", fakeList.Id, fakeUser.Uid):
			w.WriteHeader(200)
			fmt.Fprintf(w, `{"status": "ok"}`)

		default:
			w.WriteHeader(500)
		}
	}))
}
