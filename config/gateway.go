package config

type Gateway struct {
	Name      string      `json:"name,omitempty"`
	Version   string      `json:"version,omitempty"`
	Host      string      `json:"host,omitempty"`
	Endpoints []*Endpoint `json:"endpoints,omitempty"`
	Plugins   []*Plugin   `json:"plugins,omitempty"`
}

type Endpoint struct {
	Pattern     string     `json:"pattern,omitempty"`
	Method      string     `json:"method,omitempty"`
	Description string     `json:"description,omitempty"`
	Protocol    string     `json:"protocol,omitempty"`
	Timeout     int        `json:"timeout,omitempty"`
	Plugins     []*Plugin  `json:"plugins,omitempty"`
	Backends    []*Backend `json:"backends,omitempty"`
}

type Backend struct {
	Target string `json:"target,omitempty"`
	Weight int    `json:"weight,omitempty"`
}

type Plugin struct {
	Name    string      `json:"name,omitempty"`
	Options interface{} `json:"options,omitempty"`
}
