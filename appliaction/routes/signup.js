var express = require('express');
var users = require('../public/js/users');
var router = express.Router();

// User Data for login
const db = new Map();
const USER_COOKIE_KEY = 'USER';

router.get('/', (req, res, next) => {
    res.render('signup', { title: 'CryptoFungi' });
});

router.post('/', (req, res, next) => {
    const { username, name, password } = req.body;
    const exists = db.get(username);

    // 이미 존재하는 username일 경우 회원 가입 실패
    if (exists) {
        res.status(400).send(`duplicate username: ${username}`);
        return;
    }

    // 아직 가입되지 않은 username인 경우 db에 저장
    // KEY = username, VALUE = { name, password }
    const newUser = {
        username,
        name,
        password,
    };
    // db.set(username, newUser);
    users.createUser(username, name, password)

    // db에 저장된 user 객체를 문자열 형태로 변환하여 쿠키에 저장
    res.cookie(USER_COOKIE_KEY, JSON.stringify(newUser));
    // 가입 완료 후, 루트 페이지로 이동
    res.redirect('/');
});
module.exports = db;
module.exports = USER_COOKIE_KEY;
module.exports = router;


