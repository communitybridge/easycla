"""
Easily access CLA models backed by DynamoDB using pynamodb.
"""

import uuid
import base64
import datetime
from pynamodb.models import Model
from pynamodb.indexes import GlobalSecondaryIndex, AllProjection
from pynamodb.attributes import UTCDateTimeAttribute, \
                                UnicodeAttribute, \
                                BooleanAttribute, \
                                NumberAttribute, \
                                ListAttribute, \
                                JSONAttribute, \
                                MapAttribute
import cla
from cla.models import model_interfaces, key_value_store_interface

def create_database():
    """
    Named "create_database" instead of "create_tables" because create_database
    is expected to exist in all database storage wrappers.
    """
    tables = [RepositoryModel, ProjectModel, AgreementModel, \
              OrganizationModel, UserModel, StoreModel]
    # Create all required tables.
    for table in tables:
        # Wait blocks until table is created.
        table.create_table(wait=True)

def delete_database():
    """
    Named "delete_database" instead of "delete_tables" because delete_database
    is expected to exist in all database storage wrappers.

    WARNING: This will delete all existing table data.
    """
    tables = [RepositoryModel, ProjectModel, AgreementModel, \
              OrganizationModel, UserModel, StoreModel]
    # Delete all existing tables.
    for table in tables:
        if table.exists():
            table.delete_table()

class EmailUserIndex(GlobalSecondaryIndex):
    """
    This class represents a global secondary index for querying users by email.
    """
    class Meta:
        """Meta class for email User index."""
        index_name = 'email-user-index'
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
        # All attributes are projected - not sure if this is necessary.
        projection = AllProjection()

    # This attribute is the hash key for the index.
    user_email = UnicodeAttribute(hash_key=True)

class GitHubUserIndex(GlobalSecondaryIndex):
    """
    This class represents a global secondary index for querying users by GitHub ID.
    """
    class Meta:
        """Meta class for GitHub User index."""
        index_name = 'github-user-index'
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
        # All attributes are projected - not sure if this is necessary.
        projection = AllProjection()

    # This attribute is the hash key for the index.
    user_github_id = NumberAttribute(hash_key=True)

class ProjectRepositoryIndex(GlobalSecondaryIndex):
    """
    This class represents a global secondary index for querying repositories by project ID.
    """
    class Meta:
        """Meta class for project repository index."""
        index_name = 'project-repository-index'
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
        # All attributes are projected - not sure if this is necessary.
        projection = AllProjection()

    # This attribute is the hash key for the index.
    repository_project_id = UnicodeAttribute(hash_key=True)

class ExternalRepositoryIndex(GlobalSecondaryIndex):
    """
    This class represents a global secondary index for querying repositories by external ID.
    """
    class Meta:
        """Meta class for external ID repository index."""
        index_name = 'external-repository-index'
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
        # All attributes are projected - not sure if this is necessary.
        projection = AllProjection()

    # This attribute is the hash key for the index.
    repository_external_id = UnicodeAttribute(hash_key=True)


class ProjectAgreementIndex(GlobalSecondaryIndex):
    """
    This class represents a global secondary index for querying agreements by project ID.
    """
    class Meta:
        """Meta class for reference Agreement index."""
        index_name = 'project-agreement-index'
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
        # All attributes are projected - not sure if this is necessary.
        projection = AllProjection()

    # This attribute is the hash key for the index.
    agreement_project_id = UnicodeAttribute(hash_key=True)

class ReferenceAgreementIndex(GlobalSecondaryIndex):
    """
    This class represents a global secondary index for querying agreements by reference.
    """
    class Meta:
        """Meta class for reference Agreement index."""
        index_name = 'reference-agreement-index'
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
        # All attributes are projected - not sure if this is necessary.
        projection = AllProjection()

    # This attribute is the hash key for the index.
    agreement_reference_id = UnicodeAttribute(hash_key=True)

class BaseModel(Model):
    """
    Base pynamodb model used for all CLA models.
    """
    date_created = UTCDateTimeAttribute(default=datetime.datetime.now())
    date_modified = UTCDateTimeAttribute(default=datetime.datetime.now())
    version = UnicodeAttribute(default='v1') # Schema version.

    def __iter__(self):
        """Used to convert model to dict for JSON-serialized string."""
        for name, attr in self._get_attributes().items():
            if isinstance(attr, ListAttribute):
                values = attr.serialize(getattr(self, name))
                if len(values) < 1:
                    yield name, []
                else:
                    key = list(values[0].keys())[0]
                    yield name, [value[key] for value in values]
            else:
                yield name, attr.serialize(getattr(self, name))

