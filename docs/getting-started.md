# What is EasyCLA?
EasyCLA helps maintainers of open source projects streamline their workflows and reduce the hassle of managing Contributor License Agreements (CLAs) and authorizing contributors. By automating many of the manual processes, this open source solution hosted by the Linux Foundation reduces delays for developers to get authorized under a CLA.

## What is a CLA?
A Contributor License Agreement (CLA) defines the terms under which intellectual property (IP) is contributed to a company or project. Typically, the intellectual property is software under an open source license. EasyCLA guards a project's outputs so that the necessary ownership or grant of rights over all contributions is distributed under the chosen license. A contract defines the legal status of the contributed code in two types of CLAs:

### Corporate Contributor License Agreement

  If the company (employer) owns the contribution, a CLA signatory signs a Corporate CLA. The Corporate CLA legally binds the corporation, so the agreement must be signed by a person with authority to enter into legal contracts on behalf of the corporation. A Corporate CLA may not remove the need for every employee (developer) to sign their own Individual CLA, which covers both contributions which are owned and those that are not owned by the corporation signing the Corporate CLA.


### Individual Contributor License Agreement

  If as an individual you own the contribution, you sign the Individual CLA. A signed Individual CLA may be required before an individual is given commit rights to a CLA-defined project. 

## How Does it Work?
This high-level diagram shows the different flows and roles that EasyCLA supports:
[EasyCLA Workflow](https://docs.linuxfoundation.org/display/DOCS/CommunityBridge+EasyCLA?preview=/4822539/7413213/CLA%20EasyCLA%20workflow.png)

## What Role are You?
How you interact with EasyCLA depends on your role. EasyCLA supports the following roles in its workflow:


### [Project Manager](https://docs.linuxfoundation.org/display/DOCS/Project+Managers)
You are a Project Manager if you are the project maintainer who is responsible for selecting the appropriate Individual and Corporate CLA.

Traditionally, the Project Manager has had to enforce whether a contributor was authorized to commit code to their project at every commit. This became especially cumbersome when getting signatures for Corporate CLAs from companies and updating each company’s whitelist of authorized developers. With EasyCLA, Project Managers handle CLA setup and management tasks whereas CLA Managers handle whitelists and company details.

### [Contributor](https://docs.linuxfoundation.org/display/DOCS/Contributors)
You are a contributor (developer) to GitHub or Gerrit projects.

With EasyCLA, you easily comply to your legal obligations for a company or contribute individually by confirming your association with a company that has a signed Corporate Contributor License Agreement or by signing an Individual Contributor License Agreement.

### [CLA Manager](https://docs.linuxfoundation.org/display/DOCS/CLA+Managers+and+CLA+Signatories)
You are the CLA Manager if you are the person authorized to manage who can contribute under your company’s Corporate CLA. With this responsibility, you use EasyCLA to add companies to a project and whitelist contributors.

### [CLA Signatory](https://docs.linuxfoundation.org/display/DOCS/CLA+Managers+and+CLA+Signatories)
You are the CLA Signatory if you are the authorized signatory of the project’s CLA for the company. Typically this is someone within the counsel’s office of the company. Within EasyCLA, you respond to emails asking you to sign the CLA, and sign a Corporate CLA on behalf of the company—as a signatory you have legal authority to sign documents on behalf of the company.

