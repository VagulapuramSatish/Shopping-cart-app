import { useState } from "react";
import { Routes, Route, Navigate, useLocation } from "react-router-dom";
import Login from "./components/Login";
import ItemsList from "./components/ItemsList";
import NavBar from "./components/NavBar";

const App = () => {
  const [token, setToken] = useState(localStorage.getItem("token") || "");
  const [action, setAction] = useState(null);
  const location = useLocation();

  const showNavBar = token && location.pathname !== "/";

  return (
    <>
      {showNavBar && (
        <NavBar
          onCart={() => setAction("cart")}
          onCheckout={() => setAction("checkout")}
          onOrders={() => setAction("orders")}
        />
      )}

      <Routes>
        <Route path="/" element={<Login setToken={setToken} />} />

        <Route
          path="/items"
          element={
            token ? (
              <ItemsList
                token={token}
                action={action}
                clearAction={() => setAction(null)}
              />
            ) : (
              <Navigate to="/" />
            )
          }
        />
      </Routes>
    </>
  );
};

export default App;
