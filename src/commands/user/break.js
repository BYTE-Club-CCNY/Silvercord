const { SlashCommandBuilder, EmbedBuilder, AttachmentBuilder } = require('discord.js');
const path = require('node:path');
const { execFile } = require('child_process');

module.exports = {
    data: new SlashCommandBuilder()
        .setName('break')
        .setDescription('break information')
        .addStringOption(option => option
            .setName('break')
            .setDescription('Returns information about CUNYs academic calendar')
            .setRequired(true)),
        
    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;
        const query = interaction.options.getString('break') ?? 'No break provided';

        await interaction.deferReply();

        const pythonScriptPath = path.resolve(__dirname, '../../../llm.py');
        const break_string = "break"
        execFile('python3', [pythonScriptPath, break_string, query], (error, stdout, stderr) => {
            if (error) {
                console.error("Error running LLM:", error);
                interaction.followUp(`Failed to get information about ${query}.`);
                return;
            }
            
            if (stderr) {
                console.error("Python script stderr:", stderr);
            }
            
            const {name, link, response} = JSON.parse(stdout.trim());
            const file = new AttachmentBuilder("./src/assets/calendar.jpg");
            embed = new EmbedBuilder()
                .setColor('#0099ff')
                .setTitle(name)
                .setURL(link)
                .setDescription(response)
                .setThumbnail("attachment://calendar.jpg");

            interaction.followUp({embeds: [ embed ], files: [ file ] });
        });
    }
};