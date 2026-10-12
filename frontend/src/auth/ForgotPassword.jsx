import { useState } from "react";
import { Link } from "react-router-dom";
import api from "../api";
import { apiErrorMessage } from "../utils/api-error";
import "./auth.css";

export default function ForgotPassword() {
  const [email, setEmail] = useState("");
  const [error, setError] = useState("");
  const [sent, setSent] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  async function submit(event) {
    event.preventDefault();
    setError("");
    setSubmitting(true);
    try {
      await api.post("/auth/forgot-password", { email });
      setSent(true);
    } catch (requestError) {
      setError(apiErrorMessage(requestError, "Unable to request a reset email. Please try again."));
    } finally {
      setSubmitting(false);
    }
  }

  return <main className="auth-page"><section className="auth-intro"><span className="brand-mark">Garago</span><h1>Get back to work.</h1><p>We’ll send a secure password reset link to your email address.</p></section><form className="auth-card" onSubmit={submit}><p className="eyebrow">PASSWORD RESET</p><h2>Forgot your password?</h2>{sent ? <p className="success">If an account exists for this email, a reset link has been sent. Check your inbox.</p> : <><p className="auth-help">Enter the email address you use to sign in.</p>{error && <p className="error">{error}</p>}<label>Email<input required type="email" autoComplete="email" value={email} onChange={(event) => setEmail(event.target.value)} /></label><button disabled={submitting}>{submitting ? "Sending…" : "Send reset link"}</button></>}<p className="auth-switch"><Link to="/login">Back to sign in</Link></p></form></main>;
}
