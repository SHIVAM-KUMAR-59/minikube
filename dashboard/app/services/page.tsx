'use client'

import { useEffect, useState } from 'react'
import { getServices, createService, deleteService } from '@/lib/api'
import { Service } from '@/lib/types'
import Modal from '@/components/Modal'
import EmptyState from '@/components/EmptyState'
import ConfirmDialog from '@/components/ConfirmDialog'
import { Network, Plus, Trash2, RefreshCw } from 'lucide-react'

export default function ServicesPage() {
  const [services, setServices] = useState<Service[]>([])
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [confirmId, setConfirmId] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)
  const [name, setName] = useState('')
  const [port, setPort] = useState('')
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState('')

  const load = async (soft = false) => {
    if (soft) setRefreshing(true)
    else setLoading(true)
    try {
      setServices(await getServices())
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
    if (!name.trim() || !port.trim()) {
      setError('Name and port are required.')
      return
    }
    setCreating(true)
    setError('')
    try {
      await createService(name.trim(), port.trim())
      setModalOpen(false)
      setName('')
      setPort('')
      load(true)
    } catch {
      setError('Failed to create service.')
    } finally {
      setCreating(false)
    }
  }

  const handleDelete = async () => {
    if (!confirmId) return
    setDeleting(true)
    try {
      await deleteService(confirmId)
      load(true)
    } finally {
      setDeleting(false)
      setConfirmId(null)
    }
  }

  const confirmSvc = services.find((s) => s.id === confirmId)

  return (
    <div className="px-6 md:px-10 py-8 max-w-6xl mx-auto">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
        <div>
          <h1 className="text-text-primary text-3xl font-semibold tracking-tight mb-1">
            Services
          </h1>
          <p className="text-text-muted text-sm">
            {services.length} service{services.length !== 1 ? 's' : ''}{' '}
            registered
          </p>
        </div>
        <div className="flex items-center gap-2 self-start sm:self-auto">
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
            className="flex items-center gap-2 px-3.5 py-2 rounded-lg bg-violet hover:bg-violet/90 text-white text-sm font-medium transition-all shadow-lg shadow-violet/20"
          >
            <Plus size={14} strokeWidth={2} />
            Create service
          </button>
        </div>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-24">
          <RefreshCw
            size={18}
            strokeWidth={1.5}
            className="animate-spin text-text-muted"
          />
        </div>
      ) : services.length === 0 ? (
        <EmptyState
          icon={Network}
          title="No services yet"
          description="Services expose your pods to the network. Create one to get started."
          action={
            <button
              onClick={() => setModalOpen(true)}
              className="flex items-center gap-2 px-4 py-2 rounded-lg bg-violet hover:bg-violet/90 text-white text-sm font-medium transition-all shadow-lg shadow-violet/20"
            >
              <Plus size={14} strokeWidth={2} /> Create service
            </button>
          }
        />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {services.map((svc) => (
            <div
              key={svc.id}
              className="bg-card border border-border-subtle rounded-xl p-5 hover:border-border-default transition-colors group relative"
            >
              {/* Delete btn */}
              <button
                onClick={() => setConfirmId(svc.id)}
                className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 p-1.5 rounded-md text-text-muted hover:text-failed-text hover:bg-failed-bg transition-all duration-150"
                aria-label="Delete service"
              >
                <Trash2 size={13} strokeWidth={1.5} />
              </button>

              {/* Icon + name */}
              <div className="flex items-center gap-3 mb-4">
                <div className="w-9 h-9 rounded-lg bg-violet/10 border border-violet/20 flex items-center justify-center shrink-0">
                  <Network
                    size={15}
                    strokeWidth={1.5}
                    className="text-violet-lt"
                  />
                </div>
                <div className="min-w-0">
                  <p className="text-text-primary text-sm font-semibold truncate">
                    {svc.name}
                  </p>
                  <p className="text-text-muted text-xs font-mono truncate">
                    {svc.id.slice(0, 8)}…
                  </p>
                </div>
              </div>

              {/* Port */}
              <div className="flex items-center justify-between mb-3 py-2.5 px-3 bg-surface rounded-lg border border-border-subtle">
                <span className="text-text-muted text-xs">Port</span>
                <span className="text-text-primary text-xs font-mono font-medium">
                  {svc.port}
                </span>
              </div>

              {/* Pods */}
              <div>
                <p className="text-text-muted text-xs font-medium mb-2">
                  {svc.pods?.length || 0} pod{svc.pods?.length !== 1 ? 's' : ''}
                </p>
                {svc.pods?.length === 0 ? (
                  <p className="text-text-muted text-xs italic">
                    No pods attached
                  </p>
                ) : (
                  <div className="flex flex-wrap gap-1.5">
                    {svc.pods?.map((podId) => (
                      <span
                        key={podId}
                        className="text-xs font-mono bg-surface border border-border-subtle text-text-secondary px-2 py-0.5 rounded-md truncate max-w-full"
                      >
                        {podId.slice(0, 8)}…
                      </span>
                    ))}
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Create modal */}
      <Modal
        open={modalOpen}
        onClose={() => {
          setModalOpen(false)
          setError('')
          setName('')
          setPort('')
        }}
        title="Create service"
        description="Expose pods through a named service."
      >
        <div className="space-y-4">
          <div>
            <label className="text-text-secondary text-xs font-medium block mb-1.5">
              Service name
            </label>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. nginx-service"
              className="w-full bg-card border border-border-subtle rounded-lg px-3 py-2.5 text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:border-violet/50 transition-colors font-mono"
            />
          </div>
          <div>
            <label className="text-text-secondary text-xs font-medium block mb-1.5">
              Port
            </label>
            <input
              value={port}
              onChange={(e) => setPort(e.target.value)}
              placeholder="e.g. 8080"
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
                setPort('')
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

      <ConfirmDialog
        open={!!confirmId}
        title="Delete service"
        description={`Are you sure you want to delete "${confirmSvc?.name}"? This will remove the service and all its routing rules.`}
        confirmLabel="Delete service"
        loading={deleting}
        onConfirm={handleDelete}
        onClose={() => setConfirmId(null)}
      />
    </div>
  )
}
