const { SlashCommandBuilder, EmbedBuilder } = require('discord.js');
const { get_online_username } = require('../../../helper');
const { get_username } = require('../../../dynamo_helper');
// const path = require('node:path');

module.exports = {
	data: new SlashCommandBuilder()
		.setName('stats_online')
		.setDescription("View your personal online LeetCode stats"),

    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;
        server_id = interaction.guild.id;
        user_id = interaction.user.id;
        await interaction.deferReply();
        try {
            const table = "leetboard"
            const username = await get_username(server_id, user_id, table);
            if (username === null) {
                interaction.followUp('You are not registered! Do /register to view your online LeetCode stats!');
                return ;
            }
            const data = await get_online_username(username);
            console.log(data);
            if (!data || !data.submitStats || !Array.isArray(data.submitStats.acSubmissionNum)) {
                await interaction.followUp(
                    `Could not retrieve ${username} LeetCode stats. Please make sure your username is correct!`
                );
                return;
            }
            const embed = new EmbedBuilder()
                .setTitle(`${username}'s Leetcode Stats`)
                .setColor("Random");

            console.log("Stats data:", data.submitStats.acSubmissionNum);
            data.submitStats.acSubmissionNum.forEach(difficulty => {
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
            interaction.followUp('Could not get stats, try again please');
        }
    }
}
