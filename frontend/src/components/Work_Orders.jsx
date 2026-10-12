import { useEffect, useState } from "react";
import api from "../api";
import { useConfirm } from "../ui/ConfirmContext";
import FormDrawer from "./FormDrawer";
import SearchableSelect from "./SearchableSelect";
import SearchInput from "./SearchInput";
import QuotationPreview from "./QuotationPreview";
function WorkOrders() {
  const [quotes, setQuotes] = useState([]), [orders, setOrders] = useState([]), [customers, setCustomers] = useState([]), [vehicles, setVehicles] = useState([]), [services, setServices] = useState([]), [dueDate, setDueDate] = useState(""), [quoteId, setQuoteId] = useState(""), [open, setOpen] = useState(false), [error, setError] = useState(""), [search, setSearch] = useState(""), [dueDateFilter, setDueDateFilter] = useState(""), [preview, setPreview] = useState(null), [editingOrder, setEditingOrder] = useState(null);
  const confirm = useConfirm();
  const load = async () => { try { const [q, o, c, v, s] = await Promise.all([api.get("/quotations"), api.get("/work-orders"), api.get("/customers"), api.get("/vehicles"), api.get("/services")]); const quotations = Array.isArray(q.data) ? q.data : []; setQuotes(quotations.filter(x => x.Status === "Approved")); setOrders(Array.isArray(o.data) ? o.data : []); setCustomers(c.data || []); setVehicles(v.data || []); setServices(s.data || []); } catch { setError("Could not load work orders."); } };
  useEffect(() => { load(); }, []);
  const close = () => { setOpen(false); setQuoteId(""); setDueDate(""); setEditingOrder(null); };
  const quotationOption = (quotation) => {
    const customer = customers.find((entry) => String(entry.ID) === String(quotation.CustomerID));
    const vehicle = vehicles.find((entry) => String(entry.ID) === String(quotation.VehicleID));
    const customerName = customer?.Name || "Customer";
    const phone = customer?.Phone || "No phone recorded";
    const vehicleName = vehicle?.RegNo || vehicle?.Vehicle || "Vehicle not recorded";
    return {
      value: quotation.ID,
      label: `${customerName} — ${vehicleName}`,
      description: `${phone} · Quoted ₹${quotation.TotalAmount}`,
    };
  };
  const formatDueDate = (dueDate) => {
    const value = typeof dueDate === "string" ? dueDate : dueDate?.Time || dueDate?.time;
    if (!value) return "—";
    const parsed = new Date(value);
    if (Number.isNaN(parsed.getTime())) return String(value).slice(0, 10);
    return new Intl.DateTimeFormat("en-IN", { day: "2-digit", month: "short", year: "numeric" }).format(parsed);
  };
  const dueDateKey = (dueDate) => {
    const value = typeof dueDate === "string" ? dueDate : dueDate?.Time || dueDate?.time;
    return value ? String(value).slice(0, 10) : "";
  };
  const submit = async event => { event.preventDefault(); try { if (editingOrder) await api.put(`/work-orders/${editingOrder.ID}`, { due_date: dueDate }); else { const quote = quotes.find(x => String(x.ID) === String(quoteId)); if (!quote) return; await api.post("/work-orders", { quotation_id: quote.ID, customer_id: quote.CustomerID, vehicle_id: quote.VehicleID, due_date: dueDate }); } await load(); close(); } catch (err) { setError(err.response?.data?.error || "Could not save work order."); } };
  const edit = (order) => { setError(""); setEditingOrder(order); setDueDate(dueDateKey(order.DueDate)); setOpen(true); };
  const setStatus = async (id, status) => { if (!await confirm({ title: "Update work order?", message: `Change work order status to ${status}?`, confirmLabel: "Update" })) return; try { await api.put(`/work-orders/${id}/status`, { status }); load(); } catch (err) { setError(err.response?.data?.error || "Could not update work order."); } };
  const remove = async id => { if (!await confirm({ title: "Delete work order?", message: "This action cannot be undone.", confirmLabel: "Delete", destructive: true })) return; try { await api.delete(`/work-orders/${id}`); load(); } catch (err) { setError(err.response?.data?.error || "Could not delete work order."); } };
  const openPreview = async (order) => { try { const { data } = await api.get(`/quotations/${order.QuotationID}`); setPreview({ order, quotation: data.quotation, items: data.items || [] }); } catch (err) { setError(err.response?.data?.error || "Could not load work order services."); } };
  const usedQuotationIDs = new Set(orders.map((order) => String(order.QuotationID)));
  const availableQuotations = quotes.filter((quotation) => !usedQuotationIDs.has(String(quotation.ID)));
  const filteredOrders = orders.filter((order) => {
    const customer = customers.find((entry) => String(entry.ID) === String(order.CustomerID));
    const vehicle = vehicles.find((entry) => String(entry.ID) === String(order.VehicleID));
    const dueKey = dueDateKey(order.DueDate);
    const matchesDue = !dueDateFilter || dueKey === dueDateFilter;
    return matchesDue && `${order.WorkOrderNumber} ${customer?.Name || ""} ${customer?.Phone || ""} ${vehicle?.RegNo || ""} ${vehicle?.Vehicle || ""} ${order.Status} ${formatDueDate(order.DueDate)} ${dueKey}`.toLowerCase().includes(search.toLowerCase());
  });
  return <section><div className="module-tools"><div className="page-toolbar"><h2>Work Orders</h2><button className="add-button" onClick={() => { setError(""); setOpen(true); }}>+ Add Work Order</button></div><div className="work-order-filters"><label>Due date<input type="date" value={dueDateFilter} onChange={(event) => setDueDateFilter(event.target.value)} /></label>{dueDateFilter && <button type="button" className="filter-clear" onClick={() => setDueDateFilter("")}>Clear</button>}<SearchInput value={search} onChange={setSearch} placeholder="Search work orders" /></div></div>{error && !open && <p className="error">{error}</p>}<FormDrawer open={open} title={editingOrder ? "Edit pending work order" : "Add work order"} onClose={close}><form className="form-grid" onSubmit={submit}>{error && <p className="error">{error}</p>}{editingOrder ? <p className="form-helper">{editingOrder.WorkOrderNumber}: only the due date can be changed after a work order is created.</p> : <><label>Approved quotation<SearchableSelect required value={quoteId} onChange={setQuoteId} options={availableQuotations.map(quotationOption)} placeholder="Search customer, phone, vehicle, or amount" /></label>{availableQuotations.length === 0 && <p className="form-helper">Every approved quotation already has a work order.</p>}</>}<label>Due date<input required type="date" value={dueDate} onChange={e => setDueDate(e.target.value)}/></label>{editingOrder ? <div className="work-order-edit-actions"><button className="work-order-update-button" disabled={!dueDate}>Update Work Order</button><button type="button" onClick={close}>Cancel</button></div> : <><button disabled={!quoteId}>Save Work Order</button><button type="button" onClick={close}>Cancel</button></>}</form></FormDrawer><div className="table-wrap"><table><thead><tr><th>SL No.</th><th>Work order</th><th>Customer</th><th>Vehicle</th><th>Due date</th><th>Status</th><th className="actions-column">Actions</th></tr></thead><tbody>{filteredOrders.length === 0 ? <tr><td colSpan="7" className="empty-cell">No work orders found.</td></tr> : filteredOrders.map((o, index) => { const customer = customers.find(c => String(c.ID) === String(o.CustomerID)); const vehicle = vehicles.find(v => String(v.ID) === String(o.VehicleID)); return <tr key={o.ID}><td>{index + 1}</td><td><button type="button" className="document-link" onClick={() => openPreview(o)}>{o.WorkOrderNumber || "—"}</button></td><td>{customer ? <><strong>{customer.Name}</strong><br /><small>{customer.Phone}</small></> : "—"}</td><td>{vehicle ? <><strong>{vehicle.RegNo || vehicle.Vehicle}</strong>{vehicle.RegNo && vehicle.Vehicle && <><br /><small>{vehicle.Vehicle}</small></>}</> : "—"}</td><td>{formatDueDate(o.DueDate)}</td><td><select value={o.Status} onChange={e => setStatus(o.ID, e.target.value)}>{["Pending", "In Progress", "Completed", "Cancelled"].map(s => <option key={s}>{s}</option>)}</select></td><td className="actions-column"><div className="work-order-row-actions">{o.Status === "Pending" && <button type="button" onClick={() => edit(o)}>Edit</button>}<button type="button" className="action-danger" onClick={() => remove(o.ID)}>Archive</button></div></td></tr>})}</tbody></table></div>{preview && <QuotationPreview quotation={preview.quotation} items={preview.items} customer={customers.find((customer) => String(customer.ID) === String(preview.order.CustomerID))} vehicle={vehicles.find((vehicle) => String(vehicle.ID) === String(preview.order.VehicleID))} services={services} documentTitle="Work order" documentCode={preview.order.WorkOrderNumber} onClose={() => setPreview(null)} />}</section>;
}
export default WorkOrders;
