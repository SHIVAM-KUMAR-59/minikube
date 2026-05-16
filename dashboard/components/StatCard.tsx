import { LucideIcon } from 'lucide-react'

type Props = {
  label: string
  value: number | string
  icon: LucideIcon
  sub?: string
}

export default function StatCard({ label, value, icon: Icon, sub }: Props) {
  return (
    <div className="bg-card border border-border-subtle rounded-xl px-5 py-4 flex items-start gap-4">
      <div className="w-9 h-9 rounded-lg bg-violet/10 border border-violet/20 flex items-center justify-center shrink-0">
        <Icon size={16} strokeWidth={1.5} className="text-violet-lt" />
      </div>
      <div className="min-w-0">
        <p className="text-text-muted text-xs font-medium mb-1">{label}</p>
        <p className="text-text-primary text-2xl font-semibold leading-none">{value}</p>
        {sub && <p className="text-text-muted text-xs mt-1.5">{sub}</p>}
      </div>
    </div>
  )
}