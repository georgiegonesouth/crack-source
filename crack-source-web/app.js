/* ── State ─────────────────────────────────────────────── */
const TOKENS = ['TARGET_IP', 'PORT', 'LHOST', 'LPORT'];

// Regex per token that matches &lt;TOKEN&gt; even when hljs wraps < and >
// in separate spans (common for bash/shell redirect syntax).
// Pattern:  optional-opening-span  &lt;  any-span-boundaries  TOKEN  any-span-boundaries  &gt;  optional-closing-span
const SPAN_BOUNDARY = '(?:(?:<\\/span>|<span[^>]*>)\\s*)*';
const TOKEN_RE = Object.fromEntries(TOKENS.map(t => {
  const esc = t.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  return [t, new RegExp(
    `(?:<span[^>]*>)?&lt;${SPAN_BOUNDARY}${esc}${SPAN_BOUNDARY}&gt;(?:<\\/span>)?`,
    'g'
  )];
}));

function buildTokenRe(token) {
  const esc = token.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  return new RegExp(
    `(?:<span[^>]*>)?&lt;${SPAN_BOUNDARY}${esc}${SPAN_BOUNDARY}&gt;(?:<\\/span>)?`,
    'g'
  );
}

const state = {
  inputs: Object.fromEntries(TOKENS.map(t => [t, ''])),
  localInputs: {},
  localTokenRe: {},
  localTokenOptions: {},
  expandedPaths: new Set(),
  currentPath: null,
};

let CONFIG = { MANIFEST_PATH: 'manifest.json', REPO_ROOT: '' };

/* ── DOM refs ──────────────────────────────────────────── */
const $content = document.getElementById('content');
const $navWrap = document.getElementById('nav-tree-wrap');
const $search  = document.getElementById('search');

/* ── Search state ──────────────────────────────────────── */
const contentIndex = new Map(); // path → lowercased text
let searchMatches  = [];
let searchFocusIdx = -1;

async function loadConfig() {
  try {
    const text = await fetch('web.config').then(r => r.ok ? r.text() : '');
    text.split('\n').forEach(line => {
      const eq = line.indexOf('=');
      if (eq > 0) CONFIG[line.slice(0, eq).trim()] = line.slice(eq + 1).trim();
    });
  } catch {}
}

/* ── Init ──────────────────────────────────────────────── */
document.addEventListener('DOMContentLoaded', async () => {
  loadInputsFromStorage();
  restoreInputFields();
  setupInputListeners();

  await loadConfig();

  const manifest = await fetch(contentUrl(CONFIG.MANIFEST_PATH))
    .then(r => { if (!r.ok) throw new Error(r.status); return r.json(); })
    .catch(() => {
      $navWrap.innerHTML = `<p class="error-msg" style="padding:16px">
        manifest.json not found.<br>Run: <code>python3 generate_manifest.py</code>
      </p>`;
      return null;
    });

  if (!manifest) return;
  buildNav(manifest.tree);
  setupSearch();
  setupLibrarySearch();
  buildContentIndex(manifest.tree);
  showWelcome();
});

/* ── LocalStorage ──────────────────────────────────────── */
function loadInputsFromStorage() {
  TOKENS.forEach(t => { state.inputs[t] = localStorage.getItem(`pt_token_${t}`) || ''; });
}

function saveInputsToStorage() {
  TOKENS.forEach(t => { localStorage.setItem(`pt_token_${t}`, state.inputs[t]); });
}

function restoreInputFields() {
  TOKENS.forEach(t => {
    const el = document.getElementById(`input_${t}`);
    if (el) el.value = state.inputs[t];
  });
}

function setupInputListeners() {
  TOKENS.forEach(t => {
    const el = document.getElementById(`input_${t}`);
    if (!el) return;
    el.addEventListener('input', () => {
      state.inputs[t] = el.value.trim();
      saveInputsToStorage();
      applySubstitution();
    });
  });
}

