if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var router = express.Router();

var cinco = require("../lib/api");

/**
* Organizations - Contacts
* Resources for getting and manipulating contacts of organizations
**/

/**
* GET /organizations/contacts/types
* Get all organization role enum values
**/
router.get('/organizations/contacts/types', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getOrganizationContactTypes(function (err, contactTypes) {
      if (err) return res.send('');
      res.send(contactTypes);
    });
  }
});

/**
* GET /organizations/{orgId}/contacts/
* Get a list of all the contacts in the organization
**/
router.get('/organizations/:organizationId/contacts', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var organizationId = req.params.organizationId;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getOrganizationContacts(organizationId, function (err, contactTypes) {
      if (err) return res.send('');
      res.send(contactTypes);
    });
  }
});

/**
* POST /organizations/{orgId}/contacts/
* Create a new contact
**/
router.post('/organizations/:organizationId/contacts', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var organizationId = req.params.organizationId;
    var contact = {
      type: req.body.type,
      givenName: req.body.givenName,
      familyName: req.body.familyName,
      title: req.body.title,
      bio: req.body.bio,
      email: req.body.email,
      phone: req.body.phone,
    }
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.createOrganizationContact(organizationId, contact, function (err, created, contactId) {
      return res.json(contactId);
    });
  }
});

/**
* GET /organizations/{orgId}/contacts/{contactId}
* Get a Contact in the organization by their contactId
**/
router.get('/organizations/:organizationId/contacts/:contactId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var organizationId = req.params.organizationId;
    var contactId = req.params.contactId;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getOrganizationContact(organizationId, contactId, function (err, contact) {
      if (err) return res.send('');
      res.send(contact);
    });
  }
});

/**
* PUT /organizations/{orgId}/contacts/{contactId}
* Update an existing Contact
**/
router.put('/organizations/:organizationId/contacts/:contactId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var organizationId = req.params.organizationId;
    var contactId = req.params.contactId;
    var contact = {
      type: req.body.type,
      givenName: req.body.givenName,
      familyName: req.body.familyName,
      title: req.body.title,
      bio: req.body.bio,
      email: req.body.email,
      phone: req.body.phone,
    }
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.updateOrganizationContact(organizationId, contactId, contact, function (err, updated, contact) {
      return res.json(contact);
    });
  }
});

module.exports = router;
