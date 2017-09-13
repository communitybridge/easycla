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
    tables = [RepositoryModel, ProjectModel, SignatureModel, \
              CompanyModel, UserModel, StoreModel, GitHubOrgModel]
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
    tables = [RepositoryModel, ProjectModel, SignatureModel, \
              CompanyModel, UserModel, StoreModel, GitHubOrgModel]
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

class ExternalProjectIndex(GlobalSecondaryIndex):
    """
    This class represents a global secondary index for querying projects by external ID.
    """
    class Meta:
        """Meta class for external ID project index."""
        index_name = 'external-project-index'
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
        # All attributes are projected - not sure if this is necessary.
        projection = AllProjection()

    # This attribute is the hash key for the index.
    project_external_id = UnicodeAttribute(hash_key=True)

class ProjectSignatureIndex(GlobalSecondaryIndex):
    """
    This class represents a global secondary index for querying signatures by project ID.
    """
    class Meta:
        """Meta class for reference Signature index."""
        index_name = 'project-signature-index'
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
        # All attributes are projected - not sure if this is necessary.
        projection = AllProjection()

    # This attribute is the hash key for the index.
    signature_project_id = UnicodeAttribute(hash_key=True)


class ReferenceSignatureIndex(GlobalSecondaryIndex):
    """
    This class represents a global secondary index for querying signatures by reference.
    """
    class Meta:
        """Meta class for reference Signature index."""
        index_name = 'reference-signature-index'
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
        # All attributes are projected - not sure if this is necessary.
        projection = AllProjection()

    # This attribute is the hash key for the index.
    signature_reference_id = UnicodeAttribute(hash_key=True)


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
    document_major_version = NumberAttribute(default=1)
    document_minor_version = NumberAttribute(default=0)
    document_author_name = UnicodeAttribute()


class Document(model_interfaces.Document):
    """
    ORM-agnostic wrapper for the DynamoDB Document model.
    """
    def __init__(self, # pylint: disable=too-many-arguments
                 document_name=None,
                 document_file_id=None,
                 document_content_type=None,
                 document_content=None,
                 document_major_version=None,
                 document_minor_version=None,
                 document_author_name=None):
        super().__init__()
        self.model = DocumentModel()
        self.model.document_name = document_name
        self.model.document_file_id = document_file_id
        self.model.document_author_name = document_author_name
        self.model.document_content_type = document_content_type
        self.model.document_content = self.set_document_content(document_content)
        if document_major_version is not None:
            self.model.document_major_version = document_major_version
        if document_minor_version is not None:
            self.model.document_minor_version = document_minor_version

    def to_dict(self):
        return {'document_name': self.model.document_name,
                'document_file_id': self.model.document_file_id,
                'document_content_type': self.model.document_content_type,
                'document_content': self.model.document_content,
                'document_author_name': self.model.document_author_name,
                'document_major_version': self.model.document_major_version,
                'document_minor_version': self.model.document_minor_version}

    def get_document_name(self):
        return self.model.document_name

    def get_document_file_id(self):
        return self.model.document_file_id

    def get_document_content_type(self):
        return self.model.document_content_type

    def get_document_author_name(self):
        return self.model.document_author_name

    def get_document_content(self):
        content_type = self.get_document_content_type()
        if content_type is None:
            cla.log.warning('Empty content type for document - not sure how to retrieve content')
        else:
            if content_type.startswith('storage+'):
                filename = self.get_document_file_id()
                return cla.utils.get_storage_service().retrieve(filename)
        return self.model.document_content

    def get_document_major_version(self):
        return self.model.document_major_version

    def get_document_minor_version(self):
        return self.model.document_minor_version

    def set_document_author_name(self, document_author_name):
        self.model.document_author_name = document_author_name

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

    def set_document_major_version(self, version):
        self.model.document_major_version = version

    def set_document_minor_version(self, version):
        self.model.document_minor_version = version


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
    project_external_id = UnicodeAttribute()
    project_name = UnicodeAttribute()
    project_individual_documents = ListAttribute(of=DocumentModel, default=[])
    project_corporate_documents = ListAttribute(of=DocumentModel, default=[])
    project_ccla_requires_icla_signature = BooleanAttribute()
    project_external_id_index = ExternalProjectIndex()

