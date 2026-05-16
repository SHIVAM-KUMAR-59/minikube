'use client'

import { useEffect, useState } from 'react'
import { getPods, createPod, deletePod } from '@/lib/api'
import { Pod } from '@/lib/types'
import StatusBadge from '@/components/StatusBadge'
import Modal from '@/components/Modal'
import EmptyState from '@/components/EmptyState'
import ConfirmDialog from '@/components/ConfirmDialog'
import { Box, Plus, Trash2, RefreshCw } from 'lucide-react'

export default function PodsPage() {
  const [pods, setPods] = useState<Pod[]>([])
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [confirmId, setConfirmId] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)
  const [name, setName] = useState('')
  const [image, setImage] = useState('')
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState('')

  const load = async (soft = false) => {
    if (soft) setRefreshing(true)
    else setLoading(true)
    try {
      setPods(await getPods())
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

  const handleCreate = async () => {
    if (!name.trim() || !image.trim()) {
      setError('Name and image are required.')
      return
    }
    setCreating(true)
    setError('')
    try {
      await createPod(name.trim(), image.trim())
      setModalOpen(false)
      setName('')
      setImage('')
      load(true)
    } catch {
      setError('Failed to create pod.')
    } finally {
      setCreating(false)
    }
  }

  const handleDelete = async () => {
    if (!confirmId) return
    setDeleting(true)
    try {
      await deletePod(confirmId)
      load(true)
    } finally {
      setDeleting(false)
      setConfirmId(null)
    }
  }

  const confirmPod = pods.find((p) => p.id === confirmId)

  return (
    <div className="px-6 md:px-10 py-8 max-w-6xl mx-auto">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
        <div>
          <h1 className="text-text-primary text-3xl font-semibold tracking-tight mb-1">
            Pods
          </h1>
          <p className="text-text-muted text-sm">
            {pods.length} pod{pods.length !== 1 ? 's' : ''} in cluster
          </p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => load(true)}
            disabled={refreshing}
            className="flex items-center gap-2 px-3.5 py-2 rounded-lg bg-card hover:bg-overlay border border-border-subtle text-text-secondary hover:text-text-primary text-sm font-medium transition-all duration-150 disabled:opacity-50"
          >
            <RefreshCw
              size={14}
              strokeWidth={1.5}
              className={refreshing ? 'animate-spin' : ''}
            />
            Refresh
          </button>
          <button
            onClick={() => setModalOpen(true)}
            className="flex items-center gap-2 px-3.5 py-2 rounded-lg bg-violet hover:bg-violet/90 text-white text-sm font-medium transition-all duration-150 shadow-lg shadow-violet/20"
          >
            <Plus size={14} strokeWidth={2} />
            Create pod
          </button>
        </div>
      </div>

      {/* Table */}
      {loading ? (
        <div className="flex items-center justify-center py-24">
          <RefreshCw
            size={18}
            strokeWidth={1.5}
            className="animate-spin text-text-muted"
          />
        </div>
      ) : pods.length === 0 ? (
        <EmptyState
          icon={Box}
          title="No pods yet"
          description="Create your first pod to get started. Pods are the smallest deployable units in your cluster."
          action={
            <button
              onClick={() => setModalOpen(true)}
              className="flex items-center gap-2 px-4 py-2 rounded-lg bg-violet hover:bg-violet/90 text-white text-sm font-medium transition-all shadow-lg shadow-violet/20"
            >
              <Plus size={14} strokeWidth={2} /> Create pod
            </button>
          }
        />
      ) : (
        <div className="bg-card border border-border-subtle rounded-xl overflow-hidden">
          {/* Table header */}
          <div className="grid grid-cols-12 gap-4 px-5 py-3 border-b border-border-subtle">
            <span className="col-span-3 text-text-muted text-xs font-medium uppercase tracking-wider">
              Name
            </span>
            <span className="col-span-3 text-text-muted text-xs font-medium uppercase tracking-wider">
              Image
            </span>
            <span className="col-span-2 text-text-muted text-xs font-medium uppercase tracking-wider">
              Status
            </span>
            <span className="col-span-3 text-text-muted text-xs font-medium uppercase tracking-wider">
              Node
            </span>
            <span className="col-span-1" />
          </div>
          {/* Rows */}
          <ul className="divide-y divide-border-subtle">
            {pods.map((pod) => (
              <li
                key={pod.id}
                className="grid grid-cols-12 gap-4 px-5 py-3.5 items-center hover:bg-overlay/40 transition-colors group"
              >
                <div className="col-span-3 flex items-center gap-2.5 min-w-0">
                  <div className="w-6 h-6 rounded-md bg-surface border border-border-subtle flex items-center justify-center shrink-0">
                    <Box
                      size={11}
                      strokeWidth={1.5}
                      className="text-text-muted"
                    />
                  </div>
                  <div className="min-w-0">
                    <p className="text-text-primary text-sm font-medium truncate">
                      {pod.name}
                    </p>
                    <p className="text-text-muted text-xs font-mono truncate">
                      {pod.id.slice(0, 8)}…
                    </p>
                  </div>
                </div>
                <p className="col-span-3 text-text-secondary text-sm font-mono truncate">
                  {pod.image}
                </p>
                <div className="col-span-2">
                  <StatusBadge status={pod.status} />
                </div>
                <p className="col-span-3 text-text-secondary text-sm font-mono truncate">
                  {pod.node_id || '—'}
                </p>
                <div className="col-span-1 flex justify-end">
                  <button
                    onClick={() => setConfirmId(pod.id)}
                    className="opacity-0 group-hover:opacity-100 p-1.5 rounded-md text-text-muted hover:text-failed-text hover:bg-failed-bg transition-all duration-150"
                    aria-label="Delete pod"
                  >
                    <Trash2 size={13} strokeWidth={1.5} />
                  </button>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Create modal */}
      <Modal
        open={modalOpen}
        onClose={() => {
          setModalOpen(false)
          setError('')
          setName('')
          setImage('')
        }}
        title="Create pod"
        description="Deploy a new pod to the cluster."
      >
        <div className="space-y-4">
          <div>
            <label className="text-text-secondary text-xs font-medium block mb-1.5">
              Pod name
            </label>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. nginx-pod"
              className="w-full bg-card border border-border-subtle rounded-lg px-3 py-2.5 text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:border-violet/50 transition-colors font-mono"
            />
          </div>
          <div>
            <label className="text-text-secondary text-xs font-medium block mb-1.5">
              Image
            </label>
            <input
              value={image}
              onChange={(e) => setImage(e.target.value)}
              placeholder="e.g. nginx:latest"
              className="w-full bg-card border border-border-subtle rounded-lg px-3 py-2.5 text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:border-violet/50 transition-colors font-mono"
            />
          </div>
          {error && <p className="text-failed-text text-xs">{error}</p>}
          <div className="flex justify-end gap-2 pt-1">
            <button
              onClick={() => {
                setModalOpen(false)
                setError('')
                setName('')
                setImage('')
              }}
              className="px-3.5 py-2 rounded-lg bg-card hover:bg-overlay border border-border-subtle text-text-secondary hover:text-text-primary text-sm font-medium transition-all"
            >
              Cancel
            </button>
            <button
              onClick={handleCreate}
              disabled={creating}
              className="flex items-center gap-2 px-3.5 py-2 rounded-lg bg-violet hover:bg-violet/90 disabled:opacity-50 text-white text-sm font-medium transition-all shadow-lg shadow-violet/20"
            >
              {creating && (
                <RefreshCw
                  size={13}
                  strokeWidth={1.5}
                  className="animate-spin"
                />
              )}
              Create
            </button>
          </div>
        </div>
      </Modal>

      {/* Confirm delete */}
      <ConfirmDialog
        open={!!confirmId}
        title="Delete pod"
        description={`Are you sure you want to delete "${confirmPod?.name}"? This will stop and remove the container.`}
        confirmLabel="Delete pod"
        loading={deleting}
        onConfirm={handleDelete}
        onClose={() => setConfirmId(null)}
      />
    </div>
  )
}
