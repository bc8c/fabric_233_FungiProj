var express = require('express');
var router = express.Router();

const cc = require('../public/js/cc');

const USER_COOKIE_KEY = 'USER';

/* GET home page. */
router.get('/', async function(req, res, next) {
  const userCookie  = req.cookies[USER_COOKIE_KEY];
  console.log(userCookie)

  // 로그인이 안되어있는 경우
  if (!userCookie) {
    res.render('users', { title: 'CryptoFungi' });
  }
  else { // 로그인이 되어있는 경우
    var result = await cc.cc_call("TestUser1", "GetFungiByOwner", "")
    console.log(result)
    res.render('index', { title: 'CryptoFungi', result: result });
  }  
});

module.exports = router;
