const { SlashCommandBuilder, EmbedBuilder, AttachmentBuilder } = require('discord.js');

module.exports = {
    data: new SlashCommandBuilder()
        .setName('professor')
        .setDescription('Replies with professor information.')
        .addStringOption(option => option
            .setName('professor')
            .setDescription('The name of the professor you want to know more about.')
            .setRequired(true)),
        
    async execute(interaction) {
        if (!interaction.isChatInputCommand()) return;
        const profName = interaction.options.getString('professor') ?? 'No professor provided';

        await interaction.deferReply();

        try {
            // Call the Go API endpoint
            const response = await fetch('http://localhost:8080/professor', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    professor_name: profName,
                    user_id: interaction.user.id
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                console.error("Error from Go API:", errorData);
                await interaction.followUp(`Failed to get information about Professor ${profName}.`);
                return;
            }

            const {name, link, response: profResponse} = await response.json();

            const file = new AttachmentBuilder("./src/assets/chicken.png");
            const embed = new EmbedBuilder()
                .setColor('#0099ff')
                .setTitle(name)
                .setURL(link)
                .setDescription(profResponse)
                .setThumbnail("attachment://chicken.png");

            await interaction.followUp({embeds: [ embed ], files: [ file ] });
        } catch (error) {
            console.error("Error calling Go API:", error);
            await interaction.followUp(`Failed to get information about Professor ${profName}.`);
        }
    }
};
