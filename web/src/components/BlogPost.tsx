import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";

interface Post {
  id: number;
  title: string;
  body: string;
  createdAt: Date;
}

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

function BlogPost() {
  const { id } = useParams();
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${API_URL}/posts/${id}`)
      .then((res) => res.json())
      .then((data) => {
        setPost({
          ...data.post,
          createdAt: new Date(data.post.created_at), 
        });
        setLoading(false);
      })
      .catch(console.error);
  }, [id]);

  if (loading) return <p>Loading...</p>;
  if (!post) return <p>Post not found.</p>;

  return (
    <div className="blog-post">
      <h1>{post.title}</h1>
      <p>Created at: {post.createdAt.toLocaleString()}</p>
      <div dangerouslySetInnerHTML={{ __html: post.body }} /> 
      <Link to={`/edit/${post.id}`}>Edit Post</Link>
    </div>
  );
}

export default BlogPost;
