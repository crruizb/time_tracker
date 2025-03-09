import { createProject, fetchProjects } from "@/services/projects";
import { createTask } from "@/services/tasks";
import { Project, Projects } from "@/types/projects";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export const useProjects = () => {
  const queryClient = useQueryClient();

  const { isLoading, isError, data, error } = useQuery({
    queryKey: ["projects"],
    queryFn: fetchProjects,
  });

  const { mutate: addProject } = useMutation({
    mutationFn: createProject,
    onSuccess: (newProject) => {
      queryClient.setQueryData(["projects"], (oldData: Projects) => ({
        ...oldData,
        report: [...(oldData?.report ?? []), { ...newProject, tasks: [] }],
      }));
    },
  });

  const { mutate: addTask } = useMutation({
    mutationFn: createTask,
    onSuccess: (newTask, variables) => {
      queryClient.setQueryData(["projects"], (oldData: Projects) => {
        const updatedProjects = oldData.report.map((project: Project) => {
          if (project.id === variables.projectId) {
            return {
              ...project,
              tasks: [...project.tasks, newTask],
            };
          }
          return project;
        });
        return { ...oldData, report: updatedProjects };
      });
    },
  });

  return {
    isLoading,
    isError,
    error,
    projects: data?.report ?? [],
    addProject,
    addTask,
  };
};
