package templates

const dashboardHTML = `<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>ZFS Dashboard</title>
<style>
/* ── Design Tokens ─────────────────────────────────── */
:root {
  --font-body: 'DM Sans', system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  --font-mono: 'DM Mono', 'JetBrains Mono', ui-monospace, 'Cascadia Code', 'Fira Code', monospace;

  --text-xs:   clamp(0.72rem,  0.68rem + 0.2vw, 0.8rem);
  --text-sm:   clamp(0.82rem,  0.78rem + 0.2vw, 0.9rem);
  --text-base: clamp(0.9rem,   0.85rem + 0.25vw, 1rem);
  --text-lg:   clamp(1.05rem,  0.95rem + 0.5vw, 1.35rem);
  --text-xl:   clamp(1.3rem,   1.1rem  + 1vw,   2rem);

  --space-1:0.25rem; --space-2:0.5rem;  --space-3:0.75rem;
  --space-4:1rem;    --space-5:1.25rem; --space-6:1.5rem;
  --space-8:2rem;    --space-10:2.5rem; --space-12:3rem;

  --radius-sm:0.375rem; --radius-md:0.5rem;
  --radius-lg:0.75rem;  --radius-xl:1rem; --radius-full:9999px;

  --transition: 180ms cubic-bezier(0.16, 1, 0.3, 1);

  /* touch-safe min tap size */
  --tap: 44px;
}

/* ── Dark Theme ─────────────────────────────────────── */
:root, [data-theme="dark"] {
  --bg:            #0d0d0f;
  --surface:       #131316;
  --surface-2:     #18181c;
  --surface-3:     #1e1e24;
  --border:        rgba(255,255,255,0.07);
  --border-hi:     rgba(255,255,255,0.13);
  --text:          #e8e8ea;
  --text-muted:    #85858f;
  --text-faint:    #44444d;
  --primary:       #4f98a3;
  --primary-dim:   rgba(79,152,163,0.14);
  --success:       #6daa45;
  --success-dim:   rgba(109,170,69,0.13);
  --warning:       #e8af34;
  --warning-dim:   rgba(232,175,52,0.13);
  --error:         #dd6974;
  --error-dim:     rgba(221,105,116,0.13);
  --shadow-sm:     0 1px 3px rgba(0,0,0,0.35);
  --shadow-md:     0 4px 20px rgba(0,0,0,0.45);

  /* unhealthy section specific */
  --warn-section-bg:     rgba(232,175,52,0.06);
  --warn-section-border: rgba(232,175,52,0.18);
  --err-section-bg:      rgba(221,105,116,0.06);
  --err-section-border:  rgba(221,105,116,0.18);
}

/* ── Light Theme ────────────────────────────────────── */
[data-theme="light"] {
  --bg:            #f2f2ee;
  --surface:       #fafaf8;
  --surface-2:     #ffffff;
  --surface-3:     #ededea;
  --border:        rgba(0,0,0,0.07);
  --border-hi:     rgba(0,0,0,0.13);
  --text:          #1a1a1f;
  --text-muted:    #62626b;
  --text-faint:    #aaaaaf;
  --primary:       #01696f;
  --primary-dim:   rgba(1,105,111,0.1);
  --success:       #437a22;
  --success-dim:   rgba(67,122,34,0.1);
  --warning:       #9c6a00;
  --warning-dim:   rgba(156,106,0,0.1);
  --error:         #a13544;
  --error-dim:     rgba(161,53,68,0.1);
  --shadow-sm:     0 1px 3px rgba(0,0,0,0.06);
  --shadow-md:     0 4px 20px rgba(0,0,0,0.09);

  --warn-section-bg:     rgba(156,106,0,0.05);
  --warn-section-border: rgba(156,106,0,0.18);
  --err-section-bg:      rgba(161,53,68,0.05);
  --err-section-border:  rgba(161,53,68,0.18);
}

/* ── Base Reset ─────────────────────────────────────── */
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
html {
  -webkit-font-smoothing: antialiased;
  text-rendering: optimizeLegibility;
  scroll-behavior: smooth;
}
body {
  min-height: 100dvh;
  background: var(--bg);
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E");
  background-blend-mode: overlay;
  background-size: 128px;
  opacity: 1;
  color: var(--text);
  font-family: var(--font-body);
  font-size: var(--text-sm);
  line-height: 1.6;
}
button { cursor: pointer; background: none; border: none; font: inherit; color: inherit; }

/* ── Top Bar ─────────────────────────────────────────── */
.topbar {
  position: sticky; top: 0; z-index: 100;
  background: color-mix(in oklab, var(--surface) 92%, transparent);
  border-bottom: 1px solid var(--border);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  padding: 0 var(--space-4);
  height: 52px;
  display: flex; align-items: center; gap: var(--space-3);
}
.logo {
  display: flex; align-items: center; gap: var(--space-2);
  text-decoration: none; color: var(--text); flex-shrink: 0;
  min-height: var(--tap);
}
.logo-svg { color: var(--primary); flex-shrink: 0; }
.logo-name {
  font-size: var(--text-sm); font-weight: 600;
  letter-spacing: -0.01em; white-space: nowrap;
}
.logo-sub {
  font-size: var(--text-xs); color: var(--text-muted);
  font-weight: 400; margin-left: var(--space-1);
}
/* hide subtitle on very small screens */
@media (max-width: 380px) { .logo-sub { display: none; } }

.topbar-spacer { flex: 1; min-width: 0; }

.topbar-meta {
  display: flex; align-items: center; gap: var(--space-2);
  font-size: var(--text-xs); color: var(--text-muted);
  font-family: var(--font-mono); flex-shrink: 0;
}
.live-dot {
  width: 6px; height: 6px; border-radius: var(--radius-full);
  background: var(--success); display: inline-block; flex-shrink: 0;
  animation: blink 2.4s ease-in-out infinite;
}
@keyframes blink { 0%,100%{opacity:1} 50%{opacity:0.35} }

/* hide refresh interval on phones */
.topbar-refresh { display: none; }
@media (min-width: 480px) { .topbar-refresh { display: inline; } }
.topbar-sep { opacity: 0.4; }
@media (max-width: 479px) { .topbar-sep { display: none; } }

.icon-btn {
  width: var(--tap); height: var(--tap); border-radius: var(--radius-md);
  display: flex; align-items: center; justify-content: center;
  color: var(--text-muted); flex-shrink: 0;
  transition: background var(--transition), color var(--transition);
}
.icon-btn:hover { background: var(--surface-3); color: var(--text); }

/* ── Main ────────────────────────────────────────────── */
.main {
  padding: var(--space-5) var(--space-4) var(--space-8);
  max-width: 1400px; margin-inline: auto; width: 100%;
}
@media (min-width: 640px) {
  .main { padding: var(--space-6) var(--space-6) var(--space-8); }
}

/* ── Page Header ─────────────────────────────────────── */
.page-header { margin-bottom: var(--space-5); }
.page-title {
  font-size: var(--text-xl); font-weight: 600;
  letter-spacing: -0.025em; line-height: 1.1;
}
.page-sub {
  font-size: var(--text-xs); color: var(--text-muted);
  margin-top: var(--space-1); font-family: var(--font-mono);
  overflow-wrap: break-word;
}

/* ── KPI Row ─────────────────────────────────────────── */
.kpi-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-2);
  margin-bottom: var(--space-6);
}
@media (min-width: 480px) {
  .kpi-row { grid-template-columns: repeat(3, 1fr); gap: var(--space-3); }
}
@media (min-width: 900px) {
  .kpi-row { grid-template-columns: repeat(6, 1fr); }
}
.kpi {
  background: var(--surface); border: 1px solid var(--border);
  border-radius: var(--radius-lg); padding: var(--space-3) var(--space-3);
  display: flex; flex-direction: column; gap: 1px;
  transition: border-color var(--transition);
}
@media (min-width: 640px) {
  .kpi { padding: var(--space-4) var(--space-4); }
}
.kpi:hover { border-color: var(--border-hi); }
.kpi-label {
  font-size: 0.62rem; color: var(--text-muted);
  text-transform: uppercase; letter-spacing: 0.07em; font-weight: 500;
}
@media (min-width: 480px) { .kpi-label { font-size: var(--text-xs); } }
.kpi-val {
  font-size: var(--text-lg); font-weight: 700;
  letter-spacing: -0.025em; line-height: 1.15;
  font-variant-numeric: tabular-nums lining-nums;
}
.kpi-val.ok      { color: var(--success); }
.kpi-val.warn    { color: var(--warning); }
.kpi-val.bad     { color: var(--error); }
.kpi-val.neutral { color: var(--text); }
.kpi-hint {
  font-size: 0.62rem; color: var(--text-faint);
  font-family: var(--font-mono);
}
@media (min-width: 480px) { .kpi-hint { font-size: var(--text-xs); } }

/* ── Node Section ────────────────────────────────────── */
.node { margin-bottom: var(--space-8); }
.node:last-child { margin-bottom: 0; }
.node-head {
  display: flex; align-items: center; gap: var(--space-3);
  padding-bottom: var(--space-3); margin-bottom: var(--space-4);
  border-bottom: 1px solid var(--border);
}
.node-icon-wrap {
  width: 32px; height: 32px; border-radius: var(--radius-md);
  background: var(--primary-dim);
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0; align-self: center;
}
.node-prompt {
  font-family: var(--font-mono); font-size: 1.1rem; font-weight: 700;
  color: var(--primary); line-height: 1;
  animation: prompt-pulse 2s ease-in-out infinite;
  user-select: none; display: flex; align-items: center; justify-content: center;
}
@keyframes prompt-pulse {
  0%, 100% { opacity: 1; }
  50%       { opacity: 0.3; }
}
.node-label-row { display: flex; align-items: center; gap: var(--space-2); flex-wrap: wrap; }
.node-label { font-size: var(--text-base); font-weight: 600; letter-spacing: -0.01em; }
.node-meta  { display: flex; flex-direction: column; gap: 4px; min-width: 0; flex: 1; }
.node-location {
  display: inline-flex; align-items: center; gap: 5px;
  padding: 2px 8px; border-radius: var(--radius-full);
  background: color-mix(in oklab, var(--primary) 14%, transparent);
  border: 1px solid color-mix(in oklab, var(--primary) 24%, transparent);
  color: var(--primary); font-size: 0.62rem; font-weight: 700;
  letter-spacing: 0.04em; text-transform: uppercase;
}
.node-location::before {
  content: ""; width: 5px; height: 5px; border-radius: 50%;
  background: currentColor; opacity: 0.9;
}
.node-url-row {
  display: flex; align-items: baseline; gap: var(--space-3); min-width: 0;
}
.node-url {
  font-size: var(--text-xs); color: var(--text-faint);
  font-family: var(--font-mono); overflow-wrap: anywhere;
  flex: 1; min-width: 0;
}
.node-ts {
  font-size: var(--text-xs); color: var(--text-faint);
  font-family: var(--font-mono); white-space: nowrap; flex-shrink: 0;
}
.node-exporter {
  display: flex; align-items: center; gap: 4px; flex-wrap: wrap; margin-top: 3px;
}
.node-exporter-badge {
  font-size: 0.60rem; font-family: var(--font-mono); font-weight: 500;
  padding: 1px 6px; border-radius: var(--radius-sm);
  background: var(--surface-3); color: var(--text-muted);
  border: 1px solid var(--border);
  white-space: nowrap; line-height: 1.6;
}

/* ── Disk Section ────────────────────────────────────── */
.disk-section {
  margin-top: var(--space-4);
  border: 1px solid var(--border); border-radius: var(--radius-lg); overflow: hidden;
}
.disk-toggle {
  width: 100%; display: flex; align-items: center; justify-content: space-between;
  padding: var(--space-2) var(--space-4);
  cursor: pointer; transition: background var(--transition);
}
.disk-toggle:hover { background: var(--surface-3); }
.disk-toggle-label {
  display: flex; align-items: center; gap: var(--space-2);
  font-size: var(--text-xs); font-weight: 500; color: var(--text-muted);
}
.disk-chevron { transition: transform var(--transition); color: var(--text-faint); flex-shrink: 0; }
.disk-chevron.open { transform: rotate(180deg); }
.disk-list { display: none; }
.disk-list.open { display: block; }
.disk-row {
  display: flex; align-items: center; flex-wrap: wrap; gap: var(--space-2) var(--space-3);
  padding: var(--space-2) var(--space-4); border-top: 1px solid var(--border);
  font-size: var(--text-xs); font-family: var(--font-mono);
  transition: background var(--transition);
}
.disk-row:hover { background: var(--surface-2); }
.disk-model { flex: 1; min-width: 120px; color: var(--text); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.disk-serial { color: var(--text-faint); white-space: nowrap; }
.disk-cap { color: var(--text-muted); white-space: nowrap; }
.disk-type-badge {
  font-size: 0.58rem; font-weight: 700; letter-spacing: 0.05em; text-transform: uppercase;
  padding: 1px 5px; border-radius: var(--radius-sm); white-space: nowrap;
  background: var(--surface-3); color: var(--text-faint);
}
.disk-type-badge.nvme { background: color-mix(in oklab, var(--primary) 12%, transparent); color: var(--primary); }
.disk-type-badge.ssd  { background: color-mix(in oklab, var(--success) 10%, transparent); color: var(--success); }
.disk-temp { white-space: nowrap; font-weight: 600; }
.disk-temp.cool { color: var(--success); }
.disk-temp.warm { color: var(--warning); }
.disk-temp.hot  { color: var(--error); }
.disk-smart {
  font-size: 0.58rem; font-weight: 700; letter-spacing: 0.04em; text-transform: uppercase;
  padding: 1px 5px; border-radius: var(--radius-sm); white-space: nowrap;
}
.disk-smart.pass { background: var(--success-dim); color: var(--success); }
.disk-smart.fail { background: var(--error-dim);   color: var(--error); }
.disk-hours { color: var(--text-faint); white-space: nowrap; }

/* ── Error Banner ────────────────────────────────────── */
.node-err {
  background: var(--error-dim); border: 1px solid color-mix(in oklab, var(--error) 30%, transparent);
  border-radius: var(--radius-lg); padding: var(--space-4) var(--space-4);
  color: var(--error); font-size: var(--text-sm);
  display: flex; align-items: flex-start; gap: var(--space-3);
}
.node-err svg { flex-shrink: 0; margin-top: 1px; }

/* ── Empty State ─────────────────────────────────────── */
.empty {
  text-align: center; padding: var(--space-10) var(--space-6);
  color: var(--text-muted); display: flex; flex-direction: column;
  align-items: center; gap: var(--space-3);
}
.empty svg { color: var(--text-faint); }
.empty h3  { color: var(--text); font-size: var(--text-base); font-weight: 600; }
.empty p   { font-size: var(--text-sm); max-width: 36ch; }

/* ═══════════════════════════════════════════════════════
   ── Unhealthy Pools Section ──────────────────────────
   ═══════════════════════════════════════════════════════ */

/* Shared collapsible wrapper */
.alert-section {
  margin-bottom: var(--space-5);
  border-radius: var(--radius-xl);
  overflow: hidden;
  border: 1px solid var(--border);
}

/* Degraded variant */
.alert-section.degraded {
  border-color: var(--warn-section-border);
  background: var(--warn-section-bg);
}

/* Faulted/unreachable variant */
.alert-section.faulted {
  border-color: var(--err-section-border);
  background: var(--err-section-bg);
}

/* Toggle header */
.alert-toggle {
  width: 100%;
  display: flex; align-items: center; justify-content: space-between; gap: var(--space-3);
  padding: var(--space-4) var(--space-4);
  cursor: pointer;
  min-height: var(--tap);
  /* background is set on .alert-section so toggle stays transparent */
  background: none; border: none;
}
@media (min-width: 480px) { .alert-toggle { padding: var(--space-4) var(--space-5); } }

.alert-toggle:focus-visible {
  outline: 2px solid var(--primary); outline-offset: -2px;
}

.alert-toggle-left { display: flex; align-items: center; gap: var(--space-3); min-width: 0; }

/* Pulsing dot */
.alert-dot {
  width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0;
  animation: blink-dot 1.6s ease-in-out infinite;
}
.degraded .alert-dot { background: var(--warning); }
.faulted  .alert-dot { background: var(--error); }
@keyframes blink-dot { 0%,100%{opacity:1} 50%{opacity:0.25} }

.alert-toggle-text { display: flex; align-items: center; gap: var(--space-2); flex-wrap: wrap; min-width: 0; }
.alert-toggle-title {
  font-size: var(--text-sm); font-weight: 600; white-space: nowrap;
}
.degraded .alert-toggle-title { color: var(--warning); }
.faulted  .alert-toggle-title { color: var(--error); }

.alert-toggle-badge {
  display: inline-flex; align-items: center; justify-content: center;
  min-width: 20px; height: 20px;
  padding: 0 6px; border-radius: var(--radius-full);
  font-size: 0.65rem; font-weight: 700; letter-spacing: 0.04em; font-family: var(--font-mono);
}
.degraded .alert-toggle-badge {
  background: rgba(232,175,52,0.2); color: var(--warning);
  border: 1px solid rgba(232,175,52,0.35);
}
.faulted .alert-toggle-badge {
  background: rgba(221,105,116,0.2); color: var(--error);
  border: 1px solid rgba(221,105,116,0.35);
}

.alert-toggle-desc {
  font-size: var(--text-xs); font-family: var(--font-mono); opacity: 0.6;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.degraded .alert-toggle-desc { color: var(--warning); }
.faulted  .alert-toggle-desc { color: var(--error); }

.alert-arrow {
  flex-shrink: 0; transition: transform 220ms ease;
}
.alert-arrow.open { transform: rotate(180deg); }
.degraded .alert-arrow { color: var(--warning); }
.faulted  .alert-arrow { color: var(--error); }

/* Divider between toggle and body */
.alert-divider {
  height: 1px; margin: 0;
  border: none;
}
.degraded .alert-divider { background: var(--warn-section-border); }
.faulted  .alert-divider  { background: var(--err-section-border); }
.alert-divider.hidden { display: none; }

/* Body */
.alert-body {
  padding: var(--space-4);
  display: grid; gap: var(--space-3);
  /* subtle inner bg to lift cards slightly */
}
.degraded .alert-body { background: rgba(232,175,52,0.03); }
.faulted  .alert-body  { background: rgba(221,105,116,0.03); }
@media (min-width: 480px) { .alert-body { padding: var(--space-4) var(--space-5) var(--space-5); } }
.alert-body.collapsed { display: none; }

/* pool cards inside an alert section get a themed border */
.alert-body .pool-card {
  border-color: var(--warn-section-border);
}
.faulted .alert-body .pool-card {
  border-color: var(--err-section-border);
}
.alert-body .pool-card:hover {
  border-color: color-mix(in oklab, var(--warning) 50%, transparent);
}
.faulted .alert-body .pool-card:hover {
  border-color: color-mix(in oklab, var(--error) 50%, transparent);
}

/* Node failures (unreachable) inside the faulted alert section */
.alert-node-err {
  background: color-mix(in oklab, var(--error-dim) 60%, transparent);
  border: 1px solid color-mix(in oklab, var(--error) 22%, transparent);
  border-radius: var(--radius-lg); padding: var(--space-4);
  display: flex; flex-direction: column; gap: var(--space-3);
}
.alert-node-err-head {
  display: flex; align-items: center; gap: var(--space-2);
}
.alert-node-err-head > svg { margin-top: 2px; }
.alert-node-err-label { font-size: var(--text-sm); font-weight: 600; color: var(--error); line-height: 1; }
.alert-node-err-loc {
  font-size: 0.62rem; font-weight: 700; letter-spacing: 0.04em; text-transform: uppercase;
  padding: 2px 7px; border-radius: var(--radius-full);
  background: rgba(221,105,116,0.15); color: var(--error);
  border: 1px solid rgba(221,105,116,0.3);
}
.alert-node-err-url {
  font-size: var(--text-xs); color: var(--text-faint); font-family: var(--font-mono);
  overflow-wrap: anywhere;
}
.alert-node-err-msg {
  display: flex; align-items: flex-start; gap: var(--space-2);
  font-size: var(--text-xs); color: var(--error); font-family: var(--font-mono);
  opacity: 0.85;
}
.alert-node-err-msg svg { flex-shrink: 0; margin-top: 1px; }

/* Pool detail Modal ───────────────────────────────────── */
.pool-card[data-has-datasets="true"] { cursor: pointer; }
.pool-detail-overlay {
  position: fixed; inset: 0; z-index: 1200;
  display: none; align-items: flex-end; justify-content: center;
  padding: 0;
  background: rgba(5, 6, 10, 0.72);
  backdrop-filter: blur(10px); -webkit-backdrop-filter: blur(10px);
}
@media (min-width: 640px) {
  .pool-detail-overlay {
    align-items: center;
    padding: var(--space-4);
  }
}
.pool-detail-overlay.open { display: flex; }

.pool-detail-modal {
  width: 100%;
  max-height: 92dvh;
  overflow: auto;
  background: color-mix(in oklab, var(--surface) 96%, black 4%);
  border: 1px solid var(--border-hi);
  border-radius: var(--radius-xl) var(--radius-xl) 0 0;
  box-shadow: 0 -8px 48px rgba(0,0,0,0.4);
  overscroll-behavior: contain;
}
@media (min-width: 640px) {
  .pool-detail-modal {
    width: min(980px, 100%);
    max-height: min(88dvh, 900px);
    border-radius: calc(var(--radius-xl) + 2px);
    box-shadow: 0 24px 80px rgba(0,0,0,0.45);
  }
}

/* drag handle hint on mobile */
.pool-detail-modal::before {
  content: "";
  display: block;
  width: 36px; height: 4px;
  border-radius: var(--radius-full);
  background: var(--border-hi);
  margin: var(--space-3) auto var(--space-1);
}
@media (min-width: 640px) { .pool-detail-modal::before { display: none; } }

.pool-detail-head {
  position: sticky; top: 0; z-index: 2;
  display: flex; align-items: flex-start; justify-content: space-between; gap: var(--space-3);
  padding: var(--space-4);
  background: color-mix(in oklab, var(--surface) 90%, transparent);
  border-bottom: 1px solid var(--border);
  backdrop-filter: blur(8px);
}
@media (min-width: 480px) { .pool-detail-head { padding: var(--space-5); } }

.pool-detail-title { font-size: var(--text-lg); font-weight: 600; letter-spacing: -0.02em; }
.pool-detail-sub { margin-top: 3px; font-size: var(--text-xs); color: var(--text-muted); font-family: var(--font-mono); }

.pool-detail-close {
  width: var(--tap); height: var(--tap); border-radius: var(--radius-md);
  display: flex; align-items: center; justify-content: center;
  color: var(--text-muted); border: 1px solid var(--border);
  flex-shrink: 0;
}
.pool-detail-close:hover { color: var(--text); background: var(--surface-3); }

.pool-detail-body { padding: var(--space-4); }
@media (min-width: 480px) { .pool-detail-body { padding: var(--space-5); } }

.pool-detail-empty { color: var(--text-muted); font-size: var(--text-sm); padding: var(--space-6) 0; }

.dataset-list { display: grid; gap: var(--space-3); }
.dataset-item {
  background: var(--surface-2); border: 1px solid var(--border);
  border-radius: var(--radius-lg); padding: var(--space-4);
}
.dataset-item-top {
  display: flex; align-items: center; justify-content: space-between;
  gap: var(--space-3); margin-bottom: var(--space-3);
}
.dataset-name {
  font-family: var(--font-mono); font-size: var(--text-xs); color: var(--text);
  word-break: break-all;
}
@media (min-width: 480px) { .dataset-name { font-size: var(--text-sm); } }
.dataset-kind {
  padding: 2px 7px; border-radius: var(--radius-full);
  background: var(--surface-3); color: var(--text-muted);
  font-size: 0.62rem; font-weight: 700; letter-spacing: 0.05em; text-transform: uppercase;
  flex-shrink: 0;
}
.dataset-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--space-2);
}
@media (min-width: 480px) { .dataset-grid { grid-template-columns: repeat(3, minmax(0, 1fr)); gap: var(--space-3); } }
.dataset-metric { display: flex; flex-direction: column; gap: 2px; }
.dataset-metric-label {
  font-size: 0.62rem; color: var(--text-faint); text-transform: uppercase; letter-spacing: 0.05em;
}
.dataset-metric-value { font-size: var(--text-xs); color: var(--text-muted); font-family: var(--font-mono); }
.dataset-metric-value.hot { color: var(--warning); }

/* ── Pool Grid ───────────────────────────────────────── */
.pools-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--space-3);
}
@media (min-width: 560px) {
  .pools-grid { grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: var(--space-4); }
}

/* ── Pool Card ───────────────────────────────────────── */
.pool-card {
  background: var(--surface); border: 1px solid var(--border);
  border-radius: var(--radius-xl); padding: var(--space-4);
  display: flex; flex-direction: column; gap: var(--space-3);
  transition: box-shadow var(--transition), border-color var(--transition), transform var(--transition);
  will-change: transform;
}
@media (min-width: 480px) { .pool-card { padding: var(--space-5); gap: var(--space-4); } }
@media (hover: hover) {
  .pool-card:hover {
    border-color: var(--border-hi);
    box-shadow: var(--shadow-md);
    transform: translateY(-2px);
  }
}
/* tap feedback on touch devices */
@media (hover: none) {
  .pool-card:active { opacity: 0.85; transform: scale(0.99); }
}

/* Pool card top row */
.pool-top {
  display: flex; align-items: flex-start;
  justify-content: space-between; gap: var(--space-3);
}
.pool-name {
  font-size: var(--text-base); font-weight: 600;
  font-family: var(--font-mono); letter-spacing: -0.01em; word-break: break-all;
}
.pool-sub {
  font-size: var(--text-xs); color: var(--text-faint);
  font-family: var(--font-mono); margin-top: 2px;
}

/* ── Health Badge ────────────────────────────────────── */
.hbadge {
  display: inline-flex; align-items: center; gap: 5px;
  padding: 3px 8px; border-radius: var(--radius-full);
  font-size: 0.65rem; font-weight: 700;
  letter-spacing: 0.06em; text-transform: uppercase; flex-shrink: 0;
}
.hbadge-dot { width: 5px; height: 5px; border-radius: 50%; flex-shrink: 0; }

.health-online   { background: var(--success-dim); color: var(--success); }
.health-online   .hbadge-dot { background: var(--success); }
.health-degraded { background: var(--warning-dim); color: var(--warning); }
.health-degraded .hbadge-dot { background: var(--warning); }
.health-faulted  { background: var(--error-dim);   color: var(--error); }
.health-faulted  .hbadge-dot { background: var(--error); }

/* ── Usage Bar ───────────────────────────────────────── */
.usage { display: flex; flex-direction: column; gap: var(--space-2); }
.usage-header {
  display: flex; justify-content: space-between; align-items: baseline;
  font-size: var(--text-xs); color: var(--text-muted);
}
.usage-pct {
  font-family: var(--font-mono); font-weight: 600;
  font-variant-numeric: tabular-nums; color: var(--text);
}
.bar-track {
  height: 5px; background: var(--surface-3);
  border-radius: var(--radius-full); overflow: hidden;
}
.bar-fill {
  height: 100%; border-radius: var(--radius-full); background: var(--primary);
  transition: width 0.65s cubic-bezier(0.16,1,0.3,1);
}
.bar-fill.warn  { background: var(--warning); }
.bar-fill.bad   { background: var(--error); }
.usage-sizes {
  display: flex; justify-content: space-between; flex-wrap: wrap;
  font-size: var(--text-xs); font-family: var(--font-mono);
  color: var(--text-faint); gap: 4px var(--space-2);
}
.usage-sizes span { white-space: nowrap; }
.usage-sizes b { color: var(--text-muted); font-weight: 500; }

/* ── Stats 2×2 ───────────────────────────────────────── */
.divider { border: none; border-top: 1px solid var(--border); }

.stats {
  display: grid; grid-template-columns: 1fr 1fr; gap: var(--space-2);
}
.stat {
  background: var(--surface-3); border-radius: var(--radius-md);
  padding: var(--space-3); display: flex; flex-direction: column; gap: 2px;
}
.stat-lbl {
  font-size: 0.62rem; color: var(--text-faint);
  text-transform: uppercase; letter-spacing: 0.06em; font-weight: 500;
}
.stat-val {
  font-size: var(--text-sm); font-family: var(--font-mono);
  font-weight: 500; color: var(--text);
  font-variant-numeric: tabular-nums lining-nums;
}

/* ── Error Chips ─────────────────────────────────────── */
.err-row { display: flex; gap: var(--space-2); flex-wrap: wrap; }
.chip {
  font-size: 0.65rem; font-family: var(--font-mono); font-weight: 500;
  padding: 3px var(--space-2); border-radius: var(--radius-sm);
  background: var(--surface-3); color: var(--text-faint);
  white-space: nowrap;
}
.chip.hot { background: var(--error-dim); color: var(--error); }

/* ── Footer ──────────────────────────────────────────── */
.site-footer {
  border-top: 1px solid var(--border);
  padding: var(--space-4) var(--space-4);
  display: flex; align-items: center; justify-content: center;
}
.footer-link {
  display: inline-flex; align-items: center; gap: var(--space-2);
  font-size: var(--text-xs); color: var(--text-faint);
  text-decoration: none; font-family: var(--font-mono);
  transition: color var(--transition);
}
.footer-link:hover { color: var(--text-muted); }

/* ── Refresh Progress Bar ────────────────────────────── */
#rbar {
  position: fixed; bottom: 0; left: 0; height: 2px;
  background: var(--primary); width: 0%; z-index: 999;
  opacity: 0.55; pointer-events: none;
}

/* ── Entrance Animations ─────────────────────────────── */
@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(12px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
@keyframes scaleIn {
  from { opacity: 0; transform: scale(0.96); }
  to { opacity: 1; transform: scale(1); }
}
.animate-in {
  animation: fadeInUp 0.55s cubic-bezier(0.22, 1, 0.36, 1) forwards;
  opacity: 0;
  transform-origin: center top;
}
.animate-in-delayed-1 { animation-delay: 0.08s; }
.animate-in-delayed-2 { animation-delay: 0.16s; }
.animate-in-delayed-3 { animation-delay: 0.24s; }
.animate-in-delayed-4 { animation-delay: 0.32s; }
.animate-in-delayed-5 { animation-delay: 0.4s; }
.animate-in-delayed-6 { animation-delay: 0.48s; }

.animate-in-modal {
  animation: scaleIn 0.3s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  opacity: 0;
}

/* Focus trap styles */
.pool-detail-modal:focus { outline: none; }
.pool-detail-modal:focus-visible {
  outline: 2px solid var(--primary); outline-offset: 2px;
}

/* Enhanced focus states */
button:focus-visible,
.pool-card:focus-visible,
.icon-btn:focus-visible,
.pool-detail-close:focus-visible {
  outline: 2px solid var(--primary); outline-offset: 2px;
}
.pool-card[data-has-datasets="true"]:focus-visible {
  outline-offset: 3px;
}

@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
  .animate-in, .animate-in-modal { animation: none; opacity: 1; }
}
</style>
</head>
<body>

<!-- ── Top Bar ────────────────────────────────────────── -->
<header class="topbar">
  <a class="logo" href="/" aria-label="ZFS Dash home">
    <svg class="logo-svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
         stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"
         aria-hidden="true">
      <rect x="2" y="2"  width="20" height="6" rx="1.5"/>
      <rect x="2" y="10" width="20" height="6" rx="1.5"/>
      <rect x="2" y="18" width="20" height="4" rx="1.5"/>
      <circle cx="5.5" cy="5"  r="0.9" fill="currentColor" stroke="none"/>
      <circle cx="5.5" cy="13" r="0.9" fill="currentColor" stroke="none"/>
      <circle cx="5.5" cy="20" r="0.9" fill="currentColor" stroke="none"/>
    </svg>
    <span class="logo-name">zfs-dash<span class="logo-sub">pool monitor</span></span>
  </a>

  <div class="topbar-spacer"></div>

  <div class="topbar-meta" aria-live="polite">
    <span class="live-dot" aria-hidden="true"></span>
    <span id="ts">{{.FetchedAt}}</span>
    <span class="topbar-sep" aria-hidden="true">·</span>
    <span class="topbar-refresh">↻&thinsp;{{.RefreshSecs}}s</span>
  </div>

  <button class="icon-btn" data-theme-toggle aria-label="Toggle light/dark mode">
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none"
         stroke="currentColor" stroke-width="2" aria-hidden="true">
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
    </svg>
  </button>
</header>

<!-- ── Main ───────────────────────────────────────────── -->
<main class="main">

  <!-- Page header -->
  <div class="page-header animate-in">
    <h1 class="page-title">Pool Overview</h1>
    <p class="page-sub">{{.TotalNodes}} node{{if gt .TotalNodes 1}}s{{end}} &middot; {{.TotalPools}} pool{{if gt .TotalPools 1}}s{{end}} &middot; {{.UnreachableNodes}} unreachable &middot; fetched {{.FetchedAt}}</p>
  </div>

  <!-- KPI summary -->
  <div class="kpi-row animate-in animate-in-delayed-1" role="region" aria-label="Fleet summary">
    <div class="kpi">
      <div class="kpi-label">Nodes</div>
      <div class="kpi-val neutral">{{.TotalNodes}}</div>
      <div class="kpi-hint">endpoints</div>
    </div>
    <div class="kpi">
      <div class="kpi-label">Pools</div>
      <div class="kpi-val neutral">{{.TotalPools}}</div>
      <div class="kpi-hint">total</div>
    </div>
    <div class="kpi">
      <div class="kpi-label">Unreachable</div>
      <div class="kpi-val {{if gt .UnreachableNodes 0}}bad{{else}}neutral{{end}}">{{.UnreachableNodes}}</div>
      <div class="kpi-hint">nodes down</div>
    </div>
    <div class="kpi">
      <div class="kpi-label">Healthy</div>
      <div class="kpi-val ok">{{.HealthyPools}}</div>
      <div class="kpi-hint">ONLINE</div>
    </div>
    <div class="kpi">
      <div class="kpi-label">Degraded</div>
      <div class="kpi-val {{if gt .DegradedPools 0}}warn{{else}}neutral{{end}}">{{.DegradedPools}}</div>
      <div class="kpi-hint">DEGRADED</div>
    </div>
    <div class="kpi">
      <div class="kpi-label">Faulted</div>
      <div class="kpi-val {{if gt .ErroredPools 0}}bad{{else}}neutral{{end}}">{{.ErroredPools}}</div>
      <div class="kpi-hint">FAULTED/UNAVAIL</div>
    </div>
  </div>

  <!-- ══ Unreachable Nodes ══ -->
  {{if gt .UnreachableNodes 0}}
  <div class="alert-section faulted animate-in animate-in-delayed-2" id="section-unreachable">
    <button class="alert-toggle" id="toggle-unreachable"
            aria-expanded="false" aria-controls="body-unreachable">
      <div class="alert-toggle-left">
        <div class="alert-dot" aria-hidden="true"></div>
        <div class="alert-toggle-text">
          <span class="alert-toggle-title">Unreachable Nodes</span>
          <span class="alert-toggle-badge">{{.UnreachableNodes}}</span>
          <span class="alert-toggle-desc">cannot be reached</span>
        </div>
      </div>
      <svg class="alert-arrow" id="arrow-unreachable"
           width="16" height="16" viewBox="0 0 24 24" fill="none"
           stroke="currentColor" stroke-width="2" aria-hidden="true">
        <path d="m6 9 6 6 6-6"/>
      </svg>
    </button>
    <hr class="alert-divider hidden" id="divider-unreachable">
    <div class="alert-body collapsed" id="body-unreachable">
      {{range $node := .Nodes}}
      {{if $node.Error}}
      <div class="alert-node-err">
        <div class="alert-node-err-head">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="var(--error)"
               stroke-width="2" stroke-linecap="round" flex-shrink="0" aria-hidden="true">
            <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          <span class="alert-node-err-label">{{$node.Label}}</span>
          {{if $node.Location}}<span class="alert-node-err-loc">{{$node.Location}}</span>{{end}}
        </div>
        <div class="alert-node-err-url">{{$node.URL}}</div>
        <div class="alert-node-err-msg">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor"
               stroke-width="2" aria-hidden="true">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
          <span>{{$node.Error}}</span>
        </div>
      </div>
      {{end}}
      {{end}}
    </div>
  </div>
  {{end}}

  <!-- ══ Degraded Pools ══ -->
  {{if gt .DegradedPools 0}}
  <div class="alert-section degraded animate-in animate-in-delayed-2" id="section-degraded">
    <button class="alert-toggle" id="toggle-degraded"
            aria-expanded="false" aria-controls="body-degraded">
      <div class="alert-toggle-left">
        <div class="alert-dot" aria-hidden="true"></div>
        <div class="alert-toggle-text">
          <span class="alert-toggle-title">Degraded Pools</span>
          <span class="alert-toggle-badge">{{.DegradedPools}}</span>
          <span class="alert-toggle-desc">require attention</span>
        </div>
      </div>
      <svg class="alert-arrow" id="arrow-degraded"
           width="16" height="16" viewBox="0 0 24 24" fill="none"
           stroke="currentColor" stroke-width="2" aria-hidden="true">
        <path d="m6 9 6 6 6-6"/>
      </svg>
    </button>
    <hr class="alert-divider hidden" id="divider-degraded">
    <div class="alert-body collapsed" id="body-degraded">
      <!-- pools grid reused inside the alert body -->
      <div class="pools-grid">
        {{range $ni, $node := .Nodes}}
        {{if not $node.Error}}
        {{range $pi, $pool := $node.Pools}}
        {{if eq $pool.Health "DEGRADED"}}
        {{template "poolCard" dict "pool" $pool "ni" $ni "pi" $pi}}
        {{end}}
        {{end}}
        {{end}}
        {{end}}
      </div>
    </div>
  </div>
  {{end}}

  <!-- ══ Faulted / Unavailable Pools ══ -->
  {{if gt .ErroredPools 0}}
  <div class="alert-section faulted animate-in animate-in-delayed-2" id="section-faulted">
    <button class="alert-toggle" id="toggle-faulted"
            aria-expanded="false" aria-controls="body-faulted">
      <div class="alert-toggle-left">
        <div class="alert-dot" aria-hidden="true"></div>
        <div class="alert-toggle-text">
          <span class="alert-toggle-title">Faulted / Unavailable Pools</span>
          <span class="alert-toggle-badge">{{.ErroredPools}}</span>
          <span class="alert-toggle-desc">immediate action needed</span>
        </div>
      </div>
      <svg class="alert-arrow" id="arrow-faulted"
           width="16" height="16" viewBox="0 0 24 24" fill="none"
           stroke="currentColor" stroke-width="2" aria-hidden="true">
        <path d="m6 9 6 6 6-6"/>
      </svg>
    </button>
    <hr class="alert-divider hidden" id="divider-faulted">
    <div class="alert-body collapsed" id="body-faulted">
      <div class="pools-grid">
        {{range $ni, $node := .Nodes}}
        {{if not $node.Error}}
        {{range $pi, $pool := $node.Pools}}
        {{if ne $pool.Health "ONLINE"}}{{if ne $pool.Health "DEGRADED"}}
        {{template "poolCard" dict "pool" $pool "ni" $ni "pi" $pi}}
        {{end}}{{end}}
        {{end}}
        {{end}}
        {{end}}
      </div>
    </div>
  </div>
  {{end}}

  <!-- ══ Per-node healthy pools ══ -->
  {{range $ni, $node := .Nodes}}
  {{if not $node.Error}}
  <section class="node animate-in {{if eq $ni 0}}animate-in-delayed-3{{else if eq $ni 1}}animate-in-delayed-4{{else if eq $ni 2}}animate-in-delayed-5{{else}}animate-in-delayed-6{{end}}" aria-label="Node {{$node.Label}}">

    <div class="node-head">
      <div class="node-icon-wrap" aria-hidden="true">
        <span class="node-prompt">›</span>
      </div>
      <div class="node-meta">
        <div class="node-label-row">
          <div class="node-label">{{.Label}}</div>
          {{if $node.Location}}<div class="node-location">{{$node.Location}}</div>{{end}}
        </div>
        <div class="node-url-row">
          <div class="node-url">{{$node.URL}}</div>
          <div class="node-ts" aria-label="Fetched at">{{fmtNodeTime $node.FetchedAt}}</div>
        </div>
        {{if $node.ExporterInfo.Version}}<div class="node-exporter">
          <span class="node-exporter-badge">zfs-exporter {{$node.ExporterInfo.Version}}</span>
          {{if $node.ExporterInfo.GoVersion}}<span class="node-exporter-badge">{{$node.ExporterInfo.GoVersion}}</span>{{end}}
        </div>{{end}}
      </div>
    </div>

    {{if not $node.Pools}}
    <div class="empty">
      <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor"
           stroke-width="1.4" aria-hidden="true">
        <rect x="2" y="3" width="20" height="18" rx="2"/>
        <path d="M8 7h8M8 12h8M8 17h4"/>
      </svg>
      <h3>No pools found</h3>
      <p>The exporter returned no ZFS pool metrics for this node.</p>
    </div>
    {{else}}
    <div class="pools-grid">
      {{range $pi, $pool := $node.Pools}}
      {{template "poolCard" dict "pool" $pool "ni" $ni "pi" $pi}}
      {{end}}
    </div>
    {{end}}

    {{if $node.Disks}}
    <div class="disk-section">
      <button class="disk-toggle" id="disk-toggle-{{$ni}}"
              aria-expanded="false" aria-controls="disk-list-{{$ni}}">
        <div class="disk-toggle-label">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <ellipse cx="12" cy="5" rx="9" ry="3"/>
            <path d="M3 5v14c0 1.66 4.03 3 9 3s9-1.34 9-3V5"/>
            <path d="M3 12c0 1.66 4.03 3 9 3s9-1.34 9-3"/>
          </svg>
          Disks ({{len $node.Disks}})
        </div>
        <svg class="disk-chevron" id="disk-chevron-{{$ni}}" width="14" height="14"
             viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
          <path d="m6 9 6 6 6-6"/>
        </svg>
      </button>
      <div class="disk-list" id="disk-list-{{$ni}}">
        {{range $disk := $node.Disks}}
        <div class="disk-row">
          <span class="disk-model">{{if $disk.ModelName}}{{$disk.ModelName}}{{else}}{{maskSerial $disk.Device}}{{end}}</span>
          {{if $disk.SerialNumber}}<span class="disk-serial">{{maskSerial $disk.SerialNumber}}</span>{{end}}
          {{if gt0 $disk.CapacityBytes}}<span class="disk-cap">{{humanBytes $disk.CapacityBytes}}</span>{{end}}
          <span class="disk-type-badge {{diskTypeClass $disk.Interface $disk.RotationRate}}">{{diskTypeLabel $disk.Interface $disk.RotationRate}}</span>
          {{if gt0 $disk.Temperature}}<span class="disk-temp {{tempClass $disk.Temperature}}">{{printf "%.0f" $disk.Temperature}}°C</span>{{end}}
          <span class="disk-smart {{if $disk.SmartPassed}}pass{{else}}fail{{end}}">{{if $disk.SmartPassed}}PASS{{else}}FAIL{{end}}</span>
          {{if gt0 $disk.PowerOnHours}}<span class="disk-hours">{{printf "%.0f" $disk.PowerOnHours}}h</span>{{end}}
        </div>
        {{end}}
      </div>
    </div>
    {{end}}

  </section>
  {{end}}
  {{end}}

</main>

<!-- ── Pool Detail Modal ──────────────────────────────── -->
<div class="pool-detail-overlay" id="pool-detail-overlay" hidden>
  <div class="pool-detail-modal" role="dialog" aria-modal="true" aria-labelledby="pool-detail-title">
    <div class="pool-detail-head">
      <div>
        <div class="pool-detail-title" id="pool-detail-title">Pool Details</div>
        <div class="pool-detail-sub" id="pool-detail-sub"></div>
      </div>
      <button class="pool-detail-close" id="pool-detail-close" aria-label="Close pool details">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
             stroke-width="2" aria-hidden="true">
          <path d="M18 6 6 18M6 6l12 12"/>
        </svg>
      </button>
    </div>
    <div class="pool-detail-body" id="pool-detail-body"></div>
  </div>
</div>

<!-- footer -->
<footer class="site-footer">
  <a class="footer-link" href="https://github.com/crazyuploader/zfs-dash" target="_blank" rel="noopener noreferrer">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
      <path d="M12 2C6.477 2 2 6.484 2 12.021c0 4.428 2.865 8.184 6.839 9.504.5.092.682-.217.682-.482 0-.237-.009-.868-.013-1.703-2.782.605-3.369-1.342-3.369-1.342-.454-1.156-1.11-1.463-1.11-1.463-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.026 2.747-1.026.546 1.378.202 2.397.1 2.65.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482C19.138 20.2 22 16.447 22 12.021 22 6.484 17.522 2 12 2z"/>
    </svg>
    crazyuploader/zfs-dash
  </a>
</footer>

<!-- auto-refresh progress bar -->
<div id="rbar" aria-hidden="true"></div>

{{/* ── Pool Card Sub-template ── */}}
{{define "poolCard"}}
<article class="pool-card"
         aria-label="Pool {{.pool.Name}}"
         data-node-index="{{.ni}}"
         data-pool-index="{{.pi}}"
         data-has-datasets="{{if .pool.Datasets}}true{{else}}false{{end}}">

  <div class="pool-top">
    <div>
      <div class="pool-name">{{.pool.Name}}</div>
      <div class="pool-sub">zfs pool</div>
    </div>
    <span class="hbadge {{healthClass .pool.Health}}">
      <span class="hbadge-dot" aria-hidden="true"></span>
      {{.pool.Health}}
    </span>
  </div>

  {{if gt .pool.Size 0.0}}
  <div class="usage">
    <div class="usage-header">
      <span>Capacity used</span>
      <span class="usage-pct">{{printf "%.1f" .pool.UsedPercent}}%</span>
    </div>
    <div class="bar-track"
         role="progressbar"
         aria-valuenow="{{printf "%.0f" .pool.UsedPercent}}"
         aria-valuemin="0" aria-valuemax="100"
         aria-label="{{printf "%.1f" .pool.UsedPercent}}% used">
      <div class="bar-fill{{if gte .pool.UsedPercent 90.0}} bad{{else if gte .pool.UsedPercent 75.0}} warn{{end}}"
           data-width="{{printf "%.2f" .pool.UsedPercent}}"
           style="width:0%"></div>
    </div>
    <div class="usage-sizes">
      <span>Used&nbsp;<b>{{humanBytes .pool.Allocated}}</b></span>
      <span>Free&nbsp;<b>{{humanBytes .pool.Free}}</b></span>
      <span>Total&nbsp;<b>{{humanBytes .pool.Size}}</b></span>
    </div>
  </div>
  {{end}}

  <hr class="divider" aria-hidden="true">

  <div class="stats">
    <div class="stat">
      <span class="stat-lbl">Dedup</span>
      <span class="stat-val">{{printf "%.2fx" .pool.DedupRatio}}</span>
    </div>
    <div class="stat">
      <span class="stat-lbl">Fragmentation</span>
      <span class="stat-val">{{printf "%.0f%%" (mul100 .pool.FragmentationRatio)}}</span>
    </div>
    <div class="stat">
      <span class="stat-lbl">Freeing</span>
      <span class="stat-val">{{humanBytes .pool.Freeing}}</span>
    </div>
    <div class="stat">
      <span class="stat-lbl">Leaked</span>
      <span class="stat-val">{{humanBytes .pool.LeakedBytes}}</span>
    </div>
  </div>

  <div class="err-row">
    <span class="chip{{if .pool.ReadOnly}} hot{{end}}">Readonly&nbsp;{{if .pool.ReadOnly}}yes{{else}}no{{end}}</span>
    <span class="chip{{if gt0 .pool.Freeing}} hot{{end}}">Freeing&nbsp;{{humanBytes .pool.Freeing}}</span>
    <span class="chip{{if gt0 .pool.LeakedBytes}} hot{{end}}">Leaked&nbsp;{{humanBytes .pool.LeakedBytes}}</span>
  </div>

</article>
{{end}}

<script>
(function () {
  'use strict';

  const nodes = {{.NodesJSON}};
  const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;

  /* ── Theme toggle ─────────────────────────────────── */
  const html    = document.documentElement;
  const themeBtn = document.querySelector('[data-theme-toggle]');
  const themeKey = 'zfs-dash-theme';
  let theme;
  try { theme = localStorage.getItem(themeKey); } catch (_) { theme = null; }
  if (theme !== 'dark' && theme !== 'light') {
    theme = window.matchMedia('(prefers-color-scheme:dark)').matches ? 'dark' : 'light';
  }

  const moonSVG = '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>';
  const sunSVG  = '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>';

  function applyTheme(t) {
    html.setAttribute('data-theme', t);
    if (themeBtn) {
      themeBtn.innerHTML = t === 'dark' ? moonSVG : sunSVG;
      themeBtn.setAttribute('aria-label', 'Switch to ' + (t === 'dark' ? 'light' : 'dark') + ' mode');
    }
  }
  applyTheme(theme);
  if (themeBtn) {
    themeBtn.addEventListener('click', function () {
      theme = theme === 'dark' ? 'light' : 'dark';
      try { localStorage.setItem(themeKey, theme); } catch (_) {}
      applyTheme(theme);
    });
  }

  /* ── Alert section toggles ────────────────────────── */
  ['unreachable', 'degraded', 'faulted'].forEach(function (id) {
    const toggle  = document.getElementById('toggle-' + id);
    const body    = document.getElementById('body-' + id);
    const arrow   = document.getElementById('arrow-' + id);
    const divider = document.getElementById('divider-' + id);
    if (!toggle || !body) return;

    toggle.addEventListener('click', function () {
      const isOpen = toggle.getAttribute('aria-expanded') === 'true';
      const nowOpen = !isOpen;
      toggle.setAttribute('aria-expanded', String(nowOpen));
      body.classList.toggle('collapsed', !nowOpen);
      if (arrow)   arrow.classList.toggle('open', nowOpen);
      if (divider) divider.classList.toggle('hidden', !nowOpen);
    });
  });

  /* ── Disk section toggles ────────────────────────── */
  document.querySelectorAll('[id^="disk-toggle-"]').forEach(function (toggle) {
    const idx     = toggle.id.replace('disk-toggle-', '');
    const list    = document.getElementById('disk-list-' + idx);
    const chevron = document.getElementById('disk-chevron-' + idx);
    if (!list) return;
    toggle.addEventListener('click', function () {
      const nowOpen = toggle.getAttribute('aria-expanded') !== 'true';
      toggle.setAttribute('aria-expanded', String(nowOpen));
      list.classList.toggle('open', nowOpen);
      if (chevron) chevron.classList.toggle('open', nowOpen);
    });
  });

  /* ── Pool detail modal with focus trap ─────────────────── */
  const overlay    = document.getElementById('pool-detail-overlay');
  const closeBtn   = document.getElementById('pool-detail-close');
  const detailTitle = document.getElementById('pool-detail-title');
  const detailSub  = document.getElementById('pool-detail-sub');
  const detailBody = document.getElementById('pool-detail-body');
  const modal     = overlay ? overlay.querySelector('.pool-detail-modal') : null;
  const focusableSelector = 'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])';
  let lastFocused = null;
  let focusBeforeModal = null;

  function fmtBytes(value) {
    const units = ['B','KB','MB','GB','TB','PB','EB'];
    let n = Number(value) || 0;
    if (n < 1024) return Math.round(n) + ' B';
    let exp = 0;
    while (n >= 1024 && exp < units.length - 1) { n /= 1024; exp++; }
    return n.toFixed(2) + ' ' + units[exp];
  }

  function esc(v) {
    return String(v)
      .replaceAll('&','&amp;').replaceAll('<','&lt;')
      .replaceAll('>','&gt;').replaceAll('"','&quot;').replaceAll("'",'&#39;');
  }

  function metric(label, value, hot) {
    return '<div class="dataset-metric">'
      + '<span class="dataset-metric-label">' + esc(label) + '</span>'
      + '<span class="dataset-metric-value' + (hot ? ' hot' : '') + '">' + esc(value) + '</span>'
      + '</div>';
  }

  function renderDataset(ds) {
    let m = '';
    m += metric('Used',     fmtBytes(ds.used) + (ds.used_percent > 0 ? ' (' + ds.used_percent.toFixed(1) + '%)' : ''), false);
    m += metric('Available', fmtBytes(ds.available), false);
    m += metric('Quota',    ds.quota > 0 ? fmtBytes(ds.quota) + ' (' + ds.quota_used_percent.toFixed(1) + '%)' : 'none', ds.quota_used_percent >= 90);
    m += metric('Written',  fmtBytes(ds.written), false);
    m += metric('Logical',  fmtBytes(ds.logical_used), false);
    m += metric('Physical', fmtBytes(ds.used_by_dataset), false);
    if (ds.volume_size > 0) m += metric('Volume Size', fmtBytes(ds.volume_size) + ' (' + ds.volume_used_percent.toFixed(1) + '%)', false);
    m += metric('Referenced', fmtBytes(ds.referenced), false);
    return '<div class="dataset-item">'
      + '<div class="dataset-item-top">'
      + '<div class="dataset-name">' + esc(ds.name) + '</div>'
      + '<div class="dataset-kind">' + esc(ds.type) + '</div>'
      + '</div><div class="dataset-grid">' + m + '</div></div>';
  }

  function closeModal() {
    if (!overlay) return;
    overlay.hidden = true;
    overlay.classList.remove('open');
    if (detailBody) detailBody.innerHTML = '';
    if (focusBeforeModal) { focusBeforeModal.focus(); focusBeforeModal = null; }
    document.removeEventListener('keydown', trapKey);
  }

  function trapKey(e) {
    if (e.key !== 'Tab' || !modal) return;
    const foci = modal.querySelectorAll(focusableSelector);
    const first = foci[0], last = foci[foci.length - 1];
    if (e.shiftKey && document.activeElement === first) {
      e.preventDefault(); last.focus();
    } else if (!e.shiftKey && document.activeElement === last) {
      e.preventDefault(); first.focus();
    }
  }

  function openPoolDetail(ni, pi) {
    const node = nodes[ni];
    const pool = node && node.pools ? node.pools[pi] : null;
    if (!pool || !overlay) return;

    focusBeforeModal = document.activeElement;
    detailTitle.textContent = pool.name + ' datasets';
    detailSub.textContent = node.label
      + (node.location ? ' · ' + node.location : '')
      + ' · ' + pool.health;

    detailBody.innerHTML = (!pool.datasets || pool.datasets.length === 0)
      ? '<div class="pool-detail-empty">No dataset metrics available for this pool.</div>'
      : '<div class="dataset-list">' + pool.datasets.map(renderDataset).join('') + '</div>';

    overlay.hidden = false;
    overlay.classList.add('open');
    document.addEventListener('keydown', trapKey);
    if (modal) {
      const firstFocus = modal.querySelector(focusableSelector);
      if (firstFocus) firstFocus.focus();
    }
  }


  document.querySelectorAll('.pool-card[data-has-datasets="true"]').forEach(function (card) {
    card.addEventListener('click', function () {
      openPoolDetail(+card.dataset.nodeIndex, +card.dataset.poolIndex);
    });
  });

  if (closeBtn)  closeBtn.addEventListener('click', closeModal);
  if (overlay)   overlay.addEventListener('click', function (e) { if (e.target === overlay) closeModal(); });
  document.addEventListener('keydown', function (e) { if (e.key === 'Escape') closeModal(); });

  /* ── Animate usage bars ───────────────────────────── */
  if (!prefersReducedMotion) {
    document.querySelectorAll('.bar-fill[data-width]').forEach(function (el) {
      const target = el.getAttribute('data-width') + '%';
      requestAnimationFrame(function () {
        requestAnimationFrame(function () { el.style.width = target; });
      });
    });
  } else {
    document.querySelectorAll('.bar-fill[data-width]').forEach(function (el) {
      el.style.width = el.getAttribute('data-width') + '%';
    });
  }

  /* ── Animate entrance ─────────────────────────────── */
  if (!prefersReducedMotion) {
    requestAnimationFrame(function () {
      requestAnimationFrame(function () {
        document.querySelectorAll('.animate-in, .animate-in-modal').forEach(function (el) {
          el.style.opacity = '';
        });
      });
    });
  } else {
    document.querySelectorAll('.animate-in, .animate-in-modal').forEach(function (el) {
      el.style.opacity = '';
    });
  }

  /* ── Auto-refresh progress bar & reload ───────────── */
  const REFRESH_MS = {{.RefreshSecs}} * 1000;
  const rbar = document.getElementById('rbar');
  const startAt = Date.now();

  function tickBar() {
    const pct = Math.min((Date.now() - startAt) / REFRESH_MS * 100, 100);
    if (rbar) rbar.style.width = pct + '%';
    if (pct < 100) requestAnimationFrame(tickBar);
  }
  requestAnimationFrame(tickBar);

  // Fallback full-page reload
  const reloadTimer = setTimeout(function () { location.reload(); }, REFRESH_MS);

  /* ── Real-time SSE ───────────────────────────────── */
  function connectSSE() {
    const es = new EventSource('/events');

    es.onmessage = function(e) {
      if (e.data === 'refresh') {
        clearTimeout(reloadTimer);
        location.reload();
      }
    };

    es.onerror = function() {
      es.close();
      setTimeout(connectSSE, 5000);
    };
  }
  connectSSE();

  /* Live clock in topbar */
  const tsEl = document.getElementById('ts');
  if (tsEl) {
    tsEl.textContent = new Date().toLocaleTimeString([], {
      hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false
    });
  }
})();
</script>
</body>
</html>`
