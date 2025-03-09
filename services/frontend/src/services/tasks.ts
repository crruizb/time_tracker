export const createTask = async (newTask: {
  name: string;
  description: string;
  projectId: string;
}) => {
  const response = await fetch(
    `http://localhost:8080/api/projects/${newTask.projectId}/tasks`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name: newTask.name,
        description: newTask.description,
      }),
      credentials: "include",
    }
  );

  if (!response.ok) throw new Error("Error adding new task");
  return response.json();
};
