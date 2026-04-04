import fs from "fs/promises";

const OWNER = "trymist";
const REPO = "mist";
const BRANCH = "main";

const URL = `https://raw.githubusercontent.com/${OWNER}/${REPO}/${BRANCH}/CHANGELOG.md`;

async function main() {
  console.log("Fetching changelog...");

  const res = await fetch(URL);

  if (!res.ok) {
    throw new Error(`Failed to fetch changelog: ${res.status}`);
  }

  let text = await res.text();

  // Optional: add frontmatter for VitePress
  text =
    `---\n` +
    `title: Changelog\n` +
    `outline: deep\n` +
    `---\n\n` +
    text;

  await fs.writeFile("./docs/changelog.md", text);

  console.log("Changelog synced");
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
