import React from "react";
import { SWRConfig } from "swr";
import "./App.css";
import CoffeeList from "../CoffeeList";
import UploadMultipart from "../UploadMultipart";

function App() {
  return (
    <SWRConfig
      value={{
        fetcher: async (...args) => {
          const res = await fetch(args);
          return res.json();
        },
      }}
    >
      <div className="App">
        <header className="App-header">Coffee Shop</header>
        <br />
        <hr />
        <CoffeeList />
        <br />
        <hr />
        <UploadMultipart />
      </div>
    </SWRConfig>
  );
}

export default App;
