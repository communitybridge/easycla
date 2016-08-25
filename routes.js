var express = require('express');
var passport = require('passport');
var dummy_data = require('./dummy_db/dummy_data');
var request = require('request');

var router = express.Router();

const integration_user = process.env['CONSOLE_INTEGRATION_USER'];
const integration_pass = process.env['CONSOLE_INTEGRATION_PASSWORD'];

var hostURL = 'http://lf-integration-platform-sandbox.us-west-2.elasticbeanstalk.com';
if(process.argv[2] == 'dev') hostURL = 'http://localhost:5000';
console.log("hostURL: " + hostURL);

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
      request.get(hostURL + '/auth/trusted/cas/LaneMeyer', function (error, response, body) {
        if(response.statusCode == 200){
          body = JSON.parse(body);
          req.session.user.keyId = body.keyId;
          req.session.user.secret = body.secret;
          return res.redirect('/');
        }
        else{
          return res.redirect('/');
        }
       }).auth(integration_user, integration_pass, false);
    });
  })(req, res, next);
});

router.get('/profile', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  // Testing user/{lfid} endpoint
  // TODO: Move to lib
  var crypto = require('crypto');
  var payload = "";
  var httpMethod = 'GET';
  var uriPath = '/user/LaneMeyer';
  var currentTime = new Date().toISOString();
  var md5 = crypto.createHash('md5').update(payload).digest('hex');
  var signatureVersion = '1';
  var toSign = httpMethod + '\n' + uriPath + '\n' + currentTime + '\n' + md5 + '\n' + signatureVersion;
  var signature = crypto.createHmac('sha1', req.session.user.secret).update(toSign).digest('base64')
  request({
    method: httpMethod,
    url: hostURL + uriPath,
    headers: {
      'Content-Type': 'application/json',
      'Date': currentTime,
      'Signature-Version': '1',
      'Content-MD5': md5,
      'Authorization': 'CINCO '+ req.session.user.keyId + ': ' + signature
      }
    }, function(error, response){
      if(!error){
        var body = JSON.parse(response.body);
        req.session.user.integration_userId = body.userId;
        req.session.user.integration_groups = body.groups;
        res.render('profile');
      }
      else {
        console.log(error);
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

module.exports = router;
