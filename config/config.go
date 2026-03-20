package config

type Config struct {
	Env        string
	WebAddr    string `yaml:"web_addr" default:":8080"`
	AppBaseURL string `yaml:"app_base_url"` // required
	PrivateKey string `yaml:"private_key"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string
		Host     string `default:"localhost"`
		Port     string `default:"3306"`
		Debug    bool
		SSLMode  string `yaml:"sslmode"`
	}
}
