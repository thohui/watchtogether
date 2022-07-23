import { NextPage } from "next";
import { useRouter } from "next/router";
import { Chat } from "../../components/room/Chat";
import { WebSocketProvider } from "../../context/websocket";
import { useGetRoom } from "../../hooks/room";

const Room: NextPage = () => {
  const router = useRouter();
  const { id } = router.query; // this is undefined initially
  const { exists, loading } = useGetRoom(id as string);
  if (loading) {
    return <p>Loading...</p>;
  }
  if (!exists) {
    return <p>Room does not exist</p>;
  }
  return (
    <WebSocketProvider roomId={id as string}>
      <Chat />
    </WebSocketProvider>
  );
};

export default Room;
