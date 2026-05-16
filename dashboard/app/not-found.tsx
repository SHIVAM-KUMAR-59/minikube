'use client'
import Link from 'next/link'
import { LayoutDashboard, RotateCcw } from 'lucide-react'

export default function NotFound() {
  return (
    <div className="min-h-screen bg-base flex flex-col items-center justify-center px-6 relative overflow-hidden">

      {/* Ambient grid background */}
      <div
        className="absolute inset-0 opacity-[0.03]"
        style={{
          backgroundImage: `
            linear-gradient(var(--violet) 1px, transparent 1px),
            linear-gradient(90deg, var(--violet) 1px, transparent 1px)
          `,
          backgroundSize: '48px 48px',
        }}
      />

      {/* Glow orb */}
      <div
        className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-150 h-150 rounded-full opacity-[0.06] pointer-events-none"
        style={{ background: 'radial-gradient(circle, var(--violet) 0%, transparent 70%)' }}
      />

      {/* Corner decorations */}
      <div className="absolute top-8 left-8 w-16 h-16 border-l border-t border-border-subtle rounded-tl-lg opacity-40" />
      <div className="absolute top-8 right-8 w-16 h-16 border-r border-t border-border-subtle rounded-tr-lg opacity-40" />
      <div className="absolute bottom-8 left-8 w-16 h-16 border-l border-b border-border-subtle rounded-bl-lg opacity-40" />
      <div className="absolute bottom-8 right-8 w-16 h-16 border-r border-b border-border-subtle rounded-br-lg opacity-40" />

      {/* Content */}
      <div className="relative z-10 flex flex-col items-center text-center max-w-lg">

        {/* Status badge */}
        <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-card border border-border-subtle mb-10">
          <span className="relative flex h-1.5 w-1.5">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-failed opacity-60" />
            <span className="relative inline-flex rounded-full h-1.5 w-1.5 bg-failed" />
          </span>
          <span className="text-text-muted text-xs font-mono tracking-wide">pod/not-found — STATUS: 404</span>
        </div>

        {/* 404 */}
        <div className="relative mb-6">
          <p
            className="text-[130px] font-bold leading-none select-none pointer-events-none"
            style={{
              color: 'transparent',
              WebkitTextStroke: '1px rgba(124, 58, 237, 0.25)',
              letterSpacing: '-0.05em',
            }}
          >
            404
          </p>
          <p
            className="absolute inset-0 text-[130px] font-bold leading-none select-none pointer-events-none flex items-center justify-center"
            style={{
              color: 'transparent',
              backgroundClip: 'text',
              WebkitBackgroundClip: 'text',
              backgroundImage: 'linear-gradient(180deg, var(--violet-lt) 0%, var(--violet) 60%, transparent 100%)',
              letterSpacing: '-0.05em',
            }}
          >
            404
          </p>
        </div>

        {/* Heading */}
        <h1 className="text-text-primary text-2xl font-semibold mb-3 tracking-tight">
          Pod not found
        </h1>

        {/* Description */}
        <p className="text-text-secondary text-sm leading-relaxed mb-2">
          The scheduler could&apos;nt locate this resource in the cluster.
          It may have been terminated, evicted, or never existed.
        </p>

        {/* Mono detail */}
        <p className="text-text-muted text-xs font-mono mb-10 opacity-70">
          ErrImageNotFound: no route matches the requested path
        </p>

        {/* Terminal block */}
        <div className="w-full bg-card border border-border-subtle rounded-lg px-4 py-3 mb-10 text-left">
          <div className="flex items-center gap-1.5 mb-3">
            <span className="w-2.5 h-2.5 rounded-full bg-failed opacity-60" />
            <span className="w-2.5 h-2.5 rounded-full bg-pending opacity-60" />
            <span className="w-2.5 h-2.5 rounded-full bg-running opacity-60" />
          </div>
          <p className="font-mono text-xs text-text-muted leading-relaxed">
            <span className="text-violet-lt">$</span>
            <span className="text-text-secondary"> minik get pod </span>
            <span className="text-text-primary">unknown-resource</span>
            <br />
            <span className="text-failed-text">Error</span>
            <span className="text-text-muted"> from server (NotFound):</span>
            <br />
            <span className="text-text-muted pl-4">pods &quot;unknown-resource&quot; not found</span>
          </p>
        </div>

        {/* Actions */}
        <div className="flex items-center gap-3">
          <Link
            href="/"
            className="flex items-center gap-2 px-4 py-2.5 rounded-lg bg-violet hover:bg-violet/90 text-white text-sm font-medium transition-all duration-150 shadow-lg shadow-violet/20"
          >
            <LayoutDashboard size={14} strokeWidth={2} />
            Back to overview
          </Link>
          <button
            onClick={() => window.history.back()}
            className="flex items-center gap-2 px-4 py-2.5 rounded-lg bg-card hover:bg-overlay border border-border-subtle text-text-secondary hover:text-text-primary text-sm font-medium transition-all duration-150"
          >
            <RotateCcw size={14} strokeWidth={2} />
            Go back
          </button>
        </div>

      </div>
    </div>
  )
}