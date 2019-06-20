# DynamoDB Tables

Copyright The Linux Foundation and each contributor to CommunityBridge.
SPDX-License-Identifier: GFDL-1.3-or-later

## Companies
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| company_id | Primary Key UUID for Companies table |
| company_acl | List of CLA Managers of the Company |
| company_name | Name of the Company |
| company_manager_id | Initial Company Manager who created the company |

## Company-Invites
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| company_invite_id | Primary Key UUID for Company Invites table |
| requested_company_id | Company that the pending invitation will be sent to* |
| user_id | user requesting invitation |

 * When a user tries to create a company that already exists in the Corporate Console, the user can send a request to the CLA Managers to be included in the Company ACL. A row in this table represents a pending request that any CLA Manager of the company can accept or decline in the Corporate Console.

## Gerrit-Instances
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| gerrit_id | Primary Key UUID for Gerrit Instances table |
| gerrit_name | Name of the Gerrit Instance |
| gerrit_url | the URL to the Gerrit Instance |
| group_id_icla | The LDAP Group ID for ICLA for this Gerrit Instance* |
| group_name_icla | The LDAP Group Name for ICLA |
| group_id_ccla | The LDAP ID for CCLA regarding this Gerrit Instance |
| group_name_ccla | The LDAP Group Name for CCLA |
| project_id | The CLA Group ID that the Gerrit Instance is configured with |

 * When you sign an ICLA or a CCLA for this Gerrit Instance, you will be added to this LDAP Group. Once you are part of this group, you will be able to make a contribution to the repositories of the Gerrit Instance.

## Github-Orgs
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| organization_name | Primary Key String Identifier for Github Orgs Table. |
| organization_installation_id | The ID assigned by Github when the CLA App has been installed to the Github Organization |
| organization_sfid | The SFDC ID that the Github Organization is configured to. |

 * Github Organizations can be configured on the Project Management Console. They will be tied to individual SFDC IDs.

## Repositories
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| repository_id | Primary Key UUID for the Repositories table (Github repositories that have been configured under a CLA Group) |
| repository_external_id | The name of the CLA Group |
| repository_name | Name of the repository configured |
| repository_organization_name | Name of the Github org that includes the repository |
| repository_project_id | The CLA Group ID that has the Github Repository configured |
| repository_sfdc_id | The SFDC ID that the repository is configured under |
| repository_type | The repository provider (currently only will have github) |
| repository_url | History of ICLA Templates that the project manager created.* |

## Projects
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| project_id | Primary Key UUID for the Projects table (Referred to as CLA Groups on Project Management Console) |
| project_name | The name of the CLA Group |
| project_acl | The Cla Group Manager ACL (Access Control List) |
| project_ccla_enabled | The boolean flag that determines whether the CLA Group will include a CCLA. |
| project_icla_enabled | The boolean flag that determines whether the CLA Group will include an ICLA. |
| project_ccla_requires_icla_signature | The boolean flag that determines when an employee has signed a CCLA, the employee needs to sign an ICLA as well to pass the CLA checks |
| project_corporate_documents | History of CCLA Templates that the project manager created.* |
| project_corporate_documents | History of ICLA Templates that the project manager created.* |

### The Templates will be saved as a list of Documents, that have the following properties
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| document_author_name | Name of the author of the document (deprecated) |
| document_creation_date | Date and time when the document was created through the Project Management Console |
| document_content_type | Document file type. Currently only supports "storage+pdf" |
| document_file_id | UUID unique to the document generated for management under S3 and Docusign API |
| document_legal_entity_name | The legal entity name for the document (deprecated) |
| document_major_version | The major version of the Document |
| document_minor_version | The minor version of the Document |
| document_name | Name of the CLA Template |
| document_preamble | Preamble of the CLA Template (deprecated) |
| document_s3_url | The URL to the generated CLA Template on the AWS S3 Bucket |
| document_tabs | Document custom tabs that a Signer must fill before completing the singing process, such as Corporate Name, E-mail, Address, etc. * |

