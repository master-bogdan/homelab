export interface BlogPost {
  slug: string
  title: string
  date: string
  content: string
}

export const loadBlogPosts = (): BlogPost[] => {
  return [
    {
      slug: "welcome",
      title: "Welcome to My Terminal",
      date: "2024-01-15",
      content: `# Welcome to My Terminal

Hey there! Welcome to my personal site. I built this terminal-style interface because I believe in keeping things simple and functional.

## Why a Terminal?

As a backend engineer who spends most of my time in the command line, this felt like the most natural way to present my work. Plus, it's just cool.

## What You'll Find Here

- Technical deep-dives into backend architecture
- DevOps best practices and war stories
- Cloud infrastructure patterns
- Occasional rants about code quality

Feel free to explore using the commands. There might be some hidden gems if you're curious enough.

Happy exploring!`,
    },
  ]
}

export const getBlogPosts = (): BlogPost[] => {
  return loadBlogPosts()
}

export const getBlogPost = (slug: string): BlogPost | undefined => {
  return getBlogPosts().find((post) => post.slug === slug)
}
