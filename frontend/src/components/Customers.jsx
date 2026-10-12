import { useState } from "react";
import { useInitialLoad } from "../hooks/use-initial-load";
import api from "../api";
import { useConfirm } from "../ui/confirm-context";
import Pagination from "./Pagination";
import SearchInput from "./SearchInput";
import FormDrawer from "./FormDrawer";

const emptyForm = { name: "", phone: "", email: "" };

// pgtype.Text can arrive as a plain JSON string or as an object depending on
// the backend serializer. Supporting both keeps the UI independent of that
// database implementation detail.
const customerEmail = (customer) => {
  if (typeof customer.Email === "string") return customer.Email;
  return customer.Email?.String || "";
};

function Customers() {
  const [customers, setCustomers] = useState([]);
  const [form, setForm] = useState(emptyForm);
  const [editingCustomer, setEditingCustomer] = useState(null);
  const [formOpen, setFormOpen] = useState(false);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const pageSize = 10;
  const confirm = useConfirm();

  const loadCustomers = async () => {
    try {
      setCustomers((await api.get("/customers")).data);
      setPage(1);
    } catch (err) {
      setError(err.response?.data?.error || "Could not load customers.");
    }
  };

  useInitialLoad(loadCustomers);

  const closeForm = () => {
    setFormOpen(false);
    setEditingCustomer(null);
    setForm(emptyForm);
  };

  const openNewForm = () => {
    setError("");
    setEditingCustomer(null);
    setForm(emptyForm);
    setFormOpen(true);
  };

  const openEditForm = (customer) => {
    setError("");
    setEditingCustomer(customer);
    setForm({ name: customer.Name, phone: customer.Phone, email: customerEmail(customer) });
    setFormOpen(true);
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const submit = async (event) => {
    event.preventDefault();
    setError("");
    setSubmitting(true);
    try {
      if (editingCustomer) {
        if (!await confirm({ title: "Update customer?", message: "Save the changes to this customer record?", confirmLabel: "Update" })) return;
        await api.put(`/customers/${editingCustomer.ID}`, form);
      }
      else await api.post("/customers", form);
      await loadCustomers();
      closeForm();
    } catch (err) {
      setError(err.response?.data?.error || "Could not save customer.");
    } finally {
      setSubmitting(false);
    }
  };

  const remove = async (id) => {
    if (!await confirm({ title: "Delete customer?", message: "This customer will no longer appear in the workspace.", confirmLabel: "Delete", destructive: true })) return;
    setError("");
    try {
      await api.delete(`/customers/${id}`);
      await loadCustomers();
    } catch (err) {
      setError(err.response?.data?.error || "Could not delete customer.");
    }
  };

  return <section>
    <div className="module-tools"><div className="page-toolbar"><div><h2> Manage Customers</h2></div><button className="add-button" onClick={openNewForm}>+ Add Customer</button></div><SearchInput value={search} onChange={(value) => { setSearch(value); setPage(1); }} placeholder="Search customers" /></div>
    {error && <p className="error">{error}</p>}
    <FormDrawer open={formOpen} title={editingCustomer ? "Edit customer" : "Add customer"} onClose={closeForm}>
      <form className="form-grid" onSubmit={submit}>
        <input required placeholder="Customer name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
        <input required placeholder="Phone number" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
        <input type="email" placeholder="Email (optional)" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
        <button disabled={submitting}>{submitting ? "Saving…" : editingCustomer ? "Update Customer" : "Save Customer"}</button>
        <button type="button" onClick={closeForm}>Cancel</button>
      </form>
    </FormDrawer>
    <div className="table-wrap"><table><thead><tr><th>SL No.</th><th>Name</th><th>Phone</th><th>Email</th><th className="actions-column">Actions</th></tr></thead><tbody>
      {customers.length === 0 ? <tr><td colSpan="5" className="empty-cell">No customers yet. Use “Add Customer” to create one.</td></tr> : customers.filter((customer) => `${customer.Name} ${customer.Phone} ${customerEmail(customer)}`.toLowerCase().includes(search.toLowerCase())).slice((page - 1) * pageSize, page * pageSize).map((customer, index) => <tr key={customer.ID}><td>{(page - 1) * pageSize + index + 1}</td><td>{customer.Name}</td><td>{customer.Phone}</td><td>{customerEmail(customer) || "—"}</td><td className="actions-column"><button onClick={() => openEditForm(customer)}>Edit</button> <button type="button" onClick={() => remove(customer.ID)}>Delete</button></td></tr>)}
    </tbody></table></div><Pagination page={page} totalItems={customers.filter((customer) => `${customer.Name} ${customer.Phone} ${customerEmail(customer)}`.toLowerCase().includes(search.toLowerCase())).length} pageSize={pageSize} onPageChange={setPage} />
  </section>;
}

export default Customers;
