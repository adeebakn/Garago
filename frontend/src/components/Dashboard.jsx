import { useEffect, useMemo, useState } from "react";
import api from "../api";
import Card from "./Card";

const rupees = (amount) => `₹${Number(amount || 0).toLocaleString("en-IN", { maximumFractionDigits: 2 })}`;

function Dashboard() {
  const [data, setData] = useState(null);
  useEffect(() => {
    Promise.all([api.get("/dashboard"), api.get("/invoices"), api.get("/payments")])
      .then(([dashboard, invoices, payments]) => setData({ dashboard: dashboard.data, invoices: Array.isArray(invoices.data) ? invoices.data : [], payments: Array.isArray(payments.data) ? payments.data : [] }))
      .catch(() => setData({ dashboard: null, invoices: [], payments: [] }));
  }, []);
  const totals = useMemo(() => {
    if (!data) return null;
    const invoiced = data.invoices.reduce((sum, invoice) => sum + invoice.TotalAmount, 0);
    const collected = data.payments.reduce((sum, payment) => sum + payment.Amount, 0);
    return { invoiced, collected, balance: Math.max(invoiced - collected, 0) };
  }, [data]);
  if (!data || !totals) return <p>Loading dashboard…</p>;
  const dashboard = data.dashboard || {};
  const collectionRate = totals.invoiced > 0 ? Math.round((totals.collected / totals.invoiced) * 100) : 0;
  const summaryDate = new Intl.DateTimeFormat("en-IN", { day: "numeric", month: "short", year: "numeric" }).format(new Date());
  return <section className="dashboard"><div className="dashboard-titlebar"><div className="dashboard-title-content"><span className="dashboard-title-mark">G</span><div><p className="dashboard-kicker">Garago</p><h2>Overview</h2></div></div><span>Updated {summaryDate}</span></div><section className="dashboard-panel"><div className="dashboard-panel-heading"><div><h3>Work orders</h3><p>Current workshop activity</p></div></div><div className="dashboard-cards"><Card title="Pending" value={dashboard.pending_work_orders || 0} caption="Waiting to begin" icon="◷" /><Card title="In progress" value={dashboard.in_progress_work_orders || 0} caption="Jobs in the workshop" icon="◌" tone="teal" /><Card title="Completed" value={dashboard.completed_work_orders || 0} caption="Ready or delivered" icon="✓" tone="teal" /></div></section><section className="dashboard-panel finance-summary"><div className="dashboard-panel-heading finance-heading"><div><p className="eyebrow">PAYMENT OVERVIEW</p><h3>Collections</h3><p>Payments received compared with invoices raised.</p></div><div className="collection-rate"><strong>{collectionRate}%</strong><span>collected</span></div></div><div className="collection-track" aria-label={`${collectionRate}% of invoice value collected`}><span style={{ width: `${collectionRate}%` }} /></div><div className="finance-grid"><Card title="Total invoiced" value={rupees(totals.invoiced)} caption="Value billed" icon="▧" /><Card title="Collected" value={rupees(totals.collected)} caption="Payments received" icon="₹" tone="teal" /><Card title="Outstanding" value={rupees(totals.balance)} caption="Still to collect" icon="!" tone="danger" /></div></section><section className="dashboard-panel"><div className="dashboard-panel-heading"><div><h3>Invoice status</h3><p>{data.invoices.length} invoice{data.invoices.length === 1 ? "" : "s"} tracked</p></div></div><div className="dashboard-cards"><Card title="Unpaid" value={dashboard.unpaid_invoices || 0} caption="No payment recorded" icon="!" tone="danger" /><Card title="Partially paid" value={dashboard.partial_invoices || 0} caption="Balance remains" icon="◐" tone="danger" /><Card title="Paid" value={dashboard.paid_invoices || 0} caption="Fully settled" icon="✓" tone="teal" /></div></section></section>;
}
export default Dashboard;
