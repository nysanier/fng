package pkgconf

// fng_daily#common配置项
type Config struct {
	Env               string `json:"env"`
	Version           string `json:"version"`
	DnsUpdateInterval int    `json:"dns_update_interval"`
}

// fng_daily#dns_config配置项
//type DnsConfig struct {
//	UpdateInterval int `json:"update_interval"`
//}
