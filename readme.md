# goggles

Read content in your terminal while simulating GoogleBot user agent.

## Install

`go get github.com/jacoduplessis/goggles`

## Usage

`goggles [-selector CSS] [-timeout DURATION] [-ua STRING] URL`

Examples:

```
goggles https://example.com
goggles -selector h1 https://example.com
goggles -timeout 30s -selector "article p" https://news.ycombinator.com
```

