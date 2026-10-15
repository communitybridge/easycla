#!/usr/bin/env python3
# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
"""Fixture writes as ID1 (CLA manager) on the QA company x CLA group, always restoring the baseline.

A  approval list: add ID1's own criteria, add six synthetic criteria types, invalid e-mail refused,
   duplicate add idempotent, remove the synthetic ones (ID1's own entries are restored to their baseline state)
AU auto-ECLA toggle on/off (restored to the starting value)
DS designee read, duplicate manager add, manager requests list
SA sanctioned company: every write refused with 403 (when SANCTIONED_* is configured)
RM/RA approval-list removal voids ID1's own valid ECLA on the group; re-adding the criterion recovers it
   (removal-only invalidation rule) - runs only when ID1 has a VALID ecla on the fixture group
F  final baseline comparison
"""
import json
import os
import sys
import time

sys.path.insert(0, os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "lib"))
import e2e  # noqa: E402

c = e2e.cfg
ID1 = c("ID1")
PROJ, GRP, CO, SFID, CCLA = (c("WRITE_PROJECT_SFID") or c("PROJECT_SFID"), c("WRITE_CLA_GROUP") or c("CLA_GROUP"), c("WRITE_COMPANY_ID") or c("COMPANY_ID"),
                             c("WRITE_COMPANY_SFID") or c("COMPANY_SFID"), c("WRITE_COMPANY_CCLA") or c("COMPANY_CCLA"))
APPROVAL = f"/v4/signatures/project/{PROJ}/company/{CO}/clagroup/{GRP}/approval-list"
AUTO = f"/v4/signatures/company/{CO}/clagroup/{GRP}/ecla-auto-create"
LISTS = ("emailApprovalList", "domainApprovalList", "githubUsernameApprovalList", "githubOrgApprovalList", "gitlabUsernameApprovalList", "gitlabOrgApprovalList")
SYN = {"email": f"{ID1.split('@')[0]}+qa-{e2e.TAG}@example.org", "domain": f"qa-{e2e.TAG}.example.org",
       "ghuser": c("QA_GITHUB_USER", "octocat"), "ghorg": c("QA_GITHUB_ORG", "github"),
       "gluser": c("QA_GITLAB_USER", "gitlab-bot"), "glorg": c("QA_GITLAB_ORG", "gitlab-org")}
SYN_KEYS = [("emailApprovalList", "email"), ("domainApprovalList", "domain"), ("githubUsernameApprovalList", "ghuser"),
            ("githubOrgApprovalList", "ghorg"), ("gitlabUsernameApprovalList", "gluser"), ("gitlabOrgApprovalList", "glorg")]

if not e2e.cfg("ALLOW_FIXTURE_WRITES", "0") == "1" or e2e.NO_WRITES:
    e2e.skip("FW-ALL", "fixture writes disabled (ALLOW_FIXTURE_WRITES / --no-writes)")
    e2e.summary_line()
    sys.exit(0)


def ccla():
    code, body = e2e.call(ID1, "GET", e2e.v4(f"/v4/signatures/project/{PROJ}/company/{CO}"))
    rows = (body or {}).get("signatures") or [] if isinstance(body, dict) else []
    return code, next((r for r in rows if r.get("signatureID") == CCLA), None)


def items(row, key):
    return sorted(x.get("approval_item") for x in (row or {}).get(key) or [])


def put_approval(data):
    return e2e.call(ID1, "PUT", e2e.v4(APPROVAL), data)


def my_ecla():
    code, body = e2e.call(ID1, "GET", e2e.v4("/v4/my-clas?pageSize=100"))
    rows = (body or {}).get("clas") or [] if isinstance(body, dict) else []
    wanted = next((r for r in rows if r.get("signatureID") == c("ID1_ECLA") and r.get("claGroupID") == GRP), None) if c("ID1_ECLA") else None
    if wanted:
        return wanted
    valid = [r for r in rows if r.get("claGroupID") == GRP and r.get("claType") == "ecla" and r.get("status") == "valid"]
    return valid[0] if valid else next((r for r in rows if r.get("claGroupID") == GRP and r.get("claType") == "ecla"), None)


