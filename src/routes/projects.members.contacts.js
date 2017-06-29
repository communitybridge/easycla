if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var router = express.Router();

var cinco = require("../lib/api");

/**
* Projects - Members - Contacts
* Resources for getting and manipulating contacts of project members
**/

/**
* GET /project/members/contacts/types
* Get all project contact role enum values
**/
router.get('/project/members/contacts/types', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getMemberContactRoles(function (err, roles) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.send('');
      res.send(roles);
    });
  }
});

/**
* GET /projects/{projectId}/members/{memberId}/contacts/
* Get the contacts for an individual project member
**/
router.get('/projects/:projectId/members/:memberId/contacts', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.projectId;
    var memberId = req.params.memberId;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getMemberContacts(projectId, memberId, function (err, contacts) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.send('');
      res.send(contacts);
    });
  }
});

/**
* POST /projects/{projectId}/members/{memberId}/contacts/{contactId}/
* Add a Contact for the Member
**/
router.post('/projects/:projectId/members/:memberId/contacts/:contactId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.projectId;
    var memberId = req.params.memberId;
    var contactId = req.params.contactId;
    var memberContact = {
      id: req.body.id,
      memberId: req.body.memberId,
      type: req.body.type,
      primaryContact: req.body.primaryContact,
      boardMember: req.body.boardMember,
      contact: {
        id: req.body.contactId,
        email: req.body.contactEmail,
        givenName: req.body.contactGivenName,
        familyName: req.body.contactFamilyName,
        title: req.body.contactTitle,
        phone: req.body.contactPhone,
        type: req.body.contactType,
        bio: req.body.contactBio,
      },
    };
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.addMemberContact(projectId, memberId, contactId, memberContact, function (err, created, obj) {
      return res.json(obj);
    });
  }
});

/**
* DELETE /projects/{projectId}/members/{memberId}/contacts/{contactId}/roles/{roleId}
* Remove a contact from the member
**/
router.delete('/projects/:projectId/members/:memberId/contacts/:contactId/roles/:roleId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.projectId;
    var memberId = req.params.memberId;
    var contactId = req.params.contactId;
    var roleId = req.params.roleId;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.removeMemberContact(projectId, memberId, contactId, roleId, function (err, removed) {
      return res.json(removed);
    });
  }
});

/**
* PUT projects/{projectId}/members/{memberId}/contacts/{contactId}/roles/{roleId}
* Update an existing contact
**/
router.put('/projects/:projectId/members/:memberId/contacts/:contactId/roles/:roleId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.projectId;
    var memberId = req.params.memberId;
    var contactId = req.params.contactId;
    var roleId = req.params.roleId;
    var memberContact = {
      id: req.body.id,
      memberId: req.body.memberId,
      type: req.body.type,
      primaryContact: req.body.primaryContact,
      boardMember: req.body.boardMember,
      contact: {
        id: req.body.contactId,
        email: req.body.contactEmail,
        givenName: req.body.contactGivenName,
        familyName: req.body.contactFamilyName,
        title: req.body.contactTitle,
        phone: req.body.contactPhone,
        type: req.body.contactType,
        bio: req.body.contactBio,
      },
    };
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.updateMemberContact(projectId, memberId, contactId, roleId, memberContact, function (err, updated, contact) {
      return res.json(contact);
    });
  }
});

module.exports = router;
