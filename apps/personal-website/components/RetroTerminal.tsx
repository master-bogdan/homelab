"use client";
import React, { useEffect, useRef, useState } from "react";
import Link from "next/link";

const THEME = {
  bg: "#0b0f10",
  green: "#00ff9c",
  text: "#b8ffc9",
  link: "#8af5ff",
};

const promptUser = "bogdan";
const promptHost = "homelab";

const projects = [
  {
    name: "OAuth2/OIDC Server (Go)",
    desc: "Stateful Paseto tokens, Auth Code + Refresh, Fiber, PostgreSQL.",
    url: "https://github.com/yourname/oidc-server",
  },
  {
    name: "Clear Cash API (Go)",
    desc: "Modular architecture, Paseto auth, Profiles & Currencies.",
    url: "https://github.com/yourname/clear-cash-api",
  },
  {
    name: "eSignature Service (TS)",
    desc: "DocuSign-like internal tool, saved ~$100k/yr.",
    url: "https://github.com/yourname/esign",
  },
];

const postsPreview = [
  {
    title: "Kubernetes CronJobs for Scheduling",
    slug: "k8s-cronjobs",
    date: "2025-03-13",
  },
  {
    title: "OAuth2 + OIDC in Go: From Scratch",
    slug: "oauth2-oidc-go",
    date: "2025-07-01",
  },
  {
    title: "Undo/Redo with RTK + Immer",
    slug: "redux-undo-immer",
    date: "2025-06-05",
  },
];

