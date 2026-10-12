import { useMemo, useState } from "react";
import { useInitialLoad } from "../hooks/use-initial-load";
import api from "../api";
import { useConfirm } from "../ui/confirm-context";
import PaymentReceipt from "./PaymentReceipt";
import FormDrawer from "./FormDrawer";
import SearchInput from "./SearchInput";
import SearchableSelect from "./SearchableSelect";

function Payments() {
  const [invoices, setInvoices] = useState([]), [payments, setPayments] = useState([]), [orders, setOrders] = useState([]), [customers, setCustomers] = useState([]), [vehicles, setVehicles] = useState([]), [invoiceId, setInvoiceId] = useState(""), [amount, setAmount] = useState(""), [method, setMethod] = useState("Cash"), [open, setOpen] = useState(false), [error, setError] = useState(""), [receipt, setReceipt] = useState(null), [search, setSearch] = useState("");
  const confirm = useConfirm();

  const load = async () => {
    try {
      const [invoiceResponse, paymentResponse, orderResponse, customerResponse, vehicleResponse] = await Promise.all([
        api.get("/invoices"), api.get("/payments"), api.get("/work-orders"), api.get("/customers"), api.get("/vehicles"),
      ]);
      setInvoices(Array.isArray(invoiceResponse.data) ? invoiceResponse.data : []);
      setPayments(Array.isArray(paymentResponse.data) ? paymentResponse.data : []);
      setOrders(Array.isArray(orderResponse.data) ? orderResponse.data : []);
      setCustomers(Array.isArray(customerResponse.data) ? customerResponse.data : []);
      setVehicles(Array.isArray(vehicleResponse.data) ? vehicleResponse.data : []);
    } catch {
      setError("Could not load payments.");
    }
  };

  useInitialLoad(load);
  const close = () => { setOpen(false); setInvoiceId(""); setAmount(""); setMethod("Cash"); setError(""); };
  const selected = useMemo(() => invoices.find((invoice) => String(invoice.ID) === String(invoiceId)), [invoices, invoiceId]);
  const paidFor = (id) => payments.filter((payment) => payment.InvoiceID === id).reduce((sum, payment) => sum + payment.Amount, 0);
  const paid = selected ? paidFor(selected.ID) : 0;

  const invoiceOptions = useMemo(() => invoices.filter((invoice) => invoice.Status !== "Paid").map((invoice) => {
    const order = orders.find((entry) => String(entry.ID) === String(invoice.WorkOrderID));
    const customer = customers.find((entry) => String(entry.ID) === String(invoice.CustomerID));
    const vehicle = vehicles.find((entry) => String(entry.ID) === String(order?.VehicleID));
    const invoiceCode = invoice.InvoiceNumber || `Invoice ${invoice.ID}`;
    const customerName = customer?.Name || "Unknown customer";
    const vehicleRegistration = vehicle?.RegNo || vehicle?.Vehicle || "No vehicle";
    const total = `₹${Number(invoice.TotalAmount || 0).toLocaleString("en-IN", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
    return { value: invoice.ID, label: `${invoiceCode} – ${customerName} / ${vehicleRegistration} – ${total}` };
  }), [invoices, orders, customers, vehicles]);

  const filteredPayments = payments.filter((payment) => `${invoices.find((invoice) => invoice.ID === payment.InvoiceID)?.InvoiceNumber || ""} ${payment.Amount} ${payment.PaymentMethod}`.toLowerCase().includes(search.toLowerCase()));
  const submit = async (event) => { event.preventDefault(); try { await api.post("/payments", { invoice_id: Number(invoiceId), amount: Number(amount), payment_method: method }); await load(); close(); } catch (requestError) { setError(requestError.response?.data?.error || "Could not record payment."); } };
  const remove = async (id) => { if (!await confirm({ title: "Archive payment?", message: "The payment stays in audit history and the invoice balance will be recalculated.", confirmLabel: "Archive", destructive: true })) return; try { await api.delete(`/payments/${id}`); load(); } catch (requestError) { setError(requestError.response?.data?.error || "Could not archive payment."); } };
  const printReceipt = (payment) => { const invoice = invoices.find((entry) => entry.ID === payment.InvoiceID); const paidToDate = paidFor(payment.InvoiceID); setReceipt({ payment, invoice, paid: paidToDate, balance: (invoice?.TotalAmount || 0) - paidToDate }); };

  return <section><div className="module-tools"><div className="page-toolbar"><h2>Payments</h2><button className="add-button" onClick={() => { setError(""); setOpen(true); }}>+ Add Payment</button></div><SearchInput value={search} onChange={setSearch} placeholder="Search payments" /></div>{error && !open && <p className="error">{error}</p>}<FormDrawer open={open} title="Record payment" onClose={close}><form className="form-grid" onSubmit={submit}>{error && <p className="error">{error}</p>}<label>Invoice<SearchableSelect required value={invoiceId} onChange={setInvoiceId} options={invoiceOptions} placeholder="Search invoice, customer, vehicle, or amount" resetKey={open} /></label><input required min="0.01" step="0.01" type="number" placeholder="Amount" value={amount} onChange={(event) => setAmount(event.target.value)} /><select value={method} onChange={(event) => setMethod(event.target.value)}>{["Cash", "UPI", "Card", "Bank Transfer"].map((entry) => <option key={entry}>{entry}</option>)}</select>{selected && <p className="payment-balance">Invoice total: ₹{selected.TotalAmount} <span>Paid: ₹{paid}</span> <strong>Balance: ₹{selected.TotalAmount - paid}</strong></p>}<button disabled={!invoiceId}>Save Payment</button><button type="button" onClick={close}>Cancel</button></form></FormDrawer><div className="table-wrap"><table><thead><tr><th>SL No.</th><th>Invoice</th><th>Amount</th><th>Method</th><th className="actions-column">Actions</th></tr></thead><tbody>{filteredPayments.length === 0 ? <tr><td colSpan="5" className="empty-cell">No payments found.</td></tr> : filteredPayments.map((payment, index) => <tr key={payment.ID}><td>{index + 1}</td><td>{invoices.find((invoice) => invoice.ID === payment.InvoiceID)?.InvoiceNumber || "—"}</td><td>₹{payment.Amount}</td><td>{payment.PaymentMethod}</td><td className="actions-column"><button type="button" onClick={() => printReceipt(payment)}>Preview</button> <button type="button" onClick={() => remove(payment.ID)}>Archive</button></td></tr>)}</tbody></table></div>{receipt && <PaymentReceipt {...receipt} onClose={() => setReceipt(null)} />}</section>;
}

export default Payments;
