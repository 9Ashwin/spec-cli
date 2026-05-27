const fs = require("fs");
const path = require("path");
const { execFileSync } = require("child_process");
const os = require("os");
const crypto = require("crypto");

const VERSION = require("../package.json").version;
const REPO = "9Ashwin/spec-cli";
const PKG = "ashwin-spec";
const NAME = "spec-cli";  // GoReleaser binary name

const PLATFORM_MAP = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const ARCH_MAP = {
  x64: "amd64",
  arm64: "arm64",
};

const platform = PLATFORM_MAP[process.platform];
const arch = ARCH_MAP[process.arch];

const isWindows = process.platform === "win32";
const ext = isWindows ? ".zip" : ".tar.gz";
const archiveName = `${NAME}-${VERSION}-${platform}-${arch}${ext}`;
const GITHUB_URL = `https://github.com/${REPO}/releases/download/v${VERSION}/${archiveName}`;

const binDir = path.join(__dirname, "..", "bin");
const dest = path.join(binDir, NAME + (isWindows ? ".exe" : ""));

function download(url, destPath) {
  const args = [
    "--fail", "--location", "--silent", "--show-error",
    "--connect-timeout", "10", "--max-time", "120",
    "--max-redirs", "3",
    "--output", destPath,
  ];
  if (isWindows) args.unshift("--ssl-revoke-best-effort");
  args.push(url);
  execFileSync("curl", args, { stdio: ["ignore", "ignore", "pipe"] });
}

function getExpectedChecksum(archiveName) {
  const checksumsPath = path.join(__dirname, "..", "checksums.txt");

  if (!fs.existsSync(checksumsPath)) {
    console.error("[spec-cli] checksums.txt not found, skipping checksum verification");
    return null;
  }

  const content = fs.readFileSync(checksumsPath, "utf8");
  for (const line of content.split("\n")) {
    const trimmed = line.trim();
    if (!trimmed) continue;
    const idx = trimmed.indexOf("  ");
    if (idx === -1) continue;
    const hash = trimmed.slice(0, idx);
    const name = trimmed.slice(idx + 2);
    if (name === archiveName) return hash;
  }

  throw new Error(`Checksum entry not found for ${archiveName}`);
}

function verifyChecksum(archivePath, expectedHash) {
  if (expectedHash === null) return;

  const hash = crypto.createHash("sha256");
  const fd = fs.openSync(archivePath, "r");
  try {
    const buf = Buffer.alloc(64 * 1024);
    let bytesRead;
    while ((bytesRead = fs.readSync(fd, buf, 0, buf.length, null)) > 0) {
      hash.update(buf.subarray(0, bytesRead));
    }
  } finally {
    fs.closeSync(fd);
  }
  const actual = hash.digest("hex");

  if (actual.toLowerCase() !== expectedHash.toLowerCase()) {
    throw new Error(
      `Checksum mismatch for ${path.basename(archivePath)}: expected ${expectedHash} but got ${actual}`
    );
  }
}

function install() {
  fs.mkdirSync(binDir, { recursive: true });

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "spec-cli-"));
  const archivePath = path.join(tmpDir, archiveName);

  try {
    download(GITHUB_URL, archivePath);

    const expectedHash = getExpectedChecksum(archiveName);
    verifyChecksum(archivePath, expectedHash);

    if (isWindows) {
      // Extract zip on Windows
      const tmpExtract = path.join(tmpDir, "extract");
      fs.mkdirSync(tmpExtract, { recursive: true });
      try {
        execFileSync("powershell.exe", [
          "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command",
          "$ErrorActionPreference='Stop';" +
          "Add-Type -AssemblyName System.IO.Compression.FileSystem;" +
          `[System.IO.Compression.ZipFile]::ExtractToDirectory('${archivePath}','${tmpExtract}')`,
        ], { stdio: "ignore" });
      } catch (_) {
        execFileSync("tar", ["-xf", archivePath, "-C", tmpExtract], { stdio: "ignore" });
      }
      const files = fs.readdirSync(tmpExtract);
      const binaryName = files.find(f => f.endsWith(".exe")) || NAME + ".exe";
      fs.copyFileSync(path.join(tmpExtract, binaryName), dest);
    } else {
      execFileSync("tar", ["-xzf", archivePath, "-C", tmpDir], { stdio: "ignore" });
      const binaryName = NAME;
      fs.copyFileSync(path.join(tmpDir, binaryName), dest);
    }

    fs.chmodSync(dest, 0o755);
    console.log(`spec-cli v${VERSION} installed successfully`);
  } finally {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  }
}

if (require.main === module) {
  if (!platform || !arch) {
    console.error(`Unsupported platform: ${process.platform}-${process.arch}`);
    process.exit(1);
  }

  // Skip binary download for npx postinstall (binary not needed yet).
  const isNpxPostinstall =
    process.env.npm_command === "exec" && !process.env.SPEC_CLI_RUN;

  if (isNpxPostinstall) {
    process.exit(0);
  }

  try {
    install();
  } catch (err) {
    console.error(`Failed to install spec-cli:`, err.message);
    process.exit(1);
  }
}
