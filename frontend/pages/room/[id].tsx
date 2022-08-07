import { NextPage } from "next";
import { useRouter } from "next/router";
import { Chat } from "../../components/room/Chat";
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
      <Chat />
    </WebSocketProvider>
  );
};

export default Room;
