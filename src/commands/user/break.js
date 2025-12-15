const { SlashCommandBuilder, EmbedBuilder, AttachmentBuilder } = require('discord.js');
require('dotenv').config();

module.exports = {
    data: new SlashCommandBuilder()
        .setName('break')
        .setDescription('Returns information about CUNYs academic calendar')
        .addStringOption(option => option
            .setName('2025-2026')
            .setDescription('2025-2026 CUNYs academic schedule')
            .setRequired(true)),
        
    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;
        const query = interaction.options.getString('2025-2026') ?? 'No break provided';

        await interaction.deferReply();

        const AGENT_API_URL = process.env.AGENT_API_URL;

        if (!AGENT_API_URL) {
            console.error("AGENT_API_URL environment variable is not set");
            await interaction.followUp({
                content: 'Configuration error: Agent API URL not set.',
                ephemeral: true
            });
            return;
        }

        try {
            const response = await fetch(`${AGENT_API_URL}/api/v1/break`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    year: query
                }),
                signal: AbortSignal.timeout(30000) // 30 second timeout
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error?.message || 'API request failed');
            }

            const { name, link, response: answerText } = await response.json();

            if (link === null || link === "None") {
                await interaction.followUp(answerText);
                return;
            }

            const file = new AttachmentBuilder("./src/assets/calendar.jpg");
            const embed = new EmbedBuilder()
                .setColor('#0099ff')
                .setTitle(name)
                .setURL(link)
                .setDescription(answerText)
                .setThumbnail("attachment://calendar.jpg");

            await interaction.followUp({ embeds: [embed], files: [file] });

        } catch (error) {
            console.error("Error calling agent API:", error);
            await interaction.followUp({
                content: 'Failed to retrieve calendar information. Please try again later!',
                ephemeral: true
            });
        }
    }
};
