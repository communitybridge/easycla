if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var router = express.Router();

var cinco = require("../lib/api");

/**
* Organizations - Projects:
* Resources for getting details about an organizations project membership
**/

/**
* GET /organizations/{ogId}/projects_member
* Get all project memberships for this organization
**/
router.get('/organizations/:organizationId/projects_member', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var organizationId = req.params.organizationId;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getOrganizationProjectMemberships(organizationId, function (err, memberships) {
      if (err) return res.send('');
      res.send(memberships);
    });
  }
});

module.exports = router;
