const { SlashCommandBuilder, EmbedBuilder, AttachmentBuilder } = require('discord.js');
require('dotenv').config();

module.exports = {
    data: new SlashCommandBuilder()
        .setName('break')
        .setDescription('Returns information about CUNYs academic calendar')
        .addStringOption(option => option
            .setName('query')
            .setDescription('2025-2026 CUNYs academic schedule')
            .setRequired(true)),
        
    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;
        const query = interaction.options.getString('query') ?? 'No break provided';

        await interaction.deferReply();

        try {
            const response = await fetch('http://localhost:8080/calendar', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    query: query,
                    user_id: interaction.user.id
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                console.error("Error from Go API:", errorData);
                await interaction.followUp(`Failed to get information about ${query}.`);
                return;
            }
            
            const {name, link, response: breakResponse} = await response.json();
            if (link === "None") {
                interaction.followUp(breakResponse);
                return;
            }
            const file = new AttachmentBuilder("./src/assets/calendar.jpg");
            const embed = new EmbedBuilder()
                .setColor('#0099ff')
                .setTitle(name)
                .setURL(link)
                .setDescription(breakResponse)
                .setThumbnail("attachment://calendar.jpg");

            interaction.followUp({embeds: [ embed ], files: [ file ] });
        } catch (error) {
            console.error("Error calling Go API:", error);
            await interaction.followUp(`Failed to get information for ${query}.`);
        }
    }
};
