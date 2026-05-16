/* eslint-disable react-hooks/exhaustive-deps */
'use client'

import { useEffect, useRef, useState } from 'react'
import { getNodes, deleteNode } from '@/lib/api'
import { Node } from '@/lib/types'
import StatusBadge from '@/components/StatusBadge'
import EmptyState from '@/components/EmptyState'
import ConfirmDialog from '@/components/ConfirmDialog'
import { Server, Trash2, RefreshCw, Clock } from 'lucide-react'
import { timeAgo } from '@/lib/utils'
import { useToast } from '@/context/ToastContext'
import ServerOfflineBanner from '@/components/ServerOfflineBanner'

export default function NodesPage() {
  const [nodes, setNodes] = useState<Node[]>([])
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [confirmId, setConfirmId] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)
  const [offlineError, setOfflineError] = useState<Error | null>(null)

  const { error: showError, success: showSuccess } = useToast()

  const offlineToastShown = useRef(false)

  const load = async (soft = false) => {
    if (soft) setRefreshing(true)
    else setLoading(true)
    try {
      setNodes(await getNodes())
      setOfflineError(null)
      offlineToastShown.current = false
    } catch (err) {
      setOfflineError(err as Error)
      if (!offlineToastShown.current) {
        showError('Failed to load nodes')
        offlineToastShown.current = true
      }
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }

  useEffect(() => {
    const init = async () => {
      await load()
    }
    void init()
  }, [])

  useEffect(() => {
    const interval = setInterval(() => {
      load()
    }, 5000)
    return () => clearInterval(interval)
  }, [])

  const handleDelete = async () => {
    if (!confirmId) return
    setDeleting(true)
    try {
      await deleteNode(confirmId)
      showSuccess('Node deleted successfully!')
      load(true)
    } catch (err) {
      setOfflineError(err as Error)
      showError('Failed to delete services')
    } finally {
      setDeleting(false)
      setConfirmId(null)
    }
  }

  const confirmNode = nodes.find((n) => n.id === confirmId)
  const healthyCount = nodes.filter((n) => n.status === 'READY').length

  return (
    <div>
      <ServerOfflineBanner error={offlineError} />

      <div className="px-6 md:px-10 py-8 max-w-6xl mx-auto">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
          <div>
            <h1 className="text-text-primary text-3xl font-semibold tracking-tight mb-1">
              Nodes
            </h1>
            <p className="text-text-muted text-sm">
              {healthyCount} of {nodes.length} node
              {nodes.length !== 1 ? 's' : ''} healthy
            </p>
          </div>
          <button
            onClick={() => load(true)}
            disabled={refreshing}
            className="flex items-center gap-2 px-3.5 py-2 rounded-lg bg-card hover:bg-overlay border border-border-subtle text-text-secondary hover:text-text-primary text-sm font-medium transition-all duration-150 disabled:opacity-50 self-start sm:self-auto"
          >
            <RefreshCw
              size={14}
              strokeWidth={1.5}
              className={refreshing ? 'animate-spin' : ''}
            />
            Refresh
          </button>
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-24">
            <RefreshCw
              size={18}
              strokeWidth={1.5}
              className="animate-spin text-text-muted"
            />
          </div>
        ) : nodes.length === 0 ? (
          <EmptyState
            icon={Server}
            title="No nodes registered"
            description="Nodes register automatically when a worker process starts. Run a worker to add a node to the cluster."
          />
        ) : (
          <div className="bg-card border border-border-subtle rounded-xl overflow-hidden">
            <div className="grid grid-cols-12 gap-4 px-5 py-3 border-b border-border-subtle">
              <span className="col-span-3 text-text-muted text-xs font-medium uppercase tracking-wider">
                Name
              </span>
              <span className="col-span-4 text-text-muted text-xs font-medium uppercase tracking-wider">
                ID
              </span>
              <span className="col-span-2 text-text-muted text-xs font-medium uppercase tracking-wider">
                Status
              </span>
              <span className="col-span-2 text-text-muted text-xs font-medium uppercase tracking-wider">
                Heartbeat
              </span>
              <span className="col-span-1" />
            </div>
            <ul className="divide-y divide-border-subtle">
              {nodes.map((node) => (
                <li
                  key={node.id}
                  className="grid grid-cols-12 gap-4 px-5 py-3.5 items-center hover:bg-overlay/40 transition-colors group"
                >
                  <div className="col-span-3 flex items-center gap-2.5 min-w-0">
                    <div className="w-6 h-6 rounded-md bg-surface border border-border-subtle flex items-center justify-center shrink-0">
                      <Server
                        size={11}
                        strokeWidth={1.5}
                        className="text-text-muted"
                      />
                    </div>
                    <p className="text-text-primary text-sm font-medium truncate">
                      {node.name}
                    </p>
                  </div>
                  <p className="col-span-4 text-text-muted text-xs font-mono truncate">
                    {node.id}
                  </p>
                  <div className="col-span-2">
                    <StatusBadge
                      status={node.status as 'READY' | 'NOT_READY'}
                    />
                  </div>
                  <div className="col-span-2 flex items-center gap-1 text-text-muted text-xs font-mono">
                    <Clock size={11} strokeWidth={1.5} className="shrink-0" />
                    {timeAgo(node.last_heartbeat)}
                  </div>
                  <div className="col-span-1 flex justify-end">
                    <button
                      onClick={() => setConfirmId(node.id)}
                      className="opacity-0 group-hover:opacity-100 p-1.5 rounded-md text-text-muted hover:text-failed-text hover:bg-failed-bg transition-all duration-150"
                      aria-label="Delete node"
                    >
                      <Trash2 size={13} strokeWidth={1.5} />
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        )}

        <ConfirmDialog
          open={!!confirmId}
          title="Delete node"
          description={`Are you sure you want to delete "${confirmNode?.name}"? Any pods running on this node may be affected.`}
          confirmLabel="Delete node"
          loading={deleting}
          onConfirm={handleDelete}
          onClose={() => setConfirmId(null)}
        />
      </div>
    </div>
  )
}
