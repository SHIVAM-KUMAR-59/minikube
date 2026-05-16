import { LucideIcon } from 'lucide-react'

type Props = {
  icon: LucideIcon
  title: string
  description: string
  action?: React.ReactNode
}

export default function EmptyState({
  icon: Icon,
  title,
  description,
  action,
}: Props) {
  return (
    <div className="flex flex-col items-center justify-center py-20 px-6 text-center">
      <div className="w-12 h-12 rounded-xl bg-card border border-border-subtle flex items-center justify-center mb-4">
        <Icon size={20} strokeWidth={1.5} className="text-text-muted" />
      </div>
      <p className="text-text-primary text-sm font-medium mb-1">{title}</p>
      <p className="text-text-muted text-xs max-w-xs leading-relaxed">
        {description}
      </p>
      {action && <div className="mt-5">{action}</div>}
    </div>
  )
}
