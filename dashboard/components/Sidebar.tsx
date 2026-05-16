'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { LayoutDashboard, Box, Server, Network } from 'lucide-react'

const links = [
  { label: 'Overview', href: '/', icon: LayoutDashboard },
  { label: 'Pods', href: '/pods', icon: Box },
  { label: 'Nodes', href: '/nodes', icon: Server },
  { label: 'Services', href: '/services', icon: Network },
]

export default function Sidebar() {
  const pathname = usePathname()

  return (
    <aside className="flex flex-col w-56 min-h-screen bg-surface border-r border-border-subtle px-3 py-5">
      {/* Logo */}
      <div className="flex items-center gap-2.5 px-3 mb-8">
        <div className="w-7 h-7 rounded-md bg-violet flex items-center justify-center">
          <svg
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="white"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
          </svg>
        </div>
        <span className="text-text-primary font-semibold text-sm tracking-wide">
          minikube
        </span>
      </div>

      {/* Nav */}
      <nav className="flex flex-col gap-1">
        <p className="text-text-muted text-xs font-medium tracking-widest uppercase px-3 mb-2">
          Cluster
        </p>
        {links.map(({ label, href, icon: Icon }) => {
          const isActive = pathname === href
          return (
            <Link
              key={href}
              href={href}
              className={`
                flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-all duration-100
                ${
                  isActive
                    ? 'bg-violet-pale text-violet-lt border border-violet/30'
                    : 'text-text-secondary hover:text-text-primary hover:bg-card border border-transparent'
                }
              `}
            >
              <Icon
                className={isActive ? 'text-violet-lt' : 'text-text-muted'}
                size={16}
              />
              {label}
              {isActive && (
                <span className="ml-auto w-1.5 h-1.5 rounded-full bg-violet-lt" />
              )}
            </Link>
          )
        })}
      </nav>

      {/* Footer */}
      <div className="mt-auto px-3 pt-4 border-t border-border-subtle">
        <p className="text-text-muted text-xs">minikube v1.0.0</p>
        <p className="text-text-muted text-xs mt-0.5 opacity-60">
          local cluster
        </p>
      </div>
    </aside>
  )
}
