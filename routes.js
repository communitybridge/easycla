var express = require('express');
var passport = require('passport');

var router = express.Router();

router.get('/', function(req,res) {
  res.render('homepage');
});

router.get('/angular', function(req,res) {
  res.render('angular');
});

router.get('/logout', function(req, res){
  req.session.me = '';
  req.logout();
  res.redirect('/');
});

router.route('/login').get(function(req, res, next) {
  passport.authenticate('cas', function (err, user, info) {
    if (err) return next(err);
    if(user)
    {
      req.session.me = user.attributes.profile_name_full;
      req.session.group = user.attributes.group;
      req.session.timezone = user.attributes.timezone;
    }
    if (!user) {
      return res.redirect('/');
    }
    req.logIn(user, function (err) {
      if (err) return next(err);
      return res.redirect('/');
    });
  })(req, res, next);
});

module.exports = router;
