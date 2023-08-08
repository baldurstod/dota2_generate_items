`go build dota2_generate_items`

linux

`SET GOOS=linux&&SET GOARCH=amd64&&go build  dota2_generate_items`

# Usage
`dota2_generate_items -o <outputdir> -i <itemsdir> -r <resourcedir> -l <lang>`

itemsdir is usually <DOTA2 INSTALL DIR>/scripts/items

resourcedir is usually <DOTA2 INSTALL DIR>/resource/localization

lang can be any supported dota2 language
