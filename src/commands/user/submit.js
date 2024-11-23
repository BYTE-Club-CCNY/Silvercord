const { SlashCommandBuilder } = require('discord.js');
const { DynamoDBClient, GetItemCommand, PutItemCommand, UpdateItemCommand } = require("@aws-sdk/client-dynamodb");
const { get_difficulty, extractProblem } = require('../../../helper');
const { dynamoConfig } = require('../../../aws-config');
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
		try {
			is_valid = true // change this when we implement validate_page
			if (!is_valid) {
				interaction.followUp(`Link ${user_link} is an invalid submission`);
				return;
			}
			// verifying link manually
			const problem_name = extractProblem(user_link)
			if (problem_name === null) {
				interaction.followUp(`Link ${user_link} is in the wrong format, resubmit the link that contains the submission ID`);
				return;
			}
			// calculate difficulty
			const difficulty = await get_difficulty(user_link)
			score = 0
			if (difficulty === "Easy") {
				score = 1
			} else if (difficulty === "Medium") {
				score = 2
			} else {
				score = 4
			}
			// below retrieves data from dynamoDB; table name, partition key, and sort key are all required in the "query"
			const client = new DynamoDBClient(dynamoConfig);
			const input = {
				"Key": {
					"server_id": {
						"S": interaction.guild.id
					},
					"user_id": {
						"S": interaction.user.id
					}
				},
				"ExpressionAttributeNames": {
        			"#links": "link",
        			"#problems": "problem"
    			},
				"UpdateExpression": `
					SET 
            #links = list_append(if_not_exists(#links, :empty_list), :new_link),
            #problems = list_append(if_not_exists(#problems, :empty_list), :new_problem)
				`,
    			"ExpressionAttributeValues": {
        			":new_link": { L: [{ S: user_link }] }, 
        			":new_problem": { L: [{ S: problem_name }] }, 
        			":empty_list": { L: [] } 
    			},

				"TableName": "leetboard"
			};
			const get_score = {
				"Key": {
					"server_id": {
						"S": interaction.guild.id
					},
					"user_id": {
						"S": interaction.user.id
					},
				},
				"TableName": "leetboard_scores"
			};
			const get_score_command = new GetItemCommand(get_score);
			const get_score_res = await client.send(get_score_command);
			const prev_score = get_score_res.Item?.score?.N ? parseInt(get_score_res.Item.score.N, 10) : 0;
			const final_score = prev_score + parseInt(score, 10);
			const score_input = {
				"Item": {
					"server_id": {
						"S": interaction.guild.id
					},
					"user_id": {
						"S": interaction.user.id
					},
					"score": {
						"N": final_score.toString()
					}
				},
				"TableName": "leetboard_scores"
			};
			const command = new UpdateItemCommand(input);
			const submit_link_res = await client.send(command);
			const add_score = new PutItemCommand(score_input);
			const add_score_res = await client.send(add_score);
			interaction.followUp(`${difficulty} problem submitted. \n${interaction.user.username}'s score is now ${final_score}!`);
		} catch (error) {
			console.log(error);
			interaction.followUp('Could not submit problem, please try again');
		}
	}
}
