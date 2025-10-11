// src/grpc_client.js
import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const PROTO_PATH = path.join(__dirname, "../proto/llm.proto");

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const llmProto = grpc.loadPackageDefinition(packageDefinition).silvercord;

const client = new llmProto.LLMService(
  "localhost:50051",
  grpc.credentials.createInsecure(),
);

export function processLLMRequest(userId, query) {
  return new Promise((resolve, reject) => {
    client.ProcessRequest({ user_id: userId, query }, (err, response) => {
      if (err) {
        reject(err);
      } else {
        resolve(response);
      }
    });
  });
}
