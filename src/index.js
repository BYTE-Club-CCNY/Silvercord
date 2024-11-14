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
    /*{
        name: 'ping',
        description: 'replies with pongg',
    },*/
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
    {
        name: 'break',
        description: 'Ask about the academic breaks',
        options: [
            {
                name: 'academic_calendar',
                type: 3,
                description: 'Breaks Info',
                required: true,
            },
        ],
    },
    {
        name: 'talktuah',
        description: 'testing purposes',
        options: [
            {
                name: 'talktuah',
                type: 3,
                description: 'testing because im confused how this works',
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

// '/class' command -> authored by jaynopponep
client.on('interactionCreate', async (interaction) => {
    if (!interaction.isChatInputCommand()) return; // <- just checks if the interaction is a slash command 

    if (interaction.commandName === 'class') {
        const profName = interaction.options.getString('professor'); // gets professor parameter

        await interaction.deferReply(); // makes it to extend runtime because client must await LLM response

        const pythonScriptPath = path.resolve(__dirname, '../llm.py');

        execFile('python3', [pythonScriptPath, profName], (error, stdout, stderr) => { // using the child_process to run Python func for LLM within here
			// below is just some error handling
            if (error) {
                console.error("Error running LLM:", error);
                interaction.followUp(`Failed to get information about Professor ${profName}.`);
                return
            }

            if (stderr) {
                console.error("Python script stderr:", stderr);
                interaction.followUp(`Could not retrieve information about Professor ${profName}.`);
                return;
            }
			// sets up and replies with the response:
            const response = stdout.trim();
            interaction.followUp(response);
        });
    }
    
    if (interaction.commandName === 'talktuah') {
        await interaction.reply('testing');
        return;
    }
});

client.login(process.env.TOKEN);
