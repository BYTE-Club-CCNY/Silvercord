const { SlashCommandBuilder, EmbedBuilder } = require('discord.js');
require('dotenv').config();

const BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';

module.exports = {
    data: new SlashCommandBuilder()
        .setName('leaderboard')
        .setDescription("View the server's LeetCode leaderboard!"),
    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;
        await interaction.deferReply();
        try {
            const server_id = interaction.guild.id;
            const response = await fetch(`${BASE_URL}/scores/leaderboard?server_id=${server_id}`);
            
            if (!response.ok) {
                throw new Error('Failed to fetch leaderboard');
            }

            const scores = await response.json();
            
            if (scores.length === 0) {
                await interaction.followUp('No leaderboard data to be retrieved'); 
                return;
            }

            const userResults = await Promise.all(
                scores.map(async item => {
                    try {
                        const user = await interaction.client.users.fetch(item.user_id);
                        return [user.username, item.score];
                    } catch (error) {
                        console.error(`Failed to fetch user ${item.user_id}:`, error.message);
                        return null;
                    }
                })
            );
            const sorted_users = userResults.filter(user => user !== null);
            
            if (sorted_users.length === 0) {
                await interaction.followUp('No valid users found on the leaderboard.');
                return;
            }
            
            let leaderboard_text = ""
            for (let i=0; i<sorted_users.length; i++) {
                if (i===0) {
                    leaderboard_text += `🌟${sorted_users[i][0]} - ${sorted_users[i][1]}\n`;
                } else if (i===1) {
                    leaderboard_text += `⭐${sorted_users[i][0]} - ${sorted_users[i][1]}\n`;
                } else if (i===2) {
                    leaderboard_text += `✨${sorted_users[i][0]} - ${sorted_users[i][1]}\n`;
                } else {
                    leaderboard_text += `${sorted_users[i][0]} - ${sorted_users[i][1]}\n`;
                }
            }
            const embed = new EmbedBuilder()
                .setColor('#0099ff')
                .setTitle("Highest LeetCode Scores (Season 3, Oct 1 - Dec 31)")
                .setDescription(leaderboard_text)

            console.log(sorted_users);
            interaction.followUp({embeds: [ embed ]});
        } catch (error) {
            console.log(error);
            interaction.followUp('Unable to retrieve leaderboard, please try again');
        }
    }
}
