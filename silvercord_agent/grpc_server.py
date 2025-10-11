from concurrent import futures

import grpc
import llm_pb2
import llm_pb2_grpc
from llm import process_query


class LLMServiceServicer(llm_pb2_grpc.LLMServiceServicer):
    def ProcessRequest(self, request, context):
        print(f"Received request from {request.user_id}: {request.query}")
        result = process_query(request.query)
        return llm_pb2.LLMResponse(
                response_text=result,
                metadata={"model": "claude", "status": "ok"}
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    llm_pb2_grpc.add_LLMServiceServicer_to_server(LLMServiceServicer(), server)
    server.add_insecure_port("[::]:50051")
    print("Python gRPC server running on port 50051")
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
