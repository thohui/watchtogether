import { NextPage } from "next";
import { useRouter } from "next/router";
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
  return <h1>You entered a room</h1>;
};

export default Room;
