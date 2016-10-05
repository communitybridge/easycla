var express = require('express');
var passport = require('passport');
var request = require('request');
var multer  = require('multer');

var dummy_data = require('./dummy_db/dummy_data');
var cinco_api = require("./lib/api");

var router = express.Router();

const integration_user = process.env['CONSOLE_INTEGRATION_USER'];
const integration_pass = process.env['CONSOLE_INTEGRATION_PASSWORD'];

var hostURL = process.env['CINCO_SERVER_URL'];
if(process.argv[2] == 'dev') hostURL = 'http://localhost:5000';
if(!hostURL.startsWith('http')) hostURL = 'http://' + hostURL;
console.log("hostURL: " + hostURL);

var cinco = cinco_api(hostURL);

router.get('/', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  res.render('homepage');
  // res.redirect('/my_projects');
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
                  return res.redirect('/');
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

router.get('/aliases', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  dummy_data.findProjectById(req.query.id, function(err, project_data) {
    if(project_data) res.render('aliases', { project_data: project_data });
    else res.redirect('/');
  });
});

router.get('/members', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  dummy_data.findProjectById(req.query.id, function(err, project_data) {
    if(project_data) res.render('members', { project_data: project_data });
    else res.redirect('/');
  });
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
      console.log(err);
      console.log(myProjects);
      req.session.myProjects = myProjects;
      res.render('my_projects', {myProjects: myProjects});
    });
  }
});

router.get('/project/:id', require('connect-ensure-login').ensureLoggedIn('/login'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var id = req.params.id;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    projManagerClient.getProject(id, function (err, project) {
      // TODO: Create 404 page for when project doesn't exist
      if (err) return res.redirect('/');
      console.log(project);
      return res.render('project-api', {project: project});
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

var storage = multer.diskStorage({
  destination: function (req, file, cb) {
    cb(null, 'public/uploads/logos')
  },
  filename: function (req, file, cb) {
    cb(null, file.originalname)
  }
});

var upload = multer({ storage: storage });

router.post('/create_project', require('connect-ensure-login').ensureLoggedIn('/login'), upload.single('file'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    var now = new Date().toISOString();
    var url = req.body.url;
    if (!/^(?:f|ht)tps?\:\/\//.test(url)) url = "http://" + url;
    var logoFileName = "";
    if(req.file) logoFileName = req.file.originalname;
    var newProject = {
      name: req.body.project_name,
      description: req.body.project_description,
      pm: req.session.user.user,
      url: url,
      startDate: now,
      logoRef: logoFileName,
      type: req.body.project_type
    };
    console.log(newProject);
    projManagerClient.createProject(newProject, function (err, created, id) {
      console.log(id);
      console.log(err);
      return res.redirect('/all_projects');
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
      console.log(project);
      return res.render('edit_project', {project: project});
    });
  }
});

router.post('/edit_project/:id', require('connect-ensure-login').ensureLoggedIn('/login'), upload.single('file'), function(req, res){
  if(req.session.user.isAdmin || req.session.user.isProjectManager){
    var id = req.params.id;
    var projManagerClient = cinco.client(req.session.user.cinco_keys);
    // var now = new Date().toISOString();
    var logoFileName = "";
    var url = req.body.url;
    if (!/^(?:f|ht)tps?\:\/\//.test(url)) url = "http://" + url;
    if(req.file) logoFileName = req.file.originalname;
    else logoFileName = req.body.old_logoRef;
    var updatedProps = {
      id: id,
      name: req.body.project_name,
      description: req.body.project_description,
      url: url,
      // startDate: now,
      logoRef: logoFileName,
      type: req.body.project_type
    };
    console.log(updatedProps);
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
