"""
Holds the model interfaces that all storage models must implement.
"""

class Project(object): # pylint: disable=too-many-public-methods
    """
    Interface to the Project model.
    """

    def to_dict(self):
        """
        Converts models to dictionaries for JSON serialization.

        :return: A dict representation of the model.
        :rtype: dict
        """
        raise NotImplementedError()

    def save(self):
        """
        Simple abstraction around the supported ORMs to save a model.

        Should also save the Documents tied to this model.
        """
        raise NotImplementedError()

    def load(self, project_id):
        """
        Simple abstraction around the supported ORMs to load a model.
        Should populate the current object and also load all the Documents.

        :param project_id: The project's ID.
        :type project_id: string
        """
        raise NotImplementedError()

    def delete(self):
        """
        Simple abstraction around the supported ORMs to delete a model.

        Should also delete the documents tied to this project.
        """
        raise NotImplementedError()

    def get_project_id(self):
        """
        Getter for the project's ID.

        :return: The project's ID.
        :rtype: string
        """
        raise NotImplementedError()

    def get_project_name(self):
        """
        Getter for the project's name.

        :return: The project's name.
        :rtype: string
        """
        raise NotImplementedError()

    def get_project_individual_documents(self):
        """
        Getter for the project's individual agreement documents.

        :return: The project ICLA documents.
        :rtype: [cla.models.model_interfaces.Document]
        """
        raise NotImplementedError()

    def get_project_corporate_documents(self):
        """
        Getter for the project's corporate agreement documents.

        :return: The project CCLA documents.
        :rtype: [cla.models.model_interfaces.Document]
        """
        raise NotImplementedError()

    def get_project_individual_document(self, revision=None):
        """
        Getter for the project's individual agreement document given a revision number.

        A revision number of None should return the latest revision document.

        :param revision: The revision requested. None for latest revision.
        :type revision: integer
        :return: The project's ICLA document corresponding to the revision requested.
        :rtype: cla.models.model_interfaces.Document
        """
        raise NotImplementedError()

    def get_project_corporate_document(self, revision=None):
        """
        Getter for the project's corporate agreement document by revision.

        A revision number of None should return the latest revision document.

        :param revision: The revision requested. None for latest revision.
        :type revision: integer
        :return: The project CCLA document requested.
        :rtype: cla.models.model_interfaces.Document
        """
        raise NotImplementedError()

    def set_project_id(self, project_id):
        """
        Setter for the project's ID.

        :param project_id: The project's ID.
        :type project_id: string
        """
        raise NotImplementedError()

    def set_project_name(self, project_name):
        """
        Setter for the project's name.

        :param project_name: The project's name.
        :type project_name: string
        """
        raise NotImplementedError()

    def set_project_individual_documents(self, documents):
        """
        Setter for the project's individual agreement documents.

        :param documents: The project's individual documents.
        :type documents: [cla.models.model_interfaces.Document]
        """
        raise NotImplementedError()

    def set_project_corporate_documents(self, documents):
        """
        Setter for the project's corporate agreement documents.

        :param document: The project's corporate documents.
        :type document: [cla.models.model_interfaces.Document]
        """
        raise NotImplementedError()

    def add_project_individual_document(self, document):
        """
        Add a single individual document to this project.

        :param document: The document to add to this project as ICLA.
        :type document: cla.models.model_interfaces.Document
        """
        raise NotImplementedError()

    def add_project_corporate_document(self, document):
        """
        Add a single corporate document to this project.

        :param document: The document to add to this project as CCLA.
        :type document: cla.models.model_interfaces.Document
        """
        raise NotImplementedError()

    def remove_project_individual_document(self, document):
        """
        Removes a single individual document from this project.

        :param document: The ICLA document to remove from this project.
        :type document: cla.models.model_interfaces.Document
        """
        raise NotImplementedError()

    def remove_project_corporate_document(self, document):
        """
        Remove a single corporate document from this project.

        :param document: The CCLA document to remove from this project.
        :type document: cla.models.model_interfaces.Document
        """
        raise NotImplementedError()

    def get_project_repositories(self):
        """
        Getter for the project's repositories.

        :return: The project's repository objects.
        :rtype: [cla.models.model_interfaces.Repository]
        """
        raise NotImplementedError()

    def get_project_agreements(self, agreement_signed=None, agreement_approved=None):
        """
        Getter for the project's agreements.

        :param agreement_signed: Whether or not to filter by signed agreements.
            None = no filter, True = only signed, False = only unsigned.
        :type agreement_signed: boolean
        :param agreement_approved: Whether or not to filter by approved agreements.
            None = no filter, True = only approved, False = only unapproved.
        :type agreement_approved: boolean
        :return: The project's agreement objects.
        :rtype: [cla.models.model_interfaces.Agreement]
        """
        raise NotImplementedError()

    def all(self, project_ids=None):
        """
        Fetches all projects in the CLA system.

        :param project_ids: List of project IDs to retrieve.
        :type project_ids: None or [string]
        :return: A list of project objects.
        :rtype: [cla.models.model_interfaces.Project]
        """
        raise NotImplementedError()