class DocumentModel(MapAttribute):
    """
    Represents a document in the project model.
    """
    document_name = UnicodeAttribute()
    document_file_id = UnicodeAttribute(null=True)
    document_content_type = UnicodeAttribute() # pdf, url+pdf, storage+pdf, etc
    document_content = UnicodeAttribute(null=True) # None if using storage service.
    document_revision = NumberAttribute(default=1)

class Document(model_interfaces.Document):
    """
    ORM-agnostic wrapper for the DynamoDB Document model.
    """
    def __init__(self, # pylint: disable=too-many-arguments
                 document_name=None,
                 document_file_id=None,
                 document_content_type=None,
                 document_content=None,
                 document_revision=None):
        super().__init__()
        self.model = DocumentModel()
        self.model.document_name = document_name
        self.model.document_file_id = document_file_id
        self.model.document_content_type = document_content_type
        self.model.document_content = self.set_document_content(document_content)
        if document_revision is not None:
            self.model.document_revision = document_revision

    def to_dict(self):
        return {'document_name': self.model.document_name,
                'document_file_id': self.model.document_file_id,
                'document_content_type': self.model.document_content_type,
                'document_content': self.model.document_content,
                'document_revision': self.model.document_revision}

    def get_document_name(self):
        return self.model.document_name

    def get_document_file_id(self):
        return self.model.document_file_id

    def get_document_content_type(self):
        return self.model.document_content_type

    def get_document_content(self):
        content_type = self.get_document_content_type()
        if content_type is None:
            cla.log.warning('Empty content type for document - not sure how to retrieve content')
        else:
            if content_type.startswith('storage+'):
                filename = self.get_document_file_id()
                return cla.utils.get_storage_service().retrieve(filename)
        return self.model.document_content

    def get_document_revision(self):
        return self.model.document_revision

    def set_document_name(self, document_name):
        self.model.document_name = document_name

    def set_document_file_id(self, document_file_id):
        self.model.document_file_id = document_file_id

    def set_document_content_type(self, document_content_type):
        self.model.document_content_type = document_content_type

    def set_document_content(self, document_content):
        content_type = self.get_document_content_type()
        if content_type is not None and content_type.startswith('storage+'):
            filename = self.get_document_file_id()
            if filename is None:
                filename = str(uuid.uuid4())
                self.set_document_file_id(filename)
            cla.log.info('Saving document content for %s to %s',
                         self.get_document_name(), filename)
            content = base64.b64decode(document_content)
            cla.utils.get_storage_service().store(filename, content)
        else:
            self.model.document_content = document_content

    def set_document_revision(self, revision):
        self.model.document_revision = revision

class ProjectModel(BaseModel):
    """
    Represents a project in the database.
    """
    class Meta:
        """Meta class for Project."""
        table_name = 'cla_projects'
        host = cla.conf['DATABASE_HOST']
        region = cla.conf['DYNAMO_REGION']
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
    project_id = UnicodeAttribute(hash_key=True)
    project_name = UnicodeAttribute()
    project_individual_documents = ListAttribute(of=DocumentModel, default=[])
    project_corporate_documents = ListAttribute(of=DocumentModel, default=[])

