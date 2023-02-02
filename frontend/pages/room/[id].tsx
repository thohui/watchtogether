import { NextPage } from "next";
import { useRouter } from "next/router";
import { Chat } from "../../components/room/Chat";
import { PauseAction } from "../../components/room/PauseAction";
import { Video } from "../../components/room/Video";
import { WebSocketProvider } from "../../context/websocket";
import { useFetch } from "../../hooks/fetch";
import { URL } from "../../utils/constants";

const Room: NextPage = () => {
  const router = useRouter();
  const { id } = router.query;
  const url = id ? `${URL}/room/${id}` : null;

  const { loading, error, data } = useFetch({
    url,
    body: JSON.stringify({ id: id }),
    method: "POST",
  });
  if (loading) {
    return <p>Loading...</p>;
  }
  if (!data || error) {
    return <p>Room does not exist</p>;
  }

  return (
    <WebSocketProvider roomId={id as string}>
      <div className="flex h-screen justify-center items-center">
        <div className="flex flex-col space-y-3">
          <Video />
          <PauseAction />
          <Chat />
        </div>
      </div>
    </WebSocketProvider>
  );
};

export default Room;
