import { useEffect, useLayoutEffect, useRef } from "react";
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

  useLayoutEffect(() => {
    if (paused) {
      document.querySelector("video")?.pause();
    } else {
      document.querySelector("video")?.play();
    }
  }, [paused]);

  return (
    <ReactPlayer
      ref={ref}
      playing={!paused}
      id="video"
      width={size.width < 640 ? size.width - 20 : 640}
      controls={true}
      url={videoId ? `https://youtube.com/watch?v=${videoId}` : undefined}
    />
  );
};
