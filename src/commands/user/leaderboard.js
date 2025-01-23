const { SlashCommandBuilder, EmbedBuilder } = require('discord.js');
const { DynamoDBClient, QueryCommand } = require("@aws-sdk/client-dynamodb");
const { dynamoConfig } = require('../../../aws-config');
const client = new DynamoDBClient(dynamoConfig);
const path = require('node:path');

module.exports = {
    data: new SlashCommandBuilder()
        .setName('leaderboard')
        .setDescription("View the server's LeetCode leaderboard!"),
    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;
        await interaction.deferReply();
        try {
            const server_id = interaction.guild.id;
            const table_scores = "leetboard_scores"
            const query = {
                "TableName": table_scores,
                "KeyConditionExpression": "server_id = :server_id",
                "ExpressionAttributeValues": {
                    ":server_id": { S: server_id }
            },
                "ProjectionExpression": "score, user_id"
            };
            const command = new QueryCommand(query);
            const response = await client.send(command);
            const items = response.Items || [];
            if (items.length === 0) {
                await interaction.followUp('No leaderboard data to be retrieved'); return;
            }
            const sorted_users = await Promise.all(
                items.map(async item => {
                    const user = await interaction.client.users.fetch(item.user_id.S);
                    return [user.username, parseInt(item.score.N, 10)];
                })
            );
            sorted_users.sort((a, b) => b[1] - a[1]);
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
            embed = new EmbedBuilder()
                .setColor('#0099ff')
                .setTitle("Highest LeetCode Scores (Season 2, Jan23 - May24)")
                .setDescription(leaderboard_text)

            console.log(sorted_users);
            interaction.followUp({embeds: [ embed ]});
        } catch (error) {
            console.log(error);
            interaction.followUp('Unable to retrieve leaderboard, please try again');
        }
    }
}