class Project(model_interfaces.Project): # pylint: disable=too-many-public-methods
    """
    ORM-agnostic wrapper for the DynamoDB Project model.
    """
    def __init__(self, project_id=None, project_name=None):
        super(Project).__init__()
        self.model = ProjectModel()
        self.model.project_id = project_id
        self.model.project_name = project_name

    def to_dict(self):
        individual_documents = []
        corporate_documents = []
        for doc in self.model.project_individual_documents:
            document = Document()
            document.model = doc
            individual_documents.append(document.to_dict())
        for doc in self.model.project_corporate_documents:
            document = Document()
            document.model = doc
            corporate_documents.append(document.to_dict())
        project_dict = dict(self.model)
        project_dict['project_individual_documents'] = individual_documents
        project_dict['project_corporate_documents'] = corporate_documents
        return project_dict

    def save(self):
        self.model.save()

    def load(self, project_id):
        try:
            project = self.model.get(project_id)
        except ProjectModel.DoesNotExist:
            raise cla.models.DoesNotExist('Project not found')
        self.model = project

    def delete(self):
        self.model.delete()

    def get_project_id(self):
        return self.model.project_id

    def get_project_name(self):
        return self.model.project_name

    def get_project_individual_documents(self):
        documents = []
        for doc in self.model.project_individual_documents:
            document = Document()
            document.model = doc
            documents.append(document)
        return documents

    def get_project_corporate_documents(self):
        documents = []
        for doc in self.model.project_corporate_documents:
            document = Document()
            document.model = doc
            documents.append(document)
        return documents

    def get_project_individual_document(self, revision=None):
        num_documents = len(self.model.project_individual_documents)
        if num_documents < 1:
            raise cla.models.DoesNotExist('No individual document exists for this project')
        # TODO Need to optimize this on the DB side.
        latest_document = None
        for document in self.model.project_individual_documents:
            if document.document_revision == revision:
                latest_document = document
                break
            if latest_document is None or \
               document.document_revision > latest_document.document_revision:
                latest_document = document
        if latest_document is None:
            raise cla.models.DoesNotExist('Document revision not found')
        document = Document()
        document.model = latest_document
        return document

    def get_project_corporate_document(self, revision=None):
        num_documents = len(self.model.project_corporate_documents)
        if num_documents < 1:
            raise cla.models.DoesNotExist('No corporate document exists for this project')
        # TODO Need to optimize this on the DB side.
        latest_document = None
        for document in self.model.project_corporate_documents:
            if document.document_revision == revision:
                latest_document = document
                break
            if latest_document is None or \
               document.document_revision > latest_document.document_revision:
                latest_document = document
        if latest_document is None:
            raise cla.models.DoesNotExist('Document revision not found')
        document = Document()
        document.model = latest_document
        return document

    def set_project_id(self, project_id):
        self.model.project_id = str(project_id)

    def set_project_name(self, project_name):
        self.model.project_name = project_name

    def add_project_individual_document(self, document):
        self.model.project_individual_documents.append(document.model)

    def add_project_corporate_document(self, document):
        self.model.project_corporate_documents.append(document.model)

    def remove_project_individual_document(self, document):
        new_documents = _remove_project_document(self.model.project_individual_documents,
                                                 document.get_document_revision())
        self.model.project_individual_documents = new_documents

    def remove_project_corporate_document(self, document):
        new_documents = _remove_project_document(self.model.project_corporate_documents,
                                                 document.get_document_revision())
        self.model.project_corporate_documents = new_documents

    def set_project_individual_documents(self, documents):
        self.model.project_individual_documents = documents

    def set_project_corporate_documents(self, documents):
        self.model.project_corporate_documents = documents

    def get_project_repositories(self):
        repository_generator = RepositoryModel.repository_project_index.query(self.get_project_id())
        repositories = []
        for repository_model in repository_generator:
            repository = Repository()
            repository.model = repository_model
            repositories.append(repository)
        return repositories

    def get_project_agreements(self, agreement_signed=None, agreement_approved=None):
        return Agreement().get_agreements_by_project(self.get_project_id(),
                                                     agreement_approved=agreement_approved,
                                                     agreement_signed=agreement_signed)

    def all(self, project_ids=None):
        if project_ids is None:
            projects = self.model.scan()
        else:
            projects = ProjectModel.batch_get(project_ids)
        ret = []
        for project in projects:
            proj = Project()
            proj.model = project
            ret.append(proj)
        return ret

def _remove_project_document(documents, revision):
    # TODO Need to optimize this on the DB side - delete directly from list of records.
    new_documents = []
    found = False
    for document in documents:
        if document.document_revision == revision:
            found = True
            if document.document_content_type.startswith('storage+'):
                cla.utils.get_storage_service().delete(document.document_file_id)
            continue
        new_documents.append(document)
    if not found:
        raise cla.models.DoesNotExist('Document revision not found')
    return new_documents

class UserModel(BaseModel):
    """
    Represents a user in the database.
    """
    class Meta:
        """Meta class for User."""
        table_name = 'cla_users'
        host = cla.conf['DATABASE_HOST']
        region = cla.conf['DYNAMO_REGION']
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
    user_id = UnicodeAttribute(hash_key=True)
    user_email = UnicodeAttribute()
    user_name = UnicodeAttribute(null=True)
    user_organization_id = UnicodeAttribute(null=True)
    user_github_id = NumberAttribute(null=True)
    user_ldap_id = UnicodeAttribute(null=True)
    user_email_index = EmailUserIndex()
    user_github_id_index = GitHubUserIndex()

