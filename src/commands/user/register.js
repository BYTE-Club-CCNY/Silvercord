const { SlashCommandBuilder } = require('discord.js');
const path = require('node:path');
const { register_lc } = require('../../../api/dynamo_helper');

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
            const table = "leetboard"
            const registered = await register_lc(server_id, user_id, input_username, table)
            if (!registered) {
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
