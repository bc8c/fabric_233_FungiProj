var express = require('express');
var router = express.Router();
const cc = require('../public/js/cc');


router.get('/', (req, res, next) => {
    res.render('createFungus', { title: 'CryptoFungi' });
});

router.post('/', async(req, res, next) => {
    const name = req.body.name
    // 생성파트
    var result = await cc.cc_call("TestUser1", "CreateRandomFungus", name)
    console.log(result)
    res.redirect('/');
});

module.exports = router;