var express = require('express');
var passport = require('passport');
var request = require('request');
var multer  = require('multer');
var async = require('async');

var dummy_data = require('../dummy_db/dummy_data');
var cinco_api = require("../lib/api");

var router = express.Router();

var hostURL = process.env['CINCO_SERVER_URL'];
if(process.argv[2] == 'dev') hostURL = 'http://localhost:5000';
if(!hostURL.startsWith('http')) hostURL = 'http://' + hostURL;

var cinco = cinco_api(hostURL);

router.get('/mailing', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  dummy_data.findProjectById(req.query.id, function(err, project_data) {
    if(project_data) res.render('mailing', { project_data: project_data });
    else res.redirect('/');
  });
});

module.exports = router;
