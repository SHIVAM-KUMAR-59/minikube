'use client'

import { useEffect } from 'react'
import { X } from 'lucide-react'

type Props = {
  title: string
  description?: string
  open: boolean
  onClose: () => void
  children: React.ReactNode
}

export default function Modal({ title, description, open, onClose, children }: Props) {
  useEffect(() => {
    const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    if (open) document.addEventListener('keydown', handler)
    return () => document.removeEventListener('keydown', handler)
  }, [open, onClose])

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center px-4">

      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-base/80 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Panel */}
      <div className="relative z-10 w-full max-w-md bg-surface border border-border-subtle rounded-xl shadow-2xl">

        {/* Header */}
        <div className="flex items-start justify-between px-6 py-5 border-b border-border-subtle">
          <div>
            <h2 className="text-text-primary text-sm font-semibold">{title}</h2>
            {description && (
              <p className="text-text-muted text-xs mt-0.5">{description}</p>
            )}
          </div>
          <button
            onClick={onClose}
            className="text-text-muted hover:text-text-primary transition-colors ml-4 mt-0.5"
            aria-label="Close"
          >
            <X size={16} strokeWidth={1.5} />
          </button>
        </div>

        {/* Body */}
        <div className="px-6 py-5">
          {children}
        </div>

      </div>
    </div>
  )
}