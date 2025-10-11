const { SlashCommandBuilder, EmbedBuilder } = require('discord.js');

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';

module.exports = {
    data: new SlashCommandBuilder()
        .setName('problems')
        .setDescription("See the problems you've completed."),
        
    async execute(interaction) {
      const server_id = interaction.guild.id;
      const user_id = interaction.user.id;
      
      const response = await fetch(`${API_BASE_URL}/problems?server_id=${server_id}&user_id=${user_id}`);
      
      if (!response.ok) {
        await interaction.reply("Failed to fetch your problems. Please try again.");
        return;
      }

      const problemsData = await response.json();
      
      if (!problemsData || problemsData.length === 0) {
        await interaction.reply("You haven't completed any problems yet!");
        return;
      }

      const problems = problemsData.map(p => p.problem);
      
      const embed = new EmbedBuilder()
        .setTitle("Problems")
        .setDescription(`List of problems done by ${interaction.user.username}`)
        .addFields({
          name: "Problems", value: problems.join("\n")
        })
      await interaction.reply({embeds: [embed]});
    },
};
