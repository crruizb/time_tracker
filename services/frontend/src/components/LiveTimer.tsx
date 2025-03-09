import React, { useEffect, useState } from "react";

const formatTime = (seconds: number) => {
  const hrs = Math.floor(seconds / 3600)
    .toString()
    .padStart(2, "0");
  const mins = Math.floor((seconds % 3600) / 60)
    .toString()
    .padStart(2, "0");
  const secs = Math.floor(seconds % 60)
    .toString()
    .padStart(2, "0");
  return `${hrs}:${mins}:${secs}`;
};

export const LiveTimer = ({ startDate }: { startDate: string }) => {
  const [elapsedTime, setElapsedTime] = useState(0);
  console.log(elapsedTime);

  useEffect(() => {
    const start = new Date(startDate).getTime();

    const interval = setInterval(() => {
      const now = new Date().getTime();
      setElapsedTime((now - start) / 1000); // En segundos
    }, 1000);

    return () => clearInterval(interval);
  }, [startDate]);

  return <span>{formatTime(elapsedTime)}</span>;
};
