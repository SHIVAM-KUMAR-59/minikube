import { PodStatus } from '@/lib/types'

type Status = PodStatus | 'READY' | 'NOT_READY'

const config: Record<Status, { label: string; dot: string; text: string; bg: string; border: string }> = {
  RUNNING: {
    label: 'Running',
    dot:    'bg-running',
    text:   'text-running-text',
    bg:     'bg-running-bg',
    border: 'border-running-border',
  },
  PENDING: {
    label: 'Pending',
    dot:    'bg-pending',
    text:   'text-pending-text',
    bg:     'bg-pending-bg',
    border: 'border-pending-border',
  },
  SCHEDULED: {
    label: 'Scheduled',
    dot:    'bg-pending',
    text:   'text-pending-text',
    bg:     'bg-pending-bg',
    border: 'border-pending-border',
  },
  STOPPED: {
    label: 'Stopped',
    dot:    'bg-failed',
    text:   'text-failed-text',
    bg:     'bg-failed-bg',
    border: 'border-failed-border',
  },
  READY: {
    label: 'Ready',
    dot:    'bg-running',
    text:   'text-running-text',
    bg:     'bg-running-bg',
    border: 'border-running-border',
  },
  NOT_READY: {
    label: 'Not Ready',
    dot:    'bg-failed',
    text:   'text-failed-text',
    bg:     'bg-failed-bg',
    border: 'border-failed-border',
  },
}

export default function StatusBadge({ status }: { status: Status }) {
  const c = config[status] ?? {
    label: status,
    dot:    'bg-unknown',
    text:   'text-unknown-text',
    bg:     'bg-unknown-bg',
    border: 'border-unknown-border',
  }

  return (
    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border ${c.bg} ${c.text} ${c.border}`}>
      <span className={`w-1.5 h-1.5 rounded-full ${c.dot}`} />
      {c.label}
    </span>
  )
}