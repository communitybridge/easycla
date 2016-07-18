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
  res.redirect('/auth/google');
});

router.get('/logout', function(req, res){
  req.session.me = '';
  req.logout();
  res.redirect('/');
});

router.get('/auth/google', passport.authenticate('google', { scope: ['profile'] }));

router.get('/auth/google/callback', passport.authenticate('google', { failureRedirect: '/login' }), function(req, res) {
  req.session.me = req.user.displayName;
  res.redirect('/');
});

module.exports = router;
