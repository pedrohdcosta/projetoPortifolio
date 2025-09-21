import { defineStore } from "pinia";
import api from "../api/axios";

export const useAuth = defineStore("auth", {
  state: () => ({
    token: localStorage.getItem("token") || "",
    user: null as null | { id: number; name: string; email: string },
  }),
  getters: { isAuthenticated: (s) => !!s.token },
  actions: {
    async login(email: string, password: string) {
      const { data } = await api.post("/auth/login", { email, password });
      this.token = data.accessToken;
      this.user = data.user;
      localStorage.setItem("token", this.token);
    },
    async signup(name: string, email: string, password: string) {
      await api.post("/auth/signup", { name, email, password });
    },
    async fetchMe() {
      const { data } = await api.get("/auth/me");
      this.user = data;
    },
    logout() {
      this.token = "";
      this.user = null;
      localStorage.removeItem("token");
    },
  },
});
