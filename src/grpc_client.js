// src/grpc_client.js
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path = require("path");

const PROTO_PATH = path.join(__dirname, "../proto/llm.proto");

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const llmProto = grpc.loadPackageDefinition(packageDefinition).llm;

const client = new llmProto.LLMService(
  "localhost:50051",
  grpc.credentials.createInsecure(),
);

function processLLMRequest(userId, command, query) {
  return new Promise((resolve, reject) => {
    client.ProcessRequest(
      { user_id: userId, command, query },
      (err, response) => {
        if (err) reject(err);
        else resolve(response);
      },
    );
  });
}

function streamLLMRequest(userId, query) {
  return client.StreamRequest({ user_id: userId, query });
}

module.exports = {
  processLLMRequest,
  streamLLMRequest,
};
