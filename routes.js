var express = require('express');
var passport = require('passport');

var router = express.Router();

router.get('/', function(req,res) {
  res.render('homepage');
});

router.get('/angular', function(req,res) {
  res.render('angular');
});

router.route('/login').get(function(req, res) {
  req.session.me = 'Beta Tester';
  // TODO: implement g auth
  res.redirect('/');
});

router.get('/logout', function(req, res){
  req.session.me = '';
  // req.logout(); TODO: implement g auth
  res.redirect('/');
});

module.exports = router;
