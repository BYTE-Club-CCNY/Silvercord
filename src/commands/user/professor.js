const { SlashCommandBuilder, EmbedBuilder, AttachmentBuilder } = require('discord.js');
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
        const venvPath = path.resolve(__dirname, '../../../silvercord_agent/venv/Scripts/python.exe');
        const pythonScriptPath = path.resolve(__dirname, '../../../silvercord_agent/agent.py');
        const prof_string = "professor"

        // example: `python3 ../agent.py professor Douglas Troeger
        execFile(venvPath, [pythonScriptPath, prof_string, profName], (error, stdout, stderr) => {
            if (error) {
                console.error("Error running LLM:", error);
                interaction.followUp(`Our agent is currently unavailable. Please try again later!`);
                return;
            }

            if (stderr) {
                console.error("Python script stderr:", stderr);
                interaction.followUp(`Internal error retrieving info about Professor ${profName}.`);
                return;
            }
            
            const {name, link, response} = JSON.parse(stdout.trim());
            const file = new AttachmentBuilder("./src/assets/chicken.png");
            const embed = new EmbedBuilder()
                .setColor('#0099ff')
                .setTitle(name)
                .setURL(link)
                .setDescription(response)
                .setThumbnail("attachment://chicken.png");

            interaction.followUp({embeds: [ embed ], files: [ file ] });
        });
    }
};