# ---- baseline ----------------------------------------------------------------------
code, row0 = ccla()
if code != 200 or row0 is None:
    e2e.check("A0-BASELINE", False, f"cannot read fixture CCLA {CCLA[:8]} (http {code})")
    e2e.summary_line()
    sys.exit(1)
baseline = {k: items(row0, k) for k in LISTS}
auto0 = row0.get("autoCreateECLA")
own_email0 = c("ID1_EMAIL") in baseline["emailApprovalList"]
own_gh0 = c("ID1_GITHUB") in baseline["githubUsernameApprovalList"]
e2e.info("A0-BASELINE", f"auto={auto0} lists={ {k: len(v) for k, v in baseline.items()} } own-email={own_email0} own-github={own_gh0} acl={[a.get('lfUsername') or a.get('username') for a in row0.get('signatureACL') or []]}")

# ---- A: approval list ----------------------------------------------------------------
code, body = put_approval({"AddEmailApprovalList": [c("ID1_EMAIL")], "AddGithubUsernameApprovalList": [c("ID1_GITHUB")]})
_, row1 = ccla()
e2e.check("A1-ADD-OWN-EMAIL-GITHUB", code == 200 and c("ID1_EMAIL") in items(row1, "emailApprovalList") and c("ID1_GITHUB") in items(row1, "githubUsernameApprovalList"),
          f"PUT add own e-mail + github -> {code}")
# one PUT per criterion type: the approval-list PUT validates GitHub/GitLab names remotely and is NOT atomic
# (a 400 on one item still persists the valid ones), so types are exercised and cleaned up independently
added = {}
for key, syn, add_field, remove_field, strict in (
        ("emailApprovalList", "email", "AddEmailApprovalList", "RemoveEmailApprovalList", True),
        ("domainApprovalList", "domain", "AddDomainApprovalList", "RemoveDomainApprovalList", True),
        ("githubUsernameApprovalList", "ghuser", "AddGithubUsernameApprovalList", "RemoveGithubUsernameApprovalList", False),
        ("githubOrgApprovalList", "ghorg", "AddGithubOrgApprovalList", "RemoveGithubOrgApprovalList", False),
        ("gitlabUsernameApprovalList", "gluser", "AddGitlabUsernameApprovalList", "RemoveGitlabUsernameApprovalList", False),
        ("gitlabOrgApprovalList", "glorg", "AddGitlabOrgApprovalList", "RemoveGitlabOrgApprovalList", False)):
    code, body = put_approval({add_field: [SYN[syn]]})
    _, r = ccla()
    present = SYN[syn] in items(r, key)
    added[key] = (syn, remove_field, present)
    if strict:
        e2e.check(f"A2-ADD-{key}", code == 200 and present, f"PUT {add_field}=[{SYN[syn]}] -> {code}, present={present}", http=code)
    elif code == 200 and present:
        e2e.check(f"A2-ADD-{key}", True, f"PUT {add_field}=[{SYN[syn]}] -> 200, present=True", http=code)
    elif code == 400 and not present:
        e2e.info(f"A2-ADD-{key}", f"PUT {add_field}=[{SYN[syn]}] -> 400 (remote name validation refused the QA value: {e2e.redact(json.dumps(body), 120)})", http=code)
    else:
        e2e.check(f"A2-ADD-{key}", False, f"PUT {add_field}=[{SYN[syn]}] -> {code}, present={present}", http=code)
