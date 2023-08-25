var express = require('express');
var router = express.Router();

const USER_COOKIE_KEY = 'USER';

/* GET home page. */
router.get('/', function(req, res, next) {
  const userCookie  = req.cookies[USER_COOKIE_KEY];
  console.log(userCookie)

  // 로그인이 안되어있는 경우
  if (!userCookie) {
    res.render('users', { title: 'CryptoFungi' });
  }
  else {
    res.render('index', { title: 'CryptoFungi' });
  }  
});

module.exports = router;
