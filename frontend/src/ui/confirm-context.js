import { createContext, useContext } from "react";

export const ConfirmContext = createContext(null);

export function useConfirm() {
  const value = useContext(ConfirmContext);
  if (!value) throw new Error("useConfirm must be used within ConfirmProvider");
  return value;
}
