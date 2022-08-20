import type { NextPage } from "next";
import { CreateRoomForm } from "../components/home/CreateRoomForm";
import { Navbar } from "../components/navigation/Navbar";

const Home: NextPage = () => {
  return (
    <div className="container mx-auto">
      <Navbar />
      <div className="grid place-items-center">
        <h1 className="text-5xl py-3">Create Room</h1>
        <div className="w-1/2">
          <CreateRoomForm />
        </div>
      </div>
    </div>
  );
};

export default Home;
