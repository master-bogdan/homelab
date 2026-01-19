# Terminal Portfolio - Bogdan Shchavinskyi

A terminal-style personal portfolio website built with Next.js 16, TypeScript, and Material UI.

## Features

- **Terminal Interface**: Authentic bash-style terminal with boot animation
- **Interactive Commands**: Type commands or use quick-command buttons
- **Blog System**: Read-only blog powered by local Markdown files
- **Hidden Easter Eggs**: Discover hidden commands for fun surprises
- **Fully Typed**: Strict TypeScript for type safety
- **Tested**: Comprehensive unit tests with Vitest and React Testing Library
- **Accessible**: ARIA-compliant and keyboard-navigable
- **Responsive**: Works on all screen sizes

## Tech Stack

- **Framework**: Next.js 16 (Client Components only)
- **Language**: TypeScript (strict mode)
- **UI Library**: Material UI v6
- **Styling**: Emotion (CSS-in-JS)
- **Markdown**: react-markdown + gray-matter
- **Testing**: Vitest + React Testing Library
- **Code Quality**: ESLint + Prettier

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

1. Clone the repository:
\`\`\`bash
git clone https://github.com/bogdan/terminal-portfolio.git
cd terminal-portfolio
\`\`\`

2. Install dependencies:
\`\`\`bash
npm install
\`\`\`

3. Run the development server:
\`\`\`bash
npm run dev
\`\`\`

4. Open [http://localhost:3000](http://localhost:3000) in your browser.

## Available Commands

### Development

- \`npm run dev\` - Start development server
- \`npm run build\` - Build for production
- \`npm run start\` - Start production server
- \`npm run lint\` - Run ESLint

### Testing

- \`npm test\` - Run tests in watch mode
- \`npm run test:ui\` - Run tests with UI
- \`npm run test:coverage\` - Generate coverage report

## Terminal Commands

### Visible Commands

- \`whoami\` - Display information about me
- \`projects\` - List featured projects
- \`blog\` - List all blog posts
- \`blog [slug]\` - Read a specific blog post
- \`contacts\` - Get contact information
- \`resume\` - Download resume
- \`clear\` - Clear the terminal
- \`help\` - Show available commands

### Hidden Commands

Try discovering these yourself! Hint: think about common Unix commands...

## Project Structure

\`\`\`
.
├── app/
│   ├── layout.tsx          # Root layout with MUI theme
│   ├── page.tsx            # Main terminal page
│   └── globals.css         # Global styles
├── components/
│   └── terminal/
│       ├── Terminal.tsx    # Main terminal component
│       ├── BootSequence.tsx
│       ├── Banner.tsx
│       ├── Prompt.tsx
│       ├── CommandInput.tsx
│       └── QuickCommands.tsx
├── lib/
│   ├── commands/           # Command implementations
│   ├── commandRegistry.tsx # Command registry
│   ├── blog.ts            # Blog post loader
│   ├── constants.ts       # App constants
│   ├── theme.ts           # MUI theme
│   └── types.ts           # TypeScript types
├── hooks/
│   ├── useTypingAnimation.ts
│   └── useCyclingWords.ts
├── content/
│   └── blog/              # Markdown blog posts
├── __tests__/             # Unit tests
└── README.md
\`\`\`

## Adding Blog Posts

1. Create a new Markdown file in \`content/blog/\`:

\`\`\`markdown
---
title: "Your Post Title"
date: "2024-01-01"
slug: "your-post-slug"
---

# Your Post Title

Your content here...
\`\`\`

2. The post will automatically appear in the blog command.

## Customization

### Colors

Edit \`lib/theme.ts\` to change the terminal color scheme.

### Commands

Add new commands in \`lib/commands/\` and register them in \`lib/commandRegistry.tsx\`.

## Testing

Tests are located in the \`__tests__/\` directory and mirror the source structure.

Run tests:
\`\`\`bash
npm test
\`\`\`

Generate coverage:
\`\`\`bash
npm run test:coverage
\`\`\`

## Accessibility

- Keyboard navigation supported
- ARIA labels on interactive elements
- Screen reader friendly
- High contrast terminal colors

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)

## License

MIT

## Contact

- Email: bogdan@example.com
- GitHub: [@bogdan](https://github.com/bogdan)
- LinkedIn: [bogdan](https://linkedin.com/in/bogdan)

---

Built with ❤️ and lots of coffee
\`\`\`
