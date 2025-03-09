export interface Task {
  id: string;
  name: string;
  description: string;
  username: string | null;
  startedAt: string | null;
  finishedAt: string | null;
}

export interface Project {
  id: string;
  name: string;
  description: string;
  tasks: Task[];
}

export interface Projects {
  report: Project[];
}
