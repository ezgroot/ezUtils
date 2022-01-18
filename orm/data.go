package orm

// SQL Config of db.
type Config struct {
	SQLName     string `json:"sqlName"`
	Account     string `json:"account"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	Database    string `json:"database"`
	MaxIdleConn int    `json:"maxIdleConn"`
	MaxOpenConn int    `json:"maxOpenConn"`
	MaxLifetime int    `json:"maxLIfetime"` // unit - second
	IsTrace     bool   `json:"isTrace"`
}
