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
      width="100%"
      playing={!paused}
      ref={ref}
      controls={true}
      url={videoId ? `https://youtube.com/watch?v=${videoId}` : undefined}
      onPause={() => {
        if (!paused) {
          ref.current?.getInternalPlayer().playVideo();
        }
      }}
      onPlay={() => {
        if (paused) {
          ref.current?.getInternalPlayer().pauseVideo();
        }
      }}
    />
  );
};
