import { getAllPosts, getPostBySlug } from "@/lib/posts";
import React from "react";

export async function generateStaticParams() {
  const posts = await getAllPosts();
  return posts.map((p) => ({ slug: p.slug }));
}

export default async function BlogPost({
  params,
}: {
  params: { slug: string };
}) {
  const { html, meta } = await getPostBySlug(params.slug);
  return (
    <main
      style={{
        maxWidth: 820,
        margin: "40px auto",
        padding: 16,
        color: "#b8ffc9",
        fontFamily: "ui-monospace, Menlo, monospace",
      }}
    >
      <h1 style={{ color: "#00ff9c" }}>{meta.title}</h1>
      <p style={{ opacity: 0.7 }}>{meta.date}</p>
      <article dangerouslySetInnerHTML={{ __html: html }} />
    </main>
  );
}
