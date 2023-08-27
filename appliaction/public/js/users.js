const fs = require('fs').promises;
const USERS_JSON_FILENAME = './public/data/users.json';

async function fetchAllUsers() {
    const data = await fs.readFile(USERS_JSON_FILENAME);
    const users = JSON.parse(data.toString());
    return users;
}

async function fetchUser(username) {
    const users = await fetchAllUsers();
    const user = users.find((user) => user.username === username);
    return user;
}

async function createUser(newUser) {
    const users = await fetchAllUsers();
    users.push(newUser);
    await fs.writeFile(USERS_JSON_FILENAME, JSON.stringify(users));
}

module.exports.fetchAllUsers = fetchAllUsers;
module.exports.fetchUser = fetchUser;
module.exports.createUser = createUser;