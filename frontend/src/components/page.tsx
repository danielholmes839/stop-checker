import React from "react";

export const Container: React.FC = ({ children }) => {
  return <div className="container mx-auto px-10 lg:w-2/3 xl:w-1/2">{children}</div>;
};

export const Card: React.FC<React.HTMLAttributes<any>> = ({
  children,
  ...rest
}) => {
  return (
    <div
      {...rest}
      className="px-5 py-3 bg-white rounded-sm border border-gray-200"
    >
      {children}
    </div>
  );
};
