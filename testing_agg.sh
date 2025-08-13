#!/bin/bash

./down_up.sh

go run . reset

go run . register Adam
go run . login Adam

go run . addfeed "techcrunch" "https://techcrunch.com/feed"
go run . addfeed "ycombinator" "https://news.ycombinator.com/rss"
go run . addfeed "boot blog" "https://blog.boot.dev/index.xml"


go run . agg 20s
exit 0