class User(object):
    """
    Interface to the User model.
    """

    def to_dict(self):
        """
        Converts models to dictionaries for JSON serialization.

        :return: A dict representation of the model.
        :rtype: dict
        """
        raise NotImplementedError()

    def save(self):
        """
        Simple abstraction around the supported ORMs to save a model.
        """
        raise NotImplementedError()

    def load(self, user_id):
        """
        Simple abstraction around the supported ORMs to load a model.
        Should populate the current object.

        :param user_id: The user's ID.
        :type user_id: string
        """
        raise NotImplementedError()

    def delete(self):
        """
        Simple abstraction around the supported ORMs to delete a model.
        """
        raise NotImplementedError()

    def get_user_id(self):
        """
        Getter for the user's ID.

        :return: The user's ID.
        :rtype: string
        """
        raise NotImplementedError()

    def get_user_email(self):
        """
        Getter for the user's email address.

        :return: The user's email.
        :rtype: string
        """
        raise NotImplementedError()

    def get_user_name(self):
        """
        Getter for the user's name.

        :return: The user's name.
        :rtype: string
        """
        raise NotImplementedError()

    def get_user_github_id(self):
        """
        Getter for the user's GitHub ID.

        :return: The user's GitHub ID.
        :rtype: integer
        """
        raise NotImplementedError()

    def get_user_ldap_id(self):
        """
        Getter for the user's LDAP ID.

        :return: The user's LDAP ID.
        :rtype: string
        """
        raise NotImplementedError()

    def get_user_organization_id(self):
        """
        Getter for the user's organization ID.

        :return: The user's organization ID.
        :rtype: string
        """
        raise NotImplementedError()

    def set_user_id(self, user_id):
        """
        Setter for the user's ID.

        :param user_id: The ID for this user.
        :type user_id: string
        """
        raise NotImplementedError()

    def set_user_email(self, user_email):
        """
        Setter for the user's email address.

        :param user_email: The email for this user.
        :type user_email: string
        """
        raise NotImplementedError()

    def set_user_name(self, user_name):
        """
        Setter for the user's name.

        :param user_name: The user's name.
        :type user_name: string
        """
        raise NotImplementedError()

    def set_user_organization_id(self, organization_id):
        """
        Setter for the user's organization ID.

        :param organization_id: The user's organization ID.
        :type organization_id: string
        """
        raise NotImplementedError()

    def set_user_github_id(self, user_github_id):
        """
        Setter for the user's GitHub ID.

        :param user_github_id: The user's GitHub ID.
        :type user_github_id: integer
        """
        raise NotImplementedError()

    def set_user_ldap_id(self, user_ldap_id):
        """
        Setter for the user's LDAP ID.

        :param user_ldap_id: The user's LDAP ID.
        :type user_ldap_id: integer
        """
        raise NotImplementedError()

    def get_user_by_email(self, user_email):
        """
        Fetches the user object that matches the email specified.

        :param user_email: The user's email.
        :type user_email: string
        :return: The user object with the matching email address - None if not found.
        :rtype: cla.models.model_interfaces.User | None
        """
        raise NotImplementedError()

    def get_user_by_github_id(self, user_github_id):
        """
        Fetches the user object that matches the GitHub ID specified.

        :param user_github_id: The user's GitHub ID.
        :type user_github_id: integer
        :return: The user object with the GitHub ID, or None if not found.
        :rtype: cla.models.model_interfaces.User | None
        """
        raise NotImplementedError()

    def get_user_agreements(self, project_id=None, agreement_signed=None, agreement_approved=None):
        """
        Fetches the agreements associated with this user.

        :param project_id: Filter for project IDs. None = no filter.
        :type project_id: string | None
        :param agreement_signed: Whether or not to filter by signed agreements.
            None = no filter, True = only signed, False = only unsigned.
        :type agreement_signed: boolean
        :param agreement_approved: Whether or not to filter by approved agreements.
            None = no filter, True = only approved, False = only unapproved.
        :type agreement_approved: boolean
        :return: The agreement objects associated with this user.
        :rtype: [cla.models.model_interfaces.Agreement]
        """
        raise NotImplementedError()

    def get_users_by_organization(self, organization_id):
        """
        Fetches the users associated with an organization.

        :param organization_id: The organization ID to filter users by.
        :type organization_id: string
        :return: The agreement objects associated with this user.
        :rtype: [cla.models.model_interfaces.Agreement]
        """
        raise NotImplementedError()

    def all(self, emails=None):
        """
        Fetches all users in the CLA system.

        :param emails: List of user emails to retrieve.
        :type emails: None or [string]
        :return: A list of user objects.
        :rtype: [cla.models.model_interfaces.User]
        """
        raise NotImplementedError()

