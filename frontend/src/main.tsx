import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { Provider as UrqlProvider } from "urql";
import { App } from "./App";
import { urqlClient } from "./urqlClient";
import "./index.css";

const rootEl = document.getElementById("root");
if (!rootEl) {
  throw new Error("root element not found");
}

createRoot(rootEl).render(
  <StrictMode>
    <UrqlProvider value={urqlClient}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </UrqlProvider>
  </StrictMode>
);
