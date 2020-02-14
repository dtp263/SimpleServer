package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_echo"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//Donation - test structure for use with json:api
type Donation struct {
	Name  string  `jsonapi:"attr,name"`
	Value float32 `jsonapi:"attr,value"`
}

func main() {
	e := echo.New()

	limiter := tollbooth.NewLimiter(10, nil)
	limiter.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(tollbooth_echo.LimitHandler(limiter))

	e.GET("/", handleGet)
	e.POST("/json", handleJSON)
	e.POST("/donation_jsonapi", handleDonationJSONAPI)

	e.Logger.Fatal(e.Start("127.0.0.1:8000"))
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// handleGet - is your basic GET endpoint, it returns "Default Get"
func handleGet(c echo.Context) error {
	return c.String(http.StatusOK, "Default Get")
}

// handleJSON - will parse a JSON message into a Go struct, then
// convert it back to a string, and return it in the Response.
func handleJSON(c echo.Context) error {
	body := c.Request().Body

	wholeBody, err := ioutil.ReadAll(body)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusNotAcceptable,
			"Please provide valid Request Body")
	}

	message := map[string]interface{}{}
	err = json.Unmarshal([]byte(wholeBody), &message)
	if err != nil {
		if !isJSON(string(wholeBody)) {
			return echo.NewHTTPError(
				http.StatusNotAcceptable,
				"Please provide valid JSON")
		}
		return echo.NewHTTPError(
			http.StatusNotAcceptable,
			`Please provide simple an not deeply nested JSON unless there is an 
			endpoint specified for unmarshalling a specific data struct`)
	}

	return c.JSON(http.StatusOK, message)
}

// handleDonationJSONAPI - will parse a Request with the json:api format using the Donation
// struct. Then it will convert it back to a string and return it in the Response.
func handleDonationJSONAPI(c echo.Context) error {
	body := c.Request().Body

	donation := new(Donation)
	err := jsonapi.UnmarshalPayload(body, donation)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusNotAcceptable,
			"Unable to parse json:api Request")
	}

	err = jsonapi.MarshalPayload(c.Response(), donation)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusNotAcceptable,
			"Unable to parse Donation object into json:api for Response")
	}

	return c.JSON(http.StatusOK, c.Response())
}
