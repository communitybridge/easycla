var express = require('express');
var request = require('request');

var cinco = require("../lib/api");

var router = express.Router();

router.get('/', function(req,res) {
  res.redirect('/');
});

module.exports = router;
