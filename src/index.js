require('dotenv').config();
const { Client, IntentsBitField, REST, Routes } = require('discord.js');
const { execFile } = require('child_process');
const path = require('path');

const client = new Client({
    intents: [
        IntentsBitField.Flags.Guilds,
        IntentsBitField.Flags.GuildMembers,
        IntentsBitField.Flags.GuildMessages,
        IntentsBitField.Flags.MessageContent,
    ],
});

const commands = [
    {
        name: 'ping',
        description: 'replies with pongg',
    },
    {
        name: 'class',
        description: 'Ask about a course based on professor name.',
        options: [
            {
                name: 'professor',
                type: 3, // STRING type for Discord
                description: 'Professor Name',
                required: true,
            },
        ],
    },
];

client.once('ready', async () => {
    console.log("Bot is ready. Registering slash commands...");

    const rest = new REST({ version: '10' }).setToken(process.env.TOKEN);
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

// class command
client.on('interactionCreate', async (interaction) => {
    if (!interaction.isChatInputCommand()) return;

    if (interaction.commandName === 'class') {
        const profName = interaction.options.getString('professor');

        await interaction.deferReply();

        const pythonScriptPath = path.resolve(__dirname, '../llm.py');

        execFile('python', [pythonScriptPath, profName], (error, stdout, stderr) => {
            if (error) {
                console.error("Error running LLM:", error);
                interaction.followUp(`Failed to get information about Professor ${profName}.`);
                return;
            }

            if (stderr) {
                console.error("Python script stderr:", stderr);
                interaction.followUp(`Could not retrieve information about Professor ${profName}.`);
                return;
            }

            const response = stdout.trim();
            interaction.followUp(response);
        });
    }
});

client.login(process.env.TOKEN);
