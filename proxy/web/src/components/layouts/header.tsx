import { ApiStatusBanner } from "@/features/health/api-status-banner";
import { Link } from "react-router-dom";

export function Header() {
  return (
    <header className="site-header">
      <div className="wrapper">
        <div className="site-header__inner">
          <Link to="/" className="logo">
            <span className="logo__main">SentiProx</span>
            <span className="logo__sub">reverse-proxy threat console</span>
          </Link>
          <ApiStatusBanner />
        </div>
      </div>
    </header>
  );
}
