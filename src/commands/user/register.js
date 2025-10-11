const { SlashCommandBuilder } = require('discord.js');

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';

module.exports = {
    data: new SlashCommandBuilder()
        .setName('register')
        .setDescription("Register your LeetCode username to view your online stats")
        .addStringOption(option => option
            .setName('username')
            .setDescription('Your LeetCode online username')
            .setRequired(true)),

    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;
        const input_username = interaction.options.getString('username') ?? 'No username provided'
        const server_id = interaction.guild.id;
        const user_id = interaction.user.id;
        await interaction.deferReply();
        try {
            const response = await fetch(`${API_BASE_URL}/users/register`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    server_id: server_id,
                    user_id: user_id,
                    username: input_username
                })
            });

            if (!response.ok) {
                interaction.followUp(`Failed to register user ${input_username}. Please try again`);
                return;
            }

            interaction.followUp(`Successfully registered user ${input_username}`);
        } catch (error) {
            console.log(error);
            interaction.followUp(`Error registering, try again`);
        }
    }
}
