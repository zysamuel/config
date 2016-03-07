#!/bin/bash
cd $SR_CODE_BASE/snaproute/src/config/docgen/
go run *.go
mv $SR_CODE_BASE/snaproute/src/config/docgen/flexApis.json $SR_CODE_BASE/snaproute/src/config/docsui/index.html
