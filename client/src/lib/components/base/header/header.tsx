import { href, Link } from "react-router"
import ThemeChanger from "./theme-changer"

export default function Header() {
  return (
    <header
      className={"flex h-15 items-center gap-3 border-b bg-background px-3"}
    >
      <Link to={href("/")} className="text-2xl">
        Muster
      </Link>

      <div className={"grow"}></div>
      <Link to={href("/account/login")}>Login</Link>
      <ThemeChanger />
    </header>
  )
}
