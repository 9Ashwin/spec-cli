#!/usr/bin/env node

const { execFileSync } = require("child_process");
const fs = require("fs");
const path = require("path");

const ext = process.platform === "win32" ? ".exe" : "";
const bin = path.join(__dirname, "..", "bin", "spec-cli" + ext);
const args = process.argv.slice(2);

// Auto-download binary if missing (e.g. npx skipped postinstall).
if (!fs.existsSync(bin)) {
  try {
    execFileSync(process.execPath, [path.join(__dirname, "install.js")], {
      stdio: "inherit",
      env: { ...process.env, SPEC_CLI_RUN: "true" },
    });
  } catch (_) {
    console.error(
      "\nspec-cli binary not found.\n" +
      "Run the install script manually:\n" +
      '  node "' + path.join(__dirname, "install.js") + '"\n'
    );
    process.exit(1);
  }
}

try {
  execFileSync(bin, args, { stdio: "inherit" });
} catch (e) {
  process.exit(e.status || 1);
}