code, body = put_approval({"AddEmailApprovalList": ["not-an-email"]})
e2e.check("A3-INVALID-EMAIL-REJECTED", code in (400, 422), f"PUT invalid e-mail -> {code}")
code, body = put_approval({"AddDomainApprovalList": [SYN["domain"]]})
_, row3 = ccla()
e2e.check("A4-DUPLICATE-ADD-IDEMPOTENT", code in (200, 400, 409) and items(row3, "domainApprovalList").count(SYN["domain"]) == 1,
          f"PUT duplicate domain -> {code}, occurrences={items(row3, 'domainApprovalList').count(SYN['domain'])}")
# approval-list adds create user rows per criterion: visible through corporate-contributors (INFO)
code, body = e2e.call(ID1, "GET", e2e.v4(f"/v4/company/external/{SFID}/cla-group/{GRP}/corporate-contributors?companyID={CO}&pageSize=100"))
rows = (body or {}).get("list") or [] if isinstance(body, dict) else []
e2e.info("A5-CORPORATE-CONTRIBUTORS", f"{len(rows)} contributor rows after the adds (http {code})")
# cleanup synthetic entries, one type at a time
left = {}
for key, (syn, remove_field, present) in added.items():
    if not present:
        continue
    code, body = put_approval({remove_field: [SYN[syn]]})
    _, r = ccla()
    if SYN[syn] in items(r, key):
        left[key] = (code, SYN[syn])
e2e.check("A9-REMOVE-SYNTHETIC", not left, f"all synthetic criteria removed; left={left}")

# ---- AU: auto-ECLA toggle ------------------------------------------------------------
code, body = e2e.call(ID1, "PUT", e2e.v4(AUTO), {"auto_create_ecla": not bool(auto0)})
_, r = ccla()
e2e.check("AU1-AUTO-TOGGLE", code == 200 and r.get("autoCreateECLA") == (not bool(auto0)), f"PUT auto_create_ecla={not bool(auto0)} -> {code}, now {r.get('autoCreateECLA')}")
code, body = e2e.call(ID1, "PUT", e2e.v4(AUTO), {"auto_create_ecla": bool(auto0)})
_, r = ccla()
e2e.check("AU2-AUTO-RESTORED", code == 200 and r.get("autoCreateECLA") == bool(auto0), f"PUT auto_create_ecla={bool(auto0)} -> {code}, now {r.get('autoCreateECLA')}")

# ---- DS: designee / manager requests -----------------------------------------------------
e2e.expect("DS1-IS-DESIGNEE-READ", ID1, 200, "GET", f"/v4/company/{SFID}/user/{ID1}/claGroupID/{GRP}/is-cla-manager-designee")
first, last = (c("ID1_FULLNAME") + " ").split(" ", 1)
e2e.expect("DS2-DUPLICATE-MANAGER-ADD", ID1, (200, 400, 409), "POST", f"/v4/company/{CO}/project/{PROJ}/cla-manager", {"firstName": first, "lastName": last.strip() or "QA", "userEmail": c("ID1_EMAIL")})
e2e.expect("DS3-MANAGER-REQUESTS-LIST", ID1, (200, 403), "GET", f"/v4/company/{CO}/project/{PROJ}/cla-manager/requests")

# ---- SA: sanctioned company --------------------------------------------------------------------
if c("SANCTIONED_SFID"):
    sf, sid = c("SANCTIONED_SFID"), c("SANCTIONED_ID")
    for cid, method, path, data in [
        ("SA-DESIGNEE-PROJECT", "POST", f"/v4/company/{sf}/project/{PROJ}/cla-manager-designee", {"userEmail": c("ID1_EMAIL")}),
        ("SA-DESIGNEE-GROUP", "POST", f"/v4/company/{sf}/claGroup/{GRP}/cla-manager-designee", {"userEmail": c("ID1_EMAIL")}),
        ("SA-REQUEST-CREATE", "POST", f"/v4/company/{sf}/project/{PROJ}/cla-manager/requests", {"contactAdmin": False, "fullName": "QA Probe", "userEmail": c("ID1_EMAIL")}),
        ("SA-APPROVAL-LIST", "PUT", f"/v4/signatures/project/{PROJ}/company/{sid}/clagroup/{GRP}/approval-list", {"RemoveEmailApprovalList": [f"absent-{e2e.TAG}@nomatch.invalid"]}),
        ("SA-AUTO-CREATE", "PUT", f"/v4/signatures/company/{sid}/clagroup/{GRP}/ecla-auto-create", {"auto_create_ecla": False}),
        ("SA-MANAGER-CREATE", "POST", f"/v4/company/{sid}/project/{PROJ}/cla-manager", {"firstName": "QA", "lastName": "Probe", "userEmail": c("ID1_EMAIL")}),
    ]:
        ok, body = e2e.expect(cid, ID1, 403, method, path, data)
