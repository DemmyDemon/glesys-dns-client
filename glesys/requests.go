package glesys

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type GlesysClient struct {
	User   string
	Key    string
	Client *http.Client
}

// NewGlesysClient creates a brand new GlesysClient with the given credentials.
func NewGlesysClient(user string, key string) GlesysClient {
	return GlesysClient{
		User:   user,
		Key:    key,
		Client: &http.Client{},
	}
}

// buildRequest builds the request will all the headers and stuff.
func (glesys GlesysClient) buildRequest(url string, rawData url.Values) (*http.Request, error) {
	data := bytes.NewBuffer([]byte(rawData.Encode()))
	rq, err := http.NewRequest("POST", url, data)
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

	responseData := GlesysResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	return responseData, err
}

// ListRecords fetches a list of records from GleSYS
func (glesys GlesysClient) ListRecords(domain string) (GlesysResponse, error) {

	rawData := url.Values{"domainname": {domain}}

	rq, err := glesys.buildRequest("https://api.glesys.com/domain/listrecords/format/json", rawData)
	if err != nil {
		return GlesysResponse{}, err
	}
	return glesys.doRequest(rq)
}

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

// UpdateRecord attempts to update the given record to the new IP
func (glesys GlesysClient) UpdateRecord(recordId int, newIP string) (GlesysResponse, error) {

	rawData := url.Values{"recordid": {strconv.Itoa(recordId)}, "data": {newIP}}

	rq, err := glesys.buildRequest("https://api.glesys.com/domain/updaterecord/format/json", rawData)
	if err != nil {
		return GlesysResponse{}, err
	}
	return glesys.doRequest(rq)
}