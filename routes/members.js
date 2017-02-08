if (process.env['NEW_RELIC_API_KEY']) require('newrelic');
var express = require('express');
var passport = require('passport');
var request = require('request');
var multer  = require('multer');
var async = require('async');

var cinco_api = require("../lib/api");

var router = express.Router();

var hostURL = process.env['CINCO_SERVER_URL'];
if(!hostURL.startsWith('http')) hostURL = 'http://' + hostURL;

var cinco = cinco_api(hostURL);

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

router.post('/create_member', require('connect-ensure-login').ensureLoggedIn('/login'), cpUploadLogoCompany, function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    var projectId = req.body.project_id;
    var startDate = "";
    var renewalDate = "";
    if (req.body.start_date) startDate = new Date(req.body.start_date).toISOString();
    if (req.body.renewal_date) renewalDate = new Date(req.body.renewal_date).toISOString();
    var logoCompanyFileName = "";
    if(req.files){
      if(req.files.logoCompany) logoCompanyFileName = req.files.logoCompany[0].originalname;
    }
    //country code must be exactly 2 Alphabetic characters or null
    var headquartersCountry = req.body.headquarters_country;
    if(headquartersCountry == "") headquartersCountry = null;

    var billingCountry = req.body.billing_country;
    if(billingCountry == "") billingCountry = null;

    var mainThoroughfare = req.body.headquarters_address_line_1 + " /// " + req.body.headquarters_address_line_2;
    var billingThoroughfare = req.body.billing_address_line_1 + " /// " + req.body.billing_address_line_2;

    var newOrganization = {
      name: req.body.company_name,
      addresses: [
        {
          type: "MAIN",
          address: {
            country: headquartersCountry,
            administrativeArea: req.body.headquarters_state,
            localityName: req.body.headquarters_city,
            postalCode: req.body.headquarters_zip_code,
            phone: req.body.headquarters_phone,
            thoroughfare: mainThoroughfare
          }
        },
        {
          type: "BILLING",
          address: {
            country: billingCountry,
            administrativeArea: req.body.billing_state,
            localityName: req.body.billing_city,
            postalCode: req.body.billing_zip_code,
            phone: req.body.billing_phone,
            thoroughfare: billingThoroughfare
          }
        }
      ],
      logoRef : logoCompanyFileName
    }

    projManagerClient.createOrganization(newOrganization, function (err, created, organizationId) {
      console.log(err);
      if(created && projectId){
        var newMember = {
          orgId: organizationId,
          tier: {
            type: req.body.membership_tier,
            qualifier: 1 // Optional Tier Level
          },
          startDate: startDate,
          renewalDate: renewalDate
        };
        projManagerClient.addMemberToProject(projectId, newMember, function (err, created, memberId) {
          var newContacts = JSON.parse(req.body.newContacts);
          async.forEach(newContacts, function (eachContact, callback){
            projManagerClient.addContactToMember(projectId, memberId, eachContact, function (err, created, contactId) {
              callback();
            });
          }, function(err) {
            // Contacts iteration done.
            return res.redirect('/project/' + projectId);
          });
        });
      }
      else{
        return res.redirect('/');
      }
    });
  }
});

router.get('/members/:id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.id;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getProject(projectId, function (err, project) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.redirect('/');
      projManagerClient.getMemberCompanies(projectId, function (err, memberCompanies) {
        async.forEach(memberCompanies, function (eachMember, callback){
          eachMember.orgName = "";
          eachMember.orgLogoRef = "";
          projManagerClient.getOrganization(eachMember.orgId, function (err, organization) {
            if(organization){
              eachMember.orgName = organization.name;
              eachMember.orgLogoRef = organization.logoRef;
            }
            callback();
          });
        }, function(err) {
          // Member Companies iteration done.
          return res.render('members', {project: project, memberCompanies:memberCompanies});
        });
      });
    });
  }
});

router.get('/member/:project_id/:member_id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
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
          return res.render('member', {project: project, memberCompany:memberCompany});
        });
      });
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

router.post('/edit_member/:project_id/:organization_id/:member_id', require('connect-ensure-login').ensureLoggedIn('/login'), cpUploadLogoCompany, function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.project_id;
    var organizationId = req.params.organization_id;
    var memberId = req.params.member_id;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    var startDate = "";
    var renewalDate = "";
    if (req.body.start_date) startDate = new Date(req.body.start_date).toISOString();
    if (req.body.renewal_date) renewalDate = new Date(req.body.renewal_date).toISOString();
    var logoCompanyFileName = "";
    if(req.files){
      if(req.files.logoCompany) logoCompanyFileName = req.files.logoCompany[0].originalname;
      else logoCompanyFileName = req.body.old_logoRef;
    }
    //country code must be exactly 2 Alphabetic characters or null
    var headquartersCountry = req.body.headquarters_country;
    if(headquartersCountry == "") headquartersCountry = null;

    var billingCountry = req.body.billing_country;
    if(billingCountry == "") billingCountry = null;

    var mainThoroughfare = req.body.headquarters_address_line_1 + " /// " + req.body.headquarters_address_line_2;
    var billingThoroughfare = req.body.billing_address_line_1 + " /// " + req.body.billing_address_line_2;

    var updatedOrganization = {
      id: organizationId,
      name: req.body.company_name,
      addresses: [
        {
          type: "MAIN",
          address: {
            country: headquartersCountry,
            administrativeArea: req.body.headquarters_state,
            localityName: req.body.headquarters_city,
            postalCode: req.body.headquarters_zip_code,
            phone: req.body.headquarters_phone,
            thoroughfare: mainThoroughfare
          }
        },
        {
          type: "BILLING",
          address: {
            country: billingCountry,
            administrativeArea: req.body.billing_state,
            localityName: req.body.billing_city,
            postalCode: req.body.billing_zip_code,
            phone: req.body.billing_phone,
            thoroughfare: billingThoroughfare
          }
        }
      ],
      logoRef : logoCompanyFileName
    }
    var updatedMember = {
      tier: {
        type: req.body.membership_tier,
        qualifier: 1 // Optional Tier Level
      },
      startDate: startDate,
      renewalDate: renewalDate
    }
    projManagerClient.updateOrganization(updatedOrganization, function (err, updated, organization) {
      projManagerClient.updateMember(projectId, memberId, updatedMember, function (err, updated, updatedMember) {
        var updatedContacts = JSON.parse(req.body.updatedContacts);
        async.forEach(updatedContacts, function (eachUpdatedContact, callback){
          // pass each contactId
          projManagerClient.updateContactFromMember(projectId, memberId, eachUpdatedContact.id, eachUpdatedContact, function (err, udpated, contactId) {
            callback();
          });
        }, function(err) {
          // Contacts iteration done.
          return res.redirect('/member/' + projectId + '/' + memberId);
        });
      });
    });
  }
});

module.exports = router;
