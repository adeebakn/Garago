import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../api";
import { useOrganization } from "../organization/OrganizationContext";

export default function OrganizationSettings() {
  const { organization, setOrganization } = useOrganization();
  const [form, setForm] = useState(organization);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const navigate = useNavigate();

  useEffect(() => { setForm(organization); }, [organization]);
  const change = (field) => (event) => { setSaved(false); setForm({ ...form, [field]: event.target.value }); };
  const save = async (event) => {
    event.preventDefault(); setError(""); setSaving(true);
    try {
      const { data } = await api.put("/organization", form);
      setOrganization(data); setForm(data); setSaved(true);
    } catch (requestError) {
      setError(requestError.response?.data?.error || "Could not save organization details.");
    } finally { setSaving(false); }
  };

  if (saved) return <section><div className="page-toolbar"><h2>Organization settings</h2></div><div className="settings-complete"><span className="settings-complete-mark">✓</span><h3>Settings saved</h3><p><strong>{organization.name}</strong> is now shown on the sign-in screen and payment receipts.</p><div><button className="add-button" onClick={() => navigate("/")}>Go to dashboard</button><button type="button" onClick={() => setSaved(false)}>Edit details</button></div></div></section>;

  return <section><div className="page-toolbar"><h2>Organization settings</h2></div><p className="settings-intro">These details identify your shop on the sign-in screen and payment receipts.</p><div className="settings-card"><form onSubmit={save}><label>Shop or organization name<input required value={form.name || ""} onChange={change("name")} placeholder="Example Motors" /></label><label>Address<textarea value={form.address || ""} onChange={change("address")} placeholder="Street, city, state and postcode" rows="3" /></label><div className="form-grid"><label>Contact number<input value={form.phone || ""} onChange={change("phone")} placeholder="Phone number" /></label><label>Email address<input type="email" value={form.email || ""} onChange={change("email")} placeholder="Email address" /></label></div>{error && <p className="error">{error}</p>}<button className="add-button" disabled={saving}>{saving ? "Saving…" : "Save settings"}</button></form></div></section>;
}
