package templates

const historyHTML = `<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
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
html { -webkit-font-smoothing: antialiased; }
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
.back-link { display:flex; align-items:center; gap:var(--space-2); color:var(--text-muted); text-decoration:none; font-size:var(--text-sm); padding:0.35rem 0.75rem; border-radius:var(--radius-full); border:1px solid var(--border); transition:var(--transition); }
.back-link:hover { color:var(--text); border-color:var(--border-hi); }
.icon-btn { display:flex; align-items:center; justify-content:center; width:36px; height:36px; border:1px solid var(--border); border-radius:var(--radius-md); background:transparent; color:var(--text-muted); cursor:pointer; transition:var(--transition); }
.icon-btn:hover { color:var(--text); border-color:var(--border-hi); }

/* ── Layout ── */
.layout { display:flex; flex:1; min-height:0; }
.sidebar {
  width:280px; flex-shrink:0; border-right:1px solid var(--border);
  overflow-y:auto; padding:var(--space-4);
  display:flex; flex-direction:column; gap:var(--space-4);
}
.main { flex:1; display:flex; flex-direction:column; padding:var(--space-6); gap:var(--space-6); overflow-y:auto; }

/* ── Sidebar sections ── */
.sidebar-label { font-size:0.6rem; font-weight:700; letter-spacing:0.1em; text-transform:uppercase; color:var(--text-faint); margin-bottom:var(--space-2); }
.range-pills { display:flex; flex-wrap:wrap; gap:var(--space-2); }
.range-pill {
  padding:0.25rem 0.75rem; border-radius:var(--radius-full); border:1px solid var(--border);
  background:transparent; color:var(--text-muted); font-size:var(--text-sm); cursor:pointer;
  font-family:var(--font-mono); transition:var(--transition);
}
.range-pill:hover { border-color:var(--border-hi); color:var(--text); }
.range-pill.active { background:var(--primary); border-color:var(--primary); color:#fff; }

.series-tree { display:flex; flex-direction:column; gap:var(--space-2); }
.series-node-group { display:flex; flex-direction:column; gap:0.25rem; }
.series-node-label { font-size:var(--text-sm); font-weight:600; color:var(--text); padding:0.2rem 0; }
.series-kind-group { margin-left:var(--space-3); display:flex; flex-direction:column; gap:0.15rem; }
.series-kind-label { font-size:0.65rem; text-transform:uppercase; letter-spacing:0.08em; color:var(--text-faint); padding:0.15rem 0; }
.series-name-group { margin-left:var(--space-3); display:flex; flex-direction:column; gap:0.1rem; }
.series-item {
  display:flex; align-items:center; gap:var(--space-2); padding:0.25rem 0.4rem;
  border-radius:var(--radius-md); cursor:pointer; transition:var(--transition);
  font-size:var(--text-sm);
}
.series-item:hover { background:var(--surface-2); }
.series-item input[type=checkbox] { accent-color:var(--primary); cursor:pointer; flex-shrink:0; }
.series-item-label { color:var(--text-muted); flex:1; font-family:var(--font-mono); font-size:0.78rem; }
.series-color-dot { width:8px; height:8px; border-radius:50%; flex-shrink:0; }
.series-empty { color:var(--text-faint); font-size:var(--text-sm); padding:var(--space-2) 0; }

/* ── Chart area ── */
.chart-card {
  background:var(--surface); border:1px solid var(--border); border-radius:var(--radius-lg);
  padding:var(--space-4) var(--space-4) var(--space-3);
  display:flex; flex-direction:column; gap:var(--space-3);
}
.chart-title { font-size:var(--text-sm); font-weight:600; color:var(--text-muted); }
.chart-canvas-wrap { position:relative; width:100%; }
canvas#chart { display:block; width:100%; }
.chart-empty { height:200px; display:flex; align-items:center; justify-content:center; color:var(--text-faint); font-size:var(--text-sm); }
.chart-legend { display:flex; flex-wrap:wrap; gap:var(--space-3); }
.legend-item { display:flex; align-items:center; gap:var(--space-2); font-size:var(--text-sm); color:var(--text-muted); }
.legend-line { width:18px; height:2px; border-radius:1px; }

/* ── Tooltip ── */
#chart-tooltip {
  position:fixed; pointer-events:none; z-index:100;
  background:var(--surface-3); border:1px solid var(--border-hi);
  border-radius:var(--radius-md); padding:0.4rem 0.75rem;
  font-size:var(--text-sm); font-family:var(--font-mono);
  box-shadow:var(--shadow-md); display:none;
  max-width:220px;
}
.tooltip-time { color:var(--text-muted); font-size:0.72rem; margin-bottom:0.2rem; }
.tooltip-row { display:flex; align-items:center; gap:var(--space-2); }
.tooltip-dot { width:7px; height:7px; border-radius:50%; flex-shrink:0; }

/* ── Status banner ── */
.status-banner { color:var(--text-faint); font-size:var(--text-sm); text-align:center; padding:var(--space-4); }

@media (max-width: 700px) {
  .sidebar { width:100%; border-right:none; border-bottom:1px solid var(--border); }
  .layout { flex-direction:column; }
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
  <button class="icon-btn" id="theme-btn" aria-label="Toggle theme">
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor"
         stroke-width="2" aria-hidden="true">
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
    </svg>
  </button>
</header>

<div class="layout">
  <!-- Sidebar -->
  <aside class="sidebar">
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
      <div class="series-tree" id="series-tree">
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
  /* ── Theme ── */
  const saved = localStorage.getItem('theme');
  if (saved) document.documentElement.setAttribute('data-theme', saved);
  document.getElementById('theme-btn').addEventListener('click', () => {
    const t = document.documentElement.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', t);
    localStorage.setItem('theme', t);
  });

  /* ── Constants ── */
  const COLORS = ['#4f98a3','#e8af34','#dd6974','#6daa45','#9b6aba','#e07b54','#5b8dd9','#d45fa6'];
  const METRIC_LABELS = {
    used_pct: 'Usage %', alloc_bytes: 'Allocated', free_bytes: 'Free',
    temp_c: 'Temp °C', wear_pct: 'NVMe Wear %', wear_lvl: 'Wear Leveling',
    pow_hrs: 'Power-On Hours',
  };
  const METRIC_UNIT = {
    used_pct: '%', wear_pct: '%', temp_c: '°C', pow_hrs: 'h',
    alloc_bytes: 'B', free_bytes: 'B',
  };

  /* ── State ── */
  let allSeries = [];
  let selected = new Set();
  let rangeHours = 24;
  let colorMap = {};
  let colorIdx = 0;
  let activeSeries = []; // [{info, points}]

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

  /* ── Build sidebar tree ── */
  function buildTree() {
    const tree = document.getElementById('series-tree');
    tree.innerHTML = '';
    if (!allSeries.length) {
      tree.innerHTML = '<div class="series-empty">No history data yet.<br>Data is recorded on each refresh cycle.</div>';
      return;
    }

    // Group: node → kind → name → [metrics]
    const byNode = {};
    for (const s of allSeries) {
      if (!byNode[s.node]) byNode[s.node] = {};
      if (!byNode[s.node][s.kind]) byNode[s.node][s.kind] = {};
      if (!byNode[s.node][s.kind][s.name]) byNode[s.node][s.kind][s.name] = [];
      byNode[s.node][s.kind][s.name].push(s);
    }

    for (const [node, kinds] of Object.entries(byNode)) {
      const ng = document.createElement('div');
      ng.className = 'series-node-group';
      const nl = document.createElement('div');
      nl.className = 'series-node-label';
      nl.textContent = node;
      ng.appendChild(nl);

      for (const [kind, names] of Object.entries(kinds)) {
        const kg = document.createElement('div');
        kg.className = 'series-kind-group';
        const kl = document.createElement('div');
        kl.className = 'series-kind-label';
        kl.textContent = kind === 'pool' ? 'Pools' : 'Disks';
        kg.appendChild(kl);

        for (const [name, metrics] of Object.entries(names)) {
          const nameg = document.createElement('div');
          nameg.className = 'series-name-group';

          for (const s of metrics) {
            if (!colorMap[s.key]) colorMap[s.key] = COLORS[colorIdx++ % COLORS.length];
            const item = document.createElement('label');
            item.className = 'series-item';

            const cb = document.createElement('input');
            cb.type = 'checkbox';
            cb.dataset.key = s.key;
            cb.addEventListener('change', () => {
              if (cb.checked) selected.add(s.key); else selected.delete(s.key);
              refresh();
            });

            const dot = document.createElement('span');
            dot.className = 'series-color-dot';
            dot.style.background = colorMap[s.key];

            const lbl = document.createElement('span');
            lbl.className = 'series-item-label';
            lbl.title = name;
            const shortName = name.replace(/^.*\/([^/]+)$/, '$1'); // last path component for devices
            lbl.textContent = shortName + ' · ' + (METRIC_LABELS[s.metric] || s.metric);

            item.appendChild(cb);
            item.appendChild(dot);
            item.appendChild(lbl);
            nameg.appendChild(item);
          }
          kg.appendChild(nameg);
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

  /* ── Bucket size based on range ── */
  function bucketSecs(hours) {
    if (hours <= 1)   return 0;     // raw
    if (hours <= 6)   return 300;   // 5-min
    if (hours <= 24)  return 900;   // 15-min
    if (hours <= 168) return 3600;  // 1-hour
    return 21600;                   // 6-hour for 30d
  }

  /* ── Fetch and render ── */
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

    // Update title
    const metrics = [...new Set(activeSeries.map(s => s.info?.metric).filter(Boolean))];
    titleEl.textContent = metrics.map(m => METRIC_LABELS[m] || m).join(' · ');

    // DPR-aware sizing
    const dpr = window.devicePixelRatio || 1;
    const W = wrap.clientWidth;
    const H = Math.max(240, Math.min(400, W * 0.38));
    canvas.width = W * dpr;
    canvas.height = H * dpr;
    canvas.style.height = H + 'px';
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

    const plotW = W - PAD.left - PAD.right;
    const plotH = H - PAD.top - PAD.bottom;

    // Data bounds
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

    // Grid + Y labels
    const isDark = document.documentElement.getAttribute('data-theme') !== 'light';
    const gridColor = isDark ? 'rgba(255,255,255,0.05)' : 'rgba(0,0,0,0.06)';
    const labelColor = isDark ? '#85858f' : '#62626b';
    const Y_STEPS = 4;
    ctx.font = '11px ' + getComputedStyle(document.body).getPropertyValue('--font-mono').trim();
    for (let i = 0; i <= Y_STEPS; i++) {
      const v = yMin + (yMax - yMin) * (i / Y_STEPS);
      const y = yScale(v);
      ctx.strokeStyle = gridColor; ctx.lineWidth = 1;
      ctx.beginPath(); ctx.moveTo(PAD.left, y); ctx.lineTo(W - PAD.right, y); ctx.stroke();
      ctx.fillStyle = labelColor; ctx.textAlign = 'right';
      ctx.fillText(fmtVal(v, metrics[0]), PAD.left - 6, y + 4);
    }

    // X labels
    const X_STEPS = Math.min(6, Math.floor(plotW / 70));
    for (let i = 0; i <= X_STEPS; i++) {
      const ts = minTs + tsSpan * (i / X_STEPS);
      const x = xScale(ts);
      ctx.fillStyle = labelColor; ctx.textAlign = 'center';
      ctx.fillText(fmtTime(ts, rangeHours), x, H - PAD.bottom + 16);
    }

    // Series lines + fill
    for (const s of activeSeries) {
      if (!s.points || !s.points.length) continue;
      const color = colorMap[s.key] || COLORS[0];
      ctx.save();
      ctx.beginPath();
      ctx.rect(PAD.left, PAD.top, plotW, plotH);
      ctx.clip();

      // Fill
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

      // Line
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

    // Legend
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

    // Store render state for tooltip
    canvas._renderState = { xScale, yScale, tsSpan, minTs, maxTs, yMin, yMax, plotW, plotH, W, H, metrics };
  }

  /* ── Tooltip ── */
  canvas.addEventListener('mousemove', e => {
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
      // Find nearest point
      let nearest = s.points[0];
      let minDist = Math.abs(s.points[0].ts - ts);
      for (const p of s.points) {
        const d = Math.abs(p.ts - ts);
        if (d < minDist) { minDist = d; nearest = p; }
      }
      const color = colorMap[s.key] || COLORS[0];
      const shortName = s.info ? s.info.name.replace(/^.*\/([^/]+)$/, '$1') : s.key;
      const row = document.createElement('div');
      row.className = 'tooltip-row';
      row.innerHTML = '<span class="tooltip-dot" style="background:' + color + '"></span>' +
        '<span style="color:var(--text-muted);font-size:0.72rem">' + escHtml(shortName) + '</span>' +
        '<span style="margin-left:auto;color:var(--text)">' + fmtValFull(nearest.v, rs.metrics[0]) + '</span>';
      ttRows.appendChild(row);
    }

    tooltip.style.display = 'block';
    const tx = e.clientX + 14;
    const ty = e.clientY - 10;
    tooltip.style.left = (tx + tooltip.offsetWidth > window.innerWidth ? e.clientX - tooltip.offsetWidth - 14 : tx) + 'px';
    tooltip.style.top = ty + 'px';
  });
  canvas.addEventListener('mouseleave', () => { tooltip.style.display = 'none'; });

  /* ── Helpers ── */
  function fmtVal(v, metric) {
    if (metric === 'alloc_bytes' || metric === 'free_bytes') return fmtBytes(v);
    if (metric === 'used_pct' || metric === 'wear_pct') return v.toFixed(1) + '%';
    if (metric === 'temp_c') return v.toFixed(1) + '°';
    if (v >= 1e6) return (v/1e6).toFixed(1) + 'M';
    if (v >= 1e3) return (v/1e3).toFixed(1) + 'K';
    return v.toFixed(1);
  }
  function fmtValFull(v, metric) {
    if (metric === 'alloc_bytes' || metric === 'free_bytes') return fmtBytes(v);
    if (metric === 'used_pct' || metric === 'wear_pct') return v.toFixed(2) + '%';
    if (metric === 'temp_c') return v.toFixed(1) + ' °C';
    if (metric === 'pow_hrs') return v.toFixed(1) + ' h';
    return v.toFixed(2);
  }
  function fmtBytes(b) {
    const u = ['B','KB','MB','GB','TB','PB'];
    let i = 0;
    while (b >= 1024 && i < u.length - 1) { b /= 1024; i++; }
    return b.toFixed(2) + ' ' + u[i];
  }
  function fmtTime(ts, hours) {
    const d = new Date(ts * 1000);
    if (hours <= 24) return d.toLocaleTimeString([], {hour:'2-digit', minute:'2-digit'});
    return d.toLocaleDateString([], {month:'short', day:'numeric'}) + ' ' + d.toLocaleTimeString([], {hour:'2-digit', minute:'2-digit'});
  }
  function fmtTimeFull(ts) {
    const d = new Date(ts * 1000);
    return d.toLocaleDateString([], {month:'short', day:'numeric'}) + ' ' +
           d.toLocaleTimeString([], {hour:'2-digit', minute:'2-digit', second:'2-digit'});
  }
  function escHtml(s) {
    return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
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
