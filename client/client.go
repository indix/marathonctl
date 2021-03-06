package client

import (
	"encoding/json"
	"errors"

	"github.com/parnurzeal/gorequest"
)

type Marathon struct {
	Url string
}

type ErrorResponse struct {
	Details []struct {
		Errors []string `json:"errors"`
		Path   string   `json:"path"`
	} `json:"details"`
	Message string `json:"message"`
}

var httpClient = gorequest.New()

func handle(response gorequest.Response, body string, errs []error) (string, error) {
	if response != nil {
		if (response.StatusCode != 200 && response.StatusCode != 201) && body != "" {
			var errorResponse ErrorResponse
			err := json.Unmarshal([]byte(body), &errorResponse)
			if err != nil {
				errs = append(errs, err)
			} else {
				if response.StatusCode == 422 {
					errs = append(errs, errors.New(errorResponse.Details[0].Errors[0]))
				} else {
					errs = append(errs, errors.New(errorResponse.Message))
				}
			}
		} else if response.StatusCode != 200 && response.StatusCode != 201 {
			errs = append(errs, errors.New(response.Status))
		}
	}
	return body, combineErrors(errs)
}

func combineErrors(errs []error) error {
	if len(errs) == 1 {
		return errs[0]
	} else if len(errs) > 1 {
		msg := "Error(s):"
		for _, err := range errs {
			msg += " " + err.Error()
		}
		return errors.New(msg)
	} else {
		return nil
	}
}
