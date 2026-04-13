import { useEffect, useState } from "react";

export default function App() {
  const [health, setHealth] = useState("Checking backend...");
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [deletingMovieID, setDeletingMovieID] = useState(null);

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

  // whatever id is passing makes an req to delete req path
  async function handleDelete(movieID) {
    if (!movieID) {
      setError("failed to delete the movie as movie id is missing!.");
      return;
    }
    // if movie id is available
    try {
      setError("");
      setDeletingMovieID(movieID);
      const fetchReq = await fetch(`/api/movies/movie/delete/${movieID}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
      });

      let backendRes = "failed to delete the movie.";

      // if fetch req on that url failed
      if (!fetchReq.ok) {
        try {
          const res = await fetchReq.json();
          backendRes = res?.status || backendRes;
        } catch (err) {
          throw new Error(
            err instanceof Error ? err.message : "unexpected error occurred",
          );
        }
        throw new Error(backendRes);
      }
      // setting Movies without this movie by filtering it out
      setMovies((prevMovies) =>
        prevMovies.filter((movie) => movie.id !== movieID),
      );
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "unexpected error occurred",
      );
    } finally {
      setDeletingMovieID(null); // set to null to wipe stored id in state
    }
  }

  return (
    <div className="page">

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
                <button
                  onClick={() => handleDelete(movie.id)}
                  className="delete-btn"
                >
                  {deletingMovieID === movie.id
                    ? "deleting movie"
                    : "Delete Movie"}
                </button>
              </li>
            ))}
          </ul>
        )}
      </main>
    </div>
  );
}