/* ── Nav building ──────────────────────────────────────── */
function buildNav(tree) {
  const sidebarNodes = tree.filter(n => n.name !== 'LIBRARY');
  const libraryNode  = tree.find(n => n.name === 'LIBRARY');

  const ul = buildNavList(sidebarNodes, 0);
  ul.className = 'nav-tree';
  $navWrap.appendChild(ul);

  if (libraryNode) {
    const libWrap = document.getElementById('library-nav-wrap');
    const libUl = buildNavList(libraryNode.children || [], 0);
    libUl.className = 'nav-tree';
    libWrap.appendChild(libUl);
  }

  const $libraryBtn   = document.getElementById('library-btn');
  const $libraryPanel = document.getElementById('library-panel');
  const $libraryClose = document.getElementById('library-close');

  const $libSearch = document.getElementById('library-search');

  function openLibrary() {
    $libraryPanel.hidden = false;
    $libraryBtn.classList.add('active');
    setTimeout(() => $libSearch.focus(), 0);
  }
  function closeLibrary() {
    $libraryPanel.hidden = true;
    $libraryBtn.classList.remove('active');
  }

  $libraryBtn.addEventListener('click', () => {
    $libraryPanel.hidden ? openLibrary() : closeLibrary();
  });

  $libraryClose.addEventListener('click', closeLibrary);

  document.addEventListener('keydown', e => {
    if (e.key === 'Escape' && !$libraryPanel.hidden) closeLibrary();
  });
}

function buildNavList(nodes, depth) {
  const ul = document.createElement('ul');
  ul.className = depth === 0 ? 'nav-tree' : 'nav-children';
  nodes.forEach(node => ul.appendChild(buildNavNode(node, depth)));
  return ul;
}

function buildNavNode(node, depth) {
  const li = document.createElement('li');
  li.dataset.path = node.path;
  li.dataset.nameLower = node.name.trim().toLowerCase();

  if (node.type === 'dir') {
    li.className = 'nav-dir';
    const children = node.children || [];

    if (children.length === 0) {
      const label = document.createElement('span');
      label.className = 'nav-dir-empty';
      label.textContent = node.name;
      li.appendChild(label);
    } else {
      const btn = document.createElement('button');
      btn.className = 'nav-dir-btn';
      btn.setAttribute('aria-expanded', 'false');
      btn.innerHTML = `<span class="chevron">▶</span><span>${node.name}</span>`;
      li.appendChild(btn);

      const childUl = buildNavList(children, depth + 1);
      childUl.hidden = true;
      li.appendChild(childUl);

      btn.addEventListener('click', () => {
        const open = btn.getAttribute('aria-expanded') === 'true';
        btn.setAttribute('aria-expanded', String(!open));
        childUl.hidden = open;
        if (!open) state.expandedPaths.add(node.path);
        else        state.expandedPaths.delete(node.path);
      });
    }
  } else {
    li.className = 'nav-file';
    const a = document.createElement('a');
    a.href = '#';
    a.textContent = displayName(node.name);
    a.dataset.path = node.path;
    a.addEventListener('click', e => { e.preventDefault(); loadFile(node.path); });
    li.appendChild(a);
  }

  return li;
}

function displayName(filename) {
  return filename.trim().replace(/\.(md|pdf)$/i, '');
}

/* ── Content index ─────────────────────────────────────── */
async function buildContentIndex(tree) {
  const paths = [];
  function collect(nodes) {
    nodes.forEach(n => {
      if (n.type === 'file' && n.path.toLowerCase().endsWith('.md')) paths.push(n.path);
      if (n.children) collect(n.children);
    });
  }
  collect(tree);
  await Promise.all(paths.map(async p => {
    try {
      const r = await fetch(contentUrl(p), { cache: 'no-store' });
      if (r.ok) contentIndex.set(p, (await r.text()).toLowerCase());
    } catch {}
  }));
}

/* ── Search ────────────────────────────────────────────── */
function isEditing() {
  const el = document.activeElement;
  return el && (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA' || el.contentEditable === 'true' || el.contentEditable === 'plaintext-only');
}

function setSearchFocus(idx) {
  searchMatches.forEach((li, i) => li.classList.toggle('search-focused', i === idx));
  if (idx >= 0 && idx < searchMatches.length) {
    searchMatches[idx].scrollIntoView({ block: 'nearest' });
  }
  searchFocusIdx = idx;
}

function setupSearch() {
  $search.addEventListener('input', () => filterNav($search.value.trim().toLowerCase()));

  $search.addEventListener('keydown', e => {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setSearchFocus(Math.min(searchFocusIdx + 1, searchMatches.length - 1));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setSearchFocus(Math.max(searchFocusIdx - 1, 0));
    } else if (e.key === 'Enter') {
      e.preventDefault();
      if (searchFocusIdx >= 0 && searchMatches[searchFocusIdx]) {
        const a = searchMatches[searchFocusIdx].querySelector('a');
        if (a) loadFile(a.dataset.path);
      }
      $search.blur();
    } else if (e.key === 'Escape') {
      $search.blur();
    }
  });

  document.addEventListener('keydown', e => {
    if (e.key === '/' && !e.metaKey && !e.ctrlKey && document.activeElement !== $search && !isEditing()) {
      e.preventDefault();
      $search.focus();
      $search.select();
    }
  });
}

