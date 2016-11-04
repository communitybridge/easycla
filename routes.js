var express = require('express');
var passport = require('passport');
var request = require('request');
var multer  = require('multer');
var async = require('async');

var dummy_data = require('./dummy_db/dummy_data');
var cinco_api = require("./lib/api");


var storage = multer.diskStorage({
  destination: function (req, file, cb) {
    cb(null, 'public/uploads/logos')
  },
  filename: function (req, file, cb) {
    cb(null, file.originalname)
  }
});
var upload = multer({ storage: storage });
var cpUpload = upload.fields([{ name: 'logo', maxCount: 1 }, { name: 'agreement', maxCount: 1 }]);

var storageLogoCompany = multer.diskStorage({
  destination: function (req, file, cb) {
    cb(null, 'public/uploads/logos')
  },
  filename: function (req, file, cb) {
    cb(null, file.originalname)
  }
});
var uploadLogoCompany = multer({ storage: storageLogoCompany });
var cpUploadLogoCompany = uploadLogoCompany.fields([{ name: 'logoCompany', maxCount: 1 }]);

var router = express.Router();

const integration_user = process.env['CONSOLE_INTEGRATION_USER'];
const integration_pass = process.env['CONSOLE_INTEGRATION_PASSWORD'];

var hostURL = process.env['CINCO_SERVER_URL'];
if(process.argv[2] == 'dev') hostURL = 'http://localhost:5000';
if(!hostURL.startsWith('http')) hostURL = 'http://' + hostURL;
console.log("hostURL: " + hostURL);

var cinco = cinco_api(hostURL);

router.get('/', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  // res.render('homepage');
  res.redirect('/all_projects');
});

router.get('/angular', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  res.render('angular');
});

router.get('/logout', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  req.session.user = '';
  req.session.destroy();
  req.logout();
  res.redirect('/');
});

router.get('/login', function(req,res) {
  res.render('login');
});

router.get('/404', function(req,res) {
  res.render('404', { lfid: "" });
});

router.get('/401', function(req,res) {
  res.render('401', { lfid: "" });
});

router.get('/login_cas', function(req, res, next) {
  passport.authenticate('cas', function (err, user, info) {
    if (err) return next(err);
    if(user)
    {
      req.session.user = user;
    }
    if (!user) {
      req.session.destroy();
      return res.redirect('/login');
    }
    req.logIn(user, function (err) {
      if (err) return next(err);
      var lfid = req.session.user.user;
      cinco.getKeysForLfId(lfid, function (err, keys) {
        if(keys){
          req.session.user.cinco_keys = keys;
          var adminClient = cinco.client(keys);
          adminClient.getUser(lfid, function(err, user) {
            if(user){
              req.session.user.cinco_groups = JSON.stringify(user.groups);
              req.session.user.isAdmin = false;
              req.session.user.isUser = false;
              req.session.user.isProjectManager = false;
              if(user.groups)
              {
                if(user.groups.length > 0){
                  for(var i = 0; i < user.groups.length; i ++)
                  {
                    if(user.groups[i].name == "ADMIN") req.session.user.isAdmin = true;
                    if(user.groups[i].name == "USER") req.session.user.isUser = true;
                    if(user.groups[i].name == "PROJECT_MANAGER") req.session.user.isProjectManager = true;
                  }
                }
              }
              if( (req.session.user.isAdmin || req.session.user.isProjectManager) && (req.session.user.isUser)) {
                var projManagerClient = cinco.client(req.session.user.cinco_keys);
                projManagerClient.getAllProjects(function (err, projects) {
                  req.session.projects = projects;
                  projManagerClient.getMyProjects(function (err, myProjects) {
                    req.session.myProjects = myProjects;
                    return res.redirect('/');
                  });
                });

              }
              else {
                req.session.destroy();
                return res.render('401', { lfid: lfid }); // User unauthorized.
              }
            }
            else {
              return res.redirect('/login');
            }
          });
        }
        if(err){
          req.session.destroy();
          console.log("getKeysForLfId err: " + err);
          if(err.statusCode == 404) return res.render('404', { lfid: lfid }); // Returned if a user with the given id is not found
          if(err.statusCode == 401) return res.render('401', { lfid: lfid }); // Unable to get keys for lfid given. User unauthorized.
        }
      });
    });
  })(req, res, next);
});

