@echo off
setlocal

:: Save current directory
set CURRENT_DIR=%cd%

:: Move up to parent directory
cd ..

:: Clone enproto-protobuf if it does not exist
if not exist "enproto-protobuf" (
    echo Cloning enproto-protobuf...
    git clone https://github.com/enproto/enproto-protobuf
)

:: Return to original directory
cd "%CURRENT_DIR%"

:: Compile enproto.proto using protoc
protoc --go_out=pb --go_opt=paths=source_relative ^
       --proto_path=../enproto-protobuf ^
       enproto.proto

protoc --go_out=pb --go_opt=paths=source_relative ^
       --proto_path=../enproto-protobuf ^
       session/rsa.proto

protoc --go_out=pb --go_opt=paths=source_relative ^
       --proto_path=../enproto-protobuf ^
       session/aes.proto
       
protoc --go_out=pb --go_opt=paths=source_relative ^
       --proto_path=../enproto-protobuf ^
       common/types.proto
endlocal
pause