class Repository(object):
    """
    Interface to the Repository model.
    """

    def to_dict(self):
        """
        Converts models to dictionaries for JSON serialization.

        :return: A dict representation of the model.
        :rtype: dict
        """
        raise NotImplementedError()

    def save(self):
        """
        Simple abstraction around the supported ORMs to save a model.
        """
        raise NotImplementedError()

    def load(self, repository_id):
        """
        Simple abstraction around the supported ORMs to load a model.
        Should populate the current object.

        :param repository_id: The repository ID of the repo to load.
        :type repository_id: string
        """
        raise NotImplementedError()

    def delete(self):
        """
        Simple abstraction around the supported ORMs to delete a model.
        """
        raise NotImplementedError()

    def get_repository_id(self):
        """
        Getter for a repository's ID.

        :return: The repository's ID.
        :rtype: string
        """
        raise NotImplementedError()

    def get_repository_project_id(self):
        """
        Getter for a repository's project ID.

        :return: The repository's project ID.
        :rtype: string
        """
        raise NotImplementedError()

    def get_repository_name(self):
        """
        Getter for a repository's name.

        :return: The repository's name.
        :rtype: string
        """
        raise NotImplementedError()

    def get_repository_type(self):
        """
        Getter for a repository's type ('github', 'gerrit', etc).

        :return: The repository type (github, gerrit, etc).
        :rtype: string
        """
        raise NotImplementedError()

    def get_repository_url(self):
        """
        Getter for a repository's accessible url.

        :return: The repository's accessible url.
        :rtype: string
        """
        raise NotImplementedError()

    def get_repository_external_id(self):
        """
        Getter for a repository's external ID. What the repository provider IDs
        this repository with.

        :return: The repository's external ID.
        :rtype: string
        """
        raise NotImplementedError()

    def set_repository_id(self, repo_id):
        """
        Setter for a repository ID.

        :param repo_id: The repo's ID.
        :type repo_id: string
        """
        raise NotImplementedError()

    def set_repository_project_id(self, project_id):
        """
        Setter for a repository's project ID.

        :param project_id: The repo's project ID.
        :type project_id: string
        """
        raise NotImplementedError()

    def set_repository_name(self, name):
        """
        Setter for a repository's name.

        :param name: The new repository name.
        :type name: string
        """
        raise NotImplementedError()

    def set_repository_type(self, repo_type):
        """
        Setter for a repository's type ('github', 'gerrit', etc).

        :param repo_type: The repository type ('github', 'gerrit', etc).
        :type repo_type: string
        """
        raise NotImplementedError()

    def set_repository_url(self, repository_url):
        """
        Setter for a repository's accessible url.

        :param repository_url: The repository url.
        :type repository_url: string
        """
        raise NotImplementedError()

    def set_repository_external_id(self, repository_external_id):
        """
        Setter for a repository's external ID. What the repository provider IDs
        this repository with.

        :param repository_external_id: The repository external ID.
        :type repository_external_id: string
        """
        raise NotImplementedError()

    def get_repository_by_external_id(self, repository_external_id, repository_type):
        """
        Loads the repository object based on the external ID specified.

        :param repository_external_id: The ID given to the repository by the external provider.
        :type repository_external_id: string
        :param repository_type: The type of repository (GitHub, Gerrit, etc).
        :type repository_type: string
        """
        raise NotImplementedError()

    def all(self, ids=None):
        """
        Fetches all repositories in the CLA system.

        :param ids: List of repository IDs to retrieve.
        :type ids: None or [string]
        :return: A list of repository objects.
        :rtype: [cla.models.model_interfaces.Repository]
        """
        raise NotImplementedError()

