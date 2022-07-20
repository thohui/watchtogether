import create from "zustand";
import { ChatMessage } from "../types/room";

interface Store {
  status: "connected" | "disconnected";
  messages: ChatMessage[];
  videoURL: string | null;
  time: number;
  actions: {
    appendMessage: (message: ChatMessage) => void;
    setStatus: (status: "connected" | "disconnected") => void;
    setVideoURL: (videoURL: string | null) => void;
    setTime: (time: number) => void;
  };
}

export const useRoomStore = create<Store>((set) => ({
  status: "disconnected",
  messages: [],
  videoURL: null,
  time: 0,
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
    setVideoURL: (videoURL: string | null) => {
      set((state) => ({
        ...state,
        videoURL,
      }));
    },
    setTime: (time: number) => {
      set((state) => ({
        ...state,
        time,
      }));
    },
  },
}));
