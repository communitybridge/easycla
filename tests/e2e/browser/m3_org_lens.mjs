// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
//
// M3: Organization lens -> EasyCLA -> signed agreement detail as ID1 (CLA manager) on the M3 fixture.
// Approval list (validation, add, duplicate, delete-cancel, remove), Auto-ECLA on/off, acknowledgment
// invalidation, CLA manager roster (validation, add/remove ID2, remove-cancel), activity log, end state.
// State-aware: suite-owned approval entries left by an earlier run are removed first, the manager
// add/remove order follows whether ID2 is already on the roster, and the end state must equal the baseline.
// Writes only under /api/orgs/<M3_ORG_SFID>/lens/cla-groups/<M3_CCLA>/ (guard-enforced).
import * as e2e from './lib/e2e.mjs';
import { m3, fixtureWrites } from './lib/m3.mjs';

const { cfg, check, info, skip, log } = e2e;
const ID1 = cfg('ID1');
const ID2 = cfg('ID2');
const fx = { sfid: cfg('M3_ORG_SFID'), name: cfg('M3_ORG_NAME'), group: cfg('M3_CLA_GROUP'), sig: cfg('M3_CCLA') };
const WRITES = e2e.gate('ALLOW_FIXTURE_WRITES') && !e2e.NO_WRITES;
const synthetic = `qa-e2e-suite-${e2e.TAG}.example.com`;
const suiteOwned = (v) => /^qa-e2e-suite-.*\.example\.com$/.test(v) || v === cfg('ID1_GITHUB') || v === cfg('ID1_EMAIL');

if (!fx.sfid || !fx.sig) { skip('M3-FIXTURE', 'M3_ORG_SFID / M3_CCLA not configured'); process.exit(e2e.finish()); }