router.get('/profile', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  var adminClient = cinco.client(req.session.user.cinco_keys);
  var lfid = req.session.user.user;
  adminClient.getUser(lfid, function(err, user) {
    if(user){
      req.session.user.cinco_groups = JSON.stringify(user.groups);
      res.render('profile');
    }
    else {
      res.render('profile');
    }
  });
});

router.get('/create_project', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  res.render('create_project');
});

router.get('/add_company', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  res.render('add_company');
});

router.get('/add_company/:project_id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  var projectId = req.params.project_id;
  res.render('add_company', {projectId: projectId});
});

router.post('/add_company', require('connect-ensure-login').ensureLoggedIn('/login'), cpUploadLogoCompany, function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    var projectId = req.body.project_id;
    var now = new Date().toISOString();
    var logoCompanyFileName = "";
    if(req.files){
      if(req.files.logoCompany) logoCompanyFileName = req.files.logoCompany[0].originalname;
    }
    //country code must be exactly 2 Alphabetic characters or null
    var headquartersCountry = req.body.headquarters_country;
    if(headquartersCountry == "") headquartersCountry = null;

    var billingCountry = req.body.billing_country;
    if(billingCountry == "") billingCountry = null;

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
            thoroughfare: req.body.headquarters_address_line_1 + " /// " + req.body.headquarters_address_line_2
          }
        },
        {
          type: "BILL",
          address: {
            country: billingCountry,
            administrativeArea: req.body.billing_state,
            localityName: req.body.billing_city,
            postalCode: req.body.billing_zip_code,
            thoroughfare: req.body.billing_address_line_1 + " /// " + req.body.billing_address_line_2
          }
        }
      ],
      logoRef : logoCompanyFileName
    }
    projManagerClient.createOrganization(newOrganization, function (err, created, organizationId) {
      console.log("organizationId: ", organizationId);
      if(created && projectId){
        var newMember = {
          orgId: organizationId,
          tier: "PLATINUM",
          startDate: now,
          renewalDate: "2017-10-24T00:00:00.000Z"
        };
        console.log("newMember: " + newMember);
        projManagerClient.addMemberToProject(projectId, newMember, function (err, created, memberId) {
          console.log("memberId: " + memberId);

          // var isNewContact = req.body.isNewContact;
          // isNewContact = (isNewContact == "true");
          // if(isNewContact){
          //   var newContact = JSON.parse(req.body.newContact);
          //   async.forEach(newContact, function (eachContact, callback){
          //     projManagerClient.addMemberToProject(projectId, eachContact, function (err, created, contactId) {
          //       callback();
          //     });
          //   }, function(err) {
          //     // Contact Members iteration done.
          //     return res.redirect('/project/' + projectId);
          //   });
          // }
          // else{
          //   return res.redirect('/project/' + projectId);
          // }

          return res.redirect('/project/' + projectId);
        });
      }
      else{
        return res.redirect('/');
      }
    });
  }
});

router.get('/project', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  dummy_data.findProjectById(req.query.id, function(err, project_data) {
    if(project_data) res.render('project', { project_data: project_data });
    else res.redirect('/');
  });
});

router.get('/mailing', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  dummy_data.findProjectById(req.query.id, function(err, project_data) {
    if(project_data) res.render('mailing', { project_data: project_data });
    else res.redirect('/');
  });
});


