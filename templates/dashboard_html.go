package templates

const dashboardHTML = `<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>ZFS Dashboard</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Geist+Mono:wght@400;500;600&family=Geist:wght@300;400;500;600&display=swap" rel="stylesheet">
<style>
/* ── Design Tokens ─────────────────────────────────── */
:root {
  --font-body: 'Geist', 'Inter', system-ui, sans-serif;
  --font-mono: 'Geist Mono', 'JetBrains Mono', monospace;

  --text-xs:   clamp(0.75rem,  0.7rem  + 0.25vw, 0.875rem);
  --text-sm:   clamp(0.875rem, 0.8rem  + 0.35vw, 1rem);
  --text-base: clamp(1rem,     0.95rem + 0.25vw, 1.125rem);
  --text-lg:   clamp(1.125rem, 1rem    + 0.75vw, 1.5rem);
  --text-xl:   clamp(1.5rem,   1.2rem  + 1.25vw, 2.25rem);

  --space-1:0.25rem; --space-2:0.5rem;  --space-3:0.75rem;
  --space-4:1rem;    --space-5:1.25rem; --space-6:1.5rem;
  --space-8:2rem;    --space-10:2.5rem; --space-12:3rem;

  --radius-sm:0.375rem; --radius-md:0.5rem;
  --radius-lg:0.75rem;  --radius-xl:1rem; --radius-full:9999px;

  --transition: 180ms cubic-bezier(0.16, 1, 0.3, 1);
}

/* ── Dark Theme (default) ───────────────────────────── */
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
  padding: var(--space-3) var(--space-6);
  display: flex; align-items: center; gap: var(--space-4);
}
.logo {
  display: flex; align-items: center; gap: var(--space-2);
  text-decoration: none; color: var(--text); flex-shrink: 0;
}
.logo-svg { color: var(--primary); }
.logo-name {
  font-size: var(--text-sm); font-weight: 600;
  letter-spacing: -0.01em;
}
.logo-sub {
  font-size: var(--text-xs); color: var(--text-muted);
  font-weight: 400; margin-left: var(--space-1);
}
.topbar-spacer { flex: 1; }
.topbar-meta {
  display: flex; align-items: center; gap: var(--space-3);
  font-size: var(--text-xs); color: var(--text-muted);
  font-family: var(--font-mono);
}
.live-dot {
  width: 6px; height: 6px; border-radius: var(--radius-full);
  background: var(--success); display: inline-block;
  animation: blink 2.4s ease-in-out infinite;
}
@keyframes blink { 0%,100%{opacity:1} 50%{opacity:0.35} }

.icon-btn {
  width: 32px; height: 32px; border-radius: var(--radius-md);
  display: flex; align-items: center; justify-content: center;
  color: var(--text-muted);
  transition: background var(--transition), color var(--transition);
}
.icon-btn:hover { background: var(--surface-3); color: var(--text); }

/* ── Main ────────────────────────────────────────────── */
.main {
  padding: var(--space-6) var(--space-6) var(--space-12);
  max-width: 1400px; margin-inline: auto; width: 100%;
}

/* ── Page Header ─────────────────────────────────────── */
.page-header { margin-bottom: var(--space-6); }
.page-title {
  font-size: var(--text-xl); font-weight: 600;
  letter-spacing: -0.025em; line-height: 1.1;
}
.page-sub {
  font-size: var(--text-xs); color: var(--text-muted);
  margin-top: var(--space-1); font-family: var(--font-mono);
}

/* ── KPI Row ─────────────────────────────────────────── */
.kpi-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: var(--space-3); margin-bottom: var(--space-8);
}
.kpi {
  background: var(--surface); border: 1px solid var(--border);
  border-radius: var(--radius-lg); padding: var(--space-4) var(--space-5);
  display: flex; flex-direction: column; gap: 2px;
  transition: border-color var(--transition);
}
.kpi:hover { border-color: var(--border-hi); }
.kpi-label {
  font-size: var(--text-xs); color: var(--text-muted);
  text-transform: uppercase; letter-spacing: 0.07em; font-weight: 500;
}
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
  font-size: var(--text-xs); color: var(--text-faint);
  font-family: var(--font-mono);
}

/* ── Node Section ────────────────────────────────────── */
.node { margin-bottom: var(--space-10); }
.node-head {
  display: flex; align-items: center; gap: var(--space-3);
  padding-bottom: var(--space-3); margin-bottom: var(--space-4);
  border-bottom: 1px solid var(--border);
}
.node-icon-wrap {
  width: 30px; height: 30px; border-radius: var(--radius-md);
  background: var(--primary-dim); color: var(--primary);
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.node-label-row { display: flex; align-items: center; gap: var(--space-2); flex-wrap: wrap; }
.node-label { font-size: var(--text-base); font-weight: 600; letter-spacing: -0.01em; }
.node-meta  { display: flex; flex-direction: column; gap: 4px; min-width: 0; }
.node-location {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 3px 9px; border-radius: var(--radius-full);
  background: color-mix(in oklab, var(--primary) 14%, transparent);
  border: 1px solid color-mix(in oklab, var(--primary) 24%, transparent);
  color: var(--primary); font-size: 0.68rem; font-weight: 700;
  letter-spacing: 0.04em; text-transform: uppercase;
}
.node-location::before {
  content: "";
  width: 5px; height: 5px; border-radius: 50%;
  background: currentColor; opacity: 0.9;
}
.node-url   { font-size: var(--text-xs); color: var(--text-faint); font-family: var(--font-mono); }
.node-ts    { margin-left: auto; font-size: var(--text-xs); color: var(--text-faint); font-family: var(--font-mono); }

/* ── Error Banner ────────────────────────────────────── */
.node-err {
  background: var(--error-dim); border: 1px solid color-mix(in oklab, var(--error) 30%, transparent);
  border-radius: var(--radius-lg); padding: var(--space-4) var(--space-5);
  color: var(--error); font-size: var(--text-sm);
  display: flex; align-items: center; gap: var(--space-3);
}

/* ── Empty State ─────────────────────────────────────── */
.empty {
  text-align: center; padding: var(--space-10) var(--space-8);
  color: var(--text-muted); display: flex; flex-direction: column;
  align-items: center; gap: var(--space-3);
}
.empty svg { color: var(--text-faint); }
.empty h3  { color: var(--text); font-size: var(--text-base); font-weight: 600; }
.empty p   { font-size: var(--text-sm); max-width: 38ch; }

/* ── Pool Grid ───────────────────────────────────────── */
.pools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(340px, 100%), 1fr));
  gap: var(--space-4);
}

/* ── Pool Card ───────────────────────────────────────── */
.pool-card {
  background: var(--surface); border: 1px solid var(--border);
  border-radius: var(--radius-xl); padding: var(--space-5);
  display: flex; flex-direction: column; gap: var(--space-4);
  transition: box-shadow var(--transition), border-color var(--transition), transform var(--transition);
  will-change: transform;
}
.pool-card:hover {
  border-color: var(--border-hi);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
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
  font-size: 0.68rem; font-weight: 700;
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
  display: flex; justify-content: space-between;
  font-size: var(--text-xs); font-family: var(--font-mono);
  color: var(--text-faint);
}
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
  font-size: 0.68rem; color: var(--text-faint);
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
  font-size: 0.68rem; font-family: var(--font-mono); font-weight: 500;
  padding: 2px var(--space-2); border-radius: var(--radius-sm);
  background: var(--surface-3); color: var(--text-faint);
}
.chip.hot { background: var(--error-dim); color: var(--error); }

/* ── Refresh Progress Bar ────────────────────────────── */
#rbar {
  position: fixed; bottom: 0; left: 0; height: 2px;
  background: var(--primary); width: 0%; z-index: 999;
  opacity: 0.55; pointer-events: none;
}

/* ── Responsive ──────────────────────────────────────── */
@media (max-width: 640px) {
  .main    { padding: var(--space-4); }
  .topbar  { padding: var(--space-3) var(--space-4); }
  .kpi-row { grid-template-columns: 1fr 1fr; }
  .pools-grid { grid-template-columns: 1fr; }
  .topbar-meta span:nth-child(n+3) { display: none; }
  .node-head { align-items: flex-start; }
  .node-meta { flex: 1; }
  .node-ts { margin-left: 0; }
}

@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
</style>
</head>
<body>

<!-- ── Top Bar ────────────────────────────────────────── -->
<header class="topbar">
  <a class="logo" href="/" aria-label="ZFS Dash home">
    <!-- SVG logo: server rack symbol -->
    <svg class="logo-svg" width="22" height="22" viewBox="0 0 24 24" fill="none"
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
    <span aria-hidden="true">·</span>
    <span>↻&thinsp;{{.RefreshSecs}}s</span>
  </div>

  <button class="icon-btn" data-theme-toggle aria-label="Toggle light/dark mode">
    <!-- moon icon (shown in dark mode) -->
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none"
         stroke="currentColor" stroke-width="2" aria-hidden="true">
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
    </svg>
  </button>
</header>

<!-- ── Main ───────────────────────────────────────────── -->
<main class="main">

  <!-- Page header -->
  <div class="page-header">
    <h1 class="page-title">Pool Overview</h1>
    <p class="page-sub">{{.TotalNodes}} node{{if gt .TotalNodes 1}}s{{end}} &middot; {{.TotalPools}} pool{{if gt .TotalPools 1}}s{{end}} &middot; fetched {{.FetchedAt}}</p>
  </div>

  <!-- KPI summary -->
  <div class="kpi-row" role="region" aria-label="Fleet summary">
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

  <!-- Per-node sections -->
  {{range .Nodes}}
  <section class="node" aria-label="Node {{.Label}}">

    <div class="node-head">
      <div class="node-icon-wrap" aria-hidden="true">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none"
             stroke="currentColor" stroke-width="2" stroke-linecap="round">
          <rect x="2" y="3" width="20" height="18" rx="2"/>
          <path d="M8 7h8M8 12h8M8 17h4"/>
        </svg>
      </div>
      <div class="node-meta">
        <div class="node-label-row">
          <div class="node-label">{{.Label}}</div>
          {{if .Location}}<div class="node-location">{{.Location}}</div>{{end}}
        </div>
        <div class="node-url">{{.URL}}</div>
      </div>
      <div class="node-ts" aria-label="Fetched at">{{fmtNodeTime .FetchedAt}}</div>
    </div>

    {{if .Error}}
    <div class="node-err" role="alert">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
           stroke-width="2" aria-hidden="true">
        <circle cx="12" cy="12" r="10"/>
        <line x1="12" y1="8" x2="12" y2="12"/>
        <line x1="12" y1="16" x2="12.01" y2="16"/>
      </svg>
      <span>{{.Error}}</span>
    </div>
    {{else}}
      {{if not .Pools}}
      <div class="empty">
        <svg width="42" height="42" viewBox="0 0 24 24" fill="none" stroke="currentColor"
             stroke-width="1.4" aria-hidden="true">
          <rect x="2" y="3" width="20" height="18" rx="2"/>
          <path d="M8 7h8M8 12h8M8 17h4"/>
        </svg>
        <h3>No pools found</h3>
        <p>The exporter returned no ZFS pool metrics for this node.</p>
      </div>
      {{else}}
      <div class="pools-grid">
        {{range .Pools}}
        <article class="pool-card" aria-label="Pool {{.Name}}">

          <!-- Name + health -->
          <div class="pool-top">
            <div>
              <div class="pool-name">{{.Name}}</div>
              <div class="pool-sub">zfs pool</div>
            </div>
            <span class="hbadge {{healthClass .Health}}">
              <span class="hbadge-dot" aria-hidden="true"></span>
              {{.Health}}
            </span>
          </div>

          <!-- Capacity bar (only when size is known) -->
          {{if gt .Size 0.0}}
          <div class="usage">
            <div class="usage-header">
              <span>Capacity used</span>
              <span class="usage-pct">{{printf "%.1f" .UsedPercent}}%</span>
            </div>
            <div class="bar-track"
                 role="progressbar"
                 aria-valuenow="{{printf "%.0f" .UsedPercent}}"
                 aria-valuemin="0" aria-valuemax="100"
                 aria-label="{{printf "%.1f" .UsedPercent}}% used">
              <div class="bar-fill{{if gte .UsedPercent 90.0}} bad{{else if gte .UsedPercent 75.0}} warn{{end}}"
                   data-width="{{printf "%.2f" .UsedPercent}}"
                   style="width:0%"></div>
            </div>
            <div class="usage-sizes">
              <span>Used&nbsp;<b>{{humanBytes .Allocated}}</b></span>
              <span>Free&nbsp;<b>{{humanBytes .Free}}</b></span>
              <span>Total&nbsp;<b>{{humanBytes .Size}}</b></span>
            </div>
          </div>
          {{end}}

          <hr class="divider" aria-hidden="true">

          <!-- Pool stats sourced from zfs_exporter pool metrics -->
          <div class="stats">
            <div class="stat">
              <span class="stat-lbl">Dedup</span>
              <span class="stat-val">{{printf "%.2fx" .DedupRatio}}</span>
            </div>
            <div class="stat">
              <span class="stat-lbl">Fragmentation</span>
              <span class="stat-val">{{printf "%.0f%%" (mul100 .FragmentationRatio)}}</span>
            </div>
            <div class="stat">
              <span class="stat-lbl">Freeing</span>
              <span class="stat-val">{{humanBytes .Freeing}}</span>
            </div>
            <div class="stat">
              <span class="stat-lbl">Leaked</span>
              <span class="stat-val">{{humanBytes .LeakedBytes}}</span>
            </div>
          </div>

          <!-- Pool state chips -->
          <div class="err-row">
            <span class="chip{{if .ReadOnly}} hot{{end}}">Readonly&nbsp;{{if .ReadOnly}}yes{{else}}no{{end}}</span>
            <span class="chip{{if gt0 .Freeing}} hot{{end}}">Freeing&nbsp;{{humanBytes .Freeing}}</span>
            <span class="chip{{if gt0 .LeakedBytes}} hot{{end}}">Leaked&nbsp;{{humanBytes .LeakedBytes}}</span>
          </div>

        </article>
        {{end}}
      </div>
      {{end}}
    {{end}}

  </section>
  {{end}}

</main>

<!-- auto-refresh progress bar -->
<div id="rbar" aria-hidden="true"></div>

<script>
(function () {
  'use strict';

  /* ── Theme toggle ─────────────────────────────────── */
  const html = document.documentElement;
  const btn  = document.querySelector('[data-theme-toggle]');
  const themeKey = 'zfs-dash-theme';
  let theme;

  try {
    theme = localStorage.getItem(themeKey);
  } catch (_) {
    theme = null;
  }

  if (theme !== 'dark' && theme !== 'light') {
    theme = window.matchMedia('(prefers-color-scheme:dark)').matches ? 'dark' : 'light';
  }

  html.setAttribute('data-theme', theme);

  const moonSVG = '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>';
  const sunSVG  = '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>';

  function applyTheme(t) {
    html.setAttribute('data-theme', t);
    if (btn) {
      btn.innerHTML = t === 'dark' ? moonSVG : sunSVG;
      btn.setAttribute('aria-label', 'Switch to ' + (t === 'dark' ? 'light' : 'dark') + ' mode');
    }
  }
  applyTheme(theme);

  if (btn) {
    btn.addEventListener('click', function () {
      theme = theme === 'dark' ? 'light' : 'dark';
      try {
        localStorage.setItem(themeKey, theme);
      } catch (_) {}
      applyTheme(theme);
    });
  }

  /* ── Animate usage bars ───────────────────────────── */
  document.querySelectorAll('.bar-fill[data-width]').forEach(function (el) {
    const target = el.getAttribute('data-width') + '%';
    requestAnimationFrame(function () {
      requestAnimationFrame(function () {
        el.style.width = target;
      });
    });
  });

  /* ── Auto-refresh progress bar ────────────────────── */
  const REFRESH_MS = {{.RefreshSecs}} * 1000;
  const rbar = document.getElementById('rbar');
  const startAt = Date.now();

  function tickBar() {
    const elapsed = Date.now() - startAt;
    const pct = Math.min(elapsed / REFRESH_MS * 100, 100);
    if (rbar) rbar.style.width = pct + '%';
    if (pct < 100) requestAnimationFrame(tickBar);
  }
  requestAnimationFrame(tickBar);

  /* Reload page after refresh interval */
  setTimeout(function () { location.reload(); }, REFRESH_MS);

  /* Update the timestamp shown in the topbar */
  const tsEl = document.getElementById('ts');
  if (tsEl) {
    tsEl.textContent = new Date().toLocaleTimeString([], {
      hour: '2-digit', minute: '2-digit', second: '2-digit'
    });
  }
})();
</script>
</body>
</html>`
