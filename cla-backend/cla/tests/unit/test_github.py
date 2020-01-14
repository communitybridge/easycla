# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

from unittest.mock import patch

import pytest

from cla.utils import get_comment_body, get_comment_badge

SUCCESS = ":white_check_mark:"
FAILED = ":x:"
SIGN_URL = "http://test.contributor.lfcla/sign"
SUPPORT_URL = "https://jira.linuxfoundation.org/servicedesk/customer/portal/4"
GITHUB_HELP_URL = (
    "https://help.github.com/en/github/committing-changes-to-your-project/why-are-my-commits-linked-to-the-wrong-user"
)
GITHUB_FAKE_SHA = "fake_sha"
GITHUB_FAKE_SHA_2 = "fake_sha_2"


def test_get_comment_body_no_user_id():
    """
    Test CLA comment body for case CLA test failure when commit has no user ids
    """
    # case with missing list with no authors
    response = get_comment_body("github", SIGN_URL, [], [(GITHUB_FAKE_SHA, [None,'foo','foo@bar.com']), (GITHUB_FAKE_SHA_2, [None,'fake','fake@gmail.com'])])
    expected = (
        f"<ul><li>"
        + FAILED
        + "The commit ("
        + " ,".join([GITHUB_FAKE_SHA, GITHUB_FAKE_SHA_2])
        + ") is missing the User's ID,preventing the EasyCLA check.[Consult GitHub Help]("
        + GITHUB_HELP_URL
        + ") to resolve."
        + "</li></ul>"
    )
    assert response == expected


def test_get_comment_body_cla_fail_no_user_id_and_user_id():
    """
    Test CLA comment body for case CLA fail check with no user id and existing user id
    """
    # case with missing list with user id existing
    author_name = "wanyaland"
    response = get_comment_body("github", SIGN_URL , [], [(GITHUB_FAKE_SHA, ['12', author_name, 'foo@gmail.com']), (GITHUB_FAKE_SHA_2, [None, author_name,' foo@gmail.com'])])
    expected = (
        f"<ul><li>"
        + "["
        + FAILED
        + "]("
        + SIGN_URL
        + ")  "
        + author_name
        + " The commit ("
        + " ,".join([GITHUB_FAKE_SHA])
        + ") is not authorized under a signed CLA. "
        + f"[Please click here to be authorized]({SIGN_URL}). For further assistance with "
        + f"EasyCLA, [please submit a support request ticket]({SUPPORT_URL})."
        + "</li>"
        + "<li>"
        + FAILED
        + "The commit ("
        + " ,".join([GITHUB_FAKE_SHA_2])
        + ") is missing the User's ID,preventing the EasyCLA check.[Consult GitHub Help]("
        + GITHUB_HELP_URL
        + ") to resolve."
        + "</li></ul>"
    )

    assert response == expected


def test_get_comment_badge_with_no_user_id():
    """
    Test CLA badge for CLA fail check with no user id
    """
    missing_id_badge = "cla-missing-id.png"
    missing_user_id = True
    all_signed = False
    response = get_comment_badge("github", all_signed, SIGN_URL, missing_user_id=missing_user_id)
    assert missing_id_badge in response