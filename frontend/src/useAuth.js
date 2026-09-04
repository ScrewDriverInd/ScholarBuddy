import { useState, useEffect } from "react";
import { apiFetch, unwrapData } from "./api";

export const useAuth = () => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const response = await apiFetch("/api/v1/user/me");
      setUser(unwrapData(response));
    } catch {
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  const login = () => {
    window.location.href = "/oauth2/authorization/google";
  };

  const logout = async () => {
    try {
      await fetch("/logout", { method: "POST", credentials: "include" });
      setUser(null);
    } catch (error) {
      console.error("Logout failed:", error);
    }
  };

  const isAdmin = user?.roles?.includes("ROLE_ADMIN") || false;

  return { user, loading, login, logout, isAdmin, checkAuth };
};