class Agreement(object): # pylint: disable=too-many-public-methods
    """
    Interface to the Agreement model.
    """

    def to_dict(self):
        """
        Converts models to dictionaries for JSON serialization.

        :return: A dict representation of the model.
        :rtype: dict
        """
        raise NotImplementedError()

    def save(self):
        """
        Simple abstraction around the supported ORMs to save a model.
        """
        raise NotImplementedError()

    def load(self, agreement_id):
        """
        Simple abstraction around the supported ORMs to load a model.
        Populates the current object.

        :param agreement_id: The agreement ID of the repo to load.
        :type agreement_id: string
        """
        raise NotImplementedError()

    def delete(self):
        """
        Simple abstraction around the supported ORMs to delete a model.
        """
        raise NotImplementedError()

    def get_agreement_id(self):
        """
        Getter for an agreement's ID.

        :return: The agreement's ID.
        :rtype: string
        """
        raise NotImplementedError()

    def get_agreement_project_id(self):
        """
        Getter for an agreement's project ID.

        :return: The agreement's project ID.
        :rtype: string
        """
        raise NotImplementedError()

    def get_agreement_document_revision(self):
        """
        Getter for an agreement's document revision.

        :return: The agreement's document revision.
        :rtype: string
        """
        raise NotImplementedError()

    def get_agreement_reference_id(self):
        """
        Getter for an agreement's user or organization ID, depending on the type
        of agreement this is (individual or corporate).

        :return: The agreement's user or organization ID.
        :rtype: string
        """
        raise NotImplementedError()

    def get_agreement_reference_type(self):
        """
        Getter for an agreement's reference type - could be 'user' or 'organization'.

        :return: The agreement's reference type.
        :rtype: string
        """
        raise NotImplementedError()

    def get_agreement_type(self):
        """
        Getter for an agreement's type ('cla' or 'dco').

        :return: The agreement type (cla or dco)
        :rtype: string
        """
        raise NotImplementedError()

    def get_agreement_signed(self):
        """
        Getter for an agreement's signed status. True for signed, False otherwise.

        :return: The agreement's signed status. True if signed, False otherwise.
        :rtype: boolean
        """
        raise NotImplementedError()

    def get_agreement_approved(self):
        """
        Getter for an agreement's approval status. True is approved, False otherwise.

        :return: The agreement's approval status. True is approved, False otherwise.
        :rtype: boolean
        """
        raise NotImplementedError()

    def get_agreement_sign_url(self):
        """
        Getter for an agreement's signing URL. The URL the user has to visit in
        order to sign the agreement.

        :return: The agreement's signing URL.
        :rtype: string
        """
        raise NotImplementedError()

    def get_agreement_return_url(self):
        """
        Getter for an agreement's return URL. The URL the user gets sent to after signing.

        :return: The agreement's return URL.
        :rtype: string
        """
        raise NotImplementedError()

    def get_agreement_callback_url(self):
        """
        Getter for an agreement's callback URL. The URL that the signing service provider should
        hit once signature has been confirmed.

        :return: The agreement's callback URL.
        :rtype: string
        """
        raise NotImplementedError()

    def set_agreement_id(self, agreement_id):
        """
        Setter for an agreement ID.

        :param agreement_id: The agreement's ID.
        :type agreement_id: string
        """
        raise NotImplementedError()

    def set_agreement_project_id(self, project_id):
        """
        Setter for an agreement's project ID.

        :param project_id: The agreement's project ID.
        :type project_id: string
        """
        raise NotImplementedError()

    def set_agreement_document_revision(self, document_revision):
        """
        Setter for an agreement's document revision.

        :param document_revision: The agreement's document revision.
        :type document_revision: string
        """
        raise NotImplementedError()

    def set_agreement_reference_id(self, reference_id):
        """
        Setter for an agreement's reference ID.

        :param reference_id: The agreement's reference ID.
        :type reference_id: string
        """
        raise NotImplementedError()

    def set_agreement_reference_type(self, reference_type):
        """
        Setter for an agreement's reference type.

        :param reference_type: The agreement's reference type ('user' or 'organization').
        :type reference_type: string
        """
        raise NotImplementedError()

    def set_agreement_type(self, agreement_type):
        """
        Setter for an agreement's type ('cla' or 'dco').

        :param agreement_type: The agreement type ('cla' or 'dco').
        :type agreement_type: string
        """
        raise NotImplementedError()

    def set_agreement_signed(self, signed):
        """
        Setter for an agreement's signed status.

        :param signed: Signed status. True for signed, False otherwise.
        :type signed: bool
        """
        raise NotImplementedError()

    def set_agreement_approved(self, approved):
        """
        Setter for an agreement's approval status.

        :param approved: Approved status. True for approved, False otherwise.
        :type approved: bool
        """
        raise NotImplementedError()

    def set_agreement_sign_url(self, sign_url):
        """
        Setter for an agreement's signing URL. Optional on agreement creation.
        The signing provider's request_signatures() method will populate the field.

        :param sign_url: The URL the user must visit in order to sign the agreement.
        :type sign_url: string
        """
        raise NotImplementedError()

    def set_agreement_return_url(self, return_url):
        """
        Setter for an agreement's return URL. Optional on agreement creation.

        If this value is not set, the CLA system will do it's best to redirect the user to the
        appropriate location once signing is complete (project or repository page).

        :param return_url: The URL the user will be redirected to once signing is complete.
        :type return_url: string
        """
        raise NotImplementedError()

    def set_agreement_callback_url(self, callback_url):
        """
        Setter for an agreement's callback URL. Optional on agreement creation.

        If this value is not set, the signing service provider will not fire a callback request
        when the user's signature has been confirmed.

        :param callback_url: The URL that will hit once the user has signed.
        :type callback_url: string
        """
        raise NotImplementedError()

    def get_agreements_by_reference(self, reference_id, reference_type, # pylint: disable=too-many-arguments
                                    project_id=None,
                                    agreement_signed=None,
                                    agreement_approved=None):
        """
        Simple abstraction around the supported ORMs to get a user's or
        orgnanization's agreements.

        :param reference_id: The reference ID (user_id or organization_id) for
            whom we'll be fetching agreements.
        :type reference_id: string
        :param reference_type: The reference type ('user' or 'organization') for
            whom we'll be fetching agreements.
        :type reference_id: string
        :param project_id: The project ID to filter by. None will not apply any filters.
        :type project_id: string or None
        :param agreement_signed: Whether or not to return only signed/unsigned agreements.
            None will not apply any filters for signed.
        :type agreement_signed: bool or None
        :param agreement_approved: Whether or not to return only approved/unapproved agreements.
            None will not apply any filters for approved.
        :type agreement_approved: bool or None
        :return: List of agreements.
        :rtype: [cla.models.model_interfaces.Agreement]
        """
        raise NotImplementedError()

    def get_agreements_by_project(self, project_id,
                                  agreement_signed=None,
                                  agreement_approved=None):
        """
        Simple abstraction around the supported ORMs to get a project's agreements.

        :param project_id: The project ID we'll be fetching agreements for.
        :type project_id: string
        :param agreement_signed: Whether or not to return only signed/unsigned agreements.
            None will not apply any filters for signed.
        :type agreement_signed: bool or None
        :param agreement_approved: Whether or not to return only approved/unapproved agreements.
            None will not apply any filters for approved.
        :type agreement_approved: bool or None
        :return: List of agreements.
        :rtype: [cla.models.model_interfaces.Agreement]
        """
        raise NotImplementedError()

    def all(self, ids=None):
        """
        Fetches all agreements in the CLA system.

        :param ids: List of agreement IDs to retrieve.
        :type ids: None or [string]
        :return: A list of agreement objects.
        :rtype: [cla.models.model_interfaces.Agreement]
        """
        raise NotImplementedError()

