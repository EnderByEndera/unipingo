from concurrent import futures
import time
import grpc
# from .protocols.
from .protocols.msg_pb2_grpc import GreeterServicer, add_GreeterServicer_to_server
from .protocols.msg_pb2 import HelloReply, HelloRequest


class Servicer(GreeterServicer):
    def SayHello(self, request: HelloRequest, context):
        print(request.name)
        return HelloReply(message="Hi!!!!")

def serve():
    # gRPC 服务器
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_GreeterServicer_to_server(Servicer(), server)
    server.add_insecure_port('[::]:50051')
    print("sever is opening ,waiting for message...")
    server.start()  # start() 不会阻塞，如果运行时你的代码没有其它的事情可做，你可能需要循环等待。
    try:
        while True:
            time.sleep(100000)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    serve()