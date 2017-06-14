const async = require('async');
const config = require('config');
const request = require('request');

const Client = require('./client');
const errors = require('./errors');
const signature = require('./signature');

const integration_user = config.get('console.auth.user');
const integration_pass = config.get('console.auth.pass');

let apiRootUrl = config.get('platform.endpoint');
if (!apiRootUrl.endsWith("/")) {
  apiRootUrl = apiRootUrl + "/";
}

module.exports = {
  apiRootUrl: apiRootUrl,

  getKeysForLfId: function (lfId, next) {
    let opts = {
      uri: apiRootUrl + "auth/trusted/cas/" + lfId,
      auth: {
        user: integration_user,
        pass: integration_pass,
        sendImmediately: true
      }
    };

    request.get(opts, function (err, res, body) {
      if (err) {
        next(err)
      } else if (res.statusCode != 200) {
        next(errors.fromResponse(res, 'Unable to get keys for LfId of [' + lfId + '].  '));
      } else {
        body = JSON.parse(body);
        next(null, { lfId, keyId: body.keyId, secret: body.secret })
      }
    });
  },

  getVersion: function (next) {
    request.get(apiRootUrl + 'about/version', function (err, res, body) {
      if (err) {
        next(err);
      } else if (res.statusCode != 200) {
        next(errors.fromResponse(res, 'Unable to get platform version.'));
      } else {
        next(null, JSON.parse(body));
      }
    });
  },

  client: function (apiKeys) {
    const client = new Client(apiKeys);

    function makeSignedRequest(reqOpts, next) {
      if (!reqOpts.uri) {
        reqOpts.uri = apiRootUrl + reqOpts.path;
        delete reqOpts.path;
      }
      client.request(reqOpts, next);
    }

    return {
      createUser: function (lfId, email, next) {
        var body = {
          "lfId": lfId,
          "email": email
        };
        var opts = {
          method: 'POST',
          path: 'users',
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
            next(errors.fromResponse(res, 'User with lfId of [' + lfId + '] not created.'));
          }
        });
      },

      getUser: function (id, next) {
        var opts = {
          method: 'GET',
          path: 'users/' + id
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err, null);
          } else if (res.statusCode == 200) {
            var user = JSON.parse(body);
            next(null, user);
          } else {
            next(errors.fromResponse(res, 'User with id of [' + id + '] could not be retrieved'));
          }
        });
      },

      addRoleToUser: function (id, role, next) {
        var opts = {
          method: 'POST',
          path: 'users/' + id + '/role',
          body: JSON.stringify(role)
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
            next(errors.fromResponse(res, 'User with id of [' + id + '] could not have role [' + role + ']  added'));
          }
        });
      },

      removeRoleFromUser: function (userId, role, next) {
        var opts = {
          method: 'DELETE',
          path: 'users/' + userId + '/role/' + role
        };
        makeSignedRequest(opts, function (err, res) {
          if (err) {
            next(err, false);
          } else if (res.statusCode == 204) {
            next(null, true);
          } else {
            next(errors.fromResponse(res, 'Unable to delete role [' + role + '] from user with id of [' +
                userId + '].'));
          }
        });
      },

      getAllRoles: function (next) {
        var opts = {
          method: 'GET',
          path: 'users/roles'
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var roles = JSON.parse(body);
            next(null, roles);
          } else {
            next(errors.fromResponse(res, 'Unable to look up user roles. '));
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
            next(errors.fromResponse(res, 'Unable to delete user with id of [' + userId + '].'));
          }
        });
      },

      /*
        Projects:
        Resources to expose and manipulate details of projects
       */

      getMyProjects: function (next) {
        var opts = {
          method: 'GET',
          path: 'project'
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var projects = JSON.parse(body);
            next(null, projects);
          } else {
            next(errors.fromResponse(res, 'Unable to get projects managed by logged in user.'));
          }
        });
      },

      getProjectStatuses: function (next) {
        var opts = {
          method: 'GET',
          path: 'project/status'
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var statuses = JSON.parse(body);
            next(null, statuses);
          } else {
            next(errors.fromResponse(res, 'Unable to get map of valid project status values.'));
          }
        });
      },

      getProjectCategories: function (next) {
        var opts = {
          method: 'GET',
          path: 'project/categories'
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var categories = JSON.parse(body);
            next(null, categories);
          } else {
            next(errors.fromResponse(res, 'Unable to get map of valid project category values.'));
          }
        });
      },

      getProjectSectors: function (next) {
        var opts = {
          method: 'GET',
          path: 'project/sectors'
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var sectors = JSON.parse(body);
            next(null, sectors);
          } else {
            next(errors.fromResponse(res, 'Unable to get map of valid project sector values.'));
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
            next(errors.fromResponse(res, 'Unable to get all projects.'));
          }
        });
      },

      createProject: function (project, next) {
        var opts = {
          method: 'POST',
          path: 'projects',
          body: JSON.stringify(project)
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err, false);
          } else if (res.statusCode == 201) {
            var obj = JSON.parse(body);
            next(null, true, obj.id);
          } else {
            next(errors.fromResponse(res, 'Project not created'), false);
          }
        });
      },

      // archiveProject: function (id, next) {
      //   var opts = {
      //     method: 'DELETE',
      //     path: 'projects/' + id
      //   };
      //   makeSignedRequest(opts, function (err, res) {
      //     if (err) {
      //       next(err, false);
      //     } else if (res.statusCode != 204) {
      //       next(errors.fromResponse(res, 'Error while archiving project with id of [' + id + ']'), false);
      //     } else {
      //       next(null, true);
      //     }
      //   });
      // },

      getProject: function (projectId, next) {
        var opts = {
          method: 'GET',
          path: 'projects/' + projectId
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var proj = JSON.parse(body);
            next(null, proj);
          } else {
            next(errors.fromResponse(res, 'Unable to get project with id of [' + projectId + ']'));
          }
        });
      },

      updateProject: function (updatedProperties, next) {
        var body = JSON.stringify(updatedProperties);
        var opts = {
          method: 'PUT',
          path: 'projects/' + updatedProperties.id,
          body: body
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var updatedProject = JSON.parse(body);
            next(null, updatedProject);
          } else {
            next(errors.fromResponse(res, "Unable to Update Project with properties: " + body));
          }
        });
      },

      getProjectConfig: function (projectId, next) {
        var opts = {
          method: 'GET',
          path: 'projects/' + projectId + '/config'
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var config = JSON.parse(body);
            next(null, config);
          } else {
            next(errors.fromResponse(res, 'Unable to get config from project with id of [' + projectId + ']'));
          }
        });
      },

      updateProjectManagers: function (projectId, managers, next) {
        var body = JSON.stringify(managers);
        var opts = {
          method: 'PUT',
          path: 'projects/' + projectId + '/managers',
          body: body,
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var projectConfig = JSON.parse(body);
            next(null, projectConfig);
          } else {
            next(errors.fromResponse(res, "Unable to Update Project Managers with data: " + body));
          }
        });
      },

      /*
        Email Aliases:
        Resources for working with email aliases of projects
       */

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
            next(errors.fromResponse(res, 'Unable to get email aliases from project with id of [' + projectId + ']'));
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
            next(errors.fromResponse(res, 'Email Aliases not created for project with id of [' + projectId + ']'), false);
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
            next(errors.fromResponse(res, 'Unable to add participant for project with id of [' + projectId +
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
            next(errors.fromResponse(res, 'Unable to remove participant [' + participantTBR + '] for project with id of [' + projectId +
            '] and Email Alias with id of [' + aliasId + ']'), false);
          }
        });
      },

      /*
        Projects - Members:
        Resources for getting details about project members
       */

       getProjectMembers: function (projectId, next) {
         var opts = {
           method: 'GET',
           path: 'projects/' + projectId + '/members/'
         };
         makeSignedRequest(opts, function (err, res, body) {
           if (err) {
             next(err);
           } else if (res.statusCode == 200) {
             var memberCompanies = JSON.parse(body);
             next(null, memberCompanies);
           } else {
             next(errors.fromResponse(res, 'Unable to get member companies from project with id of [' + projectId + ']'));
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
            next(errors.fromResponse(res, 'Unable to get member company with id of [' + memberId + '] from project with id of [' + projectId + ']'));
          }
        });
      },

      /*
        Projects - Members - Contacts:
        Resources for getting and manipulating contacts of project members
       */

      getMemberContactRoles: function (next) {
        var opts = {
          method: 'GET',
          path: 'project/members/contacts/types'
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var roles = JSON.parse(body);
            next(null, roles);
          } else {
            next(errors.fromResponse(res, 'Unable to get member contact roles.'));
          }
        });
      },

      getMemberContacts: function (projectId, memberId, next) {
        var opts = {
          method: 'GET',
          path: 'projects/' + projectId + '/members/' + memberId + '/contacts/',
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err, false);
          } else if (res.statusCode == 200) {
            var contacts = JSON.parse(body);
            next(null, contacts);
          } else {
            next(errors.fromResponse(res, 'Unable to get contacts from member with id of [' + memberId + '] from Project [' + projectId + ']'), false);
          }
        });
      },

      addMemberContact: function (projectId, memberId, contactId, contact, next) {
        var opts = {
          method: 'POST',
          path: 'projects/' + projectId + '/members/' + memberId + '/contacts/' + contactId,
          body: JSON.stringify(contact)
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err, false);
          } else if (res.statusCode == 201) {
            var obj = JSON.parse(body);
            next(null, true, obj);
          } else {
            next(errors.fromResponse(res, 'Unable to add contact with id of [' + contactId + '] to member with id of [' + memberId + '] from Project [' + projectId + ']'), false);
          }
        });
      },

      removeMemberContact: function (projectId, memberId, contactId, roleId, next) {
        var opts = {
          method: 'DELETE',
          path: 'projects/' + projectId + '/members/' + memberId + '/contacts/' + contactId + '/roles/' + roleId,
        };
        makeSignedRequest(opts, function (err, res) {
          if (err) {
            next(err, false);
          } else if (res.statusCode == 204) {
            next(null, true);
          } else {
            next(errors.fromResponse(res, 'Unable to remove contact [' + contactId + '] with role [' + roleId + '] from member with id of [' + memberId + '] for project with id of [' + projectId + ']'), false);
          }
        });
      },

      updateMemberContact: function (projectId, memberId, contactId, roleId, contact, next) {
        var opts = {
          method: 'PUT',
          path: 'projects/' + projectId + '/members/' + memberId + '/contacts/' + contactId + '/roles/' + roleId,
          body: JSON.stringify(contact)
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err, false);
          } else if (res.statusCode == 200) {
            var obj = JSON.parse(body);
            next(null, true, obj);
          } else {
            next(errors.fromResponse(res, 'Unable to add contact with id of [' + contactId + '] to member with id of [' + memberId + '] from Project [' + projectId + ']'), false);
          }
        });
      },

      /*
        Organizations:
        Resources to expose and manipulate organizations
       */

      createOrganization: function (organization, next) {
        var opts = {
          method: 'POST',
          path: 'organizations',
          body: JSON.stringify(organization)
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err, false);
          } else if (res.statusCode == 201) {
            var obj = JSON.parse(body);
            next(null, true, obj.id);
          } else {
            next(errors.fromResponse(res, 'Organization not created'), false);
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
            next(errors.fromResponse(res, 'Unable to get organization with id of [' + organizationId + ']'));
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
            next(errors.fromResponse(res, "Unable to Update Organization with properties: " + body), false);
          }
        });
      },

      /*
        Organizations - Contacts:
        Resources for getting and manipulating contacts of organizations
       */

      getOrganizationContactTypes: function (next) {
        var opts = {
          method: 'GET',
          path: 'organizations/contacts/types'
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var contactTypes = JSON.parse(body);
            next(null, contactTypes);
          } else {
            next(errors.fromResponse(res, 'Unable to get organization contact types.'));
          }
        });
      },

      getOrganizationContacts: function (organizationId, next) {
        var opts = {
          method: 'GET',
          path: 'organizations/' + organizationId + '/contacts/'
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var contacts = JSON.parse(body);
            next(null, contacts);
          } else {
            next(errors.fromResponse(res, 'Unable to get contacts from organization with id of.[' + organizationId + ']'));
          }
        });
      },

      createOrganizationContact: function (organizationId, contact, next) {
        var opts = {
          method: 'POST',
          path: 'organizations/' + organizationId + '/contacts/',
          body: JSON.stringify(contact)
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err, false);
          } else if (res.statusCode == 201) {
            var obj = JSON.parse(body);
            next(null, true, obj.id);
          } else {
            next(errors.fromResponse(res, 'Organization not created'), false);
          }
        });
      },

      getOrganizationContact: function (organizationId, contactId, next) {
        var opts = {
          method: 'GET',
          path: 'organizations/' + organizationId + '/contacts/' + contactId
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var contact = JSON.parse(body);
            next(null, contact);
          } else {
            next(errors.fromResponse(res, 'Unable to get contacts from organization with id of.[' + organizationId + ']'));
          }
        });
      },

      updateOrganizationContact: function (organizationId, contactId, contact, next) {
        var opts = {
          method: 'PUT',
          path: 'organizations/' + organizationId + '/contacts/' + contactId,
          body: JSON.stringify(contact)
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err, false);
          } else if (res.statusCode == 200) {
            var contact = JSON.parse(body);
            next(null, true, contact);
          } else {
            next(errors.fromResponse(res, 'Organization not created'), false);
          }
        });
      },

      /*
        Organizations - Projects:
        Resources for getting details about an organizations project membership
       */

       getOrganizationProjectMemberships: function (organizationId, next) {
         var opts = {
           method: 'GET',
           path: 'organizations/' + organizationId + '/projects_member'
         };
         makeSignedRequest(opts, function (err, res, body) {
           if (err) {
             next(err);
           } else if (res.statusCode == 200) {
             var memberships = JSON.parse(body);
             next(null, memberships);
           } else {
             next(errors.fromResponse(res, 'Unable to get project memberships from organization with id of.[' + organizationId + ']'));
           }
         });
       },

      /*
        Mailing Lists:
        Resources for working with mailing lists of projects
       */

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
            next(errors.fromResponse(res, 'Unable to get mailing lists from project with id of [' + projectId + ']'));
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
            next(errors.fromResponse(res, 'Mailing List not created for project with id of [' + projectId + ']'), false);
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
            next(errors.fromResponse(res, 'Unable to delete mailing list with id of [' + mailinglistId + '] from project with id of [' +
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
            next(errors.fromResponse(res, 'Unable to get mailing list with id of [' + mailinglistId + '] from project with id of [' + projectId + ']'));
          }
        });
      },

      addParticipantToMailingList: function (projectId, mailinglistId, newParticipant , next) {
        var opts = {
          method: 'POST',
          path: 'projects/' + projectId + '/mailinglists/' + mailinglistId + '/participants/',
          body: JSON.stringify(newParticipant)
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err, false);
          } else if (res.statusCode == 201) {
            var obj = JSON.parse(body);
            var participantEmail = obj.address;
            next(null, true, participantEmail);
          } else {
            next(errors.fromResponse(res, 'Unable to add participant to mailinglists with id of [' + mailinglistId + '] from Project [' + projectId + ']'), false);
          }
        });
      },

      removeParticipantFromMailingList: function (projectId, mailinglistId, participantEmail, next) {
        var opts = {
          method: 'DELETE',
          path: 'projects/' + projectId + '/mailinglists/' + mailinglistId + '/participants/' + participantEmail
        };
        makeSignedRequest(opts, function (err, res) {
          if (err) {
            next(err, false);
          } else if (res.statusCode == 204) {
            next(null, true);
          } else {
            next(errors.fromResponse(res, 'Unable to delete participant [' + participantEmail + '] from mailing list with id of [' + mailinglistId + '] from project with id of [' +
                projectId + '].'));
          }
        });
      },

      getParticipantsFromMailingList: function (projectId, mailinglistId, next) {
        var opts = {
          method: 'GET',
          path: 'projects/' + projectId + '/mailinglists/' + mailinglistId + '/participants'
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var participants = JSON.parse(body);
            next(null, participants);
          } else {
            next(errors.fromResponse(res, 'Unable to get participants from mailing list with id of [' + mailinglistId + '] and project with id of [' + projectId + ']'));
          }
        });
      },

      getMailingListsAndParticipants: function (projectId, next) {
        var opts = {
          method: 'GET',
          path: 'projects/' + projectId + '/mailinglists'
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var mailingLists = JSON.parse(body);
            async.forEach(mailingLists, function (eachMailingList, callback){
              eachMailingList.participants = "";
              if(eachMailingList.memberCount > 0){
                var mailinglistId = eachMailingList.name;
                var optsParticipants = {
                  method: 'GET',
                  path: 'projects/' + projectId + '/mailinglists/' + mailinglistId + '/participants'
                };
                makeSignedRequest(optsParticipants, function (err, res, body) {
                  if (err) {
                    callback();
                  } else if (res.statusCode == 200) {
                    var participants = JSON.parse(body);
                    eachMailingList.participants = participants;
                    callback();
                  } else {
                    callback();
                  }
                });
              }
              else {
                callback();
              }
            }, function(err) {
              // Mailing Lists iteration done.
              next(null, mailingLists);
            });
          } else {
            next(errors.fromResponse(res, 'Unable to get mailing lists & participants from project with id of [' + projectId + ']'));
          }
        });
      },

      /*
        Users:
        Resources to manage internal LF users and roles
       */

      getAllUsers: function (next) {
        var opts = {
          method: 'GET',
          path: 'users/'
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var users = JSON.parse(body);
            next(null, users);
          } else {
            next(errors.fromResponse(res, 'Unable to get all users.'));
          }
        });
      },

      getUser: function (userId, next) {
        var opts = {
          method: 'GET',
          path: 'users/' + userId,
        };
        makeSignedRequest(opts, function (err, res, body) {
          if (err) {
            next(err);
          } else if (res.statusCode == 200) {
            var user = JSON.parse(body);
            next(null, user);
          } else {
            next(errors.fromResponse(res, 'Unable to get user with id [' + userId + '].'));
          }
        });
      },

    };
  }
}
