import { useState } from "react";
import { useInitialLoad } from "../hooks/use-initial-load";
import api from "../api";
import { useConfirm } from "../ui/confirm-context";
import SearchInput from "./SearchInput";
import FormDrawer from "./FormDrawer";

const blank = { name: "", description: "", price: "" };
const descriptionOf = (service) => typeof service.Description === "string" ? service.Description : service.Description?.String || "";

function Services() {
  const [services, setServices] = useState([]), [form, setForm] = useState(blank), [editing, setEditing] = useState(null), [open, setOpen] = useState(false), [error, setError] = useState(""), [search, setSearch] = useState("");
  const confirm = useConfirm();
  const load = async () => { try { setServices((await api.get("/services")).data); } catch (err) { setError(err.response?.data?.error || "Could not load services."); } };
  useInitialLoad(load);
  const close = () => { setOpen(false); setEditing(null); setForm(blank); };
  const add = () => { setError(""); setEditing(null); setForm(blank); setOpen(true); };
  const edit = (service) => { setError(""); setEditing(service); setForm({ name: service.Name, description: descriptionOf(service), price: service.Price }); setOpen(true); window.scrollTo({ top: 0, behavior: "smooth" }); };
  const submit = async (event) => { event.preventDefault(); try { if (editing && !await confirm({ title: "Update service?", message: "Save the changes to this service?", confirmLabel: "Update" })) return; if (editing) await api.put(`/services/${editing.ID}`, { ...form, price: Number(form.price) }); else await api.post("/services", { ...form, price: Number(form.price) }); await load(); close(); } catch (err) { setError(err.response?.data?.error || "Could not save service."); } };
  const remove = async (id) => { if (!await confirm({ title: "Delete service?", message: "This action cannot be undone.", confirmLabel: "Delete", destructive: true })) return; try { await api.delete(`/services/${id}`); load(); } catch (err) { setError(err.response?.data?.error || "Could not delete service."); } };
  const filteredServices = services.filter((service) => `${service.Name} ${descriptionOf(service)} ${service.Price}`.toLowerCase().includes(search.toLowerCase()));
  return <section><div className="module-tools"><div className="page-toolbar"><h2>Services</h2><button className="add-button" onClick={add}>+ Add Service</button></div><SearchInput value={search} onChange={setSearch} placeholder="Search services" /></div>{error && !open && <p className="error">{error}</p>}<FormDrawer open={open} title={editing ? "Edit service" : "Add service"} onClose={close}><form className="form-grid" onSubmit={submit}>{error && <p className="error">{error}</p>}<input required placeholder="Service name" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })}/><input placeholder="Description" value={form.description} onChange={e => setForm({ ...form, description: e.target.value })}/><input required min="0.01" step="0.01" type="number" placeholder="Price" value={form.price} onChange={e => setForm({ ...form, price: e.target.value })}/><button>{editing ? "Update Service" : "Save Service"}</button><button type="button" onClick={close}>Cancel</button></form></FormDrawer><div className="table-wrap"><table><thead><tr><th>SL No.</th><th>Name</th><th>Description</th><th>Price</th><th className="actions-column">Actions</th></tr></thead><tbody>{filteredServices.length === 0 ? <tr><td colSpan="5" className="empty-cell">No services found.</td></tr> : filteredServices.map((s, index) => <tr key={s.ID}><td>{index + 1}</td><td>{s.Name}</td><td>{descriptionOf(s) || "—"}</td><td>₹{s.Price}</td><td className="actions-column"><button onClick={() => edit(s)}>Edit</button> <button type="button" onClick={() => remove(s.ID)}>Delete</button></td></tr>)}</tbody></table></div></section>;
}
export default Services;
