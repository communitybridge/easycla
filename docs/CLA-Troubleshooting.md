# CLA Troubleshooting
Troubleshooting helps you solve problematic symptoms in your CLA implementation.

If you are having issues with EasyCLA, go to: <https://support.linuxfoundation.org> and file a ticket.

## CLA Management Console Data Does Not Load

The CLA Management Console data may not load due to a bug in the Auth0 implementation.

**Solution**:

1. Open a Chrome window, and then type `command + option + i`.

   The Chrome developer panel appears.

2. Select the **Application** tab.

3. Select **Clear storage** under Application in the left pane.

4. Select **Clear site data** from the bottom of the developer console.

5. Sign out of the CLA Management Console.

6. Sign back in.

If the issue persists, try using an incognito browser window.

## CLA Manager Does Not Receive Email Notifications

The CLA manager does not receive email notifications.

**Solution**:

Go to GitHub and make sure your company has an email address.