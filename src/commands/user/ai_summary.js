const { SlashCommandBuilder,EmbedBuilder, AttachmentBuilder } = require('discord.js');
const path = require('node:path');
const { execFile } = require('child_process');
const os = require('os');

module.exports = {
    data: new SlashCommandBuilder()
        .setName('ai-summary')
        .setDescription('Summarizes the conversation between two given messages.')
        .addStringOption(option => option
            .setName('bottom-message')
            .setDescription('ID of the newest message.')
            .setRequired(true))
        .addStringOption(option => option
            .setName('top-message')
            .setDescription('ID of the oldest message')
            .setRequired(true)),
        
    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;

        const bottom_message_id = interaction.options.getString('bottom-message');
        const top_message_id = interaction.options.getString('top-message'); 

        await interaction.deferReply();
        const fetched_messages = [];
        
        let before_id = bottom_message_id; 
        let top_reached = false; 

        const max_fetches = 15; 
        let counter = 0;

        while (!top_reached && counter < max_fetches) {
            counter++;

            const messages = await interaction.channel.messages.fetch({limit: 100, before: before_id});

            if (messages.size === 0) {
                break;
            }

            for (const [id, message] of messages) {
                if (id === top_message_id) { 
                    fetched_messages.push(message); 
                    top_reached = true;
                    break;
                }
                
                fetched_messages.push(message);
            }
            
            before_id = messages.last().id;
        }
        
        try {
            const bottom_message = await interaction.channel.messages.fetch(bottom_message_id);
            fetched_messages.push(bottom_message);
            fetched_messages.reverse();

        } catch (error) {
            console.error('Error fetching bottom message:', error);
        }
        
        // turning it into a single string for use in gpt        
        const message_strings = fetched_messages.map(message => {
            return `[${message.createdAt.toISOString()}] ${message.author.tag}: ${message.content}`;
        }).filter(string => string.trim().length > 0); 
        const conversation_text = message_strings.join('\n');

        const isWindows = os.platform() === 'win32';
        const venvSubPath = isWindows? 'Scripts/python.exe' : 'bin/python';
        const venvPath = path.resolve(__dirname, '../../../silvercord_agent/venv', venvSubPath);
        const pythonScriptPath = path.resolve(__dirname, '../../../silvercord_agent/summary_helper.py');

        execFile(venvPath, [pythonScriptPath, conversation_text], (error, stdout, stderr) => {
            if (error) {
                console.error("Error running LLM:", error);
                interaction.followUp(`Our agent is currently unavailable. Please try again later!`);
                return;
            }

            if (stderr) {
                console.error("Python script stderr:", stderr);
                interaction.followUp(`Internal error occured while summarizing conversation.`);
                return;
            }
            
            console.log(stdout)
            const file = new AttachmentBuilder("./src/assets/chicken.png");
            const embed = new EmbedBuilder()
                .setColor('#0099ff')
                .setTitle("Silvercord")
                .setDescription(stdout)
                .setThumbnail("attachment://chicken.png");

            interaction.followUp({embeds: [ embed ], files: [ file ] })
        })

    },
};