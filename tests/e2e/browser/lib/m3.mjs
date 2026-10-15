// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
//
// Organization lens -> EasyCLA (M3) page helpers shared by m3_org_lens.mjs, m3_chain.mjs and bridge.mjs.
// All selectors are the data-testid hooks of the Self Serve front end.
import * as e2e from './e2e.mjs';

const { SS, sleep, squash } = e2e;

// writes a module must declare for the fixture signature (guard key: "METHOD host/path")
export const fixtureWrites = (sfid, sig) => [
  new RegExp(`^(PUT|POST|DELETE) ${e2e.SS_HOST}/api/orgs/${sfid}/lens/cla-groups/${sig}(/|$)`),
];

export function m3(session, { sfid, name, group, sig }) {
  const { page, context, net } = session;
  const prefix = `/api/orgs/${sfid}/lens/cla-groups/${sig}`;
  const paths = {
    list: `/api/orgs/${sfid}/lens/cla-groups`,
    approval: `${prefix}/approval-list`,
    auto: `${prefix}/ecla-auto-create`,
    managers: `${prefix}/managers`,
    acks: `${prefix}/acknowledgments`,
    activity: `${prefix}/activity`,
  };
  const get = (p) => e2e.bff(context, 'GET', p);
  const writesSoFar = () => net.writes.filter((w) => w.url.includes(prefix)).length;
  const waitWrite = (test, method, timeout = 90000) => page.waitForResponse((r) => {
    let u; try { u = new URL(r.url()); } catch { return false; }
    return u.hostname === e2e.SS_HOST && test(u.pathname) && r.request().method() === method;
  }, { timeout });

  async function state() {
    const [list, managers, approval, acks] = await Promise.all([get(paths.list), get(paths.managers), get(paths.approval), get(paths.acks)]);
    const row = (list.json?.claGroups || []).find((r) => r.id === sig) || null;
    return {
      http: { list: list.status, managers: managers.status, approval: approval.status, acks: acks.status },
      row,
      managers: (managers.json?.managers || []).map((m) => m.lfUsername).sort(),
      entries: (approval.json?.entries || []).map((x) => ({ kind: x.kind, value: x.value })),
      acks: acks.json?.list || [],
    };
  }

  async function openDetail(cid = 'M3-OPEN-DETAIL') {
    await page.goto(`${SS}/org/easycla`, { waitUntil: 'domcontentloaded', timeout: 60000 });
    const selected = await e2e.ensureOrgSelected(page, sfid, name);
    await page.keyboard.press('Escape').catch(() => {});
    await page.locator('.org-selector-panel').waitFor({ state: 'hidden', timeout: 10000 }).catch(() => {});
    await e2e.hideOverlays(page);
    const card = page.locator(`[data-testid="org-easycla-card-link"][href*="${group}"]`).first();
    // the list is loaded once any card renders; only then decide whether the fixture card is on this page
    await page.locator('[data-testid="org-easycla-card-link"]').first().waitFor({ state: 'visible', timeout: 60000 }).catch(() => {});
    let seen = await e2e.visible(card, 10000);
    if (!seen) {
      // paginated / long lists (8 cards per page): filter through the native search input (the lfx-input-text host
      // element carries the placeholder too, so target the input itself), then fall back to paging with "Next"
      const search = page.locator('input#org-easycla-search-input, input[placeholder*="Search CLA"]').first();
      if (await search.count()) { await search.fill(e2e.cfg('M3_CLA_GROUP_NAME', e2e.cfg('CLA_GROUP_NAME'))); await sleep(3000); seen = await e2e.visible(card, 15000); }
      for (let i = 0; i < 10 && !seen; i++) {
        const next = page.getByTestId('org-easycla-next-page').getByRole('button');
        if (!(await next.count()) || !(await next.isEnabled().catch(() => false))) break;
        await next.click(); await sleep(2500);
        seen = await e2e.visible(card, 10000);
      }
    }
    await card.waitFor({ state: 'visible', timeout: 60000 });
    await card.click();
    await page.getByTestId('org-easycla-detail-status').waitFor({ timeout: 60000 });
    await e2e.hideOverlays(page);
    await e2e.closeOverlay(page);
    const status = squash(await page.getByTestId('org-easycla-detail-status').innerText());
    e2e.check(cid, /signed/i.test(status), `org=${selected} group card opened; status="${status}" url=${new URL(page.url()).pathname}`);
    return status;
  }

  async function tab(t) {
    await e2e.closeOverlay(page);
    await page.getByTestId(`org-easycla-detail-tab-${t}`).click();
    await sleep(2500);
    await e2e.closeOverlay(page);
  }

  const approvalDialog = () => page.getByTestId('org-easycla-approval-dialog');
  async function openAdd() {
    await e2e.closeOverlay(page);
    const add = page.getByTestId('org-easycla-approval-add').getByRole('button');
    await add.waitFor({ state: 'visible', timeout: 60000 });
    await add.click();
    await approvalDialog().waitFor({ state: 'visible', timeout: 20000 });
  }
  async function pickKind(re) {
    await page.getByTestId('org-easycla-approval-dialog-kind-0').getByRole('combobox').click();
    await page.getByRole('option', { name: re }).first().click();
  }
  const valueBox = () => approvalDialog().getByRole('textbox').first();
  const submitApproval = () => page.getByTestId('org-easycla-approval-dialog-submit').getByRole('button').click();
  async function cancelDialog() {
    await page.getByTestId('org-easycla-approval-dialog-cancel').getByRole('button').click();
    await approvalDialog().waitFor({ state: 'hidden', timeout: 10000 }).catch(() => {});
    await sleep(800);
  }

  // Add one approval-list entry through the dialog; returns the PUT response body.
  async function addEntry(kindRe, value, cid) {
    await openAdd();
    await pickKind(kindRe);
    await valueBox().fill(value);
    const [res] = await Promise.all([waitWrite((p) => p === paths.approval, 'PUT'), submitApproval()]);
    const body = await res.json().catch(() => ({}));
    const ok = res.status() === 200 && (body.entries || []).some((x) => x.value === value);
    await approvalRow(value).waitFor({ state: 'visible', timeout: 30000 }).catch(() => {});
    await e2e.shot(page, `${cid.toLowerCase()}`);
    e2e.check(cid, ok, `PUT approval-list (${value}) -> ${res.status()}; entries=${(body.entries || []).length}`);
    return body;
  }

  // Try the dialog with a value and expect an inline validation error and no write.
  async function rejectEntry(kindRe, value, cid, msgRe) {
    const w0 = writesSoFar();
    await openAdd();
    await pickKind(kindRe);
    await valueBox().fill(value);
    await submitApproval();
    const msg = squash(await page.getByTestId('org-easycla-approval-dialog-error-text-0').innerText({ timeout: 15000 }).catch(() => ''));
    await e2e.shot(page, `${cid.toLowerCase()}`);
    await cancelDialog();
    const ok = !!msg && (!msgRe || msgRe.test(msg)) && writesSoFar() === w0;
    e2e.check(cid, ok, `"${value}" -> inline error "${msg.slice(0, 100)}"; writes=${writesSoFar() - w0}`);
    return msg;
  }

  // the approval-list row whose value cell is exactly `value` (hasText would also match e.g. the e-mail row for a user name)
  const approvalRow = (value) => page.getByTestId('org-easycla-approval-row').filter({ has: page.getByText(value, { exact: true }) });

  async function removeEntry(value, cid) {
    const row = approvalRow(value);
    await row.waitFor({ state: 'visible', timeout: 30000 });
    await row.getByTestId('org-easycla-approval-delete').getByRole('button').click();
    const [res] = await Promise.all([waitWrite((p) => p === paths.approval, 'PUT'), page.getByRole('button', { name: 'Remove entry', exact: true }).click()]);
    const body = await res.json().catch(() => ({}));
    await row.waitFor({ state: 'hidden', timeout: 30000 }).catch(() => {});
    const ok = res.status() === 200 && !(body.entries || []).some((x) => x.value === value);
    e2e.check(cid, ok, `remove ${value} -> ${res.status()}; entries=${(body.entries || []).length}`);
    return body;
  }

  async function toggleAuto(expected, cid) {
    await tab('overview');
    const toggle = page.getByRole('switch', { name: /auto.?ecla/i }).first();
    await toggle.waitFor({ state: 'visible', timeout: 30000 });
    await e2e.closeOverlay(page);
    const before = await toggle.getAttribute('aria-checked');
    if (String(expected) === before) {
      e2e.info(cid, `switch already aria-checked=${before}; nothing to toggle`);
      return { autoCreateEcla: expected, unchanged: true };
    }
    const [res] = await Promise.all([waitWrite((p) => p === paths.auto, 'PUT'), toggle.click()]);
    const body = await res.json().catch(() => ({}));
    await sleep(1500);
    await e2e.shot(page, `${cid.toLowerCase()}`);
    e2e.check(cid, res.status() === 200 && body.autoCreateEcla === expected, `PUT ecla-auto-create -> ${res.status()} autoCreateEcla=${body.autoCreateEcla} (aria before=${before} after=${await toggle.getAttribute('aria-checked')})`);
    return body;
  }

  async function pollAcks(predicate, tries = 30, every = 3000) {
    let last = [];
    for (let i = 0; i < tries; i++) {
      const acks = await get(paths.acks);
      last = acks.json?.list || [];
      if (acks.status === 200 && predicate(last)) return last;
      await sleep(every);
    }
    return last;
  }

  const ackValid = (a) => a.approved !== false && !a.invalidatedOn && !a.dateInvalidated && !/invalid/i.test(String(a.state || a.status || ''));
  const ackOf = (a, who, github) => a.lfLogin === who || (github && a.githubUsername === github) || JSON.stringify(a).includes(`"${who}"`);
  const ackId = (a) => a.signatureId || a.id;
  const ackState = (a) => (a ? (a.state || a.status || (a.approved === false ? `invalidated(by ${a.invalidatedBy || '?'} at ${a.invalidatedAt || '?'})` : 'approved')) : 'missing');

  async function invalidateAck(id, cid) {
    await tab('acknowledgments');
    await sleep(3000);
    const row = page.getByTestId(`org-easycla-acknowledgment-row-${id}`);
    const seen = await e2e.visible(row, 60000);
    if (!seen) { e2e.check(cid, false, `acknowledgment row ${id} not rendered`); return null; }
    await row.getByTestId('org-easycla-acknowledgment-invalidate').getByRole('button').click();
    const dialog = page.getByTestId('org-easycla-invalidate-dialog');
    await dialog.waitFor({ state: 'visible', timeout: 20000 });
    // The dialog first checks the approval list for entries added specifically for this contributor; when it finds
    // some it offers (default ON) to remove them together with the invalidation - that is a second PUT approval-list.
    await dialog.getByTestId('org-easycla-invalidate-dialog-checking').waitFor({ state: 'hidden', timeout: 30000 }).catch(() => {});
    const matches = [];
    for (const li of await dialog.getByTestId('org-easycla-invalidate-dialog-match').all()) matches.push((await li.innerText()).replace(/\s+/g, ' ').trim());
    const alsoRemoveBox = dialog.getByTestId('org-easycla-invalidate-dialog-also-remove').locator('input[type="checkbox"]').first();
    const alsoRemoveOffered = (await alsoRemoveBox.count()) > 0;
    const alsoRemove = alsoRemoveOffered ? await alsoRemoveBox.isChecked().catch(() => false) : false;
    await e2e.shot(page, `${cid.toLowerCase()}-dialog`);
    const [res, removal] = await Promise.all([
      waitWrite((p) => p.startsWith(paths.acks), 'POST').catch(() => null),
      matches.length && alsoRemove ? waitWrite((p) => p === paths.approval, 'PUT', 30000).catch(() => null) : Promise.resolve(null),
      dialog.getByTestId('org-easycla-invalidate-dialog-confirm').getByRole('button').click(),
    ]);
    await sleep(3000);
    await e2e.shot(page, `${cid.toLowerCase()}`);
    const badge = await row.getByTestId('org-easycla-acknowledgment-state-invalidated').count().catch(() => 0);
    const after = ((await get(paths.acks)).json?.list || []).find((a) => ackId(a) === id);
    e2e.check(cid, !!res && res.status() < 300 && badge >= 1 && !!after && !ackValid(after), `POST invalidate -> ${res ? res.status() : 'no response'}; badge=${badge}; state=${ackState(after)}; dialog matches=${JSON.stringify(matches)} alsoRemove=${alsoRemoveOffered ? alsoRemove : 'n/a'}`);
    const removed = removal ? (await removal.json().catch(() => ({}))) : null;
    return { after, matches, alsoRemoveOffered, alsoRemove, removal: removal ? { status: removal.status(), entries: (removed?.entries || []).length } : null };
  }

  return { prefix, paths, get, state, openDetail, tab, approvalRow, openAdd, pickKind, valueBox, submitApproval, cancelDialog, addEntry, rejectEntry, removeEntry, toggleAuto, pollAcks, ackValid, ackOf, ackId, ackState, invalidateAck, waitWrite, writesSoFar };
}
