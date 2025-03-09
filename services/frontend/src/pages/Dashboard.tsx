import { useProjects } from "@/hooks/useProjects";
import { Navigate, useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import { ProjectsList } from "@/components/dashboard/ProjectsList";
import { useState } from "react";
import { DashboardLayout } from "@/components/dashboard/DashboardLayout";
import { Button } from "@/components/ui/button";
import { ArrowLeft, Plus } from "lucide-react";
import { TasksList } from "@/components/dashboard/TasksList";
import { CreateProjectDialog } from "@/components/dashboard/CreateProjectDialog";
import { CreateTaskDialog } from "@/components/dashboard/CreateTaskDialog";

const Dashboard: React.FC = () => {
  const username = Cookies.get("username");
  if (username === undefined) {
    return <Navigate to="/" />;
  }

  const navigate = useNavigate();
  const [isProjectDialogOpen, setIsProjectDialogOpen] = useState(false);
  const [selectedProjectId, setSelectedProjectId] = useState<string | null>(
    null
  );
  const [isTaskDialogOpen, setIsTaskDialogOpen] = useState(false);

  const { isLoading, isError, error, projects } = useProjects();
  if (isLoading) return <div>Loading...</div>;
  if (isError) return <div>{error?.message}</div>;

  return (
    <DashboardLayout>
      <div className="flex justify-between items-center mb-8">
        <div className="flex items-center gap-2">
          <h1 className="text-2xl font-bold">Dashboard</h1>
        </div>
        <div className="flex gap-2">
          <Button onClick={() => setIsProjectDialogOpen(true)} size="sm">
            <Plus className="h-4 w-4 mr-1" />
            New Project
          </Button>
          <Button
            onClick={() => setIsTaskDialogOpen(true)}
            size="sm"
            variant="outline"
            disabled={!selectedProjectId}
          >
            <Plus className="h-4 w-4 mr-1" />
            New Task
          </Button>
        </div>
      </div>

      <div>
        <ProjectsList
          onSelectProject={(id) => setSelectedProjectId(id)}
          selectedProjectId={selectedProjectId}
          projects={projects}
          isLoading={isLoading}
        />
      </div>

      {selectedProjectId && (
        <TasksList
          tasks={projects.filter((p) => p.id === selectedProjectId)[0].tasks}
        />
      )}

      <CreateProjectDialog
        isOpen={isProjectDialogOpen}
        onClose={() => setIsProjectDialogOpen(false)}
      />

      <CreateTaskDialog
        isOpen={isTaskDialogOpen}
        onClose={() => setIsTaskDialogOpen(false)}
        projectId={selectedProjectId}
      />
    </DashboardLayout>
  );
};

export default Dashboard;
