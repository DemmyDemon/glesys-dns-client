package glesys_test

import (
	"testing"

	"github.com/DemmyDemon/glesys-dns-client/config"
	"github.com/DemmyDemon/glesys-dns-client/glesys"
)

func TestUpdate(t *testing.T) {
	ip := "127.0.0.1"
	hosts := []config.GlesysHost{
		{
			Domain:     "example.com",
			Subdomains: []string{"test"},
		},
	}

	server := stubServer(t)
	defer server.Close()

	client := glesys.NewGlesysClient("Testing", "Testing")
	client.Base = server.URL

	client.Update(ip, hosts)
}
