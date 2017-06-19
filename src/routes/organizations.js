if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var router = express.Router();

var cinco = require("../lib/api");

/**
* Organizations:
* Resources to expose and manipulate organizations
**/

/**
* GET /organizations/{orgId}
* Get a single organization by Id
**/
router.get('/organizations/:organizationId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res) {
  if (req.session.user.isAdmin || req.session.user.isProjectManager) {
    var organizationId = req.params.organizationId;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getOrganization(organizationId, function(err, organization) {
      if (err) return res.send('');
      res.send(organization);
    });
  }
});

module.exports = router;
