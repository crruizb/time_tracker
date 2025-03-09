import { Task } from "@/types/projects";
import { ListChecks, Check, PlayCircle, PauseCircle } from "lucide-react";
import { Button } from "../ui/button";
import { LiveTimer } from "../LiveTimer";

interface TasksListProps {
  tasks: Task[];
}
export const TasksList = ({ tasks }: TasksListProps) => {
  const formatTime = (seconds: number) => {
    const hrs = Math.floor(seconds / 3600)
      .toString()
      .padStart(2, "0");
    const mins = Math.floor((seconds % 3600) / 60)
      .toString()
      .padStart(2, "0");
    const secs = Math.floor(seconds % 60)
      .toString()
      .padStart(2, "0");
    return `${hrs}:${mins}:${secs}`;
  };

  const getTimeDifference = (start: string | null, end: string | null) => {
    if (start === null || end == null) {
      return "00:00:00";
    }
    const startDate = new Date(start);
    const endDate = new Date(end);
    const diffInSeconds = (endDate.getTime() - startDate.getTime()) / 1000;
    return formatTime(diffInSeconds);
  };

  console.log(tasks);
  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-100 p-6">
      <h2 className="text-lg font-semibold mb-4 flex items-center">
        <ListChecks className="mr-2 h-5 w-5 text-brand-600" />
        Tasks
      </h2>

      <div className="space-y-3">
        {tasks.length === 0 ? (
          <p className="text-gray-500 text-center py-4">
            No tasks yet. Add tasks to this project to get started.
          </p>
        ) : (
          tasks.map((task) => (
            <div
              className={`${
                task.finishedAt !== null
                  ? "border-black-600 bg-gray-50/50"
                  : "border-gray-100"
              }`}
              key={task.id}
            >
              <div
                className={`flex justify-between items-center`}
                key={task.id}
              >
                <span>{task.name}</span>
                <div className="flex items-center space-x-2">
                  <span className="text-sm text-gray-500 font-mono">
                    {task.startedAt === null && "00:00:00"}
                    {task.startedAt !== null &&
                      task.finishedAt !== null &&
                      getTimeDifference(task.startedAt, task.finishedAt)}
                    {task.startedAt !== null && task.finishedAt === null && (
                      <LiveTimer startDate={task.startedAt} />
                    )}
                  </span>
                  <Button size="sm" variant="ghost" className="h-8 w-8 p-0" />
                  {task.startedAt !== null && task.finishedAt !== null ? (
                    <Check className="h-4 w-4" />
                  ) : (
                    <PlayCircle className="h-4 w-4 text-brand-600" />
                  )}
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};
