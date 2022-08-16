import { useContext, useEffect, useRef } from "react";
import { WebSocketContext } from "../../context/websocket";
import { useRoomStore } from "../../hooks/store";

export const Chat = () => {
  return (
    <div className="flex flex-col rounded-lg w-1/3 h-full border pt-5">
      <Messages />
      <TextArea />
    </div>
  );
};

const Messages = () => {
  const messages = useRoomStore((state) => state.messages);
  const ref = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (ref.current) {
      ref.current?.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);
  return (
    <div className="min-h-16 max-h-28 overflow-auto">
      <div className="flex flex-col">
        {messages.map((message, index) => {
          return (
            <span className="font-black" key={index}>
              {message.sender}: {message.message}
            </span>
          );
        })}
        <div ref={ref}></div>
      </div>
    </div>
  );
};

// TODO: input validation
const TextArea = () => {
  const context = useContext(WebSocketContext);
  const submit = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter") {
      const message = e.target.value;
      context.actions.sendChatMessage(message);
      e.target.value = "";
    }
  };
  return (
    <textarea
      className="textarea textarea-bordered resize-none whitespace-nowrap overflow-x-scroll mt-5"
      rows={1}
      maxLength={100}
      onKeyDown={submit}
    ></textarea>
  );
};
