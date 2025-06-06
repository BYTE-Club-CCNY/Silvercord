async function fetchNodeFetch() {
  const fetch = (await import('node-fetch')).default; 
  return fetch;
}

async function get_online_username(username) {
    try {
        const fetch = await fetchNodeFetch(); 
        const url = "https://leetcode.com/graphql";
        
        const query = `
            query getUserProfile($username: String!) {
                matchedUser(username: $username) {
                    username
                    submitStats: submitStatsGlobal {
                        acSubmissionNum {
                            difficulty
                            count
                            submissions
                        }
                    }
                }
            }
        `;
        
        const vars = { username };
        const response = await fetch(url, { 
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ query, variables: vars }), 
        });

        if (response.ok) {
            const data = await response.json();
            console.log(data);
            return data.data.matchedUser; 
        } else {
            const errorMsg = await response.text();
            console.log(`Error occurred: ${errorMsg}`);
            return `Error occurred: ${errorMsg}`; 
        }
    } catch (error) {
        return `Request failed: ${error.message}`; 
    }
}


async function get_difficulty(link) {
    try {
        const fetch = await fetchNodeFetch(); 
        const titleSlug = link.split('/problems/')[1].split('/')[0];
        const url = "https://leetcode.com/graphql";
        const query = `
            query getProblemDetails($titleSlug: String!) {
                question(titleSlug: $titleSlug) {
                    title
                    difficulty
                }
            }`;
        const vars = { titleSlug };

        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ query, variables: vars }),
        });

        if (response.ok) {
            const data = await response.json();
            const difficulty = data.data.question.difficulty;
            return difficulty;
        } else {
            const errorMsg = await response.text();
            return `Error occurred: ${errorMsg}`;
        }
    } catch (error) {
        return error;
    }
}

function extractProblem(link) {
    try {
        if (link.includes('/submissions') && link.includes('/problems/')) {
            const problem = link.split('/problems/')[1].split('/submissions')[0].replace(/-/g, ' ');
            return problem;
        } else {
            return null;
        }
    } catch (error) {
        return error;
    }
}

module.exports = { extractProblem, get_difficulty, get_online_username };
