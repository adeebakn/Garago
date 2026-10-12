import "./Sidebar.css";
import { NavLink } from "react-router-dom";

const links = [
  ["/", "Dashboard", "⌂"],
  ["/customers", "Customers", "◉"],
  ["/vehicles", "Vehicles", "▣"],
  ["/services", "Services", "✦"],
  ["/quotations", "Quotations", "▤"],
  ["/work-orders", "Work Orders", "✓"],
  ["/invoices", "Invoices", "▧"],
  ["/payments", "Payments", "₹"],
  ["/settings", "Settings", "⚙"],
];

function Sidebar() {
  return (
    <aside className="sidebar">
      <p className="nav-label">Workspace</p>
      <nav>{links.map(([to, label, icon]) => <NavLink key={to} to={to} end={to === "/"} className={({ isActive }) => isActive ? "active" : ""}><span className="nav-icon">{icon}</span>{label}</NavLink>)}</nav>
      <div className="sidebar-footer"><span className="status-dot" />Admin</div>
    </aside>
  );
}

export default Sidebar;
