

<p align="center">
  <a href="https://lfcla.com/">
    <img src="https://project.dev.lfcla.com/assets/logo/cp_app_easycla.svg" alt="CommunityBridge EasyCLA logo" width=25%>
  </a>
</p>



# Testplan
This document list all the features of EasyCLA consoles and status of workingness.
<!-- https://lfcla.com/static/logo.svg -->


## Table of Contents

- [Testplan](#testplan)
  * [Table of Contents](#table-of-contents)
  * [Test Requirements](#test-requirements)
  * [Project Management console](#project-management-console)
    + [Project Manager Creates/Edits a Contract Group](#project-manager-createsedits-a-contract-group)
  * [Corporate console](#corporate-console)
    + [CLA Manager Signs CCLA](#cla-manager-signs-ccla)
  * [Contributor console](#contributor-console)
    + [CCLA Contributor - company exists, contributor email is in whitelist](#ccla-contributor---company-exists-contributor-email-is-in-whitelist)
    + [CCLA Contributor - company exists, contributor email is not in whitelist](#ccla-contributor---company-exists-contributor-email-is-not-in-whitelist)
    + [CCLA Contributor - company does not exist](#ccla-contributor---company-does-not-exist)
    + [ICLA Contributor](#icla-contributor)
    + [Contributor Has Already Signed CLA](#contributor-has-already-signed-cla)
    + [Gerrit CLA Flow](#gerrit-cla-flow)


## Test Requirements

- App URLS
  - Landing page: https://test.lfcla.com
  - Project Management Console : https://project.test.lfcla.com
  - Corporate Console : https://corporate.test.lfcla.com
- Few LFIDs to play different roles.



## Project Management console

### Project Manager Creates/Edits a Contract Group

| Test name                   | Description                                                                                                                                                                                             | Status  |
|-----------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:-------:|
| Project manager login       | Users(Project Manager) with Valid LFID and required permissions should be able to login.                                                                                                                | ✅      |
| Project selection           | After successful login Project Manager should see list of projects configured for them and able to select a project of interest.                                                                        | ✅      |
| Project detail page         | Open project detail page which shows options to configure CLA group, CLA Templates and Github Repositories etc                                                                                          | ✅      |
| Create CLA group            | Create CLA group by clicking **ADD CLA GROUP** button. <br> On form submission project detail page list CLA group with selected CLA types.                                                              | ✅      |
| Connect Github Organization | Project manager should be able connect their Github organization to CLA project by installing GH app in their GH org. <br> Once installed **GitHub Organizations** section should list the organization.| ✅      |
| Configure CCLA/ICLA         | Configured templates should appear on CLA group.                                                                                                                                                        | ✅      |
| Configure Repositories      | Clicking **Configure Github Repositories** should list repos from org and appear on CLA group when added.                                                                                               | ✅      |
| Connect Gerrit Instance     | Gerrit instance config file detail should appear on CLA group once Gerrit connection details configured.                                                                                                | ✅      |
| View signatures             | View signatures button should list all the signatures signed for that CLA group.                                                                                                                        | ✅      |



## Corporate console



### CLA Manager Signs CCLA

| Test name                             | Description                                                                                                                 | Status |
|:--------------------------------------|:----------------------------------------------------------------------------------------------------------------------------|:------:|
| Login to Corporate management console | A user with valid LFID should be able to login to Corporate console.                                                        |   ✅    |
| Create/Find Company Company           | CLA manager should be able to search for existing companies and join it or create new company.                              |   ✅    |
| Accept CLA manager join requests      | CLA manager can request to join exiting company. It should get listed in *Pending Invites* section of Company detail page.  |   ✅    |
| Select Contract Group                 | CLA manager should be able to select a Contract group to sign.                                                              |   ✅    |
| Sign CCLA                             | Sign a CCLA for CLA group on behalf of a company                                                                            |   ✅    |
| Domain whitelists                     | White list corporate contributors by domain names                                                                           |   ✅    |
| Email whitelists                      | White list corporate contributors by emails.                                                                                |   ✅    |
| Github username whitelists            | White list corporate contributors by Github usernames                                                                       |   ✅    |
| Github organization whitelists        | Connect Github Orgs to the Corporate console and all the public members of the org should be whitelisted.                   |   ✅    |
| Add CLA manager for a CCLA            | Multiple manager can manage a signed CCLA and they can be added in *CLA Managers* section of CCLA detail page.              |   ✅    |
| Employee Acknowledgments              | List of corporate contributors who confirmed association with company should be list in *Employee Acknowledgments* section  |   ✅    |



## Contributor console

### CCLA Contributor - company exists, contributor email is in whitelist

| Test name                                                           | Description                                                                                                                                      | Status |
|:--------------------------------------------------------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------|:------:|
| Create PR - CLA check fails                                         | Corporate contributor with no previous CLA sign should get CLA check fail on their PR.                                                           |   ✅    |
| CLA Link opens to contributor console.Contributor selects "Company" | CLA sign link from PR should open Contributor console,<br> after Github oauth flow, they should get option to select   *Individual* or *Company* |   ✅    |
| Contributor selects their company from the company list             | Contributor should be able to select the Company from list of companies.                                                                         |   ✅    |
| Contributor acknowledges their employment for the company           | Contributor should be able to confirm the association with company by selecting checkbox.                                                        |   ✅    |
| Return to PR. CLA check succeeds                                    | After confirming association corp contributor should be redirected to PR.                                                                        |   ✅    |





### CCLA Contributor - company exists, contributor email is not in whitelist

| Test name                                                                                | Description                                                                                                                                      | Status |
|:-----------------------------------------------------------------------------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------|:------:|
| Create PR - CLA check fails                                                              | Corporate contributor with no previous CLA sign should get CLA check fail on their PR.                                                           |   ✅    |
| CLA Link opens to contributor console. Contributor selects "Company"                     | CLA sign link from PR should open Contributor console,<br> after Github oauth flow, they should get option to select   *Individual* or *Company* |   ✅    |
| Contributor selects their company from the company list                                  | Contributor should be able to select the Company from list of companies.                                                                         |   ✅    |
| Contributor's email address fails whitelist check. Presented with "troubleshoot" screen. | Contributor is not whitelisted in corp console, they should get informational page to get whitelisted.                                           |   ✅    |
| Contributor selects the option to email their CLA Manager                                | Contact CLA manager form should send email to CLA manager.                                                                                       |   ✅    |


### CCLA Contributor - company does not exist

| Test name                                                           | Description                                                                                                                                      | Status |
|:--------------------------------------------------------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------|:------:|
| Create PR - CLA check fails                                         | Corporate contributor with no previous CLA sign should get CLA check fail on their PR.                                                           |   ✅    |
| CLA Link opens to contributor console.Contributor selects "Company" | CLA sign link from PR should open Contributor console,<br> after Github oauth flow, they should get option to select   *Individual* or *Company* |   ✅    |
| Contributor confirms that they are not a CLA Manager                | Contributors employer is not listed the companies list. Should be able select NO from the confirm modal.                                         |   ✅    |
| Contributor fills out CLA Manager info form                         | Form to inform the CLA manager of the contributor company to create Company in corp console.                                                     |   ✅    |


### CCLA Contributor - company exists but CCLA not signed.

| Test name                                                           | Description                                                                                                                                      | Status |
|:--------------------------------------------------------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------|:------:|
| Create PR - CLA check fails                                         | Corporate contributor with no previous CLA sign should get CLA check fail on their PR.                                                           |   ✅    |
| CLA Link opens to contributor console.Contributor selects "Company" | CLA sign link from PR should open Contributor console,<br> after Github oauth flow, they should get option to select   *Individual* or *Company* |   ✅    |
| Contributor chooses their company                                   | Contributors employer chooses company name from list they should get modal saying "You company has not signed CLA yet"                           |   ✅    |
| Contributor fills out CLA Manager info form                         | Form to send email to the CLA manager of the contributor company to sign a CCLA for the project and whitelist the contributor                           |   ✅    |



### ICLA Contributor

| Test name                                                               | Description                                                                                    | Status |
|:------------------------------------------------------------------------|:-----------------------------------------------------------------------------------------------|:------:|
| Create PR - CLA check fails                                             | Corporate contributor with no previous CLA sign should get CLA check fail on their PR.         |   ✅    |
| CLA Link opens to contributor console. Contributor selects "Individual" | Individuals should be redirected docusign to fill details and sign.                            |   ✅    |
| Return to PR. CLA check succeeds                                        | Once sign completed contributor should get redirected to PR and  check should show green tick. |   ✅    |


### Contributor Has Already Signed CLA

| Test name                            | Description                                                          | Status |
|:-------------------------------------|:---------------------------------------------------------------------|:------:|
| Create PR                            | Send PR to repo which has a valid CLA signature for the contributor. |   ✅    |
| CLA check marks PR as ready to merge | CLA check should pass automatically.                                 |   ✅    |


### Gerrit CLA Flow

| Test name                                                                         | Description                                                                                          | Status |
|:----------------------------------------------------------------------------------|:-----------------------------------------------------------------------------------------------------|:------:|
| Contributor tries to push to a Gerrit Instance                                    | Patches pushed to the Gerrit repo should be rejected with message asking sign CLA agreements.        |   ✅    |
| Contributor opens Gerrit agreements page                                          | Contributor opens agreements section from Gerrit profile settings                                    |   ✅    |
| Contributor views available agreements                                            | Configured Agreements (ICLA and CCLA )should be listed here                                          |   ✅    |
| Contributor selects CLA agreement link - is redirected to CLA Contributor Console | Contributor can choose agreement type and it should generate a link to sign CLA using docusign flow. |   ✅    |
| Retry pushing patch after CLA sing                                                | Contributor should be able to submit patches to Gerrit repo as a CLA is signed.                      |   ✅    |


