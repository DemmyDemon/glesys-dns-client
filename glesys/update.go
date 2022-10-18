package glesys

import (
	"log"
	"sync"
)

func (glesys GlesysClient) Update(ip string, hosts map[string][]string) {
	wg := sync.WaitGroup{}
	for currentdomain, currentsubdomains := range hosts {
		wg.Add(1)
		go func(domain string, subdomains []string) {
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
			wg.Done()
		}(currentdomain, currentsubdomains)
	}
	wg.Wait()
}
