module github.com/cipriancraciun/z-run

go 1.12

require (
	github.com/Pallinder/go-randomdata v1.2.0 // indirect
	github.com/colinmarc/cdb v0.0.0-00000000000000-000000000000
	github.com/junegunn/fzf v0.0.0-00000000000000-000000000000
	github.com/mattn/go-isatty v0.0.10
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mattn/go-shellwords v1.0.6 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
)

replace (
	github.com/colinmarc/cdb => github.com/cipriancraciun/go-cdb-lib v0.0.0-20190809203657-d959ce9cc674
	github.com/junegunn/fzf => github.com/cipriancraciun/fzf v0.0.0-20190804152507-0f31e2393c3c
)
