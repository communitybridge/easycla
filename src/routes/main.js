if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var passport = require('passport');
var router = express.Router();

var cinco = require("../lib/api");

router.get('/', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager) {
    res.redirect('/pmc')
  }
});

router.get('/logout', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  req.session.user = '';
  req.session.destroy();
  req.logout();
  res.redirect('/');
});

router.get('/login', function(req,res) {
  res.render('login');
});

router.get('/404', function(req,res) {
  res.render('404', { lfid: "" });
});

router.get('/401', function(req,res) {
  res.render('401', { lfid: "" });
});

router.get('/login_cas', function(req, res, next) {
  passport.authenticate('cas', function (err, user, info) {
    if (err) return next(err);
    if(user) {
      req.session.user = user;
    }
    if (!user) {
      req.session.destroy();
      return res.redirect('/login');
    }
    req.logIn(user, function (err) {
      if (err) return next(err);
      var lfid = req.session.user.user;
      cinco.getKeysForLfId(lfid, function (err, keys) {
        if(keys) {
          req.session.user.cinco_keys = keys;
          var adminClient = cinco.client(keys);
          adminClient.getUser(lfid, function(err, user) {
            if(user) {
              req.session.user.isAdmin = false;
              req.session.user.isUser = false;
              req.session.user.isProjectManager = false;
              if(user.roles) {
                if(user.roles.length > 0) {
                  for(var i = 0; i < user.roles.length; i ++) {
                    if(user.roles[i] == "ADMIN") req.session.user.isAdmin = true;
                    if(user.roles[i] == "USER") req.session.user.isUser = true;
                    if(user.roles[i] == "PROGRAM_MANAGER") req.session.user.isProjectManager = true;
                  }
                }
              }
              if( (req.session.user.isAdmin || req.session.user.isProjectManager) && (req.session.user.isUser)) {
                return res.redirect('/');
              }
              else {
                req.session.destroy();
                return res.render('401', { lfid: lfid }); // User unauthorized.
              }
            }
            else {
              return res.redirect('/login');
            }
          });
        }
        if(err) {
          req.session.destroy();
          console.log("getKeysForLfId err: " + err);
          if(err.statusCode == 404) return res.render('404', { lfid: lfid }); // Returned if a user with the given id is not found
          if(err.statusCode == 401) return res.render('401', { lfid: lfid }); // Unable to get keys for lfid given. User unauthorized.
        }
      });
    });
  })(req, res, next);
});

router.get('/session_data', function(req, res) {
  res.send({
    isAdmin: req.session.user.isAdmin,
    isProjectManager: req.session.user.isProjectManager,
    isUser: req.session.user.isUser,
    user: req.session.user.user,
  });
});

module.exports = router;
