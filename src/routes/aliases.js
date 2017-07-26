if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var router = express.Router();

var cinco = require("../lib/api");

/**
* Email Aliases
* Resources for working with email aliases of projects
**/

router.get('/aliases/:id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var id = req.params.id;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getProject(id, function (err, project) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.redirect('/');
      project.domain = "";
      if(project.url) project.domain = project.url.replace('http://www.','').replace('https://www.','').replace('http://','').replace('/','');
      projManagerClient.getEmailAliases(id, function (err, emailAliases) {
        var isAlias = {
          "isContactAlias": false,
          "isEventsAlias": false,
          "isPrAlias": false,
          "isLegalAlias": false,
          "isMembershipAlias": false,
        };
        if( emailAliases.length >= 1 ){
          for (var i = 0; i < emailAliases.length; i++) {
            var email_alias = emailAliases[i].address;
            var alias_name;
            if(email_alias) alias_name = email_alias.replace(/@.*$/,"");
            if(alias_name == "contact") isAlias.isContactAlias = true;
            if(alias_name == "events") isAlias.isEventsAlias = true;
            if(alias_name == "pr") isAlias.isPrAlias = true;
            if(alias_name == "legal") isAlias.isLegalAlias = true;
            if(alias_name == "membership") isAlias.isMembershipAlias = true;
          }
        }
        return res.render('aliases', {project: project, emailAliases: emailAliases, isAlias:isAlias });
      });
    });
  }
});

router.post('/aliases/:id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    var projectId = req.params.id;
    var email_alias = req.body.email_alias;
    var participant_email = req.body.participant_email;
    var newAlias = {
      "address": email_alias,
      "participants": [
        {
          "address": participant_email
        }
      ]
    };
    projManagerClient.createEmailAliases(projectId, newAlias, function (err, created, aliasId) {
      return res.redirect('/aliases/' + projectId);
    });
  }
});

router.post('/addParticipantToEmailAlias/:projectId/', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    var projectId = req.params.projectId;
    var aliasId = req.body.alias_id;
    var participant_email = req.body.participant_email;
    var newParticipant = {
      "address": participant_email
    };
    projManagerClient.addParticipantToEmailAlias(projectId, aliasId, newParticipant, function (err, created, response) {
      return res.redirect('/aliases/' + projectId);
    });
  }
});

router.get('/removeParticipantFromEmailAlias/:projectId/:aliasId/:participantEmail', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    var projectId = req.params.projectId;
    var aliasId = req.params.aliasId;
    var participantTBR = req.params.participantEmail;
    projManagerClient.removeParticipantFromEmailAlias(projectId, aliasId, participantTBR, function (err, removed) {
      return res.redirect('/aliases/' + projectId);
    });
  }
});

module.exports = router;
