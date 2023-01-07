import {
  FavouriteIcon,
  FavouriteIconByName,
  TravelLocationInput,
} from "components/travel";
import { Container } from "components/util";
import { FavouriteIconName, FavouriteTravelLocation, useStorage } from "core";
import React, { useState, useEffect } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

const FavouriteIconInput: React.FC<{
  icon: FavouriteIconName;
  activeIcon: FavouriteIconName;
  setActiveIcon: React.Dispatch<FavouriteIconName>;
}> = ({ icon, activeIcon, setActiveIcon }) => {
  return (
    <button
      className={
        icon === activeIcon
          ? "h-10 w-10 bg-primary-100 rounded-full"
          : "h-10 w-10 bg-gray-100 rounded-full hover:bg-gray-200"
      }
      onClick={() => setActiveIcon(icon)}
    >
      <FavouriteIconByName icon={icon} />
    </button>
  );
};

const Favourite: React.FC<{
  favourite: FavouriteTravelLocation;
}> = ({ favourite }) => {
  const nav = useNavigate();
  return (
    <div className="px-3 py-2 bg-gray-50 rounded border-b">
      <div>
        <div className="inline-block align-middle">
          <FavouriteIcon placeId={favourite.id} />
        </div>
        <div
          className="pl-2 inline-block align-middle"
          style={{ maxWidth: "85%" }}
        >
          <span>{favourite.title}</span>
          <p className="text-xs">{favourite.description}</p>
        </div>
      </div>
      <div className="mt-2">
        <button
          className="text-primary-700 bg-primary-100 hover:bg-primary-200 px-2 rounded text-sm mr-2"
          onClick={() => nav(`/travel/p/${favourite.id}`)}
        >
          Directions
        </button>
        <button
          className="text-primary-700 bg-primary-100 hover:bg-primary-200 px-2 rounded text-sm mr-2"
          onClick={() => nav(`/dashboard/favourite/${favourite.id}`)}
        >
          Edit
        </button>
      </div>
    </div>
  );
};

export const DashboardAddFavourite: React.FC = () => {
  const nav = useNavigate();
  const { addFavourite } = useStorage();
  return (
    <Container>
      <h1 className="text-3xl font-bold font mt-3 mb-1">Add Favourite</h1>
      <TravelLocationInput
        suggestFavourites={false}
        setTravelLocation={(location) => {
          addFavourite(location);
          nav(`/dashboard/favourite/${location.id}`);
        }}
      />
    </Container>
  );
};

export const DashboardEditFavourite: React.FC = () => {
  const nav = useNavigate();
  const { id } = useParams();
  const { updateFavourite, getFavourite, deleteFavourite } = useStorage();

  const favourite = getFavourite(id ? id : "");
  const [icon, setIcon] = useState<FavouriteIconName>(
    favourite ? favourite.icon : "saved"
  );
  const [title, setTitle] = useState(favourite ? favourite.title : "");
  const [description, setDescription] = useState(
    favourite ? favourite.description : ""
  );

  useEffect(() => {
    if (favourite === undefined) {
      nav("/dashboard");
    }
  }, [favourite, nav]);

  if (favourite === undefined) {
    return <></>;
  }

  return (
    <Container>
      <h1 className="text-3xl font-bold font mt-3 mb-1">Edit Favourite</h1>
      <div>
        <h2 className="font-semibold">Name</h2>
        <input
          className="mt-1 bg-gray-50 border-b rounded w-full p-3 focus:outline-none focus:border-b focus:border-gray-200 focus:border-0 focus:shadow text-sm"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
      </div>
      <div className="mt-3">
        <h2 className="font-semibold">Description</h2>
        <input
          className="mt-1 bg-gray-50 border-b rounded w-full p-3 focus:outline-none focus:border-b focus:border-gray-200 focus:border-0 focus:shadow text-sm"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />
      </div>
      <div className="mt-3">
        <h2 className="font-semibold">Icon</h2>
        <div className="mt-1">
          <span className="mr-3">
            <FavouriteIconInput
              icon="saved"
              activeIcon={icon}
              setActiveIcon={setIcon}
            />
          </span>
          <span className="mr-3">
            <FavouriteIconInput
              icon="home"
              activeIcon={icon}
              setActiveIcon={setIcon}
            />
          </span>
          <span className="mr-3">
            <FavouriteIconInput
              icon="school"
              activeIcon={icon}
              setActiveIcon={setIcon}
            />
          </span>
          <span className="mr-3">
            <FavouriteIconInput
              icon="office"
              activeIcon={icon}
              setActiveIcon={setIcon}
            />
          </span>
        </div>
        <div className="mt-3 pt-3 border-t">
          <button
            className="text-primary-700 bg-primary-100 hover:bg-primary-200 px-10 py-2 rounded text-sm mr-2"
            onClick={() => {
              updateFavourite({
                ...favourite,
                title: title,
                description: description,
                icon: icon,
              });
              nav("/dashboard");
            }}
          >
            Save
          </button>
          <button
            className="text-red-700 bg-red-100 hover:bg-red-200 px-10 py-2 rounded text-sm mr-2"
            onClick={() => {
              deleteFavourite(favourite.id);
              nav("/dashboard");
            }}
          >
            Delete
          </button>
        </div>
      </div>
    </Container>
  );
};

export const Dashboard: React.FC = () => {
  const { clear, clearHistory, favourites } = useStorage();
  return (
    <Container>
      <h1 className="text-3xl font-bold font mt-3">Dashboard</h1>
      <h2 className="text font-medium text-gray-800 text-xl mt-2">
        Favourites
      </h2>
      <Link
        to="/dashboard/favourite/add"
        className="text-primary-700 text-sm underline"
      >
        New Favourite
      </Link>

      <div className="grid md:grid-cols-2 sm:grid-cols-1 gap-3 mt-2">
        {favourites.map((fav) => (
          <Favourite favourite={fav} />
        ))}
      </div>

      <div className="mt-10">
        <button onClick={clear}>Clear All</button>
        <button onClick={clearHistory}>Clear History</button>
      </div>
    </Container>
  );
};
