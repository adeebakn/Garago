import { useCallback, useRef, useState } from "react";
import "./confirm.css";
import { ConfirmContext } from "./confirm-context";

export function ConfirmProvider({ children }) {
  const [dialog, setDialog] = useState(null);
  const resolver = useRef(null);

	const confirm = useCallback((options) => new Promise((resolve) => {
    resolver.current = resolve;
    setDialog({ title: "Please confirm", message: "Are you sure you want to continue?", destructive: false, confirmLabel: "Continue", ...options });
  }), []);

  const close = (result) => {
    resolver.current?.(result);
    resolver.current = null;
    setDialog(null);
  };

  return <ConfirmContext.Provider value={confirm}>{children}{dialog && <div className="confirm-backdrop" role="presentation" onMouseDown={() => close(false)}><section className="confirm-sheet" role="dialog" aria-modal="true" aria-labelledby="confirm-title" onMouseDown={(event) => event.stopPropagation()}><div className={`confirm-icon ${dialog.destructive ? "destructive" : ""}`}>{dialog.destructive ? "!" : "?"}</div><h2 id="confirm-title">{dialog.title}</h2><p>{dialog.message}</p><div className="confirm-actions"><button type="button" onClick={() => close(false)}>Cancel</button><button type="button" className={dialog.destructive ? "danger-button" : ""} onClick={() => close(true)}>{dialog.confirmLabel}</button></div></section></div>}</ConfirmContext.Provider>;
}