export default function RetroTerminal() {
  const [output, setOutput] = useState<string[]>([]);
  const [input, setInput] = useState("");
  const [locked, setLocked] = useState(false);
  const [booted, setBooted] = useState(false);
  const [caretVisible, setCaretVisible] = useState(true);
  const hasBootedRef = useRef(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const iv = setInterval(() => setCaretVisible((v) => !v), 580);
    return () => clearInterval(iv);
  }, []);

  useEffect(() => {
    containerRef.current?.focus();
  }, [booted]);

  // ✅ Boot animation once only
  useEffect(() => {
    if (hasBootedRef.current) return;
    hasBootedRef.current = true;

    (async () => {
      setLocked(true);
      setOutput([]);

      const lines = [
        "boot: initializing CRT…",
        "boot: loading /usr/bin/portfolio …",
        "boot: mounting blog …",
        "boot: OK",
        "",
        "System ready. Type 'help' or click a command below.",
      ];

      const acc: string[] = [];
      for (const line of lines) {
        acc.push(line);
        setOutput([...acc]);
        await delay(250);
      }

      setLocked(false);
      setBooted(true);
    })();
  }, []);

  function delay(ms: number) {
    return new Promise((r) => setTimeout(r, ms));
  }

  // Typing animation – overwrites output with each new command
  async function typeOutput(lines: string[], speed = 12) {
    setLocked(true);
    const acc: string[] = [];
    for (const line of lines) {
      let typed = "";
      for (const ch of line) {
        typed += ch;
        setOutput([...acc, typed]);
        await delay(speed);
      }
      acc.push(line);
      setOutput([...acc]);
      await delay(60);
    }
    setLocked(false);
  }

  async function runCommand(cmdLine: string) {
    const [cmd] = cmdLine.trim().split(/\s+/);

    switch (cmd) {
      case "help":
        await typeOutput([
          "bogdan@homelab:~$: help",
          "Available commands:",
          "whoami, projects, blog, resume, contact, clear",
          "",
          "Click a command below to execute.",
        ]);
        break;
      case "whoami":
        await typeOutput([
          "bogdan@homelab:~$: whoami",
          "Senior Software Engineer — Node.js, TypeScript & Go",
          "Backend, DevOps & Cloud (AWS, Kubernetes). Platform Engineer.",
        ]);
        break;
      case "projects":
        await typeOutput(
          projects.flatMap((p, i) => [
            `${i + 1}. ${p.name}`,
            `   ${p.desc}`,
            `   → ${p.url}`,
          ]),
        );
        break;
      case "blog":
        await typeOutput([
          "bogdan@homelab:~$: blog",
          "Recent posts:",
          ...postsPreview.map(
            (p) => `${p.date} — ${p.title} — open: /blog/${p.slug}`,
          ),
          "",
          "Click below to open posts 👇",
        ]);
        break;
      case "resume":
        await typeOutput([
          "bogdan@homelab:~$: resume",
          "Opening resume… → /resume.pdf",
        ]);
        break;
      case "contact":
        await typeOutput([
          "bogdan@homelab:~$: contact",
          "Email: bshchavinskyi@gmail.com",
          "GitHub: https://github.com/yourname",
          "LinkedIn: https://www.linkedin.com/in/yourname",
        ]);
        break;
      case "clear":
        setOutput([]);
        break;
      default:
        await typeOutput([`command not found: ${cmd}`]);
    }
  }

  // ✅ Echo command first, clear screen before each new one
  async function onKeyDown(e: React.KeyboardEvent<HTMLDivElement>) {
    if (locked) return;
    if (e.key === "Backspace") {
      setInput((s) => s.slice(0, -1));
      e.preventDefault();
      return;
    }
    if (e.key === "Enter") {
      const cmd = input.trim();
      setInput("");
      // Clear screen and echo new command
      setOutput([`${promptUser}@${promptHost}:~$ ${cmd}`]);
      await runCommand(cmd);
      e.preventDefault();
      return;
    }
    if (e.key.length === 1 && !e.ctrlKey && !e.metaKey && !e.altKey) {
      setInput((s) => s + e.key);
      e.preventDefault();
    }
  }

  // same behavior for quick-click commands
  function clickCmd(cmd: string) {
    if (locked) return;
    setInput("");
    setOutput([`${promptUser}@${promptHost}:~$ ${cmd}`]);
    void runCommand(cmd);
  }

  const styles: { [k: string]: React.CSSProperties } = {
    wrap: {
      position: "relative",
      minHeight: "100vh",
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
      background: THEME.bg,
      overflow: "hidden",
    },
    frame: {
      width: "min(1000px, 95vw)",
      height: "min(700px, 85vh)",
      borderRadius: 18,
      border: `2px solid ${THEME.green}55`,
      boxShadow: `0 0 20px ${THEME.green}22, 0 0 120px ${THEME.green}11`,
      overflow: "hidden",
      position: "relative",
    },
    screen: {
      position: "absolute",
      inset: 0,
      padding: 24,
      color: THEME.text,
      fontFamily:
        "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace",
      fontSize: 16,
      lineHeight: 1.6,
      textShadow: `0 0 10px ${THEME.green}33`,
      whiteSpace: "pre-wrap",
      outline: "none",
      display: "flex",
      flexDirection: "column",
      justifyContent: "flex-start",
      gap: 8,
      overflowY: "auto",
    },
    link: {
      color: THEME.link,
      textDecoration: "none",
      borderBottom: `1px dotted ${THEME.link}66`,
      cursor: "pointer",
    },
    caret: {
      display: "inline-block",
      width: 9,
      height: 18,
      marginLeft: 2,
      background: THEME.green,
      boxShadow: `0 0 12px ${THEME.green}`,
      verticalAlign: "-3px",
      visibility: caretVisible ? "visible" : "hidden",
    },
  };

  return (
    <div style={styles.wrap}>
      <div style={styles.frame}>
        <div
          ref={containerRef}
          tabIndex={0}
          onKeyDown={onKeyDown}
          style={styles.screen}
        >
          <Banner />
          {output.map((line, i) => {
            if (line.includes("/blog/")) {
              const slug = line.split("/blog/")[1];
              return (
                <div key={i}>
                  <Link href={`/blog/${slug}`} style={styles.link}>
                    {line}
                  </Link>
                </div>
              );
            }
            return <div key={i}>{line}</div>;
          })}
          {booted && (
            <div>
              <span style={{ color: THEME.green, fontWeight: 600 }}>
                {`${promptUser}@${promptHost}:~$ `}
              </span>
              <span>{input}</span>
              <span style={styles.caret} />
            </div>
          )}
          <div style={{ marginTop: 10 }}>
            {[
              "whoami",
              "projects",
              "blog",
              "resume",
              "contact",
              "help",
              "clear",
            ].map((cmd) => (
              <span
                key={cmd}
                style={{ ...styles.link, marginRight: 14 }}
                onClick={() => clickCmd(cmd)}
              >
                {cmd}
              </span>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

function Banner() {
  const color = THEME.green;
  const boxStyle: React.CSSProperties = {
    color,
    fontFamily: "inherit",
    fontSize: 13,
    lineHeight: 1.3,
    marginBottom: 4,
    textShadow: `0 0 8px ${color}55, 0 0 20px ${color}33`,
  };
  return (
    <pre
      style={boxStyle}
    >{String.raw`┌──────────────────────────────────────────────┐
│  Bogdan Shchavinskyi                         │
│  Senior Software Engineer                    │
│  Go · TypeScript · DevOps · Kubernetes       │
└──────────────────────────────────────────────┘`}</pre>
  );
}
