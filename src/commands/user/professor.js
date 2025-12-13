const { SlashCommandBuilder, EmbedBuilder, AttachmentBuilder, ThreadAutoArchiveDuration } = require('discord.js');
const path = require('node:path');
const { execFile } = require('child_process');
const os = require('os');

module.exports = {
    data: new SlashCommandBuilder()
        .setName('professor')
        .setDescription('Replies with professor information.')
        .addStringOption(option => option
            .setName('professor')
            .setDescription('The name of the professor you want to know more about.')
            .setRequired(true))
		.addStringOption(option => option
			.setName('begin_thread')
			.setDescription('True or False')
			.setRequired(false)),
        
    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;
        const profName = interaction.options.getString('professor') ?? 'No professor provided';
        const beginThread = interaction.options.getString('begin_thread')?.toLowerCase() === 'true';

        await interaction.deferReply();

        const isWindows = os.platform() === 'win32';
        const venvSubPath = isWindows? 'Scripts/python.exe' : 'bin/python';
        const venvPath = path.resolve(__dirname, '../../../silvercord_agent/venv', venvSubPath);
        const pythonScriptPath = path.resolve(__dirname, '../../../silvercord_agent/agent.py');
        const prof_string = "professor"

        // example: `python3 ../agent.py professor Douglas Troeger
        execFile(venvPath, [pythonScriptPath, prof_string, profName], (error, stdout, stderr) => {
            if (error) {
                console.error("Error running LLM:", error);
                interaction.followUp(`Our agent is currently unavailable. Please try again later!`);
                return;
            }

            if (stderr) {
                console.error("Python script stderr:", stderr);
                interaction.followUp(`Internal error retrieving info about Professor ${profName}.`);
                return;
            }
            
            const {name, link, response} = JSON.parse(stdout.trim());
            const file = new AttachmentBuilder("./src/assets/chicken.png");
            const embed = new EmbedBuilder()
                .setColor('#0099ff')
                .setTitle(name)
                .setDescription(response)
                .setThumbnail("attachment://chicken.png");
            
            if (link !== null && link !== "None") {
                embed.setURL(link);
            }

            interaction.followUp({embeds: [ embed ], files: [ file ] }).then(async () => {
                if (beginThread && interaction.channel) {
                    try {
                        const thread = await interaction.channel.threads.create({
                            name: `Discussion: Professor ${name}`,
                            autoArchiveDuration: ThreadAutoArchiveDuration.OneHour,
                            reason: `Thread created for discussion about ${name}`
                        });
                        await thread.send(`Discussion thread for ${name}. Feel free to ask questions or share experiences!`);
                    } catch (error) {
                        console.error('Error creating thread:', error);
                    }
                }
            });
        });
    }
};
