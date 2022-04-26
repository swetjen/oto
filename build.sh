#!/bin/bash

go build || exit 1
cp oto ~/go/bin/ || exit 1