router.get('/aliases/:id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var id = req.params.id;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getProject(id, function (err, project) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.redirect('/');
      project.domain = "";
      if(project.url) project.domain = project.url.replace('http://www.','').replace('https://www.','').replace('http://','').replace('/','');
      projManagerClient.getEmailAliases(id, function (err, emailAliases) {
        var isAlias = {
          "isContactAlias": false,
          "isEventsAlias": false,
          "isPrAlias": false,
          "isLegalAlias": false,
          "isMembershipAlias": false,
        };
        if( emailAliases.length >= 1 ){
          for (var i = 0; i < emailAliases.length; i++) {
            var email_alias = emailAliases[i].address;
            var alias_name;
            if(email_alias) alias_name = email_alias.replace(/@.*$/,"");
            if(alias_name == "contact") isAlias.isContactAlias = true;
            if(alias_name == "events") isAlias.isEventsAlias = true;
            if(alias_name == "pr") isAlias.isPrAlias = true;
            if(alias_name == "legal") isAlias.isLegalAlias = true;
            if(alias_name == "membership") isAlias.isMembershipAlias = true;
          }
        }
        return res.render('aliases', {project: project, emailAliases: emailAliases, isAlias:isAlias });
      });
    });
  }
});

router.post('/aliases/:id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    var projectId = req.params.id;
    var email_alias = req.body.email_alias;
    var participant_email = req.body.participant_email;
    var newAlias = {
      "address": email_alias,
      "participants": [
        {
          "address": participant_email
        }
      ]
    };
    projManagerClient.createEmailAliases(projectId, newAlias, function (err, created, aliasId) {
      return res.redirect('/aliases/' + projectId);
    });
  }
});

router.post('/addParticipantToEmailAlias/:projectId/', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    var projectId = req.params.projectId;
    var aliasId = req.body.alias_id;
    var participant_email = req.body.participant_email;
    var newParticipant = {
      "address": participant_email
    };
    projManagerClient.addParticipantToEmailAlias(projectId, aliasId, newParticipant, function (err, created, response) {
      return res.redirect('/aliases/' + projectId);
    });
  }
});

router.get('/removeParticipantFromEmailAlias/:projectId/:aliasId/:participantEmail', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    var projectId = req.params.projectId;
    var aliasId = req.params.aliasId;
    var participantTBR = req.params.participantEmail;
    projManagerClient.removeParticipantFromEmailAlias(projectId, aliasId, participantTBR, function (err, removed) {
      return res.redirect('/aliases/' + projectId);
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
            eachMember.orgName = organization.name;
            eachMember.orgLogoRef = organization.logoRef;
            callback();
          });
        }, function(err) {
          // Member Companies iteration done.
          console.log(memberCompanies);
          return res.render('members', {project: project, memberCompanies:memberCompanies});
        });
      });
    });
  }

});

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

router.get('/all_projects', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getAllProjects(function (err, projects) {
      req.session.projects = projects;
      res.render('all_projects', {projects: projects});
    });
  }
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
              eachMember.orgName = organization.name;
              eachMember.orgLogoRef = organization.logoRef;
              callback();
            });
          }, function(err) {
            // Member Companies iteration done.
            return res.render('project-api', {project: project, emailAliases: emailAliases, memberCompanies:memberCompanies});
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
      return res.redirect('/all_projects');
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

router.get('/edit_project/:id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var id = req.params.id;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getProject(id, function (err, project) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.redirect('/');
      project.domain = "";
      if(project.url) project.domain = project.url.replace('http://www.','').replace('https://www.','').replace('http://','').replace('/','');
      projManagerClient.getEmailAliases(id, function (err, emailAliases) {
        return res.render('edit_project', {project: project, emailAliases: emailAliases });
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
      url: url,
      // startDate: now,
      logoRef: logoFileName,
      agreementRef: agreementFileName,
      type: req.body.project_type
    };
    projManagerClient.updateProject(updatedProps, function (err, updatedProject) {
      console.log(err);
      console.log(updatedProject);
      return res.redirect('/project/' + id);
    });
  }
});

router.get('*', function(req, res) {
    res.redirect('/');
});

module.exports = router;
