import axios from "axios"
import { Logger } from "~/lib/log"
import { requestTryCatch } from "~/lib/try-catch"
import { StatusCodes } from "~/types/statusCodes"

import { href, redirect } from "react-router"

const API_URL = import.meta.env.VITE_API_URL
if (!API_URL) {
  throw new Error("VITE_API_URL must be provided to run the application")
}

export type ApiErrorResponse = {
  statusCode: StatusCodes
  message: string
}

interface Headers {
  "Content-Type": "application/json"
  Authorization?: string
  "Access-Control-Allow-Origin"?: string
  "Accept-Language"?: string
}

interface ApiRequestProps {
  url: string
  method?: "GET" | "POST"
  headers?: Partial<Headers>
  body?: Record<string, unknown>
  withCredentials?: boolean
}
export function skipAuth(): boolean {
  const path = window.location.pathname
  return (
    path == href("/account/login") || path == href("/account/create-account")
  )
}

function getServerUrl(): string {
  return API_URL
}

function getTimeout(): number {
  if (import.meta.env.DEV) {
    return 5000 // 5 seconds in dev
  }
  return 10000 // 10 seconds
}

export class BaseRequest {
  base_path?: string
  signal?: AbortSignal
  constructor(options?: { signal?: AbortSignal }) {
    this.signal = options?.signal
  }

  private async buildHeaders(
    headers: Partial<Headers>
  ): Promise<Partial<Headers> | StatusCodes> {
    const new_headers: Partial<Headers> = {
      ...headers,
      "Content-Type": "application/json",
    }

    return new_headers
  }

  async makeJsonRequest<T>({
    url,
    headers = {},
    body,
    method,
    withCredentials,
  }: ApiRequestProps): Promise<T> {
    const new_headers = await this.buildHeaders(headers)
    if (typeof new_headers === "number") {
      return Promise.reject({
        statusCode: new_headers,
        message: "Authorization failed",
      })
    }
    let serverUrl = getServerUrl()
    let signal = this.signal
    let controller: AbortController
    if (this.base_path) {
      // If a base_path is defined, append it to the server URL
      serverUrl = `${serverUrl}${this.base_path}`
    } else {
      serverUrl = getServerUrl()
    }

    if (signal === undefined) {
      // If no signal is provided, create a new AbortController for this request
      controller = new AbortController()
      signal = controller.signal
    }
    try {
      if (import.meta.env.DEV) {
        Logger.info("Making API request:", {
          url: `${serverUrl}${url}`,
          method,
          headers: new_headers,
          body,
        })
      }
      const timeoutId = setTimeout(() => {
        controller?.abort()
      }, getTimeout())

      const buildRequest = () => {
        if (method === "POST") {
          return axios.post(`${serverUrl}${url}`, body, {
            headers: new_headers,
            withCredentials: withCredentials,
            signal: signal,
          })
        }
        return axios.get(`${serverUrl}${url}`, {
          headers: new_headers,
          data: body,
          withCredentials: withCredentials,
          signal: signal,
        })
      }
      const res = await buildRequest()
      clearTimeout(timeoutId)

      if (res.status === StatusCodes.Ok) {
        return Promise.resolve(res.data.data)
      } else {
        // Rejects the promise for non-OK status codes from successful HTTP requests
        return Promise.reject({
          statusCode: res.status as StatusCodes,
          message: res.data.data,
        })
      }
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const r = error.response
        if (r?.status === 401 && withCredentials && !skipAuth()) {
          Logger.warn("Redirecting")
          redirect(href("/account/login"))
        }
        if (r && r.data) {
          // Rejects the promise with structured error from axios response
          return Promise.reject({
            statusCode: r.status as StatusCodes,
            message: r.data.error || "An unexpected error occurred",
          })
        }
        if (error.code === "ERR_CANCELED") {
          // Rejects the promise for request cancellations
          return Promise.reject({
            statusCode: StatusCodes.Cancelled,
            message: "The request was cancelled or timed out",
          })
        }
        // Rejects the promise for other axios errors
        return Promise.reject({
          statusCode: error.response?.status || 500,
          message: error.message || "An unexpected error occurred",
        })
      }
      // Rejects the promise for non-axios errors
      return Promise.reject({
        statusCode: 500,
        message: "An unexpected error occurred",
      })
    }
  }

  async useTryCatchJsonRequest<T>({
    url,
    headers = {
      "Content-Type": "application/json",
    },
    body,
    method,
    withCredentials,
  }: ApiRequestProps) {
    return requestTryCatch<T>(
      this.makeJsonRequest({
        url,
        headers,
        body,
        method,
        withCredentials,
      })
    )
  }
}
