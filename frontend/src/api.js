const API_BASE = (import.meta.env.VITE_API_BASE || "").replace(/\/+$/, "");

export const apiURL = (path) => `${API_BASE}${path}`;

const readCookie = (name) => {
  const prefix = `${name}=`;
  const cookie = document.cookie
    .split(";")
    .map((part) => part.trim())
    .find((part) => part.startsWith(prefix));
  return cookie ? decodeURIComponent(cookie.slice(prefix.length)) : null;
};

export const apiFetch = async (path, options = {}) => {
  const body =
    options.body && typeof options.body !== "string"
      ? JSON.stringify(options.body)
      : options.body;

  const response = await fetch(apiURL(path), {
    ...options,
    body,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(options.method && !["GET", "HEAD", "OPTIONS"].includes(options.method.toUpperCase())
        ? { "X-XSRF-TOKEN": readCookie("XSRF-TOKEN") || "" }
        : {}),
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
