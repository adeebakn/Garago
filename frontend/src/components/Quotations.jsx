import { useMemo, useState } from "react";
import { useInitialLoad } from "../hooks/use-initial-load";
import api from "../api";
import { useConfirm } from "../ui/confirm-context";
import SearchableSelect from "./SearchableSelect";
import FormDrawer from "./FormDrawer";
import QuotationPreview from "./QuotationPreview";
import SearchInput from "./SearchInput";

const emptyItem = () => ({ service_id: "", quantity: 1 });

function Quotations() {
  const [customers, setCustomers] = useState([]);
  const [vehicles, setVehicles] = useState([]);
  const [services, setServices] = useState([]);
  const [quotations, setQuotations] = useState([]);
  const [customerId, setCustomerId] = useState("");
  const [vehicleId, setVehicleId] = useState("");
  const [items, setItems] = useState([emptyItem()]);
  const [advanceRequested, setAdvanceRequested] = useState("");
  const [open, setOpen] = useState(false);
  const [error, setError] = useState("");
  const [quickType, setQuickType] = useState("");
  const [quickForm, setQuickForm] = useState({});
  const [preview, setPreview] = useState(null);
  const [search, setSearch] = useState("");
  const [editingQuotation, setEditingQuotation] = useState(null);
  const confirm = useConfirm();

  const load = async () => {
    try {
      const [customerResponse, vehicleResponse, serviceResponse, quotationResponse] = await Promise.all([
        api.get("/customers"),
        api.get("/vehicles"),
        api.get("/services"),
        api.get("/quotations"),
      ]);
      setCustomers(customerResponse.data);
      setVehicles(vehicleResponse.data);
      setServices(serviceResponse.data);
      setQuotations(quotationResponse.data);
    } catch {
      setError("Could not load quotation data.");
    }
  };

  useInitialLoad(load);

  const close = () => {
    setOpen(false);
    setError("");
    setCustomerId("");
    setVehicleId("");
    setItems([emptyItem()]);
    setAdvanceRequested("");
    setEditingQuotation(null);
  };

  const filteredVehicles = vehicles.filter((vehicle) => String(vehicle.CustomerID) === String(customerId));
  const total = useMemo(() => items.reduce((sum, item) => {
    const service = services.find((entry) => String(entry.ID) === String(item.service_id));
    return sum + (service ? service.Price * Number(item.quantity || 0) : 0);
  }, 0), [items, services]);

  const changeItem = (index, key, value) => setItems(items.map((item, itemIndex) => (
    itemIndex === index ? { ...item, [key]: value } : item
  )));

  const submit = async (event) => {
    event.preventDefault();
    try {
      const payload = {
        customer_id: Number(customerId),
        vehicle_id: Number(vehicleId),
        advance_requested: Number(advanceRequested || 0),
        items: items.map((item) => ({ service_id: Number(item.service_id), quantity: Number(item.quantity) })),
      };
      if (editingQuotation) await api.put(`/quotations/${editingQuotation.ID}`, payload);
      else await api.post("/quotations", payload);
      await load();
      close();
    } catch (requestError) {
      setError(requestError.response?.data?.error || "Could not create quotation.");
    }
  };

  const openQuick = (type) => {
    setError("");
    setQuickType(type);
    setQuickForm(type === "customer"
      ? { name: "", phone: "", email: "" }
      : type === "vehicle"
        ? { reg_no: "", vehicle: "" }
        : { name: "", description: "", price: "" });
  };

  const saveQuick = async (event) => {
    event.preventDefault();
    try {
      if (quickType === "customer") {
        const { data } = await api.post("/customers", quickForm);
        setCustomerId(String(data.ID));
        setVehicleId("");
      }
      if (quickType === "vehicle") {
        const { data } = await api.post("/vehicles", { ...quickForm, customer_id: Number(customerId) });
        setVehicleId(String(data.ID));
      }
      if (quickType === "service") {
        const { data } = await api.post("/services", { ...quickForm, price: Number(quickForm.price) });
        setItems((currentItems) => {
          const emptyIndex = currentItems.findIndex((item) => !item.service_id);
          if (emptyIndex === -1) return [...currentItems, { service_id: String(data.ID), quantity: 1 }];
          return currentItems.map((item, index) => (index === emptyIndex ? { ...item, service_id: String(data.ID) } : item));
        });
      }
      await load();
      setQuickType("");
    } catch (requestError) {
      setError(requestError.response?.data?.error || "Could not save record.");
    }
  };

  const updateStatus = async (id, value) => {
    const isApproval = value === "Approved";
    const approved = await confirm({
      title: isApproval ? "Approve quotation?" : "Update quotation status?",
      message: isApproval ? "Approving this quotation allows a work order to be created." : `Change quotation status to ${value}?`,
      confirmLabel: isApproval ? "Approve" : "Update",
    });
    if (!approved) return;
    try {
      await api.put(`/quotations/${id}/status`, { status: value });
      load();
    } catch (requestError) {
      setError(requestError.response?.data?.error || "Could not update quotation.");
    }
  };

  const remove = async (id) => {
    const approved = await confirm({ title: "Delete quotation?", message: "This action cannot be undone.", confirmLabel: "Delete", destructive: true });
    if (!approved) return;
    try {
      await api.delete(`/quotations/${id}`);
      load();
    } catch (requestError) {
      setError(requestError.response?.data?.error || "Could not delete quotation.");
    }
  };

  const openPreview = async (quotation) => {
    try {
      const { data } = await api.get(`/quotations/${quotation.ID}`);
      setPreview({ quotation: data.quotation, items: data.items || [] });
    } catch (requestError) {
      setError(requestError.response?.data?.error || "Could not load quotation preview.");
    }
  };

  const openEdit = async (quotation) => {
    try {
      const { data } = await api.get(`/quotations/${quotation.ID}`);
      setEditingQuotation(quotation);
      setCustomerId(String(data.quotation.CustomerID));
      setVehicleId(String(data.quotation.VehicleID));
      setAdvanceRequested(String(data.quotation.AdvanceRequested || ""));
      setItems((data.items || []).map((item) => ({ service_id: String(item.ServiceID), quantity: item.Quantity })) || [emptyItem()]);
      setError("");
      setOpen(true);
    } catch (requestError) {
      setError(requestError.response?.data?.error || "Could not load quotation for editing.");
    }
  };

  const filteredQuotations = quotations.filter((quotation) => {
    const customer = customers.find((entry) => String(entry.ID) === String(quotation.CustomerID));
    const vehicle = vehicles.find((entry) => String(entry.ID) === String(quotation.VehicleID));
    return `${customer?.Name || ""} ${customer?.Phone || ""} ${vehicle?.RegNo || ""} ${vehicle?.Vehicle || ""} ${quotation.TotalAmount} ${quotation.Status}`.toLowerCase().includes(search.toLowerCase());
  });

  return (
    <section>
      <div className="module-tools"><div className="page-toolbar">
        <h2>Quotations</h2>
        <button className="add-button" onClick={() => { setError(""); setOpen(true); }}>+ Add Quotation</button>
      </div><SearchInput value={search} onChange={setSearch} placeholder="Search quotations" /></div>
      {error && !open && !quickType && <p className="error">{error}</p>}

      <FormDrawer open={open} title={editingQuotation ? "Edit draft quotation" : "Add quotation"} onClose={close}>
        <form onSubmit={submit}>
          {error && !quickType && <p className="error">{error}</p>}
          <div className="form-grid">
            <div className="picker-with-add">
              <SearchableSelect required value={customerId} onChange={(value) => { setCustomerId(value); setVehicleId(""); setItems([emptyItem()]); }} options={customers.map((customer) => ({ value: customer.ID, label: customer.Name, description: customer.Phone }))} placeholder="Search customer" />
              <button type="button" className="picker-add" disabled={Boolean(customerId)} onClick={() => openQuick("customer")} aria-label="Add customer">+</button>
            </div>
            <div className="picker-with-add">
              <SearchableSelect required disabled={!customerId} resetKey={customerId} value={vehicleId} onChange={(value) => { setVehicleId(value); setItems([emptyItem()]); }} options={filteredVehicles.map((vehicle) => ({ value: vehicle.ID, label: vehicle.RegNo || vehicle.Vehicle, description: vehicle.Vehicle }))} placeholder={!customerId ? "Choose a customer first" : filteredVehicles.length === 0 ? "No vehicles added for this customer" : "Search vehicle"} />
              <button type="button" className="picker-add" disabled={!customerId || Boolean(vehicleId)} onClick={() => openQuick("vehicle")} aria-label="Add vehicle">+</button>
            </div>
          </div>
          {customerId && filteredVehicles.length === 0 && <p className="form-helper">No vehicles have been added for this customer. Add a vehicle before creating the quotation.</p>}
          <h3>Services</h3>
          {items.map((item, index) => (
            <div className="line-item" key={index}>
              <div className="picker-with-add">
                <SearchableSelect required disabled={!customerId || !vehicleId} resetKey={`${customerId}-${vehicleId}`} value={item.service_id} onChange={(value) => changeItem(index, "service_id", value)} options={services.map((service) => ({ value: service.ID, label: service.Name, description: `₹${service.Price}` }))} placeholder={!customerId ? "Choose a customer first" : !vehicleId ? "Choose a vehicle first" : "Search service"} />
                <button type="button" className="picker-add" disabled={!customerId || !vehicleId || Boolean(item.service_id)} onClick={() => openQuick("service")} aria-label="Add service">+</button>
              </div>
              <input required min="1" type="number" value={item.quantity} onChange={(event) => changeItem(index, "quantity", event.target.value)} />
              {items.length > 1 && <button type="button" className="remove-button" onClick={() => setItems(items.filter((_, itemIndex) => itemIndex !== index))}>Remove</button>}
            </div>
          ))}
          <button type="button" className="quotation-add-service" disabled={!customerId || !vehicleId} onClick={() => setItems([...items, emptyItem()])}>+ Add another service</button>
    
          <div className="quotation-summary">
            <span>Quotation total</span>
            <strong>₹{Number(total).toLocaleString("en-IN", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</strong>
          </div>
          <label className="quotation-advance">Advance requested (optional)<input type="number" min="0" max={total} step="any" placeholder="₹0.00" value={advanceRequested} onChange={(event) => setAdvanceRequested(event.target.value)} /><small>Balance after advance: ₹{Math.max(total - Number(advanceRequested || 0), 0).toLocaleString("en-IN", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</small></label>
          
          <div className="quotation-actions">
            <button type="submit" className="quotation-save">{editingQuotation ? "Update quotation" : "Save quotation"}</button>
            <button type="button" className="quotation-cancel" onClick={close}>Cancel</button>
          </div>
        </form>
      </FormDrawer>

      <FormDrawer open={Boolean(quickType)} title={`Add ${quickType || "record"}`} onClose={() => { setQuickType(""); setError(""); }}>
        <form onSubmit={saveQuick}>
          {error && <p className="error">{error}</p>}
          {quickType === "customer" && <>
            <input required placeholder="Customer name" value={quickForm.name || ""} onChange={(event) => setQuickForm({ ...quickForm, name: event.target.value })} />
            <input required placeholder="Phone number" value={quickForm.phone || ""} onChange={(event) => setQuickForm({ ...quickForm, phone: event.target.value })} />
            <input placeholder="Email (optional)" value={quickForm.email || ""} onChange={(event) => setQuickForm({ ...quickForm, email: event.target.value })} />
          </>}
          {quickType === "vehicle" && <>
            <input required placeholder="Registration number" value={quickForm.reg_no || ""} onChange={(event) => setQuickForm({ ...quickForm, reg_no: event.target.value })} />
            <input required placeholder="Vehicle" value={quickForm.vehicle || ""} onChange={(event) => setQuickForm({ ...quickForm, vehicle: event.target.value })} />
          </>}
          {quickType === "service" && <>
            <input required placeholder="Service name" value={quickForm.name || ""} onChange={(event) => setQuickForm({ ...quickForm, name: event.target.value })} />
            <input placeholder="Description (optional)" value={quickForm.description || ""} onChange={(event) => setQuickForm({ ...quickForm, description: event.target.value })} />
            <input required type="number" min="0.01" step="any" placeholder="Price" value={quickForm.price || ""} onChange={(event) => setQuickForm({ ...quickForm, price: event.target.value })} />
          </>}
          <button>Save</button>
          <button type="button" onClick={() => { setQuickType(""); setError(""); }}>Cancel</button>
        </form>
      </FormDrawer>

      <div className="table-wrap">
        <table>
          <thead><tr><th>SL No.</th><th>Customer</th><th>Vehicle</th><th>Total</th><th>Status</th><th className="actions-column">Actions</th></tr></thead>
          <tbody>
            {filteredQuotations.length === 0 ? <tr><td colSpan="6" className="empty-cell">No quotations found.</td></tr> : filteredQuotations.map((quotation, index) => (
              <tr key={quotation.ID}>
                <td>{index + 1}</td><td>{(() => { const customer = customers.find((entry) => String(entry.ID) === String(quotation.CustomerID)); return customer ? <><strong>{customer.Name}</strong><br /><small>{customer.Phone}</small></> : "—"; })()}</td><td>{(() => { const vehicle = vehicles.find((entry) => String(entry.ID) === String(quotation.VehicleID)); return vehicle ? <><strong>{vehicle.RegNo || vehicle.Vehicle}</strong>{vehicle.RegNo && vehicle.Vehicle && <><br /><small>{vehicle.Vehicle}</small></>}</> : "—"; })()}</td><td>₹{quotation.TotalAmount}</td>
                <td><select value={quotation.Status} onChange={(event) => updateStatus(quotation.ID, event.target.value)}><option>Draft</option><option>Approved</option><option>Rejected</option></select></td>
                <td className="actions-column"><div className="quotation-row-actions">{quotation.Status === "Draft" && <button type="button" onClick={() => openEdit(quotation)}>Edit</button>}<button type="button" onClick={() => openPreview(quotation)}>Preview</button>{quotation.Status === "Draft" && <button type="button" className="action-danger" onClick={() => remove(quotation.ID)}>Delete</button>}</div></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {preview && <QuotationPreview quotation={preview.quotation} items={preview.items} customer={customers.find((customer) => String(customer.ID) === String(preview.quotation.CustomerID))} vehicle={vehicles.find((vehicle) => String(vehicle.ID) === String(preview.quotation.VehicleID))} services={services} onClose={() => setPreview(null)} />}
    </section>
  );
}

export default Quotations;
