import { $ } from "zx";
import { readFileSync } from "fs";

// Don't print commands before running
$.verbose = true;

async function release() {
  // Check for uncommitted changes first
  const status = await $`git status --porcelain`.quiet();
  if (status.stdout.trim()) {
    console.log("⚠️  You have uncommitted changes:");
    console.log(status.stdout);
    console.log("Please commit or stash them before releasing.");
    process.exit(1);
  }

  // Check for pending changesets
  const changesetDir = ".changeset";
  const files = await $`ls ${changesetDir}/*.md 2>/dev/null || true`.quiet();
  const hasChangesets = files.stdout
    .split("\n")
    .filter((f) => f && !f.endsWith("README.md")).length > 0;

  if (!hasChangesets) {
    console.log("❌ No changesets found. Run 'pnpm changeset' first to create one.");
    process.exit(1);
  }

  console.log("📦 Running changeset version...");
  await $`pnpm changeset version`;

  console.log("🏷️  Creating git tag...");
  await $`pnpm changeset tag`;

  // Get the new version from package.json
  const pkg = JSON.parse(readFileSync("package.json", "utf-8"));
  const version = pkg.version;

  console.log(`📝 Committing release v${version}...`);
  await $`git add .`;
  await $`git commit -m ${"chore: release v" + version}`;

  console.log("🚀 Pushing to origin...");
  await $`git push`;
  await $`git push --tags`;

  console.log(`\n✅ Released v${version}!`);
  console.log("🔨 GitHub Actions will now build and publish the binaries.");
  console.log(`📋 Check: https://github.com/alexanderchan/alex-runner/actions`);
}

release().catch((err) => {
  console.error("❌ Release failed:", err.message);
  process.exit(1);
});

