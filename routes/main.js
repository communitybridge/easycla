if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var passport = require('passport');
var request = require('request');
var multer  = require('multer');
var async = require('async');

var cinco_api = require("../lib/api");

var router = express.Router();

const integration_user = process.env['CONSOLE_INTEGRATION_USER'];
const integration_pass = process.env['CONSOLE_INTEGRATION_PASSWORD'];

var hostURL = process.env['CINCO_SERVER_URL'];
if(!hostURL.startsWith('http')) hostURL = 'http://' + hostURL;
console.log("hostURL: " + hostURL);

var cinco = cinco_api(hostURL);

router.get('/', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getAllProjects(function (err, projects) {
      req.session.projects = projects;
      res.render('homepage', {projects: projects});
    });
  }
});

router.get('/angular', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  res.render('angular');
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
    if(user)
    {
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
        if(keys){
          req.session.user.cinco_keys = keys;
          var adminClient = cinco.client(keys);
          adminClient.getUser(lfid, function(err, user) {
            if(user){
              req.session.user.cinco_groups = JSON.stringify(user.groups);
              req.session.user.isAdmin = false;
              req.session.user.isUser = false;
              req.session.user.isProjectManager = false;
              if(user.groups)
              {
                if(user.groups.length > 0){
                  for(var i = 0; i < user.groups.length; i ++)
                  {
                    if(user.groups[i].name == "ADMIN") req.session.user.isAdmin = true;
                    if(user.groups[i].name == "USER") req.session.user.isUser = true;
                    if(user.groups[i].name == "PROJECT_MANAGER") req.session.user.isProjectManager = true;
                  }
                }
              }
              if( (req.session.user.isAdmin || req.session.user.isProjectManager) && (req.session.user.isUser)) {
                var projManagerClient = cinco.client(req.session.user.cinco_keys);
                projManagerClient.getAllProjects(function (err, projects) {
                  req.session.projects = projects;
                  projManagerClient.getMyProjects(function (err, myProjects) {
                    req.session.myProjects = myProjects;
                    return res.redirect('/');
                  });
                });

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
        if(err){
          req.session.destroy();
          console.log("getKeysForLfId err: " + err);
          if(err.statusCode == 404) return res.render('404', { lfid: lfid }); // Returned if a user with the given id is not found
          if(err.statusCode == 401) return res.render('401', { lfid: lfid }); // Unable to get keys for lfid given. User unauthorized.
        }
      });
    });
  })(req, res, next);
});

router.get('/profile', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  var adminClient = cinco.client(req.session.user.cinco_keys);
  var lfid = req.session.user.user;
  adminClient.getUser(lfid, function(err, user) {
    if(user){
      req.session.user.cinco_groups = JSON.stringify(user.groups);
      res.render('profile');
    }
    else {
      res.render('profile');
    }
  });
});

module.exports = router;
