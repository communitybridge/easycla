if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var passport = require('passport');
var request = require('request');
var multer  = require('multer');
var async = require('async');

var cinco = require("../lib/api");

var router = express.Router();

var storage = multer.diskStorage({
  destination: function (req, file, cb) {
    cb(null, 'public/uploads')
  },
  filename: function (req, file, cb) {
    cb(null, file.originalname)
  }
});
var upload = multer({ storage: storage });
var cpUpload = upload.fields([{ name: 'logo', maxCount: 1 }, { name: 'agreement', maxCount: 1 }]);

/*
  Users:
  Resources to manage internal LF users and roles
 */

router.get('/users', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getAllUsers(function (err, users) {
      res.send(users);
    });
  }
});

module.exports = router;
