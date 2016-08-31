var express = require('express');
var passport = require('passport');
var dummy_data = require('./dummy_db/dummy_data');
var request = require('request');
var cinco_api = require("./lib/api");

var router = express.Router();

const integration_user = process.env['CONSOLE_INTEGRATION_USER'];
const integration_pass = process.env['CONSOLE_INTEGRATION_PASSWORD'];

var hostURL = 'http://lf-integration-platform-sandbox.us-west-2.elasticbeanstalk.com';
if(process.argv[2] == 'dev') hostURL = 'http://localhost:5000';
console.log("hostURL: " + hostURL);

var cinco = cinco_api(hostURL);

router.get('/', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  res.render('homepage');
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
          req.session.user.keyId = keys.keyId;
          req.session.user.secret = keys.secret;
          req.session.user.keys = keys;
          return res.redirect('/');
        }
        if(err){
          req.session.destroy();
          if(err.statusCode == 404) { // Returned if a user with the given id is not found
            return res.render('404', { lfid: lfid });
          }
        }
      });
    });
  })(req, res, next);
});

router.get('/profile', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  var adminClient = cinco.client(req.session.user.keys);
  var lfid = req.session.user.user;
  adminClient.getUser(lfid, function(err, user) {
    if(user){
      req.session.user.integration_userId = user.userId;
      req.session.user.integration_groups = JSON.stringify(user.groups);
      res.render('profile');
    }
    else {
      res.render('profile');
    }
  });
});

router.get('/create_project', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  res.render('create_project');
});

router.get('/project', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  dummy_data.findProjectById(req.query.id, function(err, project_data) {
    if(project_data) res.render('project', { project_data: project_data });
    else res.redirect('/');
  });
});

router.get('/mailing', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  dummy_data.findProjectById(req.query.id, function(err, project_data) {
    if(project_data) res.render('mailing', { project_data: project_data });
    else res.redirect('/');
  });
});

router.get('/aliases', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  dummy_data.findProjectById(req.query.id, function(err, project_data) {
    if(project_data) res.render('aliases', { project_data: project_data });
    else res.redirect('/');
  });
});

router.get('/members', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  dummy_data.findProjectById(req.query.id, function(err, project_data) {
    if(project_data) res.render('members', { project_data: project_data });
    else res.redirect('/');
  });
});

router.get('/admin', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  var adminClient = cinco.client(req.session.user.keys);
  var lfid = req.session.user.user;
  adminClient.getUser(lfid, function(err, user) {
    console.log(user);
  });
  res.render('admin');
});

module.exports = router;
