module github.com/gustav1105/epg_approved

go 1.23.4

require (
	github.com/pelletier/go-toml v1.9.5
	github.com/spf13/cobra v1.8.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/gustav1105/epg_approved/internal/api => ./internal/api

replace github.com/gustav1105/epg_approved/internal/config => ./internal/config

replace github.com/gustav1105/epg_approved/internal/handlers => ./internal/handlers
