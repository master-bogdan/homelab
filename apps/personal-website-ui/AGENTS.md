# AGENTS.md

## Project

Personal Website UI is a Next.js 16 + React + TypeScript terminal-style portfolio site using Material UI, local Markdown blog content, and Vitest tests.

## Architecture

- `app/` contains Next.js app routes and layout.
- `lib/` owns command registry, blog loading, constants, theme, and shared types.
- `hooks/` contains reusable UI behavior hooks.
- `styles/` contains global styling.
- `public/` contains static assets and downloadable files.
- Preserve the terminal-style interaction model and command behavior.
- Blog content is local Markdown and should continue to load through existing blog utilities.

## Commands

- Dev server: `npm run dev`
- Build: `npm run build`
- Start production server: `npm run start`
- Lint: `npm run lint`
- Tests: `npm test`
- Coverage: `npm run test:coverage`
- Test UI: `npm run test:ui`

## Rules

- Preserve keyboard accessibility and responsive behavior.
- Keep command implementations registered through the existing command registry pattern.
- Do not hardcode generated build output into source changes.
- Do not edit `.next/` or `build/` artifacts as source.
- Keep TypeScript strict and avoid weakening types to bypass errors.
- Update tests when changing terminal commands, blog behavior, or visible portfolio interactions.

## Verification

For UI/source changes, run:

```bash
npm run build
npm run lint
```

For behavior changes, also run:

```bash
npm test
```
