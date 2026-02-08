import { useState } from "react";
import api from "../../api"; // make sure path is correct
import "./index.css"; // your CSS

const Login = ({ setToken }) => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [isSignup, setIsSignup] = useState(false); // toggle login/signup

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      if (isSignup) {
        // Sign Up
        const res = await api.post("/users", { username, password });
        alert("User created! Logging you in...");

        // Automatically login after signup
        const loginRes = await api.post("/users/login", { username, password });
        const token = loginRes.data.token || loginRes.data.Token;

        if (!token) {
          return alert("Login failed after signup. Please try logging in manually.");
        }

        localStorage.setItem("token", token);
        setToken(token);
        window.location.href = "/items";
      } else {
        // Login
        const res = await api.post("/users/login", { username, password });
        const token = res.data.token || res.data.Token;

        if (!token) {
          return alert(res.data.error || "Invalid username/password");
        }

        localStorage.setItem("token", token);
        setToken(token);
        window.location.href = "/items";
      }
    } catch (err) {
      console.error(err);
      alert(err.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="login-page">
      <div className="login-card">
        {/* Logo letters */}
        <div className="login-logo">
          <span className="login-logo-span">A</span>
          <span className="login-logo-span">B</span>
          <span className="login-logo-span">C</span>
          <span className="login-logo-span">D</span>
          <span className="login-logo-span">E</span>
        </div>

        <h2 style={{ textAlign: "center", marginBottom: "16px" }}>
          {isSignup ? "Sign Up" : "Login"}
        </h2>

        <form className="login-card-form" onSubmit={handleSubmit}>
          <label className="login-card-label">Username</label>
          <input
            className="login-card-input"
            type="text"
            placeholder="Enter username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />

          <label className="login-card-label">Password</label>
          <input
            className="login-card-input"
            type="password"
            placeholder="Enter password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />

          <button className="login-btn" type="submit">
            {isSignup ? "Sign Up" : "Login"}
          </button>
        </form>

        <p
          style={{
            textAlign: "center",
            marginTop: "12px",
            cursor: "pointer",
            color: "#0000ff",
          }}
          onClick={() => setIsSignup(!isSignup)}
        >
          {isSignup
            ? "Already have an account? Login"
            : "Don't have an account? Sign Up"}
        </p>
      </div>
    </div>
  );
};

export default Login;
