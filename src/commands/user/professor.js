const { SlashCommandBuilder } = require('discord.js');
const path = require('node:path');
const { execFile } = require('child_process');

module.exports = {
    data: new SlashCommandBuilder()
        .setName('professor')
        .setDescription('Replies with professor information.')
        .addStringOption(option => option
            .setName('professor')
            .setDescription('The name of the professor you want to know more about.')
            .setRequired(true)),
        
    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;
        const profName = interaction.options.getString('professor') ?? 'No professor provided';

        await interaction.deferReply();

        const pythonScriptPath = path.resolve(__dirname, '../../../llm.py');

        execFile('python', [pythonScriptPath, profName], (error, stdout, stderr) => {
            if (error) {
                console.error("Error running LLM:", error);
                interaction.followUp(`Failed to get information about Professor ${profName}.`);
                return;
            }
            
            // TODO: Handle stderr
            if (stderr) {
                console.error("Python script stderr:", stderr);
                // interaction.followUp(`Could not retrieve information about Professor ${profName}.`);
                // return;
            }

            const response = stdout.trim();
            interaction.followUp(response);
        });
    }
};