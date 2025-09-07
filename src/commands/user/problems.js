const { SlashCommandBuilder, EmbedBuilder } = require('discord.js');
const { get_problems } = require('../../../api/dynamo_helper.js');

module.exports = {
    data: new SlashCommandBuilder()
        .setName('problems')
        .setDescription("See the problems you've completed."),
        
    async execute(interaction) {
      const server_id = interaction.guild.id;
      const user_id = interaction.user.id;
      const table = "leetboard"
      const problems = await get_problems(server_id, user_id, table);
      if (!problems || problems.length == 0) {
	await interaction.reply("You haven't completed any problems yet!");
      }
      const embed = new EmbedBuilder()
	.setTitle("Problems")
	.setDescription(`List of problems done by ${interaction.user.username}`)
	.addFields({
	  name: "Problems", value: problems.join("\n")
	})
      await interaction.reply({embeds: [embed]});
    },
};
