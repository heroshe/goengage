# GoEngage

GoEngage is a Golang API wrapper for [engage.so's v1 API](https://engage.so).
This package implements the User and List resources as covered in the [v1 api docs](https://engage.so/docs/)

## Install
```
$ go get github.com/heroshe/goengage
```

## Usage
It's quite simple and straightforward. The package supports two configuration styles.


### Config and Credentials

You can read credentials from environment variables or provide them as arguments.
* EnvCredentials: This looks for the following environment variables `ENGAGE_SO_PUBLIC_KEY` and `ENGAGE_SO_PRIVATE_KEY`. 
  If the keys aren't found or, they contain blank values, no error will be returned at this step. An error
  will be returned when truing to initialize a new client using this config
  
* StaticCredentials: Sets the public and private keys based on provided arguments.

```go
package main

import (
	"github.com/heroshe/goengage"
	"log"
	"net/http"
	"time"
)

func main() {
	// Reading Credentials From Environment Variables
	cfg := goengage.NewConfig().WithCredentials(goengage.NewEnvCredentials())

	// Using Provided Credentials
	cfg = goengage.NewConfig().WithCredentials(goengage.NewStaticCredentials("your_public_key", "your_private_key"))

	// Optionally: Use your own *http.Client
	cfg = goengage.NewConfig().WithCredentials(goengage.NewEnvCredentials()).WithHttpClient(&http.Client{
		Timeout: 5 * time.Second,
	})

	// Crete a new client
	client, err := goengage.New(cfg)
	if err != nil {
		//handle error
	}

	users, err := client.Users.List(&goengage.PaginatorInput{
		Limit: goengage.Int(50),
	})
	if err != nil {
		// handle error
	}

	// Do something with users.Data
	log.Print(users.Data)
}

```

## Resources
All resources are interfaces. That means you can create mocks or fake resources that can be used for testing.

### Users
The following endpoints are supported on the user resource. Documentation Link: https://engage.so/docs/api/users
1. `Create()`: Creates a new user
2. `Get()`: Retrieves a single user
3. `List()`: Retrieves a list of users
4. `UpdateAttributes()`: Updates user data and attributes
5. `AddEvent()`: Add user events

### Lists
The following endpoints are supported on the list resource. Documentation Link: https://engage.so/docs/api/lists
1. `CreateList()`: creates a new
2. `GetAllLists()`: returns as array of lists
3. `GetList()`: retrieves the details of a list
4. `UpdateList()`: updates properties of the list
5. `ArchiveList()`:  archives list
6. `SubscribeList()`: creates a user and subscribes to a list
7. `UnsubscribeList()`: Remove subscribers from list

## Integration Testing
The resources in package are both interfaces which mean you can create your custom client struct that have fake implementation
of the resources.

## Run Tests
go test --race -cover -coverprofile=cover.out -v ./...

## Contributing
Contributors and contributions are welcome. Open and issue or PR :)
