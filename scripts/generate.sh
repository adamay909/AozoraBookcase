#!/usr/bin/bash

TEMPL='//Generated file. Do not edit as it will be overwritten.

package main

const revHash = "REV"
' 

echo "$TEMPL" > rev.go

HASH=$(git rev-parse --short HEAD)
sed -i "s/REV/$HASH/" rev.go
