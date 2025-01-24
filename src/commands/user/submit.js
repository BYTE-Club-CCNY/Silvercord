const { SlashCommandBuilder } = require('discord.js');
const { get_difficulty, extractProblem } = require('../../../helper');
const { get_score, update_score, add_problem, get_problems } = require('../../../dynamo_helper');
const {dynamoConfig} = require("../../../aws-config");
const {DynamoDBClient, QueryCommand} = require("@aws-sdk/client-dynamodb");
const client = new DynamoDBClient(dynamoConfig);
const path = require('node:path');

module.exports = {
	data: new SlashCommandBuilder()
		.setName('submit')
		.setDescription("Submit a LeetCode problem to the server's leaderboard!")
		.addStringOption(option => option
			.setName('link')
			.setDescription('Your submission link for a solved problem')
			.setRequired(true)),
	async execute(interaction) {
		if (!interaction.isChatInputCommand()) return;
		const user_link = interaction.options.getString('link') ?? 'No link provided'
		await interaction.deferReply();
		let nextUserScore;
		let nextUserID;
		try {
			is_valid = true // change this when we implement validate_page
			if (!is_valid) {
				interaction.followUp(`Link ${user_link} is an invalid submission`);
				return;
			}
			const server_id = interaction.guild.id;
			const user_id = interaction.user.id;
			const table = "leetboard";
			const table_scores = "leetboard_scores"
			// verifying link manually
			const problem_name = extractProblem(user_link);
			if (problem_name === null) {
				interaction.followUp(`Link ${user_link} is in the wrong format, resubmit the link that contains the submission ID`);
				return;
			}
			const problems = await get_problems(server_id, user_id, table);
			if (problems.includes(problem_name)) {
				interaction.followUp(`You've already done this problem. Nice try.`)
				return;
			}
			// calculate difficulty
			const difficulty = await get_difficulty(user_link);
			score = 0
			if (difficulty === "Easy") {
				score = 1
			} else if (difficulty === "Medium") {
				score = 2
			} else {
				score = 4
			}
			// below is to get scores to compare, see if user has now beat the next user
			const query = {
				"TableName": table_scores,
				"KeyConditionExpression": "server_id = :server_id",
				"ExpressionAttributeValues": {
					":server_id": {S: server_id}
				},
				"ProjectionExpression": "score, user_id"
			};
			const command = new QueryCommand(query);
			const response = await client.send(command);
			const items = response.Items || [];
			const prev_score = await get_score(server_id, user_id, table_scores);
			const final_score = prev_score + parseInt(score, 10);
            try {
                if (items && items.length > 0) {
                    const sorted_users = await Promise.all(
                            items.map(item => {
                                return [item.user_id.S, parseInt(item.score.N, 10)];
                                })
                            );
                    sorted_users.sort((a, b) => b[1] - a[1]);
                    console.log("sorted users", sorted_users);
                    if (sorted_users.length > 1) {
                        for (let i = 0; i < sorted_users.length-1; i++) {
                            if (sorted_users[i + 1][1] === prev_score) {
                                nextUserScore = sorted_users[i][1]
                                    nextUserID = sorted_users[i][0]
                                    console.log('next user score & id: ', nextUserScore, nextUserID)
                                    break
                            }
                        }
                    }
                } else {
                    nextUserScore = prev_score
                }
            } catch (error) {
                console.log(error);
                interaction.followUp('Internal error with data retrieval');
            }

            // update the score:
            await update_score(server_id, user_id, final_score, table_scores);
            await add_problem(server_id, user_id, user_link, problem_name, table);
            console.log(`prev score: ${prev_score}, new score: ${final_score}`)
                if (prev_score <= nextUserScore < final_score) {
                    try {
                        console.log('fetching with ID:', nextUserID);
                        const nextUser = await interaction.client.users.fetch(nextUserID);
                        interaction.followUp(`${difficulty} problem submitted. \n${interaction.user.username} has taken <@${nextUser}>'s place with new score ${final_score}!`);
                    } catch (error) {
                        console.error('error fetching username', error);
                        interaction.followUp('Error with username fetching');
                    }
                } else {
                    interaction.followUp(`${difficulty} problem submitted. \n${interaction.user.username}'s score is now ${final_score}!`);
                }
        } catch (error) {
            console.log(error);
            interaction.followUp('Could not submit problem, please try again');
        }
    }
}
