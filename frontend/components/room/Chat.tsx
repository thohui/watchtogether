import { useContext, useEffect, useRef } from "react";
import { WebSocketContext } from "../../context/websocket";
import { useRoomStore } from "../../hooks/store";

export const Chat = () => {
  return (
    <div className="flex flex-col rounded-lg border border-neutral pt-5">
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
            <span
              className="max-w-xl text-gray-300
            "
              key={index}
            >
              <span
                className={
                  message.owner ? "font-bold text-red-500" : "font-bold"
                }
              >
                {message.sender}:
              </span>{" "}
              {message.message}
            </span>
          );
        })}
        <div ref={ref}></div>
      </div>
    </div>
  );
};

const TextArea = () => {
  const context = useContext(WebSocketContext);
  const submit = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter") {
      e.preventDefault();
      const message = e.target.value;
      if (message.length === 0) return;
      context.actions.sendChatMessage(message);
      e.target.value = "";
      e.target.selectionStart = 0;
    }
  };
  return (
    <textarea
      className="textarea textarea-bordered resize-none mt-3"
      maxLength={100}
      rows={1}
      onKeyDown={submit}
    ></textarea>
  );
};
