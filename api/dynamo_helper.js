const { DynamoDBClient, GetItemCommand, PutItemCommand, UpdateItemCommand } = require("@aws-sdk/client-dynamodb");
const { dynamoConfig } = require('./aws-config');
const client = new DynamoDBClient(dynamoConfig);

async function get_score(server_id, user_id, table) {
    const get_score = {
        "Key": {
            "server_id": {
                "S": server_id 
            },
            "user_id": {
                "S": user_id 
            },
        },
        "TableName": table
    };
    const command = new GetItemCommand(get_score);
    const get_score_res = await client.send(command);
    return get_score_res.Item?.score?.N ? parseInt(get_score_res.Item.score.N, 10) : 0;
}

async function get_username(server_id, user_id, table) {
    const get_score = {
        "Key": {
            "server_id": {
                "S": server_id 
            },
            "user_id": {
                "S": user_id 
            },
        },
        "TableName": table
    };
    const command = new GetItemCommand(get_score);
    const get_username = await client.send(command);
    return get_username.Item?.username?.S ? get_username.Item.username.S : null;
}

async function update_score(server_id, user_id, score, table) {
    // doesn't need to return anything; sort of an in-place function
    const update_score = {
        "Item": {
            "server_id": {
                "S": server_id
            },
            "user_id": {
                "S": user_id
            },
            "score": {
                "N": score.toString()
            }
        },
        "TableName": table
    };
    const command = new PutItemCommand(update_score);
    const res = await client.send(command);
}

async function register_lc(server_id, user_id, username, table) {
    // returns true to ensure successfully registered 
    try {
        const register_lc = {
            "Key": {
                "server_id": {
                    "S": server_id
                },
                "user_id": {
                    "S": user_id
                }
            },
            "ExpressionAttributeNames": {
                "#username": "username"
            },
            "ExpressionAttributeValues": {
                ":username": { "S": username }
            },
            "UpdateExpression": "SET #username = :username",
            "TableName": table
        };
        const command = new UpdateItemCommand(register_lc)
        const res = await client.send(command);
        return true;
    } catch (error) {
        console.error("Error registering user:", error);
        return false;
    }
}

module.exports = { get_score, update_score, add_problem, get_problems, get_username, register_lc };
