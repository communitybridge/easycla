if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var passport = require('passport');
var request = require('request');
var multer  = require('multer');
var async = require('async');

var cinco = require("../lib/api");

var upload = multer(); // for parsing multipart/form-data

var router = express.Router();

var storageLogoCompany = multer.diskStorage({
  destination: function (req, file, cb) {
    cb(null, 'public/uploads')
  },
  filename: function (req, file, cb) {
    cb(null, file.originalname)
  }
});
var uploadLogoCompany = multer({ storage: storageLogoCompany });
var cpUploadLogoCompany = uploadLogoCompany.fields([
  {name: 'logoCompany', maxCount: 1},
  {name: 'board_headshot', maxCount: 1 },
  {name: 'technical_headshot', maxCount: 1 },
  {name: 'marketing_headshot', maxCount: 1 },
  {name: 'finance_headshot', maxCount: 1 },
  {name: 'other_headshot', maxCount: 1 }
]);

/*
  Organizations:
  Resources to expose and manipulate organizations
 */

 router.get('/organizations', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
   if(req.session.user.isAdmin || req.session.user.isProjectManager){
     var projManagerClient = cinco.client(req.session.user.cinco_keys);
     projManagerClient.getAllOrganizations(function (err, organizations) {
       if (err) return res.send('');
       res.send(organizations);
     });
   }
 });

 router.get('/organizations/:organizationId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
   if(req.session.user.isAdmin || req.session.user.isProjectManager){
     var organizationId = req.params.organizationId;
     var projManagerClient = cinco.client(req.session.user.cinco_keys);
     projManagerClient.getOrganization(organizationId, function (err, organization) {
       if (err) return res.send('');
       res.send(organization);
     });
   }
 });

/*
  Organizations - Contacts:
  Resources for getting and manipulating contacts of organizations
 */

router.get('/organizations/contacts/types', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getOrganizationContactTypes(function (err, contactTypes) {
      if (err) return res.send('');
      res.send(contactTypes);
    });
  }
});

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

router.post('/organizations/:organizationId/contacts', require('connect-ensure-login').ensureLoggedIn('/login'), cpUploadLogoCompany, function(req, res){
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

router.put('/organizations/:organizationId/contacts/:contactId', require('connect-ensure-login').ensureLoggedIn('/login'), cpUploadLogoCompany, function(req, res){
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
    projManagerClient.updateOrganizationContact(organizationId, contactId, contact, function (err, created, contact) {
      return res.json(contact);
    });
  }
});

/*
  Organizations - Projects:
  Resources for getting details about an organizations project membership
 */

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
