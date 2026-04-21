import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { getAuthHeader } from "../utils/auth";

//   * Creating Shared Context to React app
const MoviesApiDataContext = createContext(null);
export function MoviesProvider({ children }) {
  // ! recieving props as children like <>children goes here</>
  const [health, setHealth] = useState("Checking backend...");
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [deletingMovieID, setDeletingMovieID] = useState(null);

  const value = useMemo(
    () => ({
      health,
      movies,
      loading,
      error,
      deletingMovieID,
      setMovies,
      setError,
      setDeletingMovieID,
      refreshMovies,
    }),
    [health, movies, loading, error, deletingMovieID],
  );

  async function refreshMovies() {
    setLoading(true);
    setError("");
    try {
      const [healthRes, moviesRes] = await Promise.all([
        fetch("/health"),
        fetch("/api/movies/all", {
          headers: {
            ...getAuthHeader(),
          },
        }),
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

  useEffect(() => {
    refreshMovies();
  }, []);

  return (
    // sharing this value as context to its children to use
    <MoviesApiDataContext.Provider value={value}>
      {children}
    </MoviesApiDataContext.Provider>
  );
}

// func that returns context to use it directly by invoking the func
export function useMovies() {
  const ctx = useContext(MoviesApiDataContext);

  if (!ctx) {
    throw new Error("useMovies must be used inside MoviesProvider");
  }
  return ctx;
}
