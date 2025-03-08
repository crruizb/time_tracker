import React from "react";

const Home = () => {
  const handleLogin = () => {
    window.location.href = "http://localhost:8080/auth/github/login";
  };

  return (
    <div className="p-4">
      <h1 className="text-2xl mb-4">OAuth2 Login with React</h1>
      <button
        onClick={handleLogin}
        className="bg-blue-500 text-white px-4 py-2 rounded"
      >
        Login with OAuth2
      </button>
    </div>
  );
};

export default Home;
