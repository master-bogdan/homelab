export interface BlogPost {
  slug: string
  title: string
  date: string
  content: string
}

// This will be populated at build time
const blogPosts: BlogPost[] = []

// Function to load blog posts (called at build time)
export const loadBlogPosts = (): BlogPost[] => {
  // In a real implementation, this would read from the file system
  // For now, we'll return the posts we've created
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
    {
      slug: "kubernetes-at-scale",
      title: "Scaling Kubernetes: Lessons Learned",
      date: "2024-02-20",
      content: `# Scaling Kubernetes: Lessons Learned

After managing Kubernetes clusters serving thousands of users, here are the key lessons I've learned about scaling container orchestration.

## 1. Resource Limits Are Not Optional

Always set resource requests and limits. Always. Your future self will thank you when debugging OOMKilled pods at 3 AM.

\`\`\`yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "100m"
  limits:
    memory: "512Mi"
    cpu: "500m"
\`\`\`

## 2. Horizontal Pod Autoscaling

HPA is your friend, but configure it wisely. Start conservative and tune based on actual metrics.

## 3. Network Policies Matter

Don't leave your cluster wide open. Implement network policies from day one.

## Conclusion

Kubernetes is powerful but complex. Start simple, monitor everything, and scale gradually.`,
    },
    {
      slug: "nodejs-performance",
      title: "Node.js Performance Tips",
      date: "2024-03-10",
      content: `# Node.js Performance Tips

Node.js is fast, but you can make it faster. Here are my top performance optimization techniques.

## 1. Use Async/Await Properly

Don't block the event loop. Use async/await for I/O operations and keep CPU-intensive tasks in worker threads.

## 2. Connection Pooling

Always use connection pooling for databases. Creating new connections is expensive.

\`\`\`typescript
const pool = new Pool({
  max: 20,
  idleTimeoutMillis: 30000,
  connectionTimeoutMillis: 2000,
});
\`\`\`

## 3. Caching Strategy

Implement multi-layer caching: in-memory, Redis, and CDN.

## 4. Monitoring

Use APM tools like New Relic or DataDog. You can't optimize what you don't measure.

Stay fast!`,
    },
  ]
}

export const getBlogPosts = (): BlogPost[] => {
  return loadBlogPosts()
}

export const getBlogPost = (slug: string): BlogPost | undefined => {
  return getBlogPosts().find((post) => post.slug === slug)
}
