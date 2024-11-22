import { extractProblem, get_difficulty } from './helper.js';

const testOne = "https://leetcode.com/problems/search-insert-position/submissions/1421477674/"
const testTwo = "https://leetcode.com/problems/plus-one/submissions/1421458821/"
async function mainCall() {
	const nameOne = extractProblem(testOne);
	const diffOne = await get_difficulty(testOne);
	const nameTwo = extractProblem(testTwo);
	const diffTwo = await get_difficulty(testTwo);
	console.log("Name and difficulty problem #1", nameOne, diffOne)
	console.log("Name and difficulty problem #2", nameTwo, diffTwo)
}

mainCall();

