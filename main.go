package main

import (
	"log"
	"os"

	"github.com/DemmyDemon/glesys-dns-client/config"
	"github.com/DemmyDemon/glesys-dns-client/glesys"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Could not load .env file:  %s", err)
	}
	user := os.Getenv("GLESYS_USER")
	key := os.Getenv("GLESYS_KEY")
	if user == "" {
		log.Fatal("Missing GLESYS_USER environment variable")
	}
	if key == "" {
		log.Fatal("Missing GLESYS_KEY environment variable")
	}
	client := glesys.NewGlesysClient(user, key)

	ip, err := config.GetExternalIPv4()
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("External IP is %s\n", ip)
	log.SetPrefix("[" + ip + "] ")

	hosts, err := config.LoadHosts("hosts.json")
	if err != nil {
		log.Fatal("Could not load hosts.json")
	}

	client.Update(ip, hosts)
}
