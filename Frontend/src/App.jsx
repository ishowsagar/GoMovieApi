import { useEffect, useState } from "react";

export default function App() {
  const [health, setHealth] = useState("Checking backend...");
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function loadData() {
      try {
        const [healthRes, moviesRes] = await Promise.all([
          fetch("/health"),
          fetch("/api/movies/all"),
        ]);

        const healthText = await healthRes.text();
        setHealth(healthText || "Backend reachable.");

        if (!moviesRes.ok) {
          throw new Error("Failed to load movies from API.");
        }

        const moviesData = await moviesRes.json();
        setMovies(moviesData?.data ?? []);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Unexpected error");
      } finally {
        setLoading(false);
      }
    }

    loadData();
  }, []);

  return (
    <div className="page">
      <video
        className="bg-video"
        autoPlay
        muted
        loop
        playsInline
        preload="auto"
      >
        <source src="/netflix.mp4" type="video/mp4" />
      </video>
      <div className="bg-overlay" aria-hidden="true" />
      <header className="hero">
        <h1>MoviesFlix 🎬</h1>
        {error && <p>{health}</p>}
      </header>

      <main className="content">
        <h2>Available movies</h2>
        {loading && <p>Loading movies...</p>}
        {error && <p className="error">{error}</p>}

        {!loading && !error && movies.length === 0 && (
          <p>No movies yet. Create one from your backend API.</p>
        )}

        {!loading && !error && movies.length > 0 && (
          <ul className="coffee-list">
            {movies.map((movie) => (
              <li key={movie.id || movie.name} className="coffee-card">
                <h3>{movie.name}</h3>
                <p>
                  Genre: <strong>{movie.genre || "Unknown"}</strong>
                </p>
                <p>
                  Description: <strong>{movie.description ?? "-"}</strong>
                </p>
                <p>
                  Ratings: <strong>{movie.ratings ?? "-"}</strong>
                </p>
                <p>
                  Created-At: <strong>{movie.created_at ?? "-"}</strong>
                </p>
                <p>
                  Updated-At: <strong>{movie.updated_at ?? "-"}</strong>
                </p>
              </li>
            ))}
          </ul>
        )}
      </main>
    </div>
  );
}