class User(model_interfaces.User): # pylint: disable=too-many-public-methods
    """
    ORM-agnostic wrapper for the DynamoDB User model.
    """
    def __init__(self, user_email=None, user_github_id=None, user_ldap_id=None):
        super(User).__init__()
        self.model = UserModel()
        self.model.user_email = user_email
        self.model.user_github_id = user_github_id
        self.model.user_ldap_id = user_ldap_id

    def to_dict(self):
        ret = dict(self.model)
        if ret['user_github_id'] == 'null':
            ret['user_github_id'] = None
        if ret['user_ldap_id'] == 'null':
            ret['user_ldap_id'] = None
        return ret

    def save(self):
        self.model.save()

    def load(self, user_id):
        try:
            repo = self.model.get(str(user_id))
        except UserModel.DoesNotExist:
            raise cla.models.DoesNotExist('User not found')
        self.model = repo

    def delete(self):
        self.model.delete()

    def get_user_id(self):
        return self.model.user_id

    def get_user_email(self):
        return self.model.user_email

    def get_user_name(self):
        return self.model.user_name

    def get_user_organization_id(self):
        return self.model.user_organization_id

    def get_user_github_id(self):
        return self.model.user_github_id

    def get_user_ldap_id(self):
        return self.model.user_ldap_id

    def set_user_id(self, user_id):
        self.model.user_id = user_id

    def set_user_email(self, user_email):
        self.model.user_email = user_email

    def set_user_name(self, user_name):
        self.model.user_name = user_name

    def set_user_organization_id(self, organization_id):
        self.model.user_organization_id = organization_id

    def set_user_github_id(self, user_github_id):
        self.model.user_github_id = user_github_id

    def set_user_ldap_id(self, user_ldap_id):
        self.model.user_ldap_id = user_ldap_id

    def get_user_by_email(self, user_email):
        user_generator = self.model.user_email_index.query(user_email)
        for user_model in user_generator:
            user = User()
            user.model = user_model
            return user
        return None

    def get_user_by_github_id(self, user_github_id):
        user_generator = self.model.user_github_id_index.query(user_github_id)
        for user_model in user_generator:
            user = User()
            user.model = user_model
            return user
        return None

    def get_user_agreements(self, project_id=None, agreement_signed=None, agreement_approved=None):
        return Agreement().get_agreements_by_reference(self.get_user_id(), 'user',
                                                       project_id=project_id,
                                                       agreement_approved=agreement_approved,
                                                       agreement_signed=agreement_signed)

    def get_users_by_organization(self, organization_id):
        user_generator = self.model.scan(user_organization_id__eq=str(organization_id))
        users = []
        for user_model in user_generator:
            user = User()
            user.model = user_model
            users.append(user)
        return users

    def all(self, emails=None):
        if emails is None:
            users = self.model.scan()
        else:
            users = UserModel.batch_get(emails)
        ret = []
        for user in users:
            usr = User()
            usr.model = user
            ret.append(usr)
        return ret

class RepositoryModel(BaseModel):
    """
    Represents a repository in the database.
    """
    class Meta:
        """Meta class for Repository."""
        table_name = 'cla_repositories'
        host = cla.conf['DATABASE_HOST']
        region = cla.conf['DYNAMO_REGION']
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
    repository_id = UnicodeAttribute(hash_key=True)
    repository_project_id = UnicodeAttribute()
    repository_name = UnicodeAttribute()
    repository_type = UnicodeAttribute() # Gerrit, GitHub, etc.
    repository_url = UnicodeAttribute()
    repository_external_id = UnicodeAttribute(null=True)
    repository_project_index = ProjectRepositoryIndex()
    repository_external_index = ExternalRepositoryIndex()

