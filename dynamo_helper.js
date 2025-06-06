const { DynamoDBClient, GetItemCommand, PutItemCommand, UpdateItemCommand } = require("@aws-sdk/client-dynamodb");
const { dynamoConfig } = require('./aws-config');
const client = new DynamoDBClient(dynamoConfig);

async function fetchNodeFetch() {
  const fetch = (await import('node-fetch')).default; 
  return fetch;
}

async function add_problem(server_id, user_id, link, problem, table) {
    // for now dont need to return anything, can be added later

    const add_problem = {
        "Key": {
            "server_id": {
                "S": server_id
            },
            "user_id": {
                "S": user_id
            }
        },
        "ExpressionAttributeNames": {
            "#links": "link",
            "#problems": "problem"
        },
        "UpdateExpression": `
        SET
        #links = list_append(if_not_exists(#links, :empty_list), :new_link),
        #problems = list_append(if_not_exists(#problems, :empty_list), :new_problem)
        `,
        "ExpressionAttributeValues": {
            ":new_link": { L: [{ S: link }]},
            ":new_problem": { L: [{ S: problem }] },
            ":empty_list": { L: [] }
        },
        "TableName": table
    };
    const command = new UpdateItemCommand(add_problem);
    const res = await client.send(command);
}


async function get_problems(server_id, user_id, table) {
  const get_problem = {
    "Key": {
      "server_id": {
	"S": server_id
      },
      "user_id": {
	"S": user_id
      }
    },
    "TableName": table
  };
  const command = new GetItemCommand(get_problem);
  const res = await client.send(command);
  console.log(res.Item);
  return res.Item?.problem?.L ? res.Item.problem.L.map(item => item.S) : [];
}


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
