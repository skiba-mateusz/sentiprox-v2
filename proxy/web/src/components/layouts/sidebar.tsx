import { NavLink } from "react-router-dom";
import { LayoutDashboard, FileText, Ban } from "lucide-react";

const navlinks = [
  {
    name: "Summary",
    to: "/",
    icon: <LayoutDashboard />,
  },
  {
    name: "Logs",
    to: "/logs",
    icon: <FileText />,
  },
  {
    name: "Blocklist",
    to: "/blocklist",
    icon: <Ban />,
  },
];

export function Sidebar() {
  return (
    <aside className="sidebar">
      <nav className="nav flow">
        <div className="heading fs-h4">Menu</div>
        <ul className="nav__list">
          {navlinks.map((navlink) => (
            <li className="nav__item">
              <NavLink
                to={navlink.to}
                className={(props) => {
                  return props.isActive
                    ? "nav__link nav__link--active"
                    : "nav__link";
                }}
              >
                {navlink.icon}
                {navlink.name}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
}
