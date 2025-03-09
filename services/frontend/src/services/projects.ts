import { Projects } from "@/types/projects";

export const fetchProjects = (): Promise<Projects> => {
  return fetch("http://localhost:8080/api/projects", {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
  })
    .then(async (res) => {
      if (res.status == 403) {
        throw new Error("Unauthorized user");
      }
      if (!res.ok) {
        throw new Error("Error fetching user projects. Try again later...");
      }

      return await res.json();
    })
    .then((res) => {
      return res;
    });
};

export const createProject = async (newProject: {
  name: string;
  description: string;
}) => {
  const response = await fetch("http://localhost:8080/api/projects", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      name: newProject.name,
      description: newProject.description,
    }),
    credentials: "include",
  });

  if (!response.ok) throw new Error("Error adding new project");
  return response.json();
};
