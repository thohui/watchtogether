import create from "zustand";
import { ChatMessage } from "../types/room";

interface Store {
  status: "connected" | "disconnected";
  messages: ChatMessage[];
  videoId: string | null;
  time: number;
  host: boolean;
  paused: boolean;
  actions: {
    appendMessage: (message: ChatMessage) => void;
    setStatus: (status: "connected" | "disconnected") => void;
    setVideoId: (videoId: string | null) => void;
    setTime: (time: number) => void;
    setHost: (host: boolean) => void;
    setPaused: (paused: boolean) => void;
  };
}

export const useRoomStore = create<Store>((set) => ({
  status: "disconnected",
  messages: [],
  videoId: null,
  time: 0,
  host: false,
  paused: false,
  actions: {
    appendMessage: (message: ChatMessage) => {
      set((state) => ({
        ...state,
        messages: [...state.messages, message],
      }));
    },
    setStatus: (status: "connected" | "disconnected") => {
      set((state) => ({
        ...state,
        status,
      }));
    },
    setVideoId: (videoId: string | null) => {
      set((state) => ({
        ...state,
        videoId,
      }));
    },
    setTime: (time: number) => {
      set((state) => ({
        ...state,
        time,
      }));
    },
    setHost: (host: boolean) => {
      set((state) => ({
        ...state,
        host,
      }));
    },
    setPaused(paused: boolean) {
      set((state) => ({
        ...state,
        paused,
      }));
    },
  },
}));
