import { useEffect, useState } from "react";
import ReactPlayer from "react-player";
import { useWindowSize } from "../../hooks/window";
interface Props {
  id: string | undefined;
}
export const VideoPreview = ({ id }: Props) => {
  const [loading, setLoading] = useState(true);
  const size = useWindowSize();
  useEffect(() => {
    setLoading(false);
  }, []);

  if (loading) {
    return null;
  }
  return (
    <ReactPlayer
      url={`https://youtube.com/watch?v=${id}`}
      width={size.width < 640 ? size.width : 640}
    />
  );
};
