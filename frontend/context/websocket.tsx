import { createContext, useEffect, useRef } from "react";
import { useRoomStore } from "../hooks/store";
import { ChatMessage, UnknownMessage } from "../types/room";

interface ContextProps {
  websocket: WebSocket | null;
  actions: {
    sendChatMessage: (message: string) => void;
  };
}

export const WebSocketContext = createContext<ContextProps>({
  websocket: null,
  actions: { sendChatMessage: () => {} },
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
    socket.onopen = () => {
      actions.setStatus("connected");
    };
    socket.onclose = () => {
      actions.setStatus("disconnected");
    };
    socket.onmessage = (event: MessageEvent) => {
      const message: UnknownMessage = JSON.parse(event.data);
      if (message.type === "chat") {
        actions.appendMessage({
          sender: message.data.sender,
          message: message.data.message,
        } as ChatMessage);
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
    },
  };

  return (
    <WebSocketContext.Provider value={properties}>
      {children}
    </WebSocketContext.Provider>
  );
};
