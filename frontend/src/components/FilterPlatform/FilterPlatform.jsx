import React from "react";

const FilterPlatform = ({ onSelect }) => {
  const selectHandler = (e) => {
    const regionName = e.target.value;
    onSelect(regionName);
  };

  return (
    <select onChange={selectHandler}>
      <option className="option">Filter by Platform</option>
      <option className="option" value="Africa">
        -- Put Platforms Here --
      </option>
    </select>
  );
};

export default FilterPlatform;