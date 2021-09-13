#!/bin/bash

counter=1
while [ $counter -le 10 ]
do
    echo $counter
    go clean -testcache && go test -gcflags=all=-d=checkptr -v -coverpkg ./... ./test/concurrency_test/... -coverprofile cover.out -testify.m Test
    ((counter++))
done