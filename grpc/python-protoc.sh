# shellcheck disable=SC2035
python -m grpc_tools.protoc -I . --python_out=../ml/proto --pyi_out=../ml/proto --grpc_python_out=../ml/proto *.proto