class Project(model_interfaces.Project): # pylint: disable=too-many-public-methods
    """
    ORM-agnostic wrapper for the DynamoDB Project model.
    """
    def __init__(self, project_id=None, project_external_id=None, project_name=None,
                 project_ccla_requires_icla_signature=None):
        super(Project).__init__()
        self.model = ProjectModel()
        self.model.project_id = project_id
        self.model.project_external_id = project_external_id
        self.model.project_name = project_name
        self.model.project_ccla_requires_icla_signature = project_ccla_requires_icla_signature

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

    def get_project_external_id(self):
        return self.model.project_external_id

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

    def get_project_individual_document(self, major_version=None, minor_version=None):
        document_models = self.get_project_individual_documents()
        num_documents = len(document_models)
        if num_documents < 1:
            raise cla.models.DoesNotExist('No individual document exists for this project')
        if major_version is None:
            major_version, minor_version = cla.utils.get_last_version(document_models)
        # TODO Need to optimize this on the DB side.
        for document in document_models:
            if document.get_document_major_version() == major_version and \
               document.get_document_minor_version() == minor_version:
                return document
        raise cla.models.DoesNotExist('Document revision not found')

    def get_project_corporate_document(self, major_version=None, minor_version=None):
        document_models = self.get_project_corporate_documents()
        num_documents = len(document_models)
        if num_documents < 1:
            raise cla.models.DoesNotExist('No corporate document exists for this project')
        if major_version is None:
            major_version, minor_version = cla.utils.get_last_version(document_models)
        # TODO Need to optimize this on the DB side.
        for document in document_models:
            if document.get_document_major_version() == major_version and \
               document.get_document_minor_version() == minor_version:
                return document
        raise cla.models.DoesNotExist('Document revision not found')

    def get_project_ccla_requires_icla_signature(self):
        return self.model.project_ccla_requires_icla_signature

    def get_project_latest_major_version(self):
        pass
        # @todo: Loop through documents for this project, return the highest version of them all.

    def set_project_id(self, project_id):
        self.model.project_id = str(project_id)

    def set_project_external_id(self, project_external_id):
        self.model.project_external_id = str(project_external_id)

    def set_project_name(self, project_name):
        self.model.project_name = project_name

    def add_project_individual_document(self, document):
        self.model.project_individual_documents.append(document.model)

    def add_project_corporate_document(self, document):
        self.model.project_corporate_documents.append(document.model)

    def remove_project_individual_document(self, document):
        new_documents = _remove_project_document(self.model.project_individual_documents,
                                                 document.get_document_major_version(),
                                                 document.get_document_minor_version())
        self.model.project_individual_documents = new_documents

    def remove_project_corporate_document(self, document):
        new_documents = _remove_project_document(self.model.project_corporate_documents,
                                                 document.get_document_major_version(),
                                                 document.get_document_minor_version())
        self.model.project_corporate_documents = new_documents

    def set_project_individual_documents(self, documents):
        self.model.project_individual_documents = documents

    def set_project_corporate_documents(self, documents):
        self.model.project_corporate_documents = documents

    def set_project_ccla_requires_icla_signature(self, ccla_requires_icla_signature):
        self.model.project_ccla_requires_icla_signature = ccla_requires_icla_signature

    def get_project_repositories(self):
        repository_generator = RepositoryModel.repository_project_index.query(self.get_project_id())
        repositories = []
        for repository_model in repository_generator:
            repository = Repository()
            repository.model = repository_model
            repositories.append(repository)
        return repositories

    def get_project_signatures(self, signature_signed=None, signature_approved=None):
        return Signature().get_signatures_by_project(self.get_project_id(),
                                                     signature_approved=signature_approved,
                                                     signature_signed=signature_signed)

    def get_project_by_external_id(self, project_external_id):
        """Currently only returns the first one found."""
        project_generator = self.model.project_external_id_index.query(project_external_id)
        for project_model in project_generator:
            project = Project()
            project.model = project_model
            return project
        return None

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

