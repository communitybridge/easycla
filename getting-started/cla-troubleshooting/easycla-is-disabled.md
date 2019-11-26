---
description: >-
  EasyCLA is disabled so the organizations that I want EasyCLA to monitor are
  not monitored.
---

# EasyCLA is Disabled

EasyCLA is disabled so the organizations that I want EasyCLA to monitor are not monitored.

**Solution:**

This is a known issue. GitHub is set up to permit administrators and organization owners to have maximum flexibility, which includes disabling apps like EasyCLA. Do the following steps to mitigate this problem immediately. Be sure to educate your administrators and organization owners about this GitHub setup and solution.

**Do these steps:**

1. As the GitHub organization owner or administrator, go to the GitHub repository that you want EasyCLA to monitor.
2. Click **Settings** from the top menu.

   ![Settings](/docs/imgs/cla-github-repository-settings.png)

   Settings appear with Options in the left pane.

3. Click **Branches** under Options.

   ![Branches](/docs/imgs/cla-github-options.png)

   Branch settings appear.

4. Select **master** for the Default branch. **Edit** or **Add rule** for Branch protection rules of your organization.

   ![Branch Protection Rules](/docs/imgs/cla-github-branch-add-rule.png)

   Branch protection rule settings appear.

5. Select the following checkboxes in Rule settings and click **Create**.

   * Require status checks to pass before merging
   * Require branches to be up to date before merging
   * Include administrators

   ![Rule Settings](/docs/imgs/cla-github-branch-protection-rule.png)

