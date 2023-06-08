package glesys

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type GlesysClient struct {
	User   string
	Key    string
	Base   string
	Client *http.Client
}

// NewGlesysClient creates a brand new GlesysClient with the given credentials.
func NewGlesysClient(user string, key string) GlesysClient {
	return GlesysClient{
		User:   user,
		Key:    key,
		Base:   "https://api.glesys.com",
		Client: &http.Client{},
	}
}

// buildRequest builds the request will all the headers and stuff.
func (glesys GlesysClient) buildRequest(url string, rawData url.Values) (*http.Request, error) {
	data := bytes.NewBuffer([]byte(rawData.Encode()))
	rq, err := http.NewRequest("POST", glesys.Base+url, data)
	if err != nil {
		return nil, err
	}
	rq.SetBasicAuth(glesys.User, glesys.Key)
	rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return rq, nil
}

// doRequest actually does the request.
func (glesys GlesysClient) doRequest(rq *http.Request) (GlesysResponse, error) {
	response, err := glesys.Client.Do(rq)
	if err != nil {
		return GlesysResponse{}, err
	}

	// Treat anything not in the 200 range as an error
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return GlesysResponse{}, errors.New(response.Status)
	}

	responseData := GlesysResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	return responseData, err
}

// ListRecords fetches a list of records from GleSYS
func (glesys GlesysClient) ListRecords(domain string) (GlesysResponse, error) {

	rawData := url.Values{"domainname": {domain}}

	rq, err := glesys.buildRequest("/domain/listrecords/format/json", rawData)
	if err != nil {
		return GlesysResponse{}, err
	}
	return glesys.doRequest(rq)
}

// RecordMap returns just the actual records from ListRecords, in a map keyed on the `Host` value of the record.
func (glesys GlesysClient) RecordMap(domain string, recordType string) (map[string]Record, error) {
	records := make(map[string]Record)
	response, err := glesys.ListRecords(domain)
	if err != nil {
		return records, err
	}
	for _, record := range response.Response.Records {
		if record.Type == recordType {
			records[record.Host] = record
		}
	}
	return records, nil
}

// CreateRecord attempts to create a record with the given parameters
func (glesys GlesysClient) CreateRecord(domain string, recordType string, recordName string, recordValue string) (GlesysResponse, error) {
	rawData := url.Values{
		"domainname": {domain},
		"host":       {recordName},
		"type":       {recordType},
		"data":       {recordValue},
		"ttl":        {"60"}, // Short TTL to avoid oopsed run DNS cache issues
	}

	rq, err := glesys.buildRequest("/domain/addrecord/format/json", rawData)
	if err != nil {
		return GlesysResponse{}, err
	}
	return glesys.doRequest(rq)
}

// UpdateRecord attempts to update the given record to the new value
func (glesys GlesysClient) UpdateRecord(recordId int, newValue string) (GlesysResponse, error) {

	rawData := url.Values{"recordid": {strconv.Itoa(recordId)}, "data": {newValue}}

	rq, err := glesys.buildRequest("/domain/updaterecord/format/json", rawData)
	if err != nil {
		return GlesysResponse{}, err
	}
	return glesys.doRequest(rq)
}
