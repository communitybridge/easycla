var api = require("../lib/api");
var assert = require('assert');
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

  describe('Admin Endoints', function () {
    var adminClient;
    var sampleUserName = randomUserName();

    before(function (done) {
      apiObj.getKeysForLfId("LaneMeyer", function (err, keys) {
        adminClient = apiObj.client(keys);
        adminClient.createUser(sampleUserName, function(err, created) {
          done();
        });
      });
    });

    it('POST user/', function (done) {
      var username = randomUserName();
      adminClient.createUser(username, function (err, created) {
        assert.ifError(err);
        assert(created, "New user with username of " + username + " should have been created");
        done();
      });
    });

    it('GET user/{id}', function (done) {
      adminClient.getUser(sampleUserName, function(err, user) {
        assert.ifError(err);
        assert.equal(user.lfId, sampleUserName, 'Username is not the same as requested');
        assert(user.userId, 'userId property should exist');
        done();
      });
    });

    it('POST user/{id}/group', function (done) {
      var adminGroup = {
        groupId: 1,
        name: 'ADMIN'
      }
      adminClient.addGroupForUser(sampleUserName, adminGroup, function(err, isUpdated, user) {
        assert.ifError(err);
        assert(isUpdated, "User resource should be updated with new group")
        assert.equal(user.lfId, sampleUserName, 'Username is not the same as requested');
        assert.deepEqual(user.groups[0], adminGroup, 'Added group is not the same');
        done();
      });
    });
  });



});