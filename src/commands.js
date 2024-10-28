require('dotenv').config();
const { REST, Routes } = require('discord.js');

const rest = new REST().setToken(process.env.TOKEN);

const commands = [
    {
        name: 'ping',
        description: 'replies with pong',
    },
];
client.once('ready', async () => {
    console.log("Bot is ready. Registering slash commands...");

    const rest = new REST().setToken(process.env.TOKEN);
    try {
        await rest.put(
            Routes.applicationGuildCommands(process.env.CLIENT_ID, process.env.GUILD_ID),
            { body: commands }
        );
        console.log("Slash commands registered successfully.");
    } catch (error) {
        console.error(`Unable to register commands: ${error.message}`);
    }
});
