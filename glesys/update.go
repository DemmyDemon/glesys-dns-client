package glesys

import (
	"fmt"
	"log"
	"sync"

	"github.com/DemmyDemon/glesys-dns-client/config"
)

// Update goes through the given records and make sure they match the given IP
func (glesys GlesysClient) Update(ip string, hosts []config.GlesysHost) {
	if len(hosts) == 0 {
		log.Println("Zero hosts specified: Skipping update")
		return
	}
	wg := sync.WaitGroup{}
	for _, host := range hosts {
		wg.Add(1)
		go func(domain string, subdomains []string) {
			defer wg.Done()
			records, err := glesys.RecordMap(domain, "A")
			if err != nil {
				log.Printf("Failed to get records for %s: %s\n", domain, err)
				return
			}
			for _, subdomain := range subdomains {
				status := "Already up to date"
				record, ok := records[subdomain]
				if ok {
					if record.Data != ip {
						_, err := glesys.UpdateRecord(record.Recordid, ip)
						if err != nil {
							status = err.Error()
						} else {
							status = "Updated"
						}
					}
				} else {
					status = "No record found!"
				}
				log.Printf("%15s.%-15s -> %s", subdomain, domain, status)
			}
		}(host.Domain, host.Subdomains)
	}
	wg.Wait()
}

// Certbot updates or creates the _acme-challenge TXT record for the given domain
func (glesys GlesysClient) Certbot(domain string, validation string) error {

	records, err := glesys.RecordMap(domain, "TXT")
	if err != nil {
		return fmt.Errorf("certbot get records: %w", err)
	}
	if record, ok := records["_acme-challenge"]; ok {
		if record.Data == validation {
			log.Printf("Certbot challenge for %s is unchanged: Skipping update\n", domain)
			return nil
		}
		_, err := glesys.UpdateRecord(record.Recordid, validation)
		if err != nil {
			return fmt.Errorf("certbot update record: %w", err)
		}
		log.Printf("Updated Certbot challenge for %s\n", domain)
		return nil
	}

	_, err = glesys.CreateRecord(domain, "TXT", "_acme-challenge", validation)
	if err != nil {
		return fmt.Errorf("certbot create record: %w", err)
	}
	log.Printf("Created Certbot challenge for %s\n", domain)
	return nil
}