def _remove_project_document(documents, major_version, minor_version):
    # TODO Need to optimize this on the DB side - delete directly from list of records.
    new_documents = []
    found = False
    for document in documents:
        if document.document_major_version == major_version and \
           document.document_minor_version == minor_version:
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
    user_external_id = UnicodeAttribute(null=True)
    user_email = UnicodeAttribute()
    user_name = UnicodeAttribute(null=True)
    user_company_id = UnicodeAttribute(null=True)
    user_github_id = NumberAttribute(null=True)
    user_ldap_id = UnicodeAttribute(null=True)
    user_email_index = EmailUserIndex()
    user_github_id_index = GitHubUserIndex()


class User(model_interfaces.User): # pylint: disable=too-many-public-methods
    """
    ORM-agnostic wrapper for the DynamoDB User model.
    """
    def __init__(self, user_external_id=None, user_email=None, user_github_id=None, user_ldap_id=None):
        super(User).__init__()
        self.model = UserModel()
        self.model.user_email = user_email
        self.model.user_external_id = user_external_id
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

    def get_user_external_id(self):
        return self.model.user_id

    def get_user_email(self):
        return self.model.user_email

    def get_user_name(self):
        return self.model.user_name

    def get_user_company_id(self):
        return self.model.user_company_id

    def get_user_github_id(self):
        return self.model.user_github_id

    def get_user_ldap_id(self):
        return self.model.user_ldap_id

    def set_user_id(self, user_id):
        self.model.user_id = user_id

    def set_user_external_id(self, user_external_id):
        self.model.user_external_id = user_external_id

    def set_user_email(self, user_email):
        self.model.user_email = user_email

    def set_user_name(self, user_name):
        self.model.user_name = user_name

    def set_user_company_id(self, company_id):
        self.model.user_company_id = company_id

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

    def get_user_signatures(self, project_id=None, signature_signed=None, signature_approved=None):
        return Signature().get_signatures_by_reference(self.get_user_id(), 'user',
                                                       project_id=project_id,
                                                       signature_approved=signature_approved,
                                                       signature_signed=signature_signed)

    def get_users_by_company(self, company_id):
        user_generator = self.model.scan(user_company_id__eq=str(company_id))
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


class SignatureModel(BaseModel): # pylint: disable=too-many-instance-attributes
    """
    Represents an signature in the database.
    """
    class Meta:
        """Meta class for Signature."""
        table_name = 'cla_signatures'
        host = cla.conf['DATABASE_HOST']
        region = cla.conf['DYNAMO_REGION']
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
    signature_id = UnicodeAttribute(hash_key=True)
    signature_external_id = UnicodeAttribute(null=True)
    signature_project_id = UnicodeAttribute()
    signature_document_minor_version = NumberAttribute()
    signature_document_major_version = NumberAttribute()
    signature_reference_id = UnicodeAttribute()
    signature_reference_type = UnicodeAttribute()
    signature_type = UnicodeAttribute(default='icla') # Only icla/ccla.
    signature_signed = BooleanAttribute(default=False)
    signature_approved = BooleanAttribute(default=False)
    signature_sign_url = UnicodeAttribute(null=True)
    signature_return_url = UnicodeAttribute(null=True)
    signature_callback_url = UnicodeAttribute(null=True)
    signature_project_index = ProjectSignatureIndex()
    signature_reference_index = ReferenceSignatureIndex()


