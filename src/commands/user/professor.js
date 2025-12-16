const { SlashCommandBuilder, EmbedBuilder, AttachmentBuilder, ThreadAutoArchiveDuration } = require('discord.js');

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
            const response = await fetch(`${AGENT_API_URL}/api/v1/professor`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    professor_name: profName,
                    question: 'How is the following professor?'
                }),
                signal: AbortSignal.timeout(30000) // 30 second timeout
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error?.message || 'API request failed');
            }

            const { name, link, response: answerText } = await response.json();

            const file = new AttachmentBuilder("./src/assets/chicken.png");
            const embed = new EmbedBuilder()
                .setColor('#0099ff')
                .setTitle(name)
                .setDescription(answerText)
                .setThumbnail("attachment://chicken.png");

            if (link !== null && link !== "None") {
                embed.setURL(link);
            }

            await interaction.followUp({ embeds: [embed], files: [file] });

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

        } catch (error) {
            console.error("Error calling agent API:", error);
            await interaction.followUp({
                content: 'Our agent is currently unavailable. Please try again later!',
                ephemeral: true
            });
        }
    }
};
