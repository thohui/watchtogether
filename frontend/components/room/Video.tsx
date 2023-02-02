import { useEffect, useRef } from "react";
import ReactPlayer from "react-player";

import { useRoomStore } from "../../hooks/store";
import { useWindowSize } from "../../hooks/window";

export const Video = () => {
  const videoId = useRoomStore((state) => state.videoId);
  const time = useRoomStore((state) => state.time);
  const paused = useRoomStore((state) => state.paused);
  const ref = useRef<ReactPlayer>(null);
  const size = useWindowSize();
  useEffect(() => {
    if (ref.current) {
      // only seek if there's a desync
      if (time - ref.current.getCurrentTime() > 1) {
        ref.current.seekTo(time);
      }
    }
  }, [time]);
  return (
    <ReactPlayer
      playing={!paused}
      ref={ref}
      width={size.width < 640 ? size.width - 20 : 640}
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
