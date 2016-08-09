var express = require('express');
var passport = require('passport');
var dummy_data = require('./dummy_db/dummy_data');

var router = express.Router();

const integration_user = process.env['CONSOLE_INTEGRATION_USER'];
const integration_pass = process.env['CONSOLE_INTEGRATION_PASSWORD'];

router.get('/', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  res.render('homepage');
});

router.get('/angular', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  res.render('angular');
});

router.get('/logout', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  req.session.user = '';
  req.logout();
  res.redirect('/');
});

router.get('/login', function(req,res) {
  res.render('login');
});

router.get('/login_cas', function(req, res, next) {
  passport.authenticate('cas', function (err, user, info) {
    if (err) return next(err);
    if(user)
    {
      req.session.user = user;
    }
    if (!user) {
      return res.redirect('/login');
    }
    req.logIn(user, function (err) {
      if (err) return next(err);
      return res.redirect('/');
    });
  })(req, res, next);
});

router.get('/profile', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  res.render('profile');
});

router.get('/create_project', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  res.render('create_project');
});

router.get('/project', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  dummy_data.findProjectById(req.query.id, function(err, project_data) {
    res.render('project', { project_data: project_data });
  });

});

// Testing integration-platform keys endpoint
// TODO: Move to lib
var serverBaseURL = 'http://lf-integration-console-sandbox.us-west-2.elasticbeanstalk.com';
if(process.argv[2] == 'dev') serverBaseURL = 'http://localhost:5000';
router.get('/keys-test', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  var request = require('request');
  request.get(serverBaseURL + '/auth/trusted/cas/LaneMeyer', function (error, response, body) {
    if(response.statusCode == 200){
      body = JSON.parse(body);
      req.session.user.keyId = body.keyId;
      req.session.user.secret = body.secret;
    }
    res.render('keys-test');
   }).auth(integration_user, integration_pass, false);
});

module.exports = router;
