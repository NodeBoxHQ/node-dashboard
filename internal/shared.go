package internal

type NodeboxConfig struct {
	Environment  string `json:"environment"`
	DataPath     string `json:"dataPath"`
	IP           string `json:"ip"`
	LineaIP      string `json:"lineaIP"`
	DuskPassword string `json:"duskPassword"`
	Port         int    `json:"port"`
	LogLevel     string `json:"logLevel"`
}
