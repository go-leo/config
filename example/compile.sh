#!/bin/bash

protoc \
--proto_path=. \
--proto_path=../proto \
--go_out=. \
--go_opt=paths=source_relative \
--config_out=. \
--config_opt=paths=source_relative \
*/*.proto
