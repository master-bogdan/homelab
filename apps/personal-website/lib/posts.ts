import fs from "node:fs/promises";
import path from "node:path";
import matter from "gray-matter";
import { remark } from "remark";
import html from "remark-html";

const postsDir = path.join(process.cwd(), "posts");

export interface PostMeta {
  title: string;
  date: string; // always normalized to string
  [key: string]: any;
}

export interface PostSummary {
  slug: string;
  meta: PostMeta;
}

export interface PostDetail {
  html: string;
  meta: PostMeta;
}

/**
 * Reads all markdown posts and returns sorted metadata
 * (newest first). Ensures all dates are strings to prevent
 * React render errors when gray-matter outputs Date objects.
 */
export async function getAllPosts(): Promise<PostSummary[]> {
  const files = await fs.readdir(postsDir);
  const posts = await Promise.all(
    files
      .filter((f) => f.endsWith(".md"))
      .map(async (file) => {
        const slug = file.replace(/\.md$/, "");
        const full = await fs.readFile(path.join(postsDir, file), "utf8");
        const { data } = matter(full);

        const date =
          typeof data.date === "string"
            ? data.date
            : new Date(data.date).toISOString().split("T")[0];

        return {
          slug,
          meta: { ...data, date } as PostMeta,
        };
      })
  );

  // sort newest first
  return posts.sort((a, b) => (a.meta.date < b.meta.date ? 1 : -1));
}

/**
 * Reads a single markdown post by slug and converts to HTML.
 * Always returns a normalized date string in metadata.
 */
export async function getPostBySlug(slug: string): Promise<PostDetail> {
  const filePath = path.join(postsDir, `${slug}.md`);
  const full = await fs.readFile(filePath, "utf8");
  const { data, content } = matter(full);
  const result = await remark().use(html).process(content);

  const date =
    typeof data.date === "string"
      ? data.date
      : new Date(data.date).toISOString().split("T")[0];

  return {
    html: result.toString(),
    meta: { ...data, date } as PostMeta,
  };
}
