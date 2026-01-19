#!/usr/bin/env bash

../oto -template server.go.txt \
	-out server.gen.go \
	-pkg main \
	./def
gofmt -w server.gen.go
echo "generated server.gen.go"

../oto -template client.js.txt \
	-out client.gen.js \
	-pkg main \
	-type-map ./type-mapping.yaml \
	./def
echo "generated client.gen.js"

../oto -template openapi.yaml.txt \
	-out openapi.yaml \
	-pkg main \
	./def
echo "generated openapi.yaml"

#../oto -template client.swift.txt \
#	-out ./swift/SwiftCLIExample/SwiftCLIExample/client.gen.swift \
#	-pkg main \
#	-type-map ./type-mapping.yaml \
#	./def
#echo "generated client.gen.swift"
