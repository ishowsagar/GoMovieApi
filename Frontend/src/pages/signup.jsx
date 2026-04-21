import { useState } from "react";
import { Link } from "react-router-dom";

export function SignupPage() {
  const formIntialState = {
    username: "",
    email: "",
    password: "",
    bio: "",
  };

  const [formData, setFormData] = useState(formIntialState);
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

    const payload = {
      username: formData.username.trim(),
      email: formData.email.trim(),
      password: formData.password,
      bio: formData.bio.trim(),
    };

    if (!payload.username || !payload.email || !payload.password) {
      setError("Username, email, and password are required.");
      return;
    }

    try {
      setSubmitting(true);
      const res = await fetch("/api/users/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        let message = "Signup failed.";
        try {
          const body = await res.json();
          message = body?.error || body?.status || message;
        } catch {
          // No-op fallback when backend does not return JSON
        }
        throw new Error(message);
      }

      setSuccess("Account created successfully.");
      setFormData(formIntialState);
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
        <h2>Sign Up</h2>
        <form className="movie-form" onSubmit={handleSubmit}>
          <label htmlFor="signup-username">Username</label>
          <input
            id="signup-username"
            name="username"
            type="text"
            value={formData.username}
            onChange={handleChange}
          />

          <label htmlFor="signup-email">Email</label>
          <input
            id="signup-email"
            name="email"
            type="email"
            value={formData.email}
            onChange={handleChange}
          />

          <label htmlFor="signup-password">Password</label>
          <input
            id="signup-password"
            name="password"
            type="password"
            value={formData.password}
            onChange={handleChange}
          />

          <label htmlFor="signup-bio">Bio</label>
          <textarea
            id="signup-bio"
            name="bio"
            value={formData.bio}
            onChange={handleChange}
          />

          <button type="submit" disabled={submitting}>
            {submitting ? "Creating account..." : "Create account"}
          </button>
        </form>

        {error && <p className="error">{error}</p>}
        {success && <p>{success}</p>}

        <p>
          Already have an account? <Link to="/login">Login</Link>
        </p>
        <p>
          <Link to="/">Back to home</Link>
        </p>
      </main>
    </div>
  );
}
