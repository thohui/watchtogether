import { createContext, useEffect, useRef } from "react";
import { useRoomStore } from "../hooks/store";
import {
  ChatMessage,
  InitMessage,
  UnknownMessage,
  VideoUpdateMessage,
} from "../types/room";

interface ContextProps {
  websocket: WebSocket | null;
  actions: {
    sendChatMessage: (message: string) => void;
    sendPauseMessage: () => void;
    sendResumeMessage: () => void;
  };
}

export const WebSocketContext = createContext<ContextProps>({
  websocket: null,
  actions: {
    sendChatMessage: () => {
      throw new Error("WebSocketContext is not initialized");
    },
    sendPauseMessage: () => {
      throw new Error("WebSocketContext is not initialized");
    },
    sendResumeMessage: () => {
      throw new Error("WebSocketContext is not initialized");
    },
  },
});

interface Props {
  roomId: string;
  children: JSX.Element | JSX.Element[];
}

export const WebSocketProvider = ({ roomId, children }: Props) => {
  const ws = useRef<WebSocket | null>(null);
  const WEBSOCKET_URL =
    process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost/ws";
  const url = `${WEBSOCKET_URL}/${roomId}`;
  const actions = useRoomStore((state) => state.actions);
  useEffect(() => {
    const socket = new WebSocket(url);
    socket.onopen = () => {};
    socket.onclose = () => {};
    socket.onmessage = (event: MessageEvent) => {
      const unknownMessage: UnknownMessage = JSON.parse(event.data);
      switch (unknownMessage.type) {
        case "chat": {
          const message: ChatMessage = unknownMessage.data;
          actions.appendMessage(message);
          break;
        }
        case "init": {
          const message: InitMessage = unknownMessage.data;
          actions.setVideoId(message.video_id);
          actions.setTime(message.time);
          actions.setHost(message.host);
          actions.setPaused(message.paused);
          break;
        }
        case "video_update":
          const message: VideoUpdateMessage = unknownMessage.data;
          actions.setTime(message.time);
          actions.setPaused(message.paused);
          break;
      }
    };
    ws.current = socket;
    return () => {
      socket.close();
    };
  });

  const properties: ContextProps = {
    websocket: ws.current,
    actions: {
      sendChatMessage: (message: string) => {
        if (ws.current) {
          const payload = JSON.stringify({ type: "chat", message: message });
          ws.current.send(payload);
        }
      },
      sendPauseMessage: () => {
        if (ws.current) {
          const payload = JSON.stringify({ type: "pause" });
          ws.current.send(payload);
        }
      },
      sendResumeMessage: () => {
        if (ws.current) {
          const payload = JSON.stringify({ type: "resume" });
          ws.current.send(payload);
        }
      },
    },
  };

  return (
    <WebSocketContext.Provider value={properties}>
      {children}
    </WebSocketContext.Provider>
  );
};
