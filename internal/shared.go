package internal

type NodeboxConfig struct {
	Environment string `json:"environment"`
	DataPath    string `json:"dataPath"`
	IP          string `json:"ip"`
	LineaIP     string `json:"lineaIP"`
	Port        int    `json:"port"`
	LogLevel    string `json:"logLevel"`
	GithubToken string `json:"githubToken"`
}
