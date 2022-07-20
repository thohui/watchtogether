import { useRoomStore } from "../../hooks/store";

export default function Chat() {
  let messages = useRoomStore((state) => state.messages);
  return (
    <div>
      {messages.map((message, index) => {
        return (
          <p key={index}>
            {message.sender} : {message.message}{" "}
          </p>
        );
      })}
    </div>
  );
}
