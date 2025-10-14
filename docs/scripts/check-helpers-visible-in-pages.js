#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// Read all markdown files in docs/data directory
const dataDir = path.join(__dirname, '../data');
const files = fs.readdirSync(dataDir).filter(f => f.endsWith('.md'));

const combinations = new Set();
const coreCategories = new Set();
const pluginCategories = new Set();

// Extract type+category combinations from each file
files.forEach(file => {
  const filePath = path.join(dataDir, file);
  const content = fs.readFileSync(filePath, 'utf8');

  const typeMatch = content.match(/^type:\s*(.+)$/m);
  const categoryMatch = content.match(/^category:\s*(.+)$/m);

  if (typeMatch && categoryMatch) {
    const type = typeMatch[1].trim();
    const category = categoryMatch[1].trim();
    const combination = `${type}|${category}`;

    combinations.add(combination);

    if (type === 'core') {
      coreCategories.add(category);
    } else if (type === 'plugin') {
      pluginCategories.add(category);
    } else {
      throw new Error(`Invalid type in ${file}: ${type}`);
    }
  }
});

console.log('=== TYPE+CATEGORY COMBINATIONS FOUND ===');
Array.from(combinations).sort().forEach(comb => console.log(comb));

console.log('\n=== CORE CATEGORIES ===');
Array.from(coreCategories).sort().forEach(cat => console.log(cat));

console.log('\n=== PLUGIN CATEGORIES ===');
Array.from(pluginCategories).sort().forEach(cat => console.log(cat));

// Check existing pages
const pluginPagesDir = path.join(__dirname, '../docs/plugins');
const operatorPagesDir = path.join(__dirname, '../docs/operator');

const existingPluginPages = new Set();
const existingCorePages = new Set();

if (fs.existsSync(pluginPagesDir)) {
  fs.readdirSync(pluginPagesDir)
    .filter(f => f.endsWith('.md'))
    .forEach(f => existingPluginPages.add(f.replace('.md', '')));
}

if (fs.existsSync(operatorPagesDir)) {
  fs.readdirSync(operatorPagesDir)
    .filter(f => f.endsWith('.md'))
    .forEach(f => existingCorePages.add(f.replace('.md', '')));
}

console.log('\n=== EXISTING PLUGIN PAGES ===');
Array.from(existingPluginPages).sort().forEach(page => console.log(page));

console.log('\n=== EXISTING CORE PAGES ===');
Array.from(existingCorePages).sort().forEach(page => console.log(page));

// Find missing pages
console.log('\n=== MISSING PLUGIN PAGES ===');
Array.from(pluginCategories).sort().forEach(category => {
  if (!existingPluginPages.has(category)) {
    console.log(`MISSING: plugin/${category}.md`);
  }
});

console.log('\n=== MISSING CORE PAGES ===');
Array.from(coreCategories).sort().forEach(category => {
  if (!existingCorePages.has(category)) {
    console.log(`MISSING: core/${category}.md`);
  }
});

// Check for duplicates
console.log('\n=== VALIDATION RESULTS ===');
let hasErrors = false;

Array.from(pluginCategories).sort().forEach(category => {
  if (!existingPluginPages.has(category)) {
    console.log(`❌ ERROR: Missing plugin page for category: ${category}`);
    hasErrors = true;
  }
});

Array.from(coreCategories).sort().forEach(category => {
  if (!existingCorePages.has(category)) {
    console.log(`❌ ERROR: Missing core page for category: ${category}`);
    hasErrors = true;
  }
});

if (!hasErrors) {
  console.log('✅ All helper categories have corresponding pages!');
} else {
  console.log('\n❌ Found missing pages. Please create them as shown above.');
  process.exit(1);
}