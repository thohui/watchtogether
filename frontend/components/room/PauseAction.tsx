import { useContext } from "react";
import { WebSocketContext } from "../../context/websocket";
import { useRoomStore } from "../../hooks/store";

export const PauseAction = () => {
  const paused = useRoomStore((state) => state.paused);
  const context = useContext(WebSocketContext);
  const host = useRoomStore((state) => state.host);

  if (!host) return null;

  if (!paused) {
    return (
      <button className="btn mt-2" onClick={context.actions.sendPauseMessage}>
        Pause
      </button>
    );
  }
  return (
    <button className="btn mt-2" onClick={context.actions.sendResumeMessage}>
      Start
    </button>
  );
};
