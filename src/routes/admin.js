if (process.env['NEWRELIC_LICENSE']) require('../../infra/newrelic/newrelic');
var express = require('express');
var passport = require('passport');
var request = require('request');
var multer  = require('multer');
var async = require('async');

var cinco = require("../lib/api");

var router = express.Router();

router.get('/admin', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin) {
    var adminClient = cinco.client(req.session.user.cinco_keys);
    adminClient.getAllUsers(function (err, users, groups) {
      res.render('admin', { users: users, groups: groups, message: req.flash('info') });
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
      var message = username + ' has been activated.';
      if (err) message = err;
      req.flash('info', message);
      return res.redirect('/admin');
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
      }
      else if (created) {
        message = 'User [' + username + '] has been created.';
      }
      else {
        message = 'User [' + username + '] already exists.';
      }
      req.flash('info', message);
      return res.redirect('/admin');
    });
  }
});

router.post('/designate_project_manager_user', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
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
    var message = '';
    adminClient.addGroupForUser(username, userGroup, function(err, isUpdated, user) {
      if (err) message = err;
      adminClient.addGroupForUser(username, projectManagerGroup, function(err, isUpdated, user) {
        if (err) message += err;
        else message = username + ' has been designated as a Project Manager.';
        req.flash('info', message);
        return res.redirect('/admin');
      });
    });
  }
});

router.post('/designate_admin_user', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
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
    var message = '';
    adminClient.addGroupForUser(username, userGroup, function(err, isUpdated, user) {
      if (err) message = err;
      adminClient.addGroupForUser(username, adminGroup, function(err, isUpdated, user) {
        if (err) message += err;
        else message = username + ' has been designated as an Admin.';
        req.flash('info', message);
        return res.redirect('/admin');
      });
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
      }
      else if(removed) {
        message = username + ' has been deactivated.';
      }
      req.flash('info', message);
      return res.redirect('/admin');
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
      }
      else if(removed) {
        message = 'Admin privileges for [' + username + ']  have been removed.';
      }
      req.flash('info', message);
      return res.redirect('/admin');
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
      }
      else if(removed) {
        message = 'Project Manager privileges for [' + username + '] have been removed.';
      }
      req.flash('info', message);
      return res.redirect('/admin');
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
      }
      else if(removed) {
        message = 'User [' + username + '] has been removed.';
      }
      req.flash('info', message);
      return res.redirect('/admin');
    });
  }
});

module.exports = router;
