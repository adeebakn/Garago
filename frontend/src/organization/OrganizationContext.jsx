import { useCallback, useMemo, useState } from "react";
import api from "../api";
import { useInitialLoad } from "../hooks/use-initial-load";
import { fallbackOrganization, OrganizationContext } from "./organization-context";

export function OrganizationProvider({ children }) {
  const [organization, setOrganization] = useState(fallbackOrganization);
	const refreshOrganization = useCallback(async () => {
    try {
      const { data } = await api.get("/organization");
      setOrganization(data);
      return data;
    } catch {
      return fallbackOrganization;
    }
  }, []);
	useInitialLoad(refreshOrganization);
	const value = useMemo(() => ({ organization, setOrganization, refreshOrganization }), [organization, refreshOrganization]);
	return <OrganizationContext.Provider value={value}>{children}</OrganizationContext.Provider>;
}
