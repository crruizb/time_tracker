import Cookies from "js-cookie";
import "./index.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { Button } from "./components/ui/button";
import Dashboard from "./pages/Dashboard";

function App() {
  const username = Cookies.get("username");
  console.log(username);

  const handleLogin = () => {
    window.location.href = "http://localhost:8080/auth/github/login";
  };

  return (
    <Router>
      <Routes>
        <Route
          path="/"
          element={
            username ? (
              <div>Welcome, {username}!</div>
            ) : (
              <div>
                <h1>OAuth2 Login</h1>
                <Button onClick={handleLogin}>Login with OAuth2</Button>
              </div>
            )
          }
        />
        <Route path="/dashboard" Component={Dashboard} />
      </Routes>
    </Router>
  );
}

export default App;
