

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

| Test name                             | Description | Status  |
|---------------------------------------|-------------|:-------:|
| Login to Corporate management console |             | ✅      |
| Create Company                        |             | ✅      |
| Select Contract Group                 |             | ✅      |
| Sign CCLA                             |             | ✅      |
| Domain whitelists                     |             | ✅      |
| Email whitelists                      |             | ✅      |
| Github username whitelists            |             | ✅      |
| Github organization whitelists        |             | ✅      |



## Contributor console

### CCLA Contributor - company exists, contributor email is in whitelist

| Test name                                                            | Description | Status  |
|----------------------------------------------------------------------|-------------|:-------:|
| Create PR - CLA check fails                                          |             | ✅      |
| CLA Link opens to contributor console. Contributor selects "Company" |             | ✅      |
| Contributor selects their company from the company list              |             | ✅      |
| Contributor acknowledges their employment for the company            |             | ✅      |
| Return to PR. CLA check succeeds                                     |             | ✅      |





### CCLA Contributor - company exists, contributor email is not in whitelist

| Test name                                                                                | Description | Status  |
|------------------------------------------------------------------------------------------|-------------|:-------:|
| Create PR - CLA check fails                                                              |             | ✅      |
| CLA Link opens to contributor console. Contributor selects "Company"                     |             | ✅      |
| Contributor selects their company from the company list                                  |             | ✅      |
| Contributor's email address fails whitelist check. Presented with "troubleshoot" screen. |             | ✅      |
| Contributor selects the option to email their CLA Manager                                |             | ✅      |


### CCLA Contributor - company does not exist

| Test name                                                            | Description | Status  |
|----------------------------------------------------------------------|-------------|:-------:|
| Create PR - CLA check fails                                          |             | ✅      |
| CLA Link opens to contributor console. Contributor selects "Company" |             | ✅      |
| Contributor confirms that they are not a CLA Manager                 |             | ✅      |
| Contributor fills out CLA Manager info form                          |             | ✅      |

### ICLA Contributor

| Test name                                                               | Description | Status  |
|-------------------------------------------------------------------------|-------------|:-------:|
| Create PR - CLA check fails                                             |             | ✅      |
| CLA Link opens to contributor console. Contributor selects "Individual" |             | ✅      |
| Return to PR. CLA check succeeds                                        |             | ✅      |


### Contributor Has Already Signed CLA

| Test name                            | Description | Status  |
|--------------------------------------|-------------|:-------:|
| Create PR                            |             | ✅      |
| CLA check marks PR as ready to merge |             | ✅      |


### Gerrit CLA Flow

| Test name                                                                         | Description | Status  |
|-----------------------------------------------------------------------------------|-------------|:-------:|
| Contributor tries to push to a Gerrit Instance                                    |             | ✅      |
| Contributor opens Gerrit agreements page                                          |             | ✅      |
| Contributor views available agreements                                            |             | ✅      |
| Contributor selects agreement                                                     |             | ✅      |
| Contributor selects CLA agreement link - is redirected to CLA Contributor Console |             | ✅      |