class Signature(model_interfaces.Signature): # pylint: disable=too-many-public-methods
    """
    ORM-agnostic wrapper for the DynamoDB Signature model.
    """
    def __init__(self, # pylint: disable=too-many-arguments
                 signature_id=None,
                 signature_external_id=None,
                 signature_project_id=None,
                 signature_document_minor_version=None,
                 signature_document_major_version=None,
                 signature_reference_id=None,
                 signature_reference_type='user',
                 signature_type=None,
                 signature_signed=False,
                 signature_approved=False,
                 signature_sign_url=None,
                 signature_return_url=None,
                 signature_callback_url=None):
        super(Signature).__init__()
        self.model = SignatureModel()
        self.model.signature_id = signature_id
        self.model.signature_external_id = signature_external_id
        self.model.signature_project_id = signature_project_id
        self.model.signature_document_minor_version = signature_document_minor_version
        self.model.signature_document_major_version = signature_document_major_version
        self.model.signature_reference_id = signature_reference_id
        self.model.signature_reference_type = signature_reference_type
        self.model.signature_type = signature_type
        self.model.signature_signed = signature_signed
        self.model.signature_approved = signature_approved
        self.model.signature_sign_url = signature_sign_url
        self.model.signature_return_url = signature_return_url
        self.model.signature_callback_url = signature_callback_url

    def to_dict(self):
        return dict(self.model)

    def save(self):
        self.model.save()

    def load(self, signature_id):
        try:
            signature = self.model.get(signature_id)
        except SignatureModel.DoesNotExist:
            raise cla.models.DoesNotExist('Signature not found')
        self.model = signature

    def delete(self):
        self.model.delete()

    def get_signature_id(self):
        return self.model.signature_id

    def get_signature_external_id(self):
        return self.model.signature_external_id

    def get_signature_project_id(self):
        return self.model.signature_project_id

    def get_signature_document_minor_version(self):
        return self.model.signature_document_minor_version

    def get_signature_document_major_version(self):
        return self.model.signature_document_major_version

    def get_signature_type(self):
        return self.model.signature_type

    def get_signature_signed(self):
        return self.model.signature_signed

    def get_signature_approved(self):
        return self.model.signature_approved

    def get_signature_sign_url(self):
        return self.model.signature_sign_url

    def get_signature_return_url(self):
        return self.model.signature_return_url

    def get_signature_callback_url(self):
        return self.model.signature_callback_url

    def get_signature_reference_id(self):
        return self.model.signature_reference_id

    def get_signature_reference_type(self):
        return self.model.signature_reference_type

    def set_signature_id(self, signature_id):
        self.model.signature_id = str(signature_id)

    def set_signature_external_id(self, signature_external_id):
        self.model.signature_external_id = str(signature_external_id)

    def set_signature_project_id(self, project_id):
        self.model.signature_project_id = str(project_id)

    def set_signature_document_minor_version(self, document_minor_version):
        self.model.signature_document_minor_version = int(document_minor_version)

    def set_signature_document_major_version(self, document_major_version):
        self.model.signature_document_major_version = int(document_major_version)

    def set_signature_type(self, signature_type):
        self.model.signature_type = signature_type

    def set_signature_signed(self, signed):
        self.model.signature_signed = bool(signed)

    def set_signature_approved(self, approved):
        self.model.signature_approved = bool(approved)

    def set_signature_sign_url(self, sign_url):
        self.model.signature_sign_url = sign_url

    def set_signature_return_url(self, return_url):
        self.model.signature_return_url = return_url

    def set_signature_callback_url(self, callback_url):
        self.model.signature_callback_url = callback_url

    def set_signature_reference_id(self, reference_id):
        self.model.signature_reference_id = reference_id

    def set_signature_reference_type(self, reference_type):
        self.model.signature_reference_type = reference_type

    def get_signatures_by_reference(self, # pylint: disable=too-many-arguments
                                    reference_id,
                                    reference_type,
                                    project_id=None,
                                    signature_signed=None,
                                    signature_approved=None):
        # TODO: Optimize this query to use filters properly.
        signature_generator = self.model.signature_reference_index.query(reference_id)
        signatures = []
        for signature_model in signature_generator:
            if signature_model.signature_reference_type != reference_type:
                continue
            if project_id is not None and \
               signature_model.signature_project_id != project_id:
                continue
            if signature_signed is not None and \
               signature_model.signature_signed != signature_signed:
                continue
            if signature_approved is not None and \
               signature_model.signature_approved != signature_approved:
                continue
            signature = Signature()
            signature.model = signature_model
            signatures.append(signature)
        return signatures

    def get_signatures_by_project(self, project_id, signature_signed=None,
                                  signature_approved=None):
        signature_generator = self.model.signature_project_index.query(project_id)
        signatures = []
        for signature_model in signature_generator:
            if signature_signed is not None and \
               signature_model.signature_signed != signature_signed:
                continue
            if signature_approved is not None and \
               signature_model.signature_approved != signature_approved:
                continue
            signature = Signature()
            signature.model = signature_model
            signatures.append(signature)
        return signatures

    def all(self, ids=None):
        if ids is None:
            signatures = self.model.scan()
        else:
            signatures = SignatureModel.batch_get(ids)
        ret = []
        for signature in signatures:
            agr = Signature()
            agr.model = signature
            ret.append(agr)
        return ret


