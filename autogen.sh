#!/bin/bash
GIT_HASH=`git rev-parse HEAD`
COMPILE_TIME=`date -u +'%Y-%m-%d %H:%M:%S UTC'`
GIT_BRANCH=`git branch | grep "^\*" | sed 's/^..//'`
cat <<EOF >version.go
package main

var hash="$GIT_HASH"
var branch="$GIT_BRANCH"
var compileTime="$COMPILE_TIME"

EOF