function filterNav(q) {
  searchMatches = [];
  searchFocusIdx = -1;

  if (!q) {
    $navWrap.classList.remove('searching');
    $navWrap.querySelectorAll('.no-match').forEach(el => el.classList.remove('no-match'));
    $navWrap.querySelectorAll('.search-focused').forEach(el => el.classList.remove('search-focused'));
    return;
  }

  $navWrap.classList.add('searching');
  const matchedDirPaths = new Set();

  $navWrap.querySelectorAll('li.nav-file').forEach(li => {
    const nameMatch    = (li.dataset.nameLower || '').includes(q);
    const contentMatch = contentIndex.has(li.dataset.path) && contentIndex.get(li.dataset.path).includes(q);
    const match = nameMatch || contentMatch;
    li.classList.toggle('no-match', !match);
    li.classList.remove('search-focused');
    if (match) {
      searchMatches.push(li);
      let node = li.parentElement;
      while (node && node !== $navWrap) {
        if (node.tagName === 'LI' && node.dataset.path) matchedDirPaths.add(node.dataset.path);
        node = node.parentElement;
      }
    }
  });

  $navWrap.querySelectorAll('li.nav-dir').forEach(li => {
    li.classList.toggle('no-match', !matchedDirPaths.has(li.dataset.path));
  });
}

/* ── Library search ────────────────────────────────────── */
let libSearchMatches  = [];
let libSearchFocusIdx = -1;

function setLibSearchFocus(idx) {
  libSearchMatches.forEach((li, i) => li.classList.toggle('search-focused', i === idx));
  if (idx >= 0 && idx < libSearchMatches.length) {
    libSearchMatches[idx].scrollIntoView({ block: 'nearest' });
  }
  libSearchFocusIdx = idx;
}

function setupLibrarySearch() {
  const $libSearch = document.getElementById('library-search');
  const $libWrap   = document.getElementById('library-nav-wrap');
  if (!$libSearch || !$libWrap) return;

  $libSearch.addEventListener('input', () => filterLibraryNav($libSearch.value.trim().toLowerCase()));

  $libSearch.addEventListener('keydown', e => {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setLibSearchFocus(Math.min(libSearchFocusIdx + 1, libSearchMatches.length - 1));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setLibSearchFocus(Math.max(libSearchFocusIdx - 1, 0));
    } else if (e.key === 'Enter') {
      e.preventDefault();
      if (libSearchFocusIdx >= 0 && libSearchMatches[libSearchFocusIdx]) {
        const a = libSearchMatches[libSearchFocusIdx].querySelector('a');
        if (a) loadFile(a.dataset.path);
      }
      $libSearch.blur();
    } else if (e.key === 'Escape') {
      $libSearch.blur();
    }
  });

  document.addEventListener('keydown', e => {
    if (e.key === '/' && (e.metaKey || e.ctrlKey) && !isEditing()) {
      e.preventDefault();
      const $libraryPanel = document.getElementById('library-panel');
      const $libraryBtn   = document.getElementById('library-btn');
      if ($libraryPanel.hidden) {
        $libraryPanel.hidden = false;
        $libraryBtn.classList.add('active');
      }
      setTimeout(() => { $libSearch.focus(); $libSearch.select(); }, 0);
    }
  });
}

function filterLibraryNav(q) {
  const $libWrap = document.getElementById('library-nav-wrap');
  if (!$libWrap) return;
  libSearchMatches  = [];
  libSearchFocusIdx = -1;

  if (!q) {
    $libWrap.classList.remove('searching');
    $libWrap.querySelectorAll('.no-match').forEach(el => el.classList.remove('no-match'));
    $libWrap.querySelectorAll('.search-focused').forEach(el => el.classList.remove('search-focused'));
    return;
  }

  $libWrap.classList.add('searching');
  const matchedDirPaths = new Set();

  $libWrap.querySelectorAll('li.nav-file').forEach(li => {
    const nameMatch    = (li.dataset.nameLower || '').includes(q);
    const contentMatch = contentIndex.has(li.dataset.path) && contentIndex.get(li.dataset.path).includes(q);
    const match = nameMatch || contentMatch;
    li.classList.toggle('no-match', !match);
    li.classList.remove('search-focused');
    if (match) {
      libSearchMatches.push(li);
      let node = li.parentElement;
      while (node && node !== $libWrap) {
        if (node.tagName === 'LI' && node.dataset.path) matchedDirPaths.add(node.dataset.path);
        node = node.parentElement;
      }
    }
  });

  $libWrap.querySelectorAll('li.nav-dir').forEach(li => {
    li.classList.toggle('no-match', !matchedDirPaths.has(li.dataset.path));
  });
}

