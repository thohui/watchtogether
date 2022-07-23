import { useEffect, useState } from "react";

const URL = process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost";

const getRoomRequest = async (id: string) => {
  // ensure we are only fetching on the client
  if (typeof window === "undefined") {
    return false;
  }
  if (!id) {
    return false;
  }
  const response = await fetch(`${URL}/room/${id}`, {
    method: "POST",
    body: JSON.stringify({ id: id }),
    headers: {
      "Content-Type": "application/json",
    },
  });
  return response.status === 200;
};

export const useGetRoom = (id: string) => {
  const [loading, setLoading] = useState(true);
  const [exists, setExists] = useState(false);
  useEffect(() => {
    async function getRoom() {
      const response = await getRoomRequest(id);
      setExists(response);
      setLoading(false);
    }
    getRoom();
  }, [id]);
  return { loading, exists };
};
