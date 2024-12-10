package config

type Config struct {
	Frontend       Frontend
	Application    Application
	Api            Api
	Discovery      Discovery
	Storage        Storage
	Representation Representation
	Selector       Selector
}

type Frontend struct {
	Http     Http
	Telegram Telegram
}

type Application struct {
	AutoUpdate bool `json:"auto-update"`
}

type Http struct {
	Enabled bool
	Host    string
	Port    uint16
}

type Telegram struct {
	Enabled bool
	ApiPath string `json:"api-path"`
}

type Api struct {
	Scheme string
	Host   string
	Port   uint16
	Domain string
}

type Discovery struct {
	ApiPath  string `json:"api-path"`
	Language string
}

type Storage struct {
	Directory   string
	Driver      string
	Limit       uint
	AddTimeout  uint `json:"add-timeout"`
	ReadTimeout uint `json:"read-timeout"`
	TTL         uint
}

type Representation struct {
	Enabled    bool
	Directory  string
	Categories struct {
		Type     bool
		Alphabet bool
		Genres   bool
		Year     bool
	}
}
type Selector struct {
	Criterion           string
	MinSeasonSize       uint `json:"min-season-size"`
	MaxSeasonSize       uint `json:"max-season-size"`
	MinSeedersThreshold uint `json:"min-seeders-threshold"`
	Quality             []string
	Voices              [][]string
}
