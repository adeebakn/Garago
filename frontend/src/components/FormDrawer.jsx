import { useEffect } from "react";
import { createPortal } from "react-dom";
import { useOrganization } from "../organization/organization-context";
import "./FormDrawer.css";

export default function FormDrawer({ open, title, onClose, children, backdropClassName = "", drawerClassName = "" }) {
  const { organization } = useOrganization();
  useEffect(() => {
    if (!open) return undefined;
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => { document.body.style.overflow = previousOverflow; };
  }, [open]);
  if (!open) return null;
  return createPortal(
    <div className={`form-drawer-backdrop ${backdropClassName}`.trim()} role="presentation" onMouseDown={onClose}>
      <aside className={`form-drawer ${drawerClassName}`.trim()} role="dialog" aria-modal="true" aria-label={title} onMouseDown={(event) => event.stopPropagation()}>
        <header className="form-drawer-header"><div><p>{organization.name || "Garago"}</p><h2>{title}</h2></div><button type="button" className="form-drawer-close" onClick={onClose} aria-label="Close form">×</button></header>
        <div className="form-drawer-body">{children}</div>
      </aside>
    </div>,
    document.body,
  );
}
