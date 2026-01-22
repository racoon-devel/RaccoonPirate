package config

type Config struct {
	Frontend       Frontend       `json:"frontend"`
	Application    Application    `json:"application"`
	Api            Api            `json:"api"`
	Discovery      Discovery      `json:"discovery"`
	Database       Database       `json:"database"`
	Torrent        Torrent        `json:"torrent"`
	Representation Representation `json:"representation"`
	Selector       Selector       `json:"selector"`
}

type Frontend struct {
	Http     Http     `json:"http"`
	Telegram Telegram `json:"telegram"`
}

type Application struct {
	AutoUpdate    bool `json:"auto-update"`
	ConfigVersion uint `json:"config-version"`
}

type Http struct {
	Enabled bool   `json:"enabled"`
	Host    string `json:"host"`
	Port    uint16 `json:"port"`
}

type Telegram struct {
	Enabled bool   `json:"enabled"`
	ApiPath string `json:"api-path"`
}

type Api struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Port   uint16 `json:"port"`
	Domain string `json:"domain"`
}

type Discovery struct {
	ApiPath  string `json:"api-path"`
	Language string `json:"language"`
}

type Database struct {
	Path   string `json:"path"`
	Driver string `json:"driver"`
}

type Torrent struct {
	Driver     string     `json:"driver"`
	Builtin    Builtin    `json:"builtin"`
	TorrServer TorrServer `json:"torr-server"`
}

type Builtin struct {
	Directory   string `json:"directory"`
	Limit       uint   `json:"limit"`
	AddTimeout  uint   `json:"add-timeout"`
	ReadTimeout uint   `json:"read-timeout"`
	TTL         uint   `json:"ttl"`
}

type TorrServer struct {
	URL      string `json:"url"`
	Fusepath string `json:"fusepath"`
}

type Representation struct {
	Enabled    bool   `json:"enabled"`
	Directory  string `json:"directory"`
	Categories struct {
		Type     bool `json:"type"`
		Alphabet bool `json:"alphabet"`
		Genres   bool `json:"genres"`
		Year     bool `json:"year"`
	} `json:"categories"`
}
type Selector struct {
	Criterion           string     `json:"criterion"`
	MinSeasonSize       uint       `json:"min-season-size"`
	MaxSeasonSize       uint       `json:"max-season-size"`
	MinSeedersThreshold uint       `json:"min-seeders-threshold"`
	Quality             []string   `json:"quality"`
	Voices              [][]string `json:"voices"`
}
