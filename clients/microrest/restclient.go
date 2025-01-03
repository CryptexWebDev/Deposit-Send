package microrest

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func NewRestClient(restBaseUrl string) *RestClient {
	return &RestClient{
		restBaseUrl: restBaseUrl,
		http:        &http.Client{},
	}
}

type RestClient struct {
	restBaseUrl string
	http        *http.Client
}

func (r *RestClient) Get(req string, responseObject interface{}) (err error) {
	reqUrl, _ := url.JoinPath(r.restBaseUrl, req)
	httpRequest, err := http.NewRequest("GET", reqUrl, nil)
	httpRequest.Header.Set("Content-Type", "application/json")
	httpResponse, err := r.http.Do(httpRequest)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusInternalServerError || httpResponse.StatusCode == http.StatusBadRequest {
		return errors.New("invalid server response: " + httpResponse.Status)
	}
	return json.NewDecoder(httpResponse.Body).Decode(responseObject)
}
