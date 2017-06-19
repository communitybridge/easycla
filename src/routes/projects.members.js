if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var router = express.Router();

var cinco = require("../lib/api");

/**
* Projects - Members:
* Resources for getting details about project members
**/

/**
* GET /projects/{projectId}/members/
* Get all project members
**/
router.get('/projects/:projectId/members', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.projectId;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getProjectMembers(projectId, function (err, memberCompanies) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.send('');
      res.send(memberCompanies);
    });
  }
});

/**
* GET /projects/{projectId}/members/{memberId}
* Get an individual project member
**/
router.get('/projects/:projectId/members/:memberId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.projectId;
    var memberId = req.params.memberId;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getMemberFromProject(projectId, memberId, function (err, memberCompany) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.send('');
      res.send(memberCompany);
    });
  }
});

module.exports = router;
