api:
goctl api go --api *.api --dir ../. 

rpc:
goctl rpc protoc *.proto --go_out=../ --go-grpc_out=../  --zrpc_out=../ --style=goZero

rpc(分组):
goctl rpc protoc *.proto --go_out=../ --go-grpc_out=../  --zrpc_out=../ --style=goZero -m

model:
goctl model mysql ddl --cachhe=true --src --dir 
