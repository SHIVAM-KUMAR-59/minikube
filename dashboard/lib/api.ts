import { Node, Pod, Service } from "./types"

const BASE_API_URL = "http://localhost:8080"

export class ServerOfflineError extends Error {
  constructor() {
    super("Cannot connect to minikube server at localhost:8080")
    this.name = "ServerOfflineError"
  }
}

const fetchWithCheck = async (url: string, options?: RequestInit) => {
  try {
    const response = await fetch(url, options)
    return response
  } catch {
    throw new ServerOfflineError()
  }
}

export const getPods = async (): Promise<Pod[]> => {
  const response = await fetchWithCheck(`${BASE_API_URL}/pods`)
  return response.json()
}

export const getNodes = async (): Promise<Node[]> => {
  const response = await fetchWithCheck(`${BASE_API_URL}/nodes`)
  return response.json()
}

export const getServices = async (): Promise<Service[]> => {
  const response = await fetchWithCheck(`${BASE_API_URL}/services`)
  return response.json()
}

export const deletePod = async (id: string) => {
  const response = await fetchWithCheck(`${BASE_API_URL}/pods/${id}`, { method: "DELETE" })
  return response.json()
}

export const deleteNode = async (id: string) => {
  const response = await fetchWithCheck(`${BASE_API_URL}/nodes/${id}`, { method: "DELETE" })
  return response.json()
}

export const deleteService = async (id: string) => {
  const response = await fetchWithCheck(`${BASE_API_URL}/services/${id}`, { method: "DELETE" })
  return response.json()
}

export const createPod = async (name: string, image: string): Promise<Pod> => {
  const response = await fetchWithCheck(`${BASE_API_URL}/pods`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, image })
  })
  return response.json()
}

export const createService = async (name: string, port: string): Promise<Service> => {
  const response = await fetchWithCheck(`${BASE_API_URL}/services`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, port })
  })
  return response.json()
}