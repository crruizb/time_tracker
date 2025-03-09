import { Briefcase } from "lucide-react";
import { Project } from "@/types/projects";
import { ProjectView } from "./ProjectView";

interface ProjectsListProps {
  onSelectProject: (id: string) => void;
  selectedProjectId: string | null;
  projects: Project[];
  isLoading: boolean;
}

export const ProjectsList = ({
  onSelectProject,
  selectedProjectId,
  projects,
  isLoading,
}: ProjectsListProps) => {
  if (isLoading) {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-100 p-6 mb-6">
        <h2 className="text-lg font-semibold mb-4 flex items-center">
          <Briefcase className="mr-2 h-5 w-5 text-brand-600" />
          Projects
        </h2>
        <div className="animate-pulse space-y-4">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="h-16 bg-gray-100 rounded-md"></div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-100 p-6 mb-6">
      <h2 className="text-lg font-semibold mb-4 flex items-center">
        <Briefcase className="mr-2 h-5 w-5 text-brand-600" />
        Projects
      </h2>

      <div className="space-y-4">
        {projects.length === 0 ? (
          <p className="text-gray-500 text-center py-4">
            No projects yet. Create your first project to get started.
          </p>
        ) : (
          projects.map((project) => (
            <ProjectView
              onSelectProject={onSelectProject}
              selectedProjectId={selectedProjectId}
              project={project}
            />
          ))
        )}
      </div>

      <div></div>
    </div>
  );
};