/* ── File loading ──────────────────────────────────────── */
function encodePath(path) {
  return path.split('/').map(encodeURIComponent).join('/');
}

function contentUrl(path) {
  const prefix = CONFIG.REPO_ROOT ? CONFIG.REPO_ROOT + '/' : '';
  return (prefix + path).split('/').map(encodeURIComponent).join('/');
}

function setActiveLink(path) {
  const $libWrap = document.getElementById('library-nav-wrap');
  const roots = [$navWrap, $libWrap].filter(Boolean);
  roots.forEach(root => root.querySelectorAll('a.active').forEach(a => a.classList.remove('active')));
  roots.forEach(root => {
    root.querySelectorAll('li.nav-file a').forEach(a => {
      if (a.dataset.path !== path) return;
      a.classList.add('active');
      let node = a.parentElement?.parentElement?.parentElement;
      while (node && node !== root) {
        if (node.tagName === 'LI' && node.classList.contains('nav-dir')) {
          const btn = node.querySelector(':scope > .nav-dir-btn');
          const ul  = node.querySelector(':scope > ul.nav-children');
          if (btn && ul && ul.hidden) {
            ul.hidden = false;
            btn.setAttribute('aria-expanded', 'true');
          }
        }
        node = node.parentElement;
      }
    });
  });
}

async function loadFile(path) {
  state.currentPath = path;
  setActiveLink(path);
  $content.innerHTML = '<div class="loading">Loading…</div>';

  const ext = path.split('.').pop().toLowerCase();
  const url  = contentUrl(path);

  if (ext === 'pdf') {
    $content.innerHTML = `<iframe class="pdf-frame" src="${url}" title="PDF viewer"></iframe>`;
    return;
  }

  try {
    const res = await fetch(url, { cache: 'no-store' });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    renderMarkdown(await res.text());
  } catch (err) {
    $content.innerHTML = `<div class="error-msg">Failed to load: ${err.message}</div>`;
  }
}

/* ── Markdown rendering ────────────────────────────────── */
function stripFrontMatter(text) {
  return text.replace(/^---[\s\S]*?---\n?/, '');
}

function renderMarkdown(text) {
  text = stripFrontMatter(text);
  $content.classList.remove('welcome-mode');
  marked.setOptions({ gfm: true, breaks: false });

  setupLocalTokens(text);

  const wrapper = document.createElement('div');
  wrapper.className = 'markdown-body';
  wrapper.innerHTML = marked.parse(text);

  wrapper.querySelectorAll('a[href]').forEach(a => {
    const href = a.getAttribute('href');
    if (href && href.endsWith('.md')) {
      a.addEventListener('click', e => {
        e.preventDefault();
        let target = href;
        if (state.currentPath && !href.startsWith('CONTENT/')) {
          const dir = state.currentPath.substring(0, state.currentPath.lastIndexOf('/') + 1);
          target = dir + href;
        }
        loadFile(target);
      });
    } else if (href && (href.startsWith('http://') || href.startsWith('https://'))) {
      a.target = '_blank';
      a.rel = 'noopener noreferrer';
    }
  });

  // Block code: highlight → store highlighted HTML → substitute
  wrapper.querySelectorAll('pre code').forEach(codeEl => {
    codeEl.dataset.origText = codeEl.textContent; // plain text for copy button
    hljs.highlightElement(codeEl);
    codeEl.dataset.origHtml = codeEl.innerHTML;   // highlighted HTML (may have split spans)
    codeEl.innerHTML = substituteHtml(codeEl.dataset.origHtml);
    makeEditable(codeEl);
    addCodeActions(codeEl.closest('pre'), codeEl);
  });

  // Inline code: marked already escaped < > so origHtml has &lt;TOKEN&gt;
  wrapper.querySelectorAll(':not(pre) > code').forEach(codeEl => {
    codeEl.dataset.origHtml = codeEl.innerHTML;
    codeEl.innerHTML = substituteHtml(codeEl.dataset.origHtml);
  });

  addOptionUseButtons(wrapper);
  makeCommandTablesInteractive(wrapper);
  const aside = extractTipsBox(wrapper);
  renderLocalInputBar(wrapper);
  makeH2Collapsible(wrapper);
  addTokenTableButtons(wrapper);
  addSectionUseButtons(wrapper);

  $content.innerHTML = '';
  $content.appendChild(wrapper);
  if (aside) $content.appendChild(aside);

  // Execute scripts embedded in rendered markdown (innerHTML doesn't run them)
  wrapper.querySelectorAll('script').forEach(orig => {
    const s = document.createElement('script');
    s.textContent = orig.textContent;
    document.head.appendChild(s).remove();
  });
}

