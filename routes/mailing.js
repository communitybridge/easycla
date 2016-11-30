var express = require('express');
var passport = require('passport');
var request = require('request');
var multer  = require('multer');
var async = require('async');

var cinco_api = require("../lib/api");

var router = express.Router();

var hostURL = process.env['CINCO_SERVER_URL'];
if(process.argv[2] == 'dev') hostURL = 'http://localhost:5000';
if(!hostURL.startsWith('http')) hostURL = 'http://' + hostURL;

var cinco = cinco_api(hostURL);

router.get('/mailing/:projectId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.projectId;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getProject(projectId, function (err, project) {
      projManagerClient.getMailingLists(projectId, function (err, mailingList) {
        console.log(mailingList);
        res.render('mailing', { mailingList: mailingList, project:project });
      });
    });
  }
});

router.post('/mailing/:projectId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    var projectId = req.params.projectId;

    var mailingName = req.body.mailing_name;
    var mailingEmailAdmin = req.body.mailing_email_admin;
    var mailingPassword = req.body.mailing_password;
    // var mailingSubsribePolicy = req.body.mailing_subsribe_policy;
    // var mailingArchivePolicy = req.body.mailing_archive_policy;

    var newMailingList = {
      "name": mailingName,
      "admin": mailingEmailAdmin,
      "password": mailingPassword
      // ,
      // "subsribePolicy": mailingSubsribePolicy,
      // "archivePolicy": mailingArchivePolicy
    };
    projManagerClient.createMailingList(projectId, newMailingList, function (err, created, mailingListId) {
      console.log("mailing list created: " + created);
      console.log("mailingListId: " + mailingListId);
      return res.redirect('/mailing/' + projectId);
    });

  }
});

module.exports = router;
