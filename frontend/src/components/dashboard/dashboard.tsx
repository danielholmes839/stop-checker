import { Container } from "components/util";
import { useStorage } from "core";
import React from "react";
import { useNavigate } from "react-router-dom";
import { Favourite } from "./favourite";
import { History } from "./history";

export const Dashboard: React.FC = () => {
  const nav = useNavigate();
  const { clearHistory, favourites, history } = useStorage();
  return (
    <Container>
      <h1 className="text-3xl font-bold font mt-3">Dashboard</h1>
      <h2 className="text font-medium text-gray-800 text-xl mt-2">
        Favourites
      </h2>

      <div className="grid md:grid-cols-2 sm:grid-cols-1 gap-3 mt-2">
        {favourites.map((fav) => (
          <Favourite favourite={fav} key={fav.id} />
        ))}{" "}
        <button
          className="flex items-center justify-center bg-gray-50 hover:bg-gray-100 py-5 rounded border-b"
          onClick={() => nav("/dashboard/favourite/add")}
        >
          <span className="font-medium text-gray-800 text-sm">
            New Favourite
          </span>
        </button>
      </div>
      {history.length > 0 && (
        <div className="mt-5">
          <h2 className="text font-medium text-gray-800 text-xl">History</h2>
          <button
            className="text-red-500 hover:text-red-700 rounded text-sm mr-2"
            onClick={clearHistory}
          >
            Delete
          </button>
          <div className="mt-2">
            {history.map((recent) => (
              <div className="mt-2">
                <History recent={recent} key={recent.id} />
              </div>
            ))}
          </div>

          <div className="mb-10"></div>
        </div>
      )}
    </Container>
  );
};
