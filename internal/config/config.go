package config

type Config struct {
	Frontend       Frontend
	Application    Application
	Discovery      Discovery
	Storage        Storage
	Representation Representation
	Selector       Selector
}

type Frontend struct {
	Http Http
}

type Application struct {
	AutoUpdate bool `json:"auto-update"`
}

type Http struct {
	Enabled bool
	Host    string
	Port    uint16
}

type Discovery struct {
	Identity string
	Scheme   string
	Host     string
	Port     uint16
	Path     string
	Language string
}

type Storage struct {
	Directory   string
	Limit       uint
	AddTimeout  uint `json:"add-timeout"`
	ReadTimeout uint `json:"read-timeout"`
	TTL         uint
}

type Representation struct {
	Categories struct {
		Enabled  bool
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
