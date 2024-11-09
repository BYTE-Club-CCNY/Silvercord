const { SlashCommandBuilder } = require('discordjs');

authorizedUsers = ["1150800954606768150", "308735159958765569", "241734169057689602", "178357540831100928"];

module.exports = {
    data: new SlashCommandBuilder()
        .setName('sync')
        .setDescription('Sync the slash commands. Only select users can use this command.'),
    async execute(interaction) {
        // check if the id of the user who triggered the command is equal to a select user
        if (interaction.user.id in authorizedUsers) {
            await interaction.reply('Syncing commands...');
            
        } else {
            await interaction.reply('Pong!');
        }
    },
};