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

router.get('/create_project', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  res.render('create_project');
});

router.get('/my_projects', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getMyProjects(function (err, myProjects) {
      req.session.myProjects = myProjects;
      res.render('my_projects', {myProjects: myProjects});
    });
  }
});

router.get('/project/:id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.id;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getProject(projectId, function (err, project) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.redirect('/');
      projManagerClient.getEmailAliases(projectId, function (err, emailAliases) {
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
            return res.render('project', {project: project, emailAliases: emailAliases, memberCompanies:memberCompanies});
          });
        });
      });
    });
  }
});

router.get('/archive_project/:id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var id = req.params.id;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.archiveProject(id, function (err) {
      console.log(err);
      return res.redirect('/');
    });
  }
});

router.post('/create_project', require('connect-ensure-login').ensureLoggedIn('/login'), cpUpload, function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    var now = new Date().toISOString();
    var url = req.body.url;
    if(url){
      if (!/^(?:f|ht)tps?\:\/\//.test(url)) url = "http://" + url;
    }
    var logoFileName = "";
    var agreementFileName = "";
    if(req.files){
      if(req.files.logo) logoFileName = req.files.logo[0].originalname;
      if(req.files.agreement) agreementFileName = req.files.agreement[0].originalname;
    }
    var newProject = {
      name: req.body.project_name,
      description: req.body.project_description,
      pm: req.session.user.user,
      url: url,
      startDate: now,
      logoRef: logoFileName,
      agreementRef: agreementFileName,
      type: req.body.project_type
    };
    projManagerClient.createProject(newProject, function (err, created, projectId) {
      var isNewAlias = req.body.isNewAlias;
      isNewAlias = (isNewAlias == "true");
      if(isNewAlias){
        var newAlias = JSON.parse(req.body.newAlias);
        async.forEach(newAlias, function (eachAlias, callback){
          projManagerClient.createEmailAliases(projectId, eachAlias, function (err, created, aliasId) {
            callback();
          });
        }, function(err) {
          // Email aliases iteration done.
          return res.redirect('/project/' + projectId);
        });
      }
      else{
        return res.redirect('/project/' + projectId);
      }
    });
  }
});

router.get('/edit_project/:projectId', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projectId = req.params.projectId;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getProject(projectId, function (err, project) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.redirect('/');
      project.domain = "";
      if(project.url) project.domain = project.url.replace('http://www.','').replace('https://www.','').replace('http://','').replace('/','');
      projManagerClient.getEmailAliases(projectId, function (err, emailAliases) {
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
            return res.render('edit_project', {project: project, emailAliases: emailAliases, memberCompanies:memberCompanies});
          });
        });
      });
    });
  }
});

router.post('/edit_project/:id', require('connect-ensure-login').ensureLoggedIn('/login'), cpUpload, function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var id = req.params.id;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    // var now = new Date().toISOString();
    var logoFileName = "";
    var agreementFileName = "";
    var url = req.body.url;
    if(url){
      if (!/^(?:f|ht)tps?\:\/\//.test(url)) url = "http://" + url;
    }
    if(req.files.logo) logoFileName = req.files.logo[0].originalname;
    else logoFileName = req.body.old_logoRef;
    if(req.files.agreement) agreementFileName = req.files.agreement[0].originalname;
    else agreementFileName = req.body.old_agreementRef;
    var updatedProps = {
      id: id,
      name: req.body.project_name,
      description: req.body.project_description,
      pm: req.body.creator_pm,
      url: url,
      // startDate: now,
      logoRef: logoFileName,
      agreementRef: agreementFileName,
      type: req.body.project_type
    };
    projManagerClient.updateProject(updatedProps, function (err, updatedProject) {
      console.log(err);
      return res.redirect('/project/' + id);
    });
  }
});

module.exports = router;
