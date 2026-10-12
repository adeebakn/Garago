import { jsPDF } from "jspdf";
import FormDrawer from "./FormDrawer";
import { useOrganization } from "../organization/organization-context";

const money = (amount) => `₹${Number(amount || 0).toLocaleString("en-IN", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
const pdfMoney = (amount) => `INR ${Number(amount || 0).toLocaleString("en-IN", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;

export default function PaymentReceipt({ payment, paid, balance, onClose }) {
  const { organization } = useOrganization();
  const date = payment.PaymentDate?.Time || payment.PaymentDate || null;
  const formattedDate = date ? new Intl.DateTimeFormat("en-IN", { day: "2-digit", month: "short", year: "numeric" }).format(new Date(date)) : "—";
  const contact = [organization.address, organization.phone, organization.email].filter(Boolean).join(" · ");

  const download = () => {
    const pdf = new jsPDF({ unit: "mm", format: "a5" });
    const teal = [8, 125, 112];
    const ink = [23, 51, 45];
    const muted = [92, 112, 106];
    const line = [220, 232, 228];
    const left = 14; const right = 134;
    pdf.setFillColor(...teal); pdf.rect(0, 0, 148, 32, "F");
    pdf.setTextColor(255, 255, 255); pdf.setFont("helvetica", "bold"); pdf.setFontSize(17); pdf.text(organization.name || "Garago", left, 15);
    pdf.setFont("helvetica", "normal"); pdf.setFontSize(8.5); pdf.text("PAYMENT RECEIPT", left, 22);
    if (contact) { pdf.setFontSize(7.5); pdf.text(pdf.splitTextToSize(contact, 118), left, 28); }
    pdf.setTextColor(...ink); pdf.setFont("helvetica", "bold"); pdf.setFontSize(16); pdf.text("Payment received", left, 46);
    pdf.setDrawColor(...line); pdf.line(left, 52, right, 52);
    const row = (label, value, y) => { pdf.setTextColor(...muted); pdf.setFont("helvetica", "normal"); pdf.setFontSize(9); pdf.text(label, left, y); pdf.setTextColor(...ink); pdf.setFont("helvetica", "bold"); pdf.setFontSize(10.5); pdf.text(value, right, y, { align: "right" }); pdf.setDrawColor(...line); pdf.line(left, y + 5, right, y + 5); };
    row("Payment date", formattedDate, 63);
    row("Payment method", payment.PaymentMethod, 75);
    pdf.setFillColor(230, 247, 242); pdf.roundedRect(left, 87, 120, 20, 3, 3, "F");
    pdf.setTextColor(...teal); pdf.setFont("helvetica", "normal"); pdf.setFontSize(9); pdf.text("Amount received", 20, 96);
    pdf.setFont("helvetica", "bold"); pdf.setFontSize(14); pdf.text(pdfMoney(payment.Amount), 128, 101, { align: "right" });
    pdf.setTextColor(...ink); pdf.setFont("helvetica", "bold"); pdf.setFontSize(10); pdf.text("Payment summary", left, 120);
    row("Total paid to date", pdfMoney(paid), 131);
    row("Remaining balance", pdfMoney(balance), 143);
    pdf.setTextColor(...muted); pdf.setFont("helvetica", "normal"); pdf.setFontSize(8.5); pdf.text("Thank you. Keep this receipt for your records.", left, 157);
    pdf.save(`payment-receipt-${String(formattedDate).replaceAll(" ", "-")}.pdf`);
  };

  return (
    <FormDrawer open title="Payment receipt" onClose={onClose} backdropClassName="receipt-backdrop" drawerClassName="receipt-sheet">
      <div className="receipt-preview">
        <article className="receipt-document">
          <header className="receipt-brand">
            <div><strong>{organization.name || "Garago"}</strong><span>{contact || "Payment receipt"}</span></div>
          </header>
          <div className="receipt-title"><div><p>PAYMENT RECEIPT</p><h1>Payment received</h1></div></div>
          <table className="receipt-table" aria-label="Payment details"><tbody>
            <tr><th>Payment date</th><td>{formattedDate}</td></tr>
            <tr><th>Payment method</th><td>{payment.PaymentMethod}</td></tr>
          </tbody></table>
          <div className="receipt-amount"><span>Amount received</span><strong>{money(payment.Amount)}</strong></div>
          <h2 className="receipt-section-title">Payment summary</h2>
          <table className="receipt-table receipt-summary-table" aria-label="Payment summary"><tbody>
            <tr><th>Total paid to date</th><td>{money(paid)}</td></tr>
            <tr><th>Remaining balance</th><td className={balance > 0 ? "receipt-balance-due" : "receipt-balance-settled"}>{money(balance)}</td></tr>
          </tbody></table>
          <p className="receipt-note">Thank you. Keep this receipt for your records.</p>
        </article>
        <div className="receipt-actions no-print"><button type="button" onClick={download}>Download PDF</button><button type="button" className="print-button" onClick={() => window.print()}>Print receipt</button></div>
      </div>
    </FormDrawer>
  );
}
