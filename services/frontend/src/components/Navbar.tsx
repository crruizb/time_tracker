import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { Clock, Menu, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import ThemeToggle from "./ThemeToggle";

export const Navbar = () => {
  const [isScrolled, setIsScrolled] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 10);
    };

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  return (
    <header
      className={cn(
        "fixed top-0 left-0 right-0 z-50 transition-all duration-300 ease-in-out py-4 bg-blue-400",
        isScrolled
          ? "bg-blue-400 shadow-sm py-3"
          : "bg-blue-400 border-b border-gray-100"
      )}
    >
      <div className="container mx-auto px-4 md:px-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <Link to="/" className="flex items-center space-x-2">
              <Clock className="h-6 w-6 text-brand-600" />
              <span className="text-xl font-semibold">TimeHaven</span>
            </Link>
          </div>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center space-x-8">
            <Link
              to="/dashboard"
              className="text-sm font-medium text-brand-600 transition-colors"
            >
              Dashboard
            </Link>
            <Button size="sm" variant="ghost">
              Log out
            </Button>
            <ThemeToggle />
          </nav>

          {/* Mobile menu button */}
          <button
            className="md:hidden flex items-center"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            aria-label="Toggle menu"
          >
            {isMobileMenuOpen ? (
              <X className="h-6 w-6" />
            ) : (
              <Menu className="h-6 w-6" />
            )}
          </button>
        </div>
      </div>

      {/* Mobile Navigation */}
      {isMobileMenuOpen && (
        <div className="md:hidden absolute top-full left-0 right-0 bg-white shadow-md animate-fade-in">
          <div className="container mx-auto px-4 py-4 flex flex-col space-y-4">
            <Link
              to="/dashboard"
              className="text-sm font-medium py-2 text-brand-600 transition-colors"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Dashboard
            </Link>
            <hr className="border-gray-100" />
            <Button
              variant="ghost"
              className="justify-start pl-0"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Log out
            </Button>
          </div>
        </div>
      )}
    </header>
  );
};
