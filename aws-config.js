const { REST, Routes } = require("discord.js");
require('dotenv').config();

const dynamoConfig = {
	region: process.env.AWS_REGION,
	credentials: {
		accessKeyId: process.env.AWS_ACCESS_KEY,
		secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY,
	},
};

module.exports = { dynamoConfig };
