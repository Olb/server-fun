import { useState } from "react";
import { useNavigate } from "react-router-dom";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

function CreatePost() {
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async () => {
    const response = await fetch(`${API_URL}/posts`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title, body }),
    });

    if (response.ok) navigate("/");
  };

  return (
    <div>
      <h2>Create New Post</h2>
      <input
        type="text"
        placeholder="Post Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
      />
      <ReactQuill value={body} onChange={setBody} />
      <button onClick={handleSubmit}>Submit</button>
    </div>
  );
}

export default CreatePost;
