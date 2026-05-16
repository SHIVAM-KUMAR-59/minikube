export type PodStatus = "PENDING" | "SCHEDULED" | "RUNNING" | "STOPPED"

export type Pod = {
    id: string
    name: string
    image: string
    status: PodStatus
    node_id: string
}

export type Node = {
    id: string
    name: string
    status: string
    last_heartbeat: string
}

export type Service = {
    id: string
    name: string
    port: string
    pods: string[]
}