class Organization(object): # pylint: disable=too-many-public-methods
    """
    Interface to the Organization model.
    """

    def to_dict(self):
        """
        Converts models to dictionaries for JSON serialization.

        :return: A dict representation of the model.
        :rtype: dict
        """
        raise NotImplementedError()

    def save(self):
        """
        Simple abstraction around the supported ORMs to save a model.
        """
        raise NotImplementedError()

    def load(self, organization_id):
        """
        Simple abstraction around the supported ORMs to load a model.
        Populates the current object.

        :param organization_id: The ID of the organization to load.
        :type organization_id: string
        """
        raise NotImplementedError()

    def delete(self):
        """
        Simple abstraction around the supported ORMs to delete a model.
        """
        raise NotImplementedError()

    def get_organization_id(self):
        """
        Getter for an organization's ID.

        :return: The organization's ID.
        :rtype: string
        """
        raise NotImplementedError()

    def get_organization_name(self):
        """
        Getter for an organization's name.

        :return: The organization's name.
        :rtype: string
        """
        raise NotImplementedError()

    def get_organization_whitelist(self):
        """
        Getter for an organization's whitelisted domain names.

        :return: The organization's configured whitelist of domains.
        :rtype: [string]
        """
        raise NotImplementedError()

    def get_organization_exclude_patterns(self):
        """
        Getter for an organization's exclude regex patterns.

        :return: The organization's configured exclude patterns.
        :rtype: [string]
        """
        raise NotImplementedError()

    def get_organization_agreements(self,
                                    agreement_type=None,
                                    agreement_signed=None,
                                    agreement_corporate=None):
        """
        Simple abstraction around the supported ORMs to fetch agreements for
        an organization.

        :param agreement_type: Filters the agreements by type ('cla' or 'dco')
            None will return all available agreement types.
        :type agreement_type: string or None
        :param agreement_signed: Whether or not to return only signed/unsigned agreements.
            None will return all agreements for this organization.
        :type agreement_signed: bool or None
        :param agreement_corporate: Filter for corporate vs individual agreement.
            True will return just corporate CLAs.
            False will return just individual CLAs.
            None will return both individual and corporate.
        :type agreement_corporate: bool or None
        :return: The filtered agreements for this organization.
        :rtype: [cla.models.model_interfaces.Agreement]
        """
        raise NotImplementedError()

    def set_organization_id(self, organization_id):
        """
        Setter for an organization ID.

        :param organization_id: The organization's ID.
        :type organization_id: string
        """
        raise NotImplementedError()

    def set_organization_name(self, organization_name):
        """
        Setter for an organization's name.

        :param organization_name: The name of the organization.
        :type organization_name: string
        """
        raise NotImplementedError()

    def set_organization_whitelist(self, whitelist):
        """
        Setter for an organization's whitelisted domain names.

        :param whitelist: The list of domain names to mark as safe.
            Example: ['ibm.com', 'ibm.ca']
        :type whitelist: list of strings
        """
        raise NotImplementedError()

    def add_organization_whitelist(self, whitelist_item):
        """
        Adds another entry in the list of whitelisted domain names.
        Does not query the DB - save() will take care of that.

        :param whitelist_item: A domain name to add to the whitelist of this organization.
        :type whitelist_item: string
        """
        raise NotImplementedError()

    def remove_organization_whitelist(self, whitelist_item):
        """
        Removes an entry from the list of whitelisted domain names.
        Does not query the DB - save() will take care of that.

        :param whitelist_item: A domain name to remove from the whitelist of this organization.
        :type whitelist_item: string
        """
        raise NotImplementedError()

    def set_organization_exclude_patterns(self, exclude_patterns):
        """
        Setter for an organization's exclude regex patterns.

        :param exclude_patterns: The list of email patterns to exlude from signing.
            Example: ['.*@ibm.co.uk$', '^info.*']
        :type exclude_patterns: list of strings

        :todo: Need to actually test out those examples.
        """
        raise NotImplementedError()

    def add_organization_exclude_pattern(self, exclude_pattern):
        """
        Adds another entry in the list of excluded patterns.
        Does not query the DB - save() will take care of that.

        :param exclude_pattern: A regex string to add to the exluded patterns of this organization.
        :type exclude_pattern: string
        """
        raise NotImplementedError()

    def remove_organization_exclude_pattern(self, exclude_pattern):
        """
        Removes an entry from the list of whitelisted domain names.
        Does not query the DB - save() will take care of that.

        :param exclude_pattern: A regex string to remove from the exluded patterns
            of this organization.
        :type exclude_pattern: string
        """
        raise NotImplementedError()

    def all(self, ids=None):
        """
        Fetches all organizations in the CLA system.

        :param ids: List of organization IDs to retrieve.
        :type ids: None or [string]
        :return: A list of organization objects.
        :rtype: [cla.models.model_interfaces.Organization]
        """
        raise NotImplementedError()

