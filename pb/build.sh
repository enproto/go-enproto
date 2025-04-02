#!/bin/bash

# Save current directory
current=$(pwd)

# Go up one directory
cd ..

# Clone proto repo if not exists
if [ ! -d "enproto-protobuf" ]; then
    git clone https://github.com/enproto/enproto-protobuf
fi

# Return to original directory
cd "$current"

# Compile proto
protoc --go_out=pb --go_opt=paths=source_relative \
       --proto_path=../enproto-protobuf \
       enproto.proto

protoc --go_out=pb --go_opt=paths=source_relative \
       --proto_path=../enproto-protobuf \
       session/rsa.proto

protoc --go_out=pb --go_opt=paths=source_relative \
       --proto_path=../enproto-protobuf \
       session/aes.proto
       
protoc --go_out=pb --go_opt=paths=source_relative \
       --proto_path=../enproto-protobuf \
       common/types.proto