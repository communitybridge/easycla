# Contribute to a Gerrit Project

As an individual contributor or a corporate \(employee\) contributor or to an open source project, you submit changed code in Gerrit to inform reviewers about the changes:

* As an individual contributor to an open-source project who is not associated to a company, you submit code in Gerrit and during the process, your CLA is verified. Also during the process, you click a link to open the CLA Contributor Console to sign the CLA Agreement. As an individual contributor, your work is owned by yourself \(Individual CLA\).
* If any of your contributions to the project are created as part of your employment, the work may not belong to you—it may be owned by your employer. In that case, the CCLA signatory represents the employer \(company\) for legal reasons, and signs the Corporate Contributor Licensing Agreement in order for your contribution to be accepted into the company. During the code submission process, your CLA status is verified.

  When all CLA management set-up tasks are complete and your CCLA signatory has signed a Corporate CLA, you simply confirm your association to the company during your first code submission. Then, you can continue. Your subsequent contributions will not require association confirmations.

**Do these steps:**

_If you are a project manager, make sure that you are logged out of the CLA Management Console before you begin._

1. In Gerrit, clone a repository under the Gerrit instance into your local machine.
2. Make a change and push the code to your Gerrit repository.

   A warning link that you need to sign a CLA appears:

   ![Gerrit Warning Link](../.gitbook/assets/cla-gerrit-sign-a-cla.png)

3. Navigate to the Gerrit instance of your project. For example, if you are contrinuting to OPNFV project, navigate to [https://gerrit.opnfv.org](https://gerrit.opnfv.org)
4. Log in using your LFID.
5. Navigate to **Settings**— the gear icon on the upper right corner, and click **Agreements** from the menu on the left:

   ![Settings Icon](/.gitbook/assets/settings-icon.png)
   
   ![Gerrit Agreements](/.gitbook/assets/agreements.png)

6. Click **New Contributor Agreement**.

   ![Gerrit Agreements](/.gitbook/assets/agreement-link.png)

   New Contributor Agreement types appear:

   ![Gerrit New Contributor Agreement](/.gitbook/assets/new-contributor-agreement.png)

7. Continue to contribute as an individual or employee \(corporate contributor\):
   * [Individual Contributor](contribute-to-a-gerrit-project.md#individual-contributor)
   * [Corporate Contributor](contribute-to-a-gerrit-project.md#corporate-contributor)

## Individual Contributor

EasyCLA presents a review agreement link to individual contributors to open a CLA and sign it.

1. Select an individual CLA type.

   ![New Contributor Agreement](../.gitbook/assets/cla-gerrit-icla-type.png)

2. Click the **Please review the agreement link** and then click the message link that appears:

   ![Gerrit Sign ICLA Link](../.gitbook/assets/cla-gerrit-icla-proceed-to-sign-cla.png)

3. Log in to EasyCLA if you are prompted.
4. Click **OPEN CLA** on the dialog that appears:

   ![Gerrit Open CLA](../.gitbook/assets/cla-gerrit-individual-cla-open-cla.png)

   DocuSign presents the agreement that you must sign. The Individual CLA is not tied to any employer you may have, so enter your @personal address in the E-Mail field.

5. Follow the instructions in the DocuSign document, sign it, and click **FINISH**.

You are redirected to Gerrit. Wait a few seconds for the CLA status to update.

## Corporate Contributor

EasyCLA presents a review agreement link where you confirm your association with the company.

1. Select **Corporate CLA**.

   ![Gerrit Corporate CLA Agreement](/.gitbook/assets/corporate-cla.png)

2. Click the **Please review the agreement link** and then click the message link that appears:

   ![Gerrit Sign ICLA Link](../.gitbook/assets/cla-gerrit-icla-proceed-to-sign-cla.png)
   
3. Sign in to EasyCLA if you are prompted.
   
4. Select **Company**.

   To contribute to this project, you must be authorized under a signed Contributor License Agreement. You are contributing on behalf of your work for a company.


5. Continue:

   * [If a **Confirmation of Association with** statement appears](contribute-to-a-gerrit-project.md#if-a-confirmation-of-association-with-statement-appears)
   * [If a **Company has not signed CCLA** window appears](contribute-to-a-gerrit-project.md#if-a-company-has-not-sigend-CCLA-window-appears)
   * [If You are not whitelsited](contribute-to-a-gerrit-project.md#if-you-are-not-whitelisted)
   * [If Company is not in the list](contribute-to-a-gerrit-project.md#if-company-is-not-in-the-list)
   
  
## If a **Confirmation of Association with** statement appears

1. Read the Confirmation of Association statement and select the checkbox.

   ![Gerrit Confirmation of Association](../.gitbook/assets/cla-gerrit-confirmation-of-association.png)

2. Click **CONTINUE**.

   A dialog appears and informs you: You are done!
   
   ![You are done](../.gitbook/assets/cla-github-you-are-done%20%281%29.png)

3. Click **RETURN TO REPO**.

   You are redirected to Gerrit. Wait a few seconds for the CLA status to update or refresh the page.
   
## If a **Company has not signed CCLA** window appears

This window appears if your comapny has not signed a Corporate CLA for the project.

To send an email notification to your company's CLA manager to sign Corporate CLA:

1. Select your email address from the **Emial to Authorize** drop-down list. This is the email address that you want your company manager to whitelist while signing the Corporate CLA.

   ![](/.gitbook/assets/company-not-signed-ccla.png)

2. Click **SEND**.

  A message shows that your email is successfully sent.
  
  ![](/.gitbook/assets/email-to-whitelist.png)

## If You are not whitelisted

This window appears if your company has not whitelisted you under their signed Corporate CLA.

To send a request to your company's CLA manager to be whitelisted:

1. Click **CONTACT**.

   ![](/.gitbook/assets/request-to-be-whitelisted.png)

    A **Request Access** window appears.
  
2. Select your email address from the **Emial to Authorize** drop-down list. This is the email address that you want your company manager to whitelist while signing the Corporate CLA.

3. Click **SEND**

   A message shows that your email is successfully sent.
  
  ![](/.gitbook/assets/email-to-whitelist.png)

## If Company is not in the list

If you don't find your company's name in the list:

1. Click **COMPANY NOT IN LIST? CLICK HERE**.

   The **Verify Your Permission of Access** dialog appears.
  
2. Click an answer: Are You a CLA Manager?

   **YES**— You will be redirected to [corporate.lfcla.com](https://corporate.lfcla.com/#/companies) to [add your company](../ccla-managers-and-ccla-signatories/add-a-company-to-a-project.md) to a project.
   
   **NO**— A Request Access form appears. Continue to next step.
   
3. Complete the form and click **SEND**.

   The CCLA manager signs a Corporate CLA and adds you to the whitelist.

You have finished signing your CLA for this Gerrit instance. You are able to submit your changes to any repository under this Gerrit instance.

