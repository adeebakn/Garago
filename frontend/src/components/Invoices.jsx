import { useEffect, useState } from "react";
import api from "../api";
import { useConfirm } from "../ui/ConfirmContext";
import FormDrawer from "./FormDrawer";
import SearchInput from "./SearchInput";
import SearchableSelect from "./SearchableSelect";

export default function Invoices() {
  const [orders, setOrders] = useState([]), [invoices, setInvoices] = useState([]), [payments, setPayments] = useState([]), [customers, setCustomers] = useState([]), [vehicles, setVehicles] = useState([]), [workOrderId, setWorkOrderId] = useState(""), [open, setOpen] = useState(false), [error, setError] = useState(""), [search, setSearch] = useState("");
  const confirm = useConfirm();
  const load = async () => { try { const [o, i, p, c, v] = await Promise.all([api.get("/work-orders"), api.get("/invoices"), api.get("/payments"), api.get("/customers"), api.get("/vehicles")]); setOrders(o.data || []); setInvoices(i.data || []); setPayments(p.data || []); setCustomers(c.data || []); setVehicles(v.data || []); } catch { setError("Could not load invoices."); } };
  useEffect(() => { load(); }, []);
  const close = () => { setOpen(false); setWorkOrderId(""); setError(""); };
  const create = async (event) => { event.preventDefault(); try { await api.post("/invoices", { work_order_id: Number(workOrderId) }); await load(); close(); } catch (err) { setError(err.response?.data?.error || "Could not create invoice."); } };
  const existing = new Set(invoices.map((invoice) => String(invoice.WorkOrderID)));
  const availableWorkOrders = orders.filter((order) => order.Status !== "Cancelled" && !existing.has(String(order.ID)));
  const workOrderOptions = availableWorkOrders.map((order) => {
    const customer = customers.find((entry) => String(entry.ID) === String(order.CustomerID));
    const vehicle = vehicles.find((entry) => String(entry.ID) === String(order.VehicleID));
    return { value: order.ID, label: order.WorkOrderNumber || "Work order", description: `${customer?.Name || "Customer"}${customer?.Phone ? ` · ${customer.Phone}` : ""} — ${vehicle?.RegNo || vehicle?.Vehicle || "Vehicle"}` };
  });
  const filteredInvoices = invoices.filter((invoice) => `${invoice.InvoiceNumber} ${orders.find((order) => String(order.ID) === String(invoice.WorkOrderID))?.WorkOrderNumber || ""} ${invoice.Status} ${invoice.TotalAmount}`.toLowerCase().includes(search.toLowerCase()));
  const paid = (invoiceID) => payments.filter((payment) => payment.InvoiceID === invoiceID).reduce((sum, payment) => sum + payment.Amount, 0);
  const remove = async (id) => { if (!await confirm({ title: "Archive invoice?", message: "Only invoices without payments can be archived.", confirmLabel: "Archive", destructive: true })) return; try { await api.delete(`/invoices/${id}`); load(); } catch (err) { setError(err.response?.data?.error || "Could not archive invoice."); } };
  return <section><div className="module-tools"><div className="page-toolbar"><h2>Invoices</h2><button className="add-button" onClick={() => { setError(""); setOpen(true); }}>+ Add Invoice</button></div><SearchInput value={search} onChange={setSearch} placeholder="Search invoices" /></div>{error && !open && <p className="error">{error}</p>}<FormDrawer open={open} title="Add invoice" onClose={close}><form className="form-grid" onSubmit={create}>{error && <p className="error">{error}</p>}<p className="form-helper">Each work order receives one invoice.</p><label>Work order<SearchableSelect required value={workOrderId} onChange={setWorkOrderId} options={workOrderOptions} placeholder="Search work order, customer, phone, or vehicle" /></label>{availableWorkOrders.length === 0 && <p className="form-helper">Every available work order already has an invoice.</p>}<div className="drawer-form-actions"><button disabled={!workOrderId}>Save Invoice</button><button type="button" onClick={close}>Cancel</button></div></form></FormDrawer><div className="table-wrap"><table><thead><tr><th>SL No.</th><th>Invoice</th><th>Work order</th><th>Total</th><th>Paid</th><th>Balance</th><th>Status</th><th className="actions-column">Actions</th></tr></thead><tbody>{filteredInvoices.length === 0 ? <tr><td colSpan="8" className="empty-cell">No invoices found.</td></tr> : filteredInvoices.map((invoice, index) => { const paidAmount = paid(invoice.ID); const balance = invoice.TotalAmount - paidAmount; const order = orders.find((entry) => String(entry.ID) === String(invoice.WorkOrderID)); return <tr key={invoice.ID}><td>{index + 1}</td><td><strong>{invoice.InvoiceNumber || "—"}</strong></td><td>{order?.WorkOrderNumber || "—"}</td><td>₹{invoice.TotalAmount}</td><td>₹{paidAmount}</td><td><span className={balance > 0 ? "balance-due" : "balance-settled"}>₹{balance}</span></td><td><span className={invoice.Status === "Paid" ? "status-paid" : "status-unpaid"}>{invoice.Status}</span></td><td className="actions-column"><button type="button" onClick={() => remove(invoice.ID)}>Archive</button></td></tr>; })}</tbody></table></div></section>;
}
