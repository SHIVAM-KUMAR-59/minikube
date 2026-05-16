'use client'

import { getPods, getNodes, getServices } from '@/lib/api'
import StatCard from '@/components/StatCard'
import StatusBadge from '@/components/StatusBadge'
import { Box, Server, Network, Activity, Clock } from 'lucide-react'
import { timeAgo } from '@/lib/utils'
import { useEffect, useRef, useState } from 'react'
import { useToast } from '@/context/ToastContext'
import ServerOfflineBanner from '@/components/ServerOfflineBanner'

export default function OverviewPage() {
  const [pods, setPods] = useState<Awaited<ReturnType<typeof getPods>>>([])
  const [nodes, setNodes] = useState<Awaited<ReturnType<typeof getNodes>>>([])
  const [services, setServices] = useState<
    Awaited<ReturnType<typeof getServices>>
  >([])
  const [error, setError] = useState<Error | null>(null)
  const { error: showError } = useToast()

  const offlineToastShown = useRef(false)

  const fetchData = async () => {
    try {
      const [p, n, s] = await Promise.all([
        getPods(),
        getNodes(),
        getServices(),
      ])
      setPods(p ?? [])
      setNodes(n ?? [])
      setServices(s ?? [])
      setError(null)
      offlineToastShown.current = false
    } catch (err) {
      setError(err as Error)
      if (!offlineToastShown.current) {
        showError('Failed to fetch cluster data')
        offlineToastShown.current = true
      }
    }
  }

  const runningPods = pods.filter((p) => p.status === 'RUNNING').length

  const pendingPods = pods.filter(
    (p) => p.status === 'PENDING' || p.status === 'SCHEDULED',
  ).length

  const healthyNodes = nodes.filter((n) => n.status === 'READY').length

  useEffect(() => {
    const init = async () => fetchData()
    void init()

    const interval = setInterval(() => {
      void fetchData()
    }, 5000)

    return () => clearInterval(interval)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <div>
      <ServerOfflineBanner error={error} />

      <div className="px-6 md:px-10 py-8 max-w-6xl mx-auto">
        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-running opacity-60" />
              <span className="relative inline-flex rounded-full h-2 w-2 bg-running" />
            </span>
            <span className="text-running-text text-xs font-mono">
              cluster online
            </span>
          </div>
          <h1 className="text-text-primary text-3xl font-semibold tracking-tight mb-2">
            Overview
          </h1>
          <p className="text-text-muted text-sm">
            Monitor your cluster resources in real time.
          </p>
        </div>

        {/* Stat cards */}
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-10">
          <StatCard
            label="Total pods"
            value={pods.length}
            icon={Box}
            sub={`${runningPods} running`}
          />
          <StatCard
            label="Total nodes"
            value={nodes.length}
            icon={Server}
            sub={`${healthyNodes} healthy`}
          />
          <StatCard
            label="Total services"
            value={services.length}
            icon={Network}
            sub="active"
          />
          <StatCard
            label="Pending pods"
            value={pendingPods}
            icon={Activity}
            sub="awaiting schedule"
          />
        </div>

        {/* Two col layout */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Recent pods */}
          <section>
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-text-primary text-sm font-semibold">
                Recent pods
              </h2>
              <a
                href="/pods"
                className="text-violet-lt text-xs hover:underline underline-offset-4"
              >
                View all →
              </a>
            </div>
            <div className="bg-card border border-border-subtle rounded-xl overflow-hidden">
              {pods.length === 0 ? (
                <p className="text-text-muted text-xs text-center py-10">
                  No pods found.
                </p>
              ) : (
                <ul className="divide-y divide-border-subtle">
                  {pods.slice(0, 5).map((pod) => (
                    <li
                      key={pod.id}
                      className="flex items-center justify-between px-4 py-3 hover:bg-overlay/40 transition-colors"
                    >
                      <div className="min-w-0">
                        <p className="text-text-primary text-sm font-medium truncate">
                          {pod.name}
                        </p>
                        <p className="text-text-muted text-xs font-mono truncate">
                          {pod.image}
                        </p>
                      </div>
                      <StatusBadge status={pod.status} />
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </section>

          {/* Nodes */}
          <section>
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-text-primary text-sm font-semibold">Nodes</h2>
              <a
                href="/nodes"
                className="text-violet-lt text-xs hover:underline underline-offset-4"
              >
                View all →
              </a>
            </div>
            <div className="bg-card border border-border-subtle rounded-xl overflow-hidden">
              {nodes.length === 0 ? (
                <p className="text-text-muted text-xs text-center py-10">
                  No nodes found.
                </p>
              ) : (
                <ul className="divide-y divide-border-subtle">
                  {nodes.map((node) => (
                    <li
                      key={node.id}
                      className="flex items-center justify-between px-4 py-3 hover:bg-overlay/40 transition-colors"
                    >
                      <div className="flex items-center gap-3 min-w-0">
                        <div className="w-7 h-7 rounded-md bg-surface border border-border-subtle flex items-center justify-center shrink-0">
                          <Server
                            size={13}
                            strokeWidth={1.5}
                            className="text-text-muted"
                          />
                        </div>
                        <div className="min-w-0">
                          <p className="text-text-primary text-sm font-medium truncate">
                            {node.name}
                          </p>
                          <div className="flex items-center gap-1 text-text-muted text-xs">
                            <Clock size={10} strokeWidth={1.5} />
                            <span className="font-mono">
                              {/* {new Date(node.last_heartbeat).toLocaleTimeString()}
                               */}
                              {timeAgo(node.last_heartbeat)}
                            </span>
                          </div>
                        </div>
                      </div>
                      <StatusBadge
                        status={node.status as 'READY' | 'NOT_READY'}
                      />
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </section>

          {/* Services */}
          <section className="lg:col-span-2">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-text-primary text-sm font-semibold">
                Services
              </h2>
              <a
                href="/services"
                className="text-violet-lt text-xs hover:underline underline-offset-4"
              >
                View all →
              </a>
            </div>
            <div className="bg-card border border-border-subtle rounded-xl overflow-hidden">
              {services.length === 0 ? (
                <p className="text-text-muted text-xs text-center py-10">
                  No services found.
                </p>
              ) : (
                <ul className="divide-y divide-border-subtle">
                  {services.map((svc) => (
                    <li
                      key={svc.id}
                      className="flex items-center justify-between px-4 py-3 hover:bg-overlay/40 transition-colors"
                    >
                      <div className="flex items-center gap-3 min-w-0">
                        <div className="w-7 h-7 rounded-md bg-surface border border-border-subtle flex items-center justify-center shrink-0">
                          <Network
                            size={13}
                            strokeWidth={1.5}
                            className="text-text-muted"
                          />
                        </div>
                        <div className="min-w-0">
                          <p className="text-text-primary text-sm font-medium truncate">
                            {svc.name}
                          </p>
                          <p className="text-text-muted text-xs font-mono">
                            port {svc.port}
                          </p>
                        </div>
                      </div>
                      <span className="text-text-muted text-xs font-mono shrink-0">
                        {svc.pods?.length} pod
                        {svc.pods?.length !== 1 ? 's' : ''}
                      </span>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </section>
        </div>
      </div>
    </div>
  )
}
