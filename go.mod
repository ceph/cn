module cn

go 1.15

require (
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137
	github.com/apcera/termtables v0.0.0-20170405184538-bcbc5dc54055
	github.com/ceph/cn v2.3.1+incompatible
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.17+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/elgs/gojq v0.0.0-20201120033525-b5293fef2759
	github.com/elgs/gosplitargs v0.0.0-20161028071935-a491c5eeb3c8 // indirect
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/moby/term v0.0.0-20220808134915-39b0c02b01ae // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.12.0
	github.com/stretchr/testify v1.8.0
	golang.org/x/crypto v0.0.0-20220817201139-bc19a97f63c8
	golang.org/x/sys v0.0.0-20220818161305-2296e01440c6
	gotest.tools/v3 v3.3.0 // indirect
)

replace github.com/ceph/cn => ./
