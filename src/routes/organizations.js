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

module.exports = router;
