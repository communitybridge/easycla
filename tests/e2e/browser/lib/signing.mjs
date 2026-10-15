// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
//
// Me lens (M1/M2) signing helpers shared by m3_chain.mjs and m2_signing.mjs: the "Sign CLA" picker,
// the hand-off to the contributor console, the corporate (ECLA click-through) and individual (ICLA,
// DocuSign sandbox) console flows and the agreements read.
import * as e2e from './e2e.mjs';

const { SS, sleep, squash } = e2e;
const host = (key) => { try { return new URL(e2e.cfg(key)).hostname; } catch { return ''; } };
const legacyHosts = ['LEGACY_API', 'LEGACY_API_ALIAS'].map(host).filter(Boolean);
const apiGwHost = host('API_GW');
const consoleHost = host('CONTRIBUTOR_CONSOLE');
const esc = (s) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

// writes the Me lens -> console flows make (guard keys are "METHOD host/path")
export const signingWrites = (kind) => {
  const out = [
    new RegExp(`^POST ${esc(e2e.SS_HOST)}/api/me/clas/prepare-sign$`),
    new RegExp(`^POST (${[...legacyHosts, apiGwHost, consoleHost].filter(Boolean).map(esc).join('|')})(/cla-service)?/v[234]/user-from-(token|session)$`),
  ];
  if (kind === 'ecla') out.push(new RegExp(`^POST (${legacyHosts.map(esc).join('|')})/v2/(check-prepare-employee-signature|request-employee-signature)$`));
  if (kind === 'icla') out.push(new RegExp(`^POST (${[...legacyHosts, apiGwHost].map(esc).join('|')})(/cla-service)?/v[24]/request-individual-signature$`));
  return out;
};
// never allowed from these flows, whatever the gates say
export const SIGNING_FORBIDDEN = /request-corporate-signature|\/company$|invite-company-admin|request-company-admin|cla-manager\/requests|\/signed\//;

export const myAgreements = async (context, groupId) => {
  const r = await e2e.bff(context, 'GET', '/api/me/clas');
  const rows = (r.json?.agreements || []).map((a) => ({ id: a.id, kind: a.kind, status: a.status, claGroupId: a.claGroupId, company: a.companyName || a.company || null }));
  return { http: r.status, rows, fixture: groupId ? rows.filter((a) => a.claGroupId === groupId) : rows };
};

// Me lens: Sign CLA -> pick the group -> (contract type) -> identity -> hand-off to the contributor console.
// Returns { picker, idText, url } ; throws when a dialog does not show up.
export async function handOffToConsole(session, identity, { groupId, groupName, githubId, contract = /employee|ecla|corporate/i }) {
  const { page } = session;
  await page.goto(`${SS}/profile/clas`, { waitUntil: 'domcontentloaded', timeout: 60000 });
  await sleep(6000);
  await e2e.hideOverlays(page);
  const signAction = page.getByTestId('sign-cla-action');
  await signAction.waitFor({ state: 'visible', timeout: 60000 });
  await signAction.click();
  const gdlg = page.getByTestId('cla-group-select-dialog');
  await gdlg.waitFor({ state: 'visible', timeout: 20000 });
  const search = page.getByTestId('cla-group-select-search').locator('input').first();
  await search.click();
  await search.pressSequentially(groupName, { delay: 60 });
  await sleep(5000);
  await e2e.shot(page, 'picker-results');
  const results = page.getByTestId('cla-group-select-results');
  const picker = squash(await results.innerText().catch(() => 'NONE')).slice(0, 300);
  await results.getByTestId(`cla-group-select-${groupId}`).click();
  await sleep(1000);
  await page.getByTestId('cla-group-select-continue').locator('button').click();
  let state = 'none';
  for (let i = 0; i < 20 && state === 'none'; i++) {
    await sleep(1500);
    if (await page.getByTestId('sign-contract-type-select-dialog').isVisible().catch(() => false)) state = 'contract';
    else if (await page.getByTestId('sign-identity-select-dialog').isVisible().catch(() => false)) state = 'identity';
  }
  await e2e.shot(page, 'picker-post-continue');
  if (state === 'none') throw new Error(`neither the contract-type nor the identity dialog appeared after the group pick (results: ${picker})`);
  if (state === 'contract') {
    const tdlg = page.getByTestId('sign-contract-type-select-dialog');
    const card = tdlg.locator('label, [class*="card"], [role="radio"], li').filter({ hasText: contract }).first();
    if (await card.count()) await card.click();
    await sleep(800);
    await tdlg.getByTestId('sign-contract-type-select-continue').locator('button').click();
    await page.getByTestId('sign-identity-select-dialog').waitFor({ state: 'visible', timeout: 20000 });
  }
  const idlg = page.getByTestId('sign-identity-select-dialog');
  const idText = squash(await idlg.innerText()).slice(0, 300);
  const gh = idlg.getByTestId(`sign-identity-select-github-${githubId}`);
  if (githubId && (await gh.count())) await gh.click();
  else await idlg.locator('[data-testid^="sign-identity-select-"]').filter({ hasText: /github/i }).first().click();
  await sleep(800);
  await e2e.shot(page, 'identity-dialog');
  await idlg.getByTestId('sign-identity-select-continue').locator('button').click();
  // SS -> (SSO again?) -> contributor console
  await page.waitForURL((u) => u.hostname === consoleHost || e2e.isAuth0Host(u.hostname), { timeout: 90000 });
  if (e2e.isAuth0Host(new URL(page.url()).hostname)) {
    await sleep(2000);
    if (e2e.isAuth0Host(new URL(page.url()).hostname)) await e2e.fillLogin(page, identity);
    await page.waitForURL((u) => u.hostname === consoleHost, { timeout: 90000 });
  }
  await sleep(9000);
  await e2e.hideOverlays(page);
  await e2e.shot(page, 'console-landing');
  return { picker, idText, url: page.url() };
}