class Document(object):
    """
    Interface to the Document model.

    Save/Load/Delete operations should be done through the Project model.
    """

    def to_dict(self):
        """
        Converts models to dictionaries for JSON serialization.

        :return: A dict representation of the model.
        :rtype: dict
        """
        raise NotImplementedError()

    def get_document_name(self):
        """
        Getter for the document's name.

        :return: The document's name.
        :rtype: string
        """
        raise NotImplementedError()

    def get_document_file_id(self):
        """
        Getter for the document's file ID used as filename for storage.

        :return: The document's file ID.
        :rtype: string
        """
        raise NotImplementedError()

    def get_document_content_type(self):
        """
        Getter for the document's content type.

        :return: The document's content type.
        :rtype: string
        """
        raise NotImplementedError()

    def get_document_content(self):
        """
        Getter for the document's content.

        If content type starts with 'storage+', should utilize the storage service to fetch
        the document content. Otherwise, return the content of this field (assumes document
        content stored in DB).

        :return: The document's content.
        :rtype: string
        """
        raise NotImplementedError()

    def get_document_revision(self):
        """
        Getter for the document's revision number.

        :return: The document's revision number.
        :rtype: integer
        """
        raise NotImplementedError()

    def set_document_name(self, document_name):
        """
        Setter for the document's name.

        :param document_name: The document's name.
        :type document_name: string
        """
        raise NotImplementedError()

    def set_document_file_id(self, document_file_id):
        """
        Setter for the document's file ID that's used as filename in storage.

        :param document_id: The document's file ID.
        :type document_id: string
        """
        raise NotImplementedError()

    def set_document_content_type(self, document_content_type):
        """
        Setter for the document's content type.

        :param document_content_type: The document's content type.
        :type document_content_type: string
        """
        raise NotImplementedError()

    def set_document_content(self, document_content):
        """
        Setter for the document's content.

        If content type starts with 'storage+', should utilize the storage service to save
        the document content. Otherwise, simply store the value provided in the DB.
        The value provided could be a URL (for content type such as 'url+pdf') or
        the raw binary data of the document content.

        NOTE: If content type starts with 'storage+', the value of document_content will
        be base64 encoded. Need to decode before sending to storage provider:

            content = base64.b64decode(document_content)
            cla.utils.get_storage_service().store(self, content)

        NOTE: document_file_id should be used as filename when storing with storage service.
        If document_file_id is None, one needs to be provided before saving (typically a UUID).

        :param document_content: The document's content.
        :type document_content: string
        """
        raise NotImplementedError()

    def set_document_revision(self, revision):
        """
        Setter for the document's revision number.

        :param document_revision: The document's revision number.
        :type document_revision: integer
        """
        raise NotImplementedError()
