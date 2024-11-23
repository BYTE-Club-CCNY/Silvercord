const { SlashCommandBuilder } = require('discord.js');
const { DynamoDBClient } = require("@aws-sdk/client-dynamodb");
const { get_difficulty, extractProblem } = require('../../../helper');
const { dynamoConfig } = require('../../../aws-config');
const path = require('node:path');

module.exports = {
	data: new SlashCommandBuilder()
		.setName('stats_online')
		.setDescription("View your personal online LeetCode stats")
		.addStringOption(option => option
			.setName('link')
			.setDescription('Your submission link for a solved problem')
			.setRequired(true)),
    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;
        const url = "https://leetcode.com/graphql"
        query = `
            query getUserProfile($username: String!) {
                matchedUser(username: $username) {
                    username
                    submitStats: submitStatsGlobal {
                    acSubmissionNum {
                        difficulty
                        count
                        submissions
                    }
                    }
                }
                }
        `
    }
}
