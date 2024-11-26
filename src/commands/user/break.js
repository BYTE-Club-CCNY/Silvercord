const { SlashCommandBuilder, EmbedBuilder, AttachmentBuilder } = require('discord.js');
const path = require('node:path');
const { execFile } = require('child_process');

module.exports = {
    data: new SlashCommandBuilder()
        .setName('break')
        .setDescription('Retrieves CUNY break info')
        .addStringOption(option => option
            .setName('academic_calendar')
            .setDescription('Breaks Info')
            .setRequired(true)),
        
    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;
        const breakQuery = interaction.options.getString('query') ?? 'your query';

        await interaction.deferReply();

        const pythonScriptPath = path.resolve(__dirname, '../../../llm.py');

        execFile('python', [pythonScriptPath, breakQuery], (error, stdout, stderr) => {
            if (error) {
                console.error("Error running LLM:", error);
                interaction.followUp(`Failed to get information about ${breakQuery}.`);
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
                .setURL("https://www.cuny.edu/academics/academic-calendars/")
                .setDescription(response)
                .setThumbnail("attachment://calendar.jpg");

            interaction.followUp({embeds: [ embed ], files: [ file ] });
        });
    }
};