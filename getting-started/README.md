# Getting Started

## What is EasyCLA?

_EasyCLA_ helps maintainers of open source projects streamline their workflows and reduce the hassle of managing Contributor License Agreements \(CLAs\) and authorizing contributors. By automating many of the manual processes, this open source solution hosted by the Linux Foundation reduces delays for developers to get authorized under a CLA.

## What is a CLA?

A _Contributor License Agreement_ \(CLA\) defines the terms under which intellectual property is contributed to a company or project. Typically, the intellectual property is software under an open source license. EasyCLA helps to ensure that contributions are not pulled into a project unless a CLA covering the contributor has been signed. CLAs typically fall into one of two categories:

* **Corporate Contributor License Agreement**

  If the company \(employer\) owns the contribution, a CCLA signatory signs a Corporate CLA. The Corporate CLA legally binds the corporation, so the agreement must be signed by a person with authority to enter into legal contracts on behalf of the corporation. A Corporate CLA may not remove the need for every employee \(developer\) to sign their own Individual CLA -- which separately covers contributions owned by the individual contributor -- if the project requires this.

* **Individual Contributor License Agreement**

  If as an individual you own the contribution, you sign the Individual CLA. A signed Individual CLA may be required before an individual's contribution can be merged into the project repository.

## How Does it Work?

This high-level diagram shows the different flows and roles that EasyCLA supports:

![CLA Diagram](../.gitbook/assets/cla_diagram_v8.png)

## EasyCLA Architechture

The following diagram explains the EasyCLA architecture. 

![](../.gitbook/assets/architecture-overview.png)

## What Role are You?

How you interact with EasyCLA depends on your role. EasyCLA supports the following roles in its workflow:

### Project Manager

You are a _project manager_ if you are the project maintainer who has responsibilities such as managing a project’s GitHub organization or Gerrit instance, members, repositories, and CLAs.

With EasyCLA, you do the following CLA set-up tasks:

1. [Install the EasyCLA Application](install-the-easycla-application.md).
2. [Add a CLA Group](add-a-cla-group.md).
3. [Add Contributor License Agreements](add-contributor-license-agreements.md).
4. Add [GitHub repositories](add-github-repositories-to-cla-monitoring-or-remove-them-from-cla-monitoring).md or [Gerrit instances](add-gerrit-instances-to-cla-monitoring-or-delete-them-from-cla-monitoring.md) to enforce CLA monitoring.

At any time, you can change the settings to manage your project CLA monitoring, and do other management tasks:

* [View Current and Previous CLA PDFs](https://lf-docs-linux-foundation.gitbook.io/easycla/getting-started/project-managers/view-current-and-previous-cla-pdfs)
* [Manage CLA Group Details](https://lf-docs-linux-foundation.gitbook.io/easycla/getting-started/project-managers/manage-cla-group-details)

### Contributor

You are a _contributor_ \(developer\) if you contribute code to GitHub or Gerrit projects. With EasyCLA, you will follow different workflows depending on whether the project is hosted on GitHub or Gerrit, and whether you contribute on behalf of a company or yourself as an individual:

* **GitHub company** contributor: [confirm your association with a company ](https://lf-docs-linux-foundation.gitbook.io/easycla/getting-started/contributors/contribute-to-a-github-company-project)that has a signed Corporate Contributor License Agreement.
* **GitHub individual** contributor: [sign an Individual Contributor License Agreement](https://lf-docs-linux-foundation.gitbook.io/easycla/getting-started/contributors/sign-a-cla-as-an-individual-contributor-to-github).
* **Gerrit company** contributor: [confirm your association with a company ](https://lf-docs-linux-foundation.gitbook.io/easycla/getting-started/contributors/contribute-to-a-gerrit-project)that has a signed Corporate Contributor License Agreement.
* **Gerrit individual** contributor: [sign an Individual Contributor License Agreement](https://lf-docs-linux-foundation.gitbook.io/easycla/getting-started/contributors/contribute-to-a-gerrit-project).

### Corporate CLA Manager

You are a _Corporate CLA manager_ \(CCLA manager\) if you are the person authorized to manage the list of approved contributors under your company’s Corporate CLA. With this responsibility, you use EasyCLA to:

* [Add companies to a project](https://lf-docs-linux-foundation.gitbook.io/easycla/getting-started/ccla-managers-and-ccla-signatories/add-a-company-to-a-project).
* [Whitelist contributors](https://lf-docs-linux-foundation.gitbook.io/easycla/getting-started/ccla-managers-and-ccla-signatories/whitelist-contributors).

### Corporate CLA Signatory

You are a _Corporate CLA signatory_ \(CCLA signatory\) if you are authorized to sign contracts, such as the project’s CLA, on behalf of the company. With EasyCLA, you can:

* [Sign a Corporate CLA on behalf of the company](https://lf-docs-linux-foundation.gitbook.io/easycla/getting-started/ccla-managers-and-ccla-signatories/sign-a-corporate-cla-on-behalf-of-the-company)—as a signatory you need to have legal authority to sign documents on behalf of the company.
* [Review and sign a Corporate CLA by request](https://lf-docs-linux-foundation.gitbook.io/easycla/getting-started/ccla-managers-and-ccla-signatories/review-and-sign-a-corporate-cla-by-request).

