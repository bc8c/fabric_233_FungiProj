var express = require('express');
var router = express.Router();
const cc = require('../public/js/cc');

const USER_COOKIE_KEY = 'USER';

router.post('/', async(req, res, next) => {
    console.log("/in transferFrom")
    const userCookie  = req.cookies[USER_COOKIE_KEY];
    const userData = JSON.parse(userCookie);

    const id = userData.username    
    const fungusid = req.body.fungusid
    const to_id =  req.body.to_id
    const from_id = req.body.from_id
    var args = [from_id, to_id, fungusid];

    var result = await cc.cc_call(id, "TransferFrom", args)
    console.log(result)
    res.redirect('/');
});

module.exports = router;