// Console corporate flow: pick the company, tick the certification, Proceed -> POST /v2/request-employee-signature.
export async function corporateAcknowledge(session, companyName) {
  const { page } = session;
  const corpBtn = page.getByRole('button', { name: /proceed as a corporate contributor/i }).first();
  await corpBtn.waitFor({ state: 'visible', timeout: 30000 });
  await corpBtn.click();
  await sleep(7000);
  await e2e.hideOverlays(page);
  await e2e.shot(page, 'console-corporate');
  const compInput = page.locator('input[type="text"]:visible, input[type="search"]:visible').first();
  if (await compInput.isVisible().catch(() => false)) {
    await compInput.click();
    await compInput.fill('');
    await compInput.pressSequentially(companyName, { delay: 70 });
    await sleep(5000);
    await e2e.shot(page, 'console-company-results');
    const option = page.getByText(new RegExp(`^\\s*${esc(companyName)}\\s*$`, 'i')).first();
    await option.waitFor({ state: 'visible', timeout: 20000 });
    await option.click();
    await sleep(1500);
  }
  const compliance = page.locator('app-consent input[type="checkbox"]');
  if ((await compliance.count()) !== 1) throw new Error(`compliance checkbox not found (count=${await compliance.count()})`);
  if (!(await compliance.isChecked())) await page.locator('app-consent label').filter({ hasText: /i hereby certify/i }).click();
  if (!(await compliance.isChecked())) throw new Error('compliance checkbox could not be ticked');
  await e2e.shot(page, 'console-company-picked');
  const proceed = page.getByRole('button', { name: /^proceed$/i, disabled: false });
  await proceed.waitFor({ state: 'visible', timeout: 10000 });
  const [res] = await Promise.all([
    page.waitForResponse((r) => new URL(r.url()).pathname.endsWith('/v2/request-employee-signature') && r.request().method() === 'POST', { timeout: 90000 }),
    proceed.click(),
  ]);
  const receipt = await res.json().catch(() => ({}));
  const done = await page.getByText('You are done!', { exact: true }).waitFor({ state: 'visible', timeout: 30000 }).then(() => true).catch(() => false);
  await e2e.shot(page, 'console-done');
  return { status: res.status(), receipt, done, signatureId: receipt.signature_id || receipt.signatureID || null };
}

// "Exit EasyCLA" back to the Self Serve Me lens.
export async function exitConsole(session) {
  const { page } = session;
  await page.getByRole('button', { name: /^exit easycla$/i }).click().catch(() => {});
  await page.waitForURL((u) => u.hostname === e2e.SS_HOST, { timeout: 60000 }).catch(() => {});
  await sleep(5000);
}

// Me lens row text for one agreement id (reloads /profile/clas).
export async function meLensRow(session, id) {
  const { page } = session;
  await page.goto(`${SS}/profile/clas`, { waitUntil: 'domcontentloaded', timeout: 60000 });
  await page.locator('[data-testid^="agreement-row-"]:not([data-testid^="agreement-row-actions-"]), [data-testid="my-clas-error"], [data-testid^="my-clas-empty"]').first().waitFor({ timeout: 90000 }).catch(() => {});
  await e2e.hideOverlays(page);
  return squash(await page.getByTestId(`agreement-row-${id}`).innerText({ timeout: 30000 }).catch(() => 'NOT RENDERED')).slice(0, 200);
}

