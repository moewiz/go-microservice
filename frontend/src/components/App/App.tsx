import React from "react";
import { SWRConfig } from "swr";
import "./App.css";
import CoffeeList from "../CoffeeList";

function App() {
  return (
    <SWRConfig
      value={{
        fetcher: (resource, init) =>
          fetch(resource, init).then((res) => res.json()),
      }}
    >
      <div className="App">
        <header className="App-header">Coffee Shop</header>
        <CoffeeList />
      </div>
    </SWRConfig>
  );
}

export default App;
