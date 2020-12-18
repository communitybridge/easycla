# Introduction to EasyCLA

[![CircleCI](https://circleci.com/gh/communitybridge/easycla.svg?style=svg)](https://circleci.com/gh/communitybridge/easycla)

The Contributor License Agreement \(CLA\) service of the Linux Foundation lets project contributors read, sign, and submit contributor license agreements easily.

This repository contains both the backend and front-end UI for supporting and managing the application.

This platform supports both GitHub and Gerrit source code repositories. Additional information can be found in the [EasyCLA documentation](https://docs.linuxfoundation.org/lfx/easycla).

## Announcements

- 06/10/2020 - We have replaced code references from whitelist to Approved List.  This includes package names, code comments, swagger specifications, API endpoints, variable names, and UI components.

## Third-party Services

[EasyCLA](#easycla-architecture)

Besides integration with Auth0 and Salesforce, the CLA system has the following third party services:

* [Docusign](https://www.docusign.com/) for CLA agreement e-sign flow
* [Docraptor](https://docraptor.com/) for converting html CLA template to PDF file

## CLA Backend

The CLA project has two backend components:

* The majority of the backend APIs are implemented in python, and can be found in the [cla-backend](cla-backend/) directory.

* Recent backend development is implemented in Golang, and can be found in the
[cla-backend-go](cla-backend-go/) directory. In particular, this backend contains APIs powering
Automated Templates, GitHub Approval Lists, and Duplicate Company handling in the
Corporate Console.

## CLA Frontend

CLA frontend consists of three independent SPA built with [Ionic](https://ionicframework.com/) framework.

* [cla-frontend-project-console](cla-frontend-project-console/) for the LinuxFoundation director/admin/user to manage project CLA
* [cla-frontend-corporate-console](cla-frontend-corporate-console/) for any concrete company CCLA manager to sign a CCLA and manage employee CLA approved list
* [cla-frontend-contributor-console](cla-frontend-contributor-console) for any project contributor to sign ICLA or CCLA

## EasyCLA Architecture

The following diagram explains the EasyCLA architecture.

![CLA Architecture](.gitbook/assets/easycla-architecture-overview.png)

## EasyCLA Release Process

The following diagram illustrates the EasyCLA release process:

![CLA Release Process](.gitbook/assets/easycla_software_development_and_release_process.png)

## License

Copyright The Linux Foundation and each contributor to CommunityBridge.

This project’s source code is licensed under the MIT License. A copy of the license is available in LICENSE.

The project includes source code from keycloak, which is licensed under the Apache License, version 2.0 \(Apache-2.0\), a copy of which is available in LICENSE-keycloak.

This project’s documentation is licensed under the Creative Commons Attribution 4.0 International License \(CC-BY-4.0\). A copy of the license is available in LICENSE-docs.

