var express = require('express');
var passport = require('passport');

var router = express.Router();

router.get('/', function(req,res) {
  res.render('homepage', { user: req.user });
});

router.get('/angular', function(req,res) {
  res.render('angular', { user: req.user });
});

module.exports = router;
