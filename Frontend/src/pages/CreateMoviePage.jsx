import { Link } from "react-router-dom";
import { useState } from "react";
import { useMovies } from "../context/MoviesContext";
import { getAuthHeader } from "../utils/auth";

export function CreateMoviePage() {
  const { refreshMovies } = useMovies();
  // initializing state to store form ingress
  const initialFormData = {
    name: "",
    genre: "",
    description: "",
    ratings: "",
  };
  const [formData, setFormData] = useState(initialFormData);
  const { name, genre, description, ratings } = formData; // accessing fields from the state

  //@ form submission is done by form and who is submitting --> button and what is invoked by form 👇 this func
  async function handleSubmit(e) {
    e.preventDefault();
    console.log("submitted payload:", formData); //* natively have access to formdata when used on form

    // post req to send this formdata
    const url = `/api/movie/create`;
    const payloadData = {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name,
        genre,
        description,
        ratings: Number(ratings),
      }),
    };

    let errMsg = "this request is unreachable to backend";

    try {
      const request = await fetch(url, {
        method: payloadData.method,
        headers: {
          ...payloadData.headers,
          ...getAuthHeader(),
        },
        body: payloadData.body,
      });

      //   ! checking if fetch req was not successfull
      if (!request.ok) {
        try {
          const resp = await request.json();
          errMsg = resp?.status || errMsg;
        } catch (err) {
          errMsg =
            err instanceof Error ? err.message : "unexpected error occurred!.";
        }
        throw new Error(errMsg);
      }

      //* otherwise successfully sent req
      await refreshMovies(); // once resp is ok, refreshing all data from movies db including this movie which is just recently created
      setFormData(initialFormData); //* resetting form back to null values
    } catch (err) {
      console.error(err instanceof Error ? err.message : errMsg);
    }
  }

  //@ Handles form change input and set state for that change
  function handleFormDataChange(e) {
    const { name, value } = e.target; //! destructuring whatever inside e.target body object
    setFormData((prevFormData) => {
      return { ...prevFormData, [name]: value }; //* keeping old state intact + as we've already set name field, extracting which input is clicked --> its name for state name --> update it with value entered on that field from event.target
    });
    console.log(e.target);
  }
  return (
    <div className="page">
      <div className="bg-overlay" aria-hidden="true" />
      <main className="content">
        <h2>Create Movie</h2>
        <p>
          <Link to="/" className="page-link">
            Back to home
          </Link>
        </p>
        <div className="movie-form-div">
          <form onSubmit={handleSubmit} className="movie-form">
            <label htmlFor="mname">Movie Name:</label>
            <input
              id="mname"
              type="text"
              name="name" //@ setting name field on input to extract name whenever target is clicked , as we set name same for whatever was in state so it would be dynamically updating that only
              onChange={(event) => handleFormDataChange(event)}
              value={formData.name}
            ></input>
            <label htmlFor="mgenre">Movie Genre:</label>
            <input
              id="mgenre"
              type="text"
              name="genre"
              onChange={(event) => handleFormDataChange(event)}
              value={formData.genre}
            ></input>
            <label htmlFor="mdesc">Movie Description:</label>
            <textarea
              id="mdesc"
              name="description"
              onChange={(event) => handleFormDataChange(event)}
              value={formData.description}
            ></textarea>
            <label htmlFor="mratings">Movie Ratings:</label>
            <input
              id="mratings"
              type="number"
              name="ratings"
              onChange={(event) => handleFormDataChange(event)}
              value={formData.ratings}
            ></input>
            {/* <label htmlFor="mtime">Created At:</label>
          <input id="mtime" type="datetime-local"></input> */}
            <button type="submit">Create Movie</button>
          </form>
        </div>
      </main>
    </div>
  );
}
