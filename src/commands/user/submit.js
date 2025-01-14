const { SlashCommandBuilder } = require('discord.js');
const { get_difficulty, extractProblem } = require('../../../helper');
const { get_score, update_score, add_problem, get_problems } = require('../../../dynamo_helper');
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
            const server_id = interaction.guild.id;
            const user_id = interaction.user.id;
			const table = "leetboard_f24";
			const table_scores = "leetboard_scores_f24"
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
			// below retrieves data from dynamoDB; table name, partition key, and sort key are all required in the "query"
            const prev_score = await get_score(server_id, user_id, table_scores);
			const final_score = prev_score + parseInt(score, 10);
            // update the score:
            await update_score(server_id, user_id, final_score, table_scores);

            await add_problem(server_id, user_id, user_link, problem_name, table);
            console.log(`prev score: ${prev_score}, new score: ${final_score}`)
			interaction.followUp(`${difficulty} problem submitted. \n${interaction.user.username}'s score is now ${final_score}!`);
		} catch (error) {
			console.log(error);
			interaction.followUp('Could not submit problem, please try again');
		}
	}
}
