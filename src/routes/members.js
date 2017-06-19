if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var passport = require('passport');
var request = require('request');
var multer  = require('multer');
var async = require('async');

var cinco = require("../lib/api");

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

router.get('/create_member', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    res.render('create_member');
  }
});

router.get('/create_member/:project_id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.project_id;
    res.render('create_member', {projectId: projectId});
  }
});

/*
  Projects - Members:
  Resources for getting details about project members
 */

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

router.get('/member/:project_id/:member_id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.project_id;
    var memberId = req.params.member_id;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getMemberFromProject(projectId, memberId, function (err, memberCompany) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.send('');
      res.send(memberCompany);
    });
  }
});

router.get('/edit_member/:project_id/:member_id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.project_id;
    var memberId = req.params.member_id;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getProject(projectId, function (err, project) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.redirect('/');
      projManagerClient.getMemberFromProject(projectId, memberId, function (err, memberCompany) {
        if(memberCompany){
          memberCompany.orgName = "";
          memberCompany.orgLogoRef = "";
          memberCompany.datepickerStartDate = "";
          memberCompany.datepickerRenewalDate = "";
          memberCompany.addresses = [];
          memberCompany.addresses.main = [];
          memberCompany.addresses.billing = [];
          memberCompany.contacts.board = [];
          memberCompany.contacts.technical = [];
          memberCompany.contacts.marketing = [];
          memberCompany.contacts.finance = [];
          memberCompany.contacts.other = [];
        }
        projManagerClient.getOrganization(memberCompany.orgId, function (err, organization) {
          if(organization){
            memberCompany.orgName = organization.name;
            memberCompany.orgLogoRef = organization.logoRef;
            memberCompany.addresses = organization.addresses;
            for (var j = 0; j < organization.addresses.length; j++){
              if (organization.addresses[j].type == 'MAIN') memberCompany.addresses.main = organization.addresses[j];
              else if (organization.addresses[j].type == 'BILLING') memberCompany.addresses.billing = organization.addresses[j];
            }
            var mainAddresses = memberCompany.addresses.main.address.thoroughfare.split(' /// ');
            memberCompany.addresses.main.address.thoroughfareLine1 = mainAddresses[0];
            memberCompany.addresses.main.address.thoroughfareLine2 = mainAddresses[1];

            var billingAddresses = memberCompany.addresses.billing.address.thoroughfare.split(' /// ');
            memberCompany.addresses.billing.address.thoroughfareLine1 = billingAddresses[0];
            memberCompany.addresses.billing.address.thoroughfareLine2 = billingAddresses[1];

            for (var j = 0; j < memberCompany.contacts.length; j++){
              if (memberCompany.contacts[j].type == 'BOARD MEMBER') memberCompany.contacts.board = memberCompany.contacts[j];
              else if (memberCompany.contacts[j].type == 'TECHNICAL') memberCompany.contacts.technical = memberCompany.contacts[j];
              else if (memberCompany.contacts[j].type == 'MARKETING') memberCompany.contacts.marketing = memberCompany.contacts[j];
              else if (memberCompany.contacts[j].type == 'FINANCE') memberCompany.contacts.finance = memberCompany.contacts[j];
              else memberCompany.contacts.other = memberCompany.contacts[j];
            }
          }
          if(memberCompany.startDate){
            // An integer number, between 0 and 11, representing the month in the given date according to local time.
            // 0 corresponds to January, 1 to February, and so on.
            var datepickerStartDate = new Date(memberCompany.startDate);
            startDateMonth = datepickerStartDate.getMonth() + 1;
            startDateDay = datepickerStartDate.getDate();
            startDateYear = datepickerStartDate.getFullYear();
            memberCompany.datepickerStartDate = startDateMonth + "/" + startDateDay + "/" + startDateYear;
          }
          if(memberCompany.renewalDate){
            var datepickerRenewalDate = new Date(memberCompany.renewalDate);
            renewalDateMonth = datepickerRenewalDate.getMonth() + 1;
            renewalDateDay = datepickerRenewalDate.getDate();
            renewalDateYear = datepickerRenewalDate.getFullYear();
            memberCompany.datepickerRenewalDate = renewalDateMonth + "/" + renewalDateDay + "/" + renewalDateYear;
          }
          return res.render('edit_member', {project: project, memberCompany:memberCompany});
        });
      });
    });
  }
});



module.exports = router;
