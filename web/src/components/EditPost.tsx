import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import ReactQuill from "react-quill";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

function EditPost() {
  const { id } = useParams();
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    fetch(`${API_URL}/posts/${id}`)
      .then((res) => res.json())
      .then((data) => {
        setTitle(data.post.title);
        setBody(data.post.body);
      });
  }, [id]);

  const handleUpdate = async () => {
    await fetch(`${API_URL}/posts/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title, body }),
    });
    navigate(`/post/${id}`);
  };

  return (
    <div>
      <h2>Edit Post</h2>
      <input value={title} onChange={(e) => setTitle(e.target.value)} />
      <ReactQuill value={body} onChange={setBody} />
      <button onClick={handleUpdate}>Update</button>
    </div>
  );
}

export default EditPost;
