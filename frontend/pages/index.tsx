import type { NextPage } from "next";
import { CreateRoomForm } from "../components/home/CreateRoomForm";

const Home: NextPage = () => {
  return (
    <div className="flex h-screen justify-center items-center">
      <CreateRoomForm />
    </div>
  );
};

export default Home;