/* ── Command tables (editable + copyable cells) ────────── */
function makeCommandTablesInteractive(wrapper) {
  wrapper.querySelectorAll('table').forEach(table => {
    const headers = [...table.querySelectorAll('th')];
    const cmdIdx = headers.findIndex(th => th.textContent.trim().toUpperCase() === 'COMMAND');
    if (cmdIdx === -1) return;

    table.querySelectorAll('tbody tr').forEach(tr => {
      const tds = tr.querySelectorAll('td');
      if (!tds[cmdIdx]) return;
      const td = tds[cmdIdx];
      tds.forEach((cell, i) => { if (i !== cmdIdx) cell.style.whiteSpace = 'nowrap'; });

      const cmdText = td.textContent.trim();
      td.innerHTML = '';
      td.classList.add('cmd-cell');

      const code = document.createElement('code');
      code.className = 'cmd-text';
      code.contentEditable = 'plaintext-only';
      code.spellcheck = false;
      code.textContent = cmdText;

      // Detect lowercase placeholders like <database>, <table>
      const placeholders = [...new Set(
        [...cmdText.matchAll(/<([a-z][a-z0-9_]*)>/g)].map(m => m[1])
      )];
      const fillValues = Object.fromEntries(placeholders.map(p => [p, '']));

      function applyFill() {
        let text = cmdText;
        placeholders.forEach(p => { if (fillValues[p]) text = text.replaceAll(`<${p}>`, fillValues[p]); });
        TOKENS.forEach(t => { const v = state.inputs[t]; if (v) text = text.replaceAll(`<${t}>`, v); });
        Object.entries(state.localInputs || {}).forEach(([t, v]) => { if (v) text = text.replaceAll(`<${t}>`, v); });
        code.textContent = text;
      }
      code._cmdFill = applyFill;

      const resetBtn = document.createElement('button');
      resetBtn.className = 'reset-btn cmd-reset-btn';
      resetBtn.title = 'Reset to original';
      resetBtn.textContent = '↺';

      const copyBtn = document.createElement('button');
      copyBtn.className = 'copy-btn cmd-copy-btn';
      copyBtn.textContent = 'Copy';
      copyBtn.addEventListener('click', () => {
        navigator.clipboard.writeText(code.textContent.trim()).then(() => {
          copyBtn.textContent = 'Copied!';
          copyBtn.classList.add('copied');
          setTimeout(() => { copyBtn.textContent = 'Copy'; copyBtn.classList.remove('copied'); }, 1500);
        });
      });

      const inner = document.createElement('div');
      inner.className = 'cmd-cell-inner';
      inner.appendChild(code);

      // Fill form for cells with placeholders
      let fillForm = null;
      let fillRow = null;
      if (placeholders.length) {
        fillForm = document.createElement('div');
        fillForm.className = 'cmd-fill-form';

        placeholders.forEach(p => {
          const group = document.createElement('div');
          group.className = 'cmd-fill-group';
          const label = document.createElement('label');
          label.textContent = p;
          const input = document.createElement('input');
          input.type = 'text';
          input.placeholder = `<${p}>`;
          input.addEventListener('input', () => { fillValues[p] = input.value; applyFill(); });
          input.addEventListener('keydown', e => {
            if (e.key !== 'Enter') return;
            e.preventDefault();
            const inputs = [...fillForm.querySelectorAll('input')];
            const next = inputs[inputs.indexOf(input) + 1];
            if (next) next.focus();
          });
          group.appendChild(label);
          group.appendChild(input);
          fillForm.appendChild(group);
        });

        const fillBtn = document.createElement('button');
        fillBtn.className = 'reset-btn cmd-fill-btn';
        fillBtn.textContent = 'Fill';
        fillBtn.addEventListener('click', () => {
          if (fillRow) {
            fillRow.classList.remove('open');
            fillRow.addEventListener('transitionend', () => { fillRow.remove(); fillRow = null; }, { once: true });
            fillBtn.classList.remove('active');
          } else {
            fillRow = document.createElement('tr');
            fillRow.className = 'cmd-fill-row';
            const fillTd = document.createElement('td');
            fillTd.colSpan = headers.length;
            fillTd.appendChild(fillForm);
            fillRow.appendChild(fillTd);
            tr.insertAdjacentElement('afterend', fillRow);
            requestAnimationFrame(() => {
              fillRow.classList.add('open');
              fillForm.querySelector('input')?.focus();
            });
            fillBtn.classList.add('active');
          }
        });
        inner.appendChild(fillBtn);
      }

      resetBtn.addEventListener('click', () => {
        code.textContent = cmdText;
        placeholders.forEach(p => { fillValues[p] = ''; });
        if (fillForm) fillForm.querySelectorAll('input').forEach(i => i.value = '');
      });

      inner.appendChild(resetBtn);
      inner.appendChild(copyBtn);
      td.appendChild(inner);
    });
  });
}

