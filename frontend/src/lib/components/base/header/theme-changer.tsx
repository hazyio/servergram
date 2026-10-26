import { Check, Moon, Sun, SunMoon } from "lucide-react"
import { useEffect, useState } from "react"
import { Button } from "~/components/ui/button"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "~/components/ui/popover"
import { useLocalStorage } from "~/lib/local-storage"
import type { ThemeType } from "~/types"
function applyTheme(theme: ThemeType) {
  const resolvedTheme =
    theme === "system"
      ? window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light"
      : theme
  const htmlElement = document.documentElement
  htmlElement.classList.remove("dark", "light", "system")
  htmlElement.classList.add(resolvedTheme)
}

function setTheme(
  newTheme: ThemeType,
  updateTheme: (theme: ThemeType) => void
) {
  applyTheme(newTheme)
  useLocalStorage().save("theme", newTheme)
  updateTheme(newTheme)
}
function getThemeIcon(
  currentTheme: ThemeType,
  updateTheme: (theme: ThemeType) => void
) {
  switch (currentTheme) {
    case "dark":
      return (
        <Sun
          onClick={() => {
            setTheme("light", updateTheme)
          }}
          className={"h-6 w-6"}
        />
      )
    case "light":
      return (
        <Moon
          onClick={() => {
            setTheme("dark", updateTheme)
          }}
          className={"h-6 w-6"}
        />
      )
    case "system":
      return (
        <SunMoon
          onClick={() => {
            setTheme("system", updateTheme)
          }}
          className={"h-6 w-6"}
        />
      )
  }
}
export default function ThemeChanger() {
  const [currentTheme, setCurrentTheme] = useState<ThemeType>("dark")
  useEffect(() => {
    const savedTheme = useLocalStorage().get<ThemeType>("theme", "dark")
    setCurrentTheme(savedTheme)
    applyTheme(savedTheme)
  }, [])

  useEffect(() => {
    if (currentTheme !== "system") return

    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)")
    const updateSystemTheme = () => applyTheme("system")
    mediaQuery.addEventListener("change", updateSystemTheme)

    return () => mediaQuery.removeEventListener("change", updateSystemTheme)
  }, [currentTheme])

  return (
    <Popover>
      <PopoverTrigger
        render={<Button variant="outline" className={"bg-muted"} />}
      >
        {getThemeIcon(currentTheme, setCurrentTheme)}
      </PopoverTrigger>
      <PopoverContent className={"max-w-40 p-1"}>
        <div className={"flex flex-col gap-2"}>
          <div
            className={
              "flex cursor-pointer items-center gap-2 rounded-md px-2 py-2 hover:bg-muted/60"
            }
            onClick={() => setTheme("light", setCurrentTheme)}
          >
            <Sun className={"h-5 w-5"} />
            <div className={"grow"}>Light</div>
            {currentTheme === "light" && <Check className={"h-5 w-5"} />}
          </div>
          <div
            className={
              "flex cursor-pointer items-center gap-2 rounded-md px-2 py-2 hover:bg-muted/60"
            }
            onClick={() => setTheme("dark", setCurrentTheme)}
          >
            <Moon className={"h-5 w-5"} />
            <div className={"grow"}>Dark</div>
            {currentTheme === "dark" && <Check className={"h-5 w-5"} />}
          </div>
          <div
            className={
              "flex cursor-pointer items-center gap-2 rounded-md px-2 py-2 hover:bg-muted/60"
            }
            onClick={() => setTheme("system", setCurrentTheme)}
          >
            <SunMoon className={"h-5 w-5"} />
            <div className={"grow"}>System</div>{" "}
            {currentTheme === "system" && <Check className={"h-5 w-5"} />}
          </div>
        </div>
      </PopoverContent>
    </Popover>
  )
}
