package glesys

import (
	"log"
	"sync"

	"github.com/DemmyDemon/glesys-dns-client/config"
)

// Update goes through the given records and make sure they match the given IP
func (glesys GlesysClient) Update(ip string, hosts []config.GlesysHost) {
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