/* ── Page-local tokens ─────────────────────────────────── */
function setupLocalTokens(rawText) {
  const re = /<([A-Z][A-Z0-9_]*)>/g;
  const found = new Set();
  let m;
  while ((m = re.exec(rawText)) !== null) {
    if (!TOKENS.includes(m[1])) found.add(m[1]);
  }

  const options = {};
  const optRe = /<!--\s*token-options:([A-Z][A-Z0-9_]*)\n([\s\S]*?)-->/g;
  let om;
  while ((om = optRe.exec(rawText)) !== null) {
    options[om[1]] = om[2].split('\n').map(l => l.trim()).filter(Boolean);
  }

  const excludeMatch = rawText.match(/<!--\s*token-exclude:(.*?)-->/);
  const excluded = excludeMatch
    ? excludeMatch[1].split(',').map(s => s.trim()).filter(Boolean)
    : [];

  excluded.forEach(t => found.delete(t));

  const tableToks = new Set();
  const tableRe2 = /<!--\s*token-table:([A-Z][A-Z0-9_]*)\s*-->/g;
  let tm;
  while ((tm = tableRe2.exec(rawText)) !== null) tableToks.add(tm[1]);
  state.localTokenTables = tableToks;

  const sectionToks = new Set();
  const sectionRe = /<!--\s*token-section:([A-Z][A-Z0-9_]*)\s*-->/g;
  let sm;
  while ((sm = sectionRe.exec(rawText)) !== null) sectionToks.add(sm[1]);
  state.localTokenSections = sectionToks;

  state.localInputs = Object.fromEntries([...found].map(t => [t, '']));
  state.localTokenRe = Object.fromEntries([...found].map(t => [t, buildTokenRe(t)]));
  state.localTokenOptions = options;
}

function addOptionUseButtons(wrapper) {
  Object.entries(state.localTokenOptions).forEach(([token, opts]) => {
    wrapper.querySelectorAll('pre code').forEach(codeEl => {
      const raw = (codeEl.dataset.origText || codeEl.textContent)
        .split('\n').map(l => l.startsWith('- ') ? l.slice(2) : l).join('\n').trim();
      if (!opts.includes(raw)) return;

      const useBtn = document.createElement('button');
      useBtn.className = 'use-btn';
      useBtn.dataset.token = token;
      useBtn.dataset.value = raw;
      useBtn.textContent = 'Use';
      useBtn.addEventListener('click', () => {
        const alreadyActive = state.localInputs[token] === raw;
        state.localInputs[token] = alreadyActive ? '' : raw;
        wrapper.querySelectorAll(`.use-btn[data-token="${token}"]`).forEach(b => {
          b.classList.toggle('active', !alreadyActive && b === useBtn);
        });
        applySubstitution();
      });

      const actions = codeEl.closest('pre').querySelector('.code-actions');
      if (actions) actions.insertBefore(useBtn, actions.firstChild);
    });
  });
}

