import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

interface Post {
    id: number;
    title: string;
    createdAt: Date;
}

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

function BlogList() {
    const [posts, setPosts] = useState<Post[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        fetch(`${API_URL}/posts`)
            .then((res) => res.json())
            .then((data) => {
                const formattedPosts = data.posts.map((post: any) => ({
                    id: post.id,
                    title: post.title,
                    createdAt: new Date(post.created_at),
                }));

                setPosts(formattedPosts);
                setLoading(false);
            })
            .catch((err) => {
                setError(err.message);
                setLoading(false);
            });
    }, []);

    return (
        <div className="blog-list">
            <h2>Blog Posts</h2>
            {loading && <p>Loading...</p>}
            {error && <p style={{ color: "red" }}>{error}</p>}
            <ul>
                {posts.map((post) => (
                    <li key={post.id}>
                        <Link to={`/post/${post.id}`}>
                            {post.title} - {post.createdAt.toLocaleDateString()}
                        </Link>
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default BlogList;
