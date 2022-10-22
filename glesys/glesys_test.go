package glesys_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const listRecordsPath = "/domain/listrecords/format/json"
const listRecordsResponse = `{"response":{"status":{"code":200,"timestamp":"2022-10-22T10:35:07+02:00","text":"OK"},"records":[{"recordid":2468,"domainname":"example.com","host":"test","type":"A","data":"192.168.1.1","ttl":120}],"debug":{"input":{"domainname":"example.com","format":"json","projectkey":"Testing","organizationnumber":1234}}}}`

const updateRecordPath = "/domain/updaterecord/format/json"
const updateRecordResponse = `{"response":{"status":{"code":200,"timestamp":"2022-10-22T12:05:50+02:00","text":"OK"},"record":{"recordid":2468,"domainname":"example.com","host":"test","type":"A","data":"127.0.0.1","ttl":120},"debug":{"input":{"recordid":"2779066","data":"192.168.1.7","format":"json","projectkey":"Testing","organizationnumber":1234}}}}`

func authOK(username string, password string, ok bool) bool {
	if !ok {
		return false
	}
	if username != "Testing" {
		return false
	}
	if password != "Testing" {
		return false
	}
	return true
}

func stubServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !authOK(r.BasicAuth()) {
				t.Errorf("Test setup problem: Credentals are wrong. Use basic auth Testing/Testing")
			}
			err := r.ParseForm()
			if err != nil {
				t.Errorf("Error encountered parsing POST data: %s", err)
			}
			var response string
			switch r.URL.Path {
			case listRecordsPath:
				response, err = listRecordsStub(r)
				if err != nil {
					t.Errorf("Error listing records: %s", err)
				}
			case updateRecordPath:
				response, err = updateRecordStub(r)
				if err != nil {
					t.Errorf("Error updating records: %s", err)
				}
			default:
				t.Errorf("Unexpected request path: %q", r.URL.Path)
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		}),
	)
}

func listRecordsStub(r *http.Request) (string, error) {
	if r.PostFormValue("domainname") != "example.com" {
		return "", fmt.Errorf(`Unexpected domainname POSTed. Expected "example.com", got %q`, r.PostFormValue("domainname"))
	}
	return listRecordsResponse, nil
}

func updateRecordStub(r *http.Request) (string, error) {
	if r.PostFormValue("recordid") != "2468" {
		return "", fmt.Errorf(`Unexpected recordid POSTed. Expected "2468", got %q`, r.PostFormValue("recordid"))
	}
	if r.PostFormValue("data") != "127.0.0.1" {
		return "", fmt.Errorf(`Unexpected data POSTed. Expected "127.0.0.1", got %q`, r.PostFormValue("data"))
	}
	return updateRecordResponse, nil
}