function addTokenTableButtons(wrapper) {
  const walker = document.createTreeWalker(wrapper, NodeFilter.SHOW_COMMENT);
  let node;
  while ((node = walker.nextNode())) {
    const m = node.nodeValue.trim().match(/^token-table:([A-Z][A-Z0-9_]*)$/);
    if (!m) continue;
    const token = m[1];
    // makeH2Collapsible uses nextElementSibling, which skips comment nodes,
    // so the table may end up in a wrapper div rather than as a direct sibling.
    // Search the comment's parent for the nearest table instead.
    const table = node.parentElement?.querySelector('table');
    if (!table) continue;
    table.querySelectorAll('td').forEach(td => {
      const value = td.textContent.trim();
      if (!value) return;
      const useBtn = document.createElement('button');
      useBtn.className = 'use-btn';
      useBtn.dataset.token = token;
      useBtn.dataset.value = value;
      useBtn.classList.toggle('active', state.localInputs[token] === value);
      useBtn.textContent = 'Use';
      useBtn.addEventListener('click', () => {
        const alreadyActive = state.localInputs[token] === value;
        state.localInputs[token] = alreadyActive ? '' : value;
        wrapper.querySelectorAll(`.use-btn[data-token="${token}"]`).forEach(b => {
          b.classList.toggle('active', !alreadyActive && b === useBtn);
        });
        applySubstitution();
      });
      td.appendChild(useBtn);
    });
  }
}

function addSectionUseButtons(wrapper) {
  const walker = document.createTreeWalker(wrapper, NodeFilter.SHOW_COMMENT);
  let node;
  while ((node = walker.nextNode())) {
    const m = node.nodeValue.trim().match(/^token-section:([A-Z][A-Z0-9_]*)$/);
    if (!m) continue;
    const token = m[1];
    node.parentElement?.querySelectorAll('pre code').forEach(codeEl => {
      const value = (codeEl.dataset.origText || codeEl.textContent)
        .split('\n').map(l => l.startsWith('- ') ? l.slice(2) : l).join('\n').trim();
      if (!value) return;
      const useBtn = document.createElement('button');
      useBtn.className = 'use-btn';
      useBtn.dataset.token = token;
      useBtn.dataset.value = value;
      useBtn.classList.toggle('active', state.localInputs[token] === value);
      useBtn.textContent = 'Use';
      useBtn.addEventListener('click', () => {
        const alreadyActive = state.localInputs[token] === value;
        state.localInputs[token] = alreadyActive ? '' : value;
        wrapper.querySelectorAll(`.use-btn[data-token="${token}"]`).forEach(b => {
          b.classList.toggle('active', !alreadyActive && b === useBtn);
        });
        applySubstitution();
      });
      const actions = codeEl.closest('pre').querySelector('.code-actions');
      if (actions) actions.insertBefore(useBtn, actions.firstChild);
    });
  }
}

function renderLocalInputBar(wrapper) {
  const tokens = Object.keys(state.localInputs);
  if (tokens.length === 0) return;

  const bar = document.createElement('div');
  bar.className = 'local-token-bar';

  tokens
    .filter(token => !(state.localTokenOptions[token]?.length) && !state.localTokenTables?.has(token) && !state.localTokenSections?.has(token))
    .forEach(token => {
      const group = document.createElement('div');
      group.className = 'input-group';

      const label = document.createElement('label');
      label.textContent = token.replace(/_/g, ' ');

      const input = document.createElement('input');
      input.type = 'text';
      input.placeholder = `<${token}>`;
      input.value = state.localInputs[token];
      input.addEventListener('input', () => {
        state.localInputs[token] = input.value.trim();
        applySubstitution();
      });

      group.appendChild(label);
      group.appendChild(input);
      bar.appendChild(group);
    });

  wrapper.insertBefore(bar, wrapper.firstChild);
}

/* ── Tips box ──────────────────────────────────────────── */
function extractTipsBox(wrapper) {
  const tipsH2 = [...wrapper.querySelectorAll('h2')].find(
    h2 => h2.textContent.trim().toLowerCase() === 'tips'
  );
  if (!tipsH2) return null;

  const aside = document.createElement('aside');
  aside.className = 'tips-box';

  const title = document.createElement('div');
  title.className = 'tips-title';
  title.textContent = 'Tips & Facts';
  aside.appendChild(title);

  let next = tipsH2.nextElementSibling;
  while (next && next.tagName !== 'H2') {
    const after = next.nextElementSibling;
    aside.appendChild(next);
    next = after;
  }

  tipsH2.remove();
  return aside;
}

/* ── Collapsible h2 sections ───────────────────────────── */
function makeH2Collapsible(wrapper) {
  makeHeadingCollapsible(wrapper, 'H2', ['H2'], 'h2-body', false);
  makeHeadingCollapsible(wrapper, 'H3', ['H3', 'H2'], 'h3-body', true);
}

