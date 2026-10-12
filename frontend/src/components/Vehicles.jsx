import { useEffect, useState } from "react";
import api from "../api";
import { useConfirm } from "../ui/ConfirmContext";
import Pagination, { pageRows } from "./Pagination";
import SearchInput from "./SearchInput";
import FormDrawer from "./FormDrawer";
import SearchableSelect from "./SearchableSelect";

const blank = { customerId: "", regNo: "", vehicle: "" };
function Vehicles() {
  const [vehicles, setVehicles] = useState([]), [customers, setCustomers] = useState([]), [form, setForm] = useState(blank), [editing, setEditing] = useState(null), [open, setOpen] = useState(false), [error, setError] = useState(""), [page, setPage] = useState(1), [search, setSearch] = useState("");
  const confirm = useConfirm();
  const load = async () => { try { const [vehicleResult, customerResult] = await Promise.all([api.get("/vehicles"), api.get("/customers")]); setVehicles(vehicleResult.data); setCustomers(customerResult.data); setPage(1); } catch (err) { setError(err.response?.data?.error || "Could not load vehicles."); } };
  useEffect(() => { load(); }, []);
  const close = () => { setOpen(false); setEditing(null); setForm(blank); };
  const add = () => { setError(""); setEditing(null); setForm(blank); setOpen(true); };
  const edit = (vehicle) => { setError(""); setEditing(vehicle); setForm({ customerId: vehicle.CustomerID, regNo: vehicle.RegNo, vehicle: vehicle.Vehicle }); setOpen(true); window.scrollTo({ top: 0, behavior: "smooth" }); };
  const submit = async (event) => { event.preventDefault(); const payload = { customer_id: Number(form.customerId), reg_no: form.regNo, vehicle: form.vehicle }; try { if (editing && !await confirm({ title: "Update vehicle?", message: "Save the changes to this vehicle?", confirmLabel: "Update" })) return; if (editing) await api.put(`/vehicles/${editing.ID}`, payload); else await api.post("/vehicles", payload); await load(); close(); } catch (err) { setError(err.response?.data?.error || "Could not save vehicle."); } };
  const remove = async (id) => { if (!await confirm({ title: "Delete vehicle?", message: "This action cannot be undone.", confirmLabel: "Delete", destructive: true })) return; try { await api.delete(`/vehicles/${id}`); load(); } catch (err) { setError(err.response?.data?.error || "Could not delete vehicle."); } };
  const filteredVehicles = vehicles.filter((vehicle) => `${vehicle.RegNo} ${vehicle.Vehicle} ${vehicle.CustomerName}`.toLowerCase().includes(search.toLowerCase()));
  return <section><div className="module-tools"><div className="page-toolbar"><h2>Vehicles</h2><button className="add-button" onClick={add}>+ Add Vehicle</button></div><SearchInput value={search} onChange={(value) => { setSearch(value); setPage(1); }} placeholder="Search vehicles" /></div>{error && <p className="error">{error}</p>}<FormDrawer open={open} title={editing ? "Edit vehicle" : "Add vehicle"} onClose={close}><form className="form-grid" onSubmit={submit}><SearchableSelect required value={form.customerId} onChange={(customerId) => setForm({ ...form, customerId })} options={customers.map((customer) => ({ value: customer.ID, label: customer.Name }))} placeholder="Search customer" /><input placeholder="Registration number (optional)" value={form.regNo} onChange={e => setForm({ ...form, regNo: e.target.value })}/><input required placeholder="Vehicle (example: Mahindra XUV700)" value={form.vehicle} onChange={e => setForm({ ...form, vehicle: e.target.value })}/><button>{editing ? "Update Vehicle" : "Save Vehicle"}</button><button type="button" onClick={close}>Cancel</button></form></FormDrawer><div className="table-wrap"><table><thead><tr><th>SL No.</th><th>Registration</th><th>Vehicle</th><th>Customer</th><th className="actions-column">Actions</th></tr></thead><tbody>{filteredVehicles.length === 0 ? <tr><td colSpan="5" className="empty-cell">No vehicles found.</td></tr> : pageRows(filteredVehicles, page).map((vehicle, index) => <tr key={vehicle.ID}><td>{(page - 1) * 10 + index + 1}</td><td>{vehicle.RegNo || "—"}</td><td>{vehicle.Vehicle}</td><td>{vehicle.CustomerName}</td><td className="actions-column"><button onClick={() => edit(vehicle)}>Edit</button> <button type="button" onClick={() => remove(vehicle.ID)}>Delete</button></td></tr>)}</tbody></table></div><Pagination page={page} totalItems={filteredVehicles.length} onPageChange={setPage} /></section>;
}
export default Vehicles;
