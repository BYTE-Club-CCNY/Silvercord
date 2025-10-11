const { SlashCommandBuilder } = require('discord.js');

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';

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
			const server_id = interaction.guild.id;
			const user_id = interaction.user.id;

			const extractResponse = await fetch(`${API_BASE_URL}/leetcode/extract-problem?link=${encodeURIComponent(user_link)}`);
			if (!extractResponse.ok) {
				interaction.followUp(`Link ${user_link} is in the wrong format, resubmit the link that contains the submission ID`);
				return;
			}
			const extractData = await extractResponse.json();
			const problem_name = extractData.problem;

			if (!problem_name || problem_name === "") {
				interaction.followUp(`Link ${user_link} is in the wrong format, resubmit the link that contains the submission ID`);
				return;
			}

			const problemsResponse = await fetch(`${API_BASE_URL}/problems?server_id=${server_id}&user_id=${user_id}`);
			if (!problemsResponse.ok) {
				interaction.followUp('Failed to fetch your problems. Please try again.');
				return;
			}
			const problemsData = await problemsResponse.json();
			console.log('problems:', problemsData);
			const problems = (problemsData || []).map(p => p.problem);

			if (problems.includes(problem_name)) {
				interaction.followUp(`You've already done this problem. Nice try.`)
				return;
			}

			const difficultyResponse = await fetch(`${API_BASE_URL}/leetcode/difficulty?link=${encodeURIComponent(user_link)}`);
			if (!difficultyResponse.ok) {
				interaction.followUp('Failed to fetch problem difficulty. Please try again.');
				return;
			}
			const difficultyData = await difficultyResponse.json();
			const difficulty = difficultyData.difficulty;

			let score = 0;
			if (difficulty === "Easy") {
				score = 1;
			} else if (difficulty === "Medium") {
				score = 2;
			} else {
				score = 4;
			}

			const scoreResponse = await fetch(`${API_BASE_URL}/scores?server_id=${server_id}&user_id=${user_id}`);
			if (!scoreResponse.ok) {
				interaction.followUp('Failed to fetch your score. Please try again.');
				return;
			}
			const scoreData = await scoreResponse.json();
			const prev_score = scoreData.score;
			const final_score = prev_score + score;

			const leaderboardResponse = await fetch(`${API_BASE_URL}/scores/leaderboard?server_id=${server_id}`);
			if (!leaderboardResponse.ok) {
				interaction.followUp('Failed to fetch leaderboard. Please try again.');
				return;
			}
			const leaderboard = await leaderboardResponse.json();

			let nextUserScore = 0;
			let nextUserID = user_id;

			if (leaderboard && leaderboard.length > 1) {
				for (let i = 1; i < leaderboard.length; i++) {
					if (leaderboard[i].user_id == user_id && leaderboard[i].score == prev_score) {
						if (leaderboard[i - 1].user_id != user_id) {
							nextUserID = leaderboard[i - 1].user_id;
							nextUserScore = leaderboard[i - 1].score;
							break;
						} else if (i + 1 < leaderboard.length && leaderboard[i + 1].user_id != user_id && leaderboard[i + 1].score == prev_score) {
							nextUserID = leaderboard[i + 1].user_id;
							nextUserScore = leaderboard[i + 1].score;
						}
					}
				}
			}

	    await fetch(`${API_BASE_URL}/scores`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({
		    server_id: server_id,
		    user_id: user_id,
		    score: final_score
		})
	    });

	    await fetch(`${API_BASE_URL}/problems`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({
		    server_id: server_id,
		    user_id: user_id,
		    link: user_link,
		    problem: problem_name
		})
	    });

			if (prev_score <= nextUserScore && nextUserScore < final_score) {
				try {
					const nextUser = await interaction.client.users.fetch(nextUserID);
					interaction.followUp(`${difficulty} problem submitted. \n${interaction.user.username} has taken ${nextUser.username}'s place with new score ${final_score}!`);
				} catch (error) {
					console.error('error fetching username', error);
					interaction.followUp(`${difficulty} problem submitted. \n${interaction.user.username}'s score is now ${final_score}!`);
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
