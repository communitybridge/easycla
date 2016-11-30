var express = require('express');
var passport = require('passport');
var request = require('request');
var multer  = require('multer');
var async = require('async');

var cinco_api = require("../lib/api");

var router = express.Router();

var hostURL = process.env['CINCO_SERVER_URL'];
if(process.argv[2] == 'dev') hostURL = 'http://localhost:5000';
if(!hostURL.startsWith('http')) hostURL = 'http://' + hostURL;

var cinco = cinco_api(hostURL);

router.get('/admin', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin) {
    var adminClient = cinco.client(req.session.user.cinco_keys);
    var username = req.body.form_lfid;
    adminClient.getAllUsers(function (err, users, groups) {
      res.render('admin', { message: "", users: users, groups: groups });
    });
  }
  else res.redirect('/');
});

router.post('/activate_user', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin){
    var adminClient = cinco.client(req.session.user.cinco_keys);
    var username = req.body.form_lfid;
    var userGroup = {
      groupId: 1,
      name: 'USER'
    }
    adminClient.addGroupForUser(username, userGroup, function(err, isUpdated, user) {
      var message = 'User has been activated.';
      if (err) message = err;
      adminClient.getAllUsers(function (err, users, groups) {
        return res.render('admin', { message: message, users: users, groups:groups });
      });
    });
  }
});

router.post('/create_user', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin){
    var adminClient = cinco.client(req.session.user.cinco_keys);
    var username = req.body.form_lfid;
    adminClient.createUser(username, function (err, created) {
      var message = '';
      if (err) {
        message = err;
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
      if(created) {
        message = 'User has been created.';
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
      else {
        message = 'User already exists. ';
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
    });
  }
});

router.post('/create_project_manager_user', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin){
    var adminClient = cinco.client(req.session.user.cinco_keys);
    var username = req.body.form_lfid;
    var projectManagerGroup = {
      groupId: 3,
      name: 'PROJECT_MANAGER'
    }
    var userGroup = {
      groupId: 1,
      name: 'USER'
    }
    adminClient.createUser(username, function (err, created) {
      var message = '';
      if (err) {
        message = err;
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
      if(created) {
        message = 'Project Manager has been created.';
        adminClient.addGroupForUser(username, userGroup, function(err, isUpdated, user) {});
        adminClient.addGroupForUser(username, projectManagerGroup, function(err, isUpdated, user) {
          if (err) message = err;
          adminClient.getAllUsers(function (err, users, groups) {
            return res.render('admin', { message: message, users: users, groups:groups });
          });
        });
      }
      else {
        message = 'User already exists.';
        adminClient.addGroupForUser(username, userGroup, function(err, isUpdated, user) {});
        adminClient.addGroupForUser(username, projectManagerGroup, function(err, isUpdated, user) {
          message = 'User already exists. ' + 'Project Manager has been created.';
          if (err) message = err;
          adminClient.getAllUsers(function (err, users, groups) {
            return res.render('admin', { message: message, users: users, groups:groups });
          });
        });
      }
    });
  }
});

router.post('/create_admin_user', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin){
    var adminClient = cinco.client(req.session.user.cinco_keys);
    var username = req.body.form_lfid;
    var adminGroup = {
      groupId: 2,
      name: 'ADMIN'
    }
    var userGroup = {
      groupId: 1,
      name: 'USER'
    }
    adminClient.createUser(username, function (err, created) {
      var message = '';
      if (err) {
        message = err;
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
      if(created) {
        message = 'Admin has been created.';
        adminClient.addGroupForUser(username, userGroup, function(err, isUpdated, user) {});
        adminClient.addGroupForUser(username, adminGroup, function(err, isUpdated, user) {
          if (err) message = err;
          adminClient.getAllUsers(function (err, users, groups) {
            return res.render('admin', { message: message, users: users, groups:groups });
          });
        });
      }
      else {
        message = 'User already exists. ' + 'Admin has been created.';
        adminClient.addGroupForUser(username, userGroup, function(err, isUpdated, user) {});
        adminClient.addGroupForUser(username, adminGroup, function(err, isUpdated, user) {
          if (err) message = err;
          adminClient.getAllUsers(function (err, users, groups) {
            return res.render('admin', { message: message, users: users, groups:groups });
          });
        });
      }
    });
  }
});

router.post('/deactivate_user', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin){
    var adminClient = cinco.client(req.session.user.cinco_keys);
    var username = req.body.form_lfid;
    var userGroup = {
      groupId: 1,
      name: 'USER'
    }
    adminClient.removeGroupFromUser(username, userGroup.groupId, function (err, removed) {
      var message = '';
      if (err) {
        message = err;
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
      if(removed) {
        message = 'User has been deactivated.';
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
    });
  }
});

router.post('/remove_admin_user', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin){
    var adminClient = cinco.client(req.session.user.cinco_keys);
    var username = req.body.form_lfid;
    var adminGroup = {
      groupId: 2,
      name: 'ADMIN'
    }
    adminClient.removeGroupFromUser(username, adminGroup.groupId, function (err, removed) {
      var message = '';
      if (err) {
        message = err;
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
      if(removed) {
        message = 'Admin has been removed.';
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
    });
  }
});

router.post('/remove_project_manager_user', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin){
    var adminClient = cinco.client(req.session.user.cinco_keys);
    var username = req.body.form_lfid;
    var projectManagerGroup = {
      groupId: 3,
      name: 'PROJECT_MANAGER'
    }
    adminClient.removeGroupFromUser(username, projectManagerGroup.groupId, function (err, removed) {
      var message = '';
      if (err) {
        message = err;
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
      if(removed) {
        message = 'Project Manager has been removed.';
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
    });
  }
});

router.post('/remove_user', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin){
    var adminClient = cinco.client(req.session.user.cinco_keys);
    var username = req.body.form_lfid;
    adminClient.removeUser(username, function (err, removed) {
      var message = '';
      if (err) {
        message = err;
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
      if(removed) {
        message = 'User has been removed.';
        adminClient.getAllUsers(function (err, users, groups) {
          return res.render('admin', { message: message, users: users, groups:groups });
        });
      }
    });
  }
});

module.exports = router;