function makeHeadingCollapsible(root, tag, stopTags, bodyClass, defaultOpen = true) {
  [...root.querySelectorAll(tag.toLowerCase())].forEach(heading => {
    const body = document.createElement('div');
    body.className = bodyClass;

    let next = heading.nextElementSibling;
    while (next && !stopTags.includes(next.tagName)) {
      const after = next.nextElementSibling;
      body.appendChild(next);
      next = after;
    }

    heading.insertAdjacentElement('afterend', body);
    heading.classList.add('h-collapsible');
    heading.setAttribute('aria-expanded', String(defaultOpen));
    body.hidden = !defaultOpen;

    heading.addEventListener('click', e => {
      e.stopPropagation();
      const open = heading.getAttribute('aria-expanded') === 'true';
      heading.setAttribute('aria-expanded', String(!open));
      body.hidden = open;
    });
  });
}

/* ── Token substitution ────────────────────────────────── */
function escapeHtml(str) {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

function substituteHtml(origHtml) {
  let html = origHtml;
  TOKENS.forEach(token => {
    const re  = TOKEN_RE[token];
    re.lastIndex = 0;
    const val = state.inputs[token];
    html = html.replace(re, val
      ? escapeHtml(val)
      : `<span class="token-placeholder">&lt;${token}&gt;</span>`
    );
  });
  Object.entries(state.localInputs).forEach(([token, val]) => {
    const re = state.localTokenRe[token];
    if (!re) return;
    re.lastIndex = 0;
    html = html.replace(re, val
      ? escapeHtml(val)
      : `<span class="token-placeholder">&lt;${token}&gt;</span>`
    );
  });
  return html;
}

function applySubstitution() {
  $content.querySelectorAll('[data-orig-html]').forEach(el => {
    if (el.dataset.dirty) return;
    el.innerHTML = substituteHtml(el.dataset.origHtml);
  });
  $content.querySelectorAll('code.cmd-text').forEach(code => {
    if (code._cmdFill) code._cmdFill();
  });
}

/* ── Editable code blocks ──────────────────────────────── */
function makeEditable(codeEl) {
  codeEl.contentEditable = 'plaintext-only';
  codeEl.spellcheck = false;
  codeEl.addEventListener('input', () => { codeEl.dataset.dirty = 'true'; });
}

/* ── Code block action buttons (copy + reset) ──────────── */
function addCodeActions(preEl, codeEl) {
  const wrap = document.createElement('div');
  wrap.className = 'code-actions';

  // Reset button — restores auto-substituted content
  const resetBtn = document.createElement('button');
  resetBtn.className = 'reset-btn';
  resetBtn.title = 'Reset to original';
  resetBtn.textContent = '↺';
  resetBtn.addEventListener('click', () => {
    delete codeEl.dataset.dirty;
    codeEl.innerHTML = substituteHtml(codeEl.dataset.origHtml);
  });

  // Copy button
  const copyBtn = document.createElement('button');
  copyBtn.className = 'copy-btn';
  copyBtn.textContent = 'Copy';
  copyBtn.addEventListener('click', () => {
    let text;
    if (codeEl.dataset.dirty) {
      text = codeEl.textContent.split('\n').map(l => l.startsWith('- ') ? l.slice(2) : l).join('\n').trim();
    } else {
      text = codeEl.dataset.origText || '';
      TOKENS.forEach(t => {
        const v = state.inputs[t];
        if (v) text = text.replaceAll(`<${t}>`, v);
      });
      Object.entries(state.localInputs || {}).forEach(([t, v]) => {
        if (v) text = text.replaceAll(`<${t}>`, v);
      });
      text = text.split('\n').map(l => l.startsWith('- ') ? l.slice(2) : l).join('\n').trim();
    }
    navigator.clipboard.writeText(text).then(() => {
      copyBtn.textContent = 'Copied!';
      copyBtn.classList.add('copied');
      setTimeout(() => { copyBtn.textContent = 'Copy'; copyBtn.classList.remove('copied'); }, 1500);
    });
  });

  wrap.appendChild(resetBtn);
  wrap.appendChild(copyBtn);
  preEl.appendChild(wrap);
}

/* ── Welcome screen ────────────────────────────────────── */
function showWelcome() {
  $content.classList.add('welcome-mode');
  $content.innerHTML = `
    <div class="welcome">
      <div class="neon-sign">
        <div class="neon-sign-line"></div>
        <div class="welcome-title">CRACK-SOURCE</div>
        <div class="neon-sign-line"></div>
      </div>
    </div>`;
}
