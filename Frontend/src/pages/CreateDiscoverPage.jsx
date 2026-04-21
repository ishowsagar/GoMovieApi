import { Link } from "react-router-dom";
export function CreateDiscoverPage() {
  return (
    <>
      <h1>Discovery page goes here...</h1>;
      <button
        style={{
          padding: "1px",
          fontStyle: "bold",
          color: "red",
          marginLeft: "10rem",
        }}
      >
        {" "}
        <Link to="/">Back🔙</Link>
      </button>
    </>
  );
}
