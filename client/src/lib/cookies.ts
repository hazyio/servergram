type Cookiekeys = "theme"
const CookieKeys = {
  theme: "theme",
}
interface CookieOptions {
  path?: string
  sameSite?: "Lax" | "Strict" | "None"
  expires?: Date
  secure?: boolean
}
export function useCookies() {
  return {
    save: <T>(key: Cookiekeys, value: T, options?: CookieOptions) => {
      if (typeof document === "undefined") {
        console.error("You cannot use cookies in backend")
        return
      }
      let valueToSave = JSON.stringify(value)
      if (typeof value === "string") {
        valueToSave = value
      }
      const keyValue = CookieKeys[key]
      const sameSite = options?.sameSite ?? "Lax"
      const parts = [
        `${encodeURIComponent(keyValue)}=${encodeURIComponent(valueToSave)}`,
        `path=${options?.path ?? "/"}`,
        `SameSite=${sameSite}`,
      ]

      if (options?.expires !== undefined) {
        parts.push(`expires=${options?.expires.toUTCString()}`)
      }
      // SameSite=None requires Secure, or the browser drops the cookie
      if (options?.secure || sameSite === "None") {
        parts.push("Secure")
      }

      document.cookie = parts.join("; ")
    },
    get: <T>(key: Cookiekeys, d: T) => {
      const keyValue = CookieKeys[key]
      if (typeof document === "undefined") {
        console.error("You cannot use cookies in backend")
        return d
      }

      const cookieString = document.cookie
        .split("; ")
        .find((cookie) => cookie.startsWith(`${encodeURIComponent(keyValue)}=`))

      if (!cookieString) return d

      const rawValue = decodeURIComponent(
        cookieString.split("=").slice(1).join("=")
      )

      if (typeof d === "string") {
        return rawValue as T
      }

      try {
        return JSON.parse(rawValue) as T
      } catch (e) {
        console.error("Failed to parse cookie value", rawValue, e)
        return rawValue as T
      }
    },
    remove: (
      key: Cookiekeys,
      options: Partial<CookieOptions> = {
        path: "/",
        sameSite: "Lax",
      }
    ) => {
      if (typeof document === "undefined") {
        console.error("You cannot use cookies in backend")
        return
      }
      const keyValue = CookieKeys[key]

      document.cookie = `${encodeURIComponent(keyValue)}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=${options?.path}; SameSite=${options.sameSite}`
    },
    exists: (key: Cookiekeys): boolean => {
      return !!useCookies().get(key, false)
    },
  }
}
