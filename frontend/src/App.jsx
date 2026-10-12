import { BrowserRouter, Navigate, Routes, Route } from "react-router-dom";
import Header from "./components/Header";
import Sidebar from "./components/Sidebar";
import Dashboard from "./components/Dashboard";
import Customers from "./components/Customers";
import Vehicles from "./components/Vehicles";
import Services from "./components/Services";
import Quotations from "./components/Quotations";
import WorkOrders from "./components/Work_Orders";
import Invoices from "./components/Invoices";
import Payments from "./components/Payments";
import { AuthProvider, useAuth } from "./auth/AuthContext";
import Login from "./auth/Login";
import Signup from "./auth/Signup";
import ForgotPassword from "./auth/ForgotPassword";
import ResetPassword from "./auth/ResetPassword";
import { ConfirmProvider } from "./ui/ConfirmContext";
import { OrganizationProvider } from "./organization/OrganizationContext";
import OrganizationSettings from "./components/OrganizationSettings";


function AppRoutes() {
  const { user, loading } = useAuth();
  if (loading) return <main className="auth-page">Checking your session…</main>;
  if (!user) return <Routes><Route path="/login" element={<Login />} /><Route path="/signup" element={<Signup />} /><Route path="/forgot-password" element={<ForgotPassword />} /><Route path="/reset-password" element={<ResetPassword />} /><Route path="*" element={<Navigate to="/login" replace />} /></Routes>;
  return (
      <div className="app">
        <Header />

        <div className="layout">
          <Sidebar />

          <main className="content">
            <div className="page-transition">
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/customers" element={<Customers />} />
              <Route path="/vehicles" element={<Vehicles />} />
              <Route path="/services" element={<Services />} />
              <Route path="/quotations" element={<Quotations />} />
              <Route path="/work-orders" element={<WorkOrders />} />
              <Route path="/invoices" element={<Invoices />} />
              <Route path="/payments" element={<Payments />} />
              <Route path="/settings" element={<OrganizationSettings />} />
            </Routes>
            </div>
          </main>
        </div>
      </div>
  );
}

function App() { return <BrowserRouter><OrganizationProvider><AuthProvider><ConfirmProvider><AppRoutes /></ConfirmProvider></AuthProvider></OrganizationProvider></BrowserRouter>; }

export default App;