class Repository(model_interfaces.Repository):
    """
    ORM-agnostic wrapper for the DynamoDB Repository model.
    """
    def __init__(self, repository_id=None, repository_project_id=None, # pylint: disable=too-many-arguments
                 repository_name=None, repository_type=None, repository_url=None,
                 repository_external_id=None):
        super(Repository).__init__()
        self.model = RepositoryModel()
        self.model.repository_id = repository_id
        self.model.repository_project_id = repository_project_id
        self.model.repository_name = repository_name
        self.model.repository_type = repository_type
        self.model.repository_url = repository_url
        self.model.repository_external_id = repository_external_id

    def to_dict(self):
        return dict(self.model)

    def save(self):
        self.model.save()

    def load(self, repository_id):
        try:
            repo = self.model.get(repository_id)
        except RepositoryModel.DoesNotExist:
            raise cla.models.DoesNotExist('Repository not found')
        self.model = repo

    def delete(self):
        self.model.delete()

    def get_repository_id(self):
        return self.model.repository_id

    def get_repository_project_id(self):
        return self.model.repository_project_id

    def get_repository_name(self):
        return self.model.repository_name

    def get_repository_type(self):
        return self.model.repository_type

    def get_repository_url(self):
        return self.model.repository_url

    def get_repository_external_id(self):
        return self.model.repository_external_id

    def set_repository_id(self, repo_id):
        self.model.repository_id = str(repo_id)

    def set_repository_project_id(self, project_id):
        self.model.repository_project_id = project_id

    def set_repository_name(self, name):
        self.model.repository_name = name

    def set_repository_type(self, repo_type):
        self.model.repository_type = repo_type

    def set_repository_url(self, repository_url):
        self.model.repository_url = repository_url

    def set_repository_external_id(self, repository_external_id):
        self.model.repository_external_id = str(repository_external_id)

    def get_repository_by_external_id(self, repository_external_id, repository_type):
        repository_generator = \
            self.model.repository_external_index.query(str(repository_external_id))
        for repository_model in repository_generator:
            if repository_model.repository_type == repository_type:
                repository = Repository()
                repository.model = repository_model
                return repository
        return None

    def all(self, ids=None):
        if ids is None:
            repositories = self.model.scan()
        else:
            repositories = RepositoryModel.batch_get(ids)
        ret = []
        for repository in repositories:
            repo = Repository()
            repo.model = repository
            ret.append(repo)
        return ret

class AgreementModel(BaseModel): # pylint: disable=too-many-instance-attributes
    """
    Represents an agreement in the database.
    """
    class Meta:
        """Meta class for Agreement."""
        table_name = 'cla_agreements'
        host = cla.conf['DATABASE_HOST']
        region = cla.conf['DYNAMO_REGION']
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
    agreement_id = UnicodeAttribute(hash_key=True)
    agreement_project_id = UnicodeAttribute()
    agreement_document_revision = NumberAttribute()
    agreement_reference_id = UnicodeAttribute()
    agreement_reference_type = UnicodeAttribute()
    agreement_type = UnicodeAttribute(default='cla') # Only CLA/DCO.
    agreement_signed = BooleanAttribute(default=False)
    agreement_approved = BooleanAttribute(default=False)
    agreement_sign_url = UnicodeAttribute(null=True)
    agreement_return_url = UnicodeAttribute(null=True)
    agreement_callback_url = UnicodeAttribute(null=True)
    agreement_project_index = ProjectAgreementIndex()
    agreement_reference_index = ReferenceAgreementIndex()

