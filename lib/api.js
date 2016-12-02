const request = require('request');
const crypto = require('crypto');
const url = require('url');

const integration_user = process.env['CONSOLE_INTEGRATION_USER'];
const integration_pass = process.env['CONSOLE_INTEGRATION_PASSWORD'];

const errorFromResponse = function (res, message) {
  if (!message) {
    message = "";
  }
  message += ' Status: [' + res.statusCode + '].';
  if (res.body) {
    message = message + 'Response Body: ' + res.body;
  }
  var error = new Error(message);
  error.statusCode = res.statusCode;
  return error;
};

module.exports = function (apiRootUrl) {
  if (!apiRootUrl.endsWith('/')) {
    apiRootUrl = apiRootUrl + '/'
  }
  return {
    apiRootUrl: apiRootUrl,

    getKeysForLfId: function (lfId, next) {
      request.get(apiRootUrl + "auth/trusted/cas/" + lfId, function (err, res, body) {
            if (err) {
              next(err)
            } else if (res.statusCode != 200) {
              next(errorFromResponse(res, 'Unable to get keys for LfId of [' + lfId + '].  '));
            } else {
              body = JSON.parse(body);
              next(null, {lfId: lfId, keyId: body.keyId, secret: body.secret})
            }
          }
      ).auth(integration_user, integration_pass, true);
    },

    getVersion: function (next) {
      request.get(apiRootUrl + 'about/version', function (err, res, body) {
        if (err) {
          next(err);
        } else if (res.statusCode != 200) {
          next(errorFromResponse(res, 'Unable to get platform version.'));
        } else {
          next(null, JSON.parse(body));
        }
      });
    },

    client: function (apiKeys) {
      var version = "1";

      var md5checksum = function () {
        var md5sum = crypto.createHash('md5');
        for (var i = 0; i < arguments.length; i++) {
          md5sum.update(arguments[i]);
        }
        return md5sum.digest('hex');
      };

      var signRequestVersionOne = function (secretKey, method, path, dateString, entityMd5) {
        var stringToSign = method + '\n' + path + '\n' + dateString + '\n';
        if (entityMd5) {
          stringToSign += entityMd5 + '\n';
        }
        stringToSign += version;

        var hmac = crypto.createHmac("sha1", secretKey);
        hmac.write(stringToSign);
        return hmac.digest('base64');
      };

      var makeSignedRequest = function (reqOpts, next) {
        if (reqOpts.body) {
          reqOpts.bodyChecksum = md5checksum(reqOpts.body);
        }
        reqOpts.uri = apiRootUrl + reqOpts.path;
        var fullPath = url.parse(reqOpts.uri).path;
        var now = new Date().toISOString();
        var signature = signRequestVersionOne(apiKeys.secret, reqOpts.method, fullPath,
            now, reqOpts.bodyChecksum);

        reqOpts.headers = {
          'Content-Type': 'application/json',
          'Date': now,
          'Signature-Version': version,
          'Authorization': 'CINCO ' + apiKeys.keyId + ':' + signature
        };
        if (reqOpts.bodyChecksum) {
          reqOpts.headers['Content-MD5'] = reqOpts.bodyChecksum;
        }

        request(reqOpts, next)
      };

      return {
        createUser: function (lfId, next) {
          var body = {
            "lfId": lfId
          };
          var opts = {
            method: 'POST',
            path: 'users/',
            body: JSON.stringify(body)
          };
          makeSignedRequest(opts, function (err, res) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 201) {
              next(null, true)
            } else if (res.statusCode == 409) {
              next(null, false);
            } else {
              next(errorFromResponse(res, 'User with lfId of [' + lfId + '] not created.'));
            }
          });
        },

        getUser: function (id, next) {
          var opts = {
            method: 'GET',
            path: 'users/' + id + '/'
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err, null);
            } else if (res.statusCode == 200) {
              var user = JSON.parse(body);
              next(null, user);
            } else {
              next(errorFromResponse(res, 'User with id of [' + id + '] could not be retrieved'));
            }
          });
        },

        addGroupForUser: function (id, group, next) {
          var opts = {
            method: 'POST',
            path: 'users/' + id + '/group/',
            body: JSON.stringify(group)
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var updatedUser = JSON.parse(body);
              next(null, true, updatedUser)
            } else if (res.statusCode == 204) {
              next(null, false, null);
            } else {
              next(errorFromResponse(res, 'User with id of [' + id + '] could not have group added'));
            }
          });
        },

        getAllGroups: function (next) {
          var opts = {
            method: 'GET',
            path: 'usergroups/'
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var groups = JSON.parse(body);
              next(null, groups);
            } else {
              next(errorFromResponse(res, 'Unable to look up usergroups. '));
            }
          });
        },

        removeGroupFromUser: function (userId, groupId, next) {
          var opts = {
            method: 'DELETE',
            path: 'users/' + userId + '/group/' + groupId + '/'
          };
          makeSignedRequest(opts, function (err, res) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 204) {
              next(null, true);
            } else {
              next(errorFromResponse(res, 'Unable to delete group with id of [' + groupId + '] from user with id of [' +
                  userId + '].'));
            }
          });
        },

        removeUser: function (userId, next) {
          var opts = {
            method: 'DELETE',
            path: 'users/' + userId + '/'
          };
          makeSignedRequest(opts, function (err, res) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 204) {
              next(null, true);
            } else {
              next(errorFromResponse(res, 'Unable to delete user with id of [' + userId + '].'));
            }
          });
        },

        getAllUsers: function (next) {
          var opts = {
            method: 'GET',
            path: 'users/'
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var groups = new Array();
              var users = JSON.parse(body);
              for (var i = 0; i < users.length; i++) {
                groups[i] = {
                  lfId: '',
                  isUser: false,
                  isAdmin: false,
                  isProjectManager: false
                };
              }
              for (var i = 0; i < users.length; i++) {
                groups[i].lfId = users[i].lfId;
                for (var j = 0; j < users[i].groups.length; j++) {
                  if (users[i].groups[j].name == 'USER') groups[i].isUser = true;
                  if (users[i].groups[j].name == 'ADMIN') groups[i].isAdmin = true;
                  if (users[i].groups[j].name == 'PROJECT_MANAGER') groups[i].isProjectManager = true;
                }
              }
              next(null, users, groups);
            } else {
              next(errorFromResponse(res, 'Unable to get all users.'));
            }
          });
        },

        getAllProjects: function (next) {
          var opts = {
            method: 'GET',
            path: 'projects/'
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var projects = JSON.parse(body);
              next(null, projects);
            } else {
              next(errorFromResponse(res, 'Unable to get all projects.'));
            }
          });
        },

        getProject: function (projectId, next) {
          var opts = {
            method: 'GET',
            path: 'projects/' + projectId + '/'
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var proj = JSON.parse(body);
              next(null, proj);
            } else {
              next(errorFromResponse(res, 'Unable to get project with id of [' + projectId + ']'));
            }
          });
        },

        createProject: function (project, next) {
          var opts = {
            method: 'POST',
            path: 'projects/',
            body: JSON.stringify(project)
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 201) {
              var obj = JSON.parse(body);
              next(null, true, obj.id);
            } else {
              next(errorFromResponse(res, 'Project not created'), false);
            }
          });
        },

        archiveProject: function (id, next) {
          var opts = {
            method: 'DELETE',
            path: 'projects/' + id + '/'
          };
          makeSignedRequest(opts, function (err, res) {
            if (err) {
              next(err, false);
            } else if (res.statusCode != 204) {
              next(errorFromResponse(res, 'Error while archiving project with id of [' + id + ']'), false);
            } else {
              next(null, true);
            }
          });
        },

        updateProject: function (updatedProperties, next) {
          var body = JSON.stringify(updatedProperties);
          var opts = {
            method: 'PATCH',
            path: 'projects/' + updatedProperties.id + '/',
            body: body
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var updatedProject = JSON.parse(body);
              next(null, updatedProject);
            } else {
              next(errorFromResponse(res, "Unable to Update Project with properties: " + body));
            }
          });
        },

        getMyProjects: function (next) {
          var opts = {
            method: 'GET',
            path: 'project/'
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var projects = JSON.parse(body);
              next(null, projects);
            } else {
              next(errorFromResponse(res, 'Unable to get projects managed by logged in user.'));
            }
          });
        },

        getEmailAliases: function (projectId, next) {
          var opts = {
            method: 'GET',
            path: 'projects/' + projectId + '/emailaliases'
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var emailAliases = JSON.parse(body);
              next(null, emailAliases);
            } else {
              next(errorFromResponse(res, 'Unable to get email aliases from project with id of [' + projectId + ']'));
            }
          });
        },

        createEmailAliases: function (projectId, newAlias, next) {
          var opts = {
            method: 'POST',
            path: 'projects/' + projectId + '/emailaliases',
            body: JSON.stringify(newAlias)
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 201) {
              var aliasId = JSON.parse(body);
              next(null, true, aliasId);
            } else {
              next(errorFromResponse(res, 'Email Aliases not created for project with id of [' + projectId + ']'), false);
            }
          });
        },

        addParticipantToEmailAlias: function (projectId, aliasId, newParticipant, next) {
          var opts = {
            method: 'POST',
            path: 'projects/' + projectId + '/emailaliases/' + aliasId + '/participants/',
            body: JSON.stringify(newParticipant)
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 201) {
              var response = JSON.parse(body);
              next(null, true, response);
            } else {
              next(errorFromResponse(res, 'Unable to add participant for project with id of [' + projectId +
              '] and Email Alias with id of [' + aliasId + ']'), false);
            }
          });
        },

        removeParticipantFromEmailAlias: function (projectId, aliasId, participantTBR, next) {
          var opts = {
            method: 'DELETE',
            path: 'projects/' + projectId + '/emailaliases/' + aliasId + '/participants/' + participantTBR,
          };
          makeSignedRequest(opts, function (err, res) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 204) {
              next(null, true);
            } else {
              next(errorFromResponse(res, 'Unable to remove participant [' + participantTBR + '] for project with id of [' + projectId +
              '] and Email Alias with id of [' + aliasId + ']'), false);
            }
          });
        },

        getMemberCompanies: function (projectId, next) {
          var opts = {
            method: 'GET',
            path: 'projects/' + projectId + '/members'
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var memberCompanies = JSON.parse(body);
              next(null, memberCompanies);
            } else {
              next(errorFromResponse(res, 'Unable to get member companies from project with id of [' + projectId + ']'));
            }
          });
        },

        addMemberToProject: function (projectId, newMember , next) {
          var opts = {
            method: 'POST',
            path: 'projects/' + projectId + '/members/',
            body: JSON.stringify(newMember)
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 201) {
              var obj = JSON.parse(body);
              next(null, true, obj.id);
            } else {
              next(errorFromResponse(res, 'Unable to add member to project with id of [' + projectId + ']'), false);
            }
          });
        },

        removeMemberFromProject: function (projectId, memberId, next) {
          var opts = {
            method: 'DELETE',
            path: 'projects/' + projectId + '/members/' + memberId,
          };
          makeSignedRequest(opts, function (err, res) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 204) {
              next(null, true);
            } else {
              next(errorFromResponse(res, 'Unable to remove member [' + memberId + '] for project with id of [' + projectId + ']'), false);
            }
          });
        },

        getMemberFromProject: function (projectId, memberId, next) {
          var opts = {
            method: 'GET',
            path: 'projects/' + projectId + '/members/' + memberId
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var memberCompany = JSON.parse(body);
              next(null, memberCompany);
            } else {
              next(errorFromResponse(res, 'Unable to get member company with id of [' + memberId + '] from project with id of [' + projectId + ']'));
            }
          });
        },

        updateMember: function (projectId, memberId, updatedProperties , next) {
          var opts = {
            method: 'PATCH',
            path: 'projects/' + projectId + '/members/' + memberId,
            body: JSON.stringify(updatedProperties)
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 200) {
              var updatedMember = JSON.parse(body);
              next(null, true, updatedMember);
            } else {
              next(errorFromResponse(res, 'Unable to update member company with id of [' + memberId + '] from project with id of [' + projectId + ']'), false);
            }
          });
        },

        addContactToMember: function (projectId, memberId, newContact , next) {
          var opts = {
            method: 'POST',
            path: 'projects/' + projectId + '/members/' + memberId + '/contacts/',
            body: JSON.stringify(newContact)
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 201) {
              var obj = JSON.parse(body);
              next(null, true, obj.id);
            } else {
              next(errorFromResponse(res, 'Unable to add contact to member with id of [' + memberId + '] from Project [' + projectId + ']'), false);
            }
          });
        },

        removeContactFromMember: function (projectId, memberId, contactId, next) {
          var opts = {
            method: 'DELETE',
            path: 'projects/' + projectId + '/members/' + memberId + '/contacts/' + contactId,
          };
          makeSignedRequest(opts, function (err, res) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 204) {
              next(null, true);
            } else {
              next(errorFromResponse(res, 'Unable to remove contact [' + contactId + '] from member with id of [' + memberId + '] for project with id of [' + projectId + ']'), false);
            }
          });
        },

        updateContactFromMember: function (projectId, memberId, contactId, updatedContact, next) {
          var body = JSON.stringify(updatedContact);
          var opts = {
            method: 'PUT',
            path: 'projects/' + projectId + '/members/' + memberId + '/contacts/' + contactId,
            body: body
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var contact = JSON.parse(body);
              next(null, true, contact);
            } else {
              next(errorFromResponse(res, "Unable to Update Contact with properties: " + body), false);
            }
          });
        },

        createOrganization: function (organization, next) {
          var opts = {
            method: 'POST',
            path: 'organizations/',
            body: JSON.stringify(organization)
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 201) {
              var obj = JSON.parse(body);
              next(null, true, obj.id);
            } else {
              next(errorFromResponse(res, 'Organization not created'), false);
            }
          });
        },

        getAllOrganizations: function (next) {
          var opts = {
            method: 'GET',
            path: 'organizations/'
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var organizations = JSON.parse(body);
              next(null, organizations);
            } else {
              next(errorFromResponse(res, 'Unable to get all organizations.'));
            }
          });
        },

        getOrganization: function (organizationId, next) {
          var opts = {
            method: 'GET',
            path: 'organizations/' + organizationId + '/'
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var organization = JSON.parse(body);
              next(null, organization);
            } else {
              next(errorFromResponse(res, 'Unable to get organization with id of [' + organizationId + ']'));
            }
          });
        },

        updateOrganization: function (updatedOrganization, next) {
          var body = JSON.stringify(updatedOrganization);
          var opts = {
            method: 'PUT',
            path: 'organizations/' + updatedOrganization.id + '/',
            body: body
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var organization = JSON.parse(body);
              next(null, true, organization);
            } else {
              next(errorFromResponse(res, "Unable to Update Organization with properties: " + body), false);
            }
          });
        },

        getMailingLists: function (projectId, next) {
          var opts = {
            method: 'GET',
            path: 'projects/' + projectId + '/mailinglists'
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var mailingLists = JSON.parse(body);
              next(null, mailingLists);
            } else {
              next(errorFromResponse(res, 'Unable to get mailing lists from project with id of [' + projectId + ']'));
            }
          });
        },

        createMailingList: function (projectId, newMailingList, next) {
          var opts = {
            method: 'POST',
            path: 'projects/' + projectId + '/mailinglists',
            body: JSON.stringify(newMailingList)
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 201) {
              var mailingListId = JSON.parse(body);
              next(null, true, mailingListId.mailinglistId);
            } else {
              next(errorFromResponse(res, 'Mailing List not created for project with id of [' + projectId + ']'), false);
            }
          });
        },

        removeMailingList: function (projectId, mailinglistId, next) {
          var opts = {
            method: 'DELETE',
            path: 'projects/' + projectId + '/mailinglists/' + mailinglistId + '/'
          };
          makeSignedRequest(opts, function (err, res) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 204) {
              next(null, true);
            } else {
              next(errorFromResponse(res, 'Unable to delete mailing list with id of [' + mailinglistId + '] from project with id of [' +
                  projectId + '].'));
            }
          });
        },

        getMailingListFromProject: function (projectId, mailinglistId, next) {
          var opts = {
            method: 'GET',
            path: 'projects/' + projectId + '/mailinglists/' + mailinglistId + '/'
          };
          makeSignedRequest(opts, function (err, res, body) {
            if (err) {
              next(err);
            } else if (res.statusCode == 200) {
              var mailingList = JSON.parse(body);
              next(null, mailingList);
            } else {
              next(errorFromResponse(res, 'Unable to get mailing list with id of [' + mailinglistId + '] from project with id of [' + projectId + ']'));
            }
          });
        }


      };
    }
  };
};
