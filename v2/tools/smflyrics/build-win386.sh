#!/bin/bash
GOOS=windows GOARCH=386 CGO_ENABLED=0 exec go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o artifacts/smflyrics.exe