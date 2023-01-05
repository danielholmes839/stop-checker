import React, { useState, useEffect } from "react";
import { TravelLocation } from "./types";

type Favourite = {
  icon: "home" | "office" | "school" | undefined;
};

export type FavouriteTravelLocation = TravelLocation & Favourite;

type StorageValue = {
  isFavourite: (id: string) => boolean;
  isHistory: (id: string) => boolean;
  history: TravelLocation[];
  favourites: FavouriteTravelLocation[];
  addHistory: (location: TravelLocation) => void;
  addFavourite: (location: TravelLocation) => void;
  updateFavourite: (favourite: FavouriteTravelLocation) => void;
  deleteFavourite: (id: string) => void;
  clearHistory: () => void;
  clear: () => void;
};

const StorageContext = React.createContext<StorageValue>({
  addFavourite: (location) => {},
  addHistory: (location) => {},
  deleteFavourite: (id) => {},
  favourites: [],
  history: [],
  isFavourite: (id) => false,
  isHistory: (id) => false,
  updateFavourite: (favourite) => {},
  clearHistory: () => {},
  clear: () => {},
});

const read = (key: string, placeholder: any): any => {
  const saved = localStorage.getItem(key);
  if (saved === null) {
    return placeholder;
  }
  const initial = JSON.parse(saved);
  if (initial) {
    return initial;
  }
  return placeholder;
};

const save = (key: string, value: any) => {
  localStorage.setItem(key, JSON.stringify(value));
};

export const StorageProvider: React.FC = ({ children }) => {
  const [history, setHistory] = useState<{ [key: string]: TravelLocation }>(
    () => read("location-history", {})
  );

  const [favourites, setFavourites] = useState<{
    [key: string]: FavouriteTravelLocation;
  }>(() => read("location-favourites", {}));

  useEffect(() => {
    save("location-history", history);
  }, [history]);

  useEffect(() => {
    save("location-favourites", favourites);
  }, [favourites]);

  const clear = () => {
    setFavourites({});
    setHistory({});
    console.log("clear");
  };

  const clearHistory = () => {
    setHistory({});
  };

  const deleteHistory = (id: string) => {
    let historyCopy = Object.assign({}, history);
    delete historyCopy[id];
    setHistory(historyCopy);
  };

  const addHistory = (location: TravelLocation) => {
    // add to history when the location is not a favourite
    if (
      favourites[location.id] === undefined &&
      history[location.id] === undefined
    ) {
      setHistory(Object.assign({}, { [location.id]: location }, history));
    }
  };

  const addFavourite = (location: TravelLocation) => {
    let favourite: FavouriteTravelLocation = {
      ...location,
      icon: undefined,
    };
    deleteHistory(favourite.id); // we don't want to show a place in history and favourite
    setFavourites(Object.assign({}, { [favourite.id]: favourite }, favourites));
  };

  // const getFavourite = (place: string): FavouriteTravelLocation => {

  // }

  const updateFavourite = (favourite: FavouriteTravelLocation) => {
    setFavourites(Object.assign({}, favourites, { [favourite.id]: favourite }));
  };

  const deleteFavourite = (id: string) => {
    let favouritesCopy = Object.assign({}, favourites);
    delete favouritesCopy[id];
    setFavourites(favouritesCopy);
  };

  return (
    <StorageContext.Provider
      value={{
        addFavourite: addFavourite,
        addHistory: addHistory,
        deleteFavourite: deleteFavourite,
        favourites: Object.values(favourites),
        history: Object.values(history),
        isFavourite: (id) => favourites[id] !== undefined,
        isHistory: (id) => history[id] !== undefined,
        updateFavourite: updateFavourite,
        clearHistory: clearHistory,
        clear: clear,
      }}
    >
      {children}
    </StorageContext.Provider>
  );
};

export const useStorage = () => {
  return React.useContext(StorageContext);
};
