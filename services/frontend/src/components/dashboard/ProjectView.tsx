import { MoreVertical } from "lucide-react";
import { Button } from "../ui/button";
import { Project } from "@/types/projects";

interface ProjectProps {
  onSelectProject: (id: string) => void;
  selectedProjectId: string | null;
  project: Project;
}

export const ProjectView = ({
  onSelectProject,
  selectedProjectId,
  project,
}: ProjectProps) => {
  return (
    <div
      key={project.id}
      className={`p-4 rounded-lg border transition-colors cursor-pointer hover:bg-gray-50 ${
        selectedProjectId === project.id
          ? "border-brand-500 bg-brand-50/50"
          : "border-gray-100"
      }`}
      onClick={() => onSelectProject(project.id)}
    >
      <div className="flex justify-between items-center">
        <div>
          <h3 className="font-medium text-gray-900">{project.name}</h3>
          <p className="text-sm text-gray-500">Client: {project.description}</p>
        </div>
        <div className="flex items-center space-x-2">
          <Button size="sm" variant="ghost" className="h-8 w-8 p-0"></Button>
          <Button size="sm" variant="ghost" className="h-8 w-8 p-0">
            <MoreVertical className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
};
