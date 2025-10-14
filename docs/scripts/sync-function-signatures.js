#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const readline = require('readline');
const { listMarkdownFiles, parseFrontmatter } = require('./utils');

const repoRoot = path.resolve(__dirname, '..', '..');
const dataDir = path.resolve(__dirname, '..', 'data');

function readFile(filePath) {
  return fs.readFileSync(filePath, 'utf8');
}

function writeFile(filePath, content) {
  fs.writeFileSync(filePath, content, 'utf8');
}

function* walkGoFiles(dir, excludeDirs = new Set()) {
  const entries = fs.readdirSync(dir, { withFileTypes: true });
  for (const entry of entries) {
    if (entry.name.startsWith('.')) continue;
    const abs = path.join(dir, entry.name);
    const rel = path.relative(repoRoot, abs);
    if (entry.isDirectory()) {
      if (excludeDirs.has(entry.name)) continue;
      yield* walkGoFiles(abs, excludeDirs);
    } else if (entry.isFile() && entry.name.endsWith('.go')) {
      // skip tests
      if (entry.name.endsWith('_test.go')) continue;
      // skip docs/ directory
      if (rel.split(path.sep)[0] === 'docs') continue;
      yield abs;
    }
  }
}

function buildFunctionRegex(name) {
  // Matches: func Name[...]( or func Name(
  const escaped = name.replace(/[-/\\^$*+?.()|[\]{}]/g, '\\$&');
  return new RegExp('^func\\s+' + escaped + '(?:\\s*\\[[^\]]*\\])?\\s*\\(', '');
}

async function findFunctionDeclaration(name, preferredPathHint) {
  const fnRegex = buildFunctionRegex(name);

  // Prefer hinted file if provided
  if (preferredPathHint) {
    const hintedAbs = path.resolve(repoRoot, preferredPathHint);
    if (fs.existsSync(hintedAbs)) {
      const hit = await scanFileForSignature(hintedAbs, fnRegex);
      if (hit) return hit;
    }
  }

  for (const abs of walkGoFiles(repoRoot)) {
    const hit = await scanFileForSignature(abs, fnRegex);
    if (hit) return hit;
  }
  return null;
}

function stripBOM(s) {
  return s.charCodeAt(0) === 0xfeff ? s.slice(1) : s;
}

async function scanFileForSignature(absPath, fnRegex) {
  const rl = readline.createInterface({
    input: fs.createReadStream(absPath, { encoding: 'utf8' }),
    crlfDelay: Infinity,
  });
  let lineNo = 0;
  for await (const rawLine of rl) {
    lineNo++;
    const line = stripBOM(rawLine);
    if (fnRegex.test(line)) {
      // Normalize multiple spaces and tabs minimally: keep original line
      const signature = line.trim();
      const rel = path.relative(repoRoot, absPath).replace(/\\/g, '/');
      return { file: rel, line: lineNo, signature };
    }
  }
  return null;
}

function updateFrontmatter(content, updates) {
  const m = content.match(/^(---[\r\n]+)([\s\S]*?)([\r\n]+---)([\s\S]*)$/);
  if (!m) return null;
  const prefix = m[1];
  const fmBody = m[2];
  const suffix = m[3];
  const rest = m[4];

  const lines = fmBody.split(/\r?\n/);

  function setKey(key, value, quote = false) {
    const idx = lines.findIndex((l) => l.startsWith(key + ':'));
    const v = quote ? `"${value}"` : value;
    const line = `${key}: ${v}`;
    if (idx >= 0) {
      lines[idx] = line;
    } else {
      lines.push(line);
    }
  }

  if (updates.sourceRef) setKey('sourceRef', updates.sourceRef, false);
  if (updates.signature) setKey('signature', updates.signature, true);

  const newFm = lines.join('\n');
  return prefix + newFm + suffix + rest;
}

function parseSourceRefFile(sourceRef) {
  if (!sourceRef) return null;
  const idx = sourceRef.indexOf('#');
  if (idx === -1) return sourceRef;
  return sourceRef.slice(0, idx);
}

async function main() {
  const args = new Set(process.argv.slice(2));
  const files = listMarkdownFiles(dataDir);
  let changed = 0;
  for (const absPath of files) {
    const content = readFile(absPath);
    const fm = parseFrontmatter(content) || {};
    const name = fm.name;
    if (!name) continue;

    const hintFile = parseSourceRefFile(fm.sourceRef);
    const hit = await findFunctionDeclaration(name, hintFile);
    if (!hit) continue;

    const newSourceRef = `${hit.file}#L${hit.line}`;
    const newSignature = hit.signature;

    const needsSourceRef = fm.sourceRef !== newSourceRef;
    const needsSignature = fm.signature !== newSignature;
    if (!needsSourceRef && !needsSignature) continue;

    const updated = updateFrontmatter(content, {
      sourceRef: newSourceRef,
      signature: newSignature,
    });
    if (updated) {
      writeFile(absPath, updated);
      changed++;
      // eslint-disable-next-line no-console
      console.log(`[updated] ${path.relative(repoRoot, absPath)} -> ${newSourceRef}`);
    }
  }

  if (args.has('--check') && changed > 0) {
    process.exitCode = 1;
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});


