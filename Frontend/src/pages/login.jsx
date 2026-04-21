import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useMovies } from "../context/MoviesContext";
import { setStoredToken } from "../utils/auth";

export function LoginPage() {
  const { refreshMovies } = useMovies();
  const navigate = useNavigate();
  const location = useLocation();
  const from = location.state?.from?.pathname || "/";

  const [formData, setFormData] = useState({
    username: "",
    password: "",
  });
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  function handleChange(e) {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setError("");
    setSuccess("");

    const username = formData.username.trim();
    const password = formData.password;

    if (!username || !password) {
      setError("Username and password are required.");
      return;
    }

    try {
      setSubmitting(true);
      const res = await fetch("/api/tokens/authentication", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      if (!res.ok) {
        let message = "Login failed.";
        try {
          const body = await res.json();
          message = body?.error || body?.status || message;
        } catch {
          // No-op fallback when backend does not return JSON
        }
        throw new Error(message);
      }

      const body = await res.json();
      const token = body?.auth_token?.token || body?.token || "";
      if (!token) {
        throw new Error("Login succeeded but token was missing in response.");
      }

      setStoredToken(token);
      setSuccess("Login successful.");
      await refreshMovies();
      navigate(from, { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unexpected error");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="page">
      <div className="bg-overlay" aria-hidden="true" />
      <main className="content">
        <h2>Login</h2>
        <form className="movie-form" onSubmit={handleSubmit}>
          <label htmlFor="login-username">Username</label>
          <input
            id="login-username"
            name="username"
            type="text"
            value={formData.username}
            onChange={handleChange}
          />

          <label htmlFor="login-password">Password</label>
          <input
            id="login-password"
            name="password"
            type="password"
            value={formData.password}
            onChange={handleChange}
          />

          <button type="submit" disabled={submitting}>
            {submitting ? "Logging in..." : "Login"}
          </button>
        </form>

        {error && <p className="error">{error}</p>}
        {success && <p>{success}</p>}

        <p>
          No account? <Link to="/signup">Create one</Link>
        </p>
        <p>
          <Link to="/">Back to home</Link>
        </p>
      </main>
    </div>
  );
}
