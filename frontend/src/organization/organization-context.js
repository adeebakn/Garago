import { createContext, useContext } from "react";

export const fallbackOrganization = { name: "Garago", address: "", phone: "", email: "" };
export const OrganizationContext = createContext(null);

export function useOrganization() {
  const value = useContext(OrganizationContext);
  if (!value) throw new Error("useOrganization must be used within OrganizationProvider");
  return value;
}
