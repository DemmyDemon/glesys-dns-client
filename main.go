package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/DemmyDemon/glesys-dns-client/config"
	"github.com/DemmyDemon/glesys-dns-client/glesys"
)

func main() {

	ip, err := config.GetExternalIPv4()
	if err != nil {
		log.Fatalf("Could not determine external IPv4: %s\n", err)
	}
	log.SetOutput(os.Stdout)
	log.SetPrefix("[" + ip + "] ")

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not determine homedir: %s\n", err)
	}
	var certBotClean string
	var configFile string
	flag.StringVar(&configFile, "cfg", homedir+"/.glesys_dns.json", "Path to where the Configuration is stored.")
	flag.StringVar(&ip, "ip", ip, "IP address to check against, and update to.")
	flag.StringVar(&certBotClean, "certbotclean", "", "Domain to clean up certbot validations for")
	flag.Parse()

	configFile = filepath.Clean(configFile)

	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("Could not load configuration file: %s", err)
	}

	client := glesys.NewGlesysClient(cfg.Credentials.User, cfg.Credentials.Key)

	if certBotClean != "" {
		err := client.CertbotCleanup(certBotClean)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0) // Because if we're cleaning, we are *just* cleaning.
	}

	certbotDomain := os.Getenv("CERTBOT_DOMAIN")
	if certbotDomain != "" {
		// That is, we're in CERTBOT mode...
		certbotValidation := os.Getenv("CERTBOT_VALIDATION")
		if certbotValidation == "" {
			certbotValidation = "dry-run"
		}
		err := client.Certbot(certbotDomain, certbotValidation)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Will idle for 60 seconds for TXT records to settle.")
		time.Sleep(60 * time.Second)
		os.Exit(0) // DO NOT do the updates if we're in certbot mode.
	}

	client.Update(ip, cfg.Hosts)
}
