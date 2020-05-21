module github.com/cipriancraciun/z-run

go 1.14

require (
	github.com/Pallinder/go-randomdata v1.2.0 // indirect
	github.com/colinmarc/cdb v0.0.0-00000000000000-000000000000
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/junegunn/fzf v0.0.0-00000000000000-000000000000
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12
	github.com/mattn/go-shellwords v1.0.10 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/peterh/liner v1.2.0
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/vbauerster/mpb/v5 v5.2.2
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37 // indirect
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a // indirect
	golang.org/x/sys v0.0.0-20200519105757-fe76b779f299
	golang.org/x/tools v0.0.0-20200521155704-91d71f6c2f04 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace (
	github.com/colinmarc/cdb => github.com/cipriancraciun/go-cdb-lib v0.0.0-20190809203657-d959ce9cc674
	github.com/junegunn/fzf => github.com/cipriancraciun/fzf v0.0.0-20200411153254-524c512952bc
	github.com/vbauerster/mpb/v5 => github.com/vbauerster/mpb/v5 v5.2.2-0.20200521170959-dc22c1ba4542
)
