import { createContext, useContext, useEffect, useState } from "react";
import api from "../api";

const fallbackOrganization = { name: "Garago", address: "", phone: "", email: "" };
const OrganizationContext = createContext({ organization: fallbackOrganization, refreshOrganization: async () => {} });

export function OrganizationProvider({ children }) {
  const [organization, setOrganization] = useState(fallbackOrganization);
  const refreshOrganization = async () => {
    try {
      const { data } = await api.get("/organization");
      setOrganization(data);
      return data;
    } catch {
      return fallbackOrganization;
    }
  };
  useEffect(() => { refreshOrganization(); }, []);
  return <OrganizationContext.Provider value={{ organization, setOrganization, refreshOrganization }}>{children}</OrganizationContext.Provider>;
}

export const useOrganization = () => useContext(OrganizationContext);
