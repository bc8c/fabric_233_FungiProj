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
    const userData = JSON.parse(userCookie);
    const id = userData.username
    console.log(userData.job)
    if (userData.job == "user"){
      var result = await cc.cc_call(id,"GetFungiByOwner", "")
      console.log(result)
      res.render('index', { title: 'CryptoFungi', result: result });
    }
    else {      
      res.render('index', { title: 'FeedFactory', result: "feed" });
    }    
    console.log(result)
    
  }  
});

router.get('/logout', (req, res) => {
  // 쿠키 삭제 후 루트 페이지로 이동
  res.clearCookie(USER_COOKIE_KEY);
  res.redirect('/');
});

module.exports = router;
