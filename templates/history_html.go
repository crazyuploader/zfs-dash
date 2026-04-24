package templates

const historyHTML = `<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta name="color-scheme" content="dark light">
<title>ZFS Dashboard — History</title>
<link rel="icon" type="image/svg+xml" href="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 32 32'%3E%3Cellipse cx='16' cy='9' rx='12' ry='4' fill='%234f8ef7'/%3E%3Cpath d='M4 9v7c0 2.21 5.37 4 12 4s12-1.79 12-4V9c0 2.21-5.37 4-12 4S4 11.21 4 9z' fill='%233a7bd5'/%3E%3Cpath d='M4 16v7c0 2.21 5.37 4 12 4s12-1.79 12-4v-7c0 2.21-5.37 4-12 4S4 18.21 4 16z' fill='%232563b0'/%3E%3C/svg%3E">
<style>
:root {
  --font-body: 'DM Sans', system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  --font-mono: 'DM Mono', 'JetBrains Mono', ui-monospace, monospace;
  --text-sm: clamp(0.82rem, 0.78rem + 0.2vw, 0.9rem);
  --text-base: clamp(0.9rem, 0.85rem + 0.25vw, 1rem);
  --text-lg: clamp(1.05rem, 0.95rem + 0.5vw, 1.35rem);
  --space-2:0.5rem; --space-3:0.75rem; --space-4:1rem; --space-6:1.5rem; --space-8:2rem;
  --radius-md:0.5rem; --radius-lg:0.75rem; --radius-full:9999px;
  --transition: 180ms cubic-bezier(0.16, 1, 0.3, 1);
}
:root, [data-theme="dark"] {
  --bg: #0d0d0f; --surface: #131316; --surface-2: #18181c; --surface-3: #1e1e24;
  --border: rgba(255,255,255,0.07); --border-hi: rgba(255,255,255,0.13);
  --text: #e8e8ea; --text-muted: #85858f; --text-faint: #44444d;
  --primary: #4f98a3; --primary-dim: rgba(79,152,163,0.14);
  --success: #6daa45; --warning: #e8af34; --error: #dd6974;
  --shadow-md: 0 4px 20px rgba(0,0,0,0.45);
}
[data-theme="light"] {
  --bg: #f2f2ee; --surface: #fafaf8; --surface-2: #ffffff; --surface-3: #ededea;
  --border: rgba(0,0,0,0.07); --border-hi: rgba(0,0,0,0.13);
  --text: #1a1a1f; --text-muted: #62626b; --text-faint: #aaaaaf;
  --primary: #01696f; --primary-dim: rgba(1,105,111,0.1);
  --success: #437a22; --warning: #9c6a00; --error: #a13544;
  --shadow-md: 0 4px 20px rgba(0,0,0,0.09);
}
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
html { -webkit-font-smoothing: antialiased; color-scheme: dark; }
html[data-theme="light"] { color-scheme: light; }
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after { transition-duration: 0.01ms !important; animation-duration: 0.01ms !important; }
}
body {
  font-family: var(--font-body); font-size: var(--text-base);
  background: var(--bg); color: var(--text);
  min-height: 100dvh; display: flex; flex-direction: column;
}

/* ── Topbar ── */
.topbar {
  position: sticky; top: 0; z-index: 40;
  display: flex; align-items: center; gap: var(--space-3);
  padding: 0 var(--space-6); height: 52px;
  background: rgba(13,13,15,0.85); backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border);
}
[data-theme="light"] .topbar { background: rgba(242,242,238,0.9); }
.logo { display:flex; align-items:center; gap:var(--space-2); text-decoration:none; color:var(--text); font-weight:600; font-size:var(--text-base); }
.logo-sub { font-weight:400; color:var(--text-muted); margin-left:0.3em; font-size:var(--text-sm); }
.topbar-spacer { flex:1; }
.back-link { display:flex; align-items:center; gap:var(--space-2); color:var(--text-muted); text-decoration:none; font-size:var(--text-sm); padding:0.35rem 0.75rem; border-radius:var(--radius-full); border:1px solid var(--border); transition:var(--transition); white-space:nowrap; }
.back-link:hover { color:var(--text); border-color:var(--border-hi); }
.icon-btn { display:flex; align-items:center; justify-content:center; width:36px; height:36px; border:1px solid var(--border); border-radius:var(--radius-md); background:transparent; color:var(--text-muted); cursor:pointer; transition:var(--transition); touch-action:manipulation; flex-shrink:0; }
.icon-btn:hover { color:var(--text); border-color:var(--border-hi); }

/* ── Mobile sidebar toggle ── */
.sidebar-toggle {
  display:none; align-items:center; justify-content:center;
  width:36px; height:36px; border:1px solid var(--border); border-radius:var(--radius-md);
  background:transparent; color:var(--text-muted); cursor:pointer; transition:var(--transition);
  touch-action:manipulation; flex-shrink:0;
}
.sidebar-toggle:hover { color:var(--text); border-color:var(--border-hi); }
@media (max-width: 700px) { .sidebar-toggle { display:flex; } }

/* ── Layout ── */
.layout { display:flex; flex:1; min-height:0; overflow:hidden; }
.sidebar {
  width:260px; flex-shrink:0; border-right:1px solid var(--border);
  overflow-y:auto; overflow-x:hidden; padding:var(--space-4);
  display:flex; flex-direction:column; gap:var(--space-4);
  transition: transform var(--transition), opacity var(--transition);
}
.main { flex:1; display:flex; flex-direction:column; padding:var(--space-4) var(--space-6); gap:var(--space-6); overflow-y:auto; min-width:0; }

/* ── Sidebar sections ── */
.sidebar-label { font-size:0.6rem; font-weight:700; letter-spacing:0.1em; text-transform:uppercase; color:var(--text-faint); margin-bottom:var(--space-2); }
.range-pills { display:flex; flex-wrap:wrap; gap:var(--space-2); }
.range-pill {
  padding:0.25rem 0.75rem; border-radius:var(--radius-full); border:1px solid var(--border);
  background:transparent; color:var(--text-muted); font-size:var(--text-sm); cursor:pointer;
  font-family:var(--font-mono); transition:var(--transition); touch-action:manipulation;
}
.range-pill:hover { border-color:var(--border-hi); color:var(--text); }
.range-pill.active { background:var(--primary); border-color:var(--primary); color:#fff; }

/* ── Hierarchical series tree ── */
/* node → kind → disk → metrics */
.series-tree { display:flex; flex-direction:column; gap:var(--space-3); }

.node-group { display:flex; flex-direction:column; gap:var(--space-2); }
.node-header {
  font-size:0.65rem; font-weight:700; letter-spacing:0.1em; text-transform:uppercase;
  color:var(--primary); padding:0.15rem 0; border-bottom:1px solid var(--border); margin-bottom:0.25rem;
}

.kind-group { display:flex; flex-direction:column; gap:0.1rem; margin-bottom:var(--space-2); }
.kind-header {
  font-size:0.6rem; font-weight:700; letter-spacing:0.08em; text-transform:uppercase;
  color:var(--text-faint); padding:0.1rem 0 0.2rem;
}

/* Disk row — collapsible */
.disk-group { display:flex; flex-direction:column; border-radius:var(--radius-md); overflow:hidden; margin-bottom:0.15rem; }
.disk-toggle {
  display:flex; align-items:center; gap:var(--space-2); padding:0.3rem 0.4rem;
  cursor:pointer; border-radius:var(--radius-md); transition:var(--transition);
  background:transparent; border:none; width:100%; text-align:left; color:var(--text);
  touch-action:manipulation;
}
.disk-toggle:hover { background:var(--surface-2); }
.disk-chevron {
  width:12px; height:12px; flex-shrink:0; color:var(--text-faint);
  transition:transform 180ms ease;
}
.disk-group.open .disk-chevron { transform:rotate(90deg); }
.disk-name {
  font-family:var(--font-mono); font-size:0.78rem; color:var(--text);
  min-width:0; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; flex:1;
}

/* Metric items under a disk */
.disk-metrics {
  display:none; flex-direction:column; gap:0.05rem;
  padding:0.15rem 0 0.15rem var(--space-4); border-left:1px solid var(--border);
  margin-left:1.25rem;
}
.disk-group.open .disk-metrics { display:flex; }

.series-item {
  display:flex; align-items:center; gap:var(--space-2); padding:0.22rem 0.4rem;
  border-radius:var(--radius-md); cursor:pointer; transition:var(--transition);
  font-size:var(--text-sm); min-width:0; overflow:hidden;
}
.series-item:hover { background:var(--surface-2); }
.series-item input[type=checkbox] { accent-color:var(--primary); cursor:pointer; flex-shrink:0; touch-action:manipulation; }
.series-item-label { min-width:0; flex:1; font-family:var(--font-mono); }
.series-item-metric { color:var(--text-muted); font-size:0.72rem; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
.series-color-dot { width:7px; height:7px; border-radius:50%; flex-shrink:0; }
.series-empty { color:var(--text-faint); font-size:var(--text-sm); padding:var(--space-2) 0; }

/* ── Chart area ── */
.chart-card {
  background:var(--surface); border:1px solid var(--border); border-radius:var(--radius-lg);
  padding:var(--space-4) var(--space-4) var(--space-3);
  display:flex; flex-direction:column; gap:var(--space-3);
}
.chart-title { font-size:var(--text-sm); font-weight:600; color:var(--text-muted); }
.chart-num { font-variant-numeric:tabular-nums; }
.chart-canvas-wrap { position:relative; width:100%; }
canvas#chart { display:block; width:100%; }
.chart-empty { height:200px; display:flex; align-items:center; justify-content:center; color:var(--text-faint); font-size:var(--text-sm); }
.chart-legend { display:flex; flex-wrap:wrap; gap:var(--space-2) var(--space-3); }
.legend-item { display:flex; align-items:center; gap:var(--space-2); font-size:0.75rem; color:var(--text-muted); min-width:0; max-width:220px; }
.legend-item span { overflow:hidden; text-overflow:ellipsis; white-space:nowrap; min-width:0; font-family:var(--font-mono); }
.legend-line { width:16px; height:2px; border-radius:1px; flex-shrink:0; }

/* ── Tooltip — FIXED ── */
#chart-tooltip {
  position:fixed; pointer-events:none; z-index:100;
  background:var(--surface-3); border:1px solid var(--border-hi);
  border-radius:var(--radius-md);
  padding:0.5rem 0.65rem;
  font-size:0.76rem; font-family:var(--font-mono);
  box-shadow:var(--shadow-md); display:none;
  /* Key fix: constrain width and allow proper layout */
  width:220px;
  font-variant-numeric:tabular-nums;
}
.tooltip-time { color:var(--text-muted); font-size:0.7rem; margin-bottom:0.35rem; font-variant-numeric:tabular-nums; }
.tooltip-row {
  display:grid;
  /* dot | name (fills space) | value (right-aligned fixed width) */
  grid-template-columns: 8px 1fr auto;
  align-items:center;
  gap:0 0.4rem;
  padding:0.18rem 0;
  border-bottom:1px solid var(--border);
}
.tooltip-row:last-child { border-bottom:none; }
.tooltip-dot { width:7px; height:7px; border-radius:50%; flex-shrink:0; }
.tooltip-name {
  color:var(--text-muted); font-size:0.7rem;
  overflow:hidden; text-overflow:ellipsis; white-space:nowrap;
  min-width:0;
}
.tooltip-val {
  color:var(--text); font-size:0.76rem;
  white-space:nowrap; text-align:right;
  font-weight:600;
}

/* ── Status banner ── */
.status-banner { color:var(--text-faint); font-size:var(--text-sm); text-align:center; padding:var(--space-4); }

/* ── Mobile ── */
@media (max-width: 700px) {
  .topbar { padding: 0 var(--space-4); }
  .layout { flex-direction:column; position:relative; }
  .sidebar {
    position: absolute; top: 0; left: 0; bottom: 0; z-index: 30;
    width: 260px; transform: translateX(-100%);
    background: var(--bg); border-right: 1px solid var(--border-hi);
    box-shadow: var(--shadow-md);
  }
  .sidebar.open { transform: translateX(0); }
  .sidebar-overlay {
    display:none; position:fixed; inset:52px 0 0 0; z-index:29;
    background:rgba(0,0,0,0.4); backdrop-filter:blur(2px);
  }
  .sidebar-overlay.open { display:block; }
  .main { padding: var(--space-3); gap: var(--space-3); }
  .chart-card { padding: var(--space-3) var(--space-3) var(--space-2); }
  #chart-tooltip { width:190px; font-size:0.72rem; }
}
</style>
</head>
<body>

<header class="topbar">
  <a class="logo" href="/" aria-label="ZFS Dash home">
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor"
         stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
      <rect x="2" y="2"  width="20" height="6" rx="1.5"/>
      <rect x="2" y="10" width="20" height="6" rx="1.5"/>
      <rect x="2" y="18" width="20" height="4" rx="1.5"/>
      <circle cx="5.5" cy="5"  r="0.9" fill="currentColor" stroke="none"/>
      <circle cx="5.5" cy="13" r="0.9" fill="currentColor" stroke="none"/>
      <circle cx="5.5" cy="20" r="0.9" fill="currentColor" stroke="none"/>
    </svg>
    <span class="logo-name">zfs-dash<span class="logo-sub">history</span></span>
  </a>
  <div class="topbar-spacer"></div>
  <a class="back-link" href="/">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor"
         stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
      <path d="m15 18-6-6 6-6"/>
    </svg>
    Dashboard
  </a>
  <button class="sidebar-toggle" id="sidebar-toggle-btn" aria-label="Toggle sidebar">
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
      <line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/>
    </svg>
  </button>
  <button class="icon-btn" id="theme-btn" aria-label="Toggle theme"></button>
</header>

<!-- Mobile overlay -->
<div class="sidebar-overlay" id="sidebar-overlay"></div>

<div class="layout">
  <!-- Sidebar -->
  <aside class="sidebar" id="sidebar">
    <div>
      <div class="sidebar-label">Time Range</div>
      <div class="range-pills" id="range-pills">
        <button class="range-pill" data-hours="1">1h</button>
        <button class="range-pill" data-hours="6">6h</button>
        <button class="range-pill active" data-hours="24">24h</button>
        <button class="range-pill" data-hours="168">7d</button>
        <button class="range-pill" data-hours="720">30d</button>
      </div>
    </div>
    <div>
      <div class="sidebar-label">Series</div>
      <div class="series-tree" id="series-tree" aria-live="polite" aria-label="Series list">
        <div class="series-empty" id="series-loading">Loading…</div>
      </div>
    </div>
  </aside>

  <!-- Main -->
  <main class="main">
    <div class="chart-card">
      <div class="chart-title" id="chart-title">Select series from the sidebar</div>
      <div class="chart-canvas-wrap" id="chart-wrap">
        <div class="chart-empty" id="chart-empty">No series selected</div>
        <canvas id="chart" style="display:none"></canvas>
      </div>
      <div class="chart-legend" id="chart-legend"></div>
    </div>
  </main>
</div>

<div id="chart-tooltip">
  <div class="tooltip-time" id="tt-time"></div>
  <div id="tt-rows"></div>
</div>

<script>
(function () {
  /* ── Mobile sidebar toggle ── */
  const sidebarEl = document.getElementById('sidebar');
  const overlayEl = document.getElementById('sidebar-overlay');
  document.getElementById('sidebar-toggle-btn').addEventListener('click', () => {
    sidebarEl.classList.toggle('open');
    overlayEl.classList.toggle('open');
  });
  overlayEl.addEventListener('click', () => {
    sidebarEl.classList.remove('open');
    overlayEl.classList.remove('open');
  });

  /* ── Theme ── */
  const THEME_KEY = 'zfs-dash-theme';
  const moonSVG = '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>';
  const sunSVG  = '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>';
  const themeBtn = document.getElementById('theme-btn');
  let theme;
  try { theme = localStorage.getItem(THEME_KEY); } catch(_) { theme = null; }
  if (theme !== 'dark' && theme !== 'light') {
    theme = window.matchMedia('(prefers-color-scheme:dark)').matches ? 'dark' : 'light';
  }
  function applyTheme(t) {
    document.documentElement.setAttribute('data-theme', t);
    document.documentElement.style.colorScheme = t;
    themeBtn.innerHTML = t === 'dark' ? moonSVG : sunSVG;
    themeBtn.setAttribute('aria-label', 'Switch to ' + (t === 'dark' ? 'light' : 'dark') + ' mode');
  }
  applyTheme(theme);
  themeBtn.addEventListener('click', () => {
    theme = theme === 'dark' ? 'light' : 'dark';
    try { localStorage.setItem(THEME_KEY, theme); } catch(_) {}
    applyTheme(theme);
    if (activeSeries.length) renderChart();
  });

  /* ── Constants ── */
  const COLORS = ['#4f98a3','#e8af34','#dd6974','#6daa45','#9b6aba','#e07b54','#5b8dd9','#d45fa6','#3db08a','#f0c060'];
  const METRIC_LABELS = {
    used_pct: 'Usage %', alloc_bytes: 'Allocated', free_bytes: 'Free',
    temp_c: 'Temp °C', wear_pct: 'NVMe Wear %', wear_lvl: 'Wear Level',
    pow_hrs: 'Pwr-On Hrs',
  };
  // Short labels for tooltip value column
  const METRIC_SHORT = {
    used_pct: '%', wear_pct: '%', temp_c: '°C', pow_hrs: 'h',
    alloc_bytes: '', free_bytes: '',
  };

  /* ── State ── */
  let allSeries = [];
  let selected = new Set();
  let rangeHours = 24;
  let colorMap = {};
  let colorIdx = 0;
  let activeSeries = [];

  /* ── Fetch series list ── */
  async function loadSeries() {
    try {
      const res = await fetch('/api/history/series');
      allSeries = await res.json();
      buildTree();
    } catch (e) {
      document.getElementById('series-loading').textContent = 'Failed to load series.';
    }
  }

  /* ── Build hierarchical sidebar tree: node → kind → disk → metrics ── */
  function buildTree() {
    const tree = document.getElementById('series-tree');
    tree.innerHTML = '';
    if (!allSeries.length) {
      tree.innerHTML = '<div class="series-empty">No history data yet.<br>Data is recorded on each refresh cycle.</div>';
      return;
    }

    // Assign colors first
    for (const s of allSeries) {
      if (!colorMap[s.key]) colorMap[s.key] = COLORS[colorIdx++ % COLORS.length];
    }

    // Group: node → kind → diskName → [series]
    const byNode = {};
    for (const s of allSeries) {
      if (!byNode[s.node]) byNode[s.node] = {};
      if (!byNode[s.node][s.kind]) byNode[s.node][s.kind] = {};
      if (!byNode[s.node][s.kind][s.name]) byNode[s.node][s.kind][s.name] = [];
      byNode[s.node][s.kind][s.name].push(s);
    }

    const chevronSVG = '<svg class="disk-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>';

    for (const [node, kinds] of Object.entries(byNode)) {
      const ng = document.createElement('div');
      ng.className = 'node-group';

      const nh = document.createElement('div');
      nh.className = 'node-header';
      nh.textContent = node;
      ng.appendChild(nh);

      for (const [kind, disks] of Object.entries(kinds)) {
        const kg = document.createElement('div');
        kg.className = 'kind-group';

        const kh = document.createElement('div');
        kh.className = 'kind-header';
        kh.textContent = kind === 'pool' ? 'Pools' : 'Disks';
        kg.appendChild(kh);

        for (const [diskName, metrics] of Object.entries(disks)) {
          const shortName = diskName.replace(/^.*\/([^/]+)$/, '$1');

          const dg = document.createElement('div');
          dg.className = 'disk-group';
          // Auto-open if any metrics of this disk are selected
          if (metrics.some(s => selected.has(s.key))) dg.classList.add('open');

          // Toggle button
          const dt = document.createElement('button');
          dt.className = 'disk-toggle';
          dt.type = 'button';
          dt.title = diskName;
          dt.innerHTML = chevronSVG +
            '<span class="disk-name">' + escHtml(shortName) + '</span>';
          dt.addEventListener('click', () => dg.classList.toggle('open'));
          dg.appendChild(dt);

          // Metric checkboxes
          const dm = document.createElement('div');
          dm.className = 'disk-metrics';

          for (const s of metrics) {
            const item = document.createElement('label');
            item.className = 'series-item';

            const cb = document.createElement('input');
            cb.type = 'checkbox';
            cb.dataset.key = s.key;
            cb.checked = selected.has(s.key);
            cb.addEventListener('change', () => {
              if (cb.checked) {
                selected.add(s.key);
                dg.classList.add('open');
              } else {
                selected.delete(s.key);
              }
              refresh();
            });

            const dot = document.createElement('span');
            dot.className = 'series-color-dot';
            dot.style.background = colorMap[s.key];

            const lbl = document.createElement('span');
            lbl.className = 'series-item-label';
            lbl.innerHTML = '<span class="series-item-metric">' +
              escHtml(METRIC_LABELS[s.metric] || s.metric) + '</span>';

            item.appendChild(cb);
            item.appendChild(dot);
            item.appendChild(lbl);
            dm.appendChild(item);
          }

          dg.appendChild(dm);
          kg.appendChild(dg);
        }
        ng.appendChild(kg);
      }
      tree.appendChild(ng);
    }
  }

  /* ── Range pills ── */
  document.getElementById('range-pills').addEventListener('click', e => {
    const pill = e.target.closest('.range-pill');
    if (!pill) return;
    document.querySelectorAll('.range-pill').forEach(p => p.classList.remove('active'));
    pill.classList.add('active');
    rangeHours = parseInt(pill.dataset.hours, 10);
    refresh();
  });

  /* ── Bucket size ── */
  function bucketSecs(hours) {
    if (hours <= 1)   return 0;
    if (hours <= 6)   return 300;
    if (hours <= 24)  return 900;
    if (hours <= 168) return 3600;
    return 21600;
  }

  /* ── Fetch & render ── */
  async function refresh() {
    if (!selected.size) {
      activeSeries = [];
      renderChart();
      return;
    }
    const now = Math.floor(Date.now() / 1000);
    const from = now - rangeHours * 3600;
    const bucket = bucketSecs(rangeHours);

    const fetches = [...selected].map(async key => {
      const url = '/api/history/query?key=' + encodeURIComponent(key) +
        '&from=' + from + '&to=' + now + '&bucket=' + bucket;
      const res = await fetch(url);
      const points = await res.json();
      return { key, info: allSeries.find(s => s.key === key), points };
    });

    try {
      activeSeries = await Promise.all(fetches);
      renderChart();
    } catch (e) {
      console.error('history fetch failed', e);
    }
  }

  /* ── Chart renderer ── */
  const canvas = document.getElementById('chart');
  const ctx = canvas.getContext('2d');
  const wrap = document.getElementById('chart-wrap');
  const empty = document.getElementById('chart-empty');
  const legendEl = document.getElementById('chart-legend');
  const titleEl = document.getElementById('chart-title');
  const tooltip = document.getElementById('chart-tooltip');
  const ttTime = document.getElementById('tt-time');
  const ttRows = document.getElementById('tt-rows');
  const PAD = { top: 16, right: 16, bottom: 40, left: 62 };

  function renderChart() {
    const hasPts = activeSeries.some(s => s.points && s.points.length > 0);
    if (!activeSeries.length || !hasPts) {
      canvas.style.display = 'none';
      empty.style.display = 'flex';
      empty.textContent = selected.size ? 'No data in selected range.' : 'No series selected.';
      legendEl.innerHTML = '';
      titleEl.textContent = 'Select series from the sidebar';
      return;
    }
    canvas.style.display = 'block';
    empty.style.display = 'none';

    const metrics = [...new Set(activeSeries.map(s => s.info?.metric).filter(Boolean))];
    titleEl.textContent = metrics.map(m => METRIC_LABELS[m] || m).join(' · ');

    const dpr = window.devicePixelRatio || 1;
    const W = wrap.clientWidth;
    const H = Math.max(200, Math.min(400, W * 0.38));
    canvas.width = W * dpr;
    canvas.height = H * dpr;
    canvas.style.height = H + 'px';
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

    const plotW = W - PAD.left - PAD.right;
    const plotH = H - PAD.top - PAD.bottom;

    let minTs = Infinity, maxTs = -Infinity, minV = Infinity, maxV = -Infinity;
    for (const s of activeSeries) {
      for (const p of (s.points || [])) {
        if (p.ts < minTs) minTs = p.ts;
        if (p.ts > maxTs) maxTs = p.ts;
        if (p.v < minV)  minV = p.v;
        if (p.v > maxV)  maxV = p.v;
      }
    }
    if (!isFinite(minTs)) return;

    const tsSpan = maxTs - minTs || 1;
    const vRange = maxV - minV || 1;
    const vPad  = vRange * 0.12;
    const yMin  = Math.max(0, minV - vPad);
    const yMax  = maxV + vPad;

    const xScale = ts => PAD.left + ((ts - minTs) / tsSpan) * plotW;
    const yScale = v  => PAD.top + plotH - ((v - yMin) / (yMax - yMin)) * plotH;

    ctx.clearRect(0, 0, W, H);

    const isDark = document.documentElement.getAttribute('data-theme') !== 'light';
    const gridColor = isDark ? 'rgba(255,255,255,0.05)' : 'rgba(0,0,0,0.06)';
    const labelColor = isDark ? '#85858f' : '#62626b';
    const Y_STEPS = 4;
    ctx.font = '11px ' + getComputedStyle(document.body).getPropertyValue('--font-mono').trim();
    const yAxisMetric = activeSeries[0]?.info?.metric || metrics[0];

    for (let i = 0; i <= Y_STEPS; i++) {
      const v = yMin + (yMax - yMin) * (i / Y_STEPS);
      const y = yScale(v);
      ctx.strokeStyle = gridColor; ctx.lineWidth = 1;
      ctx.beginPath(); ctx.moveTo(PAD.left, y); ctx.lineTo(W - PAD.right, y); ctx.stroke();
      ctx.fillStyle = labelColor; ctx.textAlign = 'right';
      ctx.fillText(fmtVal(v, yAxisMetric), PAD.left - 6, y + 4);
    }

    const X_STEPS = Math.min(6, Math.floor(plotW / 70));
    for (let i = 0; i <= X_STEPS; i++) {
      const ts = minTs + tsSpan * (i / X_STEPS);
      const x = xScale(ts);
      ctx.fillStyle = labelColor; ctx.textAlign = 'center';
      ctx.fillText(fmtTime(ts, rangeHours), x, H - PAD.bottom + 16);
    }

    for (const s of activeSeries) {
      if (!s.points || !s.points.length) continue;
      const color = colorMap[s.key] || COLORS[0];
      ctx.save();
      ctx.beginPath();
      ctx.rect(PAD.left, PAD.top, plotW, plotH);
      ctx.clip();

      ctx.beginPath();
      let first = true;
      for (const p of s.points) {
        const x = xScale(p.ts), y = yScale(p.v);
        if (first) { ctx.moveTo(x, y); first = false; } else ctx.lineTo(x, y);
      }
      const lastPt = s.points[s.points.length - 1];
      const firstPt = s.points[0];
      ctx.lineTo(xScale(lastPt.ts), PAD.top + plotH);
      ctx.lineTo(xScale(firstPt.ts), PAD.top + plotH);
      ctx.closePath();
      ctx.fillStyle = color + '18';
      ctx.fill();

      ctx.beginPath();
      first = true;
      for (const p of s.points) {
        const x = xScale(p.ts), y = yScale(p.v);
        if (first) { ctx.moveTo(x, y); first = false; } else ctx.lineTo(x, y);
      }
      ctx.strokeStyle = color; ctx.lineWidth = 2; ctx.lineJoin = 'round';
      ctx.stroke();
      ctx.restore();
    }

    legendEl.innerHTML = '';
    for (const s of activeSeries) {
      if (!s.info) continue;
      const color = colorMap[s.key] || COLORS[0];
      const shortName = s.info.name.replace(/^.*\/([^/]+)$/, '$1');
      const li = document.createElement('div');
      li.className = 'legend-item';
      li.innerHTML = '<div class="legend-line" style="background:' + color + '"></div>' +
        '<span>' + escHtml(shortName) + ' · ' + escHtml(METRIC_LABELS[s.info.metric] || s.info.metric) + '</span>';
      legendEl.appendChild(li);
    }

    canvas._renderState = { xScale, yScale, tsSpan, minTs, maxTs, yMin, yMax, plotW, plotH, W, H };
  }

  /* ── Tooltip ── */
  function showTooltip(e) {
    const rs = canvas._renderState;
    if (!rs || !activeSeries.length) return;
    const rect = canvas.getBoundingClientRect();
    const mouseX = e.clientX - rect.left;
    if (mouseX < PAD.left || mouseX > rs.W - PAD.right) { tooltip.style.display = 'none'; return; }

    const ts = rs.minTs + ((mouseX - PAD.left) / rs.plotW) * rs.tsSpan;

    ttTime.textContent = fmtTimeFull(ts);
    ttRows.innerHTML = '';

    for (const s of activeSeries) {
      if (!s.points || !s.points.length) continue;
      let nearest = s.points[0];
      let minDist = Math.abs(s.points[0].ts - ts);
      for (const p of s.points) {
        const d = Math.abs(p.ts - ts);
        if (d < minDist) { minDist = d; nearest = p; }
      }
      const color = colorMap[s.key] || COLORS[0];
      // Short name: last segment, trimmed to ~12 chars
      const rawName = s.info ? s.info.name.replace(/^.*\/([^/]+)$/, '$1') : s.key;
      const shortName = rawName.length > 14 ? rawName.slice(0, 12) + '…' : rawName;
      const metricLabel = s.info ? (METRIC_SHORT[s.info.metric] !== undefined
        ? fmtValFull(nearest.v, s.info.metric)
        : fmtValFull(nearest.v, s.info.metric)) : nearest.v.toFixed(2);

      const row = document.createElement('div');
      row.className = 'tooltip-row';
      row.innerHTML =
        '<span class="tooltip-dot" style="background:' + color + '"></span>' +
        '<span class="tooltip-name">' + escHtml(shortName) + '</span>' +
        '<span class="tooltip-val">' + escHtml(fmtValFull(nearest.v, s.info?.metric)) + '</span>';
      ttRows.appendChild(row);
    }

    tooltip.style.display = 'block';

    // Smart positioning: keep tooltip within viewport
    const tw = tooltip.offsetWidth;
    const th = tooltip.offsetHeight;
    const vw = window.innerWidth;
    const vh = window.innerHeight;
    let tx = e.clientX + 16;
    let ty = e.clientY - 10;
    if (tx + tw > vw - 8) tx = e.clientX - tw - 16;
    if (ty + th > vh - 8) ty = vh - th - 8;
    if (ty < 8) ty = 8;
    tooltip.style.left = tx + 'px';
    tooltip.style.top = ty + 'px';
  }

  canvas.addEventListener('mousemove', showTooltip);
  canvas.addEventListener('mouseleave', () => { tooltip.style.display = 'none'; });

  // Touch support for tooltip
  canvas.addEventListener('touchmove', e => {
    e.preventDefault();
    const t = e.touches[0];
    showTooltip({ clientX: t.clientX, clientY: t.clientY });
  }, { passive: false });
  canvas.addEventListener('touchend', () => { tooltip.style.display = 'none'; });

  /* ── Helpers ── */
  function fmtVal(v, metric) {
    if (metric === 'alloc_bytes' || metric === 'free_bytes') return fmtBytes(v * (1 << 20));
    if (metric === 'used_pct' || metric === 'wear_pct') return v.toFixed(1) + '%';
    if (metric === 'temp_c') return v.toFixed(1) + '°';
    if (v >= 1e6) return (v/1e6).toFixed(1) + 'M';
    if (v >= 1e3) return (v/1e3).toFixed(1) + 'K';
    return v.toFixed(1);
  }
  function fmtValFull(v, metric) {
    if (!metric) return v.toFixed(2);
    if (metric === 'alloc_bytes' || metric === 'free_bytes') return fmtBytes(v * (1 << 20));
    if (metric === 'used_pct' || metric === 'wear_pct') return v.toFixed(2) + '%';
    if (metric === 'temp_c') return v.toFixed(1) + '°C';
    if (metric === 'pow_hrs') return v.toFixed(0) + 'h';
    return v.toFixed(2);
  }
  function fmtBytes(b) {
    const u = ['B','KB','MB','GB','TB','PB'];
    let i = 0;
    while (b >= 1024 && i < u.length - 1) { b /= 1024; i++; }
    return b.toFixed(1) + ' ' + u[i];
  }
  function fmtTime(ts, hours) {
    const d = new Date(ts * 1000);
    if (hours <= 24) return d.toLocaleTimeString([], {hour:'2-digit', minute:'2-digit', hour12:false});
    return d.toLocaleDateString([], {month:'short', day:'numeric'}) + ' ' +
           d.toLocaleTimeString([], {hour:'2-digit', minute:'2-digit', hour12:false});
  }
  function fmtTimeFull(ts) {
    const d = new Date(ts * 1000);
    return d.toLocaleDateString([], {month:'short', day:'numeric'}) + ' ' +
           d.toLocaleTimeString([], {hour:'2-digit', minute:'2-digit', second:'2-digit', hour12:false});
  }
  function escHtml(s) {
    if (s == null) return '';
    return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
  }

  /* ── Resize ── */
  const ro = new ResizeObserver(() => { if (activeSeries.length) renderChart(); });
  ro.observe(wrap);

  /* ── Init ── */
  loadSeries();
})();
</script>
</body>
</html>`
