import { useEffect, useState } from "react";

interface FetchArgs {
  url: string | null;
  body?: any;
  method: string;
}

export const useFetch = ({ url, body, method }: FetchArgs) => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<boolean>(false);
  const [data, setData] = useState<any>();
  useEffect(() => {
    async function fetchData() {
      if (typeof window === "undefined" || !url) {
        return;
      }
      try {
        const response = await fetch(url, {
          method,
          body,
          headers: {
            "Content-Type": "application/json",
          },
        });
        const json = await response.json();
        setData(json);
        setLoading(false);
      } catch {
        setError(true);
        setLoading(false);
      }
    }
    fetchData();
  }, [url, body, method]);
  return { loading, error, data };
};
