import { streamLLMRequest } from "../../grpc_client.js";

export async function handleTestGrpcCommand(interaction) {
  const userId = interaction.user.id;
  const userQuery = interaction.options.getString("query");

  try {
    await interaction.reply(`🤖 LLM is thinking...`);
    const stream = streamLLMRequest(userId, userQuery);

    let message = "";
    stream.on("data", (response) => {
      message += response.response_text + "\n";
      interaction.editReply(`🤖 LLM says: ${message}`);
    });

    stream.on("end", () => {
      interaction.editReply(`🤖 LLM says: ${message}\nStream ended.`);
    });

    stream.on("error", (err) => {
      console.error("gRPC error:", err);
      interaction.editReply("⚠️ Error contacting LLM service.");
    });
  } catch (err) {
    console.error("gRPC error:", err);
    await interaction.reply("⚠️ Error contacting LLM service.");
  }
}
