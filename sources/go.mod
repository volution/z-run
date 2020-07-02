module github.com/cipriancraciun/z-run

go 1.14

require (
	github.com/Pallinder/go-randomdata v1.2.0 // indirect
	github.com/colinmarc/cdb v0.0.0-20190223170904-60f317823f70
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/junegunn/fzf v0.0.0-20200630121719-d1676776aaf1
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12
	github.com/mattn/go-shellwords v1.0.10 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/peterh/liner v1.2.0
	github.com/saracen/walker v0.1.1 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/vbauerster/mpb/v5 v5.2.2
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208 // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae
	golang.org/x/text v0.3.3 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace (
	github.com/colinmarc/cdb => github.com/cipriancraciun/go-cdb-lib v0.0.0-20190809203657-d959ce9cc674
	github.com/junegunn/fzf => github.com/cipriancraciun/fzf v0.0.0-20200411153254-524c512952bc
)