class Agreement(model_interfaces.Agreement): # pylint: disable=too-many-public-methods
    """
    ORM-agnostic wrapper for the DynamoDB Agreement model.
    """
    def __init__(self, # pylint: disable=too-many-arguments
                 agreement_id=None,
                 agreement_project_id=None,
                 agreement_document_revision=None,
                 agreement_reference_id=None,
                 agreement_reference_type='user',
                 agreement_type=None,
                 agreement_signed=False,
                 agreement_approved=False,
                 agreement_sign_url=None,
                 agreement_return_url=None,
                 agreement_callback_url=None):
        super(Agreement).__init__()
        self.model = AgreementModel()
        self.model.agreement_id = agreement_id
        self.model.agreement_project_id = agreement_project_id
        self.model.agreement_document_revision = agreement_document_revision
        self.model.agreement_reference_id = agreement_reference_id
        self.model.agreement_reference_type = agreement_reference_type
        self.model.agreement_type = agreement_type
        self.model.agreement_signed = agreement_signed
        self.model.agreement_approved = agreement_approved
        self.model.agreement_sign_url = agreement_sign_url
        self.model.agreement_return_url = agreement_return_url
        self.model.agreement_callback_url = agreement_callback_url

    def to_dict(self):
        return dict(self.model)

    def save(self):
        self.model.save()

    def load(self, agreement_id):
        try:
            agreement = self.model.get(agreement_id)
        except AgreementModel.DoesNotExist:
            raise cla.models.DoesNotExist('Agreement not found')
        self.model = agreement

    def delete(self):
        self.model.delete()

    def get_agreement_id(self):
        return self.model.agreement_id

    def get_agreement_project_id(self):
        return self.model.agreement_project_id

    def get_agreement_document_revision(self):
        return self.model.agreement_document_revision

    def get_agreement_type(self):
        return self.model.agreement_type

    def get_agreement_signed(self):
        return self.model.agreement_signed

    def get_agreement_approved(self):
        return self.model.agreement_approved

    def get_agreement_sign_url(self):
        return self.model.agreement_sign_url

    def get_agreement_return_url(self):
        return self.model.agreement_return_url

    def get_agreement_callback_url(self):
        return self.model.agreement_callback_url

    def get_agreement_reference_id(self):
        return self.model.agreement_reference_id

    def get_agreement_reference_type(self):
        return self.model.agreement_reference_type

    def set_agreement_id(self, agreement_id):
        self.model.agreement_id = str(agreement_id)

    def set_agreement_project_id(self, project_id):
        self.model.agreement_project_id = str(project_id)

    def set_agreement_document_revision(self, document_revision):
        self.model.agreement_document_revision = int(document_revision)

    def set_agreement_type(self, agreement_type):
        self.model.agreement_type = agreement_type

    def set_agreement_signed(self, signed):
        self.model.agreement_signed = bool(signed)

    def set_agreement_approved(self, approved):
        self.model.agreement_approved = bool(approved)

    def set_agreement_sign_url(self, sign_url):
        self.model.agreement_sign_url = sign_url

    def set_agreement_return_url(self, return_url):
        self.model.agreement_return_url = return_url

    def set_agreement_callback_url(self, callback_url):
        self.model.agreement_callback_url = callback_url

    def set_agreement_reference_id(self, reference_id):
        self.model.agreement_reference_id = reference_id

    def set_agreement_reference_type(self, reference_type):
        self.model.agreement_reference_type = reference_type

    def get_agreements_by_reference(self, # pylint: disable=too-many-arguments
                                    reference_id,
                                    reference_type,
                                    project_id=None,
                                    agreement_signed=None,
                                    agreement_approved=None):
        # TODO: Optimize this query to use filters properly.
        agreement_generator = self.model.agreement_reference_index.query(reference_id)
        agreements = []
        for agreement_model in agreement_generator:
            if agreement_model.agreement_reference_type != reference_type:
                continue
            if project_id is not None and \
               agreement_model.agreement_project_id != project_id:
                continue
            if agreement_signed is not None and \
               agreement_model.agreement_signed != agreement_signed:
                continue
            if agreement_approved is not None and \
               agreement_model.agreement_approved != agreement_approved:
                continue
            agreement = Agreement()
            agreement.model = agreement_model
            agreements.append(agreement)
        return agreements

    def get_agreements_by_project(self, project_id, agreement_signed=None,
                                  agreement_approved=None):
        agreement_generator = self.model.agreement_project_index.query(project_id)
        agreements = []
        for agreement_model in agreement_generator:
            if agreement_signed is not None and \
               agreement_model.agreement_signed != agreement_signed:
                continue
            if agreement_approved is not None and \
               agreement_model.agreement_approved != agreement_approved:
                continue
            agreement = Agreement()
            agreement.model = agreement_model
            agreements.append(agreement)
        return agreements

    def all(self, ids=None):
        if ids is None:
            agreements = self.model.scan()
        else:
            agreements = AgreementModel.batch_get(ids)
        ret = []
        for agreement in agreements:
            agr = Agreement()
            agr.model = agreement
            ret.append(agr)
        return ret

class OrganizationModel(BaseModel):
    """
    Represents an organization in the database.
    """
    class Meta:
        """Meta class for Organization."""
        table_name = 'cla_organizations'
        host = cla.conf['DATABASE_HOST']
        region = cla.conf['DYNAMO_REGION']
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
    organization_id = UnicodeAttribute(hash_key=True)
    organization_name = UnicodeAttribute()
    organization_whitelist = ListAttribute()
    organization_exclude_patterns = ListAttribute()