class CompanyModel(BaseModel):
    """
    Represents an company in the database.
    """
    class Meta:
        """Meta class for Company."""
        table_name = 'cla_companies'
        host = cla.conf['DATABASE_HOST']
        region = cla.conf['DYNAMO_REGION']
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
    company_id = UnicodeAttribute(hash_key=True)
    company_external_id = UnicodeAttribute(null=True)
    company_name = UnicodeAttribute()
    company_whitelist = ListAttribute()
    company_whitelist_patterns = ListAttribute()


class Company(model_interfaces.Company): # pylint: disable=too-many-public-methods
    """
    ORM-agnostic wrapper for the DynamoDB Company model.
    """
    def __init__(self, # pylint: disable=too-many-arguments
                 company_id=None,
                 company_external_id=None,
                 company_name=None,
                 company_whitelist_patterns=None,
                 company_whitelist=None):
        super(Company).__init__()
        self.model = CompanyModel()
        self.model.company_id = company_id
        self.model.company_external_id = company_external_id
        self.model.company_name = company_name
        self.model.company_whitelist = company_whitelist
        self.model.company_whitelist_patterns = company_whitelist_patterns

    def to_dict(self):
        return dict(self.model)

    def save(self):
        self.model.save()

    def load(self, company_id):
        try:
            company = self.model.get(company_id)
        except CompanyModel.DoesNotExist:
            raise cla.models.DoesNotExist('Company not found')
        self.model = company

    def delete(self):
        self.model.delete()

    def get_company_id(self):
        return self.model.company_id

    def get_company_external_id(self):
        return self.model.company_external_id

    def get_company_name(self):
        return self.model.company_name

    def get_company_whitelist(self):
        return self.model.company_whitelist

    def get_company_whitelist_patterns(self):
        return self.model.company_whitelist_patterns

    def set_company_id(self, company_id):
        self.model.company_id = company_id

    def set_company_external_id(self, company_external_id):
        self.model.company_external_id = company_external_id

    def set_company_name(self, company_name):
        self.model.company_name = str(company_name)

    def set_company_whitelist(self, whitelist):
        self.model.company_whitelist = [str(wl) for wl in whitelist]

    def add_company_whitelist(self, whitelist_item):
        if self.model.company_whitelist is None:
            self.model.company_whitelist = [str(whitelist_item)]
        else:
            self.model.company_whitelist.append(str(whitelist_item))

    def remove_company_whitelist(self, whitelist_item):
        if str(whitelist_item) in self.model.company_whitelist:
            self.model.company_whitelist.remove(str(whitelist_item))

    def set_company_whitelist_patterns(self, whitelist_patterns):
        self.model.company_whitelist_patterns = [str(wp) for wp in whitelist_patterns]

    def add_company_whitelist_pattern(self, whitelist_pattern):
        if self.model.company_whitelist_patterns is None:
            self.model.company_whitelist_patterns = [str(whitelist_pattern)]
        else:
            self.model.company_whitelist_patterns.append(str(whitelist_pattern))

    def remove_company_whitelist_pattern(self, whitelist_pattern):
        if str(whitelist_pattern) in self.model.company_whitelist_patterns:
            self.model.company_whitelist_patterns.remove(str(whitelist_pattern))

    def get_company_signatures(self, # pylint: disable=arguments-differ
                               signature_signed=None,
                               signature_approved=None):
        return Signature().get_signatures_by_reference(self.get_company_id(), 'company',
                                                       signature_approved=signature_approved,
                                                       signature_signed=signature_signed)

    def all(self, ids=None):
        if ids is None:
            companies = self.model.scan()
        else:
            companies = CompanyModel.batch_get(ids)
        ret = []
        for company in companies:
            org = Company()
            org.model = company
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


