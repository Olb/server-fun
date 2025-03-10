import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import BlogList from "./components/BlogList";
import BlogPost from "./components/BlogPost";
import CreatePost from "./components/CreatePost";
import EditPost from "./components/EditPost"; 
import Sidebar from "./components/Sidebar";

function App() {
  return (
    <Router>
      <div className="app-container">
        <header>
          <h1>My Blog</h1>
          <Link to="/create" className="create-button">+ Create Post</Link>
        </header>
        <div className="content">
          <Sidebar />
          <Routes>
            <Route path="/" element={<BlogList />} />
            <Route path="/post/:id" element={<BlogPost />} />
            <Route path="/create" element={<CreatePost />} />
            <Route path="/edit/:id" element={<EditPost />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
}

export default App;
