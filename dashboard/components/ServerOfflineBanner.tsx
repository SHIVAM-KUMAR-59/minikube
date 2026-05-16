'use client'

import { ServerOfflineError } from '@/lib/api'
import { WifiOff, Terminal } from 'lucide-react'

export default function ServerOfflineBanner({ error }: { error: Error | null }) {
  if (!error || !(error instanceof ServerOfflineError)) return null

  return (
    <div className="mx-6 md:mx-10 mt-6 rounded-xl border border-failed-border bg-failed-bg px-5 py-4">
      <div className="flex items-start gap-3">
        <WifiOff size={16} strokeWidth={1.5} className="text-failed-text mt-0.5 shrink-0" />
        <div>
          <p className="text-failed-text text-sm font-semibold mb-1">
            Cannot connect to minikube server
          </p>
          <p className="text-text-secondary text-xs leading-relaxed mb-3">
            The dashboard cannot reach the server at <span className="font-mono text-text-primary">localhost:8080</span>. Make sure your cluster is running.
          </p>
          <div className="bg-card border border-border-subtle rounded-lg px-4 py-3">
            <div className="flex items-center gap-1.5 mb-2">
              <Terminal size={11} strokeWidth={1.5} className="text-text-muted" />
              <span className="text-text-muted text-xs">Run this to start the cluster</span>
            </div>
            <p className="font-mono text-xs text-violet-lt">
              minik cluster start --workers 2
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}