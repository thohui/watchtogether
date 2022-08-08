import { useEffect, useState } from "react";
import ReactPlayer from "react-player";
interface Props {
  id: string | undefined;
}
export const VideoPreview = ({ id }: Props) => {
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    setLoading(false);
  }, []);

  if (loading) {
    return null;
  }
  return (
    <ReactPlayer
      url={`https://youtube.com/watch?v=${id}`}
      width="100%"
    ></ReactPlayer>
  );
};