### Each Document will have a list of document_tabs, that have the following properties
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| document_tab_anchor_string | The anchor string (Docusign looks for this string in the document to create the custom tab) |
| document_tab_anchor_ignore_if_present | the flag to ignore the custom tab when Docusign is not able to locate the anchor string |
| document_tab_anchor_x_offset | The anchor string X offsets for the custom tab (in pixels) |
| document_tab_anchor_y_offset | The anchor string Y offsets for the custom tab (in pixels) |
| document_tab_id | The custom tab ID identifier |
| document_tab_is_locked | Flag to determine if the custom tab is editable |
| document_tab_is_required | Flag to determine if the custom tab is optional |
| document_tab_name | custom tab name that displays when you hover the mouse on the tab |
| document_tab_type | type of custom tab (e.g. "text", "sign", "date") |
| document_tab_position_x | the X coordinates that the custom tab will be located (deprecated due to anchor string) |
| document_tab_position_y | the Y coordinates that the custom tab will be located (deprecated due to anchor string) |
| document_tab_page | The page number of the custom tab |
| document_tab_width | width of the custom tab |
| document_tab_height | height of the custom tab |

## Signature
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| signature_id | Primary Key UUID for the Signatures table |
| signature_acl | The Signature Access Control List |
| signature_approved | Flag that determines if signature is approved |
| signature_signed | Flag that determines if signature has been signed |
| signature_document_major_version | The major version number for the signed document |
| signature_document_minor_version | The minor version number for the signed document |
| signature_project_id | The CLA Group ID that the signature refers to |
| signature_reference_id | The ID of the owner of the signature (user_id or company_id) |
| signature_reference_type | The type of the owner (user/company) |
| signature_user_ccla_company_id | The company ID that the user is signing the employee CCLA for |
| signature_return_url | The return url that the user will be directed to when the user has completed the signing process (The Github Pull Request Page, the Gerrit Instance page) |
| signature_type | Determines whether the signature is for an ICLA (cla) or a CCLA (ccla) |
| signature_callback_url | When a user completes the signing process, Docusign sends the CLA System a webhook that contains information about the document (i.e. user finished signing, user opened document). The URL that the Docusign calls on the CLA System  so that the system can mark the signature as signed |
| signature_envelope_id | the envelope_id generated by Docusign (used for voiding signatures that are no longer valid) |
| signature_sign_url | The Docusign URL that the user will be directed to fill the custom tabs and sign the document |
| domain_whitelist | The domain whitelist for the CCLA signature. Accepts wildcards on URLs |
| email_whitelist | The email whitelist for the CCLA signature |
| github_whitelist | The Github username whitelist for the CCLA signature |
| github_org_whitelist | The Github organization whitelist for the CCLA signature |

## User-Permissions
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| username | Primary Key Username for the User permissions |
| projects | List of SFDC projects that the user has access to in the Project Management Console |

## Users
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| user_id | Primary Key UUID for the Users Table (users within the CLA system) |
| user_name | The Full name of the user |
| lf_email | The LF Email for the user |
| lf_username | The LFID of the user |
| user_emails | The verified list of emails for a Github User |
| user_github_id | The github ID (i.e. 23942335) for a GIthub user |
| user_github_username | The github username for user |
| user_company_id | The company ID that the user will have when the user completes an employee CCLA |

## Session-Store
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id | Primary Key ID for the Session Stores table (Session stored for Github Oauth) |
| options | The options for the session (Domain, HttpOnly, MaxAge, Path, SameSite, Secure) flags |
| values | The hashed session token for this session |

## Store
| Column Name | Description |
|-----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| key | Primary Key UUID for the Store table (Temporary Key-value storage for user information during the signing process) |
| expire (TTL) | Expiration date time for the key-value pair in epoch |
| value | Temporary values stored for the key-value |