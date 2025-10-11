import { processLLMRequest } from "../../grpc_client.js";

export async function handleTestGrpcCommand(interaction) {
  const userId = interaction.user.id;
  const userQuery = interaction.options.getString("query");

  try {
    const response = await processLLMRequest(userId, userQuery);
    await interaction.reply(`🤖 LLM says: ${response.response_text}`);
  } catch (err) {
    console.error("gRPC error:", err);
    await interaction.reply("⚠️ Error contacting LLM service.");
  }
}
