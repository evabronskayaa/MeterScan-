if not exist "protoc.exe" (
    echo "Download protoc.exe"
    powershell -Command "Invoke-WebRequest https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.24.4/protoc-3.24.4-windows-x86_64.exe -OutFile protoc.exe"
)
protoc --go_out=../backend/internal/proto --go_opt=paths=source_relative --go-grpc_out=../backend/internal/proto --go-grpc_opt=paths=source_relative *.proto