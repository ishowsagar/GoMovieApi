import { Link } from "react-router-dom";
import { useMovies } from "./context/MoviesContext";
import { clearStoredToken, getAuthHeader, isAuthenticated } from "./utils/auth";
export default function App() {
  const {
    health,
    movies,
    loading,
    error,
    setError,
    setMovies,
    deletingMovieID,
    setDeletingMovieID,
  } = useMovies();

  const loggedIn = isAuthenticated();

  function handleLogout() {
    clearStoredToken();
    setError("");
  }

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
          ...getAuthHeader(),
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

      <main className="form">
        <div className="form-div">
          <nav>
            <ul className="form-div-ul">
              <li>
                <Link to="/" className="form-div-ul-a">
                  Home
                </Link>
              </li>
              <li>
                <Link to="/discover-movies" className="form-div-ul-a">
                  Discover
                </Link>
              </li>
              <li>
                <Link to="/create-movie" className="form-div-ul-a">
                  Create
                </Link>
              </li>
              <li>
                <Link to="/login" className="form-div-ul-a">
                  Login
                </Link>
              </li>
              {loggedIn && (
                <li>
                  <button
                    type="button"
                    className="form-div-ul-a"
                    onClick={handleLogout}
                  >
                    Logout
                  </button>
                </li>
              )}
              <li>
                <Link to="/signup" className="form-div-ul-a">
                  Signup
                </Link>
              </li>
            </ul>
          </nav>
        </div>
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