const session = await e2e.launch({ identity: ID1, writes: WRITES ? fixtureWrites(fx.sfid, fx.sig) : [] });
const { page, net } = session;
const h = m3(session, fx);
try {
  const how = await e2e.login(page, ID1, `${e2e.SS}/org/easycla`);
  check('M3-LOGIN', new URL(page.url()).hostname === e2e.SS_HOST, `${ID1} login ${how}`);

  const base = await h.state();
  e2e.evidence('before.private.json', base);
  check('M3-BEFORE-STATE', !!base.row?.signed && Object.values(base.http).every((s) => s === 200),
    `row=${JSON.stringify(base.row && { status: base.row.status, claManagersCount: base.row.claManagersCount, approvalCriteriaCount: base.row.approvalCriteriaCount, autoCreateEcla: base.row.autoCreateEcla, viewerIsClaManager: base.row.viewerIsClaManager })} managers=${JSON.stringify(base.managers)} entries=${base.entries.length} acks=${base.acks.length}`);
  check('M3-VIEWER-IS-CLA-MANAGER', base.row?.viewerIsClaManager === true, `viewerIsClaManager=${base.row?.viewerIsClaManager}`);
  const expectedEntries = base.entries.filter((x) => !suiteOwned(x.value)).map((x) => `${x.kind}:${x.value}`).sort();
  const id2OnRoster = base.managers.includes(ID2);
  info('M3-BASELINE-EXPECTATION', `end state must be: entries=${JSON.stringify(expectedEntries)} managers=${JSON.stringify(base.managers)} auto=${base.row?.autoCreateEcla}; ID2 on roster at start=${id2OnRoster}`);

  await h.openDetail();
  for (const t of ['overview', 'approval', 'acknowledgments', 'managers', 'activity']) {
    await h.tab(t);
    const alert = e2e.squash(await page.locator('[role="alert"]:visible').allInnerTexts().then((a) => a.join(' ')).catch(() => ''));
    const text = await e2e.bodyText(page);
    check(`M3-TAB-${t.toUpperCase()}`, text.length > 0 && !/something went wrong|failed to load|error loading|unexpected error/i.test(alert + ' ' + text), `tab ${t} renders; alert="${alert.slice(0, 80)}"`);
  }

  if (!WRITES) {
    skip('M3-WRITES', `ALLOW_FIXTURE_WRITES=${cfg('ALLOW_FIXTURE_WRITES')} E2E_NO_WRITES=${e2e.NO_WRITES ? 1 : 0}: approval/auto-ECLA/ack/manager writes skipped`);
  } else {
    // ---- approval list ----
    await h.tab('approval');
    const stale = base.entries.filter((x) => suiteOwned(x.value));
    for (const x of stale) await h.removeEntry(x.value, `M3-P0-STALE-ENTRY-REMOVED-${x.value.replace(/[^a-z0-9]/gi, '').slice(0, 24).toUpperCase()}`);
    if (stale.length) info('M3-P0-STALE-ENTRIES', `${stale.length} suite-owned entr(ies) from an earlier run removed first: ${JSON.stringify(stale.map((x) => x.value))}`);
    await e2e.shot(page, 'approval-start');
    await h.rejectEntry(/^email domain$/i, 'not a domain', 'M3-APPROVAL-INVALID-DOMAIN');
    await h.addEntry(/^email domain$/i, synthetic, 'M3-APPROVAL-ADD-DOMAIN');
    await h.rejectEntry(/^email domain$/i, synthetic, 'M3-APPROVAL-DUPLICATE', /already/i);
    await h.addEntry(/github user/i, cfg('ID1_GITHUB'), 'M3-APPROVAL-ADD-GITHUB-USERNAME');
    await h.addEntry(/^email$/i, cfg('ID1_EMAIL'), 'M3-APPROVAL-ADD-EMAIL');
    const w = h.writesSoFar();
    const domainRow = h.approvalRow(synthetic);
    await domainRow.getByTestId('org-easycla-approval-delete').getByRole('button').click();
    await page.getByRole('button', { name: 'Cancel', exact: true }).click();
    await e2e.sleep(1000);
    check('M3-APPROVAL-DELETE-CANCEL', h.writesSoFar() === w && (await domainRow.count()) === 1, `Cancel keeps the row (rows=${await domainRow.count()}) and makes no write (writes=${h.writesSoFar() - w})`);
    const afterAdds = await h.state();
    check('M3-APPROVAL-STATE-3-NEW-ENTRIES', afterAdds.entries.length === expectedEntries.length + 3, `entries=${JSON.stringify(afterAdds.entries.map((x) => `${x.kind}:${x.value}`))}`);

    // ---- Auto-ECLA ----
    const acks0 = afterAdds.acks;
    const mine0 = acks0.filter((a) => h.ackOf(a, ID1, cfg('ID1_GITHUB')));
    const validBefore = mine0.some(h.ackValid);
    const invalidatedBefore = mine0.length > 0 && !validBefore;
    await h.toggleAuto(true, 'M3-AUTO-ECLA-ON');
    const acks1 = await h.pollAcks((list) => list.length > acks0.length, invalidatedBefore || validBefore ? 8 : 30);
    e2e.evidence('acks-after-auto-on.private.json', acks1);
    const mine1 = acks1.filter((a) => h.ackOf(a, ID1, cfg('ID1_GITHUB')));
    if (invalidatedBefore) {
      // signatures/service.go: an explicitly invalidated ECLA is left alone, it needs an explicit re-approval
      check('M3-AUTO-ECLA-LEAVES-INVALIDATED-ACK-ALONE', acks1.length === acks0.length && mine1.every((a) => !h.ackValid(a)), `acks before=${acks0.length} after=${acks1.length}; ID1's acks stay invalidated`);
    } else if (validBefore) {
      check('M3-AUTO-ECLA-NO-DUPLICATE-ACK', acks1.length === acks0.length, `ID1 already holds a valid ack; acks before=${acks0.length} after=${acks1.length}`);
    } else {
      check('M3-AUTO-ECLA-CREATES-ACK', mine1.some(h.ackValid), `acks before=${acks0.length} after=${acks1.length}; ID1 ack created=${mine1.some(h.ackValid)}`);
    }
    await h.toggleAuto(false, 'M3-AUTO-ECLA-OFF');

    // ---- acknowledgment invalidation (ID1's own valid ack when there is one) ----
    const acksNow = (await h.state()).acks;
    const valid = acksNow.filter((a) => h.ackOf(a, ID1, cfg('ID1_GITHUB')) && h.ackValid(a));
    if (valid.length) {
      await h.invalidateAck(h.ackId(valid[0]), 'M3-ACK-INVALIDATE');
    } else {
      info('M3-ACK-INVALIDATE', `no valid ID1 acknowledgment on the fixture (${acksNow.length} acks, ID1's=${acksNow.filter((a) => h.ackOf(a, ID1, cfg('ID1_GITHUB'))).length} all invalidated) - invalidation exercised by m3_chain/interactions`);
      await h.tab('acknowledgments');
      await e2e.shot(page, 'acks');
    }

    // ---- CLA managers ----
    await h.tab('managers');
    await page.getByTestId('org-easycla-managers-row').first().waitFor({ state: 'visible', timeout: 60000 });
    await e2e.sleep(1500);
    await e2e.shot(page, 'managers-start');
    const selfRow = page.getByTestId('org-easycla-managers-row').filter({ hasText: cfg('ID1_EMAIL') });
    const blockedBtn = selfRow.getByTestId('org-easycla-managers-remove-blocked');
    const rosterNow = (await h.state()).managers;
    if (rosterNow.length === 1) check('M3-MANAGERS-SELF-REMOVE-BLOCKED', (await blockedBtn.count()) >= 1, `sole manager: self-remove control blocked=${await blockedBtn.count()}`);
    else info('M3-MANAGERS-SELF-REMOVE-BLOCKED', `roster has ${rosterNow.length} managers (blocked control only shown for the sole manager); blocked=${await blockedBtn.count()} remove=${await selfRow.getByTestId('org-easycla-managers-remove').count()}`);

    const id2Row = () => page.getByTestId('org-easycla-managers-row').filter({ hasText: cfg('ID2_EMAIL') });
    async function removeId2(cid) {
      await id2Row().waitFor({ state: 'visible', timeout: 60000 });
      const removeBtn = id2Row().getByTestId('org-easycla-managers-remove').getByRole('button');
      await removeBtn.waitFor({ state: 'visible', timeout: 60000 });
      const w1 = h.writesSoFar();
      await removeBtn.click();
      await page.getByRole('button', { name: 'Cancel', exact: true }).click();
      await e2e.sleep(800);
      check('M3-MANAGER-REMOVE-CANCEL', h.writesSoFar() === w1, `Cancel in the remove dialog makes no write (writes=${h.writesSoFar() - w1})`);
      await removeBtn.click();
      const [removed] = await Promise.all([h.waitWrite((p) => p.startsWith(h.paths.managers + '/'), 'DELETE'), page.getByRole('button', { name: 'Remove', exact: true }).click()]);
      await e2e.sleep(2000);
      const roster = (await h.state()).managers;
      await e2e.shot(page, cid.toLowerCase());
      check(cid, [200, 204].includes(removed.status()) && !roster.includes(ID2), `DELETE managers/${ID2} -> ${removed.status()}; roster=${JSON.stringify(roster)}`);
    }
    async function addId2(cid) {
      await e2e.closeOverlay(page);
      await page.getByTestId('org-easycla-managers-add').getByRole('button').click();
      const dialog = page.getByTestId('org-easycla-add-manager-dialog');
      const submit = page.getByTestId('org-easycla-add-manager-submit').getByRole('button');
      await dialog.waitFor({ state: 'visible', timeout: 20000 });
      const fields = dialog.locator('input');
      const w0 = h.writesSoFar();
      await submit.click();
      const errs = await Promise.all(['first-name', 'last-name', 'email'].map((f) => e2e.visible(page.getByTestId(`org-easycla-add-manager-${f}-error`), 10000)));
      await fields.nth(0).fill(cfg('ID2_FIRST')); await fields.nth(1).fill(cfg('ID2_LAST')); await fields.nth(2).fill('not-an-email');
      await submit.click();
      const emailErr = await e2e.visible(page.getByTestId('org-easycla-add-manager-email-error'), 10000);
      await e2e.shot(page, 'manager-validation');
      check('M3-MANAGER-VALIDATION', errs.every(Boolean) && emailErr && h.writesSoFar() === w0, `required errors=${JSON.stringify(errs)} email-format error=${emailErr} writes=${h.writesSoFar() - w0}`);
      await fields.nth(2).fill(cfg('ID2_EMAIL'));
      const [added] = await Promise.all([h.waitWrite((p) => p === h.paths.managers, 'POST'), submit.click()]);
      const body = await added.json().catch(() => ({}));
      await e2e.sleep(3000);
      const roster = (await h.state()).managers;
      await e2e.shot(page, cid.toLowerCase());
      check(cid, [200, 201].includes(added.status()) && roster.includes(ID2), `POST managers (${cfg('ID2_EMAIL')}) -> ${added.status()} lfUsername=${body.lfUsername}; roster=${JSON.stringify(roster)}`);
    }
    if (id2OnRoster) { await removeId2('M3-MANAGER-REMOVE'); await addId2('M3-MANAGER-ADD-RESTORED'); }
    else { await addId2('M3-MANAGER-ADD'); await removeId2('M3-MANAGER-REMOVE'); }

    // ---- cleanup of the entries added above ----
    await h.tab('approval');
    for (const v of [synthetic, cfg('ID1_GITHUB'), cfg('ID1_EMAIL')]) {
      if ((await h.state()).entries.some((x) => x.value === v)) await h.removeEntry(v, `M3-APPROVAL-REMOVE-${v === synthetic ? 'DOMAIN' : v === cfg('ID1_EMAIL') ? 'EMAIL' : 'GITHUB-USERNAME'}`);
    }
  }

  // ---- activity + end state ----
  await h.tab('activity');
  await page.locator('[data-testid^="org-easycla-activity-log"]').first().waitFor({ timeout: 60000 }).catch(() => {});
  await e2e.shot(page, 'activity-end');
  const activity = await h.get(h.paths.activity);
  check('M3-ACTIVITY-LOG', activity.status === 200 && (activity.json?.list || []).length > 0, `GET activity -> ${activity.status}; entries=${(activity.json?.list || []).length}; newest=${JSON.stringify((activity.json?.list || []).slice(0, 3).map((a) => (a.summary || '').slice(0, 80)))}`);
  const end = await h.state();
  e2e.evidence('end.private.json', end);
  const endEntries = end.entries.map((x) => `${x.kind}:${x.value}`).sort();
  if (WRITES) {
    check('M3-END-STATE', JSON.stringify(endEntries) === JSON.stringify(expectedEntries) && JSON.stringify(end.managers) === JSON.stringify(base.managers) && end.row?.autoCreateEcla === base.row?.autoCreateEcla,
      `entries=${JSON.stringify(endEntries)} (expected ${JSON.stringify(expectedEntries)}) managers=${JSON.stringify(end.managers)} (baseline ${JSON.stringify(base.managers)}) auto=${end.row?.autoCreateEcla} (baseline ${base.row?.autoCreateEcla})`);
  }
  check('M3-NO-BLOCKED-WRITES', net.blocked.length === 0, `guard blocked ${net.blocked.length} request(s) ${JSON.stringify(net.blocked.slice(0, 3))}`);
  info('M3-WRITES-MADE', `${net.writes.length} business write(s): ${JSON.stringify([...new Set(net.writes.map((x) => `${x.method} ${x.url.replace(e2e.SS, '').replace(h.prefix, '…')}`))])}`);
} catch (err) {
  check('M3-HARNESS', false, `error: ${err.message} at ${new URL(page.url()).pathname}`);
  await e2e.shot(page, 'error');
} finally {
  await e2e.close(session);
}
process.exit(e2e.finish());