// ---- DocuSign sandbox (dev only; the guard aborts every DocuSign host unless launch({allowDocusign}) on dev) ----
const DS_HOSTS = new Set(['demo.docusign.net', 'apps-d.docusign.com', 'account-d.docusign.com']);
export const onDocusign = (page) => DS_HOSTS.has(new URL(page.url()).hostname);

// Walk a DocuSign sandbox envelope: disclosure -> fill the empty text fields -> sign tab -> adopt -> finish.
// Returns { signed, finished, finalUrl }; throws when the page is not a DocuSign sandbox page.
export async function docusignSign(page, values = ['NONBINDING E2E TEST - 123 Test Street', 'Test City 00000', 'United States', 'Test', 'Test2']) {
  await page.waitForURL((u) => DS_HOSTS.has(u.hostname), { timeout: 90000 });
  if (!onDocusign(page)) throw new Error(`not on a DocuSign sandbox host: ${page.url().slice(0, 80)}`);
  await page.setViewportSize({ width: 1440, height: 1000 }).catch(() => {});
  await sleep(12000);
  await e2e.shot(page, 'docusign-landing');
  for (let a = 0; a < 5; a++) {
    const agreeLbl = page.getByText(/i agree to use electronic records/i).first();
    if (!(await agreeLbl.isVisible().catch(() => false))) break;
    const box = page.locator('#disclosureAccepted').first();
    if (!(await box.isChecked().catch(() => false))) {
      await agreeLbl.click({ force: true }).catch(() => {});
      await sleep(700);
      if (!(await box.isChecked().catch(() => false))) { await box.click({ force: true }).catch(() => {}); await sleep(700); }
    }
    const cont = page.getByRole('button', { name: /^continue$/i }).first();
    if (await cont.isVisible().catch(() => false)) { await cont.click().catch(() => {}); await sleep(4000); }
  }
  const inputs = page.locator('input[type="text"]:visible, textarea:visible');
  const n = await inputs.count();
  let vi = 0;
  for (let i = 0; i < n && vi < values.length; i++) {
    const ti = inputs.nth(i);
    if ((await ti.inputValue().catch(() => 'x')) === '' && (await ti.isEditable().catch(() => false))) {
      await ti.click().catch(() => {});
      await ti.fill(values[vi]).catch(() => {});
      vi++;
      await sleep(400);
    }
  }
  let signed = false;
  for (const sel of ['[data-qa^="signature-tab-required"]', '[data-qa^="signature-tab"]', '[data-qa="tab-SignHere"]', 'button[aria-label*="Sign Here"]', '.signature-tab']) {
    const t = page.locator(sel).first();
    if (await t.isVisible().catch(() => false)) { await t.click({ force: true }).catch(() => {}); signed = true; break; }
  }
  if (!signed) { const t = page.getByRole('button', { name: /^sign$/i }).first(); if (await t.isVisible().catch(() => false)) { await t.click({ force: true }); signed = true; } }
  await sleep(3000);
  const adopt = page.locator('[data-qa="adopt-dialog-adopt-and-sign-button"], button:has-text("Adopt and Sign")').first();
  if (await adopt.isVisible().catch(() => false)) {
    if (!onDocusign(page)) throw new Error('adopt dialog outside the DocuSign sandbox');
    await adopt.click();
    await sleep(4000);
  }
  await e2e.shot(page, 'docusign-signed');
  let finished = false;
  for (let k = 0; k < 6 && !finished; k++) {
    const fin = page.locator('[data-qa="action-bar-btn-finish"], [data-qa="envelope-finish"]').first();
    if ((await fin.isVisible().catch(() => false)) && (await fin.isEnabled().catch(() => false))) {
      if (!onDocusign(page)) throw new Error('finish button outside the DocuSign sandbox');
      await fin.click().catch(() => {});
      finished = true;
      await sleep(10000);
    } else await sleep(3000);
  }
  await page.waitForURL((u) => !DS_HOSTS.has(u.hostname), { timeout: 90000 }).catch(() => {});
  await sleep(8000);
  await e2e.shot(page, 'docusign-returned');
  return { signed, finished, filled: vi, finalUrl: page.url() };
}
