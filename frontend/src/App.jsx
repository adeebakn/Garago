import { lazy, Suspense } from "react";
import { BrowserRouter, Navigate, Routes, Route } from "react-router-dom";
import Header from "./components/Header";
import Sidebar from "./components/Sidebar";
import { AuthProvider } from "./auth/AuthContext";
import { useAuth } from "./auth/auth-context";
import { ConfirmProvider } from "./ui/ConfirmContext";
import { OrganizationProvider } from "./organization/OrganizationContext";

const Dashboard = lazy(() => import("./components/Dashboard"));
const Customers = lazy(() => import("./components/Customers"));
const Vehicles = lazy(() => import("./components/Vehicles"));
const Services = lazy(() => import("./components/Services"));
const Quotations = lazy(() => import("./components/Quotations"));
const WorkOrders = lazy(() => import("./components/Work_Orders"));
const Invoices = lazy(() => import("./components/Invoices"));
const Payments = lazy(() => import("./components/Payments"));
const OrganizationSettings = lazy(() => import("./components/OrganizationSettings"));
const Login = lazy(() => import("./auth/Login"));
const Signup = lazy(() => import("./auth/Signup"));
const ForgotPassword = lazy(() => import("./auth/ForgotPassword"));
const ResetPassword = lazy(() => import("./auth/ResetPassword"));


function AppRoutes() {
  const { user, loading } = useAuth();
  if (loading) return <main className="auth-page">Checking your session…</main>;
  if (!user) return <Suspense fallback={<main className="auth-page">Loading…</main>}><Routes><Route path="/login" element={<Login />} /><Route path="/signup" element={<Signup />} /><Route path="/forgot-password" element={<ForgotPassword />} /><Route path="/reset-password" element={<ResetPassword />} /><Route path="*" element={<Navigate to="/login" replace />} /></Routes></Suspense>;
  return (
      <div className="app">
        <Header />

        <div className="layout">
          <Sidebar />

          <main className="content">
            <div className="page-transition">
            <Suspense fallback={<p>Loading module…</p>}><Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/customers" element={<Customers />} />
              <Route path="/vehicles" element={<Vehicles />} />
              <Route path="/services" element={<Services />} />
              <Route path="/quotations" element={<Quotations />} />
              <Route path="/work-orders" element={<WorkOrders />} />
              <Route path="/invoices" element={<Invoices />} />
              <Route path="/payments" element={<Payments />} />
              <Route path="/settings" element={<OrganizationSettings />} />
            </Routes></Suspense>
            </div>
          </main>
        </div>
      </div>
  );
}

function App() { return <BrowserRouter><OrganizationProvider><AuthProvider><ConfirmProvider><AppRoutes /></ConfirmProvider></AuthProvider></OrganizationProvider></BrowserRouter>; }

export default App;
