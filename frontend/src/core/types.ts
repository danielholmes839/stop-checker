export type TravelLocation = {
  id: string;
  title: string;
  description: string;
  position: Position;
};

export type Position = {
  latitude: number;
  longitude: number;
};
