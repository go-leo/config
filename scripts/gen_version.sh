#!/bin/bash

VERSION=$(git describe --tags --abbrev=0)
cat <<EOF > cmd/protoc-gen-leo-config/version.go
package main

const Version = "$VERSION"
EOF

git add cmd/protoc-gen-leo-config/version.go
git commit -m "chore: update version.go for $VERSION"