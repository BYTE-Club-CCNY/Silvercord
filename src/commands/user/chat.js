const { SlashCommandBuilder, EmbedBuilder, AttachmentBuilder } = require('discord.js');
const path = require('node:path');
const { execFile } = require('child_process');

module.exports = {
	data: new SlashCommandBuilder()
		.setName('chat')
		.setDescription('Chat with our Claude AI about anything!')
		.addStringOption(option => option
			.setName('message')
			.setDescription('Your message to Claude 3 Sonnet')
			.setRequired(true)),
	async execute(interaction) {
		if (!interaction.isChatInputCommand()) return;
		const message = interaction.options.getString('message') ?? 'No message received';

		await interaction.deferReply();

		const pythonScriptPath = path.resolve(__dirname, '../../../llm.py');
		// note that this means that I will have to change llm.py to receive multiple times of commands
		// which will require different functions

		execFile('python', [pythonScriptPath, message], (error, stdout, stderr) => {
			if (error) {
				console.error("Error calling the LLM", error);
				interaction.followUp(`Failed to get a response from LLM.`);
				return;
			}

			if (stderr) {
				console.error("Python script stderr:", stderr);
			}

			const {message, response} = JSON.parse(stdout.trim());
			const file = new AttachmentBuilder("./src/assets/chicken.png");
			embed = new EmbedBuilder()
				.setColor('#0099ff')
				.setTitle('Claude')
				.addFields(
					{ name: 'You:', value: message },
					{ name: 'Claude', value: response },
				)
			interaction.followUp({embeds: [ embed ], files: [ file ] });
		});
	}
}