class GitHubOrgModel(BaseModel):
    """
    Represents a user in the database.
    """
    class Meta:
        """Meta class for User."""
        table_name = 'cla_github_orgs'
        host = cla.conf['DATABASE_HOST']
        region = cla.conf['DYNAMO_REGION']
        write_capacity_units = cla.conf['DYNAMO_WRITE_UNITS']
        read_capacity_units = cla.conf['DYNAMO_READ_UNITS']
    organization_name = UnicodeAttribute(hash_key=True)
    organization_company_id = UnicodeAttribute(null=True)
    organization_installation_id = NumberAttribute(null=True)
    organization_project_id = UnicodeAttribute(null=True)


class GitHubOrg(model_interfaces.GitHubOrg): # pylint: disable=too-many-public-methods
    """
    ORM-agnostic wrapper for the DynamoDB GitHubOrg model.
    """
    def __init__(self, organization_name=None, organization_company_id=None, organization_installation_id=None, organization_project_id=None):
        super(User).__init__()
        self.model = GitHubOrgModel()
        self.model.organization_name = organization_name
        self.model.organization_company_id = organization_company_id
        self.model.organization_installation_id = organization_installation_id
        self.model.organization_project_id = organization_project_id

    def to_dict(self):
        ret = dict(self.model)
        if ret['organization_installation_id'] == 'null':
            ret['organization_installation_id'] = None
        if ret['organization_project_id'] == 'null':
            ret['organization_project_id'] = None
        return ret

    def save(self):
        self.model.save()

    def load(self, organization_name):
        try:
            organization = self.model.get(str(organization_name))
        except GitHubOrgModel.DoesNotExist:
            raise cla.models.DoesNotExist('GitHub Org not found')
        self.model = organization

    def delete(self):
        self.model.delete()

    def get_organization_name(self):
        return self.model.organization_name

    def get_organization_company_id(self):
        return self.model.organization_company_id

    def get_organization_installation_id(self):
        return self.model.organization_installation_id

    def get_organization_project_id(self):
        return self.model.organization_project_id

    def set_organization_name(self, organization_name):
        self.model.organization_name = organization_name

    def set_organization_company_id(self, organization_company_id):
        self.model.organization_company_id = organization_company_id

    def set_organization_installation_id(self, organization_installation_id):
        self.model.organization_installation_id = organization_installation_id

    def set_organization_project_id(self, organization_project_id):
        self.model.organization_project_id = organization_project_id

    def get_organizations_by_company(self, company_id):
        organization_generator = self.model.scan(organization_company_id__eq=str(company_id))
        organizations = []
        for org_model in organization_generator:
            org = GitHubOrg()
            org.model = org_model
            organizations.append(org)
        return organizations

    def get_organization_by_installation_id(self, installation_id):
        organization_generator = self.model.scan(organization_installation_id__eq=installation_id)
        for org_model in organization_generator:
            org = GitHubOrg()
            org.model = org_model
            return org
        return None

    def all(self):
        orgs = self.model.scan()
        ret = []
        for organization in orgs:
            org = GitHubOrg()
            org.model = organization
            ret.append(org)
        return ret
