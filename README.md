# Begger
This GoLang package facilitates to make HTTP request with retry options.


## Why Is It Named So?
In Bangladesh, the beggers ask to people for some money or some other goods. If the people don't respond, the beggers ask for the same thing repeatedly for multiple times.

This package hehaves exactly like those beggers. It requests to some service for some data to process. If the service fails to process the request, it retries the same request.


## How To Use
### Example 1
Lets say you need to create an order for an user. Assume the following example data for requests and responses.

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MhmdRyhn/begger"
)

type OrderItem struct {
	Id       int32 `json:"id"`
	Quantity int32 `json:"quantity"`
	Price    int32 `json:"price"`
}

type CreateOrderRequest struct {
	Items []OrderItem `json:"items"`
}

type CreateOrderResponse struct {
	OrderId int32  `json:"order_id"`
	Status  string `json:"status"`
}

type OrderSummary struct {
	OrderId   int32  `json:"order_id"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type OrderDetailsResponse struct {
	Orders []OrderSummary `json:"orders"`
}

func main() {
	userId := 1267
	url := fmt.Sprintf("https://my-example-shop.com/api/v1/users/%d/orders", userId)
	order := CreateOrderRequest{
		Items: []OrderItem{
			{
				Id:       1,
				Quantity: 2,
				Price:    33500,
			},
		},
	}
	body, err := json.Marshal(order)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req := begger.Request{
		Client: &http.Client{Timeout: 1 * time.Second},
		Components: begger.RequestComponents{
			Url:        begger.Url{Actual: &url},
			HTTPMethod: http.MethodPost,
			Body:       body,
			Headers: begger.Headers{
				"Api-Key": "132mrn34tb9193qnje43t5ijr",
			},
		},
	}
	resp, err2 := req.Do()
	if err2 != nil {
		fmt.Println(err2.Message)
		return
	}
	// Because, the package is not responsible for closing response body
	defer resp.Body.Close()
	respBody := CreateOrderResponse{}
	parser := begger.NewResponseParser(resp)
	fmt.Println("Status code: ", parser.HTTPStatusCode())
	if parser.HTTPStatusCode() == http.StatusOK {
		parser.LoadBody(&respBody)
		fmt.Println(fmt.Sprintf("Body: %+v", respBody))
	}
}
```

### Example 2
We want to get all PENDING orders of an user. Lets use the above data structures here too.
```go
func main() {
    // Desired URL is: `https://my-example-shop.com/api/v1/users/1267/orders?status=PENDING`
    urlPathFormat := "/api/v1/users/{UserId}/orders"
	req := begger.Request{
		Client: &http.Client{Timeout: 1 * time.Second},
		Components: begger.RequestComponents{
			Url: begger.Url{
				Components: &begger.UrlComponents{
					Host:        "https://my-example-shop.com",
					PathFormat:  urlPathFormat,
					PathParams:  begger.PathParams{"{UserId}": "1267"},
					QueryParams: begger.QueryParams{"status": "PENDING"},
				},
			},
			HTTPMethod: http.MethodGet,
			Headers: begger.Headers{
				"Api-Key": "132mrn34tb9193qnje43t5ijr",
			},
		},
	}
	resp, err2 := req.Do()
	if err2 != nil {
		fmt.Println(err2.Message)
		return
	}
	// Because, the package is not responsible for closing response body
	defer resp.Body.Close()
	respBody := OrderDetailsResponse{}
	parser := begger.NewResponseParser(resp)
	fmt.Println("Status code: ", parser.HTTPStatusCode())
	if parser.HTTPStatusCode() == http.StatusOK {
		parser.LoadBody(&respBody)
		fmt.Println(fmt.Sprintf("Body: %+v", respBody))
	}
}
```
### Example 3
Lets use some retry.
```go
func main() {
    // Desired URL is: `https://my-example-shop.com/api/v1/users/1267/orders?status=PENDING`
    urlPathFormat := "/api/v1/users/{UserId}/orders"
	req := begger.Request{
		Client: &http.Client{Timeout: 1 * time.Second},
		Components: begger.RequestComponents{
			Url: begger.Url{
				Components: &begger.UrlComponents{
					Host:        "https://my-example-shop.com",
					PathFormat:  urlPathFormat,
					PathParams:  begger.PathParams{"{UserId}": "1267"},
					QueryParams: begger.QueryParams{"status": "PENDING"},
				},
			},
			HTTPMethod: http.MethodGet,
			Headers: begger.Headers{
				"Api-Key": "132mrn34tb9193qnje43t5ijr",
			},
		},
		Retry: &begger.RetryOptions{
			MaxAttempt: 3, WaitInterval: 2 * time.Second,
		},
	}
	resp, err2 := req.Do()
	if err2 != nil {
		fmt.Println(err2.Message)
		return
	}
	// Because, the package is not responsible for closing response body
	defer resp.Body.Close()
	respBody := OrderDetailsResponse{}
	parser := begger.NewResponseParser(resp)
	fmt.Println("Status code: ", parser.HTTPStatusCode())
	if parser.HTTPStatusCode() == http.StatusOK {
		parser.LoadBody(&respBody)
		fmt.Println(fmt.Sprintf("Body: %+v", respBody))
	}
}
```
