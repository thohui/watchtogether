import { useEffect, useRef } from "react";
import ReactPlayer from "react-player";
import { useRoomStore } from "../../hooks/store";

export const Video = () => {
  const videoId = useRoomStore((state) => state.videoId);
  const time = useRoomStore((state) => state.time);
  const paused = useRoomStore((state) => state.paused);
  const ref = useRef<ReactPlayer>(null);
  useEffect(() => {
    if (ref.current) {
      ref.current.seekTo(time);
    }
  }, [time]);

  return (
    <ReactPlayer
      url={videoId ? `https://www.youtube.com/watch?v=${videoId}` : undefined}
      controls={true}
      ref={ref}
      playing={!paused}
    ></ReactPlayer>
  );
};
