import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

interface Post {
  id: number;
  title: string;
}

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

function Sidebar() {
  const [recentPosts, setRecentPosts] = useState<Post[]>([]);

  useEffect(() => {
    fetch(`${API_URL}/posts`)
      .then((res) => res.json())
      .then((data) => {
        setRecentPosts(data.posts.slice(-5)); // Last 5 posts
      })
      .catch(console.error);
  }, []);

  return (
    <aside className="sidebar">
      <h3>Recent Posts</h3>
      <ul>
        {recentPosts.map((post) => (
          <li key={post.id}>
            <Link to={`/post/${post.id}`}>{post.title}</Link>
          </li>
        ))}
      </ul>
    </aside>
  );
}

export default Sidebar;
