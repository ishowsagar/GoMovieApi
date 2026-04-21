import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { MoviesProvider } from "./context/MoviesContext.jsx";
import { RequireAuth } from "./components/RequireAuth.jsx";
import { CreateMoviePage } from "./pages/CreateMoviePage.jsx";
import { CreateDiscoverPage } from "./pages/CreateDiscoverPage.jsx";
import { LoginPage } from "./pages/login.jsx";
import { SignupPage } from "./pages/signup.jsx";
import App from "./App.jsx";
import "./styles/globals.css";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <BrowserRouter>
      <MoviesProvider>
        {/* rendering children and giving them access to the context */}
        <Routes>
          <Route path="/" element={<App />}></Route>
          <Route path="/login" element={<LoginPage />}></Route>
          <Route path="/signup" element={<SignupPage />}></Route>
          <Route element={<RequireAuth />}>
            <Route path="/create-movie" element={<CreateMoviePage />}></Route>
            <Route
              path="/discover-movies"
              element={<CreateDiscoverPage />}
            ></Route>
          </Route>
        </Routes>
      </MoviesProvider>
    </BrowserRouter>
  </StrictMode>,
);
