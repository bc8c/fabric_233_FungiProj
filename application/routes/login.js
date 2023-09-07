var express = require('express');
var users = require('../public/js/users');
var router = express.Router();

// User Data for login
const USER_COOKIE_KEY = 'USER';

router.get('/', (req, res, next) => {
    res.render('login', { title: 'CryptoFungi' });
});

router.post('/', async(req, res, next) => {
    const { username, password } = req.body;
    const user = await users.fetchUser(username);

    // 가입 안 된 username인 경우
    if (!user) {
        res.status(400).send(`가입되지 않은 사용자입니다.: ${username}`);
        return;
    }
    // 비밀번호가 틀렸을 경우
    if (password !== user.password) {
        res.status(400).send('비밀번호가 틀렸습니다.');
        return;
    }

    // db에 저장된 user 객체를 문자열 형태로 변환하여 쿠키에 저장
    res.cookie(USER_COOKIE_KEY, JSON.stringify(user));
    // 로그인(쿠키 발급) 후, 루트 페이지로 이동
    res.redirect('/');
});
module.exports = router;


