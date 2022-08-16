import { RouteType } from "client/types";

export const formatTime = (datetime: string): string => {
  let date = new Date(datetime);

  let hours = date.getHours();
  let minutes = date.getMinutes();
  let meridean = hours < 12 ? "AM" : "PM";

  // humans: 0 should be 12
  hours = hours % 12;
  if (hours === 0) {
    hours = 12;
  }

  // padding minutes
  let minutesPadded = minutes < 10 ? `0${minutes}` : minutes;
  return `${hours}:${minutesPadded} ${meridean}`;
};

export const formatDistance = (distance: number): string => {
  if (distance >= 1000) {
    return (distance / 1000).toFixed(1) + " km";
  }

  return (Math.round(distance / 10) * 10).toFixed(0) + " meters";
};

export const formatRouteType = (rt: RouteType): string => {
  return rt.toLowerCase();
};
