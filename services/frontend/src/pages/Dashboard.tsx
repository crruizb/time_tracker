import React, { useEffect, useState } from "react";

// Define the structure of the response data
interface Task {
  id: string;
  name: string;
  description: string;
  username: string | null;
  startedAt: string | null;
  finishedAt: string | null;
}

interface Project {
  id: string;
  name: string;
  description: string;
  tasks: Task[];
}

interface ApiResponse {
  report: Project[];
}

const FetchDataComponent: React.FC = () => {
  // State to store the fetched data and loading/error states
  const [data, setData] = useState<ApiResponse | null>(null);
  console.log(data?.report);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch("http://localhost:8080/api/projects", {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include", // If you need to send cookies with the request
        });

        if (!response.ok) {
          throw new Error("Network response was not ok");
        }

        // Parse the JSON response
        const result: ApiResponse = await response.json();
        console.log(result);
        setData(result); // Set the data in the state
      } catch (error) {
        setError("There was a problem with the fetch operation");
        console.error("Error:", error);
      } finally {
        setLoading(false); // Set loading to false after the request is done
      }
    };

    fetchData();
  }, []); // Empty dependency array, so it only runs on mount

  if (loading) return <div>Loading...</div>;
  if (error) return <div>{error}</div>;

  return (
    <div>
      <h1>Fetched Data:</h1>
      {data?.report.map((p) => {
        console.log(p);
        return (
          <div key={p.id}>
            <h1>Project ID: {p.id}</h1>
            <p>Project Name: {p.name}</p>
          </div>
        );
      })}
    </div>
  );
};

export default FetchDataComponent;
