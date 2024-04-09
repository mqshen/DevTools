package types

type DevLocation struct {
	Country  string `json:"country" yaml:"country"`
	Province string `json:"province" yaml:"province"`
	City     string `json:"city" yaml:"city"`
	ISP      string `json:"isp" yaml:"isp"`
	ASN      string `json:"asn" yaml:"asn"`
}

type DevIP struct {
	Source   string       `json:"source" yaml:"source"`
	IP       string       `json:"ip" yaml:"ip"`
	Location *DevLocation `json:"location" yaml:"location"`
}
