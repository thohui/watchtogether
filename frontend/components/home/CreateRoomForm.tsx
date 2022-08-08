import { useRouter } from "next/router";
import { useEffect, useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { URL } from "../../utils/constants";
import { VideoPreview } from "./VideoPreview";

export const CreateRoomForm = () => {
  const [videoId, setVideoId] = useState<string | undefined>(undefined);
  const { handleSubmit, register, watch } = useForm();
  const router = useRouter();

  const onSubmit = async () => {
    const response = await fetch(`${URL}/room/create`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        video_id: videoId,
      }),
    });
    if (response.ok) {
      const data = await response.json();
      data.id && router.push(`/room/${data.id}`);
    }
  };
  const url = watch("video_url");

  const regex = useMemo(
    () =>
      /^(?:https?:\/\/)?(?:www\.)?(?:youtu\.be\/|youtube\.com\/(?:embed\/|v\/|watch\?v=|watch\?.+&v=))((\w|-){11})(?:\S+)?$/,
    []
  );

  useEffect(() => {
    if (url) {
      const match = url.match(regex);
      if (match) {
        setVideoId(match[1]);
      } else {
        setVideoId(undefined);
      }
    }
  }, [url, regex]);

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="flex flex-col items-center py-5">
        <label className="text-2xl label" htmlFor="url">
          YouTube URL
        </label>
        <input
          type="text"
          className="input input-bordered w-full lg:max-w-lg max-w-sm"
          {...register("video_url", {
            required: true,
            pattern: {
              value: regex,
              message: "Invalid YouTube URL",
            },
          })}
        />
      </div>

      <div className="flex flex-col items-center py-5">
        <VideoPreview id={videoId} />
        <button type="submit" className="btn my-5 w-1/2" disabled={!videoId}>
          Submit
        </button>
      </div>
    </form>
  );
};
