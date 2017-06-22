if (process.env['NEWRELIC_LICENSE']) require('newrelic');
var express = require('express');
var router = express.Router();

var cinco = require("../lib/api");

/**
* Users
* Resources to manage internal LF users and roles
**/

/**
* GET /user
* Get current loggedIn user
**/
router.get('/user', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res) {
  if (req.session.user.isAdmin || req.session.user.isProjectManager) {
    var userId = req.session.user.user;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getUser(userId, function(err, user) {
      res.send(user);
    });
  }
});

/**
* GET /users
* Get All Users
**/
router.get('/users', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res) {
  if (req.session.user.isAdmin || req.session.user.isProjectManager) {
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getAllUsers(function(err, users) {
      res.send(users);
    });
  }
});

/**
* GET /users/{userId}
* Get an existing User
**/
router.get('/users/:userId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res) {
  if (req.session.user.isAdmin || req.session.user.isProjectManager) {
    var userId = req.params.userId;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getUser(userId, function(err, user) {
      res.send(user);
    });
  }
});

/**
* GET /users/roles
* Get all role enum values
**/
router.get('/users/roles', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getUserRoles(function (err, roles) {
      res.send(roles);
    });
  }
});

/**
* PUT /users/{userId}
* Update an existing User
**/
router.put('/users/:userId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res) {
  if (req.session.user.isAdmin || (req.session.user.isProjectManager && req.session.user.user == req.params.userId)) {
    var userId = req.params.userId;
    var user = {
      userId: req.body.userId,
      email: req.body.email,
      calendar: req.body.calendar
    };
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.updateUser(userId, user, function(err, user) {
      res.send(user);
    });
  }
});

module.exports = router;
