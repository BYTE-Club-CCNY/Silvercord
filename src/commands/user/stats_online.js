const { SlashCommandBuilder, EmbedBuilder } = require("discord.js");
require("dotenv").config();

const API_BASE_URL = process.env.API_BASE_URL || "http://api:8080";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("stats_online")
    .setDescription("View your personal online LeetCode stats"),

  async execute(interaction) {
    if (!interaction.isChatInputCommand()) return;
    const server_id = interaction.guild.id;
    const user_id = interaction.user.id;
    await interaction.deferReply();
    try {
      const usernameResponse = await fetch(
        `${API_BASE_URL}/users/username?server_id=${server_id}&user_id=${user_id}`,
      );

      if (!usernameResponse.ok) {
        interaction.followUp(
          "Failed to fetch your username. Please try again.",
        );
        return;
      }

      const usernameData = await usernameResponse.json();
      const username = usernameData.username;

      if (!username || username === "") {
        interaction.followUp(
          "You are not registered! Do /register to view your online LeetCode stats!",
        );
        return;
      }

      const statsResponse = await fetch(
        `${API_BASE_URL}/leetcode/user?username=${username}`,
      );

      if (!statsResponse.ok) {
        await interaction.followUp(
          `Could not retrieve ${username} LeetCode stats. Please make sure your username is correct!`,
        );
        return;
      }

      const data = await statsResponse.json();
      console.log(data);
      if (
        !data ||
        !data.submitStats ||
        !Array.isArray(data.submitStats.acSubmissionNum)
      ) {
        await interaction.followUp(
          `Could not retrieve ${username} LeetCode stats. Please make sure your username is correct!`,
        );
        return;
      }
      const embed = new EmbedBuilder()
        .setTitle(`${username}'s Leetcode Stats`)
        .setColor("Random");

      console.log("Stats data:", data.submitStats.acSubmissionNum);
      data.submitStats.acSubmissionNum.forEach((difficulty) => {
        const difficultyName = difficulty.difficulty || "Unknown Difficulty";
        const submissions = difficulty.submissions?.toString() || "0";
        const count = difficulty.count?.toString() || "0";

        if (difficultyName && submissions && count) {
          embed.addFields({
            name: `${difficultyName} problems`,
            value: `Submissions: ${submissions}, Count: ${count}`,
            inline: false,
          });
        } else {
          console.log("Skipped invalid difficulty:", difficulty);
        }
      });
      await interaction.followUp({ embeds: [embed] });
    } catch (error) {
      console.log(error);
      interaction.followUp("Could not get stats, try again please");
    }
  },
};