class Organization(model_interfaces.Organization): # pylint: disable=too-many-public-methods
    """
    ORM-agnostic wrapper for the DynamoDB Organization model.
    """
    def __init__(self, # pylint: disable=too-many-arguments
                 organization_id=None,
                 organization_name=None,
                 organization_exclude_patterns=None,
                 organization_whitelist=None):
        super(Organization).__init__()
        self.model = OrganizationModel()
        self.model.organization_id = organization_id
        self.model.organization_name = organization_name
        self.model.organization_whitelist = organization_whitelist
        self.model.organization_exclude_patterns = organization_exclude_patterns

    def to_dict(self):
        return dict(self.model)

    def save(self):
        self.model.save()

    def load(self, organization_id):
        try:
            organization = self.model.get(organization_id)
        except OrganizationModel.DoesNotExist:
            raise cla.models.DoesNotExist('Organization not found')
        self.model = organization

    def delete(self):
        self.model.delete()

    def get_organization_id(self):
        return self.model.organization_id

    def get_organization_name(self):
        return self.model.organization_name

    def get_organization_whitelist(self):
        return self.model.organization_whitelist

    def get_organization_exclude_patterns(self):
        return self.model.organization_exclude_patterns

    def set_organization_id(self, organization_id):
        self.model.organization_id = organization_id

    def set_organization_name(self, organization_name):
        self.model.organization_name = str(organization_name)

    def set_organization_whitelist(self, whitelist):
        self.model.organization_whitelist = [str(wl) for wl in whitelist]

    def add_organization_whitelist(self, whitelist_item):
        if self.model.organization_whitelist is None:
            self.model.organization_whitelist = [str(whitelist_item)]
        else:
            self.model.organization_whitelist.append(str(whitelist_item))

    def remove_organization_whitelist(self, whitelist_item):
        if str(whitelist_item) in self.model.organization_whitelist:
            self.model.organization_whitelist.remove(str(whitelist_item))

    def set_organization_exclude_patterns(self, exclude_patterns):
        self.model.organization_exclude_patterns = [str(ep) for ep in exclude_patterns]

    def add_organization_exclude_pattern(self, exclude_pattern):
        if self.model.organization_exclude_patterns is None:
            self.model.organization_exclude_patterns = [str(exclude_pattern)]
        else:
            self.model.organization_exclude_patterns.append(str(exclude_pattern))

    def remove_organization_exclude_pattern(self, exclude_pattern):
        if str(exclude_pattern) in self.model.organization_exclude_patterns:
            self.model.organization_exclude_patterns.remove(str(exclude_pattern))

    def get_organization_agreements(self, # pylint: disable=arguments-differ
                                    agreement_signed=None,
                                    agreement_approved=None):
        return Agreement().get_agreements_by_reference(self.get_organization_id(), 'organization',
                                                       agreement_approved=agreement_approved,
                                                       agreement_signed=agreement_signed)

    def all(self, ids=None):
        if ids is None:
            organizations = self.model.scan()
        else:
            organizations = OrganizationModel.batch_get(ids)
        ret = []
        for organization in organizations:
            org = Organization()
            org.model = organization
            ret.append(org)
        return ret

class StoreModel(Model):
    """
    Represents a key-value store in a DynamoDB.
    """
    class Meta:
        """Meta class for Store."""
        table_name = 'cla_store'
        host = cla.conf['DATABASE_HOST']
        region = cla.conf['DYNAMO_REGION']
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
    key = UnicodeAttribute(hash_key=True)
    value = JSONAttribute()

class Store(key_value_store_interface.KeyValueStore):
    """
    ORM-agnostic wrapper for the DynamoDB key-value store model.
    """
    def __init__(self):
        super(Store).__init__()

    def set(self, key, value):
        model = StoreModel()
        model.key = key
        model.value = value
        model.save()

    def get(self, key):
        model = StoreModel()
        try:
            return model.get(key).value
        except StoreModel.DoesNotExist:
            raise cla.models.DoesNotExist('Key not found')

    def delete(self, key):
        model = StoreModel()
        model.key = key
        model.delete()

    def exists(self, key):
        # May want to find a better way. Maybe using model.count()?
        try:
            self.get(key)
            return True
        except cla.models.DoesNotExist:
            return False
