python3 -m grpc_tools.protoc -Irpc/protocols --python_out=rpc/protocols --grpc_python_out=rpc/protocols --pyi_out=rpc/protocols rpc/protocols/*.proto 
sed -i 's/^import .*_pb2 as/from . \0/' rpc/protocols/*.py 

# sed -i 's/^import .*_pb2 as/from . \0/' rpc/protocols/service_pb2_grpc.py 
# sed -i 's/^import .*_pb2 as/from . \0/' rpc/protocols/scheduler_pb2_grpc.py 