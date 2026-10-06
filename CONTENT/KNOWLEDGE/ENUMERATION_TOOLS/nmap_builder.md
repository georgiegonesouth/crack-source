# NMAP COMMAND BUILDER

<style>
#nmap-cmd {
  position: sticky;
  top: 0;
  z-index: 10;
  background: var(--bg-header);
  border: 1px solid var(--border-subtle);
  border-radius: 7px;
  padding: 12px 16px;
  margin-bottom: 28px;
}
#nmap-cmd-label {
  font-family: var(--mono);
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.14em;
  color: var(--text-dim);
  margin-bottom: 8px;
}
#nmap-cmd-row {
  display: flex;
  align-items: center;
  gap: 10px;
}
#nmap-cmd-text {
  flex: 1;
  background: var(--bg-app);
  border: 1px solid var(--border);
  border-radius: 5px;
  padding: 8px 12px;
  font-family: var(--mono);
  font-size: 13px;
  color: var(--accent);
  word-break: break-all;
  line-height: 1.5;
  min-height: 38px;
  cursor: text;
  outline: none;
  transition: border-color .12s, box-shadow .12s;
}
#nmap-cmd-text:focus {
  border-color: var(--accent);
  box-shadow: 0 0 0 2px rgba(88,166,255,.12);
}
#nmap-reset {
  background: var(--bg-app);
  border: 1px solid var(--border-subtle);
  border-radius: 5px;
  padding: 6px 10px;
  font-family: var(--mono);
  font-size: 13px;
  cursor: pointer;
  color: var(--text-muted);
  line-height: 1;
  transition: color .12s, border-color .12s;
  display: none;
}
#nmap-reset:hover { color: var(--accent); border-color: var(--accent); }
#nmap-copy, #nmap-clear {
  background: var(--bg-app);
  border: 1px solid var(--border-subtle);
  border-radius: 5px;
  padding: 6px 14px;
  font-family: var(--mono);
  font-size: 11px;
  cursor: pointer;
  white-space: nowrap;
  transition: color .12s, border-color .12s;
}
#nmap-copy { color: var(--accent); }
#nmap-copy:hover { border-color: var(--accent); }
#nmap-clear { color: var(--accent-red); }
#nmap-clear:hover { border-color: var(--accent-red); }
#nmap-builder .nmap-sec { margin-bottom: 28px; }
#nmap-builder .nmap-sec-title {
  font-family: var(--mono);
  font-size: 13px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-primary);
  padding-bottom: 7px;
  margin-bottom: 2px;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 7px;
  user-select: none;
}
#nmap-builder .nmap-sec-title .nmap-chev {
  font-size: 8px;
  color: var(--text-dim);
  transition: transform .15s;
  display: inline-block;
  flex-shrink: 0;
}
#nmap-builder .nmap-sec-title.open .nmap-chev { transform: rotate(90deg); }
#nmap-builder .nmap-tbl { display: table !important; width: 100%; }
#nmap-builder .nmap-tbl.nmap-tbl-hidden { display: none !important; }
#nmap-builder .nmap-fi {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
#nmap-builder .nmap-dc { color: var(--text-muted); font-size: 13px; }
.markdown-body tr.nmap-on td,
.markdown-body tr.nmap-on:hover td { background: rgba(63,185,80,.07); }
#nmap-builder .nmap-arg {
  font-family: var(--mono);
  font-size: 12px;
  background: var(--bg-app);
  border: 1px solid var(--border-subtle);
  border-radius: 4px;
  color: #d2a8ff;
  padding: 3px 8px;
  min-width: 90px;
  outline: none;
  transition: border-color .12s;
}
#nmap-builder select.nmap-arg { min-width: unset; width: auto; cursor: pointer; }
#nmap-builder .nmap-arg:focus { border-color: var(--accent); }
#nmap-builder .nmap-use {
  font-family: var(--mono);
  font-size: 11px;
  background: transparent;
  border: 1px solid var(--border-subtle);
  border-radius: 5px;
  color: var(--accent-green);
  padding: 3px 10px;
  cursor: pointer;
  white-space: nowrap;
  flex-shrink: 0;
  transition: color .12s, border-color .12s, background .12s;
}
#nmap-builder .nmap-use:hover { border-color: var(--accent-green); }
#nmap-builder .nmap-use.on {
  background: rgba(63,185,80,.12);
  border-color: var(--accent-green);
}
#nmap-builder .nmap-use.on:hover {
  color: var(--accent-red);
  border-color: var(--accent-red);
  background: rgba(255,123,114,.12);
}
#nmap-templates {
  display: flex;
  gap: 10px;
  margin-bottom: 28px;
}
#nmap-builder .nmap-tpl {
  flex: 1;
  background: var(--bg-header);
  border: 1px solid var(--border-subtle);
  border-left: 3px solid var(--border-subtle);
  border-radius: 6px;
  padding: 11px 13px;
  cursor: pointer;
  transition: border-color .12s, background .12s;
  user-select: none;
}
#nmap-builder .nmap-tpl:hover { background: rgba(255,255,255,.02); }
#nmap-builder .nmap-tpl-evasive:hover,
#nmap-builder .nmap-tpl-evasive.nmap-tpl-on { border-left-color: var(--accent-green); }
#nmap-builder .nmap-tpl-neutral:hover,
#nmap-builder .nmap-tpl-neutral.nmap-tpl-on { border-left-color: var(--accent); }
#nmap-builder .nmap-tpl-aggressive:hover,
#nmap-builder .nmap-tpl-aggressive.nmap-tpl-on { border-left-color: var(--accent-red); }
#nmap-builder .nmap-tpl.nmap-tpl-on { background: rgba(255,255,255,.03); }
#nmap-builder .nmap-tpl-label {
  font-family: var(--mono);
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: var(--text-muted);
  margin-bottom: 3px;
}
#nmap-builder .nmap-tpl-evasive.nmap-tpl-on  .nmap-tpl-label { color: var(--accent-green); }
#nmap-builder .nmap-tpl-neutral.nmap-tpl-on  .nmap-tpl-label { color: var(--accent); }
#nmap-builder .nmap-tpl-aggressive.nmap-tpl-on .nmap-tpl-label { color: var(--accent-red); }
#nmap-builder .nmap-tpl-desc {
  font-size: 12px;
  color: var(--text-muted);
  margin-bottom: 9px;
  line-height: 1.4;
}
#nmap-builder .nmap-tpl-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}
#nmap-builder .nmap-tpl-badges code {
  font-size: 11px;
  opacity: .55;
}
#nmap-builder .nmap-tpl.nmap-tpl-on .nmap-tpl-badges code { opacity: .85; }
</style>

<div id="nmap-builder">
  <div id="nmap-cmd">
    <div id="nmap-cmd-label">Command</div>
    <div id="nmap-cmd-row">
      <div id="nmap-cmd-text" contenteditable="plaintext-only" spellcheck="false"></div>
      <button id="nmap-reset" title="Reset to builder output">↺</button>
      <button id="nmap-copy">Copy</button>
      <button id="nmap-clear">Clear</button>
    </div>
  </div>
  <div id="nmap-templates"></div>
  <div id="nmap-sections"></div>
</div>

<script>
(function () {
  const SECS = [
    { title: 'Scan Types', flags: [
      { f: '-sS', d: 'SYN scan — stealthy, no full TCP connection' },
      { f: '-sT', d: 'TCP scan — full three-way handshake' },
      { f: '-sU', d: 'UDP scan' },
      { f: '-sA', d: 'ACK scan — harder for firewalls to filter' }
    ]},
    { title: 'Ports', flags: [
      { f: '-p',          d: 'Specify ports (e.g. 22,80 or 22-443)', a: '22,80,443' },
      { f: '--top-ports', d: 'Scan N most commonly used ports',       a: '1000' },
      { f: '-p-',         d: 'Scan all 65535 ports' },
      { f: '-F',          d: 'Fast scan — top 100 ports only' }
    ]},
    { title: 'Packets', flags: [
      { f: '--packet-trace',     d: 'Show all packets sent and received' },
      { f: '-n',                 d: 'Disable DNS resolution' },
      { f: '--disable-arp-ping', d: 'Disable ARP ping discovery' },
      { f: '-Pn',                d: 'Skip host discovery (treat all hosts as up)' }
    ]},
    { title: 'Options', flags: [
      { f: '--stats-every', d: 'Print scan progress on an interval (e.g. 5s, 1m)', a: '5s' },
      { f: '-A',            d: 'OS/service detection, traceroute & default scripts' },
      { f: '-O',            d: 'OS detection' },
      { f: '--traceroute',  d: 'Trace hop path to each host' }
    ]},
    { title: 'Output', flags: [
      { f: '-oN', d: 'Save in normal .nmap format',     a: 'scan' },
      { f: '-oG', d: 'Save in greppable .gnmap format', a: 'scan' },
      { f: '-oX', d: 'Save in XML format',              a: 'scan' },
      { f: '-oA', d: 'Save in all three formats',       a: 'scan' }
    ]},
    { title: 'Scripts', flags: [
      { f: '-sC',      d: 'Run default NSE scripts' },
      { f: '--script', d: 'Run scripts by category or name', a: 'vuln' }
    ]},
    { title: 'Performance', flags: [
      { f: '--initial-rtt-timeout', d: 'Initial RTT timeout',               a: '50ms' },
      { f: '--max-rtt-timeout',     d: 'Maximum RTT timeout',               a: '100ms' },
      { f: '--max-retries',         d: 'Max probe retransmissions per port', a: '0' },
      { f: '--min-rate',            d: 'Minimum packets per second to send', a: '300' },
      { f: '-T', d: 'Timing template (0=paranoid → 5=insane)', a: '4',
        sel: ['0 — paranoid','1 — sneaky','2 — polite','3 — normal','4 — aggressive','5 — insane'],
        selV: ['0','1','2','3','4','5'], selD: 4
      }
    ]},
    { title: 'Evasion', flags: [
      { f: '-D',            d: 'Spoof with N random decoy IPs', a: 'RND:5' },
      { f: '-S',            d: 'Spoof source IP address',        a: '10.10.14.1' },
      { f: '-e',            d: 'Send via specified interface',   a: 'tun0' },
      { f: '--dns-server',  d: 'Use a specific DNS server',      a: '8.8.8.8' },
      { f: '--source-port', d: 'Spoof source port number',       a: '53' }
    ]}
  ];

  const TEMPLATES = [
    {
      key: 'evasive', label: 'Evasive', cls: 'nmap-tpl-evasive',
      desc: 'Low and slow — minimize detection footprint',
      flags: [
        { f: '-sS' }, { f: '-T', val: '1' }, { f: '-Pn' }, { f: '-n' },
        { f: '--disable-arp-ping' }, { f: '-D', val: 'RND:5' },
        { f: '--source-port', val: '53' }, { f: '--max-retries', val: '1' }, { f: '-F' }
      ]
    },
    {
      key: 'neutral', label: 'Neutral', cls: 'nmap-tpl-neutral',
      desc: 'Balanced — service & OS detection on common ports',
      flags: [
        { f: '-sS' }, { f: '-T', val: '3' }, { f: '-A' }, { f: '--top-ports', val: '1000' }
      ]
    },
    {
      key: 'aggressive', label: 'Aggressive', cls: 'nmap-tpl-aggressive',
      desc: 'Fast and loud — full port sweep with detection',
      flags: [
        { f: '-sS' }, { f: '-T', val: '4' }, { f: '-A' },
        { f: '-p-' }, { f: '--min-rate', val: '300' }
      ]
    }
  ];

  const st = new Map();
  const g   = id => document.getElementById(id);
  const fid = f  => 'nb' + f.replace(/[^a-zA-Z0-9]/g, '_');
  let dirty = false;

  function activateFlag(f, val) {
    let fd = null;
    for (const sec of SECS) { fd = sec.flags.find(x => x.f === f); if (fd) break; }
    if (!fd) return;
    const v = val !== undefined ? String(val) : (fd.sel ? fd.selV[fd.selD] : '');
    st.set(f, { on: true, val: v });
    g(fid(f) + 'r')?.classList.add('nmap-on');
    const btn = g(fid(f) + 'b');
    if (btn) { btn.textContent = 'Remove'; btn.classList.add('on'); }
    const arg = g(fid(f) + 'a');
    if (arg) {
      arg.style.display = '';
      if (arg.tagName === 'INPUT')  arg.value = v;
      if (arg.tagName === 'SELECT') arg.value = v;
    }
  }

  function applyTemplate(tpl, cardEl) {
    clearAll();
    tpl.flags.forEach(({ f, val }) => activateFlag(f, val));
    document.querySelectorAll('.nmap-tpl').forEach(c => c.classList.remove('nmap-tpl-on'));
    cardEl.classList.add('nmap-tpl-on');
    buildCmd();
  }

  function setDirty(val) {
    dirty = val;
    g('nmap-reset').style.display = val ? 'inline-block' : 'none';
  }

  function buildCmd() {
    if (dirty) return;
    const tgt = document.getElementById('input_TARGET_IP')?.value?.trim() || '<TARGET_IP>';
    const parts = ['sudo nmap'];
    SECS.forEach(sec => sec.flags.forEach(fd => {
      const s = st.get(fd.f);
      if (!s?.on) return;
      if (fd.a !== undefined) {
        const v = (s.val || '').trim();
        parts.push(v ? `${fd.f} ${v}` : fd.f);
      } else {
        parts.push(fd.f);
      }
    }));
    parts.push(tgt);
    g('nmap-cmd-text').textContent = parts.join(' ');
  }

  function toggle(fd) {
    const cur  = st.get(fd.f) ?? { on: false, val: fd.sel ? fd.selV[fd.selD] : '' };
    const next = !cur.on;
    st.set(fd.f, { on: next, val: cur.val });

    g(fid(fd.f) + 'r')?.classList.toggle('nmap-on', next);

    const btn = g(fid(fd.f) + 'b');
    if (btn) { btn.textContent = next ? 'Remove' : 'Use'; btn.classList.toggle('on', next); }

    const arg = g(fid(fd.f) + 'a');
    if (arg) {
      arg.style.display = next ? '' : 'none';
      if (next && arg.tagName === 'INPUT') { arg.value = cur.val || ''; arg.focus(); }
    }
    buildCmd();
  }

  function clearAll() {
    setDirty(false);
    document.querySelectorAll('.nmap-tpl').forEach(c => c.classList.remove('nmap-tpl-on'));
    SECS.forEach(sec => sec.flags.forEach(fd => {
      st.set(fd.f, { on: false, val: fd.sel ? fd.selV[fd.selD] : '' });
      g(fid(fd.f) + 'r')?.classList.remove('nmap-on');
      const btn = g(fid(fd.f) + 'b');
      if (btn) { btn.textContent = 'Use'; btn.classList.remove('on'); }
      const arg = g(fid(fd.f) + 'a');
      if (arg) {
        arg.style.display = 'none';
        if (arg.tagName === 'INPUT')  arg.value = '';
        if (arg.tagName === 'SELECT') arg.selectedIndex = fd.selD ?? 0;
      }
    }));
    buildCmd();
  }

  // Build template cards
  const tplWrap = g('nmap-templates');
  TEMPLATES.forEach(tpl => {
    const card = document.createElement('div');
    card.className = `nmap-tpl ${tpl.cls}`;

    const label = document.createElement('div');
    label.className = 'nmap-tpl-label';
    label.textContent = tpl.label;
    card.appendChild(label);

    const desc = document.createElement('div');
    desc.className = 'nmap-tpl-desc';
    desc.textContent = tpl.desc;
    card.appendChild(desc);

    const badges = document.createElement('div');
    badges.className = 'nmap-tpl-badges';
    tpl.flags.forEach(({ f, val }) => {
      const badge = document.createElement('code');
      badge.textContent = val !== undefined ? `${f} ${val}` : f;
      badges.appendChild(badge);
    });
    card.appendChild(badges);

    card.addEventListener('click', () => applyTemplate(tpl, card));
    tplWrap.appendChild(card);
  });

  // Build tables
  const wrap = g('nmap-sections');
  SECS.forEach(sec => {
    const div = document.createElement('div');
    div.className = 'nmap-sec';

    const secTitle = document.createElement('div');
    secTitle.className = 'nmap-sec-title';
    const chev = document.createElement('span');
    chev.className = 'nmap-chev';
    chev.textContent = '▶';
    secTitle.appendChild(chev);
    secTitle.appendChild(document.createTextNode(sec.title));
    div.appendChild(secTitle);

    const tbl = document.createElement('table');
    tbl.className = 'nmap-tbl';
    tbl.innerHTML = '<thead><tr><th>Flag</th><th>Usage</th></tr></thead>';
    const tbody = document.createElement('tbody');

    sec.flags.forEach(fd => {
      const tr = document.createElement('tr');
      tr.id = fid(fd.f) + 'r';

      // Flag cell
      const td1 = document.createElement('td');
      const fi  = document.createElement('div');
      fi.className = 'nmap-fi';

      const code = document.createElement('code');
      code.textContent = fd.f;
      fi.appendChild(code);

      // Argument input / select
      if (fd.a !== undefined) {
        let argEl;
        if (fd.sel) {
          argEl = document.createElement('select');
          argEl.className = 'nmap-arg';
          fd.sel.forEach((lbl, i) => {
            argEl.appendChild(new Option(lbl, fd.selV[i], i === fd.selD, i === fd.selD));
          });
          argEl.addEventListener('change', () => {
            const s = st.get(fd.f) ?? { on: true, val: '' };
            st.set(fd.f, { on: s.on, val: argEl.value });
            buildCmd();
          });
          st.set(fd.f, { on: false, val: fd.selV[fd.selD] });
        } else {
          argEl = document.createElement('input');
          argEl.type = 'text';
          argEl.className = 'nmap-arg';
          argEl.placeholder = fd.a;
          argEl.spellcheck = false;
          argEl.addEventListener('input', () => {
            const s = st.get(fd.f) ?? { on: true, val: '' };
            st.set(fd.f, { on: s.on, val: argEl.value });
            buildCmd();
          });
        }
        argEl.id = fid(fd.f) + 'a';
        argEl.style.display = 'none';
        fi.appendChild(argEl);
      }

      // Use / Remove button
      const btn = document.createElement('button');
      btn.className = 'nmap-use';
      btn.id = fid(fd.f) + 'b';
      btn.textContent = 'Use';
      btn.addEventListener('click', () => toggle(fd));
      fi.appendChild(btn);

      td1.appendChild(fi);
      tr.appendChild(td1);

      // Description cell
      const td2 = document.createElement('td');
      td2.className = 'nmap-dc';
      td2.textContent = fd.d;
      tr.appendChild(td2);

      tbody.appendChild(tr);
    });

    tbl.appendChild(tbody);
    tbl.classList.add('nmap-tbl-hidden');
    secTitle.addEventListener('click', () => {
      const open = secTitle.classList.toggle('open');
      tbl.classList.toggle('nmap-tbl-hidden', !open);
    });
    div.appendChild(tbl);
    wrap.appendChild(div);
  });

  // Manual edits to command block → dirty mode
  g('nmap-cmd-text').addEventListener('input', () => setDirty(true));

  // Live-sync with global Target IP in the header
  document.getElementById('input_TARGET_IP')?.addEventListener('input', () => buildCmd());

  // Reset button — exit dirty mode and rebuild
  g('nmap-reset').addEventListener('click', () => { setDirty(false); buildCmd(); });

  // Copy button
  g('nmap-copy').addEventListener('click', () => {
    navigator.clipboard.writeText(g('nmap-cmd-text').textContent).then(() => {
      const btn = g('nmap-copy');
      btn.textContent = 'Copied!';
      btn.style.cssText = 'color:var(--accent-green);border-color:var(--accent-green)';
      setTimeout(() => { btn.textContent = 'Copy'; btn.style.cssText = ''; }, 1500);
    });
  });

  g('nmap-clear').addEventListener('click', clearAll);
  buildCmd();
})();
</script>
