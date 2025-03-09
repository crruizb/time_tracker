import { ReactNode } from "react";
import { Navbar } from "../Navbar";

interface DashboardLayoutProps {
  children: ReactNode;
}

export const DashboardLayout = ({ children }: DashboardLayoutProps) => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-[var(--secondary)]">
      <Navbar />
      <div className="container mx-auto px-4 py-8 pt-24">{children}</div>
    </div>
  );
};
