import { jsPDF } from "jspdf";
import FormDrawer from "./FormDrawer";
import { useOrganization } from "../organization/organization-context";

const money = (amount) => `₹${Number(amount || 0).toLocaleString("en-IN", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;

export default function QuotationPreview({ quotation, items, customer, vehicle, services, onClose, documentTitle = "Quotation", documentCode = "" }) {
  const { organization } = useOrganization();
  const date = quotation.CreatedAt ? new Intl.DateTimeFormat("en-IN", { day: "2-digit", month: "short", year: "numeric" }).format(new Date(quotation.CreatedAt)) : "—";
  const contact = [organization.address, organization.phone, organization.email].filter(Boolean).join(" · ");
  const lines = items.map((item) => ({ ...item, service: services.find((service) => String(service.ID) === String(item.ServiceID)) }));

  const download = () => {
    const pdf = new jsPDF({ unit: "mm", format: "a4" });
    const teal = [8, 125, 112]; const ink = [23, 51, 45]; const muted = [92, 112, 106]; const line = [220, 232, 228];
    const left = 18; const right = 192;
    pdf.setFillColor(...teal); pdf.rect(0, 0, 210, 34, "F");
    pdf.setTextColor(255, 255, 255); pdf.setFont("helvetica", "bold"); pdf.setFontSize(18); pdf.text(organization.name || "Garago", left, 16);
    pdf.setFont("helvetica", "normal"); pdf.setFontSize(8.5); pdf.text(documentTitle.toUpperCase(), left, 23);
    if (contact) { pdf.setFontSize(7.5); pdf.text(pdf.splitTextToSize(contact, 165), left, 29); }
    pdf.setTextColor(...ink); pdf.setFont("helvetica", "bold"); pdf.setFontSize(18); pdf.text(documentTitle, left, 50);
    pdf.setFont("helvetica", "normal"); pdf.setTextColor(...muted); pdf.setFontSize(9); pdf.text(`${documentCode || `Issued: ${date}`}`, right, 50, { align: "right" });
    pdf.setDrawColor(...line); pdf.line(left, 56, right, 56);
    pdf.setTextColor(...ink); pdf.setFont("helvetica", "bold"); pdf.setFontSize(10); pdf.text("Customer", left, 68); pdf.text("Vehicle", 110, 68);
    pdf.setFont("helvetica", "normal"); pdf.setFontSize(9); pdf.text(customer?.Name || "—", left, 75); pdf.text(vehicle?.RegNo || "—", 110, 75);
    pdf.setTextColor(...muted); pdf.setFontSize(8); pdf.text(customer?.Phone || "", left, 81); pdf.text(vehicle?.Vehicle || "", 110, 81);
    let y = 99;
    pdf.setFillColor(244, 248, 247); pdf.rect(left, y - 7, 174, 10, "F");
    pdf.setTextColor(...muted); pdf.setFont("helvetica", "bold"); pdf.setFontSize(8); pdf.text("SERVICE", left + 4, y - 1); pdf.text("QTY", 132, y - 1, { align: "right" }); pdf.text("AMOUNT", right - 4, y - 1, { align: "right" });
    y += 10;
    lines.forEach((item) => { pdf.setTextColor(...ink); pdf.setFont("helvetica", "normal"); pdf.setFontSize(9); pdf.text(item.service?.Name || "Service", left + 4, y); pdf.text(String(item.Quantity), 132, y, { align: "right" }); pdf.text(money(item.Amount), right - 4, y, { align: "right" }); pdf.setDrawColor(...line); pdf.line(left, y + 5, right, y + 5); y += 11; });
    y += 8; pdf.setTextColor(...ink); pdf.setFont("helvetica", "bold"); pdf.setFontSize(11); pdf.text(`${documentTitle} total`, 130, y); pdf.setTextColor(...teal); pdf.setFontSize(13); pdf.text(money(quotation.TotalAmount), right, y, { align: "right" });
    pdf.setTextColor(...muted); pdf.setFont("helvetica", "normal"); pdf.setFontSize(8.5); pdf.text(documentTitle === "Work order" ? "Services to be completed for this work order." : "This quotation is subject to confirmation before work begins.", left, y + 18);
    pdf.save(`${documentTitle.toLowerCase().replaceAll(" ", "-")}-${String(documentCode || date).replaceAll(" ", "-")}.pdf`);
  };

  return <FormDrawer open title={`${documentTitle} preview`} onClose={onClose} backdropClassName="quotation-backdrop" drawerClassName="quotation-sheet"><div className="quotation-preview"><article className="quotation-document"><header className="quotation-brand"><span className="quotation-mark">◆</span><div><strong>{organization.name || "Garago"}</strong><span>{contact || documentTitle}</span></div></header><div className="quotation-document-title"><div><p>{documentTitle.toUpperCase()}</p><h1>{documentTitle}</h1></div><span>{documentCode || date}</span></div><div className="quotation-party-grid"><div><span>Customer</span><strong>{customer?.Name || "—"}</strong><small>{customer?.Phone || customer?.Email || ""}</small></div><div><span>Vehicle</span><strong>{vehicle?.RegNo || vehicle?.Vehicle || "—"}</strong><small>{vehicle?.Vehicle || ""}</small></div></div><table className="quotation-table"><thead><tr><th>Service</th><th>Qty</th><th>Amount</th></tr></thead><tbody>{lines.map((item) => <tr key={item.ID}><td>{item.service?.Name || "Service"}</td><td>{item.Quantity}</td><td>{money(item.Amount)}</td></tr>)}</tbody></table><div className="quotation-total"><span>{documentTitle} total</span><strong>{money(quotation.TotalAmount)}</strong></div>{Number(quotation.AdvanceRequested || 0) > 0 && <div className="quotation-advance-preview"><span>Advance requested</span><strong>{money(quotation.AdvanceRequested)}</strong><span>Balance after advance</span><strong>{money(quotation.TotalAmount - quotation.AdvanceRequested)}</strong></div>}<p className="quotation-note">{documentTitle === "Work order" ? "Services to be completed for this work order." : "This quotation is subject to confirmation before work begins."}</p></article><div className="receipt-actions no-print"><button type="button" onClick={download}>Download PDF</button><button type="button" className="print-button" onClick={() => window.print()}>Print {documentTitle.toLowerCase()}</button></div></div></FormDrawer>;
}