else:
    e2e.skip("SA-SANCTIONED", f"SANCTIONED_* not configured for {e2e.ENV}")

# ---- RM/RA: removal-only invalidation + recovery on ID1's own ECLA -------------------------------
ack = my_ecla()
if ack and ack.get("status") == "valid":
    sid = ack["signatureID"]
    code, body = put_approval({"RemoveEmailApprovalList": [c("ID1_EMAIL")], "RemoveGithubUsernameApprovalList": [c("ID1_GITHUB")]})
    time.sleep(4)
    after = my_ecla() or {}
    e2e.check("RM1-REMOVAL-VOIDS-OWN-ECLA", code == 200 and after.get("signatureID") == sid and after.get("status") != "valid",
              f"remove own criteria -> {code}; ecla {sid[:8]} status={after.get('status')}")
    code, body = put_approval({"AddEmailApprovalList": [c("ID1_EMAIL")], "AddGithubUsernameApprovalList": [c("ID1_GITHUB")]})
    time.sleep(4)
    after2 = my_ecla() or {}
    e2e.check("RA1-READD-RECOVERS-OWN-ECLA", code == 200 and after2.get("signatureID") == sid and after2.get("status") == "valid",
              f"re-add own criteria -> {code}; ecla {sid[:8]} status={after2.get('status')} (removal-voided acks are resurrected by a re-add)")
elif ack:
    e2e.skip("RM-RA-OWN-ECLA", f"{ID1}'s ECLA on the fixture group is {ack.get('status')} - removal/re-add rule needs a valid one")
else:
    e2e.skip("RM-RA-OWN-ECLA", f"{ID1} has no ECLA on the fixture group {GRP[:8]} on {e2e.ENV}")

# ---- restore ID1's own entries to the baseline ------------------------------------------------------
if not e2e.KEEP_STATE:
    rm = {}
    if not own_email0:
        rm["RemoveEmailApprovalList"] = [c("ID1_EMAIL")]
    if not own_gh0:
        rm["RemoveGithubUsernameApprovalList"] = [c("ID1_GITHUB")]
    if rm:
        code, body = put_approval(rm)
        e2e.info("F0-RESTORE-OWN-ENTRIES", f"removed own entries that were not in the baseline -> {code}")
# ---- final -------------------------------------------------------------------------------------------
_, rowF = ccla()
final = {k: items(rowF, k) for k in LISTS}
diff = {k: sorted(set(final[k]) ^ set(baseline[k])) for k in LISTS if set(final[k]) != set(baseline[k])}
if e2e.KEEP_STATE:
    e2e.info("F1-LISTS-STATE", f"--keep-state: diff vs baseline {diff}")
else:
    e2e.check("F1-LISTS-BASELINE", not diff, f"approval lists back to baseline; diff={diff}")
e2e.check("F2-AUTO-BASELINE", rowF.get("autoCreateECLA") == auto0, f"autoCreateECLA={rowF.get('autoCreateECLA')} (baseline {auto0})")
e2e.summary_line()
sys.exit(0 if e2e.COUNTS["FAIL"] == 0 else 1)
