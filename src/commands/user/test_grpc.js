const { streamLLMRequest } = require("../../grpc_client");
const { execute } = require("./leaderboard");
const { SlashCommandBuilder } = require("discord.js");
require("dotenv").config();

module.exports = {
  data: new SlashCommandBuilder()
    .setName("test_grpc")
    .setDescription(
      "Test to see if the grpc endpoint is properly working or not",
    ),
  async execute(interaction) {
    const userId = interaction.user.id;
    const userQuery = interaction.options.getString("query");

    try {
      await interaction.reply(`LLM is thinking...`);
      const stream = streamLLMRequest(userId, userQuery);

      let message = "";
      stream.on("data", (response) => {
        message += response.response_text + "\n";
        interaction.editReply(`LLM says: ${message}`);
      });

      stream.on("end", () => {
        interaction.editReply(`LLM says: ${message}\nStream ended.`);
      });

      stream.on("error", (err) => {
        console.error("gRPC error:", err);
        interaction.editReply("Error contacting LLM service.");
      });
    } catch (err) {
      console.error("gRPC error:", err);
      await interaction.reply("Error contacting LLM service.");
    }
  },
};
