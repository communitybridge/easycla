var api = require("../lib/api");
var assert = require('assert');
var _ = require('underscore');
var randomstring = require('randomstring');


function randomUserName() {
  return randomstring.generate({
    length: 10,
    charset: 'alphabetic'
  }).toLowerCase();
}

describe('api', function () {

  var apiObj = api("http://localhost:5000");

  describe('Properties', function () {
    describe('apiUrlRoot', function () {
      it('The passed in api root parameter should be available on the returned object', function () {
        assert.equal(apiObj.apiRootUrl, "http://localhost:5000/");
      });
    });
  });

  describe('Public Endpoints', function () {
    it('about/version', function (done) {
      apiObj.getVersion(function (err, body) {
        assert.equal(body['Application-Name'], 'CINCO');
        done();
      });
    });
  });

  describe('Trusted Auth Endpoints', function () {
    describe('keysForLfId', function () {
      it('Calling keysForLfId with an lfId returns an object with keys', function (done) {
        apiObj.getKeysForLfId("LaneMeyer", function (err, keys) {
          assert.ifError(err);
          assert.equal(keys.keyId.length, 20, "keyId length should be 20");
          assert.equal(keys.secret.length, 40, "secret length should be 40");
          done();
        });

      });
    });
  });

  describe('Admin Endpoints', function () {
    var adminClient;
    var sampleUserName = randomUserName();

    before(function (done) {
      apiObj.getKeysForLfId("LaneMeyer", function (err, keys) {
        adminClient = apiObj.client(keys);
        adminClient.createUser(sampleUserName, function (err, created) {
          done();
        });
      });
    });

    it('POST users/', function (done) {
      var username = randomUserName();
      adminClient.createUser(username, function (err, created) {
        assert.ifError(err);
        assert(created, "New user with username of " + username + " should have been created");
        done();
      });
    });

    it('GET user/{id}', function (done) {
      adminClient.getUser(sampleUserName, function (err, user) {
        assert.ifError(err);
        assert.equal(user.lfId, sampleUserName, 'Username is not the same as requested');
        assert(user.userId, 'userId property should exist');
        done();
      });
    });

    it('POST user/{id}/group', function (done) {
      var adminGroup = {
        groupId: 2,
        name: 'ADMIN'
      };
      var expected = [{groupId: 1, name: 'USER'}, adminGroup];
      adminClient.addGroupForUser(sampleUserName, adminGroup, function (err, isUpdated, user) {
        assert.ifError(err);
        assert(isUpdated, "User resource should be updated with new group")
        assert.equal(user.lfId, sampleUserName, 'Username is not the same as requested');
        assert.equal(user.groups.length, 2, 'User must have 2 groups');
        assert(_.some(user.groups, function (g) {
          return (g.groupId === adminGroup.groupId) && (g.name === adminGroup.name);
        }));
        assert(_.some(user.groups, function (g) {
          return (g.groupId === 1) && (g.name === 'USER');
        }));
        done();
      });
    });

    it('DELETE users/{id}/group/{groupId}', function (done) {
      var adminGroup = {
        groupId: 2,
        name: 'ADMIN'
      };
      adminClient.addGroupForUser(sampleUserName, adminGroup, function (err, isUpdated, user) {
        assert.ifError(err);
        adminClient.removeGroupFromUser(sampleUserName, adminGroup.groupId, function (err, isUpdated) {
          assert.ifError(err);
          assert(isUpdated);
          adminClient.getUser(sampleUserName, function (err, user) {
            assert.ifError(err);
            assert(!_.some(user.groups, function (g) {
              return g.groupId == adminGroup.groupId;
            }));
            done();
          });
        });
      });
    });

    it('GET usergroups/', function (done) {
      var expected = [{groupId: 1, name: 'USER'}, {groupId: 2, name: 'ADMIN'},
        {groupId: 3, name: 'PROJECT_MANAGER'}];

      adminClient.getAllGroups(function (err, groups) {
        assert.ifError(err);

        _.each(expected, function (eg) {
          var found = _.find(groups, function (g) {
            return (eg.groupId === g.groupId) && (eg.name === g.name);
          });
          assert(found, "Expected group [" + eg + "] not found in returned groups");
        });
        done();
      });
    });

    it('GET users/', function (done) {
      adminClient.getAllUsers(function (err, users, groups) {
        assert.ifError(err);
        assert(users.length >= 2);
        assert(_.some(users, function (u) {
          return u.lfId == 'fvega';
        }));
        assert(_.some(users, function (u) {
          return u.lfId == 'LaneMeyer';
        }));
        done();
      });
    });
  });


  describe('Projects Endpoints', function () {
    var projManagerClient;
    var projUserName = randomUserName();
    var adminClient;

    before(function (done) {
      apiObj.getKeysForLfId("LaneMeyer", function (err, keys) {
        adminClient = apiObj.client(keys);
        adminClient.createUser(projUserName, function (err, created) {
          var projManagementGroup = {
            groupId: 3,
            name: 'PROJECT_MANAGER'
          };
          adminClient.addGroupForUser(projUserName, projManagementGroup, function (err, updated, user) {
            apiObj.getKeysForLfId(projUserName, function (err, keys) {
              projManagerClient = apiObj.client(keys);
              done();
            });
          });
        });
      });


    });

    it('GET /projects', function (done) {
      projManagerClient.getAllProjects(function (err, projects) {
        assert.ifError(err);
        var sampleProj = projects[0];
        assert(sampleProj, "A single project should exist in the returned response array");
        assert(sampleProj.id, "id property should exist");
        assert(sampleProj.name, "name property should exist");
        assert(sampleProj.description, "description property should exist");
        assert(sampleProj.pm, "pm property should exist");
        assert(sampleProj.url, "url property should exist");
        assert(sampleProj.startDate, "startDate property should exist");
        assert(_.contains(['DIRECT_FUNDED', 'INCORPORATED', 'UNSPECIFIED'], sampleProj.type),
            "type should be one of: ['DIRECT_FUNDED','INCORPORATED','UNSPECIFIED']. was: " + sampleProj.type);
        done();
      });
    });

    it('GET /projects/{id}', function (done) {
      projManagerClient.getAllProjects(function (err, projects) {
        assert.ifError(err);
        var id = projects[0].id;
        projManagerClient.getProject(id, function (err, project) {
          assert.ifError(err);
          assert(project)
          done();
        });
      })
    });

    it('GET /projects/{id} 404', function (done) {
      projManagerClient.getAllProjects(function (err, projects) {
        projManagerClient.getProject("not_a_real_id", function (err, project) {
          assert.equal(err.statusCode, 404);
          done();
        });
      })
    });

    it('POST /projects', function (done) {
      var sampleProj = {
        name: 'Sample Project',
        description: 'Sample Project Description',
        pm: projUserName,
        url: 'http://www.sample.org/',
        type: 'DIRECT_FUNDED',
        startDate: '2016-09-26T09:26:36Z'
      };
      projManagerClient.createProject(sampleProj, function (err, created) {
        assert.ifError(err);
        assert(created);
        done();
      });
    });

    it('POST /projects 403 ', function (done) {
      var username = randomUserName();
      adminClient.createUser(username, function (err) {
        assert.ifError(err);
        apiObj.getKeysForLfId(username, function (err, keys) {
          assert.ifError(err);
          var client = apiObj.client(keys);
          client.createProject({}, function (err) {
            assert.equal(err.statusCode, 403);
            done();
          });
        });
      });
    });

    it("POST then GET project", function (done) {
      var sampleProj = {
        name: 'Sample Project',
        description: 'Sample Project Description',
        pm: projUserName,
        url: 'http://www.sample.org/',
        type: 'DIRECT_FUNDED',
        startDate: new Date().toISOString()
      };
      projManagerClient.createProject(sampleProj, function (err, created, id) {
        assert.ifError(err);
        projManagerClient.getProject(id, function (err, returnedProject) {
          assert.ifError(err);
          assert(returnedProject, "A single project should exist in the returned response array");
          assert.equal(returnedProject.name, sampleProj.name);
          assert.equal(returnedProject.description, sampleProj.description);
          assert.equal(returnedProject.pm, sampleProj.pm);
          assert.equal(returnedProject.url, sampleProj.url);
          assert.equal(returnedProject.startDate, sampleProj.startDate);
          done();
        });
      });
    });

    it('DELETE /projects/{id}', function (done) {
      projManagerClient.getMyProjects(function (err, projects) {
        assert.ifError(err);
        var proj = _.last(projects);
        projManagerClient.archiveProject(proj.id, function (err) {
          assert.ifError(err);
          projManagerClient.getProject(proj.id, function (err) {
            assert.equal(err.statusCode, 404);
            done();
          });
        });
      });
    });

    it('PATCH /projects/{id}', function (done) {
      projManagerClient.getMyProjects(function (err, projects) {
        assert.ifError(err);
        var project = _.last(projects);
        var updatedProps = {
          id: project.id,
          name: randomstring.generate({
            length: 20,
            charset: 'alphabetic'
          }),
          description: randomstring.generate({
            length: 200,
            charset: 'alphabetic'
          })
        };
        projManagerClient.updateProject(updatedProps, function (err, updatedProject) {
          assert.ifError(err);

          assert.equal(updatedProject.id, project.id);
          assert.equal(updatedProject.url, project.url);

          assert.equal(updatedProject.name, updatedProps.name);
          assert.equal(updatedProject.description, updatedProps.description);

          done();
        });
      });
    });

    it('GET /project', function (done) {
      projManagerClient.getMyProjects(function (err, projects) {
        assert.ifError(err);
        var sampleProj = projects[0];
        assert(sampleProj, "A single project should exist in the returned response array");
        assert(sampleProj.id, "id property should exist");
        assert(sampleProj.name, "name property should exist");
        assert(sampleProj.description, "description property should exist");
        assert(sampleProj.pm, "pm property should exist");
        assert(sampleProj.url, "url property should exist");
        assert(sampleProj.startDate, "startDate property should exist");
        assert(_.contains(['DIRECT_FUNDED', 'INCORPORATED', 'UNSPECIFIED'], sampleProj.type),
            "type should be one of: ['DIRECT_FUNDED','INCORPORATED','UNSPECIFIED']. was: " + sampleProj.type);
        done();
      });
    });

    describe('Email Aliases Endpoints', function () {
      it('GET /projects/{id}/emailaliases', function (done) {
        projManagerClient.getAllProjects(function (err, projects) {
          assert.ifError(err);
          var id = projects[0].id;
          projManagerClient.getEmailAliases(id, function (err, emailAliases) {
            assert.ifError(err);
            assert(emailAliases)
            done();
          });
        })
      });

      it('GET /projects/{id}/emailaliases 404', function (done) {
        projManagerClient.getAllProjects(function (err, projects) {
          projManagerClient.getEmailAliases("not_a_real_id", function (err, emailAliases) {
            assert.equal(err.statusCode, 404);
            done();
          });
        })
      });

      it('POST /projects/{id}/emailaliases', function (done) {
        var sampleAlias = {
          "address": "ab@cd.org",
          "participants": [
            {
              "address": "foo@bar.com"
            },
            {
              "address": "foo2@bar.com"
            },
            {
              "address": "foo3@bar.com"
            }
          ]
        };
        projManagerClient.getMyProjects(function (err, projects) {
          assert.ifError(err);
          var projectId = projects[0].id;
          projManagerClient.createEmailAliases(projectId, sampleAlias, function (err, created, aliasId) {
            assert.ifError(err);
            assert(created);
            done();
          });
        });
      });

    });


  });
});
