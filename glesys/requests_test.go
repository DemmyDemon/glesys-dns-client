package glesys_test

import (
	"testing"

	"github.com/DemmyDemon/glesys-dns-client/glesys"
)

func TestListRecords(t *testing.T) {
	server := stubServer(t)
	defer server.Close()
	client := glesys.NewGlesysClient("Testing", "Testing")
	client.Base = server.URL
	response, err := client.ListRecords("example.com")
	if err != nil {
		t.Errorf("Unexpected error during request: %s", err)
	}
	if len(response.Response.Records) != 1 {
		t.Errorf("Unexpected number of records in response. Expected 1, Got %d", len(response.Response.Records))
	}
	if response.Response.Records[0].Host != "test" {
		t.Errorf("Unexpected response format")
	}
}

func TestRecordMap(t *testing.T) {
	server := stubServer(t)
	defer server.Close()
	client := glesys.NewGlesysClient("Testing", "Testing")
	client.Base = server.URL
	records, err := client.RecordMap("example.com", "A")
	if err != nil {
		t.Errorf("Unexpected error during request: %s", err)
	}
	if len(records) != 1 {
		t.Errorf("Unexpected number of records. Expected 1, got %d", len(records))
	}
	record, ok := records["test"]
	if !ok {
		t.Error("Test record somehow missing from response")
	}
	if record.Data != "192.168.1.1" {
		t.Error("Test record came back formatted weirdly")
	}
}

// TODO: TestRecordMapMultiple

func TestUpdateRecord(t *testing.T) {
	server := stubServer(t)
	defer server.Close()
	client := glesys.NewGlesysClient("Testing", "Testing")
	client.Base = server.URL
	_, err := client.UpdateRecord(2468, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
}

// TODO: TestDeleteRecord
