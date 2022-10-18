// Everlasting gratitude to https://mholt.github.io/json-to-go/
package glesys

import "time"

type GlesysResponse struct {
	Response Response `json:"response"`
}
type Status struct {
	Code      int       `json:"code"`
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
}
type Record struct {
	Recordid   int    `json:"recordid"`
	Domainname string `json:"domainname"`
	Host       string `json:"host"`
	Type       string `json:"type"`
	Data       string `json:"data"`
	TTL        int    `json:"ttl"`
}
type Input struct {
	Domainname         string `json:"domainname"`
	Format             string `json:"format"`
	Projectkey         string `json:"projectkey"`
	Organizationnumber int    `json:"organizationnumber"`
	Recordid           string `json:"recordid"`
	Data               string `json:"data"`
}
type Debug struct {
	Input Input `json:"input"`
}
type Response struct {
	Status  Status   `json:"status"`
	Records []Record `json:"records"`
	Record  Record   `json:"record"`
	Debug   Debug    `json:"debug"`
}
