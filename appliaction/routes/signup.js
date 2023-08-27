var express = require('express');
var users = require('../public/js/users');
var cert = require('../public/js/cert');
var router = express.Router();

// User Data for login
const USER_COOKIE_KEY = 'USER';

router.get('/', (req, res, next) => {
    res.render('signup', { title: 'CryptoFungi' });
});

router.post('/', async(req, res, next) => {
    const { username, name, password } = req.body;
    const exists = await users.fetchUser(username);
    const job = "user"

    // 이미 존재하는 username일 경우 회원 가입 실패
    if (exists) {
        res.status(400).send(`이미 존재하는 사용자 입니다.: ${username}`);
        return;
    }

    // 아직 가입되지 않은 username인 경우 db에 저장
    // KEY = username, VALUE = { name, password }
    const newUser = {
        username,
        name,
        password,
        job,
    };
    // db.set(username, newUser);
    await users.createUser(newUser)

    // make wallet for user
    await cert.makeUserWallet(username, "org1")

    // db에 저장된 user 객체를 문자열 형태로 변환하여 쿠키에 저장
    console.log(JSON.stringify(newUser))
    res.cookie(USER_COOKIE_KEY, JSON.stringify(newUser));
    // 가입 완료 후, 루트 페이지로 이동
    res.redirect('/');
});
router.post('/feedfactory', async(req, res, next) => {
    const { username, name, password } = req.body;
    const exists = await users.fetchUser(username);
    const job = "feedfactory"

    // 이미 존재하는 username일 경우 회원 가입 실패
    if (exists) {
        res.status(400).send(`이미 존재하는 사용자 입니다.: ${username}`);
        return;
    }

    // 아직 가입되지 않은 username인 경우 db에 저장
    // KEY = username, VALUE = { name, password }
    const newUser = {
        username,
        name,
        password,
        job,
    };
    // db.set(username, newUser);
    await users.createUser(newUser)

    // make wallet for user
    await cert.makeUserWallet(username, "org2")

    // db에 저장된 user 객체를 문자열 형태로 변환하여 쿠키에 저장
    console.log(JSON.stringify(newUser))
    res.cookie(USER_COOKIE_KEY, JSON.stringify(newUser));
    // 가입 완료 후, 루트 페이지로 이동
    res.redirect('/');
});
module.exports = router;


