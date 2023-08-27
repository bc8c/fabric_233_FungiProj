var express = require('express');
var router = express.Router();
const cc = require('../public/js/cc');

const USER_COOKIE_KEY = 'USER';

router.get('/', (req, res, next) => {
    res.render('createFeed', { title: 'FeedFactory' });
});

router.post('/', async(req, res, next) => {
    const userCookie  = req.cookies[USER_COOKIE_KEY];
    const userData = JSON.parse(userCookie);

    const id = userData.username    
    const name = req.body.name
    // 생성파트
    var result = await cc.cc_call(id,"CreateRandomFeed", name)
    console.log(result)
    res.redirect('/');
});

module.exports = router;