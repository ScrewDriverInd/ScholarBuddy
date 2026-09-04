const API_BASE = (import.meta.env.VITE_API_BASE || "").replace(/\/+$/, "");

export const apiURL = (path) => `${API_BASE}${path}`;

export const apiFetch = async (path, options = {}) => {
  const response = await fetch(apiURL(path), {
    ...options,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error?.error?.message || `HTTP ${response.status}`);
  }

  if (response.status === 204) {
    return null;
  }

  return response.json();
};

export const unwrapData = (payload) => payload